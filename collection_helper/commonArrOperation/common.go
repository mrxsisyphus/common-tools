package commonArrOperation

import (
	"fmt"
	"github.com/mrxtryagain/common-tools/collection_helper"
	"github.com/mrxtryagain/common-tools/optional"
)

/**
数组的
增 删 改 查
*/

// 这里支持的T的维度不到compare 因为有些数组可以装一些可以不compare的
type ArrOperation[T any] struct {
	input *[]T
}

func (a *ArrOperation[T]) Add(t T) {
	// refer:https://segmentfault.com/a/1190000039943498
	// a.input 是切片的指针 他指向切片 切片中有个指针 指向底层数组
	// append 后 底层数组扩容,产生新的切片 旧的切片 和新的切片一起指向新的数组
	// 那 旧切片的指针当然也没问题
	*a.input = append(*a.input, t)

}

func (a *ArrOperation[T]) Adds(ts ...T) {
	a.AddAll(&ts)
}

func (a *ArrOperation[T]) AddAll(ts *[]T) {
	*a.input = append(*a.input, *ts...)
	//i := *a.input
	//i = append(i, *ts...)
	//a.input = &i

}

func (a *ArrOperation[T]) Remove(t T) {
	index := a.Index(t)
	if index == -1 {
		return
	}
	a.Pop(index)
}

func (a *ArrOperation[T]) RemoveWithFunc(fn func(t T) bool) {
	for i, t := range *a.input {
		if fn(t) {
			a.Pop(i)
		}
	}
}

// DeleteRange 参考 slices.Delete
func (a *ArrOperation[T]) DeleteRange(i, j int) {
	*a.input = append((*a.input)[:i], (*a.input)[j:]...)
}

func (a *ArrOperation[T]) Pop(index int) T {
	i := *a.input
	poped := i[index]
	*a.input = append(i[:index], i[index+1:]...)
	return poped
}

func (a *ArrOperation[T]) PopFirst() T {
	return a.Pop(0)
}

func (a *ArrOperation[T]) PopLast() T {
	return a.Pop(len(*a.input) - 1)
}

