package collection_helper

// Contains 包含 入参input 是否包含el K 必须是可比较对象
// - InputType any 初始的数组元素类型
func Contains[InputType comparable](input *[]InputType, els ...InputType) bool {
	raw := *input
	if len(raw) <= 0 {
		return false
	}
	// slices maps 包
	for _, el := range els {
		index := Index(input, el)
		if index == -1 {
			return false
		}
	}
	return true

}

// ContainsAny 包含 入参input 是否包含el K 必须是可比较对象
// - InputType any 初始的数组元素类型
func ContainsAny[InputType any](input *[]InputType, els ...InputType) bool {
	raw := *input
	if len(raw) <= 0 {
		return false
	}
	// slices maps 包
	for _, el := range els {
		index := IndexAny(input, el)
		if index == -1 {
			return false
		}
	}
	return true

}

// Index
// - InputType comparable 初始的数组元素类型
func Index[InputType comparable](input *[]InputType, el InputType) int {
	raw := *input
	if len(raw) <= 0 {
		return -1
	}
	for i, inputType := range raw {
		if inputType == el {
			return i
		}
	}
	return -1

}

// IndexAny
// - InputType comparable 初始的数组元素类型
func IndexAny[InputType any](input *[]InputType, el InputType) int {
	raw := *input
	if len(raw) <= 0 {
		return -1
	}
	for i, inputType := range raw {
		if any(inputType) == any(el) {
			return i
		}
	}
	return -1

}

// ContainsWithFunc 包含 入参input 是否包含el K 可以是任意对象
// - InputType any 初始的数组元素类型
// - judgeIsContains func(index int, x, y InputType) bool) bool 判断这个el 是否算包含(包不包含自己定)
func ContainsWithFunc[InputType any](input *[]InputType, el InputType, judgeIsContains func(index int, x, y InputType) bool) bool {
	raw := *input
	if len(raw) <= 0 {
		return false
	}
	for index, item := range raw {
		if judgeIsContains(index, el, item) {
			return true
		}
	}
	return false
}

// FindWithFunc 通过func找到第一个 如果找到会返回 res InputType, isExist bool 类似于 slices.IndexFunc()
// - InputType any 初始的数组元素类型
// - findFunc findFunc func(InputType) bool) 判断是对应的元素是否是想要的
func FindWithFunc[InputType any](input *[]InputType, findFunc func(InputType) bool) (res InputType, isExist bool) {
	raw := *input
	if len(raw) <= 0 {
		return
	}
	for _, item := range raw {
		if findFunc(item) {
			return item, true
		}
	}
	return
}

// AllMatch 测试所有的情况是否都满足,有一个不满足就为false
// - InputType any 初始的数组元素类型
// -match func(InputType) 符不符合条件
func AllMatch[InputType any](input *[]InputType, match func(InputType) bool) bool {
	raw := *input
	if len(raw) <= 0 {
		return false
	}
	for _, item := range raw {
		if !match(item) {
			return false
		}
	}
	return true
}

// NoneMatch 测试所有的情况是否都不满足,有一个满足则为false
// - InputType any 初始的数组元素类型
// -match func(InputType) 符不符合条件
func NoneMatch[InputType any](input *[]InputType, match func(InputType) bool) bool {
	raw := *input
	if len(raw) <= 0 {
		return true
	}
	for _, item := range raw {
		if match(item) {
			return false
		}
	}
	return true
}

// AnyMatch 测试是否有任意情况满足,有一个满足则为true
// - InputType any 初始的数组元素类型
// -match func(InputType) 符不符合条件
func AnyMatch[InputType any](input *[]InputType, match func(InputType) bool) bool {
	raw := *input
	if len(raw) <= 0 {
		return false
	}
	for _, item := range raw {
		if match(item) {
			return true
		}
	}
	return false
}
