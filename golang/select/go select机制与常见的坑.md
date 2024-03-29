## go select机制与常见的坑
go select思想来源于网络IO模型中的select，本质上也是IO多路复用，只不过这里的IO是基于channel而不是基于网络，同时go select也有一些自己不同的特性，这里简单探讨下。

go select 的特性:
1. 每个case都必须是一个通信
2. 所有channel表达式都会被求值
3. 所有被发送的表达式都会被求值
4. 如果任意某个通信可以进行，它就执行；其他被忽略。
5. 如果有多个case都可以运行，select会随机公平地选出一个执行。其他不会执行。否则执行default子句(如果有)
6. 如果没有default字句，select将阻塞，直到某个通信可以运行；Go不会重新对channel或值进行求值。

下面通过几个例子来理解这些特性:

#### 1.select closed/nil channel
```
for {
	select {
	case v1, ok := <-c1:
        // 如果c1被关闭(ok==false)，每次从c1读取都会立即返回，将导致死循环
        // 可以通过将c1置为nil来让select ignore掉这个case，继续评估其它case
		if !ok {
			c1 = nil
		}
	}
	
	case v2 := <- c2:
	    // 同样，如果c2被关闭，每次从c1读取都会立即返回对应元素类型的零值(如空字符串)，导致死循环
	    // 解决方案仍然是置c2为nil，但是有可能误判(写入方是写入了一个零值而不是关闭channel，比如整数0)
	    
	case c3 <- v3:
	    // 如果c3已经关闭，则panic
	    // 如果c3为nil，则ignore该case	    
}
```
#### 2.实现非阻塞读写

结合特性5,6，可以通过带 default 语句的 select 实现非阻塞读写，在实践中还是比较有用的，比如 GS 尝试给玩家推送某条消息，可能并不希望 GS 阻塞在该玩家的 writeChan 上。
```
select {
    case writeChan <- msg:
        // do something write successed
    default:
        // drop msg, or log err
}
```
需要注意，一些同学可能将select与switch搞混，习惯先把default写好，然后加上外层的for循环导致死循环。**使用select语句，for和default基本不会同时出现。**

#### 3.实现定时任务

结合特性2，每次 select 都会对所有通信表达式求值，因此可通过 time.After简洁实现定时器功能，并且定时任务可通过 done channel 停止:
```
for {
	select {
	case <- time.After(time.Second):
	    // do something per second
	case <- donec:
		return	
	}
}
```
现在我们稍微变更一下:
```
donec := make(chan bool, 1)
close(donec)
for {
	select {
	case <- time.After(time.Second):
		fmt.Println("timer")
	case <- donec:
	}
}
```
现在这段代码会输出什么？还是 panic？答案是什么也不会，因为:
1. donec close 了，每次 select 都会执行到 case <- donec，并读出零值(false)
2. 每次执行了 case <- donec1 后，select 再次对 case1 的 timer.After 求值，返回一个新的下一秒超时的 Timer
3. 再次执行到 case <- donec ….

因此，case <- timer.After(time.Second) 不应该解释为每一秒执行一次，而是其它 case 如果有一秒都没有执行，那么就执行这个 case。
#### 4.多个case满足读写条件

结合特性4，如果多个case满足读写条件，select会随机选择一个语句执行：
```
func main() {
	ch := make(chan int, 1024)
	go func(ch chan int) {
		for {
			val := <-ch
			fmt.Printf("val:%d\n", val)
		}
	}(ch)
    
	tick := time.NewTicker(1 * time.Second)
	for i := 0; i < 5; i++ {
		select {
		case ch <- i:
		case <-tick.C:
			fmt.Printf("%d: case <-tick.C\n", i)
		}
    
		time.Sleep(500 * time.Millisecond)
	}
	close(ch)
	tick.Stop()
}
```
输出:
```
val:0
val:1
2: case <-tick.C
val:3
4: case <-tick.C
```
可以看到向ch写入的2和4”不见”了，因为当tick.C和ch同时满足读写条件时，select随机选择了一个执行，导致看起来一些数据丢了，其实这个例子是比较极端的，因为向ch写入的数据本身就与外部for循环计数耦合了，导致依赖于select的随机结果(本次没随机到，放到下次，但此时写入的数据已经变更了)，因此实际不是数据丢了，而是代码设计时没有考虑到每次select只会执行一条读写语句(并且是随机选取的)，导致结果不如预期。

总的来说，go select还是比较容易踩坑的，比如加了不该加的default，没有考虑到channel关闭的情况，没有理解随机性等等，在使用的时候还是要小心。

参考：
1: https://wudaijun.com/2017/10/go-select/
2: http://www.codingbygolang.com/select-io-case/#select-%E5%88%86%E6%94%AF%E7%9A%84%E9%9A%8F%E6%9C%BA%E6%80%A7

