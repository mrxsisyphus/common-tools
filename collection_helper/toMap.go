package collection_helper

/*
数组转化为 Map
数组转化为 Group
*/

// ToMap 将数组 变为 map_helper,处理key重复的情况
// - KeyType comparable key的类型
// - InputType any 初始的数组元素类型
// - ReturnType any 返回的map的值的类型
// - keyFunc func(InputType) KeyType 确定key的函数
// - valueFunc func(InputType) ReturnType 确定value的函数
// - duplicateHandlerFunc func(v1, v2 ReturnType)(ReturnType,error) 如果value 出现重复了的做法,可以抛出错误,默认不加 后者覆盖前者(v1 old,v2 new)
func ToMap[InputType any, KeyType comparable, ReturnType any](input *[]InputType, keyFunc func(InputType) KeyType, valueFunc func(InputType) ReturnType, duplicateHandlerFunc func(v1, v2 ReturnType) (ReturnType, error)) (map[KeyType]ReturnType, error) {
	raw := *input
	res := make(map[KeyType]ReturnType, len(raw))
	for _, item := range raw {
		key := keyFunc(item)
		value := valueFunc(item)
		oldValue, exist := res[key]
		if exist {
			if duplicateHandlerFunc != nil {
				// 如果已经存在了,调用函数判断
				lastValue, err := duplicateHandlerFunc(oldValue, value)
				if err != nil {
					return nil, err
				}
				// 目前 lastValue 与 value oldValue的值要相同
				res[key] = lastValue
			} else {
				//后面覆盖前面
				res[key] = value
			}
		} else {
			res[key] = value
		}
	}
	return res, nil
}

// ToSimpleMap 简易Map
// 与ToMap 类似 只是不处理key相同的情况(相同直接覆盖),忽略可能的错误
func ToSimpleMap[InputType any, KeyType comparable, ReturnType any](input *[]InputType, keyFunc func(InputType) KeyType, valueFunc func(InputType) ReturnType) map[KeyType]ReturnType {
	toSimpleMap, _ := ToMap[InputType, KeyType, ReturnType](input, keyFunc, valueFunc, nil)
	return toSimpleMap
}

// ToIndexMap indexMap 使用数组下标和数组元素变为一个map
func ToIndexMap[InputType any](input *[]InputType) map[int]InputType {
	raw := *input
	res := make(map[int]InputType, len(raw))
	for index, item := range raw {
		res[index] = item
	}
	return res
}

// ToItselfMap itselfMap 与 简易map类似 不处理key相同的情况(相同直接覆盖),忽略可能的错误
// - InputType any 初始的数组元素类型
// map的value 就是元素自身
func ToItselfMap[InputType any, KeyType comparable](input *[]InputType, keyFunc func(InputType) KeyType) map[KeyType]InputType {
	// value 返回自己
	valueFunc := func(input InputType) InputType {
		return input
	}
	toItSelfMap, _ := ToMap[InputType, KeyType, InputType](input, keyFunc, valueFunc, nil)
	return toItSelfMap
}

// ToGroup 将数组 变为 group,
// - KeyType comparable key的类型
// - InputType any 初始的数组元素类型
// - ReturnType any 返回的值的类型 这个将会变为改值的切片
// - keyFunc func(InputType) KeyType 确定key的函数
// - valueFunc func(InputType) ReturnType 确定value的函数
func ToGroup[InputType any, KeyType comparable, ReturnType any](input *[]InputType, keyFunc func(InputType) KeyType, valueFunc func(InputType) ReturnType) map[KeyType][]ReturnType {
	raw := *input
	res := make(map[KeyType][]ReturnType)
	for _, item := range raw {
		key := keyFunc(item)
		value := valueFunc(item)
		oldValues, exist := res[key]
		if exist {
			oldValues = append(oldValues, value)
			res[key] = oldValues
		} else {
			newValues := make([]ReturnType, 0, 1)
			newValues = append(newValues, value)
			res[key] = newValues
		}
	}
	return res
}

// ToItselfGroup 将数组 变为 group,
// - KeyType comparable key的类型
// - InputType any 初始的数组元素类型
// - keyFunc func(InputType) KeyType 确定key的函数
// - valueFunc func(InputType) ReturnType 确定value的函数
// map的value 就是元素自身
func ToItselfGroup[InputType any, KeyType comparable](input *[]InputType, keyFunc func(InputType) KeyType) map[KeyType][]InputType {
	// value 返回自己
	valueFunc := func(input InputType) InputType {
		return input
	}
	return ToGroup[InputType, KeyType, InputType](input, keyFunc, valueFunc)
}

// ToGroupCount  以某个维度做count,返回维度:count的map
// -KeyType comparable key的类型
// - countKeyFunc func(InputType) KeyType 按照那个维度进行来count
// - InputType any 初始的数组元素类型
func ToGroupCount[InputType any, KeyType comparable](input *[]InputType, countKeyFunc func(InputType) KeyType) map[KeyType]int64 {
	raw := *input
	if len(raw) <= 0 {
		return nil
	}
	// 只用key部分
	countKeyMap := make(map[KeyType]int64)
	for _, item := range raw {
		key := countKeyFunc(item)
		if _, exist := countKeyMap[key]; !exist {
			countKeyMap[key] = 1
		} else {
			countKeyMap[key]++
		}
	}
	return countKeyMap
}
