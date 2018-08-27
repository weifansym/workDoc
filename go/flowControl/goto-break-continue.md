## 循环控制Goto、Break、Continue
循环控制语句

循环控制语句可以控制循环体内语句的执行过程。
GO 语言支持以下几种循环控制语句：

Goto、Break、Continue
```
1.三个语句都可以配合标签(label)使用
2.标签名区分大小写，定以后若不使用会造成编译错误
3.continue、break配合标签(label)可用于多层循环跳出
4.goto是调整执行位置，与continue、break配合标签(label)的结果并不相同
```
goto 语句 将控制转移到被标记的语句。

Go 语言的 goto 语句可以无条件地转移到过程中指定的行。
goto语句通常与条件语句配合使用。可用来实现条件转移， 构成循环，跳出循环体等功能。
但是，在结构化程序设计中一般不主张使用goto语句， 以免造成程序流程的混乱，使理解和调试程序都产生困难。
语法
goto 语法格式如下：
```
goto label;
..
.
label: statement;
```
Golang支持在函数内 goto 跳转。标签名区分大小写，未使用标签引发错误。
```
func main() {
    var i int
    for {
        println(i)
        i++
        if i > 2 { goto BREAK }
    }
BREAK:
    println("break")
EXIT:                 // Error: label EXIT defined and not used
}
```
goto 实例：
```
package main

import "fmt"

func main() {
   /* 定义局部变量 */
   var a int = 10

   /* 循环 */
   LOOP: for a < 20 {
      if a == 15 {
         /* 跳过迭代 */
         a = a + 1
         goto LOOP
      }
      fmt.Printf("a的值为 : %d\n", a)
      a++     
   }  
}
```
以上实例执行结果为：
```
a的值为 : 10
a的值为 : 11
a的值为 : 12
a的值为 : 13
a的值为 : 14
a的值为 : 16
a的值为 : 17
a的值为 : 18
a的值为 : 19
```
控制语句

break 语句 经常用于中断当前 for 循环或跳出 switch 语句

Go 语言中 break 语句用于以下两方面：
1.用于循环语句中跳出循环，并开始执行循环之后的语句。
2.break在switch（开关语句）中在执行一条case后跳出语句的作用。

实例：
```
package main

import "fmt"

func main() {
   /* 定义局部变量 */
   var a int = 10

   /* for 循环 */
   for a < 20 {
      fmt.Printf("a 的值为 : %d\n", a)
      a++
      if a > 15 {
         /* 使用 break 语句跳出循环 */
         break
      }
   }
}
```
