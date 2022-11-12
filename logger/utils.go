package logger

import (
	"fmt"
	"go.uber.org/zap"
)

/*
refer: https://colobu.com/2018/11/03/get-function-name-in-go/ runtime.Caller的原理
*/

// 通用方法与api保持一致
func Debug(msg string, args ...zap.Field) {
	globalLogger.zapLogger.Debug(msg, args...)
}
func Info(msg string, args ...zap.Field) {
	globalLogger.zapLogger.Info(msg, args...)
}

func Warn(msg string, args ...zap.Field) {
	globalLogger.zapLogger.Warn(msg, args...)
}

func Error(msg string, args ...zap.Field) {
	globalLogger.zapLogger.Error(msg, args...)
}

func Panic(msg string, args ...zap.Field) {
	globalLogger.zapLogger.Panic(msg, args...)
}

func Fatal(msg string, args ...zap.Field) {
	globalLogger.zapLogger.Fatal(msg, args...)
}

func Debugf(format string, args ...interface{}) {
	LoggerMsg := fmt.Sprintf(format, args...)
	globalLogger.zapLogger.Debug(LoggerMsg)
}

func Infof(format string, args ...interface{}) {
	LoggerMsg := fmt.Sprintf(format, args...)
	globalLogger.zapLogger.Info(LoggerMsg)
}

func Warnf(format string, args ...interface{}) {
	LoggerMsg := fmt.Sprintf(format, args...)
	globalLogger.zapLogger.Warn(LoggerMsg)
}

func Errorf(format string, args ...interface{}) {
	LoggerMsg := fmt.Sprintf(format, args...)
	globalLogger.zapLogger.Error(LoggerMsg)
}

func Panicf(format string, args ...interface{}) {
	LoggerMsg := fmt.Sprintf(format, args...)
	globalLogger.zapLogger.Panic(LoggerMsg)
}

func Fatalf(format string, args ...interface{}) {
	LoggerMsg := fmt.Sprintf(format, args...)
	globalLogger.zapLogger.Fatal(LoggerMsg)
}
func (l *Logger) Debug(msg string, args ...zap.Field) {
	l.zapLogger.Debug(msg, args...)
}

func (l *Logger) Info(msg string, args ...zap.Field) {
	l.zapLogger.Info(msg, args...)
}

func (l *Logger) Warn(msg string, args ...zap.Field) {
	l.zapLogger.Warn(msg, args...)
}

func (l *Logger) Error(msg string, args ...zap.Field) {
	l.zapLogger.Error(msg, args...)
}

func (l *Logger) Panic(msg string, args ...zap.Field) {
	l.zapLogger.Panic(msg, args...)
}

func (l *Logger) Fatal(msg string, args ...zap.Field) {
	l.zapLogger.Fatal(msg, args...)
}

func (l *Logger) Debugf(format string, args ...interface{}) {
	LoggerMsg := fmt.Sprintf(format, args...)
	l.zapLogger.Debug(LoggerMsg)
}

func (l *Logger) Infof(format string, args ...interface{}) {
	LoggerMsg := fmt.Sprintf(format, args...)
	l.zapLogger.Info(LoggerMsg)
}

func (l *Logger) Warnf(format string, args ...interface{}) {
	LoggerMsg := fmt.Sprintf(format, args...)
	l.zapLogger.Warn(LoggerMsg)
}

func (l *Logger) Errorf(format string, args ...interface{}) {
	LoggerMsg := fmt.Sprintf(format, args...)
	l.zapLogger.Error(LoggerMsg)
}

func (l *Logger) Panicf(format string, args ...interface{}) {
	LoggerMsg := fmt.Sprintf(format, args...)
	l.zapLogger.Panic(LoggerMsg)
}

func (l *Logger) Fatalf(format string, args ...interface{}) {
	LoggerMsg := fmt.Sprintf(format, args...)
	l.zapLogger.Fatal(LoggerMsg)
}
