package io_helpder

import (
	"fmt"
	"path/filepath"
	"testing"
)

func TestFolder(t *testing.T) {
	names, err := ReadDirNames(".")
	if err != nil {
		panic(err)
	}
	fmt.Println(names)
	abs, err := filepath.Abs(".")
	if err != nil {
		return
	}
	fmt.Printf(abs)
	fullNames, err := ReadDirFullNames(".")
	if err != nil {
		return
	}
	fmt.Println(fullNames)
}
