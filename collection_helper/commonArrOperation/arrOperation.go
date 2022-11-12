package commonArrOperation

import (
	mapset "github.com/deckarep/golang-set/v2"
	"github.com/monitor1379/yagods/lists/arraylist"
	"github.com/monitor1379/yagods/lists/doublylinkedlist"
	"github.com/monitor1379/yagods/lists/singlylinkedlist"
	"github.com/monitor1379/yagods/sets/hashset"
	"github.com/monitor1379/yagods/sets/linkedhashset"
	"github.com/monitor1379/yagods/sets/treeset"
	"github.com/monitor1379/yagods/stacks/arraystack"
	"github.com/monitor1379/yagods/stacks/linkedliststack"
	"github.com/monitor1379/yagods/trees/binaryheap"
	"github.com/mrxtryagin/common-tools/collection_helper"
	"github.com/mrxtryagin/common-tools/optional"
)

// ArrAbsOperation  数组操作,有些情况会直接panic
type ArrAbsOperation[T any] interface {

	// Add 增加
	Add(T)
	//Adds 批量增加
	Adds(...T)
	// AddAll 增加一个同类型的数组
	AddAll(*[]T)
	// Remove 删除
	Remove(T)
	//RemoveWithFunc 删除的是满足条件的
	RemoveWithFunc(func(T) bool)
	DeleteRange(i, j int)
	//Pop pop一个元素,index执行位置,返回对应的泛型,如果存在,不存在 则返回nil false
	Pop(index int) T
	//PopFirst pop第一个
	PopFirst() T
	//PopLast pop最后一个
	PopLast() T

	//Insert 将对象插入列表 将某个元素插入某个地方
	Insert(index int, ts ...T)

	//Index 找出某个值第一个匹配的位置
	Index(t T) int
	//IndexWithFunc 找出某个值第一个匹配的位置,使用某个函数
	IndexWithFunc(func(T) bool) int
	//Contains 判断是否包含某个元素
	Contains(...T) bool
	//FindWithFunc 按照func 查询满足条件的T,如果满足条件返回对应的T和是否存在
	FindWithFunc(func(T) bool) (T, bool)
	//ContainsWithFunc 判断 这个T 是否算包含 x是元素,y是迭代值
	ContainsWithFunc(T, func(index int, x, y T) bool) bool
	//AllMatch 全是则是
	AllMatch(match func(T) bool) bool
	//NoneMatch 全部是则是
	NoneMatch(match func(T) bool) bool
	//AnyMatch 一个是则是
	AnyMatch(match func(T) bool) bool

	//Sort 排序,不会修改原数组,不保证原始lice位置
	Sort(func(i, j T) bool) *[]T
	//Sorted 排序,会修改原数组,不保证原始lice位置
	Sorted(func(i, j T) bool)
	//SortStable 排序,不会修改原数组,保证原始lice位置
	SortStable(func(i, j T) bool) *[]T
	//SortedStable 排序,会修改原数组,保证原始lice位置
	SortedStable(func(i, j T) bool)
	//Reverse 反转数组 返回新数组
	Reverse() *[]T
	// Reversed 反转数组自身
	Reversed()

	//Count 直接计数
	Count(T) int64
	//CountWithFunc 按照某个条件计数
	CountWithFunc(func(T) bool) int64
	//Filter 过滤
	Filter(func(index int, t T) bool) *[]T
	//Reduce 同类型Reduce
	Reduce(func(x, y T) T) *optional.Optional[T]
	//Distinct 去重
	Distinct() *[]T
	//Join 结合
	Join(sep string) string

	//Accumulate 累算器
	Accumulate(func(index int, new, old T) T) *[]T

	//Data 返回自身
	Data() *[]T
	Modify(*[]T)
	//Iter 返回迭代器 类似于python 迭代器
	Iter() <-chan T
	//Copy copy 一份新的
	Copy() *[]T
	//Clear 清空
	Clear()
	//GetIterator 获得可控制的迭代器
	GetIterator() *collection_helper.Iterator[T]
	//Each 遍历 返回false 停止
	Each(func(T) bool)
	//String string化表示
	String() string
	//Size 容器的大小
	Size() int

	IsEmpty() bool
	//LastIndex 从后向前找
	LastIndex(T) int
	//Removes 删除一些元素 返回删除了哪些
	Removes(...T) int
	RemoveAll(*[]T) int
	//Retains 返回还剩下哪些元素
	Retains(...T) int
	RetainAll(*[]T) int
	Get(int) T
	//Set 设置某个值 返回旧的值
	Set(int, T) T
	//Map 自己换成自己
	Map(func(T) T)
}

type ConvertorToList[V comparable] interface {
	ToArrayList() *arraylist.List[V]
	ToSinglyLinkedList() *singlylinkedlist.List[V]
	ToDoublyLinkedList() *doublylinkedlist.List[V]
}

type ConvertorToSet[V comparable] interface {
	ToMapSet() mapset.Set[V]
	ToThreadSafeMapSet() mapset.Set[V]
	ToHashSet() *hashset.Set[V]
	ToLinkedHashSet() *linkedhashset.Set[V]
	ToTreeSet() *treeset.Set[V]
}

type ConvertorToStack[V comparable] interface {
	TOArrayStack() *arraystack.Stack[V]
	ToLinkedListStack() *linkedliststack.Stack[V]
}
type ConvertorToTree[V comparable] interface {
	ToBinaryHeap() *binaryheap.Heap[V]
}
