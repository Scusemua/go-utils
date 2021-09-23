package mapreduce

import (
	"errors"
	"reflect"
)

// Reducer Error Collection
var (
	ErrNilReducer          = errors.New("reducer cannot be nil")
	ErrReducerNotFunc      = errors.New("reducer must be a function")
	ErrMissingInitialValue = errors.New("missing initial value")
	ErrInvalidReducer      = errors.New("invalid reducer signature, must be \"reducer(TypeA, TypeB, int) TypeA\" or \"reducer(TypeA, TypeB) TypeA\"")
)

type Initiator interface {
	Initiate() interface{}
}

// Reduce reduces iteratible objects into a scalar value.
// The reducer can be reducer(TypeA, TypeB, int) TypeA or reducer(TypeA, TypeB) TypeA
// The source accept array, slice, channel or Iterator.
func Reduce(reducer, source interface{}, initials ...interface{}) (interface{}, error) {
	var initialValue interface{}
	if len(initials) > 0 {
		initialValue = initials[0]
	}

	// Normalize source as iterator
	iterator, ok := source.(Iterator)
	if !ok {
		var err error
		iterator, err = NewIterator(source)
		if err != nil {
			return initialValue, err
		}
	}

	// Validate reducer
	if reducer == nil {
		return initialValue, ErrNilReducer
	}

	rk := reflect.TypeOf(reducer)
	if rk.Kind() != reflect.Func {
		return initialValue, ErrReducerNotFunc
	} else if rk.NumIn() < 2 || rk.NumIn() > 3 || rk.NumOut() != 1 || rk.In(0) != rk.Out(0) || (rk.NumIn() == 3 && rk.In(2) != reflect.TypeOf(0)) {
		return initialValue, ErrInvalidReducer
	}

	// Try generate initial value
	if len(initials) == 0 {
		initiator, ok := iterator.(Initiator)
		if !ok {
			return initialValue, ErrMissingInitialValue
		}
		initialValue = initiator.Initiate()
	}

	accumulator := reflect.ValueOf(initialValue)
	rv := reflect.ValueOf(reducer)
	for iterator.Next() {
		i, entry := iterator.Value()

		// call reducer via reflection
		if rk.NumIn() > 2 {
			accumulator = rv.Call([]reflect.Value{
				accumulator,            // send accumulator value
				reflect.ValueOf(entry), // send current source entry
				reflect.ValueOf(i),     // send current loop index
			})[0]
		} else {
			accumulator = rv.Call([]reflect.Value{
				accumulator,            // send accumulator value
				reflect.ValueOf(entry), // send current source entry
			})[0]
		}
	}

	return accumulator.Interface(), nil
}
