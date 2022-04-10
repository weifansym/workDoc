## golang判断map中key是否存在的方法
```
import "fmt"

func main() {
    dict := map[string]int{"key1": 1, "key2": 2}
    value, ok := dict["key1"]
    if ok {
        fmt.Printf(value)
    } else {
        fmt.Println("key1 不存在")
    }
}
```
以上就是golang中判断map中key是否存在的方法

还有一种简化的写法是
```
import "fmt"

func main() {
    dict := map[string]int{"key1": 1, "key2": 2}
    if value, ok := dict["key1"]; ok {
        fmt.Printf(value)
    } else {
        fmt.Println("key1 不存在")
    }
}
```
之所以能这么写是因为，这是if判断的一种高级用法

上面这种写法的意思是，在 if 里先运行表达式
```
value, ok := dict["key1"]
```
，得到变量后，再对这个变量进行判断
