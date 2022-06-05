## Go语言的Channel遍历方法详解
### 先来看看基本的定义：
channel是Go语言中的一个核心类型，可以把它看成管道。并发核心单元通过它就可以发送或者接收数据进行通讯，这在一定程度上又进一步降低了编程的难度。

channel是一个数据类型，主要用来解决go程的同步问题以及协程之间数据共享（数据传递）的问题。
* channle 本质上是一个数据结构——（队列），数据是先进先出。
* 具有线程安全机制，多个go程访问时，不需要枷锁，也就是说channel本身是线程安全的。
* channel是有类型的，如一个string类型的channel只能存放string类型数据。

### Channel遍历主要分为3种：
#### 1）简单的读 data:=<-ch （如果读多次，需要用循环）
```
var ch8 = make(chan int, 6)
func mm1() {
	for i := 0; i < 10; i++ {
		ch8 <- 8 * i
	}
 
}
func main() {
	go mm1()
	for i:=0;i<10;i++{
		fmt.Print(<-ch8, "\t")
	}
}
```
输出：
![image](https://user-images.githubusercontent.com/6757408/172056135-64f1820e-1dca-4956-aff3-6385ba6320c0.png)

> 注：
>（1）写入的次数与读取的次数需要一致（本例是10）；
>（2）如果读的次数多于写的次数会发生：fatal error: all goroutines are asleep - deadlock! ，若 在mm1中对ch8进行关闭（执行 close(ch8) ），
多于的次数读到的数据为0（数据默认值）。
>（3）读的次数少于写的次数，会读取出次数对应的内容，不会报错。

#### 2）断言方式
if value, ok := <-ch; ok == true {
* 1) 如果写端没有写数据，也没有关闭。<-ch; 会阻塞 ---【重点】
* 2）如果写端写数据， value 保存 <-ch 读到的数据。 ok 被设置为 true
* 3）如果写端关闭。 value 为数据类型默认值。ok 被设置为 false
```
var ch8 = make(chan int, 6)
func mm1() {
	for i := 0; i < 10; i++ {
		ch8 <- 8 * i
	}
	close(ch8)
 
}
func main() {
	go mm1()
	for {
		if data, ok := <-ch8; ok {
			fmt.Print(data,"\t")
		} else {
			break
		}
	}
}
```
![image](https://user-images.githubusercontent.com/6757408/172056243-84b3cb0c-1c59-47c5-b9c1-ca771263ac0e.png)
> 注：写完之后一定要关闭（ 执行：close(ch8) ），否则会出现以下运行结果：
![image](https://user-images.githubusercontent.com/6757408/172056271-a90bfdc6-f597-4a0f-b5de-4b34dadb653f.png)
#### 3）通过range方法
```
for num := range ch {
               }
```
```
var ch8 = make(chan int, 6)
func mm1() {
	for i := 0; i < 10; i++ {
		ch8 <- 8 * i
	}
	close(ch8)
}
func main() {
 
	go mm1()
	for {
		for data := range ch8 {
			fmt.Print(data,"\t")
		}
		break
	}
}
```
> 注：写完之后一定要关闭（ 执行：close(ch8) ），否则会出现以下运行结果：
![image](https://user-images.githubusercontent.com/6757408/172056479-37ec59d0-a253-405f-b343-e4c9203e22ac.png)

> 特别说明：以上实例都是子go程写，主go程读。如在子go程中写，另一个子go程中读，不管哪种方法，都不会出现以上错误问题。（多次实例验证）
```
var ch8 = make(chan int, 6)
func mm1() {
	for i := 0; i < 10; i++ {
		ch8 <- 8 * i
	}
	//close(ch8)
}
func mm2() {
	for {
		for data:=range ch8{
			fmt.Print(data,"\t")
		}
	}
}
func main() {
	go mm1()
	go mm2()
	for{
		runtime.GC()
	}
}
```
![image](https://user-images.githubusercontent.com/6757408/172056573-77a93b71-b0c3-4527-893e-1dddbfd92be6.png)
### 总结：
通过以上验证，为了保证程序的健壮性，在设计程序时，最好将channel的读、写分别在子go程中进行。写完数据之后，记得关闭channel。

补充一点：

channel不像文件一样需要经常去关闭，只有当你确实没有任何发送数据了，或者你想显式的结束range循环之类的，才去关闭channel；关闭channel后，无法向channel 再发送数据
(引发 panic 错误后导致接收立即返回零值)；关闭channel后，可以继续从channel接收数据；对于nil channel，无论收发都会被阻塞。

以上为个人经验，希望能给大家一个参考，也希望大家多多支持我们。如有错误或未考虑完全的地方，望不吝赐教。

