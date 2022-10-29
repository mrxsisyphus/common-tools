package parallel_helper

import (
	"bytes"
	"runtime"
	"strconv"
)

/**
设计:
关键设计:
gather,wait
1. 无需返回值的函数
直接执行然后recover 防止意外
wg 分批次执行
2. 有返回值的函数通过数组组装返回值


*/

// RoutineId 获得gorountineid for debug
func RoutineId() uint64 {
	b := make([]byte, 64)
	b = b[:runtime.Stack(b, false)]
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	b = b[:bytes.IndexByte(b, ' ')]
	// if error, just return 0
	n, _ := strconv.ParseUint(string(b), 10, 64)

	return n
}
