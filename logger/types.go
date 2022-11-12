package logger

import (
	"errors"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
	"time"
)

const ()

type (
	EncoderType int
)

var (
	ErrNoZapCore = errors.New("zap core not found")
)

const (
	JsonEncoder EncoderType = 0
	TextEncoder EncoderType = 1
)

type LogOptions struct {
	ZapOptions
	RotateLogOptions
}

// ZapCore zap的核心 单个核心的配置
type ZapCore struct {
	EncoderConfig *zapcore.EncoderConfig //输出的配置
	EncoderType   EncoderType            //encoder类型 目前只有text和json
	Writers       []io.Writer            //输出的地方(可以是多个,可以是任意对象,包括了日志分割创造的writer)
	LevelEnabler  zapcore.LevelEnabler   //level优先级设置
}

// ZapOptions zap本身的配置,一个zap可以是多个核心
type ZapOptions struct {
	Cores        []*ZapCore
	OtherOptions []zap.Option
}

// RotateLogOptions 日志分割相关的配置
// 目前分为两个 一个是lumberjack(使用大小) 另一个是 rotatelogs(使用时间),目前这里只是整合两者
type RotateLogOptions struct {
	FileName string //fileName
}

type RotateLogsWithSize struct {
	RotateLogOptions
	MaxSize    int  //在进行切割之前，日志文件的最大大小（以 MB 为单位）
	MaxBackups int  //保留旧文件的最大个数
	MaxAge     int  //保留旧文件的最大天数
	Compress   bool //是否压缩 / 归档旧文件
}

// RotateLogsWithTime https://github.com/lestrrat-go/file-rotatelogs
type RotateLogsWithTime struct {
	RotateLogOptions                  // RotateLogsWithTime的时间 可以包含格式
	Clock            rotatelogs.Clock //clock 默认情况下是当前时间
	Location         *time.Location   //时区 是 Clock的代替品,可以用于指名某个时区
	LinkName         string           //最新日志文件的软链接名字,默认不创建
	MaxAge           int              //保留旧文件的最大天数(默认7days)
	RotationTime     time.Duration    //rotation周期 默认86400s(-1表示不限制)
	RotationCount    int              // rotation的数量 默认为-1 表示不限制

}

// RotateLogsMixed RotateLogsWithSize 和 RotateLogsWithTime 结合(暂时不做)
type RotateLogsMixed struct {
}

// NewZapCore ZapCore -> zapcore.Core
func (z *ZapCore) NewZapCore() zapcore.Core {
	//encoder
	var encoder zapcore.Encoder
	switch z.EncoderType {
	case JsonEncoder:
		encoder = zapcore.NewJSONEncoder(*z.EncoderConfig)
	case TextEncoder:
		encoder = zapcore.NewConsoleEncoder(*z.EncoderConfig)
	default:
		encoder = zapcore.NewConsoleEncoder(*z.EncoderConfig)
	}
	//writer
	multiWriter := io.MultiWriter(z.Writers...)
	writeSyncer := zapcore.AddSync(multiWriter)
	//newCore
	return zapcore.NewCore(encoder, writeSyncer, z.LevelEnabler)
}

// NewZapLogger logger 可以相互转换
func (zo *ZapOptions) NewZapLogger() (*zap.Logger, error) {
	cores := zo.Cores
	if len(cores) <= 0 {
		return nil, ErrNoZapCore
	}
	res := make([]zapcore.Core, len(cores))
	for i, core := range cores {
		res[i] = core.NewZapCore()
	}
	//最终的core
	lastCore := zapcore.NewTee(res...)
	_logger := zap.New(lastCore, zo.OtherOptions...)
	return _logger, nil
}
