package optional

import (
	"errors"
	"fmt"
	"github.com/mrxtryagin/common-tools/common_helper"
)

type Optional[T any] struct {
	data *T //处理为指针可以更好地解决零值问题 refer:https://colobu.com/2022/04/26/Crimes-with-Go-Generics/
}

var (
	// ErrAbsent used to panic when call Optional.Get on empty Optional
	ErrAbsent = errors.New("absent value")
	// ErrNil used to panic when input a nil value to the Of method
	ErrNil = errors.New("nil value")
)

// Empty 空Optional 尝试Get他会报错
func Empty[T any]() *Optional[T] {
	return &Optional[T]{}
}

// Of Optional初始化
func Of[T any](value T) *Optional[T] {
	return &Optional[T]{
		data: &value,
	}
}

// OfNullable value如果为null,返回空optional,否则是和Of一样
func OfNullable[T any](value T) *Optional[T] {
	if common_helper.IsNil(value) {
		return Empty[T]()
	}
	return Of[T](value)
}

// Get 尝试获取,如果没有获取到会抛出 ErrNil错误
func (op *Optional[T]) Get() (res T, err error) {
	// 判断是否是nil 是nil的话 报错
	if common_helper.IsNil(op.data) {
		return res, ErrNil
	}
	return *op.data, nil
}

// IsPresent 如果value 不为null 则返回true 否则false
func (op *Optional[T]) IsPresent() bool {
	return !common_helper.IsNil(op.data)
}

// IfPresent 如果value 不为null 则执行fn
func (op *Optional[T]) IfPresent(consumer func(T)) {
	if !common_helper.IsNil(op.data) {
		consumer(*op.data)
	}
}

// Filter 如果value 存在 并且满足条件,返回它自己,否则返回empty
func (op *Optional[T]) Filter(predicate func(T) bool) *Optional[T] {
	if !op.IsPresent() {
		return op
	} else {
		if predicate(*op.data) {
			return Empty[T]()
		}
		return op
	}
}

// OrElse 如果value存在,返回value 否则返回给定的其他值(这里没办法泛型只能返回any)
func (op *Optional[T]) OrElse(other any) any {
	if op.IsPresent() {
		return op.data
	}
	return other
}

// OrElseGet 如果value存在,返回value 否则返回一个可以返回该值的函数
func (op *Optional[T]) OrElseGet(supplier func() T) T {
	if op.IsPresent() {
		return *op.data
	}
	return supplier()
}

// OrElseThrow 如果value存在,返回value 否则返回指定的错误
func (op *Optional[T]) OrElseThrow(err error) (res T, e error) {
	if op.IsPresent() {
		return *op.data, nil
	}
	return res, err
}

// OrElsePanic 如果value存在,返回value 否则panic
func (op *Optional[T]) OrElsePanic() T {
	if op.IsPresent() {
		return *op.data
	}
	panic(ErrAbsent)
}

// String 表示属性
func (op *Optional[T]) String() string {
	if op.IsPresent() {
		return fmt.Sprintf("Optional[%v]", *op.data)
	}
	return "Optional.empty"
}
