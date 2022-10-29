package collection_helper

import (
	"fmt"
	"testing"
)

func TestArrOperation(t *testing.T) {
	temp1 := []int{1, 2, 3, 4, 5}
	temp2 := []int{}
	o := NewSetOperation[int](&temp2)
	fmt.Println(o.Union(&temp1))
	fmt.Println(o.Intersect(&temp1))
	fmt.Println(o.Difference(&temp1))
}
