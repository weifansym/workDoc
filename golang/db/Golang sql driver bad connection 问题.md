## 背景
公司核心业务库之前想接入ucloud读写分离中间件，其过程中对此核心库做过压测，压测期间使用线上的配置，但当并发过高时，会出现bad connection的报错，sql请求会反馈Invalid connection，
当时比较忙就简单查了下，现在得空就深入探究了下原因。
![image](https://user-images.githubusercontent.com/6757408/205683014-0e7c9062-9e1a-44ce-9774-ef41101a7119.png)

### 诊断
看报错，是go-sql-driver反馈的无效链接，查明为mysql服务端断开链接后，客户端没有断开，导致链接池中产生无效链接，sql请求从链接池中捞出了无效链接来执行，触发此报错。

### 问题复现
环境：docker-mysql-5.7 、 go-sql-driver-1.4.1  、go-1.13

思路：初始化mysql后，sleep 60s , 之后执行select 操作，sleep过程中，重启mysql, 模拟mysql主动关闭链接情景，观察select反馈信息

过程：
```
package main
 
 
import (
    "database/sql"
    "fmt"
    "log"
    "time"
 
    _ "github.com/go-sql-driver/mysql"
)
 
var (
    dsn string = "root:root@tcp(127.0.0.1:3307)/demo"
    db  *sql.DB
    err error
)
 
type user struct {
    Id   int
    Name string
    Age  int
}
 
//init 初始化
func init() {
    db, err = sql.Open("mysql", dsn)
    failOnError(err, "fail to open sql")
    err = db.Ping()
    failOnError(err, "fail to ping database")
    fmt.Println("connect mysql success!")
}
 
//failOnError 打印错误信息
func failOnError(err error, msg string) {
    if err != nil {
        log.Fatalf("%s: %s", msg, err)
    }
}
 
func main() {
    defer db.Close()
    fmt.Println("开始睡眠...")
    time.Sleep(time.Second * 30)
    fmt.Println("睡眠结束...")
    //查库
    sql := "SELECT * FROM users"
    rows, err := db.Query(sql)
    failOnError(err, "fail to query")
    for rows.Next() {
        u := user{}
        err := rows.Scan(&u.Id, &u.Name, &u.Age)
        if err != nil {
            fmt.Printf("批量查询scan失败 :%+v\n", err)
            return
        }
        fmt.Printf("批量查询成功 :%+v\n", u)
    }
}
```
运行应用后，重启mysql, 应用输出如下
```
connect mysql success!
开始睡眠...
睡眠结束...
[mysql] 2021/02/01 16:35:10 packets.go:36: unexpected EOF
2021/02/01 16:35:10 fail to query: invalid connection
exit status 1
```
我们成功复现出问题，invalid connection，由此我们可以确认，出现bad connection的原因之一可能就是mysql服务端主动或因异常中断了链接，客户端没有中断，导致连接池中产生了无效链接，
sql请求捞起了此无效链接执行sql时产生的报错。

### 问题
1.为什么mysql服务端会关闭链接？原则上服务端是不允许随意关闭链接的。
2.我们如何避免，当服务端关闭链接、客户端未关闭链接时产生的无效链接对应用的影响？

## 探究
1. 为什么mysql服务端会关闭链接？
通过查询相关资料得知，mysql默认配置里，可以配置长连接的时间，默认为8h，也就是说，出生的链接就会被计时，到了8H服务端就会主动关闭。
* wait_timeout：服务器关闭非交互连接之前等待活动的秒数。
* interactive_timeout：服务器关闭交互连接之前等待活动的秒数

还有一种情况就是mysql服务端异常时，会直接触发链接异常中断，当服务端CPU,IO, 内存等状态有明显跑高时，就会可能会产生链接异常中断。

2.我们如何避免这种因各种原因产生的无效链接给业务带来的影响呢？
设置SetConnMaxLifetime
方法之一是调整客户端的配置, SetConnMaxLifetime, 它设置了连接可重用的最大时间长度, 原理就是在生成链接的时候给链接加了个定时器，每秒会检查一次链接的有效性，超过设置的时间，
会把此链接从链接池中删除，在YouhuiDataDb压测过程中，我调整了此参数为2s(公司框架中默认不设置，即默认链接永久有效)，有效避免了bad connection的产生，但是也付出了代价，就是会频繁
的从0开始创建时间，一定程度上丢失了连接池的优势，此外还会给msql增加一定的set_options命令，根据dba所说，此命令也会一定程度上影响db性能。

要注意的有以下几点:
* 这并不能保证连接将在池中存在完整的2s，很可能由于某种原因连接将变得不可用，并且在此之前自动关闭。
* 一个连接在创建后仍可以使用2s，只是说2s后不能再被重用了。
* 这不是空闲超时。连接将在第一次创建后2s后过期，而不是2s后变成空闲。
* 每秒自动运行一次清理操作以便从池中删除“过期”连接。

### go-sql-driver升级1.5
目前zdm-go-base中依旧使用的1.4.1，但官方已经在1.5版本修复了这个bug
![image](https://user-images.githubusercontent.com/6757408/205683848-1acb896a-9649-49e8-b270-e23cbed90235.png)

实际上就是，在遇到bad connection时 ，1.5版本的驱动会标记这个链接并且重新创建一个链接来使用，我们来操作一下看看

环境：docker-mysql-5.7 、 go-sql-driver-1.5  、go-1.13

思路：初始化mysql后，sleep 60s , 之后执行select 操作，sleep过程中，重启mysql, 模拟mysql主动关闭链接情景，观察select反馈信息

过程：应用代码如上问题复现代码，我们只是升级了go-sql-driver,  我们直接来看结果，并与上述问题浮现时的结果做对比

![image](https://user-images.githubusercontent.com/6757408/205683943-3206b7cb-a74c-4b83-b8d4-5a7148019621.png)

可以看出，1.5版本的驱动，会在遇到无效链接时自动给关闭并且使用一个有效的链接去继续执行sql， 此时对业务就无影响，不会因无效链接而不响应或者返回异常。

### 总结
1. bad connection 产生的原因就是使用了无效链接，产生无效链接的场景可能多种多样，但几率很低，并且大多都是服务资源上的问题。
2. 解决方案里，个人推荐第二种，因为官方已经更新过了，比较香，但是需要我们还是有必要测试环境下看看是否有坑，新版本是否稳定。

转自：https://www.jiangailang.cn/151.html

https://www.jiangailang.cn/281.html




