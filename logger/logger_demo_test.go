package logger

import (
	"errors"
	"github.com/mrxtryagin/common-tools/time_helper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
	"os"
)

//refer: https://segmentfault.com/a/1190000040984954 日志框架 + 日志分割

//refer: https://www.liwenzhou.com/posts/Go/zap/

var (
	logger       *zap.Logger
	sugarLogger  *zap.SugaredLogger
	customLogger *zap.SugaredLogger
)

// getJSONEncoder JSON 编码器 以json的方式输出
func getJSONEncoder() zapcore.Encoder {
	return zapcore.NewJSONEncoder(*getEncoderConfig())
}

// getTextEncoder 文本编码器,以文本的方式输出
func getTextEncoder() zapcore.Encoder {
	return zapcore.NewConsoleEncoder(*getEncoderConfig())
}

// getEncoderConfig 输出哪些内容
func getEncoderConfig() *zapcore.EncoderConfig {
	// 使用生产环境自带的作为模板
	productionEncoderTemplate := zapcore.EncoderConfig{
		TimeKey:          "ts",                                                                    // 生成json时,time 所对应的key值
		LevelKey:         "level",                                                                 // 生成json时,level 所对应的key值
		NameKey:          "logger",                                                                // 生成json时,name 所对应的key值
		CallerKey:        "caller",                                                                // 生成json时,caller(调用者) 所对应的key值
		FunctionKey:      zapcore.OmitKey,                                                         // 生成json时,当key 被remove 的时候用什么代替
		MessageKey:       "msg",                                                                   // 生成json时,msg 对应的key值(msg 就是内容)
		StacktraceKey:    "stacktrace",                                                            // 生成json时,堆栈信息 对应的key值(msg 就是内容)
		LineEnding:       zapcore.DefaultLineEnding,                                               // 写log 的时候默认的换行符
		EncodeLevel:      zapcore.CapitalColorLevelEncoder,                                        // 对日志等级的配置
		EncodeTime:       zapcore.TimeEncoderOfLayout(time_helper.TimeFormatter_Default_DateTime), // 时间格式的配置
		EncodeDuration:   zapcore.SecondsDurationEncoder,                                          // 序列化 time.Duration格式(?)暂时用不到
		EncodeCaller:     zapcore.ShortCallerEncoder,                                              //时间部分,
		ConsoleSeparator: "\t",                                                                    // 每一块的分割符(console专用)
	}
	zap.NewProductionEncoderConfig()
	// 修改时间部分的表示  如何自定义?(主要是时区) zapcore.TimeEncoderOfLayout() 可以自定义
	//productionEncoderTemplate.EncodeTime = func(time time.Time, encoder zapcore.PrimitiveArrayEncoder) {
	//	// 自定义追加时间的格式
	//	encoder.AppendInt64(1)
	//}
	// 修改日志等级部分的表示(使用大写字母记录)
	return &productionEncoderTemplate
}

// getLogWriter 配置日志输出的位置
func getLogWriter() zapcore.WriteSyncer {
	file, _ := os.Create("./test.logger")
	// 利用io.MultiWriter 同时输出到多个位置,也可以用 zapcore.NewMultiWriteSyncer()
	writer := io.MultiWriter(os.Stdout, file)
	//zapcore.NewMultiWriteSyncer()
	return zapcore.AddSync(writer)
}

func init() {

	// 配置好的
	logger, _ = zap.NewProduction()
	sugarLogger = logger.Sugar() // 与普通logger区别在于 sugarLogger类似于fmt.sprintf可以接受格式化字符串
	custom()

}

func custom() {
	writer := getLogWriter()
	jsonEncoder := getJSONEncoder()
	//testEncoder := getTextEncoder()
	//zapcore.Core需要三个配置——Encoder，WriteSyncer，LogLevel
	// Encoder 编码器,用什么方式写入
	//WriteSyncer 写到那里
	//LogLevel 日志等级,>=这种等级的就会被记录,所以可以通过写自定义函数的方式 进行特殊的处理 而不让他们都记录
	core := zapcore.NewCore(jsonEncoder, writer, zapcore.WarnLevel)
	//当我们不是直接使用初始化好的logger实例记录日志，
	//而是将其包装成一个函数等，此时日录日志的函数调用链会增加，想要获得准确的调用信息就需要通过AddCallerSkip函数来跳过。
	// 比如我们不使用logger这个实例来打印 而使用一个带有logger的方法来打印(logger在里面)就可以用zap.AddCallerSkip(1) 来准确定位
	// New中传的core 可以是一个复合的core 通过 NewTee来组合
	customd := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zap.ErrorLevel))
	customLogger = customd.Sugar()

}

func main() {
	logger.Info("123", zap.String("url", "2323")) // 看上去只是为了更多的放到json里面去
	sugarLogger.Infof("123 %d", 234)
	sugarLogger.Errorf("Error fetching URL %s : Error = %s", "www", errors.New("error url")) // 会把堆栈信息也打出来
	customLogger.Warnf("Error fetching URL %s : Error = %s", "www", errors.New("error url")) // 会把堆栈信息也打出来
}
