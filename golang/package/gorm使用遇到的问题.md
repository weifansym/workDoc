## gorm使用遇到的问题
### gorm解析go time.Time数据
```
Scan error on column index 3: unsupported Scan, storing driver.Value type []uint8 into type *time.Time
```
使用gorm框架，数据库使用的mysql 直接上解决办法 最后加上这个即可
```
parseTime=true

db,err:=gorm.Open("mysql","用户名:密码@tcp(localhost:3306)/数据库名?charset=utf8&parseTime=true")
```
### gorm自动更新创建时间及更新时间，自动更新时间戳
平时写代码，总是要处理更新时间和创建时间，要写不少的代码，而且还容易忘记。
针对于这个问题研究了一下有没有什么比较好的方式。下面说一下如何摆脱体力劳动。
假设场景，需要改分好毕业。。。
```
CREATE TABLE `user` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(32) DEFAULT NULL,
  `score` int(11) DEFAULT NULL,
  `createtime` datetime DEFAULT CURRENT_TIMESTAMP,
  `updatetime` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`userid`)
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8
```
此时创建数据库的语句中要加入 createtime 跟 updatetime 两个，一个是创建的时候填上，一个是创建的时候填上并且更新的时候自动更新。
这样就配置好数据库部分了。
```
type TableUser struct {
  UserID int `gorm:"column:id" json:"id"`
  Name string `gorm:"column:name" json:"name"`
  Score int `gorm:"column:score" json:"score"`
  CreateTime time.Time `gorm:"column:createtime;default:null" json:"createtime"`
  UpdateTime time.Time `gorm:"column:updatetime;default:null" json:"updatetime"`
}
```
这样gorm需要的结构体也创建好了，最重要的是其中的**default:null**，这样以后在创建或者更新的时候都不需要传递CreuateTime跟UpdateTime两个了。
```
var u = TableUser {
  Name: "小明",
  District: 59,
}u
db.Create(&u)
```
Nice~
#### 留意不合法的时间值

MySQL的DATE/DATATIME类型可以对应Golang的time.Time。但是，如果DATE/DATATIME不慎插入了一个无效值，例如2016-00-00 00:00:00, 那么这条记录是无法查询出来的。
会返回gorm.RecordNotFound类型错误。零值0000-00-00 00:00:00是有效值，不影响正常查询。

参见：
https://www.cnblogs.com/mrylong/p/11326792.html
http://www.iwtt.xyz/articles/2020/05/17/1589725477701.html





