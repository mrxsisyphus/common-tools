package io_helpder

import (
	"fmt"
	"testing"
)

func TestPath(t *testing.T) {
	name := "path_test.go"
	absname := RelativePath(name)
	fmt.Println(absname)
}
