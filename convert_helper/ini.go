package convert_helper

import (
	"gopkg.in/ini.v1"
)

// IniUnMarshalToAny 将 bytes 按照json反序列化给T,
// oldT 必须是一个指针类型 如果oldT为nil 则创一个新的T
func IniUnMarshalToAny[T any](bs []byte, oldT *T) (*T, error) {
	if oldT == nil {
		oldT = new(T)
	}
	cfg, err := ini.LoadSources(
		ini.LoadOptions{
			// 配置 忽略不可解析的行
			SkipUnrecognizableLines: true,
		}, bs)
	if err != nil {
		return oldT, err
	}
	cfg.BlockMode = false
	err = cfg.MapTo(oldT)
	if err != nil {
		return oldT, err
	}
	return oldT, err
}
