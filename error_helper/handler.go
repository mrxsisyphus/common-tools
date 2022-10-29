package error_helper

import "fmt"

// PanicWithStr 带字符串占位符的panic
func PanicWithStr(format string, args ...any) {
	panic(fmt.Sprintf(format, args))
}
