package convert_helper

import (
	"encoding/json"
)

// JsonUnMarshalToAny 将 bytes 按照json反序列化给T,
// oldT 必须是一个指针类型 如果oldT为nil 则创一个新的T
func JsonUnMarshalToAny[T any](bs []byte, oldT *T) (*T, error) {
	if oldT == nil {
		oldT = new(T)
	}
	err := json.Unmarshal(bs, oldT)
	if err != nil {
		return oldT, err
	}
	return oldT, nil
}
