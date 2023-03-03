## golang在Kubernetes中异常重启问题排查
### 问题
Kubernetes上有些服务的某个 Pod总是时不时的重启一下，通过查业务日志根本查不到原因，我分析了一下肯定是哪里代码不严谨造成panic且没有recover才会挂掉的，但是容器重启后之前输出到控制台的
日志会被清空的，所以需要将重启前的日志输出到文件里，这样就能通过容器的volume持久化日志文件的目录方式保留程序崩溃时的信息。

![image](https://user-images.githubusercontent.com/6757408/222646846-a0cb673e-6536-4618-9ba6-ba8012cbc900.png)

### 解决方案
#### 代码
该方式在没有容器的情况下可以使用这种方式，有容器就没必要用这种方式了
服务启动入口，将控制台日志打印到文件中。
> 注意Windows系统不支持的syscall.Dup2这个函数，所以该方法只能在Linux下运行。

```
const stdErrFile = "/tmp/go-panic.log"
var stdErrFileHandler *os.File
func RewriteStderrFile() error {
	if runtime.GOOS == "windows" {
		return nil
	}

	file, err := os.OpenFile(stdErrFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println(err)
		return err
	}
	stdErrFileHandler = file //把文件句柄保存到全局变量，避免被GC回收

	if err = syscall.Dup2(int(file.Fd()), int(os.Stderr.Fd())); err != nil {
		fmt.Println(err)
		return err
	}
	// 内存回收前关闭文件描述符
	runtime.SetFinalizer(stdErrFileHandler, func(fd *os.File) {
		fd.Close()
	})

	return nil
}
```
程序重新部署上去之后，等待服务异常重启，收集日志就能看到刚才程序崩溃时的panic信息，以及导致panic时整个调用栈的信息：
```
fatal error: concurrent map writes

goroutine 12509 [running]:
runtime.throw(0x1622efd, 0x15)
	C:/Program Files/Go/src/runtime/panic.go:1117 +0x72 fp=0xc001803c78 sp=0xc001803c48 pc=0x437092
runtime.mapassign_faststr(0x13b2a40, 0xc0010ca570, 0x16168c3, 0xb, 0x0)
	C:/Program Files/Go/src/runtime/map_faststr.go:211 +0x3f1 fp=0xc001803ce0 sp=0xc001803c78 pc=0x414491
crazyfox-micro/Services/ActivityService/domain/entity/act_base.(*ActivityBase).GetActivityList.func1.1()
	G:/crazyfox-micro/Services/ActivityService/domain/entity/act_base/act_base.go:436 +0x474 fp=0xc001803e98 sp=0xc001803ce0 pc=0x11d3df4
crazyfox-micro/tkpkg/exp.Try(0xc000d49708, 0xc000d496d8)
	G:/crazyfox-micro/tkpkg/exp/exception.go:77 +0x4f fp=0xc001803ec8 sp=0xc001803e98 pc=0x11c8eef
crazyfox-micro/Services/ActivityService/domain/entity/act_base.(*ActivityBase).GetActivityList.func1(0xc00050db30, 0x17cf030, 0xc0010ca390, 0x26acc5, 0xc001730f80, 0xc0014bade0, 0xc04c6f9e7659f073, 0x88221ea5d3, 0x2193700, 0xc0010ca570, ...)
	G:/crazyfox-micro/Services/ActivityService/domain/entity/act_base/act_base.go:423 +0x15a fp=0xc001803f80 sp=0xc001803ec8 pc=0x11d421a
runtime.goexit()
	C:/Program Files/Go/src/runtime/asm_amd64.s:1371 +0x1 fp=0xc001803f88 sp=0xc001803f80 pc=0x46d501
created by crazyfox-micro/Services/ActivityService/domain/entity/act_base.(*ActivityBase).GetActivityList
	G:/crazyfox-micro/Services/ActivityService/domain/entity/act_base/act_base.go:422 +0x1de
```
从日志可以看出是因为并发读写map导致数据竞争的问题。
一些 Go 语言系统级别的错误，比如发生死锁，数据竞争，这种错误程序会立刻报错，无法 recover。
所以在开发过程中需要注意这些问题。
并发安全的map组件可以用：concurrent-map 并发安全且性能高

#### 容器
k8s本身会收集容器日志，可以直接使用指令查看，容器重启前的日志
```
#指令：kubectl logs --previous=true pod名称 -c 容器组名称 -n 命名空间
#--previous=true表示查看之前的日志
kubectl logs --previous=true activityservice-7b9884c749-q6r27 -c activityservice -n crazyfox-production>> log.txt
```


