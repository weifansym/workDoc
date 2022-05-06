
这里整理一下go开发当中用到了并发协程多任务，同时收集返回多任务结果，go 协程没有直接返回，只能通过chan返回收集，其中用到几个特性

缓存管道是当满的时候是阻塞的，这个特性可以用到并发控制

需要用到&sync.WaitGroup{} 也就是说并发请求中的执行时间跟最长的有关，需要所有的计数器都消耗完了然后结束
> go语言在执行goroutine的时候、是没有返回值的、这时候我们要用到go语言中特色的channel来获取返回值。
通过channel拿到返回值有两种处理方式，一种形式是具有go风格特色的，即发送给一个for channel 或 select channel 的独立goroutine中，由该独立的goroutine来处理函数的返回值。
还有一种传统的做法，就是将所有goroutine的返回值都集中到当前函数，然后统一返回给调用函数。

### 第一种不用函数中统一返回，通过独立的goroutine接收返回值
```
package main

import (
    "fmt"
    "strconv"
    "sync"
    "time"
)

var responseChannel = make(chan string, 15)

func main() {
    fmt.Println(time.Now())
    go response()
    wg := &sync.WaitGroup{}
    //并发10
    limiter := make(chan bool, 10)
    for i := 0; i < 100; i++ {
        wg.Add(1)
        limiter <- true
        go httpGet(strconv.Itoa(i), limiter, wg)
    }
    wg.Wait()
    fmt.Println("all Done")
    fmt.Println(time.Now())
}

func httpGet(url string, limiter chan bool, wg *sync.WaitGroup) {

    defer wg.Done() //释放一个锁
    //do something
    time.Sleep(1 * time.Second)
    responseChannel <- fmt.Sprintf("Hello Go %s", url)
    <-limiter //释放一个hold
}
func response() {
    for rc := range responseChannel {
        fmt.Println("response:", rc)
    }
}
```
### 第二种：需要封装成一个函数的,即在当前函数中聚合返回
```
package main

import (
    "fmt"
    "strconv"
    "sync"
    "time"
)

func httpGet(url string, response chan string, limiter chan bool, wg *sync.WaitGroup) {
    //计数器-1
    defer wg.Done()
    //coding do business
    time.Sleep(1 * time.Second)
    //结果数据传入管道
    response <- fmt.Sprintf("http get:%s", url)
    //释放一个并发
    <-limiter
}

func collect(urls []string) []string {
    var result []string
    //执行的 这里要注意  需要指针类型传入  否则会异常
    wg := &sync.WaitGroup{}
    //并发控制
    limiter := make(chan bool, 10)
    defer close(limiter)

    response := make(chan string, 20)
    wgResponse := &sync.WaitGroup{}
    //处理结果 接收结果
    go func() {
        wgResponse.Add(1)
        for rc := range response {
            result = append(result, rc)
        }
        wgResponse.Done()
    }()
    //开启协程处理请求
    for _, url := range urls {
        //计数器
        wg.Add(1)
        //并发控制 10
        limiter <- true
        go httpGet(url, response, limiter, wg)
    }
    //发送任务
    wg.Wait()
    close(response) //关闭 并不影响接收遍历
    //处理接收结果
    wgResponse.Wait()
    return result

}

func main() {
    var urls []string
    for i := 0; i < 100; i++ {
        urls = append(urls, strconv.Itoa(i))
    }
    fmt.Println(time.Now())
    result := collect(urls)
    fmt.Println(time.Now())
    fmt.Println(result)
}
```

参考：https://yar999.gitbook.io/gopl-zh/ch1/ch1-06
