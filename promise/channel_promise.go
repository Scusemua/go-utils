package promise

import (
	"sync"
	"sync/atomic"
	"time"
)

type ChannelPromise struct {
	AbstractPromise

	cond chan struct{}
	mu   sync.Mutex
}

func ResolvedChannel(rets ...interface{}) *ChannelPromise {
	promise := NewChannelPromiseWithOptions(nil)
	if promise.ResolveRets(rets...) {
		promise.resolved = time.Now().UnixNano()
	} else {
		promise.resolved = int64(1) // Differentiate with PromiseInit
	}
	close(promise.cond)
	return promise
}

func NewChannelPromise() *ChannelPromise {
	return NewChannelPromiseWithOptions(nil)
}

func NewChannelPromiseWithOptions(opts interface{}) *ChannelPromise {
	promise := &ChannelPromise{}
	promise.ResetWithOptions(opts)
	promise.SetProvider(promise)
	return promise
}

func (p *ChannelPromise) Reset() {
	p.ResetWithOptions(nil)
}

func (p *ChannelPromise) ResetWithOptions(opts interface{}) {
	p.AbstractPromise.ResetWithOptions(opts)
	p.cond = make(chan struct{})
}

func (p *ChannelPromise) Resolve(rets ...interface{}) (Promise, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	select {
	case <-p.cond:
		return p, ErrResolved
	default:
		p.AbstractPromise.ResolveRets(rets...)
		atomic.StoreInt64(&p.resolved, time.Now().UnixNano())
		close(p.cond)
	}
	return p, nil
}

func (p *ChannelPromise) Timeout() error {
	ch, err := p.TimeoutC()
	if err == ErrResolved {
		return nil
	} else if err != nil {
		return err
	}

	select {
	case <-ch:
		return ErrTimeout
	case <-p.cond:
		return nil
	}
}

// PromiseProvider
func (p *ChannelPromise) Wait() {
	<-p.cond
}

func (p *ChannelPromise) Lock() {
	p.mu.Lock()
}

func (p *ChannelPromise) Unlock() {
	p.mu.Unlock()
}
