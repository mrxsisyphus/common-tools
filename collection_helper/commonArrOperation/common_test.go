package commonArrOperation

import (
	"common-tools/collection_helper"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"unsafe"
)

type Person struct {
	name string
	id   int
	age  int
}

func TestArrOperation2(t *testing.T) {
	a := NewEmptyArrOperation[int]()
	fmt.Println(a)
	a.Adds(1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16)
	//a.Add(1)
	fmt.Println(a)

}

// TestDelete 所以我们通过%p打印的slice变量ages的地址其实就是内部存储数组元素的地址
// refer:https://blog.csdn.net/u012279631/article/details/80271408
// refer: https://halfrost.com/go_slice/
func TestDelete(t *testing.T) {
	//旧切片
	old := []int{1, 2, 3, 4, 5, 6, 7}
	//旧切片的地址
	fmt.Println(unsafe.Pointer(&old)) //0x1400012c090
	//旧切片 指向的底层数组的信息
	fmt.Printf("old 底层数组: %v, len: %d cap: %d ,pointer: %p\n", old, len(old), cap(old), old) //0x1400011e1c0
	//fmt.Infof("old底层数组地址:%d", getSliceOrigianalAddress(old))
	//对切片做删除 删除index =2的元素
	// 分成前半后半,它们与old一样 指向底层数组(软拷贝)
	//前半的切片地址
	before := old[:2]
	//前半的切片地址
	fmt.Println(unsafe.Pointer(&before))                                                                    //0x1400012c0d8
	fmt.Printf("before 底层数组: %v, len: %d cap: %d ,pointer: %p\n", before, len(before), cap(before), before) //0x1400011e1c0 与旧的一致
	//fmt.Infof("before底层数组地址:%d", getSliceOrigianalAddress(before))
	// refer: https://zhuanlan.zhihu.com/p/526731603
	// 直接切片可能会触发产生新的底层数组(因为cap不一样)
	//新切片的 容量(cap) 即为开始下标到原 slice 数据容量结束, 即"cap(a) - low"
	after := old[3:]
	after[0] = 10  // 改变这个地方确实会映射到old上 也就是说还是会 改变的?
	full := old[:] // 改变cap大小就不行
	fmt.Println(unsafe.Pointer(&full))
	fmt.Printf("full 底层数组: %v, len: %d cap: %d ,pointer: %p\n", full, len(full), cap(full), full) //0x1400011e1c0 与旧的一致
	//后半的切片地址
	fmt.Println(unsafe.Pointer(&after))
	fmt.Printf("after 底层数组: %v, len: %d cap: %d ,pointer: %p\n", after, len(after), cap(after), after) //0x1400011e1d8 后半的底层数组已经发生了改变? 还是同一个数组
	//fmt.Infof("after底层数组地址:%d", getSliceOrigianalAddress(after))
	// append 会触发啥? 这个时候 cap 没变,所以底层数组会不会变呢?
	newAppend := append(before, after...)
	// newAppend切片地址 这个地址变了吗?
	fmt.Println(unsafe.Pointer(&newAppend))                                                                                //0x1400012c168
	fmt.Printf("newAppend 底层数组: %v, len: %d cap: %d ,pointer: %p\n", newAppend, len(newAppend), cap(newAppend), newAppend) //0x1400011e1c0 与旧的一致
	//原切片的地址
	fmt.Println(unsafe.Pointer(&old)) //0x1400012c090
	//原切片 指向的底层数组的信息
	fmt.Printf("old 底层数组: %v, len: %d cap: %d ,pointer: %p\n", old, len(old), cap(old), old) //0x1400011e1c0 与旧的一致
	// 后面的部分参考TestAdd

}

