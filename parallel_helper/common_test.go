package parallel_helper

import (
	"common-tools/collection_helper/commonArrOperation"
	"common-tools/convert_helper"
	"context"
	"fmt"
	"github.com/imroc/req/v3"
	"github.com/panjf2000/ants/v2"
	"io"
	"net/http"
	_ "net/http/pprof"
	"runtime"
	"sync"
	"testing"
	"time"
)

/*
go leak:
https://daryeon.github.io/post/go-library-goleak/
可以用来搞内存泄漏
*/

var (
	client *req.Client
)

func init() {

	//go func() {
	//	logger.Println(http.ListenAndServe(":6060", nil))
	//}()
	client = req.C().
		DisableAutoDecode()
}

func TestDefaultParallelWorkReq_Wait(t *testing.T) {

	workers := make([]DefaultWorker[string], 0)
	for i := 0; i < 10000; i++ {
		index := i
		handler := func(ctx context.Context) (string, error) {
			return GetBangumiSubject(ctx, index)
			//time.Sleep(300 * time.Millisecond) // 就假设一个请求平均300ms
			//fmt.Println(index)
			//return convert_helper.AnyToString(index), nil
		}
		workers = append(workers, handler)
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	req, _ := NewDefaultParallelWorkReq[string](ctx, workers,
		WithParallelSize(100),
		WithIsMixed(false),
	)
	fmt.Printf("%+v\n", *req)
	before := runtime.NumGoroutine()

	fmt.Println("NumGoroutine:", before)
	//result, Err := req.Wait()
	//if Err != nil {
	//	fmt.Println("Err", Err)
	//}
	result, err := req.WaitWithPreemptive()
	if err != nil {
		fmt.Println(err)
	}
	for i, k := range result {
		fmt.Println(i, k.Worker)
	}
	fmt.Println(len(result), cap(result))
	//time.Sleep(5 * time.Second)
	//beforeTime := time.Now()
	//for {
	//	now := runtime.NumGoroutine()
	//	fmt.Println("NumGoroutine:", now)
	//	if now == before {
	//		fmt.Infof("恢复正常的协程数花费了: %f s\n", time.Since(beforeTime).Seconds())
	//		break
	//	}
	//	time.Sleep(2 * time.Second)
	//
	//}

	//o := commonArrOperation.NewArrComparableOperation[*ResultUnit[string]](&result)
	//o.Each(func(unit *ResultUnit[string]) bool {
	//	//// 如果超时退出可能为null
	//	////fmt.Println(unit)
	//	//if unit.IsError() {
	//	//	fmt.Println(unit.Err)
	//	//}
	//	return false
	//})
	//for {
	//
	//}
}

func TestDefaultParallelWorkReq_WaitWithPreemptiveForFirstN(t *testing.T) {
	//defer goleak.VerifyNone(t)

	workers := make([]DefaultWorker[string], 0)
	for i := 0; i < 1000000; i++ {
		index := i
		handler := func(ctx context.Context) (string, error) {
			return GetBangumiSubject(ctx, index)
			//time.Sleep(300 * time.Millisecond) // 就假设一个请求平均300ms
			////fmt.Println(index)
			//return "", nil

			//select {
			//case <-ctx.Done():
			//	return "", ctx.Err()
			//case <-time.After(300 * time.Millisecond):
			//	return convert_helper.AnyToString(index), nil
			//}

		}
		workers = append(workers, handler)
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()
	req, _ := NewDefaultParallelWorkReq[string](ctx, workers,
		WithParallelSize(100),
		WithIsMixed(false),
	)
	//fmt.Infof("%+v\n", *req)
	before := runtime.NumGoroutine()

	fmt.Println("NumGoroutine:", before)
	//result, Err := req.Wait()
	//if Err != nil {
	//	fmt.Println("Err", Err)
	//}
	result, err := req.WaitWithPreemptiveForFirstN(10)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(len(result), cap(result))
	for _, item := range result {
		fmt.Println(item.GetIndex(), item.GetCostTime(), item.GetWorker())
	}
	beforeTime := time.Now()
	for {
		now := runtime.NumGoroutine()
		fmt.Println("NumGoroutine:", now)
		if now == before {
			fmt.Printf("恢复正常的协程数花费了: %f s\n", time.Since(beforeTime).Seconds())
			break
		}
		time.Sleep(2 * time.Second)
	}

	//o := commonArrOperation.NewArrComparableOperation[*ResultUnit[string]](&result)
	//o.Each(func(unit *ResultUnit[string]) bool {
	//	//// 如果超时退出可能为null
	//	////fmt.Println(unit)
	//	//if unit.IsError() {
	//	//	fmt.Println(unit.Err)
	//	//}
	//	return false
	//})
	//for {
	//
	//}
}

func TestAnts(t *testing.T) {
	newPool, err := ants.NewPool(5, ants.WithNonblocking(false))
	if err != nil {
		panic(err)
	}
	ctx := context.TODO()
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		index := i
		newPool.Submit(func() {
			defer wg.Done()
			GetBangumiSubject(ctx, index)
		})
	}
	wg.Wait()
	fmt.Println("finish...")
}

func TestDefaultParallelWorkReq_Gather2(t *testing.T) {
	workers := make([]DefaultWorker[int], 0)
	for i := 0; i < 100; i++ {
		index := i
		handler := func(ctx context.Context) (int, error) {
			fmt.Println(index)
			return 0, nil
		}
		workers = append(workers, handler)
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*10)
	defer cancel()
	req, _ := NewDefaultParallelWorkReq[int](ctx, workers)
	//fmt.Println(req)
	result := req.Gather()
	o := commonArrOperation.NewArrComparableOperation[*ResultUnit[int]](&result)
	o.Each(func(unit *ResultUnit[int]) bool {
		fmt.Println(unit)
		return false
	})

}

func TestGetBangumiSubject(t *testing.T) {
	fmt.Println(runtime.NumGoroutine())
	fmt.Println(GetBangumiSubject(context.Background(), 1))
	fmt.Println(runtime.NumGoroutine())
	//fmt.Println(GetBangumiSubject2(context.Background(), 1))
}

/*
*
http内存泄露的问题:
https://wanghe4096.github.io/2019/03/13/avoiding-memory-leak-in-golang-api/
*/
func GetBangumiSubject(ctx context.Context, index int) (string, error) {
	url := fmt.Sprintf("%s/v0/subjects/%d", "https://api.bgm.tv", index)
	//fmt.Infof("request %s \n", url)
	get, err := client.R().
		SetContext(ctx).
		Get(url)
	if get != nil {
		defer get.Body.Close() // MUST CLOSED THIS
	}
	if err != nil {
		return "", err
	}

	return "", err
	//s := get.String()
	//return s, nil
}

//	func TestMain(m *testing.M) {
//		goleak.VerifyTestMain(m)
//	}
func GetBangumiSubject2(ctx context.Context, index int) (string, error) {
	url := fmt.Sprintf("%s/v0/subjects/%d", "https://api.bgm.tv", index)
	fmt.Printf("request %s \n", url)
	withContext, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return "", err
	}
	client2 := &http.Client{}
	resp, err := client2.Do(withContext)
	if resp != nil {
		defer resp.Body.Close() // MUST CLOSED THIS
	}
	if err != nil {
		return "", err
	}
	all, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return convert_helper.BytesToStr(all), nil

}
