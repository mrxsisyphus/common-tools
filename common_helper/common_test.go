package common_helper

import (
	"fmt"
	"testing"
)

func Test(t *testing.T) {
	a := 1
	b := 0
	// 带泛型的三目,直接规定了返回值的类型
	withType := IfWithType[int](a > b, a, b)
	fmt.Println(withType)
	// a/b 直接这么写会/0异常,然后再返回,没有三目的效果
	//IfWithType[int](a > b, a, a/b)
	// 使用func 就不会有问题,因为func不会执行,只有在 a>b的时候才会执行
	withFunc := IfWithFunc[int](a > b, func() int {
		return a
	}, func() int {
		return a / b // 虽然a/b有除0异常 但是这里压根不会执行
	})
	fmt.Println(withFunc)

}
