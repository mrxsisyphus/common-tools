package parallel_helper

import (
	"context"
	"errors"
	"fmt"
	"github.com/mrxtryagin/common-tools/collection_helper"
	"github.com/mrxtryagin/common-tools/random_helper"
	"github.com/panjf2000/ants/v2"
	"golang.org/x/sync/semaphore"
	"math"
	"sync"
	"time"
)

/**

//3: 返回前n个最快的协程(搞定)
//4: wait 现在还是有顺序性,如何保证无序性(搞定)
//5: 批量协程本身是否要使用池(ants等)[后期可以考虑]
//6: 单元的执行函数抽象出来,方便进行一些异构处理(搞定)
*/

var (
	pool *ants.Pool
)

func init() {
	//pool, _ = ants.NewPool(100, ants.WithMaxBlockingTasks(1000000000))
}

// NewDefaultParallelWorkReq 后期可以用option进行进行设置改造
func NewDefaultParallelWorkReq[R any](ctx context.Context, workers []DefaultWorker[R], ops ...Option) (*DefaultParallelWorkReq[R], error) {
	if len(workers) <= 0 {
		return nil, ErrNoWorker
	}
	// 处理options
	opts := loadOptions(ops...)
	// 处理常规req
	req := &DefaultParallelWorkReq[R]{
		Ctx:     ctx,
		Options: opts,
	}
	// 处理id
	switch req.Options.UniqueIdStrategy {
	case TIME_STAMP:
		req.Id = random_helper.GetIDWithTimeStamp()
	case UUID:
		uuid, err := random_helper.GetUUID()
		if err != nil {
			return nil, err
		}
		req.Id = uuid
	case UUID_WITH_TIME_STAMP:
		uuid, err := random_helper.GetUUIDWithTimeStamp()
		if err != nil {
			return nil, err
		}
		req.Id = uuid
	}
	//  处理worker
	workerArgs := make([]*DefaultWorkerArgs[R], len(workers))
	for i := 0; i < len(workers); i++ {
		workerArgs[i] = &DefaultWorkerArgs[R]{
			index:    i, //记录原index
			worker:   workers[i],
			workerId: fmt.Sprintf("%s_worker_%d", req.Id, i),
		}
	}
	req.OldWorkers = workerArgs

	// 各种初始化
	if req.Options.Desc != "" {
		req.Desc = req.Options.Desc
	}

	if req.Options.ParallelSize == 0 {
		req.ParallelSize = suitableWorkers
	} else {
		if req.Options.ParallelSize == -1 {
			req.ParallelSize = defaultMaxParallelSize
		} else {
			req.ParallelSize = req.Options.ParallelSize
		}
	}
	if req.Options.ParallelSameWithLength {
		req.ParallelSize = len(workers)
	}
	if req.Options.Logger == nil {
		req.Options.Logger = defaultLogger2
	}
	if req.Options.BeforeStartHook == nil {
		req.Options.BeforeStartHook = defaultBeforeStartHook
	}
	if req.Options.FinishHook == nil {
		req.Options.FinishHook = defaultFinishHook
	}
	if req.Options.BeforeWorkerStartHook == nil {
		req.Options.BeforeWorkerStartHook = defaultBeforeWorkerStartHook
	}
	if req.Options.ErrWorkerHook == nil {
		req.Options.ErrWorkerHook = defaultErrWorkerHook
	}
	if req.Options.PanicWorkerHook == nil {
		req.Options.PanicWorkerHook = defaultPanicWorkerHook
	}
	if req.Options.FinishWorkerHook == nil {
		req.Options.FinishWorkerHook = defaultFinishWorkerHook
	}
	//处理 混杂
	if req.Options.IsMixed {
		//使用hash散列打散
		// 后期可以自己设置打散函数/策略
		req.Workers = *collection_helper.Distinct[*DefaultWorkerArgs[R]](&req.OldWorkers)
	} else {
		// 新的就等于旧的
		req.Workers = req.OldWorkers
	}
	return req, nil
}

