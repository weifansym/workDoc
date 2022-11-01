接下来会一共会有 12 - 15 篇文章讲解 Go 并发编程，并发编程本身是一个挺大的话题，在第四周的两节课，毛老师花了将近 7 个小时讲解这些内容，我也结合自己的一些微不足道的经验，再加上一些
大神们的文章，整理出了这一部分的笔记。
当然这里更多的是抛砖引玉的作用，更多的还是我们自己要有相关的意识避免踩坑，在各个坑的边缘反复横跳，可能我们有缘会在同一个坑中发现，咦，原来你也在这里 😄

### 请对你创建的 goroutine 负责
**不要创建一个你不知道何时退出的 goroutine**
请阅读下面这段代码，看看有什么问题？
> 为什么先从下面这段代码出发，是因为在之前的经验里面我们写了大量类似的代码，之前没有意识到这个问题，并且还因为这种代码出现过短暂的事故
```
// Week03/blog/01/01.go
package main

import (
	"log"
	"net/http"
	_ "net/http/pprof"
)

func setup() {
	// 这里面有一些初始化的操作
}

func main() {
	setup()

	// 主服务
	server()

	// for debug
	pprof()

	select {}
}

func server() {
	go func() {
		mux := http.NewServeMux()
		mux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("pong"))
		})

		// 主服务
		if err := http.ListenAndServe(":8080", mux); err != nil {
			log.Panicf("http server err: %+v", err)
			return
		}
	}()
}

func pprof() {
	// 辅助服务，监听了其他端口，这里是 pprof 服务，用于 debug
	go http.ListenAndServe(":8081", nil)
}

```
灵魂拷问来了，请问：
* 如果 server  是在其他包里面，如果没有特殊说明，你知道这是一个异步调用么？
* main  函数当中最后在哪里空转干什么？会不会存在浪费？
* 如果线上出现事故，debug 服务已经退出，你想要 debug 这时你是否很茫然？
* 如果某一天服务突然重启，你却找不到事故日志，你是否能想起这个 8081  端口的服务？

**请将选择权留给对方，不要帮别人做选择**
请把是否并发的选择权交给你的调用者，而不是自己就直接悄悄的用上了 goroutine
下面这次改动将两个函数是否并发操作的选择权留给了 main 函数
```
package main

import (
	"log"
	"net/http"
	_ "net/http/pprof"
)

func setup() {
	// 这里面有一些初始化的操作
}

func main() {
	setup()

	// for debug
	go pprof()

	// 主服务
	go server()

	select {}
}

func server() {
	mux := http.NewServeMux()
	mux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})

	// 主服务
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Panicf("http server err: %+v", err)
		return
	}
}

func pprof() {
	// 辅助服务，监听了其他端口，这里是 pprof 服务，用于 debug
	http.ListenAndServe(":8081", nil)
}
```
**请不要作为一个旁观者**

