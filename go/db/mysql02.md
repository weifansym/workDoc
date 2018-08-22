## mysql操作解读
上一篇说明了mysql的基本操作，其实在database/sql中还有好多其他的操作呢，下面我们来看下常用的一些操作。
### 查询操作
查询数据的时候我们可以使用Query和QueryRow，其实Query用来查询多结果集，QueryRow查询单条数据。我们前面说过database/sql连接创建是惰性的，
所以当我们通过Query查询数据的时候主要分为三个步骤：
1. 从连接池中请求一个连接
2. 执行查询的sql语句
3. 将数据库连接的所属权传递给Result结果集
4. 结果集调用Next方法后，自动把链接归还给连接池，或者手动调用结果集的close方法，把链接归还给连接池

#### 多条结果集
前面说了Query返回的结果集是sql.Rows类型。它有一个Next方法，可以迭代数据库的游标，进而获取每一行的数据，使用方法如下：
```
//  执行sql查询
	rows,err := db.Query("SELECT username FROM userinfo WHERE uid>=?", 5)
	checkErr(err)
	// rows中包含了从数据库查的满足uid>=5的所有行的username，
	// rows.Next(), 用于循环迭代获取所有数据
	for rows.Next(){
		var s string
		err = rows.Scan(&s)
		checkErr(err)
		fmt.Println("username: ", s)
	}
	rows.Close()
```
其实当我们通过for循环迭代数据库的时候，当迭代到最后一条数据的时候，会出发一个io.EOF的信号，引发一个错误，同时go会自动调用rows.Close方法释放连接，
然后返回false，此时循环将会结束退出。

通常你会正常迭代完数据然后退出循环。可是如果并没有正常的循环而因其他错误导致退出了循环。此时rows.Next处理结果集的过程并没有完成，
归属于rows的连接不会被释放回到连接池。因此十分有必要正确的处理rows.Close事件。如果没有关闭rows连接，将导致大量的连接被占用，得不到释放，
最终将导致数据库连接池或连接数用法，数据库无法使用。

所以为了避免这种情况的发生，最好的办法就是显示的调用rows.Close方法，确保连接释放，又或者使用defer指令在函数退出的时候释放连接，即使连接已经释放了，
rows.Close仍然可以调用多次，是无害的。

rows.Next循环迭代的时候，因为触发了io.EOF而退出循环。为了检查是否是迭代正常退出还是异常退出，需要检查rows.Err。例如上面的代码应该改成：
```
//  执行sql查询
	rows,err := db.Query("SELECT username FROM userinfo WHERE uid>=?", 5)
	checkErr(err)
	// rows中包含了从数据库查的满足uid>=5的所有行的username，
	// rows.Next(), 用于循环迭代获取所有数据
	for rows.Next(){
		var s string
		err = rows.Scan(&s)
		checkErr(err)
		fmt.Println("username: ", s)
	}
	rows.Close()

	err = rows.Err();
	checkErr(err)
```
#### 单条结果集
#### 空值的处理
#### 查询字段的自动匹配
### 插入操作
#### 单个插入
#### 批量插入
