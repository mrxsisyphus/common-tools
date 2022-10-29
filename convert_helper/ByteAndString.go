package convert_helper

import (
	"encoding/json"
	"fmt"
	"html/template"
	"reflect"
	"strconv"
	"unsafe"
)

// BytesToStr unsafe 转换, 将 []byte 转换为 string
func BytesToStr(p []byte) string {
	return *(*string)(unsafe.Pointer(&p))
}

// StrToBytes unsafe 转换, 将 string 转换为 []byte,利用SliceHeader
// refer: https://medium.com/@kevinbai/golang-%E4%B8%AD-string-%E4%B8%8E-byte-%E4%BA%92%E8%BD%AC%E4%BC%98%E5%8C%96-6651feb4e1f2
func StrToBytes(str string) []byte {
	strHeader := (*reflect.StringHeader)(unsafe.Pointer(&str))
	return *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Data: strHeader.Data,
		Len:  strHeader.Len,
		Cap:  strHeader.Len,
	}))
}

// StrToBytesUnsafe unsafe 转换, 请确保转换后的 []byte 不涉及 cap() 操作, 将 string 转换为 []byte
// 与StrToBytes相比 没有设置cap, 所以这里的cap是不对的
// 同时: 对byte 进行append 会报错
// unexpected fault address 0x1043fef0a
// fatal error: fault
// [signal SIGBUS: bus error code=0x1 addr=0x1043fef0a pc=0x1043fe4e8]
func StrToBytesUnsafe(str string) []byte {
	return *(*[]byte)(unsafe.Pointer(&str))
}

// AnyToString string化函数
func AnyToString(i any) string {
	i = indirectToStringerOrError(i)

	switch s := i.(type) {
	case string:
		return s
	case bool:
		return strconv.FormatBool(s)
	case float64:
		return strconv.FormatFloat(s, 'f', -1, 64)
	case float32:
		return strconv.FormatFloat(float64(s), 'f', -1, 32)
	case int:
		return strconv.Itoa(s)
	case int64:
		return strconv.FormatInt(s, 10)
	case int32:
		return strconv.Itoa(int(s))
	case int16:
		return strconv.FormatInt(int64(s), 10)
	case int8:
		return strconv.FormatInt(int64(s), 10)
	case uint:
		return strconv.FormatUint(uint64(s), 10)
	case uint64:
		return strconv.FormatUint(s, 10)
	case uint32:
		return strconv.FormatUint(uint64(s), 10)
	case uint16:
		return strconv.FormatUint(uint64(s), 10)
	case uint8:
		return strconv.FormatUint(uint64(s), 10)
	case json.Number:
		return s.String()
	case []byte:
		return string(s)
	case template.HTML:
		return string(s)
	case template.URL:
		return string(s)
	case template.JS:
		return string(s)
	case template.CSS:
		return string(s)
	case template.HTMLAttr:
		return string(s)
	case nil:
		return ""
	case fmt.Stringer:
		return s.String()
	case error:
		return s.Error()
	default:
		return fmt.Sprintf("%v", s)
	}
}
