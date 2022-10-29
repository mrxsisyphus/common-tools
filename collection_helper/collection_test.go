package collection_helper

import (
	"common-tools/stream"
	"common-tools/stream/types"
	"fmt"
	"math"
	"strconv"
	"testing"
)

func Test(t *testing.T) {
	people := make([]Person, 0, 10)
	for i := 0; i < 10; i++ {
		people = append(people, Person{
			id:   i,
			name: "p" + strconv.Itoa(i),
		})
	}
	fmt.Println(people)
	//go-zero方案,用channel 来做 可以并行,缺点就是 总结的时候没办法直接接受
	//https://github.com/chenquan/stream
	//first, err := stream.Of(people).
	//	Filter(func(item interface{}) bool {
	//		person := item.(Person)
	//		return person.id > 2
	//	}).Map(func(item interface{}) interface{} {
	//	person := item.(Person)
	//	return person.name
	//}).FindFirst()
	//if err != nil {
	//	return
	//}
	//fmt.Println(first)
	//github.com/youthlin/stream 的方案 只能串行,
	//res := stream.OfSlice(people).
	//	Filter(func(item types.T) bool {
	//		person := item.(Person)
	//		return person.id > 2
	//	}).Map(func(item types.T) types.R {
	//	person := item.(Person)
	//	return person.name
	//}).Map()
	//fmt.Println(res)
	// 支持泛型的方案 局限性过大,map_helper 转换值也只有一个
	//slice := stream.NewSliceByMapping[Person, string, string](people).
	//	Filter(func(p Person) bool { return p.id >= 2 }).
	//	Map(func(p Person) string { return p.name }).
	//	ToSlice()
	//fmt.Println(slice)
	res := stream.OfSlice(people).
		Filter(func(item types.T) bool {
			person := item.(Person)
			return person.id > 2
		}).Group(func(val types.T) types.R {
		person := val.(Person)
		return person.id
	})
	fmt.Println(res)

}

func TestCollectionUtils(t *testing.T) {
	people := make([]Person, 0, 10)
	for i := 10; i > 0; i-- {
		people = append(people, Person{
			id:   i,
			name: "p" + strconv.Itoa(i),
		})
	}
	fmt.Println(people)
	peoplePointer := Sort[Person](&people, func(x, y Person) bool {
		// order by id asc
		return x.id < y.id
	})
	fmt.Println(peoplePointer)
	temp := Map[Person, int](&people, func(index int, p Person) int {
		return p.id
	})
	fmt.Println(temp)

	temp = Filter[int](temp, func(index, x int) bool {
		return x <= 2
	})
	fmt.Println(temp)
	res := Reduce[int](temp, func(x, y int) int {
		return x + y
	})
	fmt.Println(res)

}

func TestCollectionUtils2(t *testing.T) {
	//test reduce
	temp := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	//累加
	sum := Reduce[int](&temp, func(x, y int) int {
		return x + y
	})
	fmt.Println(sum)

	// 累积
	multi := Reduce[int](&temp, func(x, y int) int {
		//fmt.Println(x, y)
		return x * y
	})
	fmt.Println(multi)

	// 求最大值
	max := Reduce[int](&temp, func(x, y int) int {
		return int(math.Max(float64(x), float64(y)))
	})
	fmt.Println(max)

	//求最小值
	min := Reduce[int](&temp, func(x, y int) int {
		return int(math.Min(float64(x), float64(y)))
	})
	fmt.Println(min)

	//Reduce空数组处理
	//如果直接返回会返回0值 这不是我们想要的
	var empty []int
	res := Reduce[int](&empty, func(x, y int) int {
		return x + y
	})
	//使用optional处理
	get, err := res.Get()
	if err != nil {
		panic(err)
	}
	fmt.Println(get)

}

func TestSort(t *testing.T) {
	people := make([]Person, 0, 10)
	for i := 10; i > 0; i-- {
		people = append(people, Person{
			id:   i,
			name: "p" + strconv.Itoa(i),
			age:  10 + (i % 2),
		})
	}
	fmt.Printf("%+v\n", people)
	//按照用户id正序排序,不会更改原people
	sortPeople1 := Sort[Person](&people, func(x, y Person) bool {
		return x.id < y.id
	})
	fmt.Printf("%v\n", *sortPeople1)
	fmt.Printf("%v\n", people)
	//按照id倒序排序,会更改原people
	Sorted[Person](sortPeople1, func(x, y Person) bool {
		return x.id > y.id
	})
	fmt.Printf("%v\n", *sortPeople1)
	fmt.Printf("%v\n", people)
	// 多阶段排序 order by age desc,id asc
	sortPeople2 := Sort[Person](&people, func(x, y Person) bool {
		// order by age desc
		if x.age != y.age {
			return x.age > y.age
		}
		//order by id asc 当x.age == y.age的时候
		if x.id != y.id {
			return x.id < y.id
		}
		return false
	})
	fmt.Printf("%v\n", *sortPeople2)

}

func TestDistinct(t *testing.T) {
	people := make([]Person, 0, 10)
	for i := 10; i > 0; i-- {
		people = append(people, Person{
			id:   i,
			name: "p" + strconv.Itoa(i),
			age:  10 + (i % 2),
		})
	}
	fmt.Printf("%+v\n", people)
	//Set 直接把people去重(people对象直接hash)
	peopleDistinct := Distinct[Person](&people)
	fmt.Printf("%+v\n", *peopleDistinct)
}

func TestGenerateRange(t *testing.T) {
	ran := IntRange(1, 80, 1)

	fmt.Println(len(*ran), cap(*ran), ran)
	ran2 := IntRandRange(0, 80)
	fmt.Println(len(*ran2), cap(*ran2), ran2)
}

func TestAccumulate(t *testing.T) {
	temp := []int{1, 2, 3, 4, 5, 6}
	res := AccumulateWithSameType(&temp, func(index, x, y int) int {
		return x + y
	})
	fmt.Println(res)
	temp2 := []float64{1.0, 2.1, 3.0, 4.0, 5.0, 6.0}
	res2 := AccumulateWithSameType[float64](&temp2, func(index int, x, y float64) float64 {
		if index == 0 {
			return y
		}
		return x + y
	})
	fmt.Println(res2)

}

func TestBatchRemove(t *testing.T) {
	d1 := []int{1, 2, 3, 4, 5}
	d2 := []int{1, 2}
	PrintSlice(&d1, "d1")
	PrintSlice(&d2, "d2")
	////d1 删除 d2
	//fmt.Println(collection_helper.BatchRemove(&d1, &d2, false))
	//collection_helper.PrintSlice(&d1, "d1")
	//collection_helper.PrintSlice(&d2, "d2")
	//d1 保留 d2
	fmt.Println(BatchRemove(&d1, &d2, true))
	//fmt.Println(d1[3])
	PrintSlice(&d1, "d1")
	PrintSlice(&d2, "d2")
}