// Gather 保证顺序的一致性
// 使用顺序调度(组内抢占),进来什么顺序,出去就是什么顺序(利用切片下标)
// 保证执行全部完成(不会主动打断循环)
func (req *DefaultParallelWorkReq[R]) Gather() []*ResultUnit[R] {
	var (
		wg sync.WaitGroup
	)
	req.Options.Logger.Infof("task[Id=%s]-%s gather start..\n", req.Id, req.Desc)
	now := time.Now()
	//beforeHook
	req.Options.BeforeStartHook(req.Ctx)
	//limiter
	limiter := make(chan struct{}, req.ParallelSize)
	for i := 0; i < req.ParallelSize; i++ {
		limiter <- struct{}{}
	}
	//因为是gather 所以要保证顺序性,需要提前new好result数组(每个work 都会有一个result)
	workers := req.Workers
	// safe for result use index
	// 使用 切片 索引下标,来进行,因为需要顺序返回
	result := make([]*ResultUnit[R], len(workers))
	// for-range 配合select, 顺序的分批次执行workders
	for i, worker := range workers {
		//wg++
		wg.Add(1)
		// 创建引用,防止问题
		handler := worker
		index := i
		select {
		// 先使用select保证并发数不至于过大
		case <-limiter:
			// 一个工作执行单元
			go func() {
				defer func() {
					wg.Done()             //处理完成
					limiter <- struct{}{} //增加limiter
				}()
				req.gatherWork0(index, handler, result)
			}()
		}
	}
	// 直到结束为止
	wg.Wait()
	cost := time.Since(now).Milliseconds()
	req.Options.Logger.Infof("task[Id=%s]-%s gather end cost:%d ms\n", req.Id, req.Desc, cost)
	req.Options.FinishHook(req.Ctx) // 结束钩子
	return result
}

// Wait 不保证顺序的一致性
// 使用顺序调度(组内抢占),但是ParallelSize内按照速度快的排序(利用切片append)
// 保证执行全部完成(在不会外部打断的情况下)
// 可以适配超时场景
// refer: 关于信号量 https://juejin.cn/post/7095701003172839432
// refer2: https://www.jianshu.com/p/e481064aeab4
func (req *DefaultParallelWorkReq[R]) Wait() ([]*ResultUnit[R], error) {

	req.Options.Logger.Infof("task[Id=%s]-%s wait start..\n", req.Id, req.Desc)
	now := time.Now() // 当前时间
	//beforeHook
	//beforeHook
	req.Options.BeforeStartHook(req.Ctx)
	//limiter
	// 限制器(由于要顺序调度,所以优先放好)
	limiter := make(chan struct{}, req.ParallelSize)
	for i := 0; i < req.ParallelSize; i++ {
		limiter <- struct{}{}
	}
	workers := req.Workers
	workerLength := len(workers)
	// 假设他们能够拿到所有的ResultUnit
	result := make([]*ResultUnit[R], 0, workerLength)
	// resultChan 是接受结果的地方,这个地方容易造成生产者和消费者的内存泄漏,所以需要巧妙处理
	//使用 有缓存chan 保证极端情况(在资源充足的情况下,缓存数量与发送者最好一致,这样的话起码能保证生产者不pending)
	resultChan := make(chan *ResultUnit[R], req.ParallelSize)
	// 结果合并结束的信号(如果正常结束,可能消费的比较慢 使用这个通道来强制保证消费结束)
	mergeDone := make(chan struct{})
	// 使用一个stopCh,来让消费者通知生产者,如果设置result pending了,可以退出了
	stopCh := make(chan struct{})
	sem := semaphore.NewWeighted(int64(workerLength))
	go func() {
		//发送合并结束的信号标志合并结束
		defer close(mergeDone)
		// 如果直接range,在没有被关闭的情况下,这里必须要先能进来才会退出(如果刚好执行的点就是执行完的时候,就会导致resultChan 一直没有数据进来了)
		// 但是直接用select,则超时的时候 直接通知退出.
		for {
			//判断超时问题,保证消费者退出
			select {
			// 超时的话
			case <-req.Ctx.Done():
				//req.Options.Logger.Infof("消费者退出...\n")
				//消费者退出前通知生产者
				close(stopCh)
				return // 退出消费者
			case res, isOk := <-resultChan: // resultChan 被关闭,直接退出;resultChan 没有被关闭,被pending很长一段时间,也会由于前面的退出
				if !isOk {
					//req.Options.Logger.Infof("resultChan 被关闭...,不再接受数据 \n")
					//被关闭,直接返回
					return
				} else {
					// 正常的话 append
					result = append(result, res)
				}

			}
		}
		// 当 resultChan被关闭时且channel中所有的值都已经被处理完毕后, 将执行到这一行
	}()
	// break 外层for
	// refer: https://blog.csdn.net/u011461385/article/details/106017483
L:
	for i, worker := range workers {
		ctxErr := sem.Acquire(req.Ctx, 1)
		if ctxErr != nil {
			return result, ctxErr
		}
		// 将worker 拿出来
		handler := worker
		index := i
		select {
		case <-req.Ctx.Done(): // Ctx 控制退出 wait的关键
			req.Options.Logger.Infof("task[Id=%s]-%s wait exit[ctx done]...\n", req.Id, req.Desc)
			break L // 退出外循环
		// 先使用select保证并发数不至于过大(保证在parallel内)
		case <-limiter:
			go func() {
				defer func() {
					sem.Release(1)        //信号量+1
					limiter <- struct{}{} // limiter补充
				}()
				req.waitWork0(index, handler, resultChan, stopCh)

			}()
		}
	}
	// 请求所有的worker,会wait 直到之前的搞完 或者接受一个错误,
	// 这就是 使用这个 而不使用waitGroup的原因,原因就是这里可以接受一个ctx的done事件
	err := sem.Acquire(req.Ctx, int64(workerLength))
	if err != nil {
		//fmt.Println("sem.Acquire Err:", Err)
		//如果发生错误 直接退出(一般是超时)
		return result, err
	}

	//close(limiter) //关闭limiter,这里可以不做,因为不需要通知 gc会帮我们完成
	close(resultChan) //关闭resultChan  这里需要做,不然上面的消费者没办法退出
	<-mergeDone       // 这里主要是等消费完成,防止消费过慢的情况
	cost := time.Since(now).Milliseconds()
	req.Options.Logger.Infof("task[Id=%s]-%s wait end cost:%d ms\n", req.Id, req.Desc, cost)
	req.Options.FinishHook(req.Ctx) // 结束钩子
	return result, nil

}

