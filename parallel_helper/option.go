package parallel_helper

type Option interface {
	apply(*Options)
}

type optionFunc func(*Options)

func (f optionFunc) apply(o *Options) {
	//调用传入的函数
	f(o)
}

type Options struct {
	//Desc 任务描述
	Desc string

	//ParallelSize 平行工作的大小,-1为不上限,0(默认)是cpu核数
	ParallelSize int

	//ParallelSameWithLength 设置平行工作大小与workers长度一样长 默认为false
	ParallelSameWithLength bool

	//UniqueIdStrategy taskId的策略
	UniqueIdStrategy UniqueIdStrategy

	//CustomUniqueIdStrategy 自定义的策略
	CustomUniqueIdStrategy CustomUniqueIdStrategyFunc

	//TimeOutDuration  超时周期(会放到ctx里面)
	//TimeOutDuration time.Duration

	//NotCollectPanic 是否不收集panic当作error处理
	NotCollectPanic bool

	//IsMixed 是否要打乱任务执行单元的顺序(注意 这样你可能不能保证严格的一致性)
	IsMixed bool

	//Logger
	Logger Logger

	//执行策略
	// FIRST_COMPLETED FIRST_EXCEPTION ALL_COMPLETED 等

	// hooks
	//BeforeStartHook 开始的时候的函数(整个任务)
	BeforeStartHook BeforeStartHook

	//FinishHook 结束的时候的函数(整个任务)
	FinishHook FinishHook

	BeforeWorkerStartHook BeforeWorkerStartHook
	ErrWorkerHook         ErrWorkerHook
	PanicWorkerHook       PanicWorkerHook
	FinishWorkerHook      FinishWorkerHook
}

func WithOptions(options Options) Option {
	return optionFunc(func(opts *Options) {
		*opts = options
	})
}

func WithDesc(desc string) Option {
	return optionFunc(func(opts *Options) {
		opts.Desc = desc
	})
}

func WithParallelSize(parallelSize int) Option {
	return optionFunc(func(opts *Options) {
		opts.ParallelSize = parallelSize
	})
}

func WithUniqueIdStrategy(uniqueIdStrategy UniqueIdStrategy) Option {
	return optionFunc(func(opts *Options) {
		opts.UniqueIdStrategy = uniqueIdStrategy
	})
}

func WithCustomUniqueIdStrategy(customUniqueIdStrategy CustomUniqueIdStrategyFunc) Option {
	return optionFunc(func(opts *Options) {
		opts.CustomUniqueIdStrategy = customUniqueIdStrategy
	})
}

//func WithTimeOutDuration(timeOutDuration time.Duration) Option {
//	return optionFunc(func(opts *Options) {
//		opts.TimeOutDuration = timeOutDuration
//	})
//}

func WithCollectPanic(collectPanic bool) Option {
	return optionFunc(func(opts *Options) {
		opts.NotCollectPanic = collectPanic
	})
}

func WithIsMixed(isMixed bool) Option {
	return optionFunc(func(opts *Options) {
		opts.IsMixed = isMixed
	})
}

func WithParallelSameWithLength(parallelSameWithLength bool) Option {
	return optionFunc(func(opts *Options) {
		opts.ParallelSameWithLength = parallelSameWithLength
	})
}

func WithLogger(logger Logger) Option {
	return optionFunc(func(opts *Options) {
		opts.Logger = logger
	})
}

func WithBeforeStartHook(hook BeforeStartHook) Option {
	return optionFunc(func(opts *Options) {
		opts.BeforeStartHook = hook
	})
}

func WithFinishHook(hook FinishHook) Option {
	return optionFunc(func(opts *Options) {
		opts.FinishHook = hook
	})
}

func WithBeforeWorkerStartHook(hook BeforeWorkerStartHook) Option {
	return optionFunc(func(opts *Options) {
		opts.BeforeWorkerStartHook = hook
	})
}

func WithErrWorkerHook(hook ErrWorkerHook) Option {
	return optionFunc(func(opts *Options) {
		opts.ErrWorkerHook = hook
	})
}

func WithPanicWorkerHook(hook PanicWorkerHook) Option {
	return optionFunc(func(opts *Options) {
		opts.PanicWorkerHook = hook
	})
}

func WithFinishWorkerHook(hook FinishWorkerHook) Option {
	return optionFunc(func(opts *Options) {
		opts.FinishWorkerHook = hook
	})
}

func loadOptions(options ...Option) *Options {
	opts := new(Options)
	for _, op := range options {
		op.apply(opts)
	}
	return opts
}
