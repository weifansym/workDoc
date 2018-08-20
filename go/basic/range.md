## range常用的操作
### 常用操作
我们经常看到range表达式用在**array, slice, string 和 map**上面，如下代码：
```
m := make(map[string]float64)
m["oneone"] = 1.1
m["twotwo"] = 2.2
for key, value := range m {
    fmt.Printf("[%s]: %.1f\n", key, value)
}
a := [3]int{3, 2, 1}
for idx, value := range a {
    fmt.Printf("[%d]: %d\n", idx, value)
}
s := []int{30, 20, 10}
for idx, value := range s {
    fmt.Printf("[%d]: %d\n", idx, value)
}
name := "Michał"
for idx, code := range name {
    fmt.Printf("[%d]: %q\n", idx, code)
}

```
输出如下：
```
[oneone]: 1.1
[twotwo]: 2.2
[0]: 3
[1]: 2
[2]: 1
[0]: 30
[1]: 20
[2]: 10
[0]: 'M'
[1]: 'i'
[2]: 'c'
[3]: 'h'
[4]: 'a'
[5]: 'ł'
```
### 遍历channel
除了遍历上面的类型之外，我们还用他来遍历发送到channel中的值。为了能够打断channel的遍历我们必须明确的关闭channel。否则range将会一直阻塞，
行为和nil（空）channel一样，让我们看一段代码：
```
package main
import "fmt"
func FibonacciProducer(ch chan int, count int) {
    n2, n1 := 0, 1
    for count >= 0 {
        ch <- n2
        count--
        n2, n1 = n1, n2+n1
    }
    close(ch)
}
func main() {
    ch := make(chan int)
    go FibonacciProducer(ch, 10)
    idx := 0
    for num := range ch {
        fmt.Printf("F(%d): \t%d\n", idx, num)
        idx++
    }
}
```
输出如下：
```
F(0): 0
F(1): 1
F(2): 1
F(3): 2
F(4): 3
F(5): 5
F(6): 8
F(7): 13
F(8): 21
F(9): 34
F(10): 55
```
使用range遍历channel的时候索引变量是不允许的，下面代码将会抛出错误：
```
ch := make(chan int)
go FibonacciProducer(ch, 10)
for idx, num := range ch {
    fmt.Printf("idx: %d F(%d): \t%d\n", idx, num)
}
```
在编译的时候将会抛出如下错误：
```
too many variables in range error 
```


