### MongoDB的模糊查询
模糊查询时数据库应用中不可缺少的一步，MySQL中使用like和或者regexp来实现实现模糊查询，而MongoDB则使用$regex操作符或直接使用正则表达式对象来实现。
| MySQL  | Mongodb |
| ------------- | ------------- |
| select * from users where name like ’%InkDP%’  | db.users.find({name: {$regex: /InkDP/}})  |
| select * from users where name regexp ’InkDP’ | db.users.find({name: /InkDP/})  |
更多相关的语法可查看官方文档：[$regex](https://www.mongodb.com/docs/manual/reference/operator/query/regex/)，就不再做多讨论。
### 使用MongoDB GO Driver进行查询
先来看看我们的数据源：
```
 db.users.find({})
{ "_id" : ObjectId("600704fffc9b483f284d0bc3"), "name" : "1InkDP" }
{ "_id" : ObjectId("600704fffc9b483f284d0bc4"), "name" : "InkDPPP" }
{ "_id" : ObjectId("600704fffc9b483f284d0bc5"), "name" : "InkDP" }
{ "_id" : ObjectId("600704fffc9b483f284d0bc6"), "name" : "inkdp123" }
{ "_id" : ObjectId("600704fffc9b483f284d0bc7"), "name" : "abcdef" }
{ "_id" : ObjectId("60070500fc9b483f284d0bc8"), "name" : "test" }
```
然后执行模糊查询：
```
db.users.find({name:{$regex: /InkDP/,$options: "i"}})
{ "_id" : ObjectId("600704fffc9b483f284d0bc3"), "name" : "1InkDP" }
{ "_id" : ObjectId("600704fffc9b483f284d0bc4"), "name" : "InkDPPP" }
{ "_id" : ObjectId("600704fffc9b483f284d0bc5"), "name" : "InkDP" }
{ "_id" : ObjectId("600704fffc9b483f284d0bc6"), "name" : "inkdp123" }

```
错误尝试
上述方式是MongoDB的命令行的执行方式，如果我们直接在Go里面直接这样写是行不通的
```
filter := bson.M{
   "name": bson.M{
      "$regex":   "/InkDP/",
      "$options": "i",
   },
}
```
当你兴高采烈地拿着上面的查询条件去查询时，你会发现它会返回一个空数组给你
#### 正确的使用方式
```
filter := bson.M{
	"name": primitive.Regex{
		Pattern:"/InkDP/",
		Options: "i",
	},
}
```
执行后发现还是没有，一番查找后才发现Pattern不再额外需要两个/，直接填写正则内容即可，所以我们改为：
```
filter := bson.M{
   "name": primitive.Regex{
      Pattern:"InkDP",
      Options: "i",
   },
}
```
执行结果为：
```
{ID:ObjectID("600704fffc9b483f284d0bc3") Name:1InkDP}
{ID:ObjectID("600704fffc9b483f284d0bc4") Name:InkDPPP}
{ID:ObjectID("600704fffc9b483f284d0bc5") Name:InkDP}
{ID:ObjectID("600704fffc9b483f284d0bc6") Name:inkdp123}
```
与命令行查找的一致，说明没有问题

转自：https://www.inkdp.cn/skill/back-end/61018.html#MongoDB%E7%9A%84%E6%A8%A1%E7%B3%8A%E6%9F%A5%E8%AF%A2
