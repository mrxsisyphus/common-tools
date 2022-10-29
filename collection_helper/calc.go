package collection_helper

/**
数组的累集器(参考python)
数组相关的 交集、并集、差集、补集、拆分、去重，排列组合.

*/

// Accumulate 累积器
// - InputType any inputType
// - AccumulateValType  累积的值的类型
// - apply func(index int, x AccumulateValType, y InputType) AccumulateValType 累积的函数
// 当index 为0是, apply的入参为 0,AccumulateValType的零值,InputType,共同决定AccumulateValType数组第一个index的最终值(因为无法动态)
// 当index 不为0时, apply的入参为 0,AccumulateValType[index-1](上一个值),InputType,共同决定AccumulateValType数组其他的值
func Accumulate[InputType any, AccumulateValType any](input *[]InputType, apply func(index int, x AccumulateValType, y InputType) AccumulateValType) *[]AccumulateValType {
	raw := *input
	if len(raw) <= 0 {
		return nil
	}
	res := make([]AccumulateValType, len(raw), len(raw))
	for index, item := range raw {
		if index == 0 {
			res[index] = apply(index, res[0], item)
		} else {
			res[index] = apply(index, res[index-1], item)
		}
	}
	return &res
}

// AccumulateWithSameType  Accumulate 的特例,转换值相同
func AccumulateWithSameType[InputType any](input *[]InputType, apply func(index int, x, y InputType) InputType) *[]InputType {
	return Accumulate[InputType, InputType](input, apply)
}

// arrOperation 数组操作 主要是 并集 交集 差集
// refer:https://blog.csdn.net/yiweiyi329/article/details/101030510
type setOperations[InputType any, KeyType comparable] struct {
	input      *[]InputType
	getkeyFunc func(InputType) KeyType
}

// Union 数组的并集 返回 slice1 并上 slice2
func (arr *setOperations[InputType, KeyType]) Union(slice2 *[]InputType) *[]InputType {
	slice1Arr := *arr.input
	slice2Arr := *slice2
	if len(slice1Arr) == 0 && len(slice2Arr) == 0 {
		return &([]InputType{})
	}
	// 用来去重的依据
	keys := make(map[KeyType]struct{}, len(slice1Arr))
	res := make([]InputType, 0, len(slice1Arr))
	for _, inputType := range slice1Arr {
		unionKey := arr.getkeyFunc(inputType)
		if _, exists := keys[unionKey]; !exists {
			keys[unionKey] = struct{}{}
			res = append(res, inputType)
		}
	}
	for _, inputType := range slice2Arr {
		unionKey := arr.getkeyFunc(inputType)
		if _, exists := keys[unionKey]; !exists {
			keys[unionKey] = struct{}{}
			res = append(res, inputType)
		}
	}
	return &res
}

// Intersect 数组的交集 返回 slice1 交上 slice2
func (arr *setOperations[InputType, KeyType]) Intersect(slice2 *[]InputType) *[]InputType {
	slice1Arr := *arr.input
	slice2Arr := *slice2
	slice1ArrLength := len(slice1Arr)
	slice2ArrLength := len(slice2Arr)
	if slice1ArrLength == 0 || slice2ArrLength == 0 {
		return &([]InputType{})
	}
	// 用来去重的依据
	keys := make(map[KeyType]struct{})
	res := make([]InputType, 0)
	for _, inputType := range slice1Arr {
		intersectKey := arr.getkeyFunc(inputType)
		if _, exists := keys[intersectKey]; !exists {
			keys[intersectKey] = struct{}{}
		}
	}
	for _, inputType := range slice2Arr {
		intersectKey := arr.getkeyFunc(inputType)
		if _, exists := keys[intersectKey]; exists {
			// 如果有存在的 才放
			res = append(res, inputType)
		}
	}
	return &res
}

// Difference 数组的差集
func (arr *setOperations[InputType, KeyType]) Difference(slice2 *[]InputType) *[]InputType {
	slice1Arr := *arr.input
	slice2Arr := *slice2
	slice1ArrLength := len(slice1Arr)
	slice2ArrLength := len(slice2Arr)
	if slice1ArrLength == 0 {
		return &([]InputType{})
	}
	if slice2ArrLength == 0 {
		return arr.input
	}

	// 用来去重的依据
	keys := make(map[KeyType]struct{})
	res := make([]InputType, 0)
	// slice1Arr 中有的 ,slice2中没有的,就放slice1Arr的
	// 所有的slice2
	for _, inputType := range slice2Arr {
		differenceKey := arr.getkeyFunc(inputType)
		if _, exists := keys[differenceKey]; !exists {
			keys[differenceKey] = struct{}{}
		}
	}
	// 循环slice1,如果 slice1 有的key 不在slice2里面,那这些key 就是slice1的
	for _, inputType := range slice1Arr {
		differenceKey := arr.getkeyFunc(inputType)
		if _, exists := keys[differenceKey]; !exists {
			res = append(res, inputType)
		}
	}
	return &res
}

// NewSetOperationWithFunc  使用input 和 keyFunc 来创建数组操作类
func NewSetOperationWithFunc[InputType any, KeyType comparable](input *[]InputType, getKeyFunc func(InputType) KeyType) *setOperations[InputType, KeyType] {
	return &setOperations[InputType, KeyType]{
		input:      input,
		getkeyFunc: getKeyFunc,
	}
}

// NewSetOperation  使用input 来创建Operation 这个时候
func NewSetOperation[InputType comparable](input *[]InputType) *setOperations[InputType, InputType] {
	return &setOperations[InputType, InputType]{
		input: input,
		getkeyFunc: func(input InputType) InputType {
			return input
		},
	}
}
