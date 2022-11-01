## Go并发编程(三) data race
### 回顾
在前两篇文章当中我们反复提到了虽然在 go 中，并发编程十分简单，我们只需要使用 go func()  就能启动一个 goroutine 去做一些事情，但是正是由于这种简单我们要十分当心，不然很容易出现一
些莫名其妙的 bug 或者是你的服务由于不知名的原因就重启了。

### 数据竞争(data race)
之前我们提到了很多次在多个 goroutine 对同一个变量的数据进行修改的时候会出现很多奇奇怪怪的问题，那我们有没有什么办法检测它呢，除了通过我们聪明的脑袋？

答案就是 data race tag，go 官方早在 1.1 版本就引入了数据竞争的检测工具，我们只需要在执行测试或者是编译的时候加上 -race  的 flag 就可以开启数据竞争的检测
```
go test -race ./...
go build -race
```
> 不建议在生产环境 build 的时候开启数据竞争检测，因为这会带来一定的性能损失(一般内存5-10倍，执行时间2-20倍)，当然 必须要 debug 的时候除外。 建议在执行单元测试时始终开启数据竞争
> 的检测。

#### 案例一
我们来直接看一下下面的这个例子，这是来自课上的一个例子，但是我稍稍做了一些改造，源代码没有跑 10w 次这个操作，会导致看起来每次跑的结果都是差不多的，我们只需要把这个次数放大就可以发现
每次结果都会不一样

**正常执行**
```
package main

import (
	"fmt"
	"sync"
)

var wg sync.WaitGroup
var counter int

func main() {
	// 多跑几次来看结果
	for i := 0; i < 100000; i++ {
		run()
	}
}

func run() {
	for i := 1; i <= 2; i++ {
		wg.Add(1)
		go routine(i)
	}
	wg.Wait()
	fmt.Printf("Final Counter: %d\n", counter)
}

func routine(id int) {
	for i := 0; i < 2; i++ {
		value := counter
		value++
		counter = value
	}
	wg.Done()
}
```
我执行了三次每次的结果都不一致，分别是:
```
Final Counter: 399996
Final Counter: 399989
Final Counter: 399988
```
为什么会导致这样的结果呢，是因为每一次执行的时候，我们都使用 go routine(i)  启动了两个 goroutine，但是我们并没有控制它的执行顺序，那就有好几种可能了，我这里描述两种情况
1. 执行一次 run() , counter + 4 这种情况下，第二个 goroutine 开始执行时，拿到了第一个 goroutine 的执行结果，也就是 value := counter  这一步时，value = 2
2. 执行一次 run() , counter + 2 这种情况下，第二个 goroutine 开始执行时，没有拿到了第一个 goroutine 的执行结果，也就是 value := counter 这一步时，counter 还是零值，这时候 value = 0

