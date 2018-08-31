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
### 设置过期时间
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
	_, err = c.Do("expire", "abc", 10)
    	if err != nil {
        	fmt.Println(err)
        	return
    	}
}
```
### list操作
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
	pws := redis.DialPassword("****")

	c, err := redis.Dial("tcp", "localhost:6379", db, pws)
	if err != nil {
		fmt.Println("conn redis failed,", err)
		return
	}
	//  关闭链接
	defer c.Close()

	_, err = c.Do("lpush", "book_list", "abc", "ceg", 300)
	if err != nil {
		fmt.Println(err)
		return
	}

	r, err := redis.String(c.Do("lpop", "book_list"))
	if err != nil {
		fmt.Println("get abc failed,", err)
		return
	}

	fmt.Println(r)
}
```
其他操作都类似。
### Pipelining
* [redigo api](https://godoc.org/github.com/gomodule/redigo/redis)
### Publish and Subscribe 
* [redigo api](https://godoc.org/github.com/gomodule/redigo/redis)
### redis连接池
首先来看一下连接池结构体
```
type Pool struct {
    // Dial是应用程序提供的功能，用于创建和配置连接
    // 从Dial中返回的链接必须不能是处在特殊状态下的，例如下面这些
    // (subscribed to pubsub channel, transaction started, ...).
    
    // Dial is an application supplied function for creating and configuring a
    // connection.
    //
    // The connection returned from Dial must not be in a special state
    // (subscribed to pubsub channel, transaction started, ...).
    Dial func() (Conn, error)

    // TestOnBorrow是一个程序提供的可选的方法，用来检查在链接在被程序使用之前，这个空闲链接是否可用。
    // 参数t是一个时间，代表链接归还给连接池的时间。如果这个方法返回err,标识这个链接已经被关闭了。
    
    // TestOnBorrow is an optional application supplied function for checking
    // the health of an idle connection before the connection is used again by
    // the application. Argument t is the time that the connection was returned
    // to the pool. If the function returns an error, then the connection is
    // closed.
    TestOnBorrow func(c Conn, t time.Time) error

    // 连接池中的最大空闲连接数
    MaxIdle int

    // 连接池在给定时间能够分配的最大连接数，当为0的时候从连接池中获取链接没有连接数限制
    // Maximum number of connections allocated by the pool at a given time.
    // When zero, there is no limit on the number of connections in the pool.
    MaxActive int
    
    // 在保持空闲一段时间后关闭链接。如果值是0，空闲链接不会被关闭。
    // 程序在设置这个值的时候应该小于服务器的超时时间
    
    // Close connections after remaining idle for this duration. If the value
    // is zero, then idle connections are not closed. Applications should set
    // the timeout to a value less than the server's timeout.
    IdleTimeout time.Duration

    // 如果设置为true且连接池的MaxActive达到了这个最大值，Get()方法在返回之前会一直等待从连接池中返回一个链接
    
    // If Wait is true and the pool is at the MaxActive limit, then Get() waits
    // for a connection to be returned to the pool before returning.
    Wait bool

    //  关闭早于此持续时间的链接，如果值为0，连接池不会根据生存期限来关闭链接
    
    // Close connections older than this duration. If the value is zero, then
    // the pool does not close connections based on age.
    MaxConnLifetime time.Duration
    //  包含已过滤或未导出的字段
    // contains filtered or unexported fields
}
```
程序通过Get方法从连接池中获取一个链接，然后通过Close方法把链接归还给连接池。

使用Dial方法配合AUTH命令来做权限验证链接，或者SELECT命令来选择一个数据库：
```
pool := &redis.Pool{
  // Other pool configuration not shown in this example.
  Dial: func () (redis.Conn, error) {
    c, err := redis.Dial("tcp", server)
    if err != nil {
      return nil, err
    }
    if _, err := c.Do("AUTH", password); err != nil {
      c.Close()
      return nil, err
    }
    if _, err := c.Do("SELECT", db); err != nil {
      c.Close()
      return nil, err
    }
    return c, nil
  },
}
```
下面来看一个具体的例子:
```
package main

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
)

var pool *redis.Pool  //  声明一个指针类型的连接池

func init()  {
	pool = &redis.Pool{     //实例化一个连接池
		MaxIdle:16,    //最初的连接数量
		// MaxActive:1000000,    //最大连接数量
		MaxActive:0,    //连接池最大连接数量,不确定可以用0（0表示自动定义），按需分配
		IdleTimeout:300,    //连接关闭时间 300秒 （300秒不使用自动关闭）
		Dial: func() (redis.Conn ,error){     //要连接的redis数据库
			c, err := redis.Dial("tcp","116.62.213.223:6379")
			if err != nil {
				return nil, err
			}
			if _, err := c.Do("AUTH", "TokenClub2018"); err != nil {
				c.Close()
				return nil, err
			}
			if _, err := c.Do("SELECT", 3); err != nil {
				c.Close()
				return nil, err
			}
			return c, nil

		},
	}
}

func main() {
	c := pool.Get() //从连接池，取一个链接
	defer c.Close() //函数运行结束 ，把连接放回连接池

	_,err := c.Do("Set","abc",200)
	if err != nil {
		fmt.Println(err)
		return
	}

	r,err := redis.Int(c.Do("Get","abc"))
	if err != nil {
		fmt.Println("get abc faild :",err)
		return
	}
	fmt.Println(r)
	pool.Close() //关闭连接池
}
```
### redis的主从与集群
redis的主从，集群需要借助其他包实现：
[redigo github地址](https://github.com/gomodule/redigo)中包含了相关的包。

### 注意事项：
* 在程序中链接被使用过后需要手动进行关闭链接
* 连接池中Dial返回的链接
