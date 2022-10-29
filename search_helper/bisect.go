package search_helper

import "golang.org/x/exp/constraints"

// BisectLeft 二分查找 定位 target 在seq中的插入点,保持原有seq的顺序不变,如果target已在seq中,那么插入点在已存在元素的左边
// - seq 需要查找的列表,被认为是有序的
// - target 需要定位的目标
// - startIndex 开始的索引 因为是索引不能小于0
// - end 结束的地方,  左闭右开
// - insertIndex 插入点所在的索引
// 返回的插入点 insertIndex 可以将数组 seq 分成两部分。
// 左侧是 all(val < x for val in a[lo:i]) ，右侧是 all(val >= x for val in a[i:hi]) 。
func BisectLeft[T constraints.Ordered](seq *[]T, target T, startIndex int, end int) (insertIndex int) {
	raw := *seq
	length := len(raw)
	if length <= 0 {
		panic("seq is empty")
	}
	if startIndex < 0 {
		panic("startIndex  must be non-negative")
	}
	if end > length {
		panic("end can not larger than seq_length")
	}
	if end == 0 {
		end = length
	}
	for startIndex < end {
		// 无符号右移
		midIndex := int(uint(startIndex+end) >> 1)

		if target > raw[midIndex] {
			// target 在中央值的右边
			startIndex = midIndex + 1
		} else {
			// 等于的情况下 startIndex 不变
			end = midIndex
		}
	}
	return startIndex
}

// BisectRight 二分查找 定位 target 在seq中的插入点,保持原有seq的顺序不变,如果target已在seq中,那么插入点在已存在元素的右边 不相等的情况与BisectLeft做法一致
// - seq 需要查找的列表,被认为是有序的
// - target 需要定位的目标
// - startIndex 开始的索引 因为是索引不能小于0
// - end 结束的地方,  左闭右开
// - insertIndex 插入点所在的索引
// 返回的插入点 insertIndex 可以将数组 seq 分成两部分。
// 左侧是 all(val <= x for val in a[lo:i])，右侧是 all(val > x for val in a[i:hi])
func BisectRight[T constraints.Ordered](seq *[]T, target T, startIndex int, end int) (insertIndex int) {
	raw := *seq
	length := len(raw)
	if length <= 0 {
		panic("seq is empty")
	}
	if startIndex < 0 {
		panic("startIndex  must be non-negative")
	}
	if end > length {
		panic("end can not larger than seq_length")
	}
	if end == 0 {
		end = length
	}
	for startIndex < end {
		// 无符号右移
		midIndex := int(uint(startIndex+end) >> 1)

		if target < raw[midIndex] {
			// target 在中央值的左侧
			end = midIndex
		} else {
			// 等于的情况下 midIndex + 1
			startIndex = midIndex + 1
		}
	}
	return startIndex
}

// Bisect 等价与BisectRight
func Bisect[T constraints.Ordered](seq *[]T, target T, startIndex int, end int) (insertIndex int) {
	return BisectRight[T](seq, target, startIndex, end)
}
