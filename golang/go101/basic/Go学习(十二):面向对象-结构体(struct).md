## Go学习(十二):面向对象-结构体(struct)
> Go没有沿袭传统面向对象编程中的诸多概念，也没有提供类(class)，但是它提供了结构体(struct)，方法(method)可以在结构体上添加。与类相似，结构体提供了捆绑数据和方法的行为。
### 1. 介绍
#### 1.1 概念
单一的数据类型已经满足不了现实开发需求，于是 Go 语言提供了结构体来定义复杂的数据类型。结构体是由一系列相同类型或不同类型的数据构成的数据集合。结构体的定义只是一种内存布局的描述，只有当结构体实例化时，才会真正地分配内存。因此必须在定义结构体并实例化后才能使用结构体的字段。
#### 1.2 语法
```
type 类型名称 struct {
   field type
   field type
   field1,field2,field3 type // 同类型变量可以写在一行
}
```
#### 1.3 注意事项
* 类型名是标识结构体的名称，在同一个包内不能重复。
* 结构体的属性，也叫字段(field)，必须唯一。
* 同类型的成员属性可以写在一行。
* 结构体是值类型。
* 只有当结构体实例化时,才能使用结构体的字段。
### 2. 实例化
#### 2.1 使用var
```
package main
import "fmt"
// 定义结构体
type student struct {
	name string
	age int
	like []string
}
func main() {
	// 使用var 实例化结构体
	var s student
	fmt.Printf("变量s--> 类型: %T 值: %v \n",s,s)
	// 给属性赋值
	s.name = "张三"
	s.age = 20
	s.like = []string{"打游戏","看动漫"}
	fmt.Printf("给属性赋值: 变量s--> 类型: %T 值: %v \n",s,s)
}
/*输出
变量s--> 类型: main.student 值: { 0 []}
给属性赋值: 变量s--> 类型: main.student 值: {张三 20 [打游戏 看动漫]}
*/
```
#### 3.2 使用简短声明 (:=)
```
package main
import "fmt"
// 定义结构体
type student struct {
	name string
	age int
	like []string
}
func main() {
	// 方式一: 先实例化结构体，后赋值
	s := student{}
	// 给属性赋值
	s.name = "张三"
	s.age = 20
	s.like = []string{"打游戏","看动漫"}
	fmt.Printf(" 变量s--> 类型: %T 值: %v \n",s,s)

	// 方式二: 声明时初始化
	s1 := student{
		name: "李四",
		age:  23,
		like: []string{"旅游","运动"},
	}
	fmt.Printf(" 变量s1--> 类型: %T 值: %v \n",s1,s1)

	// 方式三: 声明时初始化，省略属性
	s2 := student{"王麻子",30,[]string{"睡觉","吃饭"}}
	fmt.Printf(" 变量s2--> 类型: %T 值: %v \n",s2,s2)
}
/**输出
 变量s--> 类型: main.student 值: {张三 20 [打游戏 看动漫]}
 变量s1--> 类型: main.student 值: {李四 23 [旅游 运动]}
 变量s2--> 类型: main.student 值: {王麻子 30 [睡觉 吃饭]}
*/
```
#### 3.3 使用new
使用内置函数new()对结构体进行实例化，结构体实例化后形成指针类型的结构体，new()内置函数会分配内存。第一个参数是类型，而不是值，返回的值是指向该类型新分配的零值的指针。
```
package main
import "fmt"
// 定义结构体
type student struct {
	name string
	age int
	like []string
}
func main() {
	// 使用new实例化
	s := new(student)
	fmt.Printf(" 变量s--> 类型: %T 值: %v \n",s,s)
	// 给属性赋值
	(*s).name = "包青天"
	(*s).age = 55
	(*s).like = []string{"判案"}
	fmt.Printf(" 变量s--> 类型: %T 值: %v \n",s,s)
	// 语法糖写法(省略*)
	s.name = "包大人"
	s.age = 99
	s.like = []string{"判案","元芳你怎么看"}
	fmt.Printf(" 变量s--> 类型: %T 值: %v \n",s,s)
}
/**输出
 变量s--> 类型: *main.student 值: &{ 0 []}
 变量s--> 类型: *main.student 值: &{包青天 55 [判案]}
 变量s--> 类型: *main.student 值: &{包大人 99 [判案 元芳你怎么看]}
*/
```
### 3. 结构体在函数中使用
结构体作为函数参数，若复制一份传递到函数中，在函数中对参数进行修改，不会影响到实际参数，证明结构体是值类型。
#### 3.1 传结构体值作为参数
```
package main
import "fmt"
// 定义结构体
type student struct {
	name string
	age int
	like []string
}
func main() {
	s := student{name:"张三",age:17}
	fmt.Printf(" 变量s--> 值: %v \n",s)
	grownUp(s)
	fmt.Printf("调用函数后,变量s--> 值: %v \n",s)
}
// 传结构体值作为参数
func grownUp( s student)  {
	s.age = 80
	s.name = "长大的 "+s.name
}
/**输出
 变量s--> 值: {张三 17 []}
 调用函数后,变量s--> 值: {张三 17 []}
*/
```
#### 3.2 传结构体指针作为参数(其实也是值传递，传递的值是地址)
```
package main
import "fmt"
// 定义结构体
type student struct {
	name string
	age int
	like []string
}
func main() {
	s := student{name:"张三",age:17}
	fmt.Printf(" 变量s--> 值: %v \n",s)
  // 取址
	grownUp(&s)
	fmt.Printf("调用函数后,变量s--> 值: %v \n",s)
}
// 传结构体指针作为参数
func grownUp( s *student)  {
	s.age = 80
	s.name = "长大的 "+s.name
}
/** 输出:
 变量s--> 值: {张三 17 []}
 调用函数后,变量s--> 值: {长大的 张三 80 []}
*/
```
#### 3.3 返回对象
```
package main
import "fmt"
// 定义结构体
type student struct {
	name string
	age int
	like []string
}
func main() {
	s := getStudent("杨过",40,[]string{"骑大雕"})
	fmt.Printf("函数返回值 s--> 值: %v  类型: %T \n",s,s)
}
// 作为值类型传递
func getStudent( name string, age int,likes []string) student  {
	return student{name,age,likes}
}
// 函数返回值 s--> 值: {杨过 40 [骑大雕]}  类型: main.student
```
#### 3.4 返回指针
```
package main
import "fmt"
// 定义结构体
type student struct {
	name string
	age int
	like []string
}
func main() {
	s := getStudent("杨过",40,[]string{"骑大雕"})
	fmt.Printf("函数返回值 s--> 值: %v  类型: %T \n",s,s)
}

// 返回指针
func getStudent( name string, age int,likes []string) *student  {
	return &student{name,age,likes}
}
// 输出: 函数返回值 s--> 值: &{杨过 40 [骑大雕]}  类型: *main.student
```
### 4.匿名结构体
#### 4.1 语法
```
变量名 := struct {
  // 定义成员属性
} { /*初始化成员属性*/ }
```
#### 4.2 使用
```
package main
import "fmt"
func main() {
	// 声明初始化匿名结构体
	s := struct {
		name, home, phone string
		age               int
	}{
		name:  "张二十",
		phone: "17600111111",
		age:   18,
	}
  // 打印
	fmt.Printf("变量 s--> 值: %v  类型: %T \n", s, s)
}
// 输出: 变量 s--> 值: {张二十  17600111111 18}  类型: struct { name string; home string; phone string; age int }
```
### 5.匿名字段
#### 5.1 定义
匿名字段就是在结构体中的字段没有名字，只包含一个没有字段名的类型。这些字段被称为匿名字段。在同一个结构体中,同一个类型只能有一个匿名字段。
#### 5.2 使用
```
package main
import "fmt"
type people struct {
	name, home string
	int        // 匿名字段
	float32    // 匿名字段
}
func main() {
	// 声明初始化匿名结构体
	s := people{name: "张三", home: "北京", int: 18, float32: 1.73}
	fmt.Printf("变量 s--> 值: %v \n", s)
	// 声明初始化匿名结构体(省略属性名)
	s2 := people{"李四", "南京", 22, 1.80}
	fmt.Printf("变量 s2--> 值: %v \n", s2)
}
/** 输出
  变量 s--> 值: {张三 北京 18 1.73}
  变量 s2--> 值: {李四 南京 22 1.8}
*/
```
### 6. 结构体嵌套
#### 6.1 定义
将一个结构体作为另一个结构体的属性（字段），这种结构就是结构体嵌套。