// WaitWithPreemptive 与 wait 类似,但是是完全抢占式
// 使用抢占式调度
// 保证执行全部完成(在不会外部打断的情况下)
// 可以适配超时场景
// 每个worker 协程都会启动,只是在拿到令牌(limiter)的情况下会被阻塞,所以内存会更大
// 更符合抢占式的要求,看哪个协程先抢到,哪个来
func (req *DefaultParallelWorkReq[R]) WaitWithPreemptive() ([]*ResultUnit[R], error) {

	req.Options.Logger.Infof("task[Id=%s]-%s waitWithPreemptive start..\n", req.Id, req.Desc)
	now := time.Now() // 当前时间
	//beforeHook
	//beforeHook
	req.Options.BeforeStartHook(req.Ctx)
	//limiter
	// 限制器
	limiter := make(chan struct{}, req.ParallelSize)
	workers := req.Workers
	workerLength := len(workers)
	// 假设他们能够拿到所有的ResultUnit
	result := make([]*ResultUnit[R], 0, workerLength)
	// resultChan 是接受结果的地方,这个地方容易造成生产者和消费者的内存泄漏,所以需要巧妙处理
	//使用 有缓存chan 保证极端情况(在资源充足的情况下,缓存数量与发送者最好一致,这样的话起码能保证生产者不pending)
	resultChan := make(chan *ResultUnit[R], req.ParallelSize)
	// 结果合并结束的信号(如果正常结束,可能消费的比较慢 使用这个通道来强制保证消费结束)
	mergeDone := make(chan struct{})
	// 使用一个stopCh,来让消费者通知生产者,如果设置result pending了,可以退出了
	stopCh := make(chan struct{})
	sem := semaphore.NewWeighted(int64(workerLength))
	go func() {
		//发送合并结束的信号标志合并结束
		defer close(mergeDone)
		// 如果直接range,在没有被关闭的情况下,这里必须要先能进来才会退出(如果刚好执行的点就是执行完的时候,就会导致resultChan 一直没有数据进来了)
		// 但是直接用select,则超时的时候 直接通知退出.
		for {
			//判断超时问题,保证消费者退出
			select {
			// 超时的话
			case <-req.Ctx.Done():
				//req.Options.Logger.Infof("消费者退出...\n")
				//消费者退出前通知生产者
				close(stopCh)
				return // 退出消费者
			case res, isOk := <-resultChan: // resultChan 被关闭,直接退出;resultChan 没有被关闭,被pending很长一段时间,也会由于前面的退出
				if !isOk {
					//req.Options.Logger.Infof("resultChan 被关闭...,不再接受数据 \n")
					//被关闭,直接返回
					return
				} else {
					// 正常的话 append
					result = append(result, res)
				}

			}
		}
		// 当 resultChan被关闭时且channel中所有的值都已经被处理完毕后, 将执行到这一行
	}()
	// break 外层for
	// refer: https://blog.csdn.net/u011461385/article/details/106017483
L:
	for i, worker := range workers {
		ctxErr := sem.Acquire(req.Ctx, 1)
		if ctxErr != nil {
			return result, ctxErr
		}
		// 将worker 拿出来
		handler := worker
		index := i
		select {
		case <-req.Ctx.Done(): // 这里意义不是很大
			req.Options.Logger.Infof("task[Id=%s]-%s wait exit[ctx done]...\n", req.Id, req.Desc)
			break L // 退出外循环
		default:
			//默认协程全部启动
			//go func() {
			//	defer func() {
			//		sem.Release(1) //sam信号量+1
			//		<-limiter      // 消耗limit
			//	}()
			//	//运行前放入(也会导致协程阻塞,临时内存泄露)
			//	limiter <- struct{}{}
			//	//执行器
			//	req.waitWork0(index, handler, resultChan, stopCh)
			//}()
			pool.Submit(func() {
				defer func() {
					sem.Release(1) //sam信号量+1
					<-limiter      // 消耗limit
				}()
				//运行前放入(也会导致协程阻塞,临时内存泄露)
				limiter <- struct{}{}
				//执行器
				req.waitWork0(index, handler, resultChan, stopCh)

			})

		}
	}
	// 请求所有的worker,会wait 直到之前的搞完 或者接受一个错误,
	// 这就是 使用这个 而不使用waitGroup的原因,原因就是这里可以接受一个ctx的done事件
	err := sem.Acquire(req.Ctx, int64(workerLength))
	if err != nil {
		//fmt.Println("sem.Acquire Err:", Err)
		//如果发生错误 直接退出(一般是超时)
		return result, err
	}

	//close(limiter) //关闭limiter,这里可以不做,因为不需要通知 gc会帮我们完成
	close(resultChan) //关闭resultChan  这里需要做,不然上面的消费者没办法退出
	<-mergeDone       // 这里主要是等消费完成,防止消费过慢的情况
	cost := time.Since(now).Milliseconds()
	req.Options.Logger.Infof("task[Id=%s]-%s waitWithPreemptive end cost:%d ms\n", req.Id, req.Desc, cost)
	req.Options.FinishHook(req.Ctx) // 结束钩子
	return result, nil

}

