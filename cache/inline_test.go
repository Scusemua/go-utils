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

func NewTypeA(chained bool, validate bool) *TypeA {
	typeA := &TypeA{}
	if chained {
		typeA.variable.Producer = cache.FormalizeChainedICProducer(typeA.chainedCostOperation)
	} else {
		typeA.variable.Producer = cache.FormalizeICProducer(typeA.costOperation)
	}
	if validate {
		typeA.variable.Validator = cache.FormalizeICValidator(typeA.validate)
	}
	return typeA
}

func (f *TypeA) GetVariable(arg int64) float64 {
	f.test = arg
	return f.variable.Value(arg).(float64)
}

func (f *TypeA) costOperation(arg int64) float64 {
	return float64(arg)
}

func (f *TypeA) chainedCostOperation(cached float64, arg int64) float64 {
	return float64(arg)
}

func (f *TypeA) validate(cached float64) bool {
	return float64(f.test) == cached
}

var _ = Describe("InlineCache", func() {
	It("should example works", func() {
		a := NewTypeA(false, false)
		Expect(a.GetVariable(1)).To(Equal(1.0))
		Expect(a.GetVariable(2)).To(Equal(1.0))

		b := NewTypeA(true, false)
		Expect(b.GetVariable(1)).To(Equal(1.0))
		Expect(b.GetVariable(1)).To(Equal(1.0))

		c := NewTypeA(false, true)
		Expect(c.GetVariable(1)).To(Equal(1.0))
		Expect(c.GetVariable(2)).To(Equal(2.0))
	})
})
