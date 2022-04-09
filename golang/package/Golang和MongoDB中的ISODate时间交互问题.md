## Golang 和 MongoDB 中的 ISODate 时间交互问题
MongoDB 中有一种时间格式数据 ISODate，参考如下：
![image](https://user-images.githubusercontent.com/6757408/162582817-fe82019a-e270-4d10-a3c0-2ea7c45e9fac.png)
如果在 Golang 中查询这条记录，Golang用什么类型的变量来保存呢？

### 查找 ISODate 时间字段
在 Golang 中可以使用 time.Time 数据类型来保存 MongoDB 中的 ISODate 时间。
```
type Model struct {
    Id   bson.ObjectId `bson:"_id,omitempty"`
    Time time.Time     `bson:"time"`
}
m := Model{}
err := c.Find(bson.M{"_id": bson.ObjectIdHex("572f3c68e43001d2c1703aa7")}).One(&m)
if err != nil {
    panic(err)
}
fmt.Printf("%+v\n", m)
// output: {Id:ObjectIdHex("572f3c68e43001d2c1703aa7") Time:2015-07-08 17:29:14.002 +0800 CST}
```
从输出中可以看到 Golang 输出的时间格式是 CST 时区，Golang 在处理的过程中将 ISO 时间转换成了 CST 时间，从时间面板上来看也比 MongoDB 中的快上 8 个小时，这个是正常的。

那么 Golang 做插入操作和或者时间比较操作的时候需要自己转换时间戳吗？答案是不需要的，来看下插入的例子。

### 插入时间

重新插入一条记录，记录的Time字段为当前时间，在golang中可以通过time.Now获取当前时间，查看输出可以看到是CST的时间格式
```
now := time.Now()
fmt.Printf("%+v\n", now)
// output: 2016-05-12 14:34:00.998011694 +0800 CST
err = c.Insert(Model{Time: now})
if err != nil {
    panic(err)
}
```
### 查看 MongoDB 中的记录
插入当前时间到 MongoDB:
![image](https://user-images.githubusercontent.com/6757408/162582865-9a8aa395-40a6-47c3-8690-8f72a69f0020.png)
可以看到存储到 MongoDB 中时间的自动转换为了 ISO 时间，时间少了 8 个小时。小结一下就是 Golang 和 MongoDB 中的时间交互不需要考虑额外的东西，因为驱动都进行了转换。
### 时间字符串转成 time.Time

有时我们会将 time.Time 的时间以字符串的形式存储，那么要和 MongoDB 交互的时候就需要转换 time.Time 格式
```
// 时间字符串转到到time.Time格式
// 使用time.Parse方法进行转换
timeString := "2016-05-12 14:34:00.998011694 +0800 CST"
t, err := time.Parse("2006-01-02 15:04:05.999999999 -0700 MST", timeString)
if err != nil {
    panic(err)
}
fmt.Printf("%+v\n", t)
```
代码中比较难理解的就是 time.Parse 的第一个参数，这个其实是 Golang 当中的定义，详细看下 time.Time.String() 的源码就会明白了。
