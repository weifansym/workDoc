## 详解 Go 语言中的 init () 函数
* 变量初始化
* 检查 / 修复状态
* 注册器
* 运行计算

### 包初始化
为了使用导入的程序包，必须首先对其进行初始化。初始化始终在单个线程中执行，并且以程序包依赖关系的顺序执行。这由 Golang 的运行时系统控制，如下图所示：
* 初始化导入的包（递归导入）
* 计算并为块中声明的变量分配初始值
* 在包中执行初始化函数

![image](https://user-images.githubusercontent.com/6757408/163231405-7f0b08dd-c9c8-4cbd-80bc-1ddd1ff5b726.png)

initial.go
```

package main
import "fmt"
var _ int64=s()
func init(){
  fmt.Println("init function --->")
}
func s() int64{
  fmt.Println("function s() --->")
  return 1
}
func main(){
  fmt.Println("main --->")
}
```
执行结果
```
function s() —>
init function —>
main —>
```
即使程序包被多次导入，初始化也只需要一次。

### 特性
init 函数不需要传入参数，也不需要返回任何值。与 main 相比，init 没有声明，因此无法引用。
```
package main
import "fmt"
func init(){
  fmt.Println("init")
}
func main(){
  init()
}
```
编译上述函数 “undefined：init” 时发生错误。

每个源文件可以包含一个以上的 init 函数，请记住，写在每个源文件中的 “行进方式” 只能包含一个 init 函数，这有点不同，因此进行下一个验证。
```
package main
import "fmt"
func init(){
  fmt.Println("init 1")
}
func init(){
  fmt.Println("init2")
}
func main(){
  fmt.Println("main")
}
/* 实施结果:
init1
init2
main */
```
从上面的示例中，您可以看到每个源文件可以包含多个 init 函数。

初始化函数的一个常见示例是设置初始表达式的值。
```
var precomputed=[20]float64{}
func init(){
  var current float64=1
  precomputed[0]=current
  for i:=1;i<len(precomputed);i++{
    precomputed[i]=precomputed[i-1]*1.2
  }
}
```
因为不可能在上面的代码 (这是一条语句) 中将 for 循环用作预先计算的值，所以可以使用 init 函数来解决此问题。

Go 套件汇入规则的副作用

Go 非常严格，不允许引用未使用的软件包。但是有时您引用包只是为了调用 init 函数进行一些初始化。空标识符 (即下划线) 的目的是解决此问题。
```
import _ "image/png"
```
原文如下：https://developpaper.com/detailed-explanation-of-init-function-in-go-language/
