package convert_helper

import "gopkg.in/yaml.v3"

// YamlUnMarshalToAny 将 bytes 按照toml反序列化给T
func YamlUnMarshalToAny[T any](bs []byte, oldT *T) (*T, error) {
	if oldT == nil {
		oldT = new(T)
	}
	err := yaml.Unmarshal(bs, oldT)
	if err != nil {
		return oldT, err
	}
	return oldT, nil
}
