package promise

import (
	"errors"
	"time"
)

const (
	PromiseInit = int64(0)
)

var (
	ErrResolved       = errors.New("resolved already")
	ErrTimeoutNoSet   = errors.New("timeout not set")
	ErrTimeout        = errors.New("timeout")
	ErrNotImplemented = errors.New("not implemented")
)

type Promise interface {
	// Reset Reset promise
	Reset()

	// ResetWithOptions Reset promise will options
	ResetWithOptions(interface{})

	// Close Close the promise
	Close()

	// IsResolved If the promise is resolved
	IsResolved() bool

	// Get the time the promise last resolved. time.Time{} if the promise is unresolved.
	ResolvedAt() time.Time

	// Resolve Resolve the promise with value or (value, error)
	Resolve(...interface{}) (Promise, error)

	// Options Get options
	Options() interface{}

	// Value Get resolved value
	Value() interface{}

	// Result Helper function to get (value, error)
	Result() (interface{}, error)

	// Error Get last error on resolving
	Error() error

	// SetTimeout Set how long the promise should timeout.
	SetTimeout(time.Duration)

	// Deadline returns the deadline if timeout is set.
	Deadline() (time.Time, bool)

	// Timeout returns ErrTimeout if timeout, or ErrTimeoutNoSet if the timeout not set.
	Timeout(timeout ...time.Duration) error

	// TimeoutC returns a channel that is closed when the promise timeout.
	TimeoutC(timeout ...time.Duration) (<-chan time.Time, error)
}

func Resolved(rets ...interface{}) Promise {
	return ResolvedChannel(rets...)
}

func NewPromise() Promise {
	return NewChannelPromiseWithOptions(nil)
}

func NewPromiseWithOptions(opts interface{}) Promise {
	return NewChannelPromiseWithOptions(opts)
}
