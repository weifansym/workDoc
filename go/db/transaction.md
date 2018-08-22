## mysql事务
事务在关系型数据库中是非常重要的特性，在好多为了保证数据一致性的地方都需要使用到事务来处理。database/sql提供了事务处理的功能。通过Tx对象实现。
db.Begin会创建tx对象，后者的Exec和Query执行事务的数据库操作，最后在tx的Commit和Rollback中完成数据库事务的提交和回滚，同时释放连接。

### tx对象
我们在之前查询以及操作数据库都是用的db对象，而事务则是使用另外一个对象。使用db.Begin 方法可以创建tx对象，tx对象同样具有数据库操作的Query,Exec方法
用法和我们前面的db操作基本一样，但是需要在查询或者操作完毕之后执行tx对象的Commit提交或者Rollback方法回滚。

一旦创建了这个tx对象，则这个事务中所有的处理都依赖于这个tx对象，这个对象会从连接池中取出一个空闲的连接，接下来的sql执行都基于这个连接，
直到调用commit或Roolback之后，才会把这个连接释放到连接池。

具体golang相关的数据库操作请看：
* https://golang.org/pkg/database/sql/
* http://go-database-sql.org/index.html

在事务处理的时候，不能使用db的查询方法，当然你如果使用也能执行语句成功，但是这和你事务里执行的操作将不是一个事务，将不会接受commit和rollback的改变，
如下面操作时：
```
tx,err := Db.Begin()
Db.Exec()
tx.Exec()
tx.Commit()
```
上面这个伪代码中，调用Db.Exec方法的时候，和tx执行Exec方法时候是不同的，只有tx的会绑定到事务中，db则是额外的一个连接，两者不是同一个事务。

创建Tx对象的时候，会从连接池中取出连接，然后调用相关的Exec方法的时候，连接仍然会绑定在该事务处理中。事务的连接生命周期从Beigin函数调用起，直到Commit和Rollback函数的调用结束。

### tx对象与db对象
对于下面的代码：
```
rows, _ := db.Query("SELECT uid FROM userinfo") 
for rows.Next() {
    var uid int
    var name string
    rows.Scan(&uid)
    db.QueryRow("SELECT username FROM userinfo WHERE uid = ?", uid).Scan(&name)
}
```
调用了Query方法之后，在Next方法中取结果的时候，rows是维护了一个连接，再次调用QueryRow的时候，db会再从连接池取出一个新的连接。
rows和db的连接两者可以并存，并且相互不影响。

对于sql.Tx对象，因为事务过程只有一个连接，事务内的操作都是顺序执行的，只有当前的数据库操作完成之后才能进行下一个数据库操作，
上面的逻辑如果在事务处理中会失效，如下代码：
```
rows, _ := tx.Query("SELECT uid FROM userinfo")
for rows.Next() {
   var uid int
   var name string
   rows.Scan(&uid)
   tx.QueryRow("SELECT username FROM userinfo WHERE uid = ?", uid).Scan(&name)
}
```
tx执行了Query方法后，连接转移到rows上，在Next方法中，tx.QueryRow将尝试获取该连接进行数据库操作。因为还没有调用rows.Close，
因此底层的连接属于busy状态，tx是无法再进行查询的。

报错信息如下：
```
[mysql] 2018/08/22 15:13:36 packets.go:427: busy buffer
2018/08/22 15:13:36 driver: bad connection
```
## 完整实例如下：
```
package main

import (
	"log"
	"fmt"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

func doSomething(){
	panic("A Panic Running Error")
}

func checkErr(err error) {
	if err != nil {
		fmt.Println("check error: ", err)
	}
}

func clearTransaction(tx *sql.Tx){
	err := tx.Rollback()
	fmt.Println("事务回滚了！！！！")
	if err != sql.ErrTxDone && err != nil{
		log.Fatalln(err)
	}
}

func main() {
	db, err := sql.Open("mysql", "root:duhuo126@/test?charset=utf8");
	checkErr(err)

	defer db.Close()

	//  开启一个事务
	tx, err := db.Begin()
	checkErr(err)

	defer clearTransaction(tx)
	
	rs, err := tx.Exec("UPDATE userinfo SET username=? WHERE uid=?", "testTx11", 8)
	checkErr(err)

	rowAffected, err := rs.RowsAffected()
	checkErr(err)
	fmt.Println("one rowAffected: ", rowAffected)

	rs, err = tx.Exec("UPDATE userinfo SET score=? WHERE uid=?", 1111, 9)
	checkErr(err)

	rowAffected, err = rs.RowsAffected()
	checkErr(err)
	fmt.Println("one rowAffected: ", rowAffected)

	doSomething()

	err = tx.Commit();
	checkErr(err)
}
```
这里定义了一个clearTransaction(tx)函数，该函数会执行rollback操作。因为我们事务处理过程中，任何一个错误都会导致main函数退出，
因此在main函数退出执行defer的rollback操作，回滚事务和释放连接。

如果不添加defer，只在最后Commit后check错误err后再rollback，那么当doSomething发生异常的时候，函数就退出了，此时还没有执行到tx.Commit。
这样就导致事务的连接没有关闭，事务也没有回滚。

tx事务环境中，只有一个数据库连接，事务内的Eexc都是依次执行的，事务中也可以使用db进行查询，但是db查询的过程会新建连接，这个连接的操作不属于该事务。
