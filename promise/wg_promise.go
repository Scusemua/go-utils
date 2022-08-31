package promise

import (
	"runtime"
	"time"

	"github.com/mason-leap-lab/go-utils/sync"
)

type WaitGroupPromise struct {
	AbstractPromise

	wg sync.WaitGroup
}

func ResolvedWaitGroup(rets ...interface{}) *WaitGroupPromise {
	promise := NewWaitGroupPromiseWithOptions(nil)
	promise.Resolve(rets...)
	return promise
}

func NewWaitGroupPromise() *WaitGroupPromise {
	return NewWaitGroupPromiseWithOptions(nil)
}

func NewWaitGroupPromiseWithOptions(opts interface{}) *WaitGroupPromise {
	promise := &WaitGroupPromise{}
	promise.resetWithOptions(opts)
	promise.SetProvider(promise)
	return promise
}

func (p *WaitGroupPromise) Reset() {
	p.ResetWithOptions(nil)
}

func (p *WaitGroupPromise) ResetWithOptions(opts interface{}) {
	p.Resolve(nil, ErrReset) // This will release all waiting goroutines
	runtime.Gosched()        // Give waiting goroutines a chance to run, so wg can be reused.

	p.resetWithOptions(opts)
}

func (p *WaitGroupPromise) Resolve(rets ...interface{}) (Promise, error) {
	if !p.wg.IsWaiting() {
		return p, ErrResolved
	}

	p.AbstractPromise.ResolveRets(rets...)
	p.wg.Done()
	return p, nil
}

func (p *WaitGroupPromise) Timeout(timeouts ...time.Duration) error {
	ch, err := p.TimeoutC(timeouts...)
	if err == ErrResolved {
		return nil
	} else if err != nil {
		return err
	}

	cond := make(chan struct{})
	go func() {
		p.wg.Wait()
		close(cond)
	}()

	select {
	case <-ch:
		return ErrTimeout
	case <-cond:
		return nil
	}
}

// PromiseProvider
func (p *WaitGroupPromise) Wait() {
	p.wg.Wait()
}

func (p *WaitGroupPromise) Lock() {
}

func (p *WaitGroupPromise) Unlock() {
}

func (p *WaitGroupPromise) resetWithOptions(opts interface{}) {
	p.AbstractPromise.ResetWithOptions(opts)
	p.wg.Add(1)
}
