package optional

import (
	"fmt"
	"testing"
)

func Test(t *testing.T) {
	o := Of[int](1)
	get, err := o.Get()
	if err != nil {
		return
	}
	fmt.Println(get)
	e := Empty[int]()
	fmt.Println(e)
	fmt.Println(e.Get())
}
