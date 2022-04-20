## 理解 Go 中的 JSON
> 本文是基于 Go 官方和 https://eager.io/blog/go-and-json/ 进行翻译整理的
JSON 是一种轻量级的数据交换格式，常用作前后端数据交换，Go 在 encoding/json 包中提供了对 JSON 的支持。

### 序列化
把 Go struct 序列化成 JSON 对象，Go 提供了 Marshal 方法，正如其含义所示表示编排序列化，函数签名如下：
```
func Marshal(v interface{}) ([]byte, error)
```
举例来说，比如下面的 Go struct：
```
type Message struct {
    Name string
    Body string
    Time int64
}
```
使用 Marshal 序列化：
```
m := Message{"Alice", "Hello", 1294706395881547000}
b, err := json.Marshal(m) 
fmt.Println(b) //{"Name":"Alice","Body":"Hello","Time":1294706395881547000}
```
在 Go 中并不是所有的类型都能进行序列化：
* JSON object key 只支持 string
* Channel、complex、function 等 type 无法进行序列化
* 数据中如果存在循环引用，则不能进行序列化，因为序列化时会进行递归
* Pointer 序列化之后是其指向的值或者是 nil

> 还需要注意的是：**只有 struct 中支持导出的 field 才能被 JSON package 序列化，即首字母大写的 field。**
### 反序列化
反序列化函数是 Unmarshal ，其函数签名如下：
```
func Unmarshal(data []byte, v interface{}) error
```
如果要进行反序列化，我们首先需要创建一个可以接受序列化数据的 Go struct：
```
var m Message
err := json.Unmarshal(b, &m)
```
JSON 对象一般都是小写表示，Marshal 之后 JSON 对象的首字母依然是大写，如果序列化之后名称想要改变如何实现，答案就是**struct tags**。
### Struct Tag
Struct tag 可以决定 Marshal 和 Unmarshal 函数如何序列化和反序列化数据。
#### 指定 JSON filed name
JSON object 中的 name 一般都是小写，我们可以通过 struct tag 来实现：
```
type MyStruct struct {
    SomeField string `json:"some_field"`
}
```
SomeField 序列化之后会变成 some_field。
#### 指定 field 是 empty 时的行为
使用**omitempty**可以告诉 Marshal 函数如果 field 的值是对应类型的 zero-value，那么序列化之后的 JSON object 中不包含此 field：
```
type MyStruct struct {
    SomeField string `json:"some_field,omitempty"`
}

m := MyStruct{}
b, err := json.Marshal(m) //{}
```
如果**SomeField == “”**，序列化之后的对象就是 {}。
#### 跳过 field

Struct tag “-” 表示跳过指定的 filed：
```
type MyStruct struct {
    SomeField string `json:"some_field"`
    Passwd string `json:"-"`
}
m := MyStruct{}
b, err := json.Marshal(m) //{"some_feild":""}
```
即序列化的时候不输出，这样可以有效保护需要保护的字段不被序列化。

#### 反序列化任意 JSON 数据
默认的 JSON 只支持以下几种 Go 类型：
* bool for JSON booleans
* float64 for JSON numbers
* string for JSON strings
* nil for JSON null
在序列化之前如果不知道 JSON 数据格式，我们使用 interface{} 来存储。interface {} 的作用详见本博的其他文章。

有如下的数据格式：
```
b := []byte(`{"Name":"Wednesday","Age":6,"Parents":["Gomez","Morticia"]}`)
```
如果我们序列化之前不知道其数据格式，我们可以使用 interface{} 来存储我们的 decode 之后的数据：
```
var f interface{}
err := json.Unmarshal(b, &f)
```
反序列化之后 f 应该是像下面这样：
```
f = map[string]interface{}{
    "Name": "Wednesday",
    "Age":  6,
    "Parents": []interface{}{
        "Gomez",
        "Morticia",
    },
}
```
key 是 string，value 是存储在 interface{} 内的。想要获得 f 中的数据，我们首先需要进行 type assertion，然后通过 range 迭代获得 f 中所有的 key ：
```
m := f.(map[string]interface{})
for k, v := range m {
    switch vv := v.(type) {
    case string:
        fmt.Println(k, "is string", vv)
    case float64:
        fmt.Println(k, "is float64", vv)
    case []interface{}:
        fmt.Println(k, "is an array:")
        for i, u := range vv {
            fmt.Println(i, u)
        }
    default:
        fmt.Println(k, "is of a type I don't know how to handle")
    }
}
```
#### 反序列化对 slice、map、pointer 的处理
我们定义一个 struct 继续对上面例子中的 b 进行反序列化：
```
type FamilyMember struct {
    Name    string
    Age     int
    Parents []string
}

var m FamilyMember
err := json.Unmarshal(b, &m)
```
这个例子是能够正常工作的，你一定也注意到了，struct 中包含一个 slice Parents ，slice 默认是 nil，之所以反序列化可以正常进行就是因为 Unmarshal 在序列化时进行了对 slice Parents 
做了初始化，同理，对 map 和 pointer 都会做类似的工作，比如序列化如果 Pointer 不是 nil 首先进行 dereference 获得其指向的值，然后再进行序列化，反序列化时首先对 nil pointer 
进行初始化

#### Stream JSON
除了 marshal 和 unmarshal 函数，Go 还提供了 Decoder 和 Encoder 对 stream JSON 进行处理，常见 request 中的 Body、文件等：
```
jsonFile, err := os.Open("post.json")
if err != nil {
    fmt.Println("Error opening json file:", err)
    return
}

defer jsonFile.Close()
decoder := json.NewDecoder(jsonFile)
for {
    var post Post
    err := decoder.Decode(&post)
    if err == io.EOF {
        break
    }

    if err != nil {
        fmt.Println("error decoding json:", err)
        return
    }

    fmt.Println(post)
}
```
#### 嵌入式 struct 的序列化
Go 支持对 nested struct 进行序列化和反序列化:
```
type App struct {
	Id string `json:"id"`
}

type Org struct {
	Name string `json:"name"`
}

type AppWithOrg struct {
	App
	Org
}

func main() {
	data := []byte(`
        {
            "id": "k34rAT4",
            "name": "My Awesome Org"
        }
    `)

	var b AppWithOrg

	json.Unmarshal(data, &b)
	fmt.Printf("%#v", b)

	a := AppWithOrg{
		App: App{
			Id: "k34rAT4",
		},
		Org: Org{
			Name: "My Awesome Org",
		},
	}
	data, _ = json.Marshal(a)
	fmt.Println(string(data))
}
```
Nested struct 虽然看起来有点怪异，有些时候它将非常有用。

#### 自定义序列化函数
Go JSON package 中定了两个 Interface Marshaler 和 Unmarshaler ，实现这两个 Interface 可以让你定义的 type 支持序列化操作。

### 错误处理
总是记得检查序列或反序列化的错误，可以让你的程序更健壮，而不是在出错之后带着错误继续执行下去。
### 参考资料
https://blog.golang.org/json-and-go
https://eager.io/blog/go-and-json/

