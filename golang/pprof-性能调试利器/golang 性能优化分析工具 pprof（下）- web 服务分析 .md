## golang 性能优化分析工具 pprof（下）- web 服务分析 
[golang 性能优化分析工具 pprof（上）篇-基础使用介绍](https://www.cnblogs.com/jiujuan/p/14588185.html)

### 四、web 服务(http server)的分析 net/http/pprof#
#### 4.1 代码例子 1
> go version go1.13.9

把上面的程序例子稍微改动下，命名为 demohttp.go:
```
package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	_ "net/http/pprof"
	"sync"
)

func main() {
	http.HandleFunc("/pprof-test", handler)

	fmt.Println("http server start")
	err := http.ListenAndServe(":8090", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func handler(resp http.ResponseWriter, req *http.Request) {
	var wg sync.WaitGroup
	wg.Add(200)

	for i := 0; i < 200; i++ {
		go cyclenum(30000, &wg)
	}

	wg.Wait()

	wb := writeBytes()
	b, err := ioutil.ReadAll(wb)
	if err != nil {
		resp.Write([]byte(err.Error()))
		return
	}
	resp.Write(b)
}

func cyclenum(num int, wg *sync.WaitGroup) {
	slice := make([]int, 0)
	for i := 0; i < num; i++ {
		for j := 0; j < num; j++ {
			j = i + j
			slice = append(slice, j)
		}
	}
	wg.Done()
}

func writeBytes() *bytes.Buffer {
	var buff bytes.Buffer

	for i := 0; i < 30000; i++ {
		buff.Write([]byte{'a' + byte(rand.Intn(10))})
	}
	return &buff
}
```
#### 4.2 开始分析#
##### 4.2.1 在 web 界面上分析#
先运行上面的 demohttp.go 程序，执行命令：
> go run demohttp.go

然后在浏览器输入：http://localhost:8090/debug/pprof/，查看服务运行情况，如下图：
![image](https://github.com/weifansym/workDoc/assets/6757408/d143d56b-e79d-43e8-a2a6-0486726b8570)

| 名称  | url | 说明 |
| ------------- | ------------- | ------------- |
| allocs  | $host/debug/pprof/allocs?debug=1  | 过去所有内存抽样情况  |
| block  | $host/debug/pprof/block?debug=1 | 同步阻塞时程序栈跟踪的一些情况 |
| heap  | $host/debug/pprof/heap?debug=1  | 活动对象的内存分配情况 |
| mutex  | $host/debug/pprof/mutex?debug=1  | 互斥锁持有者的栈帧情况  |
| profile  | $host/debug/pprof/profile | cpu profile，点击时会得到一个文件，然后可以用 go tool pprof 命令进行分析  |
| threadcreate  | $host/debug/pprof/threadcreate?debug=1  | 创建新 OS 线程的堆栈跟踪情况  |
| trace  | $host/debug/pprof/trace  | 当前程序执行的追踪情况，点击时会得到一个文件，可以用 go tool trace 命令来分析这个文件  |

点击上面的链接，就可以查看具体的分析情况。
不断刷新网页，可以看到数据在不断变化。

#### 4.2.2 命令行交互分析#
在命令行上运行 demohttp.go 程序，执行命令:
> go run demohttp.go

##### A. 分析 cpu profile

在开启另外一个命令行终端，执行如下命令：
> go tool pprof http://localhost:8090/debug/pprof/profile?seconds=70
![image](https://github.com/weifansym/workDoc/assets/6757408/c93cdfcb-df39-493d-bd5f-a4c914a64801)

参数 seconds = 70：进行 70s 的数据样本采集，这个参数可以根据实际情况调整。

上面的命令执行后，会等待 70s ， 然后才会进入命令交互界面，如上图

输入**top**命令：
![image](https://github.com/weifansym/workDoc/assets/6757408/0fd9fae4-d38e-4cb1-b2df-b8aed8cda967)

大家发现没，其实与上面 runtime/pprof 在命令行交互时是一样的操作，可以参考上面的字段参数说明。

找出耗时代码部分，也可以用命令：**list**。

在**top**命令执行后，发现什么问题没？这个 top 命令显示的信息都是系统调用信息耗时，没有用户定义的函数。为什么？下面进行分析。

##### B. 分析 memory profile
执行命令：
> go tool pprof http://localhost:8090/debug/pprof/heap

然后同样输入 top 命令查看函数使用情况，如下图：
![image](https://github.com/weifansym/workDoc/assets/6757408/a402a55a-18e1-4a21-9a85-fcd31098dfea)

其余的跟踪分析命令类似，就不一一分析了。

把上面在终端命令行下交互分析的数据进行可视化分析。

#### 4.2.3 图形可视化分析#
##### A. pprof 图形可视化
在前面可视化分析中，我们了解到可视化最重要有 2 步：1.采集数据 2.图形化采集的数据。

在上面第三节 runtime/pprof 中，进入终端命令行交互操作，然后输入 web 命令，就可以生成一张 svg 格式的图片，用浏览器可以直接查看该图片。我们用同样的方法来试一试。
1. 输入命令：
> go tool pprof http://localhost:8090/debug/pprof/profile?seconds=30

2. 等待 30s 后输入 web 命令, 如下图：
![image](https://github.com/weifansym/workDoc/assets/6757408/1a3beac4-c3b2-493e-8f82-9dc317598e87)

果然生成了一个 svg 文件，在浏览器查看该图片文件，啥有用信息也没有，如下图：
![image](https://github.com/weifansym/workDoc/assets/6757408/37f36f48-6aa7-484d-8790-5a3d2e5696ae)

为什么没有有用信息？前面有讲到过，没有用户访问 http server ，需要的程序没有运行，一直阻塞在那里等待客户端的访问连接，所以 go tool pprof 只能采集部分代码运行的信息，而这部分代码又没有消耗多少 cpu。

那怎么办？

一个方法就是用 http 测试工具模拟用户访问。这里用 https://github.com/rakyll/hey 这个工具。
安装 hey：
> go get -u github.com/rakyll/hey

安装完成后，进行 http 测试：
> hey -n 1000 http://localhost:8090/pprof-test

同时开启另一终端执行命令：
> go tool pprof http://localhost:8090/debug/pprof/profile?seconds=120

等待 120s 后，采集信息完成，如下图：
![image](https://github.com/weifansym/workDoc/assets/6757408/6f48298f-d06e-4c9e-83ee-809e4e24251e)

输入**top**命令查看统计信息：

![image](https://github.com/weifansym/workDoc/assets/6757408/a133bde1-1f8f-4e78-a9a6-49c14ab22a3a)

可以看到用户定义的一个最耗时函数是：**main.cyclenum**。如果要查看这个函数最耗时部分代码，可以用**list cyclenum**命令查看。

我们这里是要生成一张图片，所以输入**web**命令生成图片：
![image](https://github.com/weifansym/workDoc/assets/6757408/8435fcdf-63b6-4ec6-84a3-1d5afea0108d)

在浏览器上查看 svg 图片：
![image](https://github.com/weifansym/workDoc/assets/6757408/23120a32-3510-4a7f-a921-8f91406aff2f)

(图片较大，只截取了部分)

这张图完整的展示了**top**命令的信息。

##### B. web 可视化
执行命令：
> go tool pprof -http=":8080" http://localhost:8090/debug/pprof/profile

同时开启另一终端执行测试命令：
> hey -n 200 -q 5 http://localhost:8090/pprof-test

上面**go tool pprof**执行完成后，会自动在浏览器打开一个 http 地址，http://localhost:8080/ui/，如下图：
![image](https://github.com/weifansym/workDoc/assets/6757408/002e2cae-3f36-4416-9e7f-a7463538c506)

(截取部分图片)

这样就可以在web浏览器上查看分析数据了。

##### C. 火焰图
用 http 测试框架 hey 访问，命令为：
> hey -n 200 -q 5 http://localhost:8090/pprof-test

在压测的同时开启另一终端执行命令：
> go-torch -u http://localhost:8090

来生成火焰图。

运行命令时在终端输出了信息 ：
> Run pprof command: go tool pprof -raw -seconds 30 http://localhost:8090/debug/pprof/profile

可以看到**go-torch**的原始命令也是用到了**go tool pprof**

上面这个命令默认生成了**torch.svg**的火焰图文件，如下：
![image](https://github.com/weifansym/workDoc/assets/6757408/b34783e6-7e4e-4f2d-894b-75dcb565cf0a)

(截取一部分图展示)

点击方块可以查看更详细信息:

![image](https://github.com/weifansym/workDoc/assets/6757408/1113c849-0419-46bf-95e5-1d78aba43395)


### 参考#
* [pprof](https://github.com/google/pprof)
  * [README](https://github.com/google/pprof/blob/main/doc/README.md)
* [Profiling Go Programs](https://go.dev/blog/pprof)
* [runtime/pprof](https://pkg.go.dev/runtime/pprof)
* [net/http/pprof](https://pkg.go.dev/net/http/pprof)
* [go-torch](https://github.com/uber-archive/go-torch)
* [Flame Graph](https://github.com/brendangregg/FlameGraph)
* [http 压测工具 hey](https://github.com/rakyll/hey)

转自：https://www.cnblogs.com/jiujuan/p/14598141.html






