package promise

import (
	"sync"
	"sync/atomic"
	"time"
)

type SyncPromise struct {
	AbstractPromise

	cond   *sync.Cond
	mu     sync.Mutex
	timers []*time.Timer
}

func ResolvedSync(rets ...interface{}) *SyncPromise {
	promise := NewSyncPromiseWithOptions(nil)
	if promise.ResolveRets(rets...) {
		promise.resolved = time.Now().UnixNano()
	} else {
		promise.resolved = int64(1) // Differentiate with PromiseInit
	}
	return promise
}

func NewSyncPromise() *SyncPromise {
	return NewSyncPromiseWithOptions(nil)
}

func NewSyncPromiseWithOptions(opts interface{}) *SyncPromise {
	promise := &SyncPromise{}
	promise.cond = sync.NewCond(&promise.mu)
	promise.timers = make([]*time.Timer, 0, 2)
	promise.ResetWithOptions(opts)
	promise.SetProvider(promise)
	return promise
}

func (p *SyncPromise) Resolve(rets ...interface{}) (Promise, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if atomic.LoadInt64(&p.resolved) != PromiseInit {
		return p, ErrResolved
	}

	p.ResolveRets(rets...)
	atomic.StoreInt64(&p.resolved, time.Now().UnixNano())
	p.cond.Broadcast()
	return p, nil
}

func (p *SyncPromise) Timeout() error {
	ch, err := p.TimeoutC()
	if err == ErrResolved {
		return nil
	} else if err != nil {
		return err
	}

	<-ch
	if p.IsResolved() {
		return nil
	} else {
		return ErrTimeout
	}
}

// PromiseProvider
func (p *SyncPromise) Wait() {
	p.mu.Lock()
	defer p.mu.Unlock()

	for atomic.LoadInt64(&p.resolved) == PromiseInit {
		p.cond.Wait()
	}

	for _, timer := range p.timers {
		timer.Stop()
	}
	p.timers = p.timers[:0]
}

func (p *SyncPromise) Lock() {
	p.mu.Lock()
}

func (p *SyncPromise) Unlock() {
	p.mu.Unlock()
}

func (p *SyncPromise) OnCreateTimerLocked(timer *time.Timer) {
	p.timers = append(p.timers, timer)
}