// SimpleWaitWithPreemptiveForFirstN 与 WaitWithPreemptive 类似
// 问题是 底下的for循环依赖于上面不阻塞,如果用线程池之类的就会阻塞了
func (req *DefaultParallelWorkReq[R]) SimpleWaitWithPreemptiveForFirstN(firstN int) ([]*ResultUnit[R], error) {
	if firstN == 0 {
		return nil, ErrFirstN
	}
	req.Options.Logger.Infof("task[Id=%s]-%s simpleWaitWithPreemptive start..\n", req.Id, req.Desc)
	now := time.Now() // 当前时间
	//beforeHook
	req.Options.BeforeStartHook(req.Ctx)
	//limiter
	// 限制器
	limiter := make(chan struct{}, req.ParallelSize)
	workers := req.Workers
	workerLength := len(workers)
	// 决定最小的firstN(firstN与workerLength 最小的一个)
	firstN = int(math.Min(float64(firstN), float64(workerLength)))
	// 假设他们能够拿到所有的firstN
	result := make([]*ResultUnit[R], 0, firstN)
	// resultChan 是接受结果的地方,这个地方容易造成生产者和消费者的内存泄漏,所以需要巧妙处理
	//使用 有缓存chan 保证极端情况(在资源充足的情况下,缓存数量与发送者最好一致,这样的话起码能保证生产者不pending)
	resultChan := make(chan *ResultUnit[R], req.ParallelSize)
	// 结果合并结束的信号(如果正常结束,可能消费的比较慢 使用这个通道来强制保证消费结束)
	//mergeDone := make(chan struct{})
	// 使用一个stopCh,来让消费者通知生产者,如果设置result pending了,可以退出了
	stopCh := make(chan struct{})
	//sem := semaphore.NewWeighted(int64(workerLength))
	// break 外层for
	// refer: https://blog.csdn.net/u011461385/article/details/106017483

	//使用一个cancel 快速取消内部的
	ctx, cancelFunc := context.WithCancel(req.Ctx)
	defer cancelFunc()
	// 改变原上下文的指向
	req.Ctx = ctx
	//L:
	for i, worker := range workers {
		//ctxErr := sem.Acquire(req.Ctx, 1)
		//if ctxErr != nil {
		//	return result, ctxErr
		//}
		// 将worker 拿出来
		handler := worker
		index := i
		select {
		default:
			//默认协程全部启动
			go func() {
				defer func() {
					//sem.Release(1) //sam信号量+1
					<-limiter // 消耗limit
				}()
				//运行前放入(也会导致协程阻塞,临时内存泄露)
				limiter <- struct{}{}
				//执行器
				req.waitWork0(index, handler, resultChan, stopCh)
				//fmt.Println("worker 执行完毕...")
			}()
		}
	}
	//req.Options.Logger.Infof("workers end \n")
	// 单消费者 不使用waitgroup手段,其实也可以用waitgroup 限定死firstN个就行
	// 但是取消逻辑再消费者这一点是不可避免的
	// 最大的问题是没有人关闭resultChan,如果resultCHan 阻塞了消费者,就没办法走了(关键点)[用信号量的话就可以解决这个点
	for {
		//判断超时问题,保证消费者退出
		select {
		// 超时的话
		case <-req.Ctx.Done():
			req.Options.Logger.Infof("消费者退出...\n")
			//消费者退出前通知生产者
			close(stopCh)
			return result, req.Ctx.Err()
			//return
		case res := <-resultChan:
			//res 一定不会被关闭
			//req.Options.Logger.Infof("running: %d \n", pool.Running())
			// 正常的话 append
			result = append(result, res)
			//判断是否满足了firstN了
			// 消费者就这一个所以可以直接比较
			if len(result) >= firstN {
				req.Options.Logger.Infof("消费者已接收到数量 %d,退出 \n", firstN)
				close(stopCh) // 退出之后,如果没办法退,使用这个
				cost := time.Since(now).Milliseconds()
				req.Options.Logger.Infof("task[Id=%s]-%s simpleWaitWithPreemptive end cost:%d ms\n", req.Id, req.Desc, cost)
				req.Options.FinishHook(req.Ctx) // 结束钩子
				return result, nil
				//return
			}
		}
	}

}

