## 用Golang做一个永久阻塞，你有哪些小技巧
Go 的运行时的当前设计，假定程序员自己负责检测何时终止一个 goroutine 以及何时终止该程序。可以通过调用 os.Exit 或从 main() 函数的返回来以正常方式终止程序。而有时候我们需要的是使程序阻塞在这一行。
#### 使用 sync.WaitGroup
一直等待直到 WaitGroup 等于 0
```
package main
import "sync"

func main() {
    var wg sync.WaitGroup
    wg.Add(1)
    wg.Wait()
}
```
#### 空 select
select{} 是一个没有任何 case 的 select，它会一直阻塞
```
package main

func main() {
    select{}
}
```
#### 死循环
虽然能阻塞，但会 100% 占用一个 cpu。不建议使用
```
package main

func main() {
    for {}
}
```

参考：https://learnku.com/articles/64799

