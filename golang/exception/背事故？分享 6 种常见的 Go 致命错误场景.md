## 背事故？分享 6 种常见的 Go 致命错误场景
大家好，我是煎鱼。
有一次事故现场，在紧急恢复后，他正在排查代码，查了好一会。我回头一看，这错误提醒很明显就是致命错误，较好定位。
但此时，他竟然在查 panic-recover 是不是哪里漏了，我表示大受震惊...
今天就由煎鱼给大家分享一下错误类型有哪几种，又在什么场景下会触发。
### 错误类型
#### error

第一种是 Go 中最标准的 error 错误，其真身是一个 interface{}。

如下：
```
func throw(s string) {
	systemstack(func() {
		print("fatal error: ", s, "\n")
	})
	gp := getg()
	if gp.m.throwing == 0 {
		gp.m.throwing = 1
	}
	fatalthrow()
	*(*int)(nil) = 0 // not reached
}
```
根据上述程序，会获取当前 G 的实例，并设置其 M 的 throwing 状态为 1。

状态设置好后，会调用**fatalthrow**方法进行真正的 crash 相关操作
```
func fatalthrow() {
	pc := getcallerpc()
	sp := getcallersp()
	gp := getg()
	
	systemstack(func() {
		startpanic_m()
		if dopanic_m(gp, pc, sp) {
			crash()
		}

		exit(2)
	})

	*(*int)(nil) = 0 // not reached
}
```
主体逻辑是发送 _SIGABRT 信号量，最后调用 exit 方法退出，所以你会发现这是拦也拦不住的 “致命” 错误。

### 致命场景
为此，作为一名 “成熟” 的 Go 工程师，除了保障自己程序的健壮性外，我也在网上收集了一些致命的错误场景，分享给大家。

一起学习和规避这些致命场景，年底争取拿个 A，不要背上 P0 事故。

#### 并发读写 map
```
func foo() {
	m := map[string]int{}
	go func() {
		for {
			m["煎鱼1"] = 1
		}
	}()
	for {
		_ = m["煎鱼2"]
	}
}
```
输出结果：
```
fatal error: concurrent map read and map write

goroutine 1 [running]:
runtime.throw(0x1078103, 0x21)
...
```
#### 堆栈内存耗尽
```
func foo() {
	var f func(a [1000]int64)
	f = func(a [1000]int64) {
		f(a)
	}
	f([1000]int64{})
}
```
输出结果：
```
runtime: goroutine stack exceeds 1000000000-byte limit
runtime: sp=0xc0200e1bf0 stack=[0xc0200e0000, 0xc0400e0000]
fatal error: stack overflow

runtime stack:
runtime.throw(0x1074ba3, 0xe)
        /usr/local/Cellar/go/1.16.6/libexec/src/runtime/panic.go:1117 +0x72
runtime.newstack()
...
```
#### 将 nil 函数作为 goroutine 启动
```
func foo() {
	var f func()
	go f()
}
```
输出结果：
```
fatal error: go of nil func value

goroutine 1 [running]:
main.foo()
...
```
#### goroutines 死锁
```
func foo() {
	select {}
}
```
输出结果：
```
fatal error: all goroutines are asleep - deadlock!

goroutine 1 [select (no cases)]:
main.foo()
...
```
#### 线程限制耗尽
如果你的 goroutines 被 IO 操作阻塞了，新的线程可能会被启动来执行你的其他 goroutines。

Go 的最大的线程数是有默认限制的，如果达到了这个限制，你的应用程序就会崩溃。

会出现如下输出结果：
```
fatal error: thread exhaustion
...
```
可以通过调用**runtime.SetMaxThreads**方法增大线程数，不过也需要考量是否程序有问题。

#### 超出可用内存
如果你执行的操作，例如：下载大文件等。导致应用程序占用内存过大，程序上涨，导致 OOM。

会出现如下输出结果：
```
fatal error: runtime: out of memory
...
```
建议处理掉一些程序，或者换新电脑了。
#### 总结
在今天这篇文章中，我们介绍了 Go 语言的三种错误类型。其中针对大家最少见，但一碰到就很容易翻车的致命错误 fatal error 进行了介绍，给出了一些经典案例。
希望大家后续能够规避，你有没有遇到过其中的场景？

转自：https://juejin.cn/post/7045122592939704333


