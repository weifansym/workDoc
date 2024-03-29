## Go学习(十四):面向对象-接口(interface)
> 虽然Go语言没有继承和多态，但是Go语言可以通过匿名字段实现继承，通过接口实现多态。
### 1.介绍
#### 1.1 概念
在Go语言中，接口是一组方法签名。接口指定了类型应该具有的方法，类型决定了如何实现这些方法。当某个类型为接口中的所有方法提供了具体的实现细节时，这个类型就被称为实现了该接口。接口定义了一组方法，如果某个
对象实现了该接口的所有方法，则此对象就实现了该接口。
#### 1.2 声明语法
```
type 接口名称 interface {
    Method1([参数列表]) [返回值列表]
    Method2([参数列表]) [返回值列表]
    ...
}
```
示例
```
// 定义一个接口
type Bird interface {
	fly() // 无参数无返回值方法
	eat(string2 string) // 有参数无返回值方法
	walk(string2 string) string // 有参数有返回值方法
}
```
### 2.定义和实现
#### 2.1 定义接口
```
// 定义一个鸟类接口
type Birder interface {
	fly()  
	eat(food string)
}
```
### 2.2 实现接口
**Go没有implements或extends关键字,类型都是隐式实现接口的。任何定义了接口中所有方法的类型都被称为隐式地实现了该接口。**
```
package main
import "fmt"
// 定义一个鸟类接口
type Birder interface {
	fly()
	eat(food string)
}
// 定义乌鸦结构体
type Crow struct {
	name string
}
// -------- 下面开始实现Bird接口 ------
func (c Crow) fly() {
	fmt.Printf("我是 %s,我会飞....\n", c.name)
}
func (c Crow) eat(food string) {
	fmt.Printf("我是 %s,我喜欢吃 %s \n", c.name,food)
}
// -------- 实现鸟类接口的所有方法，就代表实现了接口 ------
func main() {
	crow := Crow{"乌鸦"}
	crow.fly()
	crow.eat("谷子")
}
/** 输出
  我是 乌鸦,我会飞....
  我是 乌鸦,我喜欢吃 谷子 
*/
```
### 3. 模拟多态
#### 3.1 什么是多态?
如果有几个相似而不完全相同的对象，有时人们要求在向它们发出同一个消息时，它们的反应各不相同，分别执行不同的操作，这种情况就是多态现象。
**Go语言中的多态性是在接口的帮助下实现的——定义接口类型，创建实现该接口的结构体对象。**
#### 3.2 使用示例
> 定义接口类型的对象，可以保存实现该接口的任何类型的值。Go语言接口变量的这个特性实现了Go语言中的多态性。

实现: 写一个函数，接收不同类型的结构体，并打印不同其方法。
```
package main
import "fmt"
// 定义一个飞行器接口
type Flying interface {
	getName() string
}
// 定义小鸟结构体
type bird struct {
  name string
}
// bird实现接口
func (b bird) getName() string {
	return b.name
}
// 定义飞机结构体
type aircraft struct {
	name string
}
// aircraft实现接口
func (a aircraft) getName() string {
	return a.name
}
// 定义ufo结构体
type ufo struct {
	name string
}
// ufo实现接口
func (u ufo) getName() string {
	return u.name
}
// 写一个函数，接收不同类型的结构体，并打印不同其方法。
func print(flyList []Flying)  {
	for _,v := range flyList {
		fmt.Printf("我是%s,我会飞.....\n",v.getName())
	}
}
func main() {
	// 定义一个接口切片
	flyList := make([]Flying,0,3)
	bird := bird{"小鸟"}
	aircraft := aircraft{"飞机"}
	ufo := ufo{"UFO"}
	flyList = append(flyList, bird,aircraft,ufo)
	fmt.Printf("len: %d cap:%d val: %v \n",len(flyList),cap(flyList),flyList)
	// 调用函数
	print(flyList)
}
/** 输出
  len: 3 cap:3 val: [{小鸟} {飞机} {UFO}] 
  我是小鸟,我会飞.....
  我是飞机,我会飞.....
  我是UFO,我会飞.....
*/
```
### 4.空接口
空接口是接口类型的特殊形式，空接口没有任何方法，因此任何类型都无须实现空接口。从实现的角度看，任何值都满足这个接口的需求。因此空接口类型可以保存任何值，也可以从空接口中取出原值。
#### 4.1 定义空接口
```
// 定义一个空接口
type A interface {}
```
#### 4.2 保存任意类型
```
package main
import "fmt"
// 定义一个空接口
type A interface {
}
func main() {
	// 声明变量
	var a A
	// 保存整型
	a = 10
	fmt.Printf("保存整型: %v \n", a)
	// 保存字符串
	a = "hello word"
	fmt.Printf("保存字符串: %v \n", a)
	// 保存数组
	a = [3]float32{1.0, 2.0, 3.0}
	fmt.Printf("保存数组: %v \n", a)
	// 保存切片
	a = []string{"您", "好"}
	fmt.Printf("保存切片: %v \n", a)
	// 保存Map
	a = map[string]int{
		"张三": 22,
		"李四": 25,
	}
	fmt.Printf("保存map: %v \n", a)
	// 保存结构体
	a = struct {
		name string
		age  int
	}{"王麻子", 40}
	fmt.Printf("保存结构体: %v \n", a)

	// 声明一个空接口切片
	var aa []A
	// 保存任意类型数据到切片中
	aa = append(aa, 23, []string{"php", "go"}, map[string]int{"a": 1, "b": 2}, struct {
		city,province string

	}{"合肥","安徽"})
	fmt.Printf("空接口切片: %v \n", aa)
}

/** 输出:
  保存整型: 10 
  保存字符串: hello word 
  保存数组: [1 2 3] 
  保存切片: [您 好] 
  保存map: map[张三:22 李四:25] 
  保存结构体: {王麻子 40} 
  空接口切片: [23 [php go] map[a:1 b:2] {合肥 安徽}] 
*/
```
#### 4.3 从空接口中取值
保存到空接口的值，如果直接取出指定类型的值时，会发生编译错误。

