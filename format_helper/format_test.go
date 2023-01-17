package format_helper

import (
	"fmt"
	"testing"
)

type Person struct {
	Name string
	Age  int
}

func TestPrettyPrintWithIndent(t *testing.T) {
	p := &Person{
		Name: "abc",
		Age:  10,
	}
	fmt.Println(PrettyPrintWithDefaultIndent(p))
}
