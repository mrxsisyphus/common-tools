package commonMapOperation

import (
	mapset "github.com/deckarep/golang-set/v2"
	"github.com/monitor1379/yagods/sets/hashset"
	"github.com/monitor1379/yagods/sets/linkedhashset"
	"github.com/monitor1379/yagods/sets/treeset"
	"github.com/mrxtryagin/common-tools/map_helper"
)

// refer: https://blog.csdn.net/qq_43679056/article/details/104976819
// 返回零值都加个bool判断存不存在
type MapAbsOperation[K comparable, V any] interface {
	//Get 获得K的V 没有返回false
	Get(K) (V, bool)
	//GetOrDefault get到了 返回get的 否则返回default
	GetOrDefault(K, V) V
	//Put 返回一个旧V,如果k不存在返回V的零值,和 false
	Put(K, V) (V, bool)
	//PutIfAbsent 和 Put 类似 只是 只有不存在的时候 才进行put,返回旧的值
	PutIfAbsent(K, V) (V, bool)
	//PutOrDefault 查找K K不存在时,设置V 存在时 不设置返回V,返回新的值,与 PutIfAbsent的区别在于不存在时,返回新的值
	PutOrDefault(K, V) V
	//Replace replace key 当前仅当key存在时,返回旧的值 和是否存在
	Replace(K, V) (V, bool)
	//ReplaceWith 替换k值,如果k和value 都存在 用 新value替换 并返回true
	ReplaceWith(K, V, V) bool
	//PutAll 将参数的map 中的 key 和value 全部放入 原map中 类似于 maps.Copy()
	PutAll(map[K]V)
	//Compute K存不存在 计算都会进行,使用K和旧的val 进行计算产生的新的val 重新赋值给k,通过bool 来判断该值原本存不存在
	Compute(K, func(k K, oldVal V) V) (V, bool)
	//ComputeIfAbsent Compute的特例 k不存在的时候,做映射,存在直接返回
	ComputeIfAbsent(K, func(K) V) (V, bool)
	//ComputeIfPresent  Compute的特例 k存在的时候,做映射,不存在直接返回
	ComputeIfPresent(K, func(k K, oldVal V) V) (V, bool)
	//Merge 和Compute 类似,不同的是,会多传入一个值,如果不存在,直接赋值,存在,会吧旧值,新值一起算一个值进行赋值, 常用于分组
	Merge(k K, val V, fn func(oldVal V, newVal V) V) (newVal V, exists bool)
	Remove(K) V
	//RemoveValue 使用kv来进行remove
	RemoveValue(K, V) bool
	//RemoveWithFunc 使用函数来remove
	RemoveWithFunc(func(K, V) bool)
	Keys() *[]K
	Values() *[]V
	ContainsKey(K) bool
	ContainsValue(V) bool
	String() string
	Size() int
	Copy() map[K]V
	Clear()
	//Update 类似于PutAll
	Update(map[K]V)

	Each(func(K, V) bool)
	Iter() <-chan *map_helper.Pair[K, V]
	GetIterator() *map_helper.Iterator[*map_helper.Pair[K, V]]

	// 转化为对应的key
	KeyMapSet() mapset.Set[K]
	KeyThreadSafeMapSet() mapset.Set[K]
	KeyHashSet() *hashset.Set[K]
	KeyLinkedHashSet() *linkedhashset.Set[K]
	KeyTreeSet() *treeset.Set[K]
}
