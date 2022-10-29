package test

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"runtime"
	"strconv"
	"sync"
	"testing"
	"time"
)

/*
refer: https://www.51cto.com/article/710755.html?u_atoken=540cfbf9-7533-487f-8cc4-7ccc1898f3cc&u_asession=01i-6Fk-YFP2hn6IjKFqmDbuMttkTfT3eGa6GBq0IXK221HfJzI_GRI5KYxyPrKM1cX0KNBwm7Lovlpxjd_P_q4JsKWYrT3W_NKPr8w6oU7K9jb2mrN1qLTwlPPXgXNSgThUF3o-sVtq6Wun3JL3SJe2BkFo3NEHBv0PZUm6pbxQU&u_asig=05qFfvfDNOQSPq9NIuzj4ViQmsRR2ujKesrZ13A9CHbEbwO8yV_xLdQf42TVvgKbPEoSR8E8atMH7s6op0sAefwR1VOKD7UeW_0AoX-3hJAOFFIWDBjF5TeF6yo818wrSvh0CRxCvy8Ckq2X8ynTEnmXuAN5G39yKk7Ob825rhAef9JS7q8ZD7Xtz2Ly-b0kmuyAKRFSVJkkdwVUnyHAIJzXLtJqPCZzLkDe3e5xvAYYNu4IgZ4yDiqg8SSeRqVLUeom7nzSzR1LP16f45fIKp-e3h9VXwMyh6PgyDIVSG1W8-OhX6VCKnU-cR_oFm7_4ebUEz969bgik9l49-8PObqGaHsG6xa_FKN2x9iRrWoXn475ZYb4kXV2Qwe4VdfkktmWspDxyAEEo4kbsryBKb9Q&u_aref=DjjZ%2ByKsLZammMuFPzusBwnBGfs%3D
refer: https://juejin.cn/post/7033671944587182087
refer: https://juejin.cn/post/7033711399041761311
refer: https://learnku.com/go/t/23459/how-to-close-the-channel-gracefully
refer: https://www.cnblogs.com/gwyy/p/13629999.html 关于select case
预防 goroutine 泄漏的核心就是：

创建 goroutine 时就要想清楚它什么时候被回收。(重点)
具体到执行层面，包括：

当 goroutine 退出时，需要考虑它使用的 channel 有没有可能阻塞对应的生产者、消费者的 goroutine。
尽量使用buffered channel使用buffered channel 能减少阻塞发生、即使疏忽了一些极端情况，也能降低 goroutine 泄漏的概率。

关不关闭chan:
结论：除非必须关闭 chan，否则不要主动关闭。关闭 chan 最优雅的方式，就是不要关闭 chan~。

当一个 chan 没有 sender 和 receiver 时，即不再被使用时，GC 会在一段时间后标记、清理掉这个 chan。
那么什么时候必须关闭 chan 呢？
比较常见的是将 close 作为一种通知机制，尤其是生产者与消费者之间是 1:M 的关系时，通过 close 告诉下游：我收工了，你们别读了。(比如for循环结束)

关闭chan的场合:
1. 不要在消费者端关闭 chan,如果要关闭,需要先同步生产者,在关闭(否则会因为chan导致pending)

一写一读：生产者关闭即可。其实就类似于通知消费者,没有更多的生产了.
一写多读：生产者关闭即可，关闭时下游全部消费者都能收到通知。 其实就类似于通知消费者,没有更多的生产了.
多写一读：多个生产者之间需要引入一个协调 channel 来处理信号。这个就是下面的做法,需要保证消费者有效或者消费者关闭前,要保证接受者都关闭.需要这么一个协调的channel
多写多读：与 3 类似，核心思路是引入一个中间层以及使用try-send 的套路来处理非阻塞的写入
(多个发送者，多个接收者：任意一方使用专用的 stop channel 关闭；发送者、接收者都使用 select 监听 stop channel 是否关闭。)

有没有办法尽可能的减少生产者的pending?
1. 使用有缓存通道 通道数量 >= 发送者数量
这样即使没有消费者,也不会pending
2. 使用select 发送,给出default,
如果发送阻塞了说明缺乏接受者(也能说明接受者可能处理不过来)所以这个方法还不是最好
但是通过这个方式确实可以保证 生产者的退出 (refer:https://geektutu.com/post/hpg-timeout-goroutine.html)

*/

