package random_helper

import (
	"fmt"
	"github.com/mrxtryagin/common-tools/collection_helper"
	"github.com/mrxtryagin/common-tools/string_helper"
	"math"
	"math/rand"
	"strconv"
	"testing"
	"time"
	"unsafe"
)

func TestRandom(t *testing.T) {
	fmt.Println(RandInt(33, 80))
}

type Person struct {
	name string
	id   int
	age  int
}

func TestChoice(t *testing.T) {
	// 初始化全局seed
	InitSeed()
	people := make([]Person, 0, 1000)
	for i := 1000; i > 0; i-- {
		people = append(people, Person{
			id:   i,
			name: "p" + strconv.Itoa(i),
			age:  10 + (i % 2),
		})
	}
	fmt.Printf("%+v\n", people)
	rand := collection_helper.Range(100)
	fmt.Println(rand)
	// 随机选一个
	for i := 0; i < 10; i++ {
		//随机选一个
		choice := Choice(rand)
		fmt.Println(choice)
	}
	// sample people
	res := Sample(&people, 1)
	fmt.Printf("%+v\n", res)
	isDistinct := collection_helper.JudgeSliceIsDistinctByFunc[Person, int](res, func(person Person) int {
		return person.id
	})
	fmt.Printf("%+v\n", isDistinct)
	// 不会改变原来的序列
	//fmt.Infof("%+v\n", people)

}

func TestGenerateRandomNumberWithoutRepeat(t *testing.T) {
	res := GenerateRandomNumberWithoutRepeat(100, 10000, 1000)
	fmt.Println(res)
}

func TestSizeOf(t *testing.T) {
	sizeofMap := unsafe.Sizeof(map[int]bool{})
	fmt.Println(sizeofMap)
	sizeofArr := unsafe.Sizeof([]int{})
	fmt.Println(sizeofArr)

}

func TestBinarySearch(t *testing.T) {
	res := collection_helper.Range(10)
	fmt.Println(res)
	res2 := BinaryTest(*res, 9)
	fmt.Println(res2)
}

func BinaryTest(els []int, target int) int {
	res := -1
	start := 0
	end := len(els) - 1
	for start <= end {
		mid := (start + end) >> 1
		if target < els[mid] {
			//左半边,end 就改成中间的值
			end = mid - 1
		} else if target == els[mid] {
			return mid
		} else {
			//右半边
			start = mid + 1
		}
	}
	return res
}

func TestDemo(t *testing.T) {
	rand.Seed(1)
	for i := 0; i < 10000; i++ {
		fmt.Println(rand.Intn(10000))
	}
}

func TestFloat64n(t *testing.T) {
	res := Float64n(0.125)
	fmt.Println(math.Floor(res))
}

func TestChoices(t *testing.T) {
	InitSeed()
	a := []string{"A", "B", "C"}
	weight := []float64{1.0, 2.0, 7.0}
	reverseWeight := collection_helper.Reverse[float64](&weight)
	fmt.Println(reverseWeight)
	choice := Choices(&a, *reverseWeight, nil, 100000000)
	//fmt.Println(*choice)
	fmt.Println(len(*choice))
	res := collection_helper.ToGroupCount[string, string](choice, func(inputType string) string {
		return inputType
	})
	fmt.Println(res)
}

func TestChoices2(t *testing.T) {
	InitSeed()
	a := [][]int{[]int{1, 2, 3}, []int{2, 3, 4}, []int{3, 4, 5}}
	weight := []float64{1.0, 2.0, 7.0}
	reverseWeight := collection_helper.Reverse[float64](&weight)
	fmt.Println(reverseWeight)
	now := time.Now()
	choice := Choices(&a, *reverseWeight, nil, 100000000)
	fmt.Println(time.Since(now).Milliseconds())
	//fmt.Println(*choice)
	fmt.Println(len(*choice))
	res := collection_helper.ToGroupCount[[]int, string](choice, func(inputType []int) string {
		return string_helper.ToString(inputType)
	})
	fmt.Println(res)
}

type Prize struct {
	PlayerId int64
	Weight   int
}

func TestChoices3(t *testing.T) {
	InitSeed()
	prizes := make([]*Prize, 4, 4)
	for i := 0; i < 4; i++ {
		prize := &Prize{
			PlayerId: int64(i) + 10000,
			Weight:   (5 - i) * 10,
		}
		prizes[i] = prize
	}
	// 50,40,30,20
	fmt.Println(prizes)
	total := collection_helper.ReduceWithInitValue[*Prize, float64](&prizes, func(x *Prize, val float64) float64 {
		return float64(x.Weight) + val
	}, 0)
	for _, prize := range prizes {
		radio := float64(prize.Weight) / total
		fmt.Printf("key: %v, val: %v, radio:%v, reverse:%v\n", prize.PlayerId, prize.Weight, radio, 1-radio)
	}
	mark := 10000
	res := ChoicesByFunc[*Prize](&prizes, func(t *Prize) float64 {
		return float64(t.Weight)
	}, func(i *[]float64) *[]float64 {
		// 逆转各个的概率
		collection_helper.Reversed(i)
		return i
	}, mark)
	res2 := collection_helper.ToGroupCount[*Prize, int64](res, func(inputType *Prize) int64 {
		return inputType.PlayerId
	})
	for key, val := range res2 {
		fmt.Printf("key: %v, val: %v, radio:%v\n", key, val, float64(val)/float64(mark))
	}
	fmt.Println(res)

}
