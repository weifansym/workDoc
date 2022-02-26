## Go语言中结构体打Tag是什么意思？
### 前言
哈喽，大家好，我是asong。今天想与大家分享Go语言中结构体标签是怎么使用的，以及怎样定制自己的结构体标签解析。

大多数初学者在看公司的项目代码时，看到的一些结构体定义会是这样的：
```
type Location struct {
 Longitude float32 `json:"lon,omitempty"`
 Latitude  float32 `json:"lat,omitempty"`
}
```
字段后面会有一个标签，这个标签有什么用呢？

上面的例子中，标签json:"lon,omitempty"代表的意思是结构体字段的值编码为json对象时，每一个导出字段变成该对象的一个成员，这个成员的名字为lon或者lat，并且当字段是空值时，不导出该字段；
总结就是lon、lat是重命名成员的名字，omitempty用来决定成员是否导出。

看到这里，有一些朋友可能会好奇，这个你是怎么知道这样使用的呢？我可以随便写标签吗？

接下来我们就一点点来揭秘，开车！！！
### 什么是标签
Go语言提供了可通过反射发现的的结构体标签，这些在标准库json/xml中得到了广泛的使用，orm框架也支持了结构体标签，上面那个例子的使用就是因为encoding/json支持了结构体标签，不过他有自己的标签规则；
但是他们都有一个总体规则，这个规则是不能更改的，具体格式如下：
```
`key1:"value1" key2:"value2" key3:"value3"...`  // 键值对用空格分隔
```
结构体标签可以有多个键值对，键与值要用冒号分隔，值要使用双引号括起来，多个键值对之间要使用一个空格分隔，千万不要使用逗号！！！

如果我们想要在一个值中传递多个信息怎么办？不同库中实现的是不一样的，在encoding/json中，多值使用逗号分隔：
```
`json:"lon,omitempty"`
```
在gorm中，多值使用分号分隔：
```
`gorm:"column:id;primaryKey"
```
具体使用什么符号分隔需要大家要看各自库的文档获取。

结构体标签是在编译阶段就和成员进行关联的，以字符串的形式进行关联，在运行阶段可以通过反射读取出来。

现在大家已经知道什么是结构体标签了，规则还是很规范的，但是很容易出错，因为Go语言在编译阶段并不会对其格式做合法键值对的检查，这样我们不小心写错了，就很难被发现，不过我们有**go vet**工具做检查，
具体使用来看一个例子：
```
type User struct {
 Name string `abc def ghk`
 Age uint16 `123: 232`
}
func main()  {
}
```
然后执行**go vet main.go，** 得出执行结果：
```
# command-line-arguments
go_vet_tag/main.go:4:2: struct field tag `abc def ghk` not compatible with reflect.StructTag.Get: bad syntax for struct tag pair
go_vet_tag/main.go:5:2: struct field tag `123: 232` not compatible with reflect.StructTag.Get: bad syntax for struct tag value
```
**bad syntax for struct tag pair**告诉我们键值对语法错误，bad syntax for struct tag value值语法错误。

所以在我们项目中引入go vet作为CI检查是很有必要的。
### 标签使用场景
Go官方已经帮忙整理了哪些库已经支持了struct tag：https://github.com/golang/go/wiki/Well-known-struct-tags。
<img width="677" alt="image" src="https://user-images.githubusercontent.com/6757408/155849673-b48c5d92-8d72-46f2-939c-fd9a2b52a835.png">

像json、yaml、gorm、validate、mapstructure、protobuf这几个库的结构体标签是很常用的，gin框架就集成了validate库用来做参数校验，方便了许多，之前写了一篇关于validate的文章：boss: 
[这小子还不会使用validator库进行数据校验开了～～～](https://mp.weixin.qq.com/s?__biz=MzkyNzI1NzM5NQ==&mid=2247484752&idx=1&sn=24a691b9305df828c24b5d9a56f25b46&scene=21#wechat_redirect)，可以关注一下。

具体这些库中是怎么使用的，大家可以看官方文档介绍，写的都很详细，具体场景具体使用哈！！！
### 自定义结构体标签
现在我们可以回答开头的一个问题了，结构体标签是可以随意写的，只要符合语法规则，任意写都可以的，但是一些库没有支持该标签的情况下，随意写的标签是没有任何意义的，如果想要我们的标签变得有意义，
就需要我们提供解析方法。可以通过反射的方式获取标签，所以我们就来看一个例子，如何使用反射获取到自定义的结构体标签。
```
type User struct {
 Name string `asong:"Username"`
 Age  uint16 `asong:"age"`
 Password string `asong:"min=6,max=10"`
}
func getTag(u User) {
 t := reflect.TypeOf(u)

 for i := 0; i < t.NumField(); i++ {
  field := t.Field(i)
  tag := field.Tag.Get("asong")
  fmt.Println("get tag is ", tag)
 }
}

func main()  {
 u := User{
  Name: "asong",
  Age: 5,
  Password: "123456",
 }
 getTag(u)
}
```
运行结果如下：
```
get tag is  Username
get tag is  age
get tag is  min=6,max=10
```
这里我们使用TypeOf方法获取的结构体类型，然后去遍历字段，每个字段StructField都有成员变量Tag：
```
// A StructField describes a single field in a struct.
type StructField struct {
 Name string
 PkgPath string
 Type      Type      // field type
 Tag       StructTag // field tag string
 Offset    uintptr   // offset within struct, in bytes
 Index     []int     // index sequence for Type.FieldByIndex
 Anonymous bool      // is an embedded field
}
```
Tag是一个内置类型，提供了Get、Loopup两种方法来解析标签中的值并返回指定键的值：
```
func (tag StructTag) Get(key string) string
func (tag StructTag) Lookup(key string) (value string, ok bool)
```
Get内部也是调用的Lookup方法。区别在于Lookup会通过返回值告知给定key是否存在与标签中，Get方法完全忽略了这个判断。

### 总结
本文主要介绍一下Go语言中的结构体标签是什么，以及如何使用反射获取到解结构体标签，在日常开发中我们更多的是使用一些库提供好的标签，很少自己开发使用，不过大家有兴趣的话可以读一下validae的源码，
看看他是如何解析结构体中的tag，也可以自己动手实现一个校验库，当作练手项目。

转自：https://mp.weixin.qq.com/s/3sz8oE8nGmba8WECXa0AWg




