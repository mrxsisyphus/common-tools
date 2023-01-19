package collection_helper

/*
切片相关的函数作用
主要是 Map filter reduce 等等

*/

import (
	"fmt"
	"github.com/mrxtryagin/common-tools/convert_helper"
	"github.com/mrxtryagin/common-tools/optional"
	"sort"
	"strings"
)

// Map 将el 做一层函数映射
// - InputType any 初始的数组元素类型
// - ReturnType any 返回的map的值的类型
// - convertor func(InputType) ReturnType 转换器
func Map[InputType any, ReturnType any](input *[]InputType, convertor func(index int, inputType InputType) ReturnType) *[]ReturnType {
	raw := *input
	if len(raw) <= 0 {
		return &([]ReturnType{})
	}
	res := make([]ReturnType, len(raw), len(raw))
	for index, item := range raw {
		res[index] = convertor(index, item)
	}
	return &res
}

// Filter 将el做过滤
// - InputType any 初始的数组元素类型
// - predicate func(InputType) bool  布尔值断言,为true的加入到最新的结果集中
func Filter[InputType any](input *[]InputType, predicate func(index int, inputType InputType) bool) *[]InputType {
	raw := *input
	if len(raw) <= 0 {
		return input
	}
	result := make([]InputType, 0)
	for index, item := range raw {
		if predicate(index, item) {
			result = append(result, item)
		}
	}
	return &result
}

// Reduce 将 el做 reduce
// - InputType any 初始的数组元素类型
// - accumulator func(x, y InputType) InputType 累加器 用来更新accumulateVal
// 返回optional
func Reduce[InputType any](input *[]InputType, accumulator func(x, y InputType) InputType) *optional.Optional[InputType] {
	raw := *input
	var accumulateVal InputType
	if len(raw) <= 0 {
		return optional.Empty[InputType]()
	}
	foundAny := false
	for _, item := range raw {
		if !foundAny {
			foundAny = true
			accumulateVal = item // 第一个元素
		} else {
			//其余元素
			accumulateVal = accumulator(item, accumulateVal)
		}

	}
	return optional.OfNullable(accumulateVal)
}

// ReduceWithInitValueInSameType 将 el做 reduce,有初始值 initValue 类型 和 元素类型 相同
// - InputType any 初始的数组元素类型
// - accumulator func(x, y InputType) InputType 累加器 用来更新accumulateVal
// - initValue InputType 初始值
func ReduceWithInitValueInSameType[InputType any](input *[]InputType, accumulator func(x, y InputType) InputType, initValue InputType) InputType {
	raw := *input
	accumulatorVal := initValue
	if len(raw) <= 0 {
		// 返回初始值
		return accumulatorVal
	}
	for _, item := range raw {
		accumulatorVal = accumulator(item, accumulatorVal)
	}
	return accumulatorVal

}

// ReduceWithInitValue 将 el做 reduce,initValue 类型 和 元素类型 不同
// - InputType any 初始的数组元素类型
// - AccumulatorValType any 累加值的类型
// - accumulator func(x InputType, y AccumulatorValType) AccumulatorValType
// - initValue AccumulatorValType 初始值
func ReduceWithInitValue[InputType any, AccumulatorValType any](initAccumulateVal AccumulatorValType, input *[]InputType, accumulator func(x InputType, accumulateVal AccumulatorValType) AccumulatorValType) AccumulatorValType {
	raw := *input
	accumulatorVal := initAccumulateVal
	if len(raw) <= 0 {
		// 返回初始值
		return accumulatorVal
	}
	for _, item := range raw {
		accumulatorVal = accumulator(item, accumulatorVal)
	}
	return accumulatorVal

}

// Join 数组拼接 sep 为参数
func Join[InputType any](input *[]InputType, sep string) string {
	// 转换为strings
	raw := *input
	switch len(raw) {
	case 0:
		return ""
	case 1:
		return convert_helper.AnyToString(raw[0])
	}
	res := Map[InputType, string](input, func(index int, input InputType) string {
		return convert_helper.AnyToString(input)
	})
	return strings.Join(*res, sep)

}

// CountSliceByFunc  满足某个条件的数量
// -apply func(InputType) bool)  满足条件的函数
func CountSliceByFunc[InputType any](input *[]InputType, apply func(InputType) bool) int64 {
	raw := *input
	if len(raw) <= 0 {
		return 0
	}
	var total int64
	for _, item := range raw {
		if apply(item) {
			total++
		}
	}
	return total
}

// Reverse 反转数组,不会改变原数组
// - InputType any 初始的数组元素类型
func Reverse[InputType any](input *[]InputType) *[]InputType {
	raw := *input
	if len(raw) <= 0 {
		return input
	}
	result := make([]InputType, len(raw), len(raw))
	// 深拷贝一份再排序 copy 只看len refer:https://blog.csdn.net/rickie7/article/details/105869252
	copy(result, raw)
	sort.SliceStable(result, func(i, j int) bool {
		return true //全交换
	})
	return &result
}

// Reversed 反转数组,会改变原数组
// - InputType any 初始的数组元素类型
func Reversed[InputType any](input *[]InputType) {
	raw := *input
	if len(raw) <= 0 {
		return
	}
	sort.SliceStable(raw, func(i, j int) bool {
		return true //全交换
	})
}

// Clear 切片变为0
func Clear[InputType any](input *[]InputType) {
	*input = (*input)[:0]
}

// BatchRemove 批量删除
// input input数组
// tempData 需要删除或者保留的数组
// removeOrRetain 删除还是保留,如果false 则是删除,true 则是保留
// 返回删除的 / 还保留的的数量
func BatchRemove[InputType any](input *[]InputType, tempData *[]InputType, removeOrRetain bool) int {
	raw := *input
	size := len(raw)
	if size <= 0 {
		return 0
	}
	temp := *tempData
	if len(temp) <= 0 {
		if removeOrRetain {
			// 意为保留0个,也就是全部删除
			Clear(input)
			return size
		} else {
			// 意为删除0个,就什么都不动
			return 0
		}
	}
	var leftNum int
	// 参考java  batchRemove
	var r, w int
	for ; r < size; r++ {
		// 如果removeOrRetain 为false(删除),且不包含,则说明  raw[r] 需要保留
		// 如果removeOrRetain 为 true(保留) 且包含, 则说明  raw[r] 需要保留
		// 把需要被移除的数据都替换掉,不需要移除的数据迁移
		if ContainsAny(tempData, raw[r]) == removeOrRetain {
			// 这里确定的总是 需要保留的元素
			raw[w] = raw[r]
			w++ //w就是要保留的个数
		}
	}
	leftNum = w
	if w != size {
		// 不等于 说明不是全部需要留下来,只需要留下来w下标之前的
		// 直接切片处理一下(底层数组并没有发生改变)
		// refer: https://zhuanlan.zhihu.com/p/526731603
		*input = raw[:w:w] // 把新切片的cap也矫正
	}
	if removeOrRetain {
		//如果保留返回保留的
		return leftNum
	} else {
		// 如果不是保留返回 所有的-保留的
		return size - leftNum
	}

}

// PrintSlice 打印slices
// 打印 切片内容 长度 容量 第一个切片的地址
func PrintSlice[InputType any](input *[]InputType, tag string) {
	raw := *input
	fmt.Printf("%s 底层数组: %v, len: %d cap: %d ,pointer: %p\n", tag, raw, len(raw), cap(raw), raw) //0x1400011e1c0
}
