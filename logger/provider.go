package logger

import (
	"github.com/mrxtryagain/common-tools/string_helper"
	"github.com/mrxtryagain/common-tools/time_helper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
	"os"
	"time"
)

var (
	_levelToColor = map[zapcore.Level]Color{
		zapcore.DebugLevel:  Magenta,
		zapcore.InfoLevel:   Blue,
		zapcore.WarnLevel:   Yellow,
		zapcore.ErrorLevel:  Red,
		zapcore.DPanicLevel: Red,
		zapcore.PanicLevel:  Red,
		zapcore.FatalLevel:  Red,
	}
	_unknownLevelColor = Red

	_levelToLowercaseColorString = make(map[zapcore.Level]string, len(_levelToColor))
	_levelToCapitalColorString   = make(map[zapcore.Level]string, len(_levelToColor))
)

func init() {
	//初始化颜色
	for level, color := range _levelToColor {
		_levelToLowercaseColorString[level] = color.Add(level.String())
		_levelToCapitalColorString[level] = color.Add(level.CapitalString())
	}
}

// NewDefaultConsoleZapLogger 生产defaultLogger,用于组件库内部使用
// defaultLogger特点
// 1. 有颜色,有代码行
// 2. writer 只到 stdout,只有一个writer
// 3. 给出callback 和堆栈
// 4. debug 以上都打
func NewDefaultConsoleZapLogger() *zap.Logger {
	writer := io.MultiWriter(os.Stdout)
	core := &ZapCore{
		EncoderConfig: defaultLogEncoderConfig(),
		EncoderType:   TextEncoder,
		Writers:       []io.Writer{writer}, // 只有一个stdout
		LevelEnabler:  zapcore.DebugLevel,
	}
	zo := &ZapOptions{
		Cores: []*ZapCore{core},
		OtherOptions: []zap.Option{
			zap.AddCaller(),
			zap.AddStacktrace(zap.ErrorLevel), //error打堆栈
			zap.AddCallerSkip(1),              //统一跳一层
		},
	}
	l, _ := zo.NewZapLogger()
	return l

}

// defaultLogEncoderConfig 默认日志encoder
func defaultLogEncoderConfig() *zapcore.EncoderConfig {
	d := NewDefaultEncoderConfig()
	//caller
	d.EncodeCaller = EncodeCallerPattern1
	// 使用大写亮色
	d.EncodeLevel = EncodeLevelPattern2
	// 时间使用自定义时间格式
	d.EncodeTime = EncodeTimePattern1(time_helper.TimeFormatter_Default_DateTime)
	return d

}

// NewDefaultEncoderConfig 默认config
func NewDefaultEncoderConfig() *zapcore.EncoderConfig {
	config := zap.NewProductionEncoderConfig()
	return &config
}

// EncodeLevelPattern1 自定义日志级别显示,[日志]
func EncodeLevelPattern1(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(string_helper.Concat("[", level.CapitalString(), "]"))
}

// EncodeLevelPattern2 自定义日志级别显示,[日志(颜色)]
func EncodeLevelPattern2(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	s, ok := _levelToCapitalColorString[level]
	if !ok {
		s = _unknownLevelColor.Add(string_helper.Concat("[", level.CapitalString(), "]"))
	}
	enc.AppendString(string_helper.Concat("[", s, "]"))
}

// EncodeTimePattern1 自定义日志时间显示,[时间]
func EncodeTimePattern1(layout string) zapcore.TimeEncoder {
	return func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(string_helper.Concat("[", t.Format(layout), "]"))
	}
}

// EncodeCallerPattern1  自定义caller,[caller]
func EncodeCallerPattern1(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(string_helper.Concat("[", caller.TrimmedPath(), "]"))
}
