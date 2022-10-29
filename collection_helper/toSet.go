package collection_helper

/*
数组转化为 MapSet
数组转化为 类HashSet(第三方)
*/
import mapset "github.com/deckarep/golang-set/v2"

// ToMapSet 因为go没有Set 用map的Set代替,使用Key作为去重标准,Key必须为可比较的,取第一个不重复的
// - distinctKeyFunc func(InputType) KeyType 按照那个维度进行去重
// - InputType comparable 初始的数组元素类型
// 返回 map_helper[KeyType]bool类型 特别的是如果是对象直接取Set的话返回的是Set自身
// 本质hashmap 无序
func ToMapSet[InputType comparable, KeyType comparable](input *[]InputType, distinctKeyFunc func(InputType) KeyType) map[InputType]struct{} {
	raw := *input
	// 只用key部分
	keys := make(map[KeyType]struct{})
	values := make(map[InputType]struct{})
	for _, item := range raw {
		key := distinctKeyFunc(item)
		if _, exist := keys[key]; !exist {
			keys[key] = struct{}{}
			values[item] = struct{}{}
		}
	}
	return values
}

// ToItSelfMapSet  与ToMapSet类似 但是 是自去重
// - InputType any 初始的数组元素类型
// 返回 map_helper[KeyType]bool类型 特别的是如果是对象直接取Set的话返回的是Set自身
// 本质hashmap 无序
func ToItSelfMapSet[InputType comparable](input *[]InputType) map[InputType]struct{} {
	raw := *input
	// 只用key部分
	res := make(map[InputType]struct{})
	for _, item := range raw {
		if _, ok := res[item]; !ok {
			//不存在 放入
			res[item] = struct{}{}
		}
	}
	return res
}

// ToSet 和  ToMapSet 类似 只是使用 包 https://github.com/deckarep/golang-set
func ToSet[InputType comparable, KeyType comparable](input *[]InputType, distinctKeyFunc func(InputType) KeyType) mapset.Set[InputType] {
	raw := *input
	// 只用key部分
	keys := make(map[KeyType]struct{})
	mySet := mapset.NewSet[InputType]()
	for _, item := range raw {
		key := distinctKeyFunc(item)
		if _, exist := keys[key]; !exist {
			keys[key] = struct{}{}
			mySet.Add(item)
		}
	}
	return mySet
}