func TestLeakOfMemory(t *testing.T) {
	fmt.Println("NumGoroutine:", runtime.NumGoroutine())
	/*
						原因:
						1. 使用无缓存channel, 一个放必须由一个取是特点.
					      发送者的发送操作将阻塞，直到接收者执行接受操作。同样接受者的接受操作将阻塞，直到发送者执行发送操作。发送者的发送操作和接受者的接受操作是同步的。
					    2.下面 在3 处超时后,2处执行完成,往errCh通道中发送元素,但是因为上述1的特点,所以 errCh 没有被人接受的情况下,会一直卡主,造成死锁
				        3. 但是由于主程序已经返回了,所以造成了死锁的协程泄露了.
			            解决思路:
		               改用有缓存的channel,同样超时之后 // (3) 之后 (2) 执行,会往errCh中发数据,errCh有缓存(至少1个,size >= 参与这事的协程数) 那么通道就能接收到数据(当然通道这个时候也不能关闭)
	*/
	chanLeakOfMemory()
	time.Sleep(time.Second * 3) // 等待 goroutine 执行，防止过早输出结果
	fmt.Println("NumGoroutine:", runtime.NumGoroutine())
}

func chanLeakOfMemory() {
	errCh := make(chan error, 10) // (1)
	go func() {                   // (5)
		time.Sleep(2 * time.Second)
		errCh <- errors.New("chan error") // (2) 如果超时退出,这里就一直走不下去了,因为
		fmt.Println("finish sending")     //(6)
	}()

	var err error
	select {
	case <-time.After(time.Second): // (3) 大家也经常在这里使用 <-ctx.Done()
		fmt.Println("超时") // 退出之后,errCh 都没有接收到结果
	case err = <-errCh: // (4)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(nil)
		}
	}
}

func TestLeakOfMemory2(t *testing.T) {
	fmt.Println("NumGoroutine:", runtime.NumGoroutine())
	/*
						原因:
						1. 使用有缓存的通道,但是缓存通道数量不够,接受者没有消耗完,就退出了,那么再往里面放就阻塞了,状况类似于1
					    解决思路:
					     1. 上面的原因主要是因为缓存数量 <了发送者,就会导致,需要依赖接受者,接受者退出后,由于有缓存通道得不到及时的消费,所以后续的发送者就会pending住, 那么最简单的思路就是
				          缓存队列的容量需要和发送次数一致. chanLeakOfMemory3,这种情况有几个不好的点,首先你通道得比较长(占内存),其次你还是需要让生产者完全走完,所以时间上停止的也更慢
			             2. 采用一个stopChan,通过这个stopChan,改变生产者的行为,让他后续不会再向通道中发送数据了,这样做也不需要很长的缓存通道 chanLeakOfMemory2
		                 3. 利用gc来回收这个不再被使用的chan
	*/
	// 使用stop 不让发送者再发送
	//chanLeakOfMemory2()
	// 让有缓存通道与发送者数量保证一致
	chanLeakOfMemory3()
	time.Sleep(time.Second * 15) // 等待 goroutine 执行，防止过早输出结果(如果是 有缓存通道与发送者数量保持一致的话,那么这个就相当于让生产者走完,会比较长)
	fmt.Println("NumGoroutine:", runtime.NumGoroutine())
}

func chanLeakOfMemory2() {
	ich := make(chan int, 100)    // (3)
	stopCh := make(chan struct{}) // 采用一个stopCh 来合理的退出发送者
	// sender
	go func() {
		defer close(ich)
		for i := 0; i < 10000; i++ {
			//ich <- i
			select {
			case <-stopCh:
				return
			case ich <- i: // 这步可能阻塞
			}
			time.Sleep(time.Millisecond) // 控制一下，别发太快
		}
	}()
	// receiver
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		for i := range ich { // (2)
			if ctx.Err() != nil { // (1) 这里的超时判断会导致receiver提前退出,receiver提前退出,那么chan 100 也是不够的
				fmt.Println(ctx.Err())
				close(stopCh) // 当超时退出的时候,直接退出发送者
				return
			}
			fmt.Println(i)
		}
	}()
}