// WaitWithPreemptiveForFirstN 与 SimpleWaitWithPreemptiveForFirstN 类似,支持池化
// 但是只返回最快的N个(到了N个就退出) N = math.min(N,len(workers)
// 思路为 利用信号量 的特点 读到N个停止
func (req *DefaultParallelWorkReq[R]) WaitWithPreemptiveForFirstN(firstN int) ([]*ResultUnit[R], error) {
	if firstN == 0 {
		return nil, ErrFirstN
	}
	req.Options.Logger.Infof("task[Id=%s]-%s waitWithPreemptiveForFirstN start..\n", req.Id, req.Desc)
	now := time.Now() // 当前时间
	//beforeHook
	req.Options.BeforeStartHook(req.Ctx)
	//limiter
	// 限制器
	limiter := make(chan struct{}, req.ParallelSize)
	workers := req.Workers
	workerLength := len(workers)
	// 决定最小的firstN(firstN与workerLength 最小的一个)
	firstN = int(math.Min(float64(firstN), float64(workerLength)))
	// 假设他们能够拿到所有的firstN
	result := make([]*ResultUnit[R], 0, firstN)
	// resultChan 是接受结果的地方,这个地方容易造成生产者和消费者的内存泄漏,所以需要巧妙处理
	//使用 有缓存chan 保证极端情况(在资源充足的情况下,缓存数量与发送者最好一致,这样的话起码能保证生产者不pending)
	resultChan := make(chan *ResultUnit[R], req.ParallelSize)
	// 结果合并结束的信号(如果正常结束,可能消费的比较慢 使用这个通道来强制保证消费结束)
	//mergeDone := make(chan struct{})
	// 使用一个stopCh,来让消费者通知生产者,如果设置result pending了,可以退出了
	stopCh := make(chan struct{})
	sem := semaphore.NewWeighted(int64(workerLength))
	ctxErr := sem.Acquire(req.Ctx, int64(firstN))
	if ctxErr != nil {
		return result, ctxErr
	}
	//使用一个cancel 快速取消内部的
	ctx, cancelFunc := context.WithCancel(req.Ctx)
	defer cancelFunc() // 快速取消内部
	// 改变原上下文的指向
	req.Ctx = ctx
	//消费者
	go func() {
		//发送合并结束的信号标志合并结束
		//defer close(mergeDone)
		// 如果直接range,在没有被关闭的情况下,这里必须要先能进来才会退出(如果刚好执行的点就是执行完的时候,就会导致resultChan 一直没有数据进来了)
		// 但是直接用select,则超时的时候 直接通知退出.
		for {
			select {
			// 超时的话
			case <-req.Ctx.Done():
				//req.Options.Logger.Infof("消费者退出...\n")
				//消费者退出前通知生产者
				close(stopCh)
				return // 退出消费者
			case res := <-resultChan:
				//resultChan 不会被关闭(关闭一定有问题)
				// 正常的话 append
				result = append(result, res)
				// 在这里release
				sem.Release(1)
				if len(result) >= firstN {
					req.Options.Logger.Infof("消费者已接收到数量 %d,退出 \n", firstN)
					//消费者退出前通知生产者
					close(stopCh)
					return
				}

			}
		}
	}()

L:
	for i, worker := range workers {

		// 将worker 拿出来
		handler := worker
		index := i
		select {
		case <-req.Ctx.Done(): // 这里意义不是很大
			req.Options.Logger.Infof("task[Id=%s]-%s wait exit[ctx done]...\n", req.Id, req.Desc)
			break L // 退出外循环
		default:
			f := func() {
				defer func() {
					//sem.Release(1) //sam信号量+1
					<-limiter // 消耗limit
				}()
				//运行前放入(也会导致协程阻塞,临时内存泄露)
				limiter <- struct{}{}
				//执行器
				req.waitWork0(index, handler, resultChan, stopCh)

			}
			go f()

		}
	}
	//req.Options.Logger.Infof("workers end \n")
	// 请求所有的worker,会wait 直到之前的搞完 或者接受一个错误,
	// 这就是 使用这个 而不使用waitGroup的原因,原因就是这里可以接受一个ctx的done事件
	err := sem.Acquire(req.Ctx, int64(firstN))
	if err != nil {
		//fmt.Println("sem.Acquire Err:", Err)
		//如果发生错误 直接退出(一般是超时)
		return result, err
	}
	//close(resultChan) //关闭resultChan  这里需要做,不然上面的消费者没办法退出
	//<-mergeDone // 这里主要是等消费完成,防止消费过慢的情况
	cost := time.Since(now).Milliseconds()
	req.Options.Logger.Infof("task[Id=%s]-%s waitWithPreemptiveForFirstN end cost:%d ms\n", req.Id, req.Desc, cost)
	req.Options.FinishHook(req.Ctx) // 结束钩子
	return result, nil

}

