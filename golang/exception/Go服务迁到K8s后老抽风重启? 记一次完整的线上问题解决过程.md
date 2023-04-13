## Go服务迁到K8s后老抽风重启? 记一次完整的线上问题解决过程
### 前言
之前把Go服务都迁到Kubernetes上后有些服务的某个 Pod总是时不时的重启一下，通过查业务日志根本查不到原因，我分析了一下肯定是哪里代码不严谨造成引用空指针导致Go发送运行时panic才会挂掉的，
但是容器重启后之前输出到stderr的panic是会被清空的，所以才有了这篇文章里后面的分析和方案解决。

### 解决思路分析
在Go编写的应用程序里无论是在主协程（main goroutine）还是其他子协程里，一旦出了运行时panic错误后，整个程序都会宕掉。一般的部署Go项目的时候都会使用supervisor监控应用程序进程，
一旦应用程序发生panic停掉后supervisor会把进程再启动起来。

那么在把项目部署到Kubernetes集群后，因为每个节点上的kubelet会对主进程崩溃的容器进行重启，所以就再引入supervisor就有些功能重叠。但是Go的panic信息是直接写到标准错误的，容器重启后
之前的panic错误就没有了，没法排查导致容器崩溃的原因。所以排查容器重启的关键点就变成了：怎么把**panic从stderr**重定向到文件，这样就能通过容器的volume持久化日志文件的目录方式保留程序崩溃
时的信息。

那么以前在supervisor里可以直接通过配置**stderr_logfile**把程序运行时的标准错误设置成一个文件：
```
[program: go-xxx...]
directory=/home/go/src...
environment=...
command=/home/go/src.../bin/app
stderr_logfile=/home/xxx/log/..../app_err.log
```
现在换成了Kubernetes，不再使用supervisor后就只能是想办法在程序里实现了。针对在Go里实现记录panic到日志文件你可能首先会考虑：**在recover里把导致panic的错误记录到文件里**，不过引用
的第三方包里也有可能panic，这个不现实。而且Go 也没有其他语言那样的Exception，未捕获的异常能由全局的**ExceptionHandler**捕获到的机制，实现不了用一个recover捕获所有的panic的功能。

最后就只有一个办法了，想办法把程序运行时的标准错误替换成日志文件，这样Go再panic的时候它还是往标准错误里写，只不过我们偷偷把标准错误的文件描述符换成了日志文件的描述符（在系统眼里stderr
也是个文件，Unix系统里一切皆文件）。

### 方案试错
按着这个思路我先用下面例子的试了一下：
```
package main

import (
    "fmt"
    "os"
)

const stdErrFile = "/tmp/go-app1-stderr.log"

func RewriteStderrFile() error {
    file, err := os.OpenFile(stdErrFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
    if err != nil {
      fmt.Println(err)
        return err
    }
    os.Stderr = file
    return nil
}


func testPanic() {
    panic("test panic")
}

func main() {
    RewriteStderrFile()
    testPanic()
}
```
这个例子，我们尝试使用 os.Stderr = file 来强制转换，但运行程序后，发现不起作用，**/tmp/go-app1-stderr.log**没有任何信息流入，panic信息照样输出到标准错误里。

### 最终方案
关于原因，搜索了一下，幸运的是 Rob Pike有专门对类似问题的解答，是这样说的：
![image](https://user-images.githubusercontent.com/6757408/231690502-a78b21f9-d680-4c87-b825-3c635d34ccad.png)

把高层包创建的变量直接赋值到底层的runtime是不行的，我们用**syscall.Dup2**实现替换描述符再试一次，并且增加一个全局变量对日志文件描述符的引用，避免常驻线程运行中文件描述符被GC回收掉：
```
var stdErrFileHandler *os.File

func RewriteStderrFile() error {
   if runtime.GOOS == "windows" {
      return nil
   }
  
    file, err := os.OpenFile(stdErrFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
    if err != nil {
      fmt.Println(err)
        return err
    }
    stdErrFileHandler = file //把文件句柄保存到全局变量，避免被GC回收
    
    if err = syscall.Dup2(int(file.Fd()), int(os.Stderr.Fd())); err != nil {
        fmt.Println(err)
        return err
    }
    // 内存回收前关闭文件描述符
   runtime.SetFinalizer(stdErrFileHandler, func(fd *os.File) {
      fd.Close()
   })
   
    return nil
}
```
因为Windows系统不支持的syscall.Dup2这个函数，所以我加了个判读，Windows环境下的Go运行时加载系统的一个dll文件也能实现这里的功能，不过我们服务器环境都是Linux的，所以我认为这部分要
兼容Windows是无用功，保证项目在Windows下能跑不受影响就行了。

再次运行程序后，打开日志文件/tmp/go-app1-stderr.log后就能看到刚才程序崩溃时的panic信息，以及导致panic时整个调用栈的信息：
```
➜  ~ cat /tmp/go-app1-stderr.log 
panic: test panic

goroutine 1 [running]:
main.testPanic(...)
        /Users/kev/Library/Application Support/JetBrains/GoLand2020.1/scratches/scratch_4.go:39
main.main()
        /Users/kev/Library/Application Support/JetBrains/GoLand2020.1/scratches/scratch_4.go:44 +0x3f
panic: test panic

goroutine 1 [running]:
main.testPanic(...)
        /Users/kev/Library/Application Support/JetBrains/GoLand2020.1/scratches/scratch_4.go:39
main.main()
        /Users/kev/Library/Application Support/JetBrains/GoLand2020.1/scratches/scratch_4.go:44 +0x3f
```
### 方案实施后的效果
目前这个方案已经在我们线上运行一个月了，已发现的Pod重启事件都能把程序崩溃时的调用栈准确记录到日志文件里，帮助我们定位了几个代码里的问题。其实问题都是空指针相关的问题，这些问题我在
之前的文章[《如何避免用动态语言的思维写Go代码》](https://cloud.tencent.com/developer/article/1700381)也提到过，项目一旦复杂起来谁写的代码也不能保证说不会发生空指针，不过我们事先做好检查很多都是能够避免的明显错误，对于特别细微条件下引发
的错误只能靠分析事故当时的日志来解决啦。

转自：https://cloud.tencent.com/developer/article/1700381


