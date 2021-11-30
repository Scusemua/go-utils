package mapreduce

import (
	"reflect"
)

type Container interface {
	Container() interface{}
}

type Enumerator interface {
	Len() int
	Item(int) interface{}
}

type ArrayEnumerator struct {
	// Array exposes backend reflect.Slice or reflect.Array
	array reflect.Value
}

func NewEnumerator(slice interface{}) (Enumerator, error) {
	srcT := reflect.TypeOf(slice)
	kind := srcT.Kind()
	if kind == reflect.Slice || kind == reflect.Array {
		return &ArrayEnumerator{array: reflect.ValueOf(slice)}, nil
	} else {
		return nil, ErrNotSupported
	}
}

func (enum *ArrayEnumerator) Len() int {
	return enum.array.Len()
}

func (enum *ArrayEnumerator) Cap() int {
	return enum.array.Cap()
}

func (enum *ArrayEnumerator) Container() interface{} {
	return enum.array.Interface()
}

func (enum *ArrayEnumerator) Item(i int) interface{} {
	return enum.array.Index(i).Interface()
}

type EnumeratorIterator struct {
	Enumerator
	index int
}

func (enum *EnumeratorIterator) Next() bool {
	enum.index++
	return enum.index < enum.Len()
}

func (enum *EnumeratorIterator) Value() (int, interface{}) {
	return enum.index, enum.Item(enum.index)
}