一般情况下，不要让主进程成为一个旁观者，明明可以干活，但是最后使用了一个 select  在那儿空跑
感谢上一步将是否异步的选择权交给了我( main )，在旁边看着也怪尴尬的
```
package main

import (
	"log"
	"net/http"
	_ "net/http/pprof"
)

func setup() {
	// 这里面有一些初始化的操作
}

func main() {
	setup()

	// for debug
	go pprof()

	// 主服务
	server()
}

func server() {
	mux := http.NewServeMux()
	mux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})

	// 主服务
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Panicf("http server err: %+v", err)
		return
	}
}

func pprof() {
	// 辅助服务，监听了其他端口，这里是 pprof 服务，用于 debug
	http.ListenAndServe(":8081", nil)
}
```
**不要创建一个你永远不知道什么时候会退出的 goroutine**
我们再做一些改造，使用 channel  来控制，解释都写在代码注释里面了
```
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"time"
)

func setup() {
	// 这里面有一些初始化的操作
}

func main() {
	setup()

	// 用于监听服务退出
	done := make(chan error, 2)
	// 用于控制服务退出，传入同一个 stop，做到只要有一个服务退出了那么另外一个服务也会随之退出
	stop := make(chan struct{}, 0)
	// for debug
	go func() {
		done <- pprof(stop)
	}()

	// 主服务
	go func() {
		done <- app(stop)
	}()

	// stoped 用于判断当前 stop 的状态
	var stoped bool
	// 这里循环读取 done 这个 channel
	// 只要有一个退出了，我们就关闭 stop channel
	for i := 0; i < cap(done); i++ {
		if err := <-done; err != nil {
			log.Printf("server exit err: %+v", err)
		}

		if !stoped {
			stoped = true
			close(stop)
		}
	}
}

func app(stop <-chan struct{}) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})

	return server(mux, ":8080", stop)
}

func pprof(stop <-chan struct{}) error {
	// 注意这里主要是为了模拟服务意外退出，用于验证一个服务退出，其他服务同时退出的场景
	go func() {
		server(http.DefaultServeMux, ":8081", stop)
	}()

	time.Sleep(5 * time.Second)
	return fmt.Errorf("mock pprof exit")
}

// 启动一个服务
func server(handler http.Handler, addr string, stop <-chan struct{}) error {
	s := http.Server{
		Handler: handler,
		Addr:    addr,
	}

	// 这个 goroutine 我们可以控制退出，因为只要 stop 这个 channel close 或者是写入数据，这里就会退出
	// 同时因为调用了 s.Shutdown 调用之后，http 这个函数启动的 http server 也会优雅退出
	go func() {
		<-stop
		log.Printf("server will exiting, addr: %s", addr)
		s.Shutdown(context.Background())
	}()

	return s.ListenAndServe()
}
```
我们看一下返回结果，这个代码启动 5s 之后就会退出程序
```
❯ go run ./01_goroutine/04
2020/12/08 21:49:43 server exit err: mock pprof exit
2020/12/08 21:49:43 server will exiting, addr: :8081
2020/12/08 21:49:43 server will exiting, addr: :8080
2020/12/08 21:49:43 server exit err: http: Server closed
```
**思考题**
虽然我们已经经过了三轮优化，但是这里还是有一些需要注意的地方，可以思考一下怎么做
* 虽然我们调用了 Shutdown  方法，但是我们其实并没有实现优雅退出，相信聪明的你可以完成这项工作。可以参考上一篇笔记：[Go 错误处理最佳实践](https://lailin.xyz/post/go-training-03.html)
* 在 server  方法中我们并没有处理 panic  的逻辑，这里需要处理么？如果需要那该如何处理呢？

**不要创建一个永远都无法退出的 goroutine [goroutine 泄漏]**
再来看下面一个例子，这也是常常会用到的操作
```
func leak(w http.ResponseWriter, r *http.Request) {
	ch := make(chan bool, 0)
	go func() {
		fmt.Println("异步任务做一些操作")
		<-ch
	}()

	w.Write([]byte("will leak"))
}
```
复用一下上面的 server 代码，我们经常会写出这种类似的代码
* http 请求来了，我们启动一个 goroutine 去做一些耗时一点的工作
* 然后返回了
* 然后之前创建的那个**goroutine 阻塞了**
* 然后就泄漏了

绝大部分的 goroutine 泄漏都是因为 goroutine 当中因为各种原因阻塞了，我们在外面也没有控制它退出的方式，所以就泄漏了，具体导致阻塞的常见原因会在接下来的 sync 包、channel 中讲到，
这里就不过多赘述了
接下来我们验证一下是不是真的泄漏了
启动之后我们访问一下: http://localhost:8081/debug/pprof/goroutine?debug=1 查看当前的 goroutine 个数为 7
```
goroutine profile: total 7
2 @ 0x43b945 0x40814f 0x407d8b 0x770998 0x470381
#	0x770997	main.server.func1+0x37	/home/ll/project/Go-000/Week03/blog/01_goroutine/05/05.go:71
```
然后我们再访问几次 http://localhost:8080/leak 可以发现 goroutine 增加到了 15 个，而且一直不会下降
```
goroutine profile: total 15
7 @ 0x43b945 0x40814f 0x407d8b 0x770ad0 0x470381
#	0x770acf	main.leak.func1+0x8f	/home/ll/project/Go-000/Week03/blog/01_goroutine/05/05.go:83
```
**确保创建出的 goroutine 的工作已经完成**
这个其实就是优雅退出的问题，我们可能启动了很多的 goroutine 去处理一些问题，但是服务退出的时候我们并没有考虑到就直接退出了。例如退出前日志没有 flush 到磁盘，我们的请求还没完全关闭，
异步 worker 中还有 job 在执行等等。
我们也来看一个例子，假设现在有一个埋点服务，每次请求我们都会上报一些信息到埋点服务上

```
// Reporter 埋点服务上报
type Reporter struct {
}

var reporter Reporter

// 模拟耗时
func (r Reporter) report(data string) {
	time.Sleep(time.Second)
	fmt.Printf("report: %s\n", data)
}

mux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
    // 在请求中异步调用
    go reporter.report("ping pong")
    fmt.Println("ping")
    w.Write([]byte("pong"))
})
```
我在发送了一次请求之后直接退出了，异步上报的逻辑根本没执行上
```
❯ go run ./01_goroutine/06
ping
^Csignal: interrupt
```
这个有两种改法，一种是给 reporter 加上 shutdown 方法，类似 http 的 shutdown，等待所有的异步上报完成之后，我们再退出，另外一种是我们直接使用 一些 worker 来执行，在当然这个 
worker 也要实现类似 shutdown 的方法。一般推荐后一种，因为这样可以避免请求量比较大时，创建大量 goroutine，当然如果请求量比较小，不会很大，用第一种也是可以的。

我们给一个第二种的简单实现，第一种可以参考 https://www.ardanlabs.com/blog/2019/04/concurrency-trap-2-incomplete-work.html

```
// Reporter 埋点服务上报
type Reporter struct {
	worker   int
	messages chan string
	wg       sync.WaitGroup
	closed   bool
}

// NewReporter NewReporter
func NewReporter(worker, buffer int) *Reporter {
	return &Reporter{worker: worker, messages: make(chan string, buffer)}
}

func (r *Reporter) run(stop <-chan struct{}) {
	go func() {
		<-stop
		r.shutdown()
	}()

	for i := 0; i < r.worker; i++ {
		r.wg.Add(1)
		go func() {
			for msg := range r.messages {
				time.Sleep(5 * time.Second)
				fmt.Printf("report: %s\n", msg)
			}
			r.wg.Done()
		}()
	}
	r.wg.Wait()
}

func (r *Reporter) shutdown() {
	r.closed = true
	// 注意，这个一定要在主服务结束之后再执行，避免关闭 channel 还有其他地方在啊写入
	close(r.messages)
}

// 模拟耗时
func (r *Reporter) report(data string) {
	if r.closed {
		return
	}
	r.messages <- data
}
```
然后在 main 函数中我们加上
```
go func() {
    reporter.run(stop)
    done <- nil
}()
```
> 留一个思考题：我们在 reporter 的实现可能会导致 panic，你是否发现了呢？如何修改可以避免这种情况？ 感谢评论区 @hddxds 的指出，我这里给出一个实现例子:
>  [点击查看](https://github.com/mohuishou/Go-000/blob/main/Week03/blog/01_goroutine/07/reporter.go)，可以看看是否和你想的一样？ 如果你对为什么会出现 panic 
>  或者为什么要这么实现感到困惑可以查看后面的这篇文章 [Go并发编程(十) 深入理解 Channel](https://lailin.xyz/post/go-training-week3-channel.html)

### 总结
总结一下这一部分讲到的几个要点，这也是我们
1. **请将是否异步调用的选择权交给调用者**，不然很有可能大家并不知道你在这个函数里面使用了 goroutine
2. 如果你要启动一个 goroutine 请对它负责
* **永远不要启动一个你无法控制它退出，或者你无法知道它何时推出的 goroutine**
* 还有上一篇提到的，启动 goroutine 时请加上 panic recovery 机制，避免服务直接不可用
* 造成 goroutine 泄漏的主要原因就是 goroutine 中造成了阻塞，并且没有外部手段控制它退出
3. **尽量避免在请求中直接启动 goroutine 来处理问题**，而应该通过启动 worker 来进行消费，这样可以避免由于请求量过大，而导致大量创建 goroutine 从而导致 oom，当然如果请求量本身非常小，那当我没说

### 参考文献
* https://dave.cheney.net/practical-go/presentations/qcon-china.html 这篇 dave 在 Qcon China 上的文章值得好好拜读几遍
* https://www.ardanlabs.com/blog/2018/11/goroutine-leaks-the-forgotten-sender.html
* https://www.ardanlabs.com/blog/2019/04/concurrency-trap-2-incomplete-work.html
* https://www.ardanlabs.com/blog/2014/01/concurrency-goroutines-and-gomaxprocs.html

转自：https://lailin.xyz/post/go-training-week3-goroutine.html



