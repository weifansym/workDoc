## for-select中的break、continue和return
### break
* select中的break，类似c系列中的break，break后的语句不执行
* **for和select一同使用，有坑**

break只能跳出select，无法跳出for
```
package test

import (
	"fmt"
	"testing"
	"time"
)

func TestBreak(t *testing.T) {
	tick := time.Tick(time.Second)
	for {
		select {
		case t := <-tick:
			fmt.Println(t)    
			break
		}
	}
	fmt.Println("end")
}
```
执行结果：
```
=== RUN   TestBreak
2019-12-19 14:43:41.7912242 +0800 CST m=+1.005627701
2019-12-19 14:43:42.0862832 +0800 CST m=+1.007127901
2019-12-19 14:43:42.7914242 +0800 CST m=+2.005754701
2019-12-19 14:43:43.0864832 +0800 CST m=+2.007254901
...
```
break无法跳出select的解决方案

1、标签
```
func TestBreak(t *testing.T) {
	tick := time.Tick(time.Second)
//FOR是标签
FOR:
	for {
		select {
		case t := <-tick:
			fmt.Println(t)
			//break出FOR标签标识的代码
			break FOR
		}
	}
	fmt.Println("end")
}
```
2、goto
```

func TestBreak(t *testing.T) {
	tick := time.Tick(time.Second)
	for {
		select {
		case t := <-tick:
			fmt.Println(t)
			//跳到指定位置
			goto END
		}
	}
END:
	fmt.Println("end")
}
```
执行结果：
```
=== RUN   TestBreak
2019-12-19 14:43:41.7912242 +0800 CST m=+1.005627701
end
```
### continue
**单独在select中是不能使用continue，会编译错误，只能用在for-select中。**
continue的语义就类似for中的语义，select后的代码不会被执行到。
```
func TestBreak(t *testing.T) {
	tick := time.Tick(time.Second)
	for {
		select {
		case t := <-tick:
			fmt.Println(t)
			continue
			fmt.Println("test")
		}
	}
	fmt.Println("end")
}
```
执行结果：
```
=== RUN   TestBreak
2019-12-19 14:43:41.7912242 +0800 CST m=+1.005627701
2019-12-19 14:43:42.0862832 +0800 CST m=+1.007127901
2019-12-19 14:43:42.7914242 +0800 CST m=+2.005754701
2019-12-19 14:43:43.0864832 +0800 CST m=+2.007254901
...
```
### return
和函数中的return一样，跳出select，和for，后续代码都不执行


