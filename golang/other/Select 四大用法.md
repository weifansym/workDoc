## Select 四大用法
本篇教学要带大家认识Go 语言的Select用法，相信大家对于switch 并不陌生，但是 select 跟 switch 有个共同特性就是都过case 的方式来处理，但是select 跟switch 处理的事情完全不同，
也完全不相容。来看看switch 有什么特性: 各种类型及型别操作，接口 interface{} 型别判断variable.(type)，重点是会依照case顺序依序执行。底下看个例子:
```
package main

var (
    i interface{}
)

func convert(i interface{}) {
    switch t := i.(type) {
    case int:
        println("i is interger", t)
    case string:
        println("i is string", t)
    case float64:
        println("i is float64", t)
    default:
        println("type not found")
    }
}

func main() {
    i = 100
    convert(i)
    i = float64(45.55)
    convert(i)
    i = "foo"
    convert(i)
    convert(float32(10.0))
}
```
结果如下：
```
i is interger 100
i is float64 +4.555000e+001
i is string foo
type not found
```
而 select 的特性就不同了，只能接channel 否则会出错，default会直接执行，所以没有 default 的select 就会遇到blocking，假设没有送value 进去Channel 就会造成panic，
底下拿几个实际例子来解说。

### Random Select
同一个channel 在select 会随机选取，底下看个例子:
```
package main

import "fmt"

func main() {
    ch := make(chan int, 1)

    ch <- 1
    select {
    case <-ch:
        fmt.Println("random 01")
    case <-ch:
        fmt.Println("random 02")
    }
}
```
执行后会发现有时候拿到 random 01 有时候拿到random 02，这就是select 的特性之一，case 是随机选取，所以当select 有两个channel 以上时，如果同时对全部channel 送资料，则会随机
选取到不同的Channel。而上面有提到另一个特性『假设没有送value 进去Channel 就会造成panic』，拿上面例子来改:
```
func main() {
    ch := make(chan int, 1)

    select {
    case <-ch:
        fmt.Println("random 01")
    case <-ch:
        fmt.Println("random 02")
    }
}
```
执行后会发现变成deadlock，造成main 主程式爆炸，这时候可以直接用 default 方式解决此问题:
```
func main() {
    ch := make(chan int, 1)

    select {
    case <-ch:
        fmt.Println("random 01")
    case <-ch:
        fmt.Println("random 02")
    default:
        fmt.Println("exit")
    }
}
```
主程式main 就不会因为读不到channel value 造成整个程式deadlock。
### Timeout 超时机制
用select 读取channle 时，一定会实作超过一定时间后就做其他事情，而不是一直blocking 在select 内。底下是简单的例子:
```
package main

import (
    "fmt"
    "time"
)

func main() {
    timeout := make(chan bool, 1)
    go func() {
        time.Sleep(2 * time.Second)
        timeout <- true
    }()
    ch := make(chan int)
    select {
    case <-ch:
    case <-timeout:
        fmt.Println("timeout 01")
    }
}
```
建立timeout channel，让其他地方可以透过trigger timeout channel 达到让select 执行结束，也或者有另一个写法是透握 time.After 机制
```
select {
    case <-ch:
    case <-timeout:
        fmt.Println("timeout 01")
    case <-time.After(time.Second * 1):
        fmt.Println("timeout 02")
    }
```
可以注意 time.After 是回传chan time.Time，所以执行select 超过一秒时，就会输出timeout 02。

### 检查channel 是否已满
直接来看例子比较快:
```
package main

import (
    "fmt"
)

func main() {
    ch := make(chan int, 1)
    ch <- 1
    select {
    case ch <- 2:
        fmt.Println("channel value is", <-ch)
        fmt.Println("channel value is", <-ch)
    default:
        fmt.Println("channel blocking")
    }
}
```
先宣告buffer size 为1 的channel，先丢值把channel 填满。这时候可以透过 select + default 方式来确保channel 是否已满，上面例子会输出channel blocking，我们再把程式改成底下
```
func main() {
    ch := make(chan int, 2)
    ch <- 1
    select {
    case ch <- 2:
        fmt.Println("channel value is", <-ch)
        fmt.Println("channel value is", <-ch)
    default:
        fmt.Println("channel blocking")
    }
}
```
把buffer size 改为2 后，就可以继续在塞value 进去channel 了，这边的buffer channel 观念可以看之前的文章[用五分钟了解什么是unbuffered vs buffered channel](https://blog.wu-boy.com/2019/04/understand-unbuffered-vs-buffered-channel-in-five-minutes/)

### select for loop 用法
如果你有多个channel 需要读取，而读取是不间断的，就必须使用for + select 机制来实现，更详细的实作可以参考[15 分钟学习Go 语言如何处理多个Channel 通道](https://blog.wu-boy.com/2019/05/handle-multiple-channel-in-15-minutes/)
```
package main

import (
    "fmt"
    "time"
)

func main() {
    i := 0
    ch := make(chan string, 0)
    defer func() {
        close(ch)
    }()

    go func() {
    LOOP:
        for {
            time.Sleep(1 * time.Second)
            fmt.Println(time.Now().Unix())
            i++

            select {
            case m := <-ch:
                println(m)
                break LOOP
            default:
            }
        }
    }()

    time.Sleep(time.Second * 4)
    ch <- "stop"
}
```
上面例子可以发现执行后如下:
```
1574474619
1574474620
1574474621
1574474622
```
其实把 default 拿掉也可以达到目的
```
select {
case m := <-ch:
    println(m)
    break LOOP
```
当没有值送进来时，就会一直停在select 区段，所以其实没有 default 也是可以正常运作的，而要结束for 或select 都需要透过break 来结束，但是要在select 区间直接结束掉for 回圈，
只能使用 break variable 来结束，这边是大家需要注意的地方。

转自：
* https://blog.wu-boy.com/2019/11/four-tips-with-select-in-golang/
* 


