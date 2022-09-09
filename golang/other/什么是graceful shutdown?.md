## 什么是graceful shutdown?
我们该如何升级Web 服务，你会说很简单啊，只要关闭服务，上程式码，再开启服务即可，可是很多时候开发者可能没有想到现在服务上面是否有正在处理的资料，像是购物车交易？也或者是说背景有正在
处理重要的事情，如果强制关闭服务，就会造成下次启动时会有一些资料上的差异，那该如何优雅地关闭服务，这就是本篇的重点了。底下先透过简单的gin http 服务范例介绍简单的web 服务
### 基本HTTPD 服务
```
package main

import (
    "log"
    "net/http"
    "time"

    "github.com/gin-gonic/gin"
)

func main() {
    router := gin.Default()
    router.GET("/", func(c *gin.Context) {
        time.Sleep(5 * time.Second)
        c.String(http.StatusOK, "Welcome Gin Server")
    })

    srv := &http.Server{
        Addr:    ":8080",
        Handler: router,
    }

    // service connections
    if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
        log.Fatalf("listen: %s\n", err)
    }

    log.Println("Server exiting")
}
```
上述程式码在我们写基本的web 服务都不会考虑到[graceful shutdown](https://go.dev/doc/go1.8#http_shutdown)，如果有重要的Job 在上面跑，我强烈建议一定要加上Go在1.8 版推出的graceful shutdown 函式，
上述程式码假设透过底下指令执行:
```
curl -v http://localhost:8080
```
接着把server 关闭，就会强制关闭client 连线，并且喷错。底下会用graceful shutdown 来解决此问题。
### 使用graceful shutdown
Go 1.8 推出graceful shutdown，让开发者可以针对不同的情境在升级过程中做保护，整个流程大致上会如下:
1. 关闭服务连接埠
2. 等待并且关闭所有连线

可以看到步骤1. 会先关闭连接埠，确保没有新的使用者连上服务，第二步骤就是确保处理完剩下的http 连线才会正常关闭，来看看底下范例
```
// +build go1.8

package main

import (
    "context"
    "log"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"

    "github.com/gin-gonic/gin"
)

func main() {
    router := gin.Default()
    router.GET("/", func(c *gin.Context) {
        time.Sleep(5 * time.Second)
        c.String(http.StatusOK, "Welcome Gin Server")
    })

    srv := &http.Server{
        Addr:    ":8080",
        Handler: router,
    }

    go func() {
        // service connections
        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatalf("listen: %s\n", err)
        }
    }()

    // Wait for interrupt signal to gracefully shutdown the server with
    // a timeout of 5 seconds.
    quit := make(chan os.Signal, 1)
    // kill (no param) default send syscall.SIGTERM
    // kill -2 is syscall.SIGINT
    // kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit
    log.Println("Shutdown Server ...")

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    if err := srv.Shutdown(ctx); err != nil {
        log.Fatal("Server Shutdown: ", err)
    }

    log.Println("Server exiting")
}
```
首先可以看到将 srv.ListenAndServe 直接丢到背景执行，这样才不会阻断后续的流程，接着宣告一个 os.Signal 讯号的Channel，并且接受系统SIGINT 及SIGTERM，
也就是只要透过kill 或者是 docker rm 就会收到讯号关闭 quit 通道
```
<-quit
```
由上面可知，整个main func 会被block 在这地方，假设按下ctrl + c 就会被系统讯号(SIGINT 及SIGTERM) 通知关闭quit 通道，通道被关闭后，就会继续往下执行
```
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()
if err := srv.Shutdown(ctx); err != nil {
    log.Fatal("Server Shutdown: ", err)
}
```
最后可以看到 srv.Shutdown 就是用来处理『1. 关闭连接埠』及『2. 等待所有连线处理结束』，可以看到传了一个context 进Shutdown 函式，目的就是让程式最多等待5 秒时间，
如果超过5 秒就强制关闭所有连线，所以您需要根据server 处理的资料时间来决定等待时间，设定太短就会造成强制关闭，建议依照情境来设定。至于服务shutdown 后可以处理哪些事情就看开发者决定。
1. 关闭Database 连线
2. 等到背景worker 处理

可以搭配上一篇提到的[『graceful shutdown with multiple workers』](https://blog.wu-boy.com/2020/02/graceful-shutdown-with-multiple-workers/)

转自：https://blog.wu-boy.com/2020/02/what-is-graceful-shutdown-in-golang/








