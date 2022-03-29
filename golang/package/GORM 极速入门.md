## GORM 极速入门
### 一、基础概念
ORM(Object Relational Mapping)，意思是对象关系映射。

数据库会提供官方客户端驱动，但是需要自己处理 SQL 和结构体的转换。

使用 ORM 框架让我们避免转换，写出一些无聊的冗余代码。理论上 ORM 框架可以让我们脱离 SQL，但实际上还是需要懂 SQL 才可以使用 ORM。

我本人是比较排斥使用 ORM 框架的，原因有两点。

一、不自由，我不能随心所欲的控制我的数据库。

二、性能差，比官方客户端驱动直接编写 SQL 的效率低 3-5 倍。

不过 ORM 也有很多优点，它可以在一定程度上让新手避免慢 SQL。

也有一些文章讨论过 ORM 的利弊。比如这篇：[orm_is_an_antipattern](https://seldo.com/posts/orm_is_an_antipattern)。

总的来说，是否使用 ORM 框架取决于一个项目的开发人员组织结构。

老手渴望自由，新手需要规则。世界上新手多，老手就要做出一些迁就。

gorm 是一款用 Golang 开发的 orm 框架，目前已经成为在 Golang Web 开发中最流行的 orm 框架之一。本文将对 gorm 中常用的 API 进行讲解，帮助你快速学会 gorm。

除了 gorm，你还有其他选择，比如 sqlx 和 sqlc。
### 二、连接 MySQL
gorm 可以连接多种数据库，只需要不同的驱动即可。官方目前仅支持 MySQL、PostgreSQL、SQlite、SQL Server 四种数据库，不过可以通过自定义的方式接入其他数据库。

下面以连接 mySQL 为例，首先需要安装两个包。
```
import (
    "gorm.io/driver/mysql" // gorm mysql 驱动包
    "gorm.io/gorm"// gorm
)
```
连接代码。
```
// MySQL 配置信息
username := "root"              // 账号
password := "xxxxxxxx" // 密码
host := "127.0.0.1"             // 地址
port := 3306                    // 端口
DBname := "gorm1"               // 数据库名称
timeout := "10s"                // 连接超时，10秒
dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local&timeout=%s", username, password, host, port, DBname, timeout)
// Open 连接
db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
if err != nil {
    panic("failed to connect mysql.")
}
```
### 三、声明模型
每一张表都会对应一个模型（结构体）。

比如数据库中有一张 goods 表。
![image](https://user-images.githubusercontent.com/6757408/160511500-6e729551-4a34-40f1-801c-00f17012feb5.png)
```
CREATE TABLE `gorm1`.`无标题`  (
  `id` int(0) NOT NULL AUTO_INCREMENT,
  `name` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin NULL DEFAULT NULL,
  `price` decimal(10, 2) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8 COLLATE = utf8_bin ROW_FORMAT = Dynamic;
```
那么就会对应如下的结构体。
```
type Goods struct {
    Id    int
    Name  string
    Price int
}
```
### 约定
gorm 制定了很多约定，并按照约定大于配置的思想工作。

比如会根据结构体的复数寻找表名，会使用 ID 作为主键，会根据 CreateAt、UpdateAt 和 DeletedAt 表示创建时间、更新时间和删除时间。

gorm 提供了一个 Model 结构体，可以将它嵌入到自己的结构体中，省略以上几个字段。
```
type Model struct {
  ID        uint           `gorm:"primaryKey"`
  CreatedAt time.Time
  UpdatedAt time.Time
  DeletedAt gorm.DeletedAt `gorm:"index"`
}
```
嵌入到 goods 结构体中。
```
type Goods struct {
    gorm.Model
    Id    int
    Name  string
    Price int
}
```
这样在每次创建不同的结构体时就可以省略创建 ID、CreatedAt、UpdatedAt、DeletedAt 这几个字段。
### 字段标签 tag
在创建模型时，可以给字段设置 tag 来对该字段一些属性进行定义。

比如创建 Post 结构体，我们希望 Title 映射为 t，设置最大长度为 256，该字段唯一。
```
type Post struct {
    Title string `gorm:"column:t; size:256; unique:true"`
}
```
等同于以下 SQL。
```
CREATE TABLE `posts` (`t, size:256; unique:true` longtext)
```
更多功能可参照下面这张表。
| 标签名  | 说明 |
| ------------- | ------------- |
| column  | 指定 db 列名  |
| type  | 列数据类型，推荐使用兼容性好的通用类型，例如：所有数据库都支持 bool、int、uint、float、string、time、bytes 并且可以和其他标签一起使用，例如：not null、size, autoIncrement… 像 varbinary(8) 这样指定数据库数据类型也是支持的。在使用指定数据库数据类型时，它需要是完整的数据库数据类型，如：MEDIUMINT UNSIGNED not NULL AUTO_INSTREMENT  |
| size  | 指定列大小，例如：size:256  |
| primaryKey  | 指定列为主键  |
| unique  | 指定列为唯一  |
| default  | 指定列的默认值  |
| precision  | 指定列的精度  |
| scale  | 指定列大小  |
| not null  | 指定列为 NOT NULL  |
| autoIncrement  | 指定列为自动增长  |
| embedded  | 嵌套字段  |
| embeddedPrefix  | 嵌入字段的列名前缀  |
| autoCreateTime  | 创建时追踪当前时间，对于 int 字段，它会追踪时间戳秒数，您可以使用 nano/milli 来追踪纳秒、毫秒时间戳，例如：autoCreateTime:nano  |
| autoUpdateTime  | 创建/更新时追踪当前时间，对于 int 字段，它会追踪时间戳秒数，您可以使用 nano/milli 来追踪纳秒、毫秒时间戳，例如：autoUpdateTime:milli  |
| index  | 根据参数创建索引，多个字段使用相同的名称则创建复合索引，查看 索引 获取详情  |
| uniqueIndex  | 与 index 相同，但创建的是唯一索引  |
| check  | 创建检查约束，例如 check:age > 13，查看 约束 获取详情  |
| <-  | 设置字段写入的权限， <-:create 只创建、<-:update 只更新、<-:false 无写入权限、<- 创建和更新权限  |
| ->  | 设置字段读的权限，->:false 无读权限  |
| -  | 忽略该字段，- 无读写权限 |
| comment  | 迁移时为字段添加注释 |
### 四、自动迁移

转自：https://www.luzhenqian.com/blog/gorm-quick-introduction/

可以参见中文官方网址：https://gorm.io/zh_CN/docs/indexes.html





