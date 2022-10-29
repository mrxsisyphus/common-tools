package time_helper

import (
	"time"
)

const (
	TimeFormatter_Default_Date      = "2006-01-02"
	TimeFormatter_Default_Date2     = "2006-1-2"
	TimeFormatter_Default_DateTime  = "2006-01-02 15:04:05"
	TimeFormatter_Default_DateTime2 = "2006-01-02_15_04_05"
)

// TimeStrToTimeArgs time字符串转time参数
type TimeStrToTimeArgs struct {
	TimeStr        string         //时间字符串
	TimeStrPattern string         // 时间字符串格式(layout)
	TimeZone       *time.Location //timeZone
}

// TimeToTimeStrArgs time转time字符串参数
type TimeToTimeStrArgs struct {
	Time           time.Time      // 时间
	TimeStrPattern string         // 时间字符串格式(layout)
	TimeZone       *time.Location //timeZone 如果有的话需要先转化为此
}

// convertToTime args:TimeStrToTimeArgs 转化为 time_helper.Time
func (args *TimeStrToTimeArgs) convertToTime() (time.Time, error) {
	// 初始化
	if args.TimeStrPattern == "" {
		args.TimeStrPattern = TimeFormatter_Default_DateTime
	}
	if args.TimeZone == nil {
		args.TimeZone = time.Local
	}
	return TimeStrToWithPatternWithCustomLocation(args.TimeStr, args.TimeStrPattern, args.TimeZone)
}

// convertToTimeStr args:TimeToTimeStrArgs 转化为 timeStr
func (args *TimeToTimeStrArgs) convertToTimeStr() string {
	// 初始化
	if args.TimeStrPattern == "" {
		args.TimeStrPattern = TimeFormatter_Default_DateTime
	}
	if args.TimeZone != nil {
		//转换Time
		args.Time = TimeZoneChange(args.Time, args.TimeZone)
	}
	return args.Time.Format(args.TimeStrPattern)
}

// TimeStrToLocationTime 将 timeStr 转换为本地时区的时间,时间的格式使用默认的 TimeFormatter_Default_DateTime
func TimeStrToLocationTime(timeStr string) (time.Time, error) {
	return time.ParseInLocation(TimeFormatter_Default_DateTime, timeStr, time.Local)
}

// TimeStrToLocationTimeWithPattern 将 timeStr 转换为本地时区的时间,时间的格式需要传入
func TimeStrToLocationTimeWithPattern(timeStr string, pattern string) (time.Time, error) {
	return time.ParseInLocation(pattern, timeStr, time.Local)
}

// TimeStrToWithPatternWithCustomLocation 将 timeStr 转换为 loc 时区,时间格式需要传入
func TimeStrToWithPatternWithCustomLocation(timeStr string, pattern string, loc *time.Location) (time.Time, error) {
	return time.ParseInLocation(pattern, timeStr, loc)
}

// TimeToTimeStr 将time 转换为 timeStr,默认使用 TimeFormatter_Default_DateTime
func TimeToTimeStr(inputTime time.Time) string {
	return TimeToStrWithPattern(inputTime, TimeFormatter_Default_DateTime)
}

// TimeToStrWithPattern 将time 转换为 timeStr, pattern 需要传入
func TimeToStrWithPattern(inputTime time.Time, pattern string) string {
	return inputTime.Format(pattern)
}

// TimeZoneChange 直接转换时区,将某个time 转换为另一个时区的time
func TimeZoneChange(inputTime time.Time, loc *time.Location) time.Time {
	return inputTime.In(loc)
}

// GetTimeZoneFromTimeZoneStr 通过timeZoneStr 加载对应的go语言使用的location 按需加载,
// 比如:
//
//		PRC(Asia/Chongqing) 就是中国时区
//	    Japan(Asia/Tokyo)  日本时区
//
// refer: https://studygolang.com/topics/2192
// 具体对应$GOROOT/lib/time_helper/zoneinfo.zip
func GetTimeZoneFromTimeZoneStr(timeZoneStr string) (loc *time.Location, err error) {
	return time.LoadLocation(timeZoneStr)
}

// GetTimeZoneOffsetUTC  通过指定时区 name是你自定义的名字,offset 是距离UTC的偏移量
// 8*60*60(seconds) 也可用 int((8 * time_helper.Hour).Seconds()) 表示
//
//	cusZone := time_helper.FixedZone("UTC+8", 8*60*60)
//
// 方便用于一些多时区不好用TimeZoneStr的场合
// refer:https://immwind.com/golang-time-timezone-and-timestamp/
func GetTimeZoneOffsetUTC(name string, offset int) (loc *time.Location) {
	return time.FixedZone(name, offset)
}

// TimeTruncate 时间truncate
// time_helper.Now().Truncate(time_helper.Hour)) 精确到当前小时
// 2022-06-28 11:10:10.944918 +0800 CST m=+0.000980066 -> 2022-06-28 11:00:00 +0800 CST
// refer: https://www.cnblogs.com/paulwhw/p/14154137.html
func TimeTruncate(inputTime time.Time, duration time.Duration) time.Time {
	return inputTime.Truncate(duration)
}