func TestAdd(t *testing.T) {
	// 旧切片
	old := []int{1, 2, 3, 4, 5, 6, 7}
	// 切片地址
	fmt.Println(unsafe.Pointer(&old)) //0x1400012c090
	//旧切片 指向的底层数组
	fmt.Printf("%p\n", old) //0x1400011e080
	// append 触发的扩容,会产生指向新的底层数组的新的切片
	newAppend := append(old, 8)
	// newAppend切片地址
	fmt.Println(unsafe.Pointer(&newAppend)) //0x1400012c0c0 因为扩容 产生了新的指向数组的切片
	//旧切片 指向的底层数组
	fmt.Printf("%p\n", old) //0x1400011e080 不变
	//新切片指向的底层数组
	fmt.Printf("%p\n", newAppend) //0x1400015e0e0  变成了新的数组地址了
	// 就切片和新切片一起指向 新的底层数组
	old = newAppend
	// 旧的指针
	fmt.Println(unsafe.Pointer(&old)) //0x1400012c090 依然是就切片的地址值 只是指向的地方不对了
	//旧切片 指向的底层数组
	fmt.Printf("%p\n", old) //0x1400015e0e0 变成了 新的数组地址了

}

func getSliceOrigianalAddress[T any](input []T) uintptr {
	return reflect.ValueOf(input).Pointer()
}

func TestSearch(t *testing.T) {
	a := NewEmptyArrOperation[[]int]()
	fmt.Println(a)
	a.Adds([]int{1, 2, 3}, []int{2, 3, 4})
	fmt.Println(a)
	//indexd := a.Index([]int{1, 2, 3})
	//fmt.Println(indexd)
	fmt.Println(a.PopLast())
	fmt.Println(a.PopLast())
	fmt.Println(a.PopLast())
	fmt.Println(a)

}

func TestArrOperationFull(t *testing.T) {
	res := collection_helper.Maker[Person](100, func(index int) Person {
		return Person{
			id:   index,
			name: "p" + strconv.Itoa(index),
		}
	})
	fmt.Println(res)
	ao := NewArrOperation(res)

	ao.Add(Person{
		id:   -1,
		name: "p-1",
	})
	ao.Adds(Person{
		id:   -2,
		name: "p-2",
	}, Person{
		id:   -3,
		name: "p-3",
	})
	fmt.Println(ao)
	ao.RemoveWithFunc(func(person Person) bool {
		return person.id == -3
	})
	fmt.Println(ao)
	ao.DeleteRange(10, 30)
	fmt.Println(ao, cap(*ao.Data()))
	ao.Insert(10, Person{
		id:   20,
		name: "p-20",
	}, Person{
		id:   22,
		name: "p-22",
	})
	fmt.Println(ao, cap(*ao.Data()))
	fmt.Println(ao.IndexWithFunc(func(p Person) bool {
		return p.id == -1
	}))
	fmt.Println(ao.FindWithFunc(func(p Person) bool {
		return p.id == -1
	}))
	fmt.Println(ao.AnyMatch(func(p Person) bool {
		return p.id == -1
	}))
	fmt.Println(ao.Sort(func(i, j Person) bool {
		return i.id < j.id
	}))
	fmt.Println(ao.SortStable(func(i, j Person) bool {
		return i.id > j.id
	}))

	fmt.Println(ao.Reverse())
	fmt.Println(ao.Count(Person{
		name: "",
		age:  0,
	}))
	fmt.Println(ao.CountWithFunc(func(person Person) bool {
		return person.id > 1
	}))
	fmt.Println(ao.Filter(func(index int, person Person) bool {
		return person.id < 0
	}))

	fmt.Println(ao.Reduce(func(p1, p2 Person) Person {
		return Person{
			id:   p1.id + p2.id,
			name: strings.Join([]string{p1.name, p2.name}, "_"),
		}
	}))
	fmt.Println(ao.Distinct())
	fmt.Println(ao.Size())
}

func TestBatchRemove2(t *testing.T) {
	d1 := []int{1, 2, 3, 4, 5}
	d2 := []int{}
	collection_helper.PrintSlice(&d1, "d1")
	collection_helper.PrintSlice(&d2, "d2")
	//d1 删除 d2
	o := NewArrOperation(&d1)

	//fmt.Println(o.RemoveAll(&d2))
	//collection_helper.PrintSlice(&d1, "d1")
	//collection_helper.PrintSlice(&d2, "d2")
	////d1 保留 d2
	fmt.Println(o.RemoveAll(&d2))
	collection_helper.PrintSlice(&d1, "d1")
	collection_helper.PrintSlice(&d2, "d2")
}
