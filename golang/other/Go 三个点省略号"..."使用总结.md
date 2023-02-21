## Go 三个点省略号"..."使用总结
### 使用在数组中
```
q := [...]int{1,2,3}
fmt.Printf("%T\n",q) //"[3]int"
```
在数组字面量中，如果省略号"..."出现在数组长度的位置，那么数组的长度由初始化数组的元素个数决定。
### 打散Slice
```
package main
 
import (
	"fmt"
)
 
func main() {
    var arr1 []int
    arr2 := []int{1,2,3}
    arr1 = append(arr1,0)	
    arr1 = append(arr1,arr2...)	 //arr2... 将切片arr2打散成 ==> arr1 = append(arr1,1,2,3)
    fmt.Printf("%v\n",arr1)
 
    var arr3 []byte
    arr3 = append(arr3,[]byte("hello")...)  
    fmt.Printf("%s\n",arr3)
}
 
```
运行结果：
```
[0,1,2,3]
hello
```
上面例子中append函数的参数后面的省略号表示如何将一个slice转换为参数列表
### 变长的函数参数
```
package main
 
import (
	"fmt"
)
 
func f1(parms ...int){
    for i,v := range parms {
	    fmt.Printf("%v %v\n",i,v)
	}
}
 
func main() {
   f1(0，1，2)
}
```
在参数列表最后的类型名称之前使用省略号“...”表示声明一个变长函数，调用这个函数的时候可以传递该类型任意数目的参数。

**尽管...type参数就像函数体内的slice，但变长函数的类型和带有普通slice参数的函数类型不相同，所以在传参的时候也是有所区别**。...type格式的类型只能作为函数的参数类型存在，
并且必须是最后一个参数。它是一个语法糖（syntactic sugar），即这种语法对语言的功能并没有影响，但是更方便程序员使用。

如果不使用...type，则必须这样：
```
package main
 
import (
	"fmt"
)
 
func f1(parms []int){
    for i,v := range parms {
	    fmt.Printf("f2:%v %v\n",i,v)
	}
}
 
func main() {
   b := []int{0,1,2}
   f1(b)  
}
```
结果与上面一样。
****

以上算是比较基本的用法，但是实际工作中可能会遇到一些别的情况。

#### 情况1
```
package main
 
import (
	"fmt"
)
 
func f1(parms ...int){
    fmt.Printf("%T\n",parms)
    for i,v := range parms {
	    fmt.Printf("f1:%v %v\n",i,v)
	}
}
 
func main() {
   b := []int{0,1,2}
   //f1(b) //error
   f1(b...)  
}
```
上面已经说明:尽管...type参数就像函数体内的slice，但变长函数的类型和带有普通slice参数的函数类型不相同。所以f1(b)会发生错误，因为类型不匹配。
而这里在slice后面加...，将slice打散，传入f1。

同时也可以知道，parms类型是[]int，所以在函数f1里面，parms就是一个int类型的切片。
#### 情况2
```
package main
 
import (
	"fmt"
)
 
func f2(a ...int){
    fmt.Printf("f2 %T\n",a)
    fmt.Printf("f2 %v\n",a)
    fmt.Printf("f2 %p\n",a)
}
 
func f1(a ...int){
    fmt.Printf("f1 %v\n",a)
	fmt.Printf("f1 %p\n",a)
	//f2(a) //error
	f2(a...)
}
 
func main() {
    f1(1,2,3)
}
```
从情况一知道，在f1中a是slice，再次传入到f2函数中，必须将slice打散操作，不然发生类型不匹配的错误。

但是我们把 ...int 改为 ...interface{} ，再看下面的例子：
```
package main
 
import (
	"fmt"
)
 
func f2(a ...interface{}){
    fmt.Printf("f2 %T\n",a)
    fmt.Printf("f2 %v\n",a)
    fmt.Printf("f2 %p\n",a)
}
 
func f1(a ...interface{}){
    fmt.Printf("f1 %v\n",a)
	fmt.Printf("f1 %p\n",a)
	f2(a)
	f2(a...)
}
 
func main() {
    f1(1,2,3)
}
```
发现在f2(a)没有发生错误，但是仔细看看打印的结果[[1 2 3]]是数组为元素的数组类型。这都归功于interface{}这个万能类型，它将f2(a)中的a作为一个参数传入，a在f2函数中，
是[]interface{}类型。比如：
```
package main
 
import (
	"fmt"
)
 
func f1(a ...interface{}){
    fmt.Printf("f1 %v\n",a)
}
 
func main() {
    arr := []int{1,2,3}
    arr2 := []int{11,22,33}
    f1(arr,arr2,111,222,333)
}
```
#### 情况3
```
package main
 
import (
	"fmt"
)
 
func f1(a ...interface{}){
    fmt.Printf("f1 %v\n",a)
}
 
func main() {
    arr := []int{1,2,3}
    f1(arr...)
}
```
这和情况一看起来差不多，但是这里怎么出错了？

我们从错误提示可以看出：[]int 不能作为 []interface{}类型传入f1。换句话说，就是类型不匹配。

但是interface{}不是万能类型嘛？从情况二可以看出。

这里可能是一个误区，虽然interface{}是万能类型，但是[]interface{}并不是万能类型，它本质是slice，只不过slice的成员都是万能类型interface{}，就像[]int和[]string是不同类型slice，
所以错误是因为类型不匹配。也可以看看官方怎么解释[]interface{}。

结合下面例子再理解。
```
package main
 
import (
	"fmt"
)
 
func f2(a interface{}){
    fmt.Printf("f2 %v\n",a)
}
 
func f3(a []interface{}){
    fmt.Printf("f3 %v\n",a)
}
 
func f4(a []int){
    fmt.Printf("f4 %v\n",a)
}
 
func main() {
    arr := []int{1,2,3}
 
    f2(arr)  //[]int -> interface{}	
 
    //f3(arr)	//error 类型不匹配[]int -> []interface{}
	
    f4(arr)
}
```
那如何修改情况三第一个例子，如下：
```
package main
 
import (
	"fmt"
)
 
func f1(a ...interface{}){
    fmt.Printf("f1 %v\n",a)
}
 
func main() {
    var arr []interface{}
    arr = append(arr,1,2,3)
    f1(arr...)
}
```
以上大概就是Go语言的省略号"..."的所有用法和遇到的问题了。

参考文章：

http://c.biancheng.net/view/60.html
https://blog.csdn.net/jeffrey11223/article/details/79166724






