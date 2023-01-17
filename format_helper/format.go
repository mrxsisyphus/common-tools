package format_helper

import (
	"encoding/json"
	"fmt"
	"github.com/mrxtryagin/common-tools/convert_helper"
)

// PrettyPrint json化输出
func PrettyPrint(v any) string {
	b, err := json.Marshal(v)
	if err != nil {
		return fmt.Sprintf("%v", b)
	}
	return convert_helper.BytesToStr(b)
}

func PrettyPrintWithPanic(v any) string {
	b, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return convert_helper.BytesToStr(b)
}

// PrettyPrintWithIndent json化输出
func PrettyPrintWithIndent(v any, prefix, indent string) string {
	b, err := json.MarshalIndent(v, prefix, indent)
	if err != nil {
		return fmt.Sprintf("%v", b)
	}
	return convert_helper.BytesToStr(b)
}

// PrettyPrintWithDefaultIndent PrettyPrintWithDefaultIndent
func PrettyPrintWithDefaultIndent(v any) string {
	return PrettyPrintWithIndent(v, "", "\t")
}