// Insert  参考slices.Insert
func (a *ArrOperation[T]) Insert(i int, v ...T) {
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

func (a *ArrOperation[T]) Index(t T) int {
	// https://stackoverflow.com/questions/68053957/go-with-generics-type-parameter-t-is-not-comparable-with 使用反射的方式
	for i, vs := range *a.input {
		// '手动类型擦除'
		// refer: https://segmentfault.com/a/1190000041634906
		if any(t) == any(vs) {
			return i
		}
	}
	return -1
}

func (a *ArrOperation[T]) IndexWithFunc(f func(T) bool) int {
	for i, vs := range *a.input {
		if f(vs) {
			return i
		}
	}
	return -1
}

func (a *ArrOperation[T]) Contains(ts ...T) bool {
	for _, t := range ts {
		index := a.Index(t)
		if index == -1 {
			return false
		}
	}
	return true
}

func (a *ArrOperation[T]) FindWithFunc(f func(T) bool) (t T, exist bool) {
	for _, vs := range *a.input {
		if f(vs) {
			return vs, true
		}
	}
	return
}

func (a *ArrOperation[T]) ContainsWithFunc(t T, f func(index int, x, y T) bool) bool {
	return collection_helper.ContainsWithFunc[T](a.input, t, f)
}

func (a *ArrOperation[T]) AllMatch(match func(T) bool) bool {
	return collection_helper.AllMatch[T](a.input, match)
}

func (a *ArrOperation[T]) NoneMatch(match func(T) bool) bool {
	return collection_helper.NoneMatch[T](a.input, match)
}

func (a *ArrOperation[T]) AnyMatch(match func(T) bool) bool {
	return collection_helper.AnyMatch[T](a.input, match)
}

func (a *ArrOperation[T]) Sort(f func(i T, j T) bool) *[]T {
	return collection_helper.Sort[T](a.input, f)
}

func (a *ArrOperation[T]) Sorted(f func(i T, j T) bool) {
	collection_helper.Sorted[T](a.input, f)
}

func (a *ArrOperation[T]) SortStable(f func(i T, j T) bool) *[]T {
	return collection_helper.SortStable[T](a.input, f)
}

func (a *ArrOperation[T]) SortedStable(f func(i T, j T) bool) {
	collection_helper.SortStable[T](a.input, f)
}

func (a *ArrOperation[T]) Reverse() *[]T {
	return collection_helper.Reverse[T](a.input)
}

func (a *ArrOperation[T]) Reversed() {
	collection_helper.Reversed[T](a.input)
}

func (a *ArrOperation[T]) Count(t T) int64 {
	var total int64
	for _, t2 := range *a.input {
		if any(t2) == any(t) {
			total++
		}
	}
	return total
}

func (a *ArrOperation[T]) CountWithFunc(f func(T) bool) int64 {
	return collection_helper.CountSliceByFunc[T](a.input, f)
}

func (a *ArrOperation[T]) Filter(f func(index int, t T) bool) *[]T {
	return collection_helper.Filter[T](a.input, f)
}

func (a *ArrOperation[T]) Reduce(f func(x T, y T) T) *optional.Optional[T] {
	return collection_helper.Reduce[T](a.input, f)
}

func (a *ArrOperation[T]) Distinct() *[]T {
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

func (a *ArrOperation[T]) Join(sep string) string {
	return collection_helper.Join[T](a.input, sep)
}

func (a *ArrOperation[T]) Accumulate(f func(index int, new T, old T) T) *[]T {
	return collection_helper.AccumulateWithSameType[T](a.input, f)
}

func (a *ArrOperation[T]) Data() *[]T {
	return a.input
}

func (a *ArrOperation[T]) Modify(ts *[]T) {
	a.input = ts
}

func (a *ArrOperation[T]) Iter() <-chan T {
	return collection_helper.Iter[T](a.input)
}

func (a *ArrOperation[T]) Copy() *[]T {
	i := *a.input
	res := make([]T, len(i))
	copy(res, i)
	return &res
}

func (a *ArrOperation[T]) Clear() {
	// 可以使用nil ,但是 nil 序列化成json 时会变为nil
	// refer:https://islishude.github.io/blog/2020/09/22/golang/Go-%E6%B8%85%E7%A9%BA-Slice-%E7%9A%84%E4%B8%A4%E7%A7%8D%E6%96%B9%E6%B3%95%EF%BC%9A-0-%E5%92%8Cnil/
	*a.input = (*a.input)[:0]
}

func (a *ArrOperation[T]) GetIterator() *collection_helper.Iterator[T] {
	return collection_helper.GetIterator[T](a.input)
}

func (a *ArrOperation[T]) Each(f func(T) bool) {
	for _, t := range *a.input {
		if f(t) {
			break
		}
	}
}

func (a *ArrOperation[T]) String() string {

	return fmt.Sprintf("[%s]", a.Join(","))
}

func (a *ArrOperation[T]) Size() int {
	return len(*a.input)
}

func (a *ArrOperation[T]) IsEmpty() bool {
	return a.Size() == 0
}

func (a *ArrOperation[T]) LastIndex(t T) int {
	raw := *a.input
	for i := len(raw) - 1; i >= 0; i-- {
		if any(raw[i]) == any(t) {
			return i
		}
	}
	return -1
}

func (a *ArrOperation[T]) Removes(ts ...T) int {
	return collection_helper.BatchRemove[T](a.input, &ts, false)
}

func (a *ArrOperation[T]) RemoveAll(ts *[]T) int {
	return collection_helper.BatchRemove[T](a.input, ts, false)
}

func (a *ArrOperation[T]) Retains(ts ...T) int {
	return collection_helper.BatchRemove[T](a.input, &ts, true)
}

func (a *ArrOperation[T]) RetainAll(ts *[]T) int {
	return collection_helper.BatchRemove[T](a.input, ts, true)
}

func (a *ArrOperation[T]) Get(index int) T {
	return (*a.input)[index]
}

func (a *ArrOperation[T]) Set(index int, t T) T {
	raw := *a.input
	old := raw[index]
	raw[index] = t
	return old
}

func (a *ArrOperation[T]) Map(fn func(T) T) {
	raw := *a.input
	for i := 0; i < len(raw); i++ {
		raw[i] = fn(raw[i])
	}
}

// NewArrOperation 从输入的数组 开始造一个数组操作类
func NewArrOperation[T any](input *[]T) *ArrOperation[T] {
	return &ArrOperation[T]{
		input: input,
	}
}

// NewEmptyArrOperation 直接一个空的数组操作类
func NewEmptyArrOperation[T any]() *ArrOperation[T] {
	return &ArrOperation[T]{
		input: new([]T),
	}
}

// NewArrOperationWithEls 用一些元素来初始化数组操作类
func NewArrOperationWithEls[T any](ts ...T) *ArrOperation[T] {
	return &ArrOperation[T]{
		input: &ts,
	}
}