错误示例:
```
package main
import "fmt"
// 定义一个空接口
type I interface {
}

func main() {
	// 声明变量num
	num := 10
	// 把变量num存到空接口中
	var i I = num
	fmt.Printf("输出变量i: %v \n", i)
  // 从空接口中取出值，赋值给新的变量 
	var c int = i // (!!! 这里会报错)
	fmt.Printf("输出变量c: %v \n", c)
}
/** 输出
 ./main.go:16:6: cannot use i (type I) as type int in assignment: need type assertion
*/
```
正确示例:
```
package main
import "fmt"
// 定义一个空接口
type I interface {
}
func main() {
	// 声明变量num
	num := 10
	// 把变量num存到空接口中
	var i I = num
	fmt.Printf("输出变量i: %v \n", i)
	// 从空接口中取出值，赋值给新的变量
	var c int = i.(int)
	fmt.Printf("输出变量c: %v \n", c)
}
/**
  输出变量i: 10 
  输出变量c: 10 
*/
```
### 5. 接口对象转换
#### 5.1 转换语法
```
// 方式一
instance,ok := 接口对象.(实际类型)
// 方式二
接口对象.(实际类型)
```
#### 5.2 使用示例
```
package main
import "fmt"
// 定义一个空接口
type I interface {
}

func main() {
	// 声明变量
	var a I
	// 保存整型
	a = 10
	printType(a)
	printType2(a)
	// 保存字符串
	a = "hello word"
	printType(a)
	printType2(a)
	// 保存数组
	a = [3]float32{1.0, 2.0, 3.0}
	printType(a)
	printType2(a)
	// 保存切片
	a = []string{"您", "好"}
	printType(a)
	printType2(a)
	// 保存Map
	a = map[string]string{
		"张三": "男",
		"小丽": "女",
	}
	printType(a)
	printType2(a)
	// 保存结构体
	a = people{"刘山", 32}
	printType(a)
	printType2(a)
}
// 定义结构体
type people struct {
	name string
	age  int
}

// 方式一
func printType(i I) {
	if t, ok := i.(int); ok {
		echo(t)
	} else if t, ok := i.(string); ok {
		echo(t)
	} else if t, ok := i.(map[string]string); ok {
		echo(t)
	} else if t, ok := i.([]int); ok {
		echo(t)
	} else if t, ok := i.([3]string); ok {
		echo(t)
	} else if t, ok := i.(people); ok {
		echo(t)
	}
}

// 方式二
func printType2(i I) {
	switch i.(type) {
	case int:
		echo2(i)
	case string:
		echo2(i)
	case map[string]string:
		echo2(i)
	case []int:
		echo2(i)
	case [3]string:
		echo2(i)
	case people:
		echo2(i)
	}
}
func echo(i interface{}) {
	fmt.Printf("方式一 ---> 变量i类型: %T 值: %v \n", i, i)
}
func echo2(i interface{}) {
	fmt.Printf("方式二 ---> 变量i类型: %T 值: %v \n", i, i)
}
/**输出
方式一 ---> 变量i类型: int 值: 10 
方式二 ---> 变量i类型: int 值: 10 
方式一 ---> 变量i类型: string 值: hello word 
方式二 ---> 变量i类型: string 值: hello word 
方式一 ---> 变量i类型: map[string]string 值: map[小丽:女 张三:男] 
方式二 ---> 变量i类型: map[string]string 值: map[小丽:女 张三:男] 
方式一 ---> 变量i类型: main.people 值: {刘山 32} 
方式二 ---> 变量i类型: main.people 值: {刘山 32} 
*/
```
#### 6. 使用注意事项
* 接口本身不能创建实例,但是可以指向一个实现了该接口的自定义类型的变量(实例)
* 接口中所有的方法都没有方法体,即都是没有实现的方法。
* 在 Go中，一个自定义类型需要将某个接口的所有方法都实现，我们说这个自定义类型实现 了该接口。
* 一个自定义类型只有实现了某个接口，才能将该自定义类型的实例(变量)赋给接口类型
* 只要是自定义数据类型，就可以实现接口，不仅仅是结构体类型
* 一个自定义类型可以实现多个接口
* Go接口中不能有任何变量
* interface类型默认是一个指针(引用类型)，如果没有对interface初始化就使用，那么会输出nil
* 空接口 interface{} 没有任何方法，所以所有类型都实现了空接口, 即我们可以把任何一个变量 赋给空接口
* 