// chanLeakOfMemory3 让缓存通道的数量与发送者数量保持一致
func chanLeakOfMemory3() {
	n := 10000
	ich := make(chan int, n) // (3)
	// sender
	go func() {
		defer close(ich)
		for i := 0; i < n; i++ {
			ich <- i
			time.Sleep(time.Millisecond) // 控制一下，别发太快
		}
	}()
	// receiver
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		for i := range ich { // (2)
			if ctx.Err() != nil { // (1) 这里的超时判断会导致receiver提前退出,
				fmt.Println(ctx.Err())
				return
			}
			fmt.Println(i)
		}
	}()
}

func TestLeakOfMemory3(t *testing.T) {
	fmt.Println("NumGoroutine:", runtime.NumGoroutine())
	/*
						原因:
						1. 使用有缓存的通道,但是缓存通道数量不够,接受者没有消耗完,就退出了,那么再往里面放就阻塞了,状况类似于1
					    解决思路:
					     1. 上面的原因主要是因为缓存数量 <了发送者,就会导致,需要依赖接受者,接受者退出后,由于有缓存通道得不到及时的消费,所以后续的发送者就会pending住, 那么最简单的思路就是
				          缓存队列的容量需要和发送次数一致. chanLeakOfMemory3,这种情况有几个不好的点,首先你通道得比较长(占内存),其次你还是需要让生产者完全走完,所以时间上停止的也更慢
			             2. 采用一个stopChan,通过这个stopChan,改变生产者的行为,让他后续不会再向通道中发送数据了,这样做也不需要很长的缓存通道 chanLeakOfMemory2
		                 3. 利用gc来回收这个不再被使用的chan
	*/
	// 使用stop 不让发送者再发送
	//chanLeakOfMemory2()
	// 让有缓存通道与发送者数量保证一致
	manyReadManyWrite()
	time.Sleep(time.Second * 15) // 等待 goroutine 执行，防止过早输出结果(如果是 有缓存通道与发送者数量保持一致的话,那么这个就相当于让生产者走完,会比较长)
	fmt.Println("NumGoroutine:", runtime.NumGoroutine())
}

// manyReadManyWrite  多个生产者,多个消费者
func manyReadManyWrite() {
	rand.Seed(time.Now().UnixNano())
	log.SetFlags(0)

	const Max = 100000
	const NumReceivers = 10
	const NumSenders = 1000

	wgReceivers := sync.WaitGroup{}
	wgReceivers.Add(NumReceivers)

	dataCh := make(chan int)
	stopCh := make(chan struct{})
	// stopCh 是额外引入的一个信号 channel.
	// 它的生产者是下面的 toStop channel，
	// 消费者是上面 dataCh 的生产者和消费者
	toStop := make(chan string, 1)
	// toStop 是拿来关闭 stopCh 用的，由 dataCh 的生产者和消费者写入
	// 由下面的匿名中介函数(moderator)消费
	// 要注意，这个一定要是 buffered channel （否则没法用 try-send 来处理了）

	var stoppedBy string

	// moderator
	go func() {
		stoppedBy = <-toStop
		close(stopCh)
	}()

	// senders
	for i := 0; i < NumSenders; i++ {
		go func(id string) {
			for {
				value := rand.Intn(Max)
				if value == 0 {
					// try-send 操作
					// 如果 toStop 满了，就会走 default 分支啥也不干，也不会阻塞
					select {
					case toStop <- "sender#" + id:
					default:
					}
					return
				}

				// try-receive 操作，尽快退出
				// 如果没有这一步，下面的 select 操作可能造成 panic
				select {
				case <-stopCh:
					return
				default:
				}

				// 如果尝试从 stopCh 取数据的同时，也尝试向 dataCh
				// 写数据，则会命中 select 的伪随机逻辑，可能会写入数据
				select {
				case <-stopCh:
					return
				case dataCh <- value:
				}
			}
		}(strconv.Itoa(i))
	}

	// receivers
	for i := 0; i < NumReceivers; i++ {
		go func(id string) {
			defer wgReceivers.Done()

			for {
				// 同上
				select {
				case <-stopCh:
					return
				default:
				}

				// 尝试读数据
				select {
				case <-stopCh:
					return
				case value := <-dataCh:
					if value == Max-1 {
						select {
						case toStop <- "receiver#" + id:
						default:
						}
						return
					}

					log.Println(value)
				}
			}
		}(strconv.Itoa(i))
	}

	wgReceivers.Wait()
	log.Println("stopped by", stoppedBy)
}
