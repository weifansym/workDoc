## 并发遇到的问题
### range在遍历信道时遇到的问题
range在遍历信道中的数据的时候，如果信道不关闭，range会一直阻塞当前协程，知道信道被关闭？？
```
package main

import (
	"fmt"
)

func producer(chnl chan int) {
	for i := 0; i < 10; i++ {
		chnl <- i
	}
	//close(chnl)
}
func main() {
	ch := make(chan int)
	go producer(ch)
	fmt.Println("AAAAAAA")
	//  信道未关闭的时候range会一直阻塞当前协程
	for v := range ch {
		fmt.Println("Received ",v)
	}
	fmt.Println("BBBBB")
}
```
