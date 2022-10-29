package time_helper

import (
	"fmt"
	"testing"
	"time"
)

func TestTime(t *testing.T) {
	now := time.Now()
	fmt.Println(now)
	utc := now.UTC()
	fmt.Println(utc)
	// 时区更改
	utc2 := TimeZoneChange(now, time.UTC)
	fmt.Println(utc2)
	//日期相等
	isDateEqual := DateEqual(utc, utc2)
	fmt.Println(isDateEqual)
	// 获取固定的时区
	customUtc := GetTimeZoneOffsetUTC("utc+5", 5*60*60)
	fmt.Println(customUtc.String())
	//时区更换
	utc3 := TimeZoneChange(utc2, customUtc)
	fmt.Println(utc3)
	cst, err := GetTimeZoneFromTimeZoneStr("PRC")
	if err != nil {
		panic(err)
	}
	fmt.Println(cst.String())
	args := &TimeToTimeStrArgs{
		Time:     utc3,
		TimeZone: cst,
	}
	resultString := args.convertToTimeStr()
	fmt.Println(resultString)

}
