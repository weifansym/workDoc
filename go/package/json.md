## json序列化，反序列化
具体序列化与反序列化的方法，请参见：官网文档：[json](https://golang.org/pkg/encoding/json/)
### 结构体转json
实例如下：
```
package main

import (
	"encoding/json"
	"fmt"
)

type User struct {
	UserName string `json:"username"`
	NickName string `json:"nickname"`
	Age      int
	Birthday string
	Sex      string
	Email    string
	Phone    string
}

/*结构体转json*/

func testStruct() {
	user1 := &User{
		UserName: "user1",
		NickName: "Murphy",
		Age:      18,
		Birthday: "2008/8/8",
		Sex:      "男",
		Email:    "123456@qq.com",
		Phone:    "15600000000",
	}

	data, err := json.Marshal(user1)
	if err != nil {
		fmt.Printf("json.marshal failed, err:", err)
		return
	}

	fmt.Printf("%s\n", string(data))
}

func main() {
	testStruct()
	fmt.Println("----")
}
```
输出结果如下：
```
{"username":"user1","nickname":"Murphy","Age":18,"Birthday":"2008/8/8","Sex":"男","Email":"123456@qq.com","Phone":"15600000000"}
----
```
### map转json
示例如下
```
package main

import (
	"encoding/json"
	"fmt"
)

type User struct {
	UserName string `json:"username"`
	NickName string `json:"nickname"`
	Age      int
	Birthday string
	Sex      string
	Email    string
	Phone    string
}

/*map转json*/

func testMap() {
	var mmp map[string]interface{}
	mmp = make(map[string]interface{})

	mmp["username"] = "Murphy"
	mmp["age"] = 19
	mmp["sex"] = "man"

	data, err := json.Marshal(mmp)
	if err != nil {
		fmt.Println("json marshal failed,err:", err)
		return
	}
	fmt.Printf("%s\n", string(data))

}

func main() {
	testMap()
	fmt.Println("----")
}
```
输出结果
```
{"age":19,"sex":"man","username":"Murphy"}
----
```
其他类型序列化类似，例如：slice序列化
```
package main

import (
	"encoding/json"
	"fmt"
)

func testSlice() {
	var m map[string]interface{}
	var s []map[string]interface{}
	m = make(map[string]interface{})
	m["username"] = "user1"
	m["age"] = 18
	m["sex"] = "man"

	s = append(s, m)

	m = make(map[string]interface{})
	m["username"] = "user2"
	m["age"] = 29
	m["sex"] = "female"
	s = append(s, m)

	data, err := json.Marshal(s)
	if err != nil {
		fmt.Printf("json.marshal failed, err:", err)
		return
	}

	fmt.Printf("%s\n", string(data))
}

func main() {
	testSlice()
	fmt.Println("--------")
}
```
下面来看一下反序列化的例子。
###json反序列化为结构体
实例如下：
```
package main

import (
	"encoding/json"
	"fmt"
)

type User struct {
	UserName string `json:"username"`
	NickName string `json:"nickname"`
	Age      int
	Birthday string
	Sex      string
	Email    string
	Phone    string
}

func testStruct() (ret string, err error) {
	user1 := &User{
		UserName: "user1",
		NickName: "Murphy",
		Age:      18,
		Birthday: "2008/8/8",
		Sex:      "男",
		Email:    "taidou008@qq.com",
		Phone:    "110",
	}

	data, err := json.Marshal(user1)
	if err != nil {
		err = fmt.Errorf("json.marshal failed, err:", err)
		return
	}

	ret = string(data)
	return
}

func test() {
	data, err := testStruct()
	if err != nil {
		fmt.Println("test struct failed, ", err)
		return
	}

	var user1 User
	err = json.Unmarshal([]byte(data), &user1)
	if err != nil {
		fmt.Println("Unmarshal failed, ", err)
		return
	}
	fmt.Println(user1)
}

func main() {
	test()
	fmt.Println("-----")
}
```
输出结果：
```
{user1 Murphy 18 2008/8/8 男 taidou008@qq.com 110}
-----
```
### json反序列化为map
```
package main

import (
	"encoding/json"
	"fmt"
	"reflect"
)

func testMap() (ret string, err error) {
	var m map[string]interface{}
	m = make(map[string]interface{})
	m["username"] = "user1"
	m["age"] = 18
	m["sex"] = "man"

	data, err := json.Marshal(m)
	if err != nil {
		err = fmt.Errorf("json.marshal failed, err:", err)
		return
	}

	ret = string(data)
	return
}

func test2() {
	data, err := testMap()
	if err != nil {
		fmt.Println("test map failed, ", err)
		return
	}

	var m map[string]interface{}
	err = json.Unmarshal([]byte(data), &m)
	if err != nil {
		fmt.Println("Unmarshal failed, ", err)
		return
	}
	fmt.Println(m)
	fmt.Println(reflect.TypeOf(m))
}

func main() {
	test2()
}
```
结果如下：
```
map[sex:man username:user1 age:18]
map[string]interface {}
```
