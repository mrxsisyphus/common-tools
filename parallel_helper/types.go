package parallel_helper

import (
	"context"
	"errors"
	"github.com/mrxtryagain/common-tools/logger"
	"runtime"
	"time"
)

const (
	Pending = iota
	Done
	Error
)

var (
	suitableWorkers = runtime.GOMAXPROCS(0) // 获取CPU核数作为worker数量

)

var (
	ErrNoWorker = errors.New("no Worker error")
	ErrFirstN   = errors.New("firstN must > 0")
)

type (
	// SimpleWorker 简易工作单元,无结果,只会返回错误
	SimpleWorker func(ctx context.Context) error
	// DefaultWorker 默认的工作单元 返回一个结果(结果类型已知)和错误(可能)
	DefaultWorker[R any] func(ctx context.Context) (R, error)
	//HighLoadWorker 多返回值工作单元,返回的是一个结果数组和错误
	HighLoadWorker func(ctx context.Context) ([]any, error)
)

type (
	//BeforeStartHook 开始前hook(整个任务)
	BeforeStartHook func(ctx context.Context)
	//FinishHook 结束时候的hook(整个任务)
	FinishHook func(ctx context.Context)
	//BeforeWorkerStartHook 开始单个子任务的函数
	BeforeWorkerStartHook func(ctx context.Context)
	//ErrWorkerHook 单个子任务出错的函数
	ErrWorkerHook func(ctx context.Context, err error)
	//PanicWorkerHook 单个子任务panic的函数
	PanicWorkerHook func(ctx context.Context, panicErr any)
	//FinishWorkerHook 单个子任务结束的函数
	FinishWorkerHook func(ctx context.Context)
)

// DefaultWorkerArgs Worker + index and more... 需要记住workerId的场景
type DefaultWorkerArgs[R any] struct {
	worker   DefaultWorker[R] //worker
	index    int              // index
	workerId string           //worker唯一标识
}

// ResultUnit 返回单元.错误和结果
type ResultUnit[R any] struct {
	Err      error
	Res      R
	Index    int
	Status   WorkerStatus          // 状态
	CostTime time.Duration         // 耗时
	Worker   *DefaultWorkerArgs[R] //关联到的Worker一对一
}

// DefaultParallelWorkReq 默认的并行任务请求 R 是返回值
type DefaultParallelWorkReq[R any] struct {
	Ctx          context.Context         // 上下文
	Id           string                  // 任务的唯一表示
	Desc         string                  // 任务的说明
	ParallelSize int                     // 最大并发数量
	OldWorkers   []*DefaultWorkerArgs[R] //原workers(保留下来)
	Workers      []*DefaultWorkerArgs[R] //新workers(最终的workers
	Options      *Options
}

type (
	UniqueIdStrategy           uint
	WorkerStatus               uint
	CustomUniqueIdStrategyFunc func() string
)

const (
	//UUID + 时间戳
	UUID_WITH_TIME_STAMP UniqueIdStrategy = 0
	//TIME_STAMP 只用时间戳(秒级)
	TIME_STAMP UniqueIdStrategy = 1
	//UUID
	UUID UniqueIdStrategy = 2
)

type Logger interface {
	// Infof must have the same semantics as logger.Infof.
	Infof(format string, args ...interface{})
}

var (
	//defaultLogger 默认日志
	//defaultLogger  = Logger(log.New(os.Stdout, "", log.LstdFlags))
	defaultLogger2 = logger.GetLogger()

	defaultBeforeStartHook = func(ctx context.Context) {
		return
	}

	defaultFinishHook = func(ctx context.Context) {
		return
	}

	defaultBeforeWorkerStartHook = func(ctx context.Context) {
		return
	}

	defaultErrWorkerHook = func(ctx context.Context, err error) {
		return
	}

	defaultPanicWorkerHook = func(ctx context.Context, panicErr any) {
		return
	}

	defaultFinishWorkerHook = func(ctx context.Context) {
		return
	}

	defaultMaxParallelSize = 10000
)