// DateEqual 2个日期是否相等(注意要在同时区,至于哪个时区 无所谓) 本质就是比较 同一时区下两个日期的年月日是否相等
// 区别与 time_helper.Equal() 这个相等是 时间戳相等(全等) 而DateEqual 只是比较日期(不管时分秒)
func DateEqual(date1, date2 time.Time) bool {
	date1 = date1.UTC()
	date2 = date2.UTC()
	y1, m1, d1 := date1.Date()
	y2, m2, d2 := date2.Date()
	return y1 == y2 && m1 == m2 && d1 == d2
}

//
//// timeBuilder timeZone 最终来源
//const (
//	defaultTimeZone = iota
//	timeZoneFromLoction
//	timeZoneFromStr
//	timeZoneFromCustom
//)
//
//// TimeBuilder Time 的 建造者模式 链式调用
//type timeBuilder struct {
//	timeZoneCustom
//	time_helper            *time_helper.Time
//	timeStr         string
//	timeStrPattern  string
//	timeZone        *time_helper.Location
//	timeZoneStr     string
//	timeZoneDefiner int
//}
//
//type timeZoneCustom struct {
//	timeZoneCustomName string
//	timeZoneUTCOffset  int
//}
//
//func NewTimeBuilder() *timeBuilder {
//	return &timeBuilder{}
//}
//
//func (builder *timeBuilder) SetTime(inputTime time_helper.Time) *timeBuilder {
//	builder.time_helper = &inputTime
//	return builder
//}
//func (builder *timeBuilder) SetTimeStr(timeStr string) *timeBuilder {
//	builder.timeStr = timeStr
//	return builder
//}
//
//func (builder *timeBuilder) SetTimeStrPattern(timePatternStr string) *timeBuilder {
//	builder.timeStrPattern = timePatternStr
//	return builder
//}
//
//func (builder *timeBuilder) ChangeZoneByLocation(loc *time_helper.Location) *timeBuilder {
//	builder.timeZone = loc
//	builder.timeZoneDefiner = timeZoneFromLoction
//	return builder
//}
//
//func (builder *timeBuilder) ChangeZoneByCustomLocation(name string, offsetUtc int) *timeBuilder {
//	builder.timeZoneCustomName = name
//	builder.timeZoneUTCOffset = offsetUtc
//	builder.timeZoneDefiner = timeZoneFromCustom
//	return builder
//}
//
//func (builder *timeBuilder) ChangeZoneByStr(locationStr string) *timeBuilder {
//	builder.timeZoneStr = locationStr
//	builder.timeZoneDefiner = timeZoneFromStr
//	return builder
//}
//
//func (builder *timeBuilder) GetTimeZone() (*time_helper.Location, error) {
//	var (
//		err      error
//		timeZone *time_helper.Location
//	)
//	switch builder.timeZoneDefiner {
//	case defaultTimeZone:
//	case timeZoneFromLoction:
//		timeZone = builder.timeZone
//	case timeZoneFromStr:
//		timeZone, err = GetTimeZoneFromTimeZoneStr(builder.timeZoneStr)
//		if err != nil {
//			return nil, err
//		}
//	case timeZoneFromCustom:
//		timeZone = GetTimeZoneOffsetUTC(builder.timeZoneCustomName, builder.timeZoneUTCOffset)
//	}
//	return timeZone, err
//}
//
//func (builder *timeBuilder) ToTime() (time_helper.Time, error) {
//	var (
//		err      error
//		timeZone *time_helper.Location
//	)
//
//	if builder.time_helper != nil {
//		// 有time 没有timeStr
//		//如果有time 导出time 那么与 timeStr 和 timeStrPattern 无关
//		// 只会做时区相关的变化
//
//		timeZone, err = builder.GetTimeZone()
//		if err != nil {
//			return time_helper.Time{}, err
//		}
//		if timeZone != nil {
//			return TimeZoneChange(*builder.time_helper, timeZone), nil
//		} else {
//			return *builder.time_helper, nil
//		}
//	} else {
//		if builder.timeStr == "" {
//			return time_helper.Time{}, errors.New("no timeStr")
//		} else {
//			// 没有time 有timeStr
//			if builder.timeStrPattern == "" {
//				builder.timeStrPattern = TimeFormatter_Default_DateTime
//			}
//			// 1. 拿到location
//			timeZone, err = builder.GetTimeZone()
//			if err != nil {
//				return time_helper.Time{}, err
//			}
//			//2. 转换成time
//			return TimeStrToWithPatternWithCustomLocation(builder.timeStr, builder.timeStrPattern, timeZone)
//		}
//
//	}
//}
//
//func (builder *timeBuilder) ToTimeStr() (string, error) {
//	var (
//		err      error
//		timeZone *time_helper.Location
//	)
//
//	if builder.time_helper != nil {
//
//	} else {
//		// 没有time 没有timeStr
//		if builder.timeStr == "" {
//			return "", nil
//		}
//		if builder.timeStrPattern == "" {
//			builder.timeStrPattern = TimeFormatter_Default_DateTime
//		}
//		// 看看时区
//		timeZone, err = builder.GetTimeZone()
//		if err != nil {
//			return "", err
//		}
//		//timeStr 先转换成time
//
//		// time_helper -> timeStr
//
//	}
//}
