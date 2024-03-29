## Go学习(十三):面向对象-方法(method)
> Go没有沿袭传统面向对象编程中的诸多概念，也没有提供类(class)，但是它提供了结构体(struct)，方法(method)可以在结构体上添加。与类相似，结构体提供了捆绑数据和方法的行为。

### 1.介绍
#### 1.1 概念
Go语言同时有函数和方法，方法的本质是函数，但是方法和函数又有所不同。
#### 1.2 方法和函数的区别
* 函数(function)是一段具有独立功能的代码，可以被反复多次调用，从而实现代码复用。
* 方法(method)是一个类的行为功能，只有该类的对象才能调用。
* 方法有接受者，而函数无接受者。
* 函数不可以重名，而方法可以重名。**只要接受者不同，方法名就可以相同。**

### 2. 使用
#### 2.1 基本语法
```
func (接收器变量 接收器类型) 方法名(参数列表) (返回参数) {
     函数体
}
```
* 接收器变量：接收器变量在命名时，官方建议使用接收器类型的第一个小写字母，而不是self、this之类的命名。例如: **Socket类型的接收器变量应该命名为s，Connector类型的接收器变量应该命名为c等。**
* 接收器类型：接收器类型和参数类似，可以是指针类型和非指针类型。
* 方法名、参数列表、返回参数：格式与函数定义一致。

#### 2.2 使用示例
```
package main
import "fmt"
// 定义一个结构体
type Student struct {
	name string
	age  int
}
// 定义一个方法(接收器为Student的指针)
func (s *Student)updateName(newName string)  {
	s.name = newName
}
// 定义一个方法(接收器为Student)
func (s Student) updateAge(newAge int)  {
	s.age = newAge
	fmt.Printf("修改结构体s的age -> %v \n",s)
}
func main() {
	// 初始化结构体
	s := Student{"张三",20}
    fmt.Printf("结构体初始化s -> %v \n",s)

	// 通过方法修改名称
	s.updateName("张三新名")
	fmt.Printf("调用updateName后 -> %v \n",s)

	// 通过方法修改年龄
	s.updateAge(22)
	fmt.Printf("调用updateAge后 -> %v \n",s)
}
/** 输出:
结构体初始化s -> {张三 20} 
调用updateName后 -> {张三新名 20} 
修改结构体s的age -> {张三新名 22} 
调用updateAge后 -> {张三新名 20} 
*/
```
**通过上述示例可以看出: ** 若方法的接受者不是指针，实际只是获取了一个拷贝，而不能真正改变接受者中原来的数据。

### 3.方法继承
方法是可以继承的，如果匿名字段实现了一个方法，那么包含这个匿名字段的struct也能调用该匿名字段中的方法。
#### 3.1 使用示例
```
package main
import "fmt"
// 定义一个人类结构体
type People struct {
	name, position string
	age            int
}
// 定义一个学生结构体
type Student struct {
	People
}
type Teacher struct {
	People
}
// 定义一个方法
func (p People) say() {
	fmt.Printf("我叫 %s  %d岁 从事: %s \n", p.name,p.age,p.position)
}
func main() {
	student := Student{People{"张三","学生",15}}
	teacher := Teacher{People{"李杨","老师",35}}
	// 调用方法(继承父类)
	student.say()
	teacher.say()
}
/** 输出
  我叫 张三  15岁 从事: 学生 
  我叫 李杨  35岁 从事: 老师
*/
```
### 4 .方法重写
在Go语言中，方法重写是指一个包含了匿名字段的struct也实现了该匿名字段实现的方法（即子类也实现了父类的方法）

#### 4.1 使用示例
```
package main
import "fmt"
// 定义一个人类结构体
type People struct {
	name, position string
	age            int
}
// 定义一个学生结构体
type Student struct {
	People
}
// 定义一个老师结构体
type Teacher struct {
	People
}
// 定义一个方法
func (p People) say() {
	fmt.Printf("我叫 %s  %d岁 从事: %s \n", p.name, p.age, p.position)
}
// 学生(子类)重写People(父类)的say方法
func (s Student) say() {
	fmt.Printf("我是一名学生,名字叫: %s 今年: %d岁 \n", s.name, s.age)
}
func main() {
	student := Student{People{"张三", "学生", 15}}
	teacher := Teacher{People{"李杨", "老师", 35}}
	// 调用方法(重写父类方法)
	student.say()
	// 调用方法(继承父类)
	teacher.say()
}
/** 输出
  我是一名学生,名字叫: 张三 今年: 15岁 
  我叫 李杨  35岁 从事: 老师 
*/
```
**当结构体存在继承关系时，方法调用按照就近原则。**

转自：http://liuqh.icu/2020/08/27/go/basic/13-method/





