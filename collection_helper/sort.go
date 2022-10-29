package collection_helper

import "sort"

// Sort 将 el做 sort,sort 会 copy 来源数组,sort后返回,不会修改原数组,不保证原始slice的位置
// - InputType any 初始的数组元素类型
// - compareFn func(a, b InputType) bool) 比较函数 < 按照 正序 排序 > 为逆序排序
// 单维度排序, 多维度排序 refer: https://learnku.com/articles/38269
func Sort[InputType any](input *[]InputType, compareFn func(a, b InputType) bool) *[]InputType {
	raw := *input
	if len(raw) <= 0 {
		return input
	}
	result := make([]InputType, len(raw), len(raw))
	// 深拷贝一份再排序 copy 只看len refer:https://blog.csdn.net/rickie7/article/details/105869252
	copy(result, raw)
	sort.Slice(result, func(i, j int) bool {
		return compareFn(result[i], result[j])
	})
	return &result
}

// Sorted 将 el做 sorted,  sorted会直接修改原数组.不保证原始slice的位置
// - InputType any 初始的数组元素类型
// - compareFn func(a, b InputType) bool) 比较函数 < 按照 正序 排序 > 为逆序排序
func Sorted[InputType any](input *[]InputType, compareFn func(a, b InputType) bool) {
	raw := *input
	if len(raw) <= 0 {
		return
	}
	sort.Slice(raw, func(i, j int) bool {
		return compareFn(raw[i], raw[j])
	})
}

// SortStable  将 el做 sort,sort 会 copy 来源数组,sort后返回,不会修改原数组,保证原始slice的位置
// - InputType any 初始的数组元素类型
// - compareFn func(a, b InputType) bool) 比较函数 < 按照 正序 排序 > 为逆序排序
// 单维度排序, 多维度排序 refer: https://learnku.com/articles/38269
func SortStable[InputType any](input *[]InputType, compareFn func(a, b InputType) bool) *[]InputType {
	raw := *input
	if len(raw) <= 0 {
		return input
	}
	result := make([]InputType, len(raw), len(raw))
	// 深拷贝一份再排序 copy 只看len refer:https://blog.csdn.net/rickie7/article/details/105869252
	copy(result, raw)
	sort.SliceStable(result, func(i, j int) bool {
		return compareFn(result[i], result[j])
	})
	return &result
}

// SortedStable 将 el做 sorted,  sorted会直接修改原数组,不保证原始slice的位置
// - InputType any 初始的数组元素类型
// - compareFn func(a, b InputType) bool) 比较函数 < 按照 正序 排序 > 为逆序排序
func SortedStable[InputType any](input *[]InputType, compareFn func(a, b InputType) bool) {
	raw := *input
	if len(raw) <= 0 {
		return
	}
	sort.SliceStable(raw, func(i, j int) bool {
		return compareFn(raw[i], raw[j])
	})
}
