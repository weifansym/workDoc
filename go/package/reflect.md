## reflect-反射
具体使用请参见官方文档：[reflect](https://golang.org/pkg/reflect/)
### 获取基本类型
使用反射获取基本类型，反射：可以在运行时动态获取变量的相关信息。
```
* reflect.TypeOf()，获取变量的类型，返回reflect.Type类型 
* reflect.ValueOf()，获取变量的值，返回reflect.Value类型 
* reflect.Value.Kind()，获取变量的类别，返回一个常量 
* reflect.Value.Interface()，转换成interface{}类型
```
```
package main

import (
	"fmt"
	"reflect"
)

func main() {
	var x float64 = 3.4
	fmt.Println("type:", reflect.TypeOf(x))
	v := reflect.ValueOf(x)
	fmt.Println("value:", v)
	fmt.Println("type:", v.Type())
	fmt.Println("kind:", v.Kind())
	fmt.Println("value:", v.Float())

	fmt.Println(v.Interface())
	fmt.Printf("value is %5.2e\n", v.Interface())
	y := v.Interface().(float64)
	fmt.Println(y)
}
```
输出如下：
```
type: float64
value: 3.4
type: float64
kind: float64
value: 3.4
3.4
value is 3.40e+00
3.4
```
### 反射获取结构体
示例如下：
```
package main

import (
	"fmt"
	"reflect"
)

type Student struct {
	Name  string
	Age   int
	Score float32
}

func test(b interface{}) {
	t := reflect.TypeOf(b)
	fmt.Println(t)

	v := reflect.ValueOf(b)
	fmt.Println(v)

	k := v.Kind()
	fmt.Println(k)

	iv := v.Interface()
	fmt.Println(iv)

	stu, ok := iv.(Student)
	if ok {
		fmt.Printf("%v %T\n", stu, stu)
	}
}

func main() {
	var a Student = Student{
		Name:  "stu01",
		Age:   18,
		Score: 92,
	}
	test(a)
}
```
输出结果：
```
main.Student
{stu01 18 92}
struct
{stu01 18 92}
{stu01 18 92} main.Student
```
### Elem反射操作基本类型
Elem反射操作基本类型，用来获取指针指向的变量，相当于： var a *int;
实例如下：
```
package main

import (
	"fmt"
	"reflect"
)

func main() {

	var b int = 1
	b = 200
	testInt(&b)
	fmt.Println(b)
}

//fv.Elem()用来获取指针指向的变量
func testInt(b interface{}) {
	val := reflect.ValueOf(b)
	val.Elem().SetInt(100)
	c := val.Elem().Int()

	fmt.Printf("get value  interface{} %d\n", c)
	fmt.Printf("string val:%d\n", val.Elem().Int())
}
```
输出结果：
```
get value  interface{} 100
string val:100
100
```
### 反射调用结构体方法
实例如下：
```
package main

import (
	"fmt"
	"reflect"
)

type Student struct {
	Name  string
	Age   int
	Score float32
}

func (s Student) Print() {
	fmt.Println(s)
}

func (s Student) Set(name string, age int, score float32) {
	s.Age = age
	s.Name = name
	s.Score = score
}

func TestStruct(a interface{}) {
	val := reflect.ValueOf(a)
	kd := val.Kind()

	fmt.Println(val, kd)
	if kd != reflect.Struct {
		fmt.Println("expect struct")
		return
	}
	//获取字段数量
	fields := val.NumField()
	fmt.Printf("struct has %d field\n", fields)
	//获取字段的类型
	for i := 0; i < fields; i++ {
		fmt.Printf("%d %v\n", i, val.Field(i).Kind())
	}
	//获取方法数量
	methods := val.NumMethod()
	fmt.Printf("struct has %d methods\n", methods)

	//反射调用的Print方法
	var params []reflect.Value
	val.Method(0).Call(params)

}

func main() {
	var a Student = Student{
		Name:  "stu01",
		Age:   18,
		Score: 92.8,
	}
	TestStruct(a)
	// fmt.Println(a)
}
```
输出结果：
```
{stu01 18 92.8} struct
struct has 3 field
0 string
1 int
2 float32
struct has 2 methods
{stu01 18 92.8}
```
### Elem反射操作结构体

```
package main

import (
	"fmt"
	"reflect"
)

type Student struct {
	Name  string
	Age   int
	Score float32
}

func (s Student) Print() {
	fmt.Println(s)
}

func TestStruct(a interface{}) {
	val := reflect.ValueOf(a)
	kd := val.Kind()

	fmt.Println(val, kd)
	if kd != reflect.Ptr && val.Elem().Kind() == reflect.Struct {
		fmt.Println("expect struct")
		return
	}
	//获取字段数量
	fields := val.Elem().NumField()
	fmt.Printf("struct has %d field\n", fields)
	//获取字段的类型
	for i := 0; i < fields; i++ {
		fmt.Printf("%d %v\n", i, val.Elem().Field(i).Kind())
	}
	//获取方法数量
	methods := val.NumMethod()
	fmt.Printf("struct has %d methods\n", methods)

	//反射调用的Print方法
	var params []reflect.Value
	val.Elem().Method(0).Call(params)
}

func main() {
	var a Student = Student{
		Name:  "stu01",
		Age:   18,
		Score: 92.8,
	}
	TestStruct(&a)
	// fmt.Println(a)
}
```
输出如下：
```
&{stu01 18 92.8} ptr
struct has 3 field
0 string
1 int
2 float32
struct has 2 methods
{stu01 18 92.8}
```
### Elem反射获取tag
```
package main

import (
	"fmt"
	"reflect"
)

type Student struct {
	Name  string `json:"stu_name"`
	Age   int
	Score float32
}

func TestStruct(a interface{}) {
	typ := reflect.TypeOf(a)

	tag := typ.Elem().Field(0).Tag.Get("json")
	fmt.Printf("Tag:%s\n", tag)
}

func main() {
	var a Student = Student{
		Name:  "stu01",
		Age:   18,
		Score: 92.8,
	}
	TestStruct(&a)
}
```
输出结果：
```
Tag:stu_name
```
获取 tag（一）
```
package main

import (
	"fmt"
	"reflect"
)

func main() {
	type User struct {
		Name   string "user name"
		Passwd string `user passsword`
	}
	u := &User{
		Name:   "Murphy",
		Passwd: "123456",
	}
	s := reflect.TypeOf(u).Elem()
	for i := 0; i < s.NumField(); i++ {
		fmt.Println(s.Field(i).Tag)
	}
}
```
获取 tag（二）
```
package main

import (
	"fmt"
	"reflect"
)

func main() {
	type User struct {
		Name string `json:"user_name" name:"user name"`
	}
	u := User{
		Name: "Murphy",
	}
	f := reflect.TypeOf(u).Field(0)
	fmt.Println(f.Tag.Get("json"))
	fmt.Println(f.Tag.Get("name"))
}
```
### 应用demo
练习：
1.定义一个结构体
2.给结构体赋值
3.用反射获取结构体的 下标、结构体名称、类型、值
4.改变结构体的值

代码如下：
```
package main

import (
    "fmt"
    "reflect"
)

type T struct {
    A int
    B string
}

func main() {
    t := T{23, "skidoo"}
    s := reflect.ValueOf(&t).Elem()
    typeOfT := s.Type()
    for i := 0; i < s.NumField(); i++ {
        f := s.Field(i)
        fmt.Printf("%d: %s %s = %v\n", i,
            typeOfT.Field(i).Name, f.Type(), f.Interface())
    }
    s.Field(0).SetInt(77)
    s.Field(1).SetString("Sunset Strip")
    fmt.Println("t is now", t)
}
```
输出结果：
```
0: A int = 23
1: B string = skidoo
t is now {77 Sunset Strip}
```



