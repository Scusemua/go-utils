package promise

import (
	"runtime"
	"sync"
	"testing"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var (
	ret = &struct{}{}
)

type TimeoutResult bool

func (t TimeoutResult) String() string {
	if t {
		return "timeout"
	} else {
		return "not timeout"
	}
}

func TestTypes(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Promise")
}

func shouldNotTimeout[V any](test func() V, expects ...TimeoutResult) V {
	expect := TimeoutResult(false)
	if len(expects) > 0 {
		expect = expects[0]
	}

	timer := time.NewTimer(time.Second)
	timeout := TimeoutResult(false)
	responeded := make(chan V)
	var ret V
	go func() {
		responeded <- test()
	}()
	select {
	case <-timer.C:
		timeout = TimeoutResult(true)
	case ret = <-responeded:
		if !timer.Stop() {
			<-timer.C
		}
	}

	Expect(timeout).To(Equal(expect))
	return ret
}

func shouldTimeout[V any](test func() V) {
	shouldNotTimeout(test, TimeoutResult(true))
}

var _ = Describe("Promise", func() {
	It("no wait if data has been available already", func() {
		promise := NewPromise()
		promise.Resolve(ret)

		Expect(shouldNotTimeout(promise.Value)).To(Equal(ret))
	})

	It("should wait if data is not available", func() {
		promise := NewPromise()

		shouldTimeout(promise.Value)
	})

	It("should wait until data is available", func() {
		promise := NewPromise()

		var done sync.WaitGroup
		done.Add(1)
		go func() {
			Expect(shouldNotTimeout(promise.Value)).To(Equal(ret))
			done.Done()
		}()
		runtime.Gosched()

		<-time.After(500 * time.Millisecond)
		promise.Resolve(ret)

		done.Wait()
	})

	It("should unblock on reset", func() {
		promise := NewPromise()

		var done sync.WaitGroup
		done.Add(1)
		go func() {
			<-time.After(500 * time.Millisecond)
			promise.Reset()
			done.Done()
		}()
		Expect(shouldNotTimeout(promise.Error)).To(Equal(ErrReset))

		// Wait for reset
		done.Wait()

		shouldTimeout(promise.Value)

		_, err := promise.Resolve(ret)
		Expect(err).To(BeNil())
		Expect(promise.Value()).To(Equal(ret))
	})

	It("should timeout as expected", func() {
		promise := NewPromise()
		promise.SetTimeout(100 * time.Millisecond)

		var done sync.WaitGroup
		done.Add(1)
		go func() {
			shouldTimeout(promise.Value)
			done.Done()
		}()
		runtime.Gosched()

		Expect(shouldNotTimeout(func() interface{} {
			return promise.Timeout()
		})).To(Equal(ErrTimeout))

		done.Wait()
	})

	It("should not timeout if value has been available", func() {
		promise := NewPromise()
		promise.Resolve(ret, nil)
		promise.SetTimeout(2000 * time.Millisecond)

		Expect(shouldNotTimeout(func() interface{} {
			return promise.Timeout()
		})).To(BeNil())

		Expect(shouldNotTimeout(promise.Value)).To(Equal(ret))
		Expect(promise.Error()).To(BeNil())
	})

	It("should not timeout if value is available", func() {
		promise := NewPromise()
		promise.SetTimeout(100 * time.Millisecond)

		var done sync.WaitGroup
		done.Add(1)
		go func() {
			promise.Resolve(ret, nil)
			done.Done()
		}()

		Expect(shouldNotTimeout(func() interface{} {
			return promise.Timeout()
		})).To(BeNil())

		Expect(shouldNotTimeout(promise.Value)).To(Equal(ret))
		Expect(promise.Error()).To(BeNil())

		done.Wait()
	})

	It("should not timeout multiple times", func() {
		promise := NewPromise()

		var done sync.WaitGroup
		done.Add(1)
		go func() {
			<-time.After(250 * time.Millisecond)
			promise.Resolve(ret, nil)
			done.Done()
		}()

		promise.SetTimeout(100 * time.Millisecond)
		Expect(shouldNotTimeout(func() interface{} {
			return promise.Timeout()
		})).To(Equal(ErrTimeout))

		promise.SetTimeout(100 * time.Millisecond)
		Expect(shouldNotTimeout(func() interface{} {
			return promise.Timeout()
		})).To(Equal(ErrTimeout))

		promise.SetTimeout(100 * time.Millisecond)
		Expect(shouldNotTimeout(func() interface{} {
			return promise.Timeout()
		})).To(BeNil())

		Expect(shouldNotTimeout(promise.Value)).To(Equal(ret))
		Expect(promise.Error()).To(BeNil())

		done.Wait()
	})

	It("should support concurrent timeout", func() {
		promise := NewPromise()
		concurrency := 10

		var done sync.WaitGroup
		done.Add(1)
		go func() {
			<-time.After(250 * time.Millisecond)
			promise.Resolve(ret, nil)
			done.Done()
		}()

		promise.SetTimeout(100 * time.Millisecond)
		for i := 0; i < concurrency; i++ {
			done.Add(1)
			go func() {
				// defer GinkgoRecover()

				Expect(shouldNotTimeout(func() interface{} {
					return promise.Timeout()
				})).To(Equal(ErrTimeout))
				done.Done()
			}()
		}
		done.Wait()

		promise.SetTimeout(200 * time.Millisecond)
		for i := 0; i < concurrency; i++ {
			done.Add(1)
			go func() {
				defer GinkgoRecover()

				Expect(shouldNotTimeout(func() interface{} {
					return promise.Timeout()
				})).To(BeNil())
				done.Done()
			}()
		}

		Expect(shouldNotTimeout(promise.Value)).To(Equal(ret))
		Expect(promise.Error()).To(BeNil())

		done.Wait()
	})
})

func BenchmarkNewChannelPromise(b *testing.B) {
	for i := 0; i < b.N; i++ {
		promise := NewChannelPromise()
		promise.Close()
	}
}

func BenchmarkNewSyncPromise(b *testing.B) {
	for i := 0; i < b.N; i++ {
		promise := NewSyncPromise()
		promise.Close()
	}
}

func BenchmarkNewWaitGroupPromise(b *testing.B) {
	for i := 0; i < b.N; i++ {
		promise := NewWaitGroupPromise()
		promise.Close()
	}
}

func BenchmarkChannelPromiseResolvedCheck(b *testing.B) {
	for i := 0; i < b.N; i++ {
		promise := NewChannelPromise()
		if !promise.IsResolved() {
			promise.Close()
		}
	}
}

func BenchmarkSyncPromiseResolvedCheck(b *testing.B) {
	for i := 0; i < b.N; i++ {
		promise := NewSyncPromise()
		if !promise.IsResolved() {
			promise.Close()
		}
	}
}

func BenchmarkWaitGroupPromiseResolvedCheck(b *testing.B) {
	for i := 0; i < b.N; i++ {
		promise := NewWaitGroupPromise()
		if !promise.IsResolved() {
			promise.Close()
		}
	}
}

func BenchmarkChannelPromiseNotification(b *testing.B) {
	for i := 0; i < b.N; i++ {
		promise := NewChannelPromise()
		go func() {
			promise.Resolve(ret)
		}()
		promise.Value()
	}
}

func BenchmarkSyncPromiseNotification(b *testing.B) {
	for i := 0; i < b.N; i++ {
		promise := NewSyncPromise()
		go func() {
			promise.Resolve(ret)
		}()
		promise.Value()
	}
}

func BenchmarkWaitGroupPromiseNotification(b *testing.B) {
	for i := 0; i < b.N; i++ {
		promise := NewWaitGroupPromise()
		go func() {
			promise.Resolve(ret)
		}()
		promise.Value()
	}
}

func BenchmarkPromiseNotificationWithRecycling(b *testing.B) {
	for i := 0; i < b.N; i++ {
		promise := NewPromise()
		go func() {
			promise.Resolve(ret)
		}()
		promise.Value()
		Recycle(promise)
	}
}
