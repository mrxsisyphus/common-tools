package logger

import "testing"

func TestLogger(t *testing.T) {
	//l := NewDefaultLogger()
	Infof("你 %s 好啊", "124")
	Info("你好啊")
	globalLogger.Infof("你 %s 好啊", "124")
}
