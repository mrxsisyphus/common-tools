package convert_helper

import (
	"bytes"
	"github.com/BurntSushi/toml"
)

// TomlUnMarshalToAny 将 bytes 按照toml反序列化给T
func TomlUnMarshalToAny[T any](bs []byte, oldT *T) (*T, error) {
	if oldT == nil {
		oldT = new(T)
	}
	td := toml.NewDecoder(bytes.NewReader(bs))
	_, err := td.Decode(oldT)
	if err != nil {
		return oldT, err
	}
	return oldT, nil
}
