## golang 性能优化分析：benchmark 结合 pprof 
前面 2 篇 golang 性能优化分析系列文章：
* [golang 性能优化分析工具 pprof (上)](https://www.cnblogs.com/jiujuan/p/14588185.html)
* [golang 性能优化分析工具 pprof (下)](https://www.cnblogs.com/jiujuan/p/14598141.html)

### 一、基准测试 benchmark 简介#
在 golang 中，可以通过 benchmark 基准测试来测试代码性能。基准测试主要是通过测试 cpu 和内存的效率问题，来评估被测试代码的性能。

基准测试的指标：
1. 程序所花费的时间
2. 内存使用的情况
3. cpu 使用情况

基准测试文件名和函数规定：
* go 基准测试文件都是以 _test.go 结尾，和单元测试用例在同一个文件中。
* 基准测试每个函数都是以 Benchmark 开头。

基准测试常用命令：
```
go test ./fib              // 不进行基准测试，对 fib 进行单元测试

go test -bench=. -run=none  // 进行基准测试，不进行单元测试，-run 表示执行哪些单元测试和测试函数，一般函数名不会是 none，所以不执行单元测试
// 上面的测试命令还可以用空格隔开，意义是一样
go test -bench . -run none

go test -bench=.    // 对所有的进行基准测试

go test -bench='fib$'     // 只运行以 fib 结尾的基准测试，-bench 可以进行正则匹配

go test -bench=. -benchtime=6s  // 基准测试默认时间是 1s，-benchtime 可以指定测试时间
go test -bench=. -benchtime=50x  // 参数 -benchtime 除了指定时间，还可以指定运行的次数

go test -bench=. -benchmem // 进行时间、内存的基准测试
```
说明：上面的命令中，**-bench**后面都有一个 . ，这个点并不是指当前文件夹，而是一个匹配所有测试的正则表达式。

更多参数说明请查看帮助：**go help testflag**

分析基准测试数据：
* cpu 使用分析：-cpuprofile=cpu.pprof
* 内存使用分析：-benchmem -memprofile=mem.pprof
* block分析：-blockprofile=block.pprof

在配合 pprof 就可以进行分析。

运行命令采样数据：
```
go test -bench=. -run=none -benchmem -memprofile=mem.pprof
go test -bench=. -run=none -blockprofile=block.pprof

go test -bench=. -run=none -benchmem -memprofile=mem.pprof -cpuprofile=cpu.pprof
```
### 二、代码示例#
#### 2.1 代码示例#
fib.go：
```
package main

func Fib(n int) int {
	if n < 2 {
		return n
	}

	return Fib(n-1) + Fib(n-2)
}
```
fib_test.go:
```
package main

import (
	"testing"
)

func BenchmarkFib(b *testing.B) {
	// 运行 Fib 函数 b.N 次
	for n := 0; n < b.N; n++ {
		Fib(20)
	}
}

func BenchmarkFib2(b *testing.B) {
	// 运行 Fib 函数 b.N 次
	for n := 0; n < b.N; n++ {
		Fib(10)
	}
}
```
#### 2.2 运行命令采集数据#
```
go test -bench=. -run=none \
-benchmem -memprofile=mem.pprof \
-cpuprofile=cpu.pprof \
-blockprofile=block.pprof
```
也可以用一个一个命令来完成采集数据，分开运行：
> go test -bench=. -run=none -benchmem -memprofile=mem.pprof

> go test -bench=. -run=none -benchmem -cpuprofile=cpu.pprof

#### 2.3 分析数据#
前面有 上，下 两篇 pprof 的文章怎么分析数据，一种方法是命令行交互分析模式，一种是可视化图形分析模式。

A. 命令行交互分析

分析 cpu：
> go tool pprof cpu.pprof

>再用 top15 命令分析，或者 top --cum 进行排序分析

如下图：

![image](https://github.com/weifansym/workDoc/assets/6757408/33d962d6-73cc-4b39-be2f-4f4316afdd28)

B. web 界面分析

命令行执行命令：
> go tool pprof -http=":8080" cpu.pprof

会自动在浏览器上打开地址：http://localhost:8080/ui/ ，然后就可以在浏览器上查看各种分析数据，如下图：

![image](https://github.com/weifansym/workDoc/assets/6757408/7cfbe788-8347-46eb-859c-c9a71ee63840)

其他数据也可以进行同样的分析，这里就略过。

[完]
转自：https://www.cnblogs.com/jiujuan/p/14604609.html


  
