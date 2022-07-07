## Go语言操作mongoDB
我们这里使用的是官方的驱动包，当然你也可以使用第三方的驱动包（如mgo等）。 mongoDB官方版的Go驱动发布的比较晚（2018年12月13号）。
### 安装mongoDB Go驱动包
```
go get github.com/mongodb/mongo-go-driver
```
### 通过Go代码连接mongoDB
```
package main

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// 设置客户端连接配置
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	// 连接到MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// 检查连接
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to MongoDB!")
}
```
连接上MongoDB之后，可以通过下面的语句处理我们上面的q1mi数据库中的student数据集了：
```
// 指定获取要操作的数据集
collection := client.Database("q1mi").Collection("student")
```
处理完任务之后可以通过下面的命令断开与MongoDB的连接：
```
// 断开连接
err = client.Disconnect(context.TODO())
if err != nil {
	log.Fatal(err)
}
fmt.Println("Connection to MongoDB closed.")
```
### 连接池模式
```
import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectToDB(uri, name string, timeout time.Duration, num uint64) (*mongo.Database, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	o := options.Client().ApplyURI(uri)
	o.SetMaxPoolSize(num)
	client, err := mongo.Connect(ctx, o)
	if err != nil {
		return nil, err
	}

	return client.Database(name), nil
```
### BSON
MongoDB中的JSON文档存储在名为BSON(二进制编码的JSON)的二进制表示中。与其他将JSON数据存储为简单字符串和数字的数据库不同，BSON编码扩展了JSON表示，使其包含额外的类型，如int、long、date、浮点数和decimal128。这使得应用程序更容易可靠地处理、排序和比较数据。

连接MongoDB的Go驱动程序中有两大类型表示BSON数据：D和Raw。

类型D家族被用来简洁地构建使用本地Go类型的BSON对象。这对于构造传递给MongoDB的命令特别有用。D家族包括四类:
* D：一个BSON文档。这种类型应该在顺序重要的情况下使用，比如MongoDB命令。
* M：一张无序的map。它和D是一样的，只是它不保持顺序。
* A：一个BSON数组。
* E：D里面的一个元素。

要使用BSON，需要先导入下面的包：
```
import "go.mongodb.org/mongo-driver/bson"
```
下面是一个使用D类型构建的过滤器文档的例子，它可以用来查找name字段与’张三’或’李四’匹配的文档:
```
bson.D{{
	"name",
	bson.D{{
		"$in",
		bson.A{"张三", "李四"},
	}},
}}
```
Raw类型家族用于验证字节切片。你还可以使用Lookup()从原始类型检索单个元素。如果你不想要将BSON反序列化成另一种类型的开销，那么这是非常有用的。这个教程我们将只使用D类型。

### CRUD

转自：https://www.liwenzhou.com/posts/Go/go_mongodb/

mongodb官方链接配置：https://www.mongodb.com/docs/drivers/go/current/fundamentals/connection/



