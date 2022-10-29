package collection_helper

import "golang.org/x/exp/maps"

// JudgeSliceIsDistinctByFunc 与DistinctByFunc类似, 用来以某个维度判断是否重复,如果有重复,返回true
// -KeyType comparable key的类型
// - distinctKeyFunc func(InputType) KeyType 重复func
// - InputType any 初始的数组元素类型
func JudgeSliceIsDistinctByFunc[InputType any, KeyType comparable](input *[]InputType, distinctKeyFunc func(InputType) KeyType) bool {
	raw := *input
	if len(raw) <= 0 {
		return false
	}
	// 只用key部分
	distinctKeyMap := make(map[KeyType]struct{})
	for _, item := range raw {
		key := distinctKeyFunc(item)
		if _, exist := distinctKeyMap[key]; !exist {
			distinctKeyMap[key] = struct{}{}
		} else {
			return true
		}
	}
	return false
}

// Distinct  与ToItSelfMapSet类似,但是返回的是slice
// - InputType any 初始的数组元素类型
// 返回的是slice
func Distinct[InputType comparable](input *[]InputType) *[]InputType {
	raw := *input
	if len(raw) <= 0 {
		return input
	}
	distinctMap := make(map[InputType]struct{})
	for _, item := range raw {
		distinctMap[item] = struct{}{}
	}
	res := maps.Keys(distinctMap)
	return &res
}

// DistinctByFunc 与ToMapSet类似,但是返回的是slice
// -KeyType comparable key的类型
// - distinctKeyFunc func(InputType) KeyType 按照那个维度进行去重
// - InputType any 初始的数组元素类型
func DistinctByFunc[InputType any, KeyType comparable](input *[]InputType, distinctKeyFunc func(InputType) KeyType) *[]InputType {
	raw := *input
	if len(raw) <= 0 {
		return input
	}
	// 只用key部分
	distinctKeyMap := make(map[KeyType]struct{})
	res := make([]InputType, 0)
	for _, item := range raw {
		key := distinctKeyFunc(item)
		if _, exist := distinctKeyMap[key]; !exist {
			distinctKeyMap[key] = struct{}{}
			res = append(res, item)
		}
	}
	return &res
}
