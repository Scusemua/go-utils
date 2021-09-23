package mapreduce

import (
	"errors"
	"reflect"
)

var (
	ErrNotSupported   = errors.New("not supported")
	ErrInvalidChanDir = errors.New("not readable channel")
)

type Iterator interface {
	Next() bool
	Value() (int, interface{})
}

type ContainerIterator struct {
	container reflect.Value
	index     int
}

type ArrayIterator struct {
	ContainerIterator
}

type ChanIterator struct {
	ContainerIterator
}

func NewIterator(sliceOrChan interface{}) (Iterator, error) {
	srcT := reflect.TypeOf(sliceOrChan)
	kind := srcT.Kind()
	if kind == reflect.Slice || kind == reflect.Array {
		return &ArrayIterator{ContainerIterator: ContainerIterator{container: reflect.ValueOf(sliceOrChan), index: -1}}, nil
	} else if kind == reflect.Chan {
		if srcT.ChanDir() != reflect.SendDir {
			return nil, ErrInvalidChanDir
		}
		return &ChanIterator{ContainerIterator: ContainerIterator{container: reflect.ValueOf(sliceOrChan), index: -1}}, nil
	} else {
		return nil, ErrNotSupported
	}
}

func (iter *ContainerIterator) Len() int {
	return iter.container.Len()
}

func (iter *ContainerIterator) Cap() int {
	return iter.container.Cap()
}

func (iter *ArrayIterator) Next() bool {
	return iter.index+1 < iter.container.Len()
}

func (iter *ArrayIterator) Value() (int, interface{}) {
	iter.index++
	return iter.index, iter.container.Index(iter.index).Interface()
}

func (iter *ChanIterator) Next() bool {
	return iter.container.Len() > 0
}

func (iter *ChanIterator) Value() (int, interface{}) {
	iter.index++
	val, _ := iter.container.Recv()
	return iter.index, val.Interface()
}
