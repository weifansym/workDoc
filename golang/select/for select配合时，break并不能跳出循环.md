## for select配合时，break并不能跳出循环
> 通常在for循环中，使用break可以跳出循环，但是注意在go语言中，for select配合时，break并不能跳出循环。

```
func testSelectFor(chExit chan bool){
	for  {
		select {
		case v, ok := <-chExit:
			if !ok {
				fmt.Println("close channel 1", v)
				break
			}
			fmt.Println("ch1 val =", v)
		}
	}
	fmt.Println("exit testSelectFor")
}
```
如下调用：
```
//尝试2 select for 跳出循环
c := make(chan bool)
go testSelectFor(c)
 
c <- true
c <- false
close(c)
 
time.Sleep(time.Duration(2) * time.Second)
```
运行结果如下，可以看到break无法跳出循环：
```
...
close channel 1 false
close channel 1 false
close channel 1 false
close channel 1 false
...
```
为了解决这个问题，**需要设置标签，break 标签或goto 便签即可跳出循环**，如下两种方法均可。
```
func testSelectFor2(chExit chan bool){
	EXIT:
	for  {
		select {
		case v, ok := <-chExit:
			if !ok {
				fmt.Println("close channel 2", v)
				break EXIT//goto EXIT2
			}
 
			fmt.Println("ch2 val =", v)
		}
	}
	//EXIT2:
	fmt.Println("exit testSelectFor2")
}
```
同样调用，输出结果如下：
```
ch2 val = true
ch2 val = false
close channel 2 false
exit testSelectFor2
```


