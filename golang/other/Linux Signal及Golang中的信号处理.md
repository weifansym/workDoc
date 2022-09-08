## Linux Signal及Golang中的信号处理
信号(Signal)是Linux, 类Unix和其它POSIX兼容的操作系统中用来进程间通讯的一种方式。一个信号就是一个异步的通知，发送给某个进程，或者同进程的某个线程，告诉它们某个事件发生了。
当信号发送到某个进程中时，操作系统会中断该进程的正常流程，并进入相应的信号处理函数执行操作，完成后再回到中断的地方继续执行。
如果目标进程先前注册了某个信号的处理程序(signal handler),则此处理程序会被调用，否则缺省的处理程序被调用。

### 发送信号
kill 系统调用(system call)可以用来发送一个特定的信号给进程。
kill 命令允许用户发送一个特定的信号给进程。
raise 库函数可以发送特定的信号给当前进程。

在Linux下运行man kill可以查看此命令的介绍和用法。
> The command kill sends the specified signal to the specified process or process group. If no signal is specified, the TERM signal is sent. 
> The TERM signal will kill processes which do not catch this signal. For other processes, it may be necessary to use the KILL (9) signal, since this 
> signal cannot be caught.Most modern shells have a builtin kill function, with a usage rather similar to that of the command described here. 
> The '-a' and '-p' options, and the possibility to specify pids by command name is a local extension. 
> If sig is 0, then no signal is sent, but error checking is still performed.


一些异常比如除以0或者 segmentation violation 相应的会产生SIGFPE和SIGSEGV信号，缺省情况下导致core dump和程序退出。
内核在某些情况下发送信号，比如在进程往一个已经关闭的管道写数据时会产生SIGPIPE信号。
在进程的终端敲入特定的组合键也会导致系统发送某个特定的信号给此进程：
* Ctrl-C 发送 INT signal (SIGINT)，通常导致进程结束
* Ctrl-Z 发送 TSTP signal (SIGTSTP); 通常导致进程挂起(suspend)
* Ctrl-\ 发送 QUIT signal (SIGQUIT); 通常导致进程结束 和 dump core.
* Ctrl-T (不是所有的UNIX都支持) 发送INFO signal (SIGINFO); 导致操作系统显示此运行命令的信息

kill -9 pid 会发送 SIGKILL信号给进程。

### 处理信号
Signal handler可以通过signal()系统调用进行设置。如果没有设置，缺省的handler会被调用，当然进程也可以设置忽略此信号。
有两种信号不能被拦截和处理: SIGKILL和SIGSTOP。

当接收到信号时，进程会根据信号的响应动作执行相应的操作，信号的响应动作有以下几种：
* 中止进程(Term)
* 忽略信号(Ign)
* 中止进程并保存内存信息(Core)
* 停止进程(Stop)
* 继续运行进程(Cont)

用户可以通过signal或sigaction函数修改信号的响应动作（也就是常说的“注册信号”）。另外，在多线程中，各线程的信号响应动作都是相同的，不能对某个线程设置独立的响应动作。
### 信号类型
个平台的信号定义或许有些不同。下面列出了POSIX中定义的信号。
Linux 使用34-64信号用作实时系统中。
命令man 7 signal提供了官方的信号介绍。

在POSIX.1-1990标准中定义的信号列表

<img width="639" alt="image" src="https://user-images.githubusercontent.com/6757408/189111242-a41d1060-9d19-4568-822b-46531971bfe8.png">

在SUSv2和POSIX.1-2001标准中的信号列表:

<img width="673" alt="截屏2022-09-08 下午7 30 17" src="https://user-images.githubusercontent.com/6757408/189111489-5ab61c3c-d245-4e22-bdc1-385eba141e0a.png">

Windows中没有SIGUSR1,可以用SIGBREAK或者SIGINT代替。

