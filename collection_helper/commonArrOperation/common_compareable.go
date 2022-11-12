package commonArrOperation

import (
	"fmt"
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
	"github.com/monitor1379/yagods/utils"
	"github.com/mrxtryagain/common-tools/collection_helper"
	"github.com/mrxtryagain/common-tools/optional"
)

/**
数组的
增 删 改 查
*/

// ArrComparableOperation 这里支持的T的维度不到compare 因为有些数组可以装一些可以不compare的
type ArrComparableOperation[T comparable] struct {
	input *[]T
}

func (a *ArrComparableOperation[T]) ToArrayList() *arraylist.List[T] {
	return arraylist.New[T](*a.input...)
}

func (a *ArrComparableOperation[T]) ToSinglyLinkedList() *singlylinkedlist.List[T] {
	return singlylinkedlist.New[T](*a.input...)
}

func (a *ArrComparableOperation[T]) ToDoublyLinkedList() *doublylinkedlist.List[T] {
	return doublylinkedlist.New[T](*a.input...)
}

func (a *ArrComparableOperation[T]) ToMapSet() mapset.Set[T] {
	return mapset.NewThreadUnsafeSet[T](*a.input...)
}
func (a *ArrComparableOperation[T]) ToThreadSafeMapSet() mapset.Set[T] {
	return mapset.NewSet[T](*a.input...)
}

func (a *ArrComparableOperation[T]) ToHashSet() *hashset.Set[T] {
	return hashset.New[T](*a.input...)
}

func (a *ArrComparableOperation[T]) ToLinkedHashSet() *linkedhashset.Set[T] {
	return linkedhashset.New[T](*a.input...)
}

func (a *ArrComparableOperation[T]) ToTreeSet(comparator utils.Comparator[T]) *treeset.Set[T] {
	return treeset.NewWith[T](comparator, *a.input...)
}

func (a *ArrComparableOperation[T]) TOArrayStack() *arraystack.Stack[T] {
	s := arraystack.New[T]()
	for _, t := range *a.input {
		s.Push(t)
	}
	return s
}

func (a *ArrComparableOperation[T]) ToLinkedListStack() *linkedliststack.Stack[T] {
	s := linkedliststack.New[T]()
	for _, t := range *a.input {
		s.Push(t)
	}
	return s
}

func (a *ArrComparableOperation[T]) ToBinaryHeap(comparator utils.Comparator[T]) *binaryheap.Heap[T] {
	h := binaryheap.NewWith[T](comparator)
	for _, t := range *a.input {
		h.Push(t)
	}
	return h
}

func (a *ArrComparableOperation[T]) Add(t T) {
	// refer:https://segmentfault.com/a/1190000039943498
	// a.input 是切片的指针 他指向切片 切片中有个指针 指向底层数组
	// append 后 底层数组扩容,产生新的切片 旧的切片 和新的切片一起指向新的数组
	// 那 旧切片的指针当然也没问题
	*a.input = append(*a.input, t)

}

func (a *ArrComparableOperation[T]) Adds(ts ...T) {
	a.AddAll(&ts)
}

func (a *ArrComparableOperation[T]) AddAll(ts *[]T) {
	*a.input = append(*a.input, *ts...)
	//i := *a.input
	//i = append(i, *ts...)
	//a.input = &i

}

func (a *ArrComparableOperation[T]) Remove(t T) {
	index := a.Index(t)
	if index == -1 {
		return
	}
	a.Pop(index)
}

func (a *ArrComparableOperation[T]) RemoveWithFunc(fn func(t T) bool) {
	for i, t := range *a.input {
		if fn(t) {
			a.Pop(i)
		}
	}
}

// DeleteRange 参考 slices.Delete
func (a *ArrComparableOperation[T]) DeleteRange(i, j int) {
	*a.input = append((*a.input)[:i], (*a.input)[j:]...)
}

func (a *ArrComparableOperation[T]) Pop(index int) T {
	i := *a.input
	poped := i[index]
	*a.input = append(i[:index], i[index+1:]...)
	return poped
}

