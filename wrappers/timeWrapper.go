package wrappers

// FuncCostTimer 函数花费成本
//
//	func FuncCostTimer(funcName string) func() {
//		start := time.Now()
//		return func() {
//			//配合defer 该函数的执行会在最后
//
//		}
//	}
//

// ErrorWrapper 错误包装,保证执行
func ErrorWrapper(errorCount int, handler func() error) error {
	var lastErr error

	for i := 0; i < errorCount; i++ {
		if err := handler(); err != nil {
			lastErr = err
		} else {
			return nil
		}
	}
	return lastErr
}