### Go中的Signal发送和处理
有时候我们想在Go程序中处理Signal信号，比如收到SIGTERM信号后优雅的关闭程序(参看下一节的应用)。
Go信号通知机制可以通过往一个channel中发送os.Signal实现。
首先我们创建一个os.Signal channel，然后使用signal.Notify注册要接收的信号。
```
package main
import "fmt"
import "os"
import "os/signal"
import "syscall"
func main() {
    // Go signal notification works by sending `os.Signal`
    // values on a channel. We'll create a channel to
    // receive these notifications (we'll also make one to
    // notify us when the program can exit).
    sigs := make(chan os.Signal, 1)
    done := make(chan bool, 1)
    // `signal.Notify` registers the given channel to
    // receive notifications of the specified signals.
    signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
    // This goroutine executes a blocking receive for
    // signals. When it gets one it'll print it out
    // and then notify the program that it can finish.
    go func() {
        sig := <-sigs
        fmt.Println()
        fmt.Println(sig)
        done <- true
    }()
    // The program will wait here until it gets the
    // expected signal (as indicated by the goroutine
    // above sending a value on `done`) and then exit.
    fmt.Println("awaiting signal")
    <-done
    fmt.Println("exiting")
}
```
go run main.go执行这个程序，敲入ctrl-C会发送SIGINT信号。 此程序接收到这个信号后会打印退出。

### Go网络服务器如果无缝重启
Go很适合编写服务器端的网络程序。DevOps经常会遇到的一个情况是升级系统或者重新加载配置文件，在这种情况下我们需要重启此网络程序，如果网络程序暂停的时间较长，则给客户的感觉很不好。
如何实现优雅地重启一个Go网络程序呢。主要要解决两个问题：
1. 进程重启不需要关闭监听的端口
2. 既有请求应当完全处理或者超时

@humblehack 在他的文章[Graceful Restart in Golang](https://grisha.org/blog/2014/06/03/graceful-restart-in-golang/)中提供了一种方式，而Florian von Bock根据此思路实现了一个框架[endless](https://github.com/fvbock/endless)。
此框架使用起来超级简单:
```
err := endless.ListenAndServe("localhost:4242", mux)
```
只需替换 http.ListenAndServe 和 http.ListenAndServeTLS。

它会监听这些信号： syscall.SIGHUP, syscall.SIGUSR1, syscall.SIGUSR2, syscall.SIGINT, syscall.SIGTERM, 和 syscall.SIGTSTP。

此文章提到的思路是：
1. 通过exec.Command fork一个新的进程，同时继承当前进程的打开的文件(输入输出，socket等)
```
file := netListener.File() // this returns a Dup()
path := "/path/to/executable"
args := []string{
    "-graceful"}
cmd := exec.Command(path, args...)
cmd.Stdout = os.Stdout
cmd.Stderr = os.Stderr
cmd.ExtraFiles = []*os.File{file}
err := cmd.Start()
if err != nil {
    log.Fatalf("gracefulRestart: Failed to launch, error: %v", err)
}
```
2. 子进程初始化
网络程序的启动代码
```
server := &http.Server{Addr: "0.0.0.0:8888"}
 var gracefulChild bool
 var l net.Listever
 var err error
 flag.BoolVar(&gracefulChild, "graceful", false, "listen on fd open 3 (internal use only)")
 if gracefulChild {
     log.Print("main: Listening to existing file descriptor 3.")
     f := os.NewFile(3, "")
     l, err = net.FileListener(f)
 } else {
     log.Print("main: Listening on a new file descriptor.")
     l, err = net.Listen("tcp", server.Addr)
 }
```
3. 父进程停止
```
if gracefulChild {
    parent := syscall.Getppid()
    log.Printf("main: Killing parent pid: %v", parent)
    syscall.Kill(parent, syscall.SIGTERM)
}
server.Serve(l)
```
同时他还提供的如何处理已经正在处理的请求。可以查看它的文章了解详细情况。

因此，处理特定的信号可以实现程序无缝的重启。
### 其它
graceful shutdown实现非常的简单，通过简单的信号处理就可以实现。本文介绍的是graceful restart,要求无缝重启，所以所用的技术相当的hack。

Facebook的工程师也提供了http和net的实现：[facebookgo](https://github.com/facebookarchive/grace)。

转自：[Linux Signal及Golang中的信号处理](https://colobu.com/2015/10/09/Linux-Signals/)

参考：https://tonybai.com/2012/09/21/signal-handling-in-go/


