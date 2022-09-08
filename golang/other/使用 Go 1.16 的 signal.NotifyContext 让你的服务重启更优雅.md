## 使用 Go 1.16 的 signal.NotifyContext 让你的服务重启更优雅
在 Go 1.16 的更新中，signal包增加了一个函数[NotifyContext](https://pkg.go.dev/os/signal#NotifyContext)，
这让我们优雅的重启服务（Graceful Restart）可以写的更加优雅。

一个服务想要优雅的重启主要包含两个方面：
1. 退出的旧服务需要 Graceful Shutdown，不强制杀进程，不泄漏系统资源。
2. 在一个集群内轮流重启服务实例，保证服务不中断。

第二个问题跟部署方式相关，改天专门写一篇讨论，今天我们主要谈怎么样优雅的退出。

首先在代码里，用了外部资源，一定要使用defer去调用Close()方法关闭。
然后我们就要拦截系统的中断信号，保证程序收到中断信号之后，主动有序退出，这样所有的 defer 才会被执行。

在以前，大概是这么写：
```
func everLoop(ctx context.Context) {
LOOP:
    for {
        select {
        case <-ctx.Done():
            // 收到信号退出无限循环
            break LOOP
        default:
            // 用一个 sleep 模拟业务逻辑
            time.Sleep(time.Second * 10)
        }
    }
}

func main() {
    // 建立一个可以手动取消的 Context
    ctx, cancel := context.WithCancel(context.Background())

    // 监控系统信号，这里只监控了 SIGINT（Ctrl+c），SIGTERM
    // 在 systemd 和 docker 中，都是先发 SIGTERM，过一段时间没退出再发 SIGKILL
    // 所以这里没捕获 SIGKILL
    sig := make(chan os.Signal, 1)
    signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
    go func() {
        <-sig
        cancel()
    }()

    // 开始无限循环，收到信号就会退出
    everLoop(ctx)
    fmt.Println("graceful shuwdown")
}

```
现在有了新的函数，这一段变得更简单了：
```
func main() {
    // 监控系统信号和创建 Context 现在一步搞定
    ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
    // 在收到信号的时候，会自动触发 ctx 的 Done ，这个 stop 是不再捕获注册的信号的意思，算是一种释放资源。
    defer stop()

    // 开始无限循环，收到信号就会退出
    everLoop(ctx)
    fmt.Println("graceful shuwdown")
}
```
感谢 Golang ，当年用别的语言需要写一大堆代码的功能，现在几行就可以轻松实现了。
让它成为你服务程序的标配吧。

最后，我是写最新的独立项目LetServerRun的时候，发现这种最新的写法的。
LetServerRun 可以让你把微信公众号当作随身的 Terminal 控制你的服务端。
在它的 Agent 的 main 函数中就有上述用法的示例，
欢迎参考。

