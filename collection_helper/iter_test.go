package collection_helper

import (
	"fmt"
	"testing"
)

func TestIter(t *testing.T) {
	temp := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	it := GetIterator(&temp)
	for val := range it.C {
		if val <= 5 {
			it.Stop()
			fmt.Println(val)
		}

	}
	fmt.Println(it)

}