结构体嵌套可以模拟面向对象编程中的以下两种关系。
* 聚合关系: 一个类作为另一个类的属性。
* 继承关系: 一个类作为另一个类的子类。子类和父类的关系。

#### 6.2 聚合场景
**模拟聚合关系时一定要采用有名字的结构体作为字段。**
```
package main
import "fmt"
// 定义学生结构体
type student struct {
	name       string
	height      float32
	schoolInfo  school
}
// 定义学习结构体
type school struct {
	schoolName, schoolAddress string
}
func main() {
	// 简短声明嵌套结构体
	s := student{"小张",1.72,school{"北京大学","北京"}}
	fmt.Printf("变量 s--> 值: %v 类型: %T \n", s,s)
	// 使用var
	var ss student
	ss.name = "小龙"
	ss.height = 1.67
	ss.schoolInfo = school{"南京大学","南京"}
	fmt.Printf("变量 ss--> 值: %v 类型: %T \n", ss,ss)
	// 使用new
	s2 := new(student)
	s2.name = "小虎"
	s2.height = 1.77
	s2.schoolInfo.schoolName = "武汉大学"
	s2.schoolInfo.schoolAddress = "武汉"
	fmt.Printf("变量 s2--> 值: %v 类型: %T \n", s2,s2)
}
/** 输出:
  变量 s--> 值: {小张 1.72 {北京大学 北京}} 类型: main.student
  变量 ss--> 值: {小龙 1.67 {南京大学 南京}} 类型: main.student
  变量 ss--> 值: &{小虎 1.77 {武汉大学 武汉}} 类型: *main.student
*/
```
#### 6.3 模拟继承
在结构体中，属于匿名结构体的字段称为提升字段，它们可以被访问，匿名结构体就像是该结构体的父类。
```
package main
import "fmt"

// 定义父类结构体
type people struct {
	name  string
	age   int
}
type student struct {
	people // 集成父类结构体
	class string
}

func main() {
	// 方式1.使用new声明结构体
	var s = new(student)
	// 集成父类成员
	s.name = "张三"
	s.age = 12
	// 自己成员
	s.class = "三年级"
	fmt.Printf("变量s -> %v \n",s)
	// 方式2.使用简短声明
	s2 := student{people{"李四",13},"四年级"}
	fmt.Printf("变量s2 -> %v \n",s2)
}
/** 输出
变量s -> &{{张三 12} 三年级}
变量s2 -> {{李四 13} 四年级}
*/
```
#### 6.4 成员冲突
```
package main
import "fmt"
type A struct {
	name string
	age  int
}
type B struct {
	name   string
	height float32
}
// 在C结构体中嵌套A和B
type C struct {
	A
	B
}
func main() {
	// 定义结构体C
	c := C{}
	// 不冲突的成员赋值
	c.age = 12
	c.height = 1.88
	// 冲突的成员赋值
	c.A.name = "这是A的成员"
	c.B.name = "这是B的成员"
	fmt.Printf("变量c -> %v \n", c)
}
// 输出: 变量c -> {{这是A的成员 12} {这是B的成员 1.88}}
```
转自：http://liuqh.icu/2020/08/26/go/basic/12-struct/

