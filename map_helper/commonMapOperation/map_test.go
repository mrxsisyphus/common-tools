package commonMapOperation

import (
	"fmt"
	"github.com/mrxtryagin/common-tools/collection_helper"
	"github.com/mrxtryagin/common-tools/collection_helper/commonArrOperation"
	"github.com/mrxtryagin/common-tools/random_helper"
	"strconv"
	"testing"
)

type Person struct {
	id    int
	name  string
	score int
	group int
}

func Test1(t *testing.T) {
	a := map[int]any{}
	a[1] = nil
	val, exist := a[1]
	fmt.Println(val, exist)
	fmt.Println(a)
}

func Test2(t *testing.T) {
	a := map[int]Person{}
	//a[1] = []int{1, 2, 3}
	fmt.Println(a[2])

}

func TestMapOperation_Merge(t *testing.T) {
	res := collection_helper.Maker[*Person](1000, func(index int) *Person {
		return &Person{
			id:    index,
			name:  "p" + strconv.Itoa(index),
			score: index + 1,
			group: random_helper.RandInt(index, 1000),
		}
	})
	//fmt.Println(res)
	mo := NewEmptyMapOperation[int, int]()
	o := commonArrOperation.NewArrComparableOperation(res)
	// 使用Each + merge 实现 求和
	o.Each(func(person *Person) bool {
		mo.Merge(person.group, person.score, func(oldScore, newScore int) int {
			return oldScore + newScore
		})
		return false
	})
	fmt.Println(mo.ToJsonString())
	mo.Each(func(k int, v int) bool {
		fmt.Println(k, v)
		return false
	})
	moIt := mo.Iter()
	for m := range moIt {
		// 对于map 来说这里就是Key V
		fmt.Println(m.First, m.Second)
	}
	fmt.Println("==========")
	moIt2 := mo.GetIterator()
	for m := range moIt2.C {
		if m.Second > 100 {
			moIt2.Stop()
			break
		}
		fmt.Println(m.First, m.Second)
	}
	sets := mo.KeyMapSet()
	fmt.Println(sets)
}
