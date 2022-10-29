package string_helper

import (
	"bytes"
)

// Reverse 反转字符串 利用i，j来控制首尾两个位置，交换对应位置的元素即可实现字符串反转
func Reverse(s string) string {
	a := []rune(s)
	for i, j := 0, len(a)-1; i < j; i, j = i+1, j-1 {
		a[i], a[j] = a[j], a[i]
	}
	return string(a)
}

// Concat 结合 string
func Concat(strs ...string) string {
	strBuf := bytes.NewBufferString("")
	for _, s := range strs {
		strBuf.WriteString(s)
	}
	return strBuf.String()
}
