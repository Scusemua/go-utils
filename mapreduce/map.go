package mapreduce

import (
	"errors"
	"reflect"
	"sync"
)

// Map Error Collection
var (
	ErrNilMapper     = errors.New("mapper function cannot be nil")
	ErrMapperNotFunc = errors.New("mapper must be a function")
	ErrInvalidMapper = errors.New("invalid Mapper signature, must be \"mapper(TypeA) TypeB\" or \"mapper(TypeA, int) TypeB\"")
)

type Length interface {
	Len() int
}

type Capacitor interface {
	Length
	Cap() int
}

// Map maps iteratible objects into another slice.
// The mapper can be mapper(TypeA, int) TypeB or mapper(TypeA) TypeB
// The source accept array, slice, channel or Iterator.
func Map(mapper interface{}, source interface{}) (interface{}, error) {
	iterator, numIn, out, err := prepareMap(mapper, source)
	if err != nil {
		return nil, err
	}

	// create a waitgroup with length = source array length
	// we'll reduce the counter each time an entry finished processing
	mv := reflect.ValueOf(mapper)
	rets := reflect.ValueOf(out)
	for iterator.Next() {
		// Ensure capacity
		rets = reflect.Append(rets, reflect.Zero(rets.Type().Elem()))
		// one go routine for each entry
		i, entry := iterator.Value()

		//Call the transformation and store the result value
		var ret reflect.Value
		if numIn == 1 {
			ret = mv.Call([]reflect.Value{reflect.ValueOf(entry)})[0]
		} else {
			ret = mv.Call([]reflect.Value{reflect.ValueOf(entry), reflect.ValueOf(i)})[0]
		}

		//Store the transformation result into array of result
		rets.Index(i).Set(ret)
	}

	return rets.Interface(), nil
}

// ParallelMap maps iteratible objects in parallel.
func ParallelMap(mapper interface{}, source interface{}) (interface{}, error) {
	iterator, numIn, out, err := prepareMap(mapper, source)
	if err != nil {
		return nil, err
	}

	wg := &sync.WaitGroup{}
	mv := reflect.ValueOf(mapper)
	rets := reflect.ValueOf(out)
	for iterator.Next() {
		// Ensure capacity
		rets = reflect.Append(rets, reflect.Zero(rets.Type().Elem()))
		// one go routine for each entry
		wg.Add(1)
		go func(i int, entry interface{}) {
			//Call the transformation and store the result value
			var ret reflect.Value
			if numIn > 1 {
				ret = mv.Call([]reflect.Value{reflect.ValueOf(entry), reflect.ValueOf(i)})[0]
			} else {
				ret = mv.Call([]reflect.Value{reflect.ValueOf(entry)})[0]
			}

			//Store the transformation result into array of result
			rets.Index(i).Set(ret)

			//this go routine is done
			wg.Done()
		}(iterator.Value())
	}

	wg.Wait()
	return rets.Interface(), nil
}

func prepareMap(mapper interface{}, source interface{}) (Iterator, int, interface{}, error) {
	// Normalize source as iterator
	iterator, ok := source.(Iterator)
	if !ok {
		var err error
		iterator, err = NewIterator(source)
		if err != nil {
			return nil, 0, nil, err
		}
	}

	// Validate mapper
	if mapper == nil {
		return iterator, 0, nil, ErrNilMapper
	}

	mk := reflect.TypeOf(mapper)
	if mk.Kind() != reflect.Func {
		return iterator, 0, nil, ErrMapperNotFunc
	} else if mk.NumIn() < 1 || mk.NumIn() > 2 || mk.NumOut() != 1 || (mk.NumIn() == 2 && mk.In(1) != reflect.TypeOf(0)) {
		return iterator, mk.NumIn(), nil, ErrInvalidMapper
	}

	cap := 0
	if capacitor, ok := iterator.(Capacitor); ok {
		cap = capacitor.Cap()
	} else if length, ok := iterator.(Length); ok {
		cap = length.Len()
	}

	return iterator, mk.NumIn(), reflect.MakeSlice(reflect.SliceOf(mk.Out(0)), 0, cap).Interface(), nil
}
