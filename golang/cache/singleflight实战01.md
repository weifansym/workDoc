singleflight，官方解释其为：singleflight提供了一个重复的函数调用抑制机制。

通俗的解释其作用是，若有多个协程运行某函数时，只让一个协程去处理，然后批量返回。非常适合来做并发控制。常见用于缓存穿透的情况。

缓存穿透即为某个热门内容Key过期，或者突然暴热，请求均没有从cache中获取到数据，就会导致大量的同进程、跨进程的数据回源到存储层，可能会引起存储过载的情况。这个时候使用singleflight就能达到一种
归并回源的效果了。

### 源码解释
普通版本，无归并操作
```
package main

import (
	"errors"
	"log"
	"sync"
)

var errorNotExist = errors.New("redis: key not found")

func main() {
	var wg sync.WaitGroup
	wg.Add(10)

	// 开启10个协程
	for i := 0; i < 10; i++ {
		go func() {
			defer wg.Done()
			data, err := getData("2000")
			if err != nil {
				log.Print(err)
				return
			}
			log.Println(data)
		}()
	}
	wg.Wait()
}

// 获取数据
func getData(key string) (string, error) {
	data, err := getDataFromCache(key)
	if err == errorNotExist {
		// 穿透到 DB 捞取数据
		data, err = getDataFromDB(key)
		if err != nil {
			log.Println(err)
			return "", err
		}

		// 回填数据到cache, 此处为模拟请求数差不多时间到达，还来不及回填cache
	} else if err != nil {
		return "", err
	}
	return data, nil
}

// 从cache中获取值，cache中无该值
func getDataFromCache(key string) (string, error) {
	return "", errorNotExist
}

// 从数据库中获取值
func getDataFromDB(key string) (string, error) {
	log.Printf("get %s from database", key)
	return "2000 in db", nil
}
```
查看打印值：
```
2021/09/29 23:45:27 get 2000 from database
2021/09/29 23:45:27 2000 in db
2021/09/29 23:45:27 get 2000 from database
2021/09/29 23:45:27 2000 in db
2021/09/29 23:45:27 get 2000 from database
2021/09/29 23:45:27 2000 in db
2021/09/29 23:45:27 get 2000 from database
2021/09/29 23:45:27 2000 in db
2021/09/29 23:45:27 get 2000 from database
2021/09/29 23:45:27 2000 in db
2021/09/29 23:45:27 get 2000 from database
2021/09/29 23:45:27 2000 in db
2021/09/29 23:45:27 get 2000 from database
2021/09/29 23:45:27 2000 in db
2021/09/29 23:45:27 get 2000 from database
2021/09/29 23:45:27 2000 in db
2021/09/29 23:45:27 get 2000 from database
2021/09/29 23:45:27 2000 in db
2021/09/29 23:45:27 get 2000 from database
2021/09/29 23:45:27 2000 in db
```
跟常见的cache-aside一个逻辑：1、先cache中拿取数据；2、取不到数据就从DB拿，再回填给cache；

但当请求来的比较快，cache是来不及回填的，也就会出现上述打印的现象，请求都从DB拿的数据。这个时候再来看看singleflight的官方介绍：
> singleflight提供了一个重复的函数调用抑制机制。

引入singleflight修改源码：
```
import "golang.org/x/sync/singleflight"

var singleFlightTest singleflight.Group

.......

//获取数据
func getDataBySingleFlight(key string) (string, error) {
	data, err := getDataFromCache(key)
	if err == errorNotExist {
		// 从db中获取数据
		v, err, _ := singleFlightTest.Do(key, func() (interface{}, error) {
			return getDataFromDB(key)
			// set cache
		})
		if err != nil {
			log.Println(err)
			return "", err
		}

		data = v.(string)
        
        // set cache
        // 可以在sl的wrapFunc中回填cache，也可以在外面回填数据，前者在go-zero的实例代码中可见，后者出现在bilibili的代码中
        // 读取的数据DB不存在，应该放置一个TTL标志位
	} else if err != nil {
		return "", err
	}
	return data, nil
}
```
结果的打印如下，可以看到一个请求去捞取完数据后，其他请求也都拿到数据了。
```
2021/09/30 00:14:23 get 2000 from database
2021/09/30 00:14:23 2000 in db
2021/09/30 00:14:23 2000 in db
2021/09/30 00:14:23 2000 in db
2021/09/30 00:14:23 2000 in db
2021/09/30 00:14:23 2000 in db
2021/09/30 00:14:23 2000 in db
2021/09/30 00:14:23 2000 in db
2021/09/30 00:14:23 2000 in db
2021/09/30 00:14:23 2000 in db
2021/09/30 00:14:23 2000 in db
```
包比较简单，网上也有很多分析代码的博客，我只在这里看下Do函数：
```
func (g *Group) Do(key string, fn func() (interface{}, error)) (v interface{}, err error, shared bool) {
	g.mu.Lock()
	// 初始化进入此单飞集合的集合
	if g.m == nil {
		g.m = make(map[string]*call)
	}
    // 相同的处理，做下累加操作，同时卡死在这里等待第一次执行的请求执行完，返回其数据
	if c, ok := g.m[key]; ok {
		c.dups++
		g.mu.Unlock()
		c.wg.Wait()
		return c.val, c.err, true
	}
    
    // 第一次过来的请求，做初始化处理
	c := new(call)
	c.wg.Add(1)
	g.m[key] = c
	g.mu.Unlock()

    // 执行请求
	g.doCall(c, key, fn)
	return c.val, c.err, c.dups > 0
}
```
可以通过分析源码看到，相同的key，只有第一次进入的请求，才会执行；后面进入的请求，都会被waitgroup给锁住，原地等待第一次执行的请求执行完成，然后统一返回请求拿到的数据以及err信息。

### 提到的资料
[go-zero缓存设计之持久层缓存](https://www.bookstack.cn/read/go-zero-1.3-zh/redis-cache.md)
https://pkg.go.dev/golang.org/x/sync/singleflight