// gatherWork0 一个gather工作单元
// index 是 当前执行的索引号
// handler 处理单元
// result 需要放入的结果集(因为是gather 所以直接放入结果集可以保证顺序性)
func (req *DefaultParallelWorkReq[R]) gatherWork0(index int, handler *DefaultWorkerArgs[R], result []*ResultUnit[R]) {
	var r R
	resultUnit := &ResultUnit[R]{
		Index:  index,
		Status: Pending,
		Err:    nil,
		Res:    r,
		Worker: handler, // 关联起来
	}
	workerStart := time.Now()
	// 额外处理
	defer func() {
		// panic 捕获
		if panicErr := recover(); panicErr != nil {
			// panic错误处理
			if !req.Options.NotCollectPanic {
				err := errors.New(fmt.Sprintf("%s panic:%v", handler.workerId, panicErr))
				resultUnit.Err = err
				resultUnit.Status = Error
				resultUnit.CostTime = time.Since(workerStart)
				result[index] = resultUnit //塞回去
			} else {
				// panic 处理
				req.Options.Logger.Infof("%s panic: %v\n", handler.workerId, panicErr)
			}
			//panic 钩子
			req.Options.PanicWorkerHook(req.Ctx, panicErr)
		}
	}()
	req.Options.BeforeWorkerStartHook(req.Ctx) // 工作开始钩子
	// 这里的index 是更新原数组的关键
	r, err := handler.worker(req.Ctx)
	if err != nil {
		resultUnit.Err = err
		resultUnit.Status = Error
		resultUnit.Res = r
		resultUnit.CostTime = time.Since(workerStart) // CostTime 记录下来
		result[index] = resultUnit                    // 塞回去
		req.Options.ErrWorkerHook(req.Ctx, err)       // 工作错误钩子
	} else {
		resultUnit.Err = err
		resultUnit.Status = Done
		resultUnit.Res = r
		resultUnit.CostTime = time.Since(workerStart) // CostTime 记录下来
		result[index] = resultUnit                    // 塞回去
	}
	req.Options.FinishWorkerHook(req.Ctx) // 工作结束钩子
}