func (a *ArrComparableOperation[T]) PopFirst() T {
	return a.Pop(0)
}

func (a *ArrComparableOperation[T]) PopLast() T {
	return a.Pop(len(*a.input) - 1)
}

// Insert  参考slices.Insert
func (a *ArrComparableOperation[T]) Insert(i int, v ...T) {
	s := *a.input
	tot := len(s) + len(v)
	if tot <= cap(s) {
		s2 := s[:tot]              //s2的切片容量
		copy(s2[i+len(v):], s[i:]) //插入点之后v长度之后的copy
		copy(s2[i:], v)            // 插入点之v长度的copy
		*a.input = s2
		return
	}
	// new个新的
	s2 := make([]T, tot)
	copy(s2, s[:i])            //插入点之前的copy
	copy(s2[i:], v)            // 插入点之v长度的copy
	copy(s2[i+len(v):], s[i:]) //插入点之后v长度之后的copy
	*a.input = s2
	return
}

func (a *ArrComparableOperation[T]) Index(t T) int {
	for i, vs := range *a.input {
		if t == vs {
			return i
		}
	}
	return -1
}

func (a *ArrComparableOperation[T]) IndexWithFunc(f func(T) bool) int {
	for i, vs := range *a.input {
		if f(vs) {
			return i
		}
	}
	return -1
}

func (a *ArrComparableOperation[T]) Contains(ts ...T) bool {
	for _, t := range ts {
		index := a.Index(t)
		if index == -1 {
			return false
		}
	}
	return true
}

func (a *ArrComparableOperation[T]) FindWithFunc(f func(T) bool) (t T, exist bool) {
	for _, vs := range *a.input {
		if f(vs) {
			return vs, true
		}
	}
	return
}

func (a *ArrComparableOperation[T]) ContainsWithFunc(t T, f func(index int, x, y T) bool) bool {
	return collection_helper.ContainsWithFunc[T](a.input, t, f)
}

func (a *ArrComparableOperation[T]) AllMatch(match func(T) bool) bool {
	return collection_helper.AllMatch[T](a.input, match)
}

func (a *ArrComparableOperation[T]) NoneMatch(match func(T) bool) bool {
	return collection_helper.NoneMatch[T](a.input, match)
}

func (a *ArrComparableOperation[T]) AnyMatch(match func(T) bool) bool {
	return collection_helper.AnyMatch[T](a.input, match)
}

func (a *ArrComparableOperation[T]) Sort(f func(i T, j T) bool) *[]T {
	return collection_helper.Sort[T](a.input, f)
}

func (a *ArrComparableOperation[T]) Sorted(f func(i T, j T) bool) {
	collection_helper.Sorted[T](a.input, f)
}

func (a *ArrComparableOperation[T]) SortStable(f func(i T, j T) bool) *[]T {
	return collection_helper.SortStable[T](a.input, f)
}

func (a *ArrComparableOperation[T]) SortedStable(f func(i T, j T) bool) {
	collection_helper.SortStable[T](a.input, f)
}

func (a *ArrComparableOperation[T]) Reverse() *[]T {
	return collection_helper.Reverse[T](a.input)
}

func (a *ArrComparableOperation[T]) Reversed() {
	collection_helper.Reversed[T](a.input)
}

func (a *ArrComparableOperation[T]) Count(t T) int64 {
	var total int64
	for _, t2 := range *a.input {
		if t2 == t {
			total++
		}
	}
	return total
}

func (a *ArrComparableOperation[T]) CountWithFunc(f func(T) bool) int64 {
	return collection_helper.CountSliceByFunc[T](a.input, f)
}

func (a *ArrComparableOperation[T]) Filter(f func(index int, t T) bool) *[]T {
	return collection_helper.Filter[T](a.input, f)
}

func (a *ArrComparableOperation[T]) Reduce(f func(x T, y T) T) *optional.Optional[T] {
	return collection_helper.Reduce[T](a.input, f)
}

