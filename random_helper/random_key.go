package random_helper

import (
	"fmt"
	"github.com/gofrs/uuid"
	"github.com/mrxtryagin/common-tools/convert_helper"
	"time"
)

// GetUUID 获得一个uuid 返回错误
func GetUUID() (string, error) {
	v4, err := uuid.NewV4()
	if err != nil {
		return "", err
	}
	return v4.String(), nil
}

// GetUUIDWithTimeStamp  获得一个uuid + time_helper 格式为 timestamp_uuid,返回错误
func GetUUIDWithTimeStamp() (string, error) {
	v4, err := uuid.NewV4()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%d_%s", time.Now().Unix(), v4.String()), nil
}

func GetIDWithTimeStamp() string {
	return convert_helper.AnyToString(time.Now().Unix())
}

func GetIDWithNanoTimeStamp() string {
	return convert_helper.AnyToString(time.Now().UnixNano())
}