当然由于种种不确定性，所有肯定不止这两种情况，但是这个不是本文讨论的重点，具体的原因可以结合上一篇文章[Week03: Go 并发编程(二)](https://lailin.xyz/post/go-training-week3-go-memory-model.html) Go 内存模型 进行思考

**data race 执行**

可以发现，写出这种代码时上线后如果出现 bug 会非常难定位，因为你不知道到底是哪里出现了问题，所以我们就要在测试阶段就结合 data race 工具提前发现问题。
我们执行以下命令
```
go run -race ./main.go
```
会发现结果会所有的都输出， data race  的报错信息，我们已经看不到了，因为终端的打印的太长了，可以发现的是，最后打印出发现了一处 data race 并且推出码为  66
```
Final Counter: 399956
Final Counter: 399960
Found 1 data race(s)
exit status 66
```
**data race 配置**

问题来了，我们有没有什么办法可以立即知道 data race 的报错呢？
答案就在官方的文档当中，我们可以通过设置 GORACE  环境变量，来控制 data race 的行为， 格式如下:
```
GORACE="option1=val1 option2=val2"
```
可选配置:
<img width="736" alt="image" src="https://user-images.githubusercontent.com/6757408/199284159-54c903d5-635a-433f-b0af-9ad879f33edd.png">

有了这个背景知识后就很简单了，在我们这个场景我们可以控制发现数据竞争后直接退出
```
GORACE="halt_on_error=1 strip_path_prefix=/home/ll/project/Go-000/Week03/blog/03_sync/01_data_race" go run -race ./main.go
```
重新执行后我们的结果
```
==================
WARNING: DATA RACE
Read at 0x00000064a9c0 by goroutine 7:
  main.routine()
      /main.go:29 +0x47

Previous write at 0x00000064a9c0 by goroutine 8:
  main.routine()
      /main.go:31 +0x64

Goroutine 7 (running) created at:
  main.run()
      /main.go:21 +0x75
  main.main()
      /main.go:14 +0x38

Goroutine 8 (finished) created at:
  main.run()
      /main.go:21 +0x75
  main.main()
      /main.go:14 +0x38
==================
exit status 66
```
这个结果非常清晰的告诉了我们在 29 行这个地方我们有一个 goroutine 在读取数据，但是呢，在 31 行这个地方又有一个 goroutine 在写入，所以产生了数据竞争。
然后下面分别说明这两个 goroutine 是什么时候创建的，已经当前是否在运行当中。

#### 典型案例
接来下我们再来看一些典型案例，这些案例都来自 go 官方的文档[Data Race Detector](https://go.dev/doc/articles/race_detector)，这些也是初学者很容易犯的错误

**案例二 在循环中启动 goroutine 引用临时变量**
```
func main() {
	var wg sync.WaitGroup
	wg.Add(5)
	for i := 0; i < 5; i++ {
		go func() {
			fmt.Println(i) // Not the 'i' you are looking for.
			wg.Done()
		}()
	}
	wg.Wait()
}
```
如果你去找一些 go 的面试题，肯定能找到类似的例子，然后会问你这里会输出什么？
常见的答案就是会输出 5 个 5，因为在 for 循环的 i++ 会执行的快一些，所以在最后打印的结果都是 5
这个答案不能说不对，因为真的执行的话大概率也是这个结果，但是不全
因为这里本质上是有数据竞争，在新启动的 goroutine 当中读取 i 的值，在 main 中写入，导致出现了 data race，这个结果应该是不可预知的，因为我们不能假定 goroutine 中 print 就一定
比外面的 i++ 慢，习惯性的做这种假设在并发编程中是很有可能会出问题的
```
func main() {
	var wg sync.WaitGroup
	wg.Add(5)
	for i := 0; i < 5; i++ {
		go func(i int) {
			fmt.Println(i)
			wg.Done()
		}(i)
	}
	wg.Wait()
}
```
这个要修改也很简单，只需要将 i 作为参数传入即可，这样每个 goroutine 拿到的都是拷贝后的数据

**案例三 一不小心就把变量共享了**
```
package main

import "os"

func main() {
	ParallelWrite([]byte("xxx"))
}

// ParallelWrite writes data to file1 and file2, returns the errors.
func ParallelWrite(data []byte) chan error {
	res := make(chan error, 2)
	f1, err := os.Create("/tmp/file1")
	if err != nil {
		res <- err
	} else {
		go func() {
			// This err is shared with the main goroutine,
			// so the write races with the write below.
			_, err = f1.Write(data)
			res <- err
			f1.Close()
		}()
	}
	f2, err := os.Create("/tmp/file2") // The second conflicting write to err.
	if err != nil {
		res <- err
	} else {
		go func() {
			_, err = f2.Write(data)
			res <- err
			f2.Close()
		}()
	}
	return res
}
```
我们使用 go run -race main.go  执行，可以发现这里报错的地方是，19 行和 24 行，有 data race，这里主要是因为共享了 err 这个变量
```
==================
WARNING: DATA RACE
Write at 0x00c0000a01a0 by goroutine 7:
  main.ParallelWrite.func1()
      /home/ll/project/Go-000/Week03/blog/03_data_race/03/main.go:19 +0x94

Previous write at 0x00c0000a01a0 by main goroutine:
  main.ParallelWrite()
      /home/ll/project/Go-000/Week03/blog/03_data_race/03/main.go:24 +0x1dd
  main.main()
      /home/ll/project/Go-000/Week03/blog/03_data_race/03/main.go:6 +0x84

Goroutine 7 (running) created at:
  main.ParallelWrite()
      /home/ll/project/Go-000/Week03/blog/03_data_race/03/main.go:16 +0x336
  main.main()
      /home/ll/project/Go-000/Week03/blog/03_data_race/03/main.go:6 +0x84
==================
Found 1 data race(s)
exit status 66
```
修改的话只需要在两个 goroutine 中使用新的临时变量就行了
```
...
_, err := f1.Write(data)
...
_, err := f2.Write(data)
...
```
细心的同学可能会有这个疑问，在 24 行不也是重新赋值了么，为什么在这里会和 19 行产生 data race 呢？
这是由于 go 的语法规则导致的，我们在初始化变量的时候如果在同一个作用域下，如下方代码，这里使用的 err 其实是同一个变量，只是 f1 f2 不同，
具体可以看[effective go 当中 Redeclaration and reassignment](https://go.dev/doc/effective_go#redeclaration)的内容
```
f1, err := os.Create("a")
f2, err := os.Create("b")
```
**案例四 不受保护的全局变量**
```
var service = map[string]string{}

// RegisterService RegisterService
func RegisterService(name, addr string) {
	service[name] = addr
}

// LookupService LookupService
func LookupService(name string) string {
	return service[name]
}
```
这个也是很容易犯的一个错，在之前写 Go 设计模式这个系列文章的时候，应该有提到过我们要写出可测性比较高的代码就要少用或者是尽量避免用全局变量，使用 map 作为全局变量比较常见的一种情况
就是配置信息。关于全局变量的话一般的做法就是加锁，就本文这个问题也可以使用 sync.Map 这个下一篇文章会讲，这里篇幅有限就不多讲了
```
var (
	service   map[string]string
	serviceMu sync.Mutex
)

func RegisterService(name, addr string) {
	serviceMu.Lock()
	defer serviceMu.Unlock()
	service[name] = addr
}

func LookupService(name string) string {
	serviceMu.Lock()
	defer serviceMu.Unlock()
	return service[name]
}
```
**案例五 未受保护的成员变量**
```
type Watchdog struct{ last int64 }

func (w *Watchdog) KeepAlive() {
	w.last = time.Now().UnixNano() // First conflicting access.
}

func (w *Watchdog) Start() {
	go func() {
		for {
			time.Sleep(time.Second)
			// Second conflicting access.
			if w.last < time.Now().Add(-10*time.Second).UnixNano() {
				fmt.Println("No keepalives for 10 seconds. Dying.")
				os.Exit(1)
			}
		}
	}()
}
```
同样成员变量也会有这个问题，这里可以用 atomic  包来解决，同样这个我们下篇文章会细讲
```
type Watchdog struct{ last int64 }

func (w *Watchdog) KeepAlive() {
	atomic.StoreInt64(&w.last, time.Now().UnixNano())
}

func (w *Watchdog) Start() {
	go func() {
		for {
			time.Sleep(time.Second)
			if atomic.LoadInt64(&w.last) < time.Now().Add(-10*time.Second).UnixNano() {
				fmt.Println("No keepalives for 10 seconds. Dying.")
				os.Exit(1)
			}
		}
	}()
}
```
**案例六 一个有趣的例子**
dava 在博客中提到过一个很有趣的例子的[Ice cream makers and data races](https://dave.cheney.net/2014/06/27/ice-cream-makers-and-data-races)
```
package main

import "fmt"

type IceCreamMaker interface {
	// Great a customer.
	Hello()
}

type Ben struct {
	name string
}

func (b *Ben) Hello() {
	fmt.Printf("Ben says, \"Hello my name is %s\"\n", b.name)
}

type Jerry struct {
	name string
}

func (j *Jerry) Hello() {
	fmt.Printf("Jerry says, \"Hello my name is %s\"\n", j.name)
}

func main() {
	var ben = &Ben{name: "Ben"}
	var jerry = &Jerry{"Jerry"}
	var maker IceCreamMaker = ben

	var loop0, loop1 func()

	loop0 = func() {
		maker = ben
		go loop1()
	}

	loop1 = func() {
		maker = jerry
		go loop0()
	}

	go loop0()

	for {
		maker.Hello()
	}
}
```
这个例子有趣的点在于，最后输出的结果会有这种例子
```
Ben says, "Hello my name is Jerry"
Ben says, "Hello my name is Jerry"
```
这是因为我们在 maker = jerry  这种赋值操作的时候并不是原子的，在上一篇文章中我们讲到过，只有对 single machine word 进行赋值的时候才是原子的，虽然这个看上去只有一行，
但是 interface 在 go 中其实是一个结构体，它包含了 type 和 data 两个部分，所以它的复制也不是原子的，会出现问题
```
type interface struct {
       Type uintptr     // points to the type of the interface implementation
       Data uintptr     // holds the data for the interface's receiver
}
```
这个案例有趣的点还在于，这个案例的两个结构体的内存布局一模一样所以出现错误也不会 panic 退出，如果在里面再加入一个 string 的字段，去读取就会导致 panic，但是这也恰恰说明这个案例很
可怕，这种错误在线上实在太难发现了，而且很有可能会很致命。
这个案例还有一个衍生案例，大家有兴趣可以点开查看一下，并不是说要看起来一样才不会 panic https://www.ardanlabs.com/blog/2014/06/ice-cream-makers-and-data-races-part-ii.html

### 总结
回顾一下，这篇文章通过一个案例讲解了 data race 的使用方法:
```
go build -race main.go
go test -race ./...
```
然后讲述了 data race 如何通过 GORACE 环境变量进行配置
最后讲解了几个典型案例，看完这篇相信你对 data race 已经有了一个基本的了解，希望可以在接下来的工作学习当中对你有有所启发
最后在重申一下关键点：
* 善用 data race 这个工具帮助我们提前发现并发错误
* 不要对未定义的行为做任何假设，虽然有时候我们写的只是一行代码，但是 go 编译器可能后面坐了很多事情，并不是说一行写完就一定是原子的
* 即使是原子的出现了 data race 也不能保证安全，因为我们还有可见性的问题，上篇我们讲到了现代的 cpu 基本上都会有一些缓存的操作。
* 所有出现了 data race 的地方都需要进行处理

### 参考文献
* https://dave.cheney.net/2014/06/27/ice-cream-makers-and-data-races
* https://www.ardanlabs.com/blog/2014/06/ice-cream-makers-and-data-races-part-ii.html
* http://blog.golang.org/race-detector
* https://golang.org/doc/articles/race_detector.html
* https://dave.cheney.net/2018/01/06/if-aligned-memory-writes-are-atomic-why-do-we-need-the-sync-atomic-package 除了考虑原子性之外，还要考虑可见性，并不是说赋值原子了，并发操作就没有问题了
* https://golang.org/doc/effective_go.html#redeclaration

转自：https://lailin.xyz/post/go-training-week3-data-race.html


















