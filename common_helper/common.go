package common_helper

import "reflect"

// If 模拟三目
// 任意类型的三目
func If(isTrue bool, trueRes, falseRes interface{}) interface{} {
	if isTrue {
		return trueRes
	}
	return falseRes
}

// IfWithType 带泛型的三目,因为三目的结果只有1个
// 指定类型的三目 可以指定的类型为任意类型
func IfWithType[T interface{}](isTrue bool, trueRes, falseRes T) T {
	if isTrue {
		return trueRes
	}
	return falseRes
}

// IfWithFunc 带泛型的func三目
// 指定的结果用func返回,不提前计算
func IfWithFunc[T interface{}](isTrue bool, trueResFunc, falseResFunc func() T) T {
	if isTrue {
		return trueResFunc()
	}
	return falseResFunc()
}

// IsNil 判断是否为Null
func IsNil(value any) bool {
	if value == nil {
		return true
	}
	reflectValue := reflect.ValueOf(value)
	switch reflectValue.Kind() {
	case reflect.Chan, reflect.Func, reflect.Map, reflect.Ptr, reflect.UnsafePointer,
		reflect.Interface, reflect.Slice:
		if reflectValue.IsNil() {
			return true
		}
	}
	return false
}