// waitWork0 一个wait的工作单元
// index 是 当前执行的索引号
// handler 处理单元
// resultChan 需要放入通道,外面用来接受,顺序可能按照时间先后
// stopCh 通知停止机制(消费者通知发送者关闭的信号)
func (req *DefaultParallelWorkReq[R]) waitWork0(index int, handler *DefaultWorkerArgs[R], resultChan chan<- *ResultUnit[R], stopCh <-chan struct{}) {
	var r R
	resultUnit := &ResultUnit[R]{
		Index:  index,
		Status: Pending,
		Err:    nil,
		Res:    r,
		Worker: handler, // 关联起来
	}
	workerStart := time.Now()
	defer func() {
		// panic 捕获
		if panicErr := recover(); panicErr != nil {
			// panic错误处理
			if !req.Options.NotCollectPanic {
				err := errors.New(fmt.Sprintf("%s panic:%v", handler.workerId, panicErr))
				resultUnit.Err = err
				resultUnit.Status = Error
				resultUnit.CostTime = time.Since(workerStart)
				// 如果超时了,那其实上就不应该放到结果集里面,因为它不被需要了
				select {
				case <-req.Ctx.Done():
					return
				case <-stopCh: // stop直接返回(消费者通知的)
					return
				case resultChan <- resultUnit:
				}
			} else {
				// panic 处理
				req.Options.Logger.Infof("%s panic: %v\n", handler.workerId, panicErr)
			}
			//panic 钩子
			req.Options.PanicWorkerHook(req.Ctx, panicErr)

		}
	}()
	req.Options.BeforeWorkerStartHook(req.Ctx) // 工作开始钩子
	// 这里的index 是更新原数组的关键
	r, err := handler.worker(req.Ctx)
	if err != nil {
		resultUnit.Err = err
		resultUnit.Status = Error
		resultUnit.Res = r
		resultUnit.CostTime = time.Since(workerStart)
		select {
		// 如果超时了,那其实上就不应该放到结果集里面,因为它不被需要了
		case <-stopCh: // stop直接返回
			return
		case <-req.Ctx.Done():
			return
		case resultChan <- resultUnit: //超时的时候也能进来,如果能放就放,但是不应该从外侧关闭,因为case 是并发执行的

		}
		req.Options.ErrWorkerHook(req.Ctx, err) // 工作错误钩子
	} else {
		resultUnit.Err = err
		resultUnit.Status = Done
		resultUnit.Res = r
		resultUnit.CostTime = time.Since(workerStart)
		select {
		// 如果超时了,那其实上就不应该放到结果集里面,因为它不被需要了
		case <-req.Ctx.Done():
			return
		case <-stopCh: // stop直接返回
			return
		case resultChan <- resultUnit:
		}
	}
	req.Options.FinishWorkerHook(req.Ctx) // 工作结束钩子
}

func (res *ResultUnit[R]) IsError() bool {
	return res.Status == Error
}

func (res *ResultUnit[R]) Result() R {
	return res.Res
}

func (res *ResultUnit[R]) Error() error {
	return res.Err
}

func (res *ResultUnit[R]) GetIndex() int {
	return res.Index
}

func (res *ResultUnit[R]) IsDone() bool {
	return res.Status == Done
}

func (res *ResultUnit[R]) GetCostTime() time.Duration {
	return res.CostTime
}

func (res *ResultUnit[R]) GetWorker() *DefaultWorkerArgs[R] {
	return res.Worker
}
