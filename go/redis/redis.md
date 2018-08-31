## Golang redis client
从redis官网我们可以看到golang可用的Redis客户端：[redis client](https://redis.io/clients#go),下面我们选择redigo作为我们的客户端。
* [redigo github地址](https://github.com/gomodule/redigo)
* [redigo api](https://godoc.org/github.com/gomodule/redigo/redis)

### redigo
在项目中引用
```
import "github.com/gomodule/redigo/redis"
```
#### 链接redis
```
package main

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
)

func main() {
	c, err := redis.Dial("tcp", "116.62.213.223:6379")
	if err != nil {
		fmt.Println("conn redis failed,", err)
		return
	}

	fmt.Println("redis conn success")

	defer c.Close()
}
```
### string类型操作
string类型的简单set,get:
```
package main

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
)

func main() {
	//  设置选择链接的db
	db := redis.DialDatabase(3)
	// 设置密码
	pws := redis.DialPassword("******")

	c, err := redis.Dial("tcp", "localhost:6379", db, pws)
	if err != nil {
		fmt.Println("conn redis failed,", err)
		return
	}
	//  关闭链接
	defer c.Close()

	_, err = c.Do("Set", "abc", 100)
	if err != nil {
		fmt.Println(err)
		return
	}

	r, err := redis.Int(c.Do("Get", "abc"))
	if err != nil {
		fmt.Println("get abc failed,", err)
		return
	}

	fmt.Println(r)
}
```
输出：
```
100
```
简单的批量操作
```
package main

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
)

func main() {
	//  设置选择链接的db
	db := redis.DialDatabase(3)
	// 设置密码
	pws := redis.DialPassword("*****")

	c, err := redis.Dial("tcp", "localhost:6379", db, pws)
	if err != nil {
		fmt.Println("conn redis failed,", err)
		return
	}
	//  关闭链接
	defer c.Close()

	_, err = c.Do("MSet", "abc", 100, "efg", 300)
	if err != nil {
		fmt.Println(err)
		return
	}

	r, err := redis.Ints(c.Do("MGet", "abc", "efg"))
	if err != nil {
		fmt.Println("get abc failed,", err)
		return
	}

	for _, v := range r {
		fmt.Println(v)
	}
}
```
注意事项：
* 在程序中链接被使用过后需要手动进行关闭链接
* 
