package search_helper

import (
	"fmt"
	"testing"
)

func TestBisect(t *testing.T) {
	temp := []int{1, 3, 5, 7, 9}
	//左插入的位置
	target := BisectLeft(&temp, 3, 0, len(temp))
	fmt.Println(target)

	//右插入的位置,只有存在的情况下,才有差别,否则都是一样的
	target2 := BisectRight(&temp, 3, 0, len(temp))
	fmt.Println(target2)
}
