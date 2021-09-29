package cache_test

import (
	"testing"

	"github.com/mason-leap-lab/go-utils/cache"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestCache(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Cache")
}

type TypeA struct {
	variable cache.InlineCache
	test     int64
}

func NewTypeA(validate bool) *TypeA {
	typeA := &TypeA{}
	typeA.variable.Producer = cache.FormalizeChainedICProducer(typeA.costOperation)
	if validate {
		typeA.variable.Validator = cache.FormalizeICValidator(typeA.validate)
	}
	return typeA
}

func (f *TypeA) GetVariable(arg int64) float64 {
	f.test = arg
	return f.variable.Value(arg).(float64)
}

func (f *TypeA) costOperation(cached float64, arg int64) float64 {
	return float64(arg)
}

func (f *TypeA) validate(cached float64) bool {
	return float64(f.test) == cached
}

var _ = Describe("InlineCache", func() {
	It("should example works", func() {
		a := NewTypeA(false)
		Expect(a.GetVariable(1)).To(Equal(1.0))
		Expect(a.GetVariable(2)).To(Equal(1.0))

		b := NewTypeA(true)
		Expect(b.GetVariable(1)).To(Equal(1.0))
		Expect(b.GetVariable(2)).To(Equal(2.0))
	})
})
