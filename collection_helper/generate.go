package collection_helper

import (
	"math"
	"math/rand"
)

// IntRange 从 start 到end的一个range,类似于python的range 但是不是生成器
// start,end 左闭右开
// 如果不满足条件 都返回 nil,默认step至少为1 为0也会变为1
func IntRange(start, end, step int) *[]int {
	if start < 0 || end < 0 || step < 0 || start >= end {
		return &([]int{})
	}
	if step == 0 {
		step = 1
	}
	length := int(math.Ceil(float64(end-start) / float64(step)))
	// 如果length 算出来是0 直接返回
	if length == 0 {
		return &([]int{})
	}
	res := make([]int, 0, length)
	for i := start; i < end; i += step {
		res = append(res, i)
	}
	return &res
}

// Range range的简化写法
func Range(end int) *[]int {
	return IntRange(0, end, 1)
}

// IntRandRange 一个seed 出来的 返回一个有n个元素的，[start, end)范围内整数的伪随机排列的切片
// 参考了 rand.Perm
func IntRandRange(start, end int) *[]int {
	if start < 0 || end < 0 || start >= end {
		return &([]int{})
	}
	//rand.Seed(time.Now().UnixNano())
	length := end - start
	// res 是从0开始的
	res := make([]int, length)
	// start 和 end 是 外面的设定的范围
	for i := start; i < end; i++ {
		//数组的位置(是偏移了这么多) 值是i
		index := i - start
		// 找一个数组的随机位置
		j := rand.Intn(index + 1)
		res[index] = res[j]
		res[j] = i
	}
	return &res
}

// RandRange 一个seed 出来的,返回一个有n个元素的，[0,n)范围内整数的伪随机排列的切片
func RandRange(n int) *[]int {
	return IntRandRange(0, n)
}

// Maker 批量制造某类型的切片
func Maker[T any](size int, supplier func(index int) T) *[]T {
	res := make([]T, size, size)
	for i := 0; i < size; i++ {
		res[i] = supplier(i)
	}
	return &res
}