func (a *ArrComparableOperation[T]) Distinct() *[]T {
	i := *a.input
	if len(i) <= 0 {
		return new([]T)
	}
	distinctMap := make(map[any]struct{})
	res := make([]T, 0)
	for _, item := range i {
		if _, ok := distinctMap[item]; !ok {
			//不存在 放入
			distinctMap[item] = struct{}{}
			res = append(res, item)
		}
	}
	return &res
}

func (a *ArrComparableOperation[T]) Join(sep string) string {
	return collection_helper.Join[T](a.input, sep)
}

func (a *ArrComparableOperation[T]) Accumulate(f func(index int, new T, old T) T) *[]T {
	return collection_helper.AccumulateWithSameType[T](a.input, f)
}

func (a *ArrComparableOperation[T]) Data() *[]T {
	return a.input
}

func (a *ArrComparableOperation[T]) Modify(ts *[]T) {
	a.input = ts
}

func (a *ArrComparableOperation[T]) Iter() <-chan T {
	return collection_helper.Iter[T](a.input)
}

func (a *ArrComparableOperation[T]) Copy() *[]T {
	i := *a.input
	res := make([]T, len(i))
	copy(res, i)
	return &res
}

func (a *ArrComparableOperation[T]) Clear() {
	// 可以使用nil ,但是 nil 序列化成json 时会变为nil
	// refer:https://islishude.github.io/blog/2020/09/22/golang/Go-%E6%B8%85%E7%A9%BA-Slice-%E7%9A%84%E4%B8%A4%E7%A7%8D%E6%96%B9%E6%B3%95%EF%BC%9A-0-%E5%92%8Cnil/
	*a.input = (*a.input)[:0]
}

func (a *ArrComparableOperation[T]) GetIterator() *collection_helper.Iterator[T] {
	return collection_helper.GetIterator[T](a.input)
}

func (a *ArrComparableOperation[T]) Each(f func(T) bool) {
	for _, t := range *a.input {
		if f(t) {
			break
		}
	}
}

func (a *ArrComparableOperation[T]) String() string {

	return fmt.Sprintf("[%s]", a.Join(","))
}

func (a *ArrComparableOperation[T]) Size() int64 {
	return int64(len(*a.input))
}

func (a *ArrComparableOperation[T]) IsEmpty() bool {
	return a.Size() == 0
}

func (a *ArrComparableOperation[T]) LastIndex(t T) int {
	raw := *a.input
	for i := len(raw) - 1; i >= 0; i-- {
		if any(raw[i]) == any(t) {
			return i
		}
	}
	return -1
}

func (a *ArrComparableOperation[T]) Removes(ts ...T) int {
	return collection_helper.BatchRemove[T](a.input, &ts, false)
}

func (a *ArrComparableOperation[T]) RemoveAll(ts *[]T) int {
	return collection_helper.BatchRemove[T](a.input, ts, false)
}

func (a *ArrComparableOperation[T]) Retains(ts ...T) int {
	return collection_helper.BatchRemove[T](a.input, &ts, true)
}

func (a *ArrComparableOperation[T]) RetainAll(ts *[]T) int {
	return collection_helper.BatchRemove[T](a.input, ts, true)
}

func (a *ArrComparableOperation[T]) Get(index int) T {
	return (*a.input)[index]
}

func (a *ArrComparableOperation[T]) Set(index int, t T) T {
	raw := *a.input
	old := raw[index]
	raw[index] = t
	return old
}

// NewArrComparableOperation 从输入的数组 开始造一个数组操作类
func NewArrComparableOperation[T comparable](input *[]T) *ArrComparableOperation[T] {
	return &ArrComparableOperation[T]{
		input: input,
	}
}

// NewEmptyArrComparableOperation 直接一个空的数组操作类
func NewEmptyArrComparableOperation[T comparable]() *ArrComparableOperation[T] {
	return &ArrComparableOperation[T]{
		input: new([]T),
	}
}

// NewArrComparableOperationWithEls 用一些元素来初始化数组操作类
func NewArrComparableOperationWithEls[T comparable](ts ...T) *ArrComparableOperation[T] {
	return &ArrComparableOperation[T]{
		input: &ts,
	}
}
