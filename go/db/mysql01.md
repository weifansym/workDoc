# 使用MySQL数据库
目前Internet上流行的网站构架方式是LAMP，其中的M即MySQL, 作为数据库，MySQL以免费、开源、使用方便为优势成为了很多Web开发的后端数据库存储引擎。

## MySQL驱动
Go中支持MySQL的驱动目前比较多，有如下几种，有些是支持database/sql标准，而有些是采用了自己的实现接口,常用的有如下几种:

- https://github.com/go-sql-driver/mysql  支持database/sql，全部采用go写。
- https://github.com/ziutek/mymysql   支持database/sql，也支持自定义的接口，全部采用go写。
- https://github.com/Philio/GoMySQL 不支持database/sql，自定义接口，全部采用go写。

接下来的例子我主要以第一个驱动为例(我目前项目中也是采用它来驱动)，也推荐大家采用它，主要理由：

- 这个驱动比较新，维护的比较好
- 完全支持database/sql接口
- 支持keepalive，保持长连接,虽然[星星](http://www.mikespook.com)fork的mymysql也支持keepalive，但不是线程安全的，这个从底层就支持了keepalive。

## 示例代码
接下来的几个小节里面我们都将采用同一个数据库表结构：数据库test，用户表userinfo，关联用户信息表userdetail。
```sql

CREATE TABLE `userinfo` (
	`uid` INT(10) NOT NULL AUTO_INCREMENT,
	`username` VARCHAR(64) NULL DEFAULT NULL,
	`department` VARCHAR(64) NULL DEFAULT NULL,
	`created` DATE NULL DEFAULT NULL,
	PRIMARY KEY (`uid`)
);

CREATE TABLE `userdetail` (
	`uid` INT(10) NOT NULL DEFAULT '0',
	`intro` TEXT NULL,
	`profile` TEXT NULL,
	PRIMARY KEY (`uid`)
)
```
如下示例将示范如何使用database/sql接口对数据库表进行增删改查操作
```Go
package main

import (
	"database/sql"
	"fmt"
	"sync"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

type user struct {
	username string
	department string
	created string
}

func userCreate(db *sql.DB, wg *sync.WaitGroup) {
	defer wg.Done()
	sqlStr := "INSERT INTO userinfo(username, department, created) VALUES (?,?,?)"
	res, err := db.Exec(sqlStr, "13", "33", "2018-03-03")
	checkErr(err)

	resultId, err := res.LastInsertId()
	checkErr(err)
	fmt.Println("resultId: ", resultId)

	affectRow, err := res.RowsAffected()
	checkErr(err)
	fmt.Println("affectRow: ", affectRow)
}

func getUserList(db *sql.DB, wg *sync.WaitGroup) []user  {
	defer wg.Done()
	userList := []user{}
	sqlStr := "SELECT username, department, created FROM userinfo"
	rows, err := db.Query(sqlStr)
	checkErr(err)

	for rows.Next() {
		var name string
		var department string
		var created string

		err := rows.Scan(&name, &department, &created)
		checkErr(err)
		fmt.Println("name: ", name)
		fmt.Println("department: ", department)
		fmt.Println("created: ", created)
		userList = append(userList, user{name, department, created})
	}
	fmt.Println("userList: ", userList)
	return userList
}

func main()  {
	db, err := sql.Open("mysql", "root:duhuo126@/test?charset=utf8");
	checkErr(err)

	//  插入数据

	/*
	//  插入方式：1
	sqlStr := "INSERT INTO userinfo(username, department, created) VALUES (?,?,?)"
	res, err := db.Exec(sqlStr, "13", "33", "2018-03-03")
	checkErr(err)

	resultId, err := res.LastInsertId()
	checkErr(err)
	fmt.Println("resultId: ", resultId)

	affectRow, err := res.RowsAffected()
	checkErr(err)
	fmt.Println("affectRow: ", affectRow)*/

	//  插入方式：2
	/*sqlStr := "INSERT INTO userinfo(username, department, created) VALUES (?,?,?)"
	stmt, err := db.Prepare(sqlStr)
	checkErr(err)

	res, err := stmt.Exec("22", "22", "2019-01-01")
	checkErr(err)

	resultId, err := res.LastInsertId()
	checkErr(err)
	fmt.Println("resultId: ", resultId)

	affectRow, err := res.RowsAffected()
	checkErr(err)
	fmt.Println("affectRow: ", affectRow)*/

	/*//  查询
	userList := []user{}
	sqlStr := "SELECT username, department, created FROM userinfo"
	rows, err := db.Query(sqlStr)
	checkErr(err)

	for rows.Next() {
		var name string
		var department string
		var created string

		err := rows.Scan(&name, &department, &created)
		checkErr(err)
		fmt.Println("name: ", name)
		fmt.Println("department: ", department)
		fmt.Println("created: ", created)
		userList = append(userList, user{name, department, created})
	}
	fmt.Println("userList: ", userList)*/


	//  更新数据
	/*
	//  更新方式：1
	sqlStr := "UPDATE userinfo SET username=? where username=?"
	rest,err := db.Exec(sqlStr, "duhuo", "11")
	checkErr(err)
	resultId, err := rest.LastInsertId()
	checkErr(err)
	fmt.Println("update resultId: ", resultId)

	affectRow, err := rest.RowsAffected()
	checkErr(err)
	fmt.Println("update affectRow: ", affectRow)*/

	/*
	//  更新方式：2
	stmt, err := db.Prepare("update userinfo set username=? where uid=?")
	checkErr(err)

	res, err := stmt.Exec("astaxieupdate", 3)
	checkErr(err)

	affect, err := res.RowsAffected()
	checkErr(err)

	fmt.Println("update RowsAffected: ", affect)*/

	//  删除数据
	/*
	//  删除方式：1
	results,err := db.Exec("DELETE from userinfo where uid=?",5)
	checkErr(err)
	affect, err := results.RowsAffected()
	checkErr(err)
	fmt.Println("delete RowsAffected: ", affect)*/

	/*
	//  删除方式：2
	stmt, err := db.Prepare("delete from userinfo where uid=?")
	checkErr(err)

	res, err := stmt.Exec(12)
	checkErr(err)

	affect, err := res.RowsAffected()
	checkErr(err)
	fmt.Println("delete RowsAffected: ", affect)*/

	//  开启不同的协程处理数据
	/*var wg sync.WaitGroup
	go userCreate(db, &wg)
	wg.Add(1)
	go getUserList(db, &wg)
	wg.Add(1)

	wg.Wait()
	fmt.Println("all goroutine has finished!!!")*/


	db.Close()
}
```

通过上面的代码我们可以看出，Go操作Mysql数据库是很方便的。

关键的几个函数我解释一下：

sql.Open()函数用来打开一个注册过的数据库驱动，go-sql-driver中注册了mysql这个数据库驱动，第二个参数是DSN(Data Source Name)，它是go-sql-driver定义的一些数据库链接和配置信息。它支持如下格式：

	user@unix(/path/to/socket)/dbname?charset=utf8
	user:password@tcp(localhost:5555)/dbname?charset=utf8
	user:password@/dbname
	user:password@tcp([de:ad:be:ef::ca:fe]:80)/dbname

db.Prepare()函数用来返回准备要执行的sql操作，然后返回准备完毕的执行状态。

db.Query()函数用来直接执行Sql返回Rows结果。

stmt.Exec()函数用来执行stmt准备好的SQL语句

我们可以看到我们传入的参数都是=?对应的数据，这样做的方式可以一定程度上防止SQL注入。

具体需要参见go的相关数据接口包：[database/sql](https://golang.org/pkg/database/sql/), 以及教程[Go database/sql tutorial](http://go-database-sql.org/index.html)

参考：https://github.com/astaxie/build-web-application-with-golang/blob/master/zh/05.2.md
