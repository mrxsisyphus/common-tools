package convert_helper

import (
	"fmt"
	"testing"
)

func TestStrToBytesUnsafe(t *testing.T) {
	a := "12345"
	//unsafe := StrToBytesUnsafe(a)
	//fmt.Println(unsafe, cap(unsafe))

	unsafe2 := StrToBytes(a)
	fmt.Println(unsafe2, cap(unsafe2))
	new1 := append(unsafe2, byte(12))
	fmt.Println(new1, cap(new1))

}
