## 高德Go生态的服务稳定性建设｜性能优化的实战总结
### 前言
go语言凭借着优秀的性能，简洁的编码风格，极易使用的协程等优点，逐渐在各大互联网公司中流行起来。而高德业务使用go语言已经有3年时间了，随着高德业务的发展，go语言生态也日
趋完善，今后会有越来越多新的go服务出现。在任何时候，保障服务的稳定性都是首要的，go服务也不例外，而性能优化作为保障服务稳定性，降本增效的重要手段之一，在高德go服务日益
普及的当下显得愈发重要。此时此刻，我们将过去go服务开发中的性能调优经验进行总结和沉淀，为您呈上这篇精心准备的go性能调优指南。

通过本文您将收获以下内容： 
1. 从理论的角度，和你一起捋清性能优化的思路，制定最合适的优化方案。
2. 推荐几款go语言性能分析利器，与你一起在性能优化的路上披荆斩棘。
3. 总结归纳了众多go语言中常用的性能优化小技巧，总有一个你能用上。
4. 基于高德go服务百万级QPS实践，分享几个性能优化实战案例，让性能优化不再是纸上谈兵。

### 一、性能调优-理论篇
#### 1.1 衡量指标
优化的第一步是先衡量一个应用性能的好坏，性能良好的应用自然不必费心优化，性能较差的应用，则需要从多个方面来考察，找到木桶里的短板，才能对症下药。那么如何衡量一个应用
的性能好坏呢？最主要的还是通过观察应用对核心资源的占用情况以及应用的稳定性指标来衡量。**所谓核心资源，就是相对稀缺的，并且可能会导致应用无法正常运行的资源**，
常见的核心资源如下：

* **cpu**：对于偏计算型的应用，cpu往往是影响性能好坏的关键，如果代码中存在无限循环，或是频繁的线程上下文切换，亦或是糟糕的垃圾回收策略，都将导致cpu被大量占用，使得应用程序无法获取到足够的cpu资源，从而响应缓慢，性能变差。

* **内存**：内存的读写速度非常快，往往不是性能的瓶颈，但是内存相对来说容量有限切价格昂贵，如果应用大量分配内存而不及时回收，就会造成内存溢出或泄漏，应用无法分配新的内存，便无法正常运行，这将导致很严重的事故。

* **带宽**：对于偏网络I/O型的应用，例如网关服务，带宽的大小也决定了应用的性能好坏，如果带宽太小，当系统遇到大量并发请求时，带宽不够用，网络延迟就会变高，这个虽然对服务端可能无感知，但是对客户端则是影响甚大。

* 磁盘：相对内存来说，磁盘价格低廉，容量很大，但是读写速度较慢，如果应用频繁的进行磁盘I/O，那性能可想而知也不会太好。

以上这些都是系统资源层面用于衡量性能的指标，除此之外还有应用本身的稳定性指标：
* **异常率**：也叫错误率，一般分两种，执行超时和应用panic。panic会导致应用不可用，虽然服务通常都会配置相应的重启机制，确保偶然的应用挂掉后能重启再次提供服务，但是经常性的panic，会导致应用频繁的重启，减少了应用正常提供服务的时间，整体性能也就变差了。**异常率是非常重要的指标，服务的稳定和可用是一切的前提，如果服务都不可用了，还谈何性能优化**。
* **响应时间(RT)**：包括平均响应时间，百分位(top percentile)响应时间。响应时间是指应用从收到请求到返回结果后的耗时，反应的是应用处理请求的快慢。通常平均响应时间无法反应服务的整体响应情况，响应慢的请求会被响应快的请求平均掉，而响应慢的请求往往会给用户带来糟糕的体验，即所谓的长尾请求，所以我们需要百分位响应时间，例如tp99响应时间，即99%的请求都会在这个时间内返回。
* **吞吐量**：主要指应用在一定时间内处理请求/事务的数量，反应的是应用的负载能力。我们当然希望在应用稳定的情况下，能承接的流量越大越好，主要指标包括QPS(每秒处理请求数)和QPM(每分钟处理请求数)。

#### 1.2 制定优化方案
明确了性能指标以后，我们就可以评估一个应用的性能好坏，同时也能发现其中的短板并对其进行优化。但是做性能优化，有几个点需要提前注意：

第一，不要反向优化。比如我们的应用整体占用内存资源较少，但是rt偏高，那我们就针对rt做优化，优化完后，rt下降了30%，但是cpu使用率上升了50%，导致一台机器负载能力下降30%，这便是反向优化。性能优化要从整体考虑，尽量在优化一个方面时，不影响其他方面，或是其他方面略微下降。

第二，不要过度优化。如果应用性能已经很好了，优化的空间很小，比如rt的tp99在2ms内，继续尝试优化可能投入产出比就很低了，不如将这些精力放在其他需要优化的地方上。

由此可见，在优化之前，明确想要优化的指标，并制定合理的优化方案是很重要的。

常见的优化方案有以下几种：
##### 1. 优化代码
有经验的程序员在编写代码时，会时刻注意减少代码中不必要的性能消耗，比如使用strconv而不是fmt.Sprint进行数字到字符串的转化，在初始化map或slice时指定合理的容量以减少内存分配等。良好的编程习惯不仅能使应用性能良好，同时也能减少故障发生的几率。总结下来，常用的代码优化方向有以下几种：

1）提高复用性，将通用的代码抽象出来，减少重复开发。

2）池化，对象可以池化，减少内存分配；协程可以池化，避免无限制创建协程打满内存。

3）并行化，在合理创建协程数量的前提下，把互不依赖的部分并行处理，减少整体的耗时。

4）异步化，把不需要关心实时结果的请求，用异步的方式处理，不用一直等待结果返回。

5）算法优化，使用时间复杂度更低的算法。

##### 2.使用设计模式
设计模式是对代码组织形式的抽象和总结，代码的结构对应用的性能有着重要的影响，结构清晰，层次分明的代码不仅可读性好，扩展性高，还能避免许多潜在的性能问题，帮助开发人员快速找到性能瓶颈，进行专项优化，为服务的稳定性提供保障。常见的对性能有所提升的设计模式例如单例模式，我们可以在应用启动时将需要的外部依赖服务用单例模式先初始化，避免创建太多重复的连接。

##### 3.空间换时间或时间换空间
在优化的前期，可能一个小的优化就能达到很好的效果。但是优化的尽头，往往要面临抉择，鱼和熊掌不可兼得。性能优秀的应用往往是多项资源的综合利用最优。为了达到综合平衡，在某些场景下，就需要做出一些调整和牺牲，常用的方法就是空间换时间或时间换空间。比如在响应时间优先的场景下，把需要耗费大量计算时间或是网络i/o时间的中间结果缓存起来，以提升后续相似请求的响应速度，便是空间换时间的一种体现。

#### 4. 使用更好的三方库
在我们的应用中往往会用到很多开源的第三方库，目前在github上的go开源项目就有173万+。有很多go官方库的性能表现并不佳，比如go官方的日志库性能就一般，下面是zap发布的基准测试信息（记录一条消息和10个字段的性能表现）。

![image](https://github.com/weifansym/workDoc/assets/6757408/17a6b54f-b2e8-4841-bdc7-10e9bea68ec1)

从上面可以看出zap的性能比同类结构化日志包更好，也比标准库更快，那我们就可以选择更好的三方库。
### 二、性能调优-工具篇
当我们找到应用的性能短板，并针对短板制定相应优化方案，最后按照方案对代码进行优化之后，我们怎么知道优化是有效的呢？直接将代码上线，观察性能指标的变化，风险太大了。此时我们需要有好用的性能分析工具，帮助我们检验优化的效果，下面将为大家介绍几款go语言中性能分析的利器。
#### 2.1 benchmark
Go语言标准库内置的 testing 测试框架提供了基准测试(benchmark)的能力，benchmark可以帮助我们评估代码的性能表现，主要方式是通过在一定时间(默认1秒)内重复运行测试代码，然后输出执行次数和内存分配结果。下面我们用一个简单的例子来验证一下，strconv是否真的比fmt.Sprint快。首先我们来编写一段基准测试的代码，如下：
```
package main

import (
    "fmt"
    "strconv"
    "testing"
)

func BenchmarkStrconv(b *testing.B) {
    for n := 0; n < b.N; n++ {
      strconv.Itoa(n)
  }
}

func BenchmarkFmtSprint(b *testing.B) {
    for n := 0; n < b.N; n++ {
      fmt.Sprint(n)
  }
}
```
我们可以用命令行**go test -bench**. 来运行基准测试，输出结果如下：
```
goos: darwin
goarch: amd64
pkg: main
cpu: Intel(R) Core(TM) i7-9750H CPU @ 2.60GHz
BenchmarkStrconv-12             41988014                27.41 ns/op
BenchmarkFmtSprint-12           13738172                81.19 ns/op
ok      main  7.039s
```
可以看到strconv每次执行只用了27.41纳秒，而fmt.Sprint则是81.19纳秒，strconv的性能是fmt.Sprint的三倍，那为什么strconv要更快呢？会不会是这次运行时间太短呢？为了公平起见，我们决定让他们再比赛一轮，这次我们延长比赛时间，看看结果如何。

通过**go test -bench . -benchtime=5s**命令，我们可以把测试时间延长到5秒，结果如下:
```
goos: darwin
goarch: amd64
pkg: main
cpu: Intel(R) Core(TM) i7-9750H CPU @ 2.60GHz
BenchmarkStrconv-12             211533207               31.60 ns/op
BenchmarkFmtSprint-12           69481287                89.58 ns/op
PASS
ok      main  18.891s
```
结果有些变化，strconv每次执行的时间上涨了4ns，但变化不大，差距仍有2.9倍。但是我们仍然不死心，我们决定让他们一次跑三轮，每轮5秒，三局两胜。

通过**go test -bench . -benchtime=5s -count=3**命令，我们可以把测试进行3轮，结果如下:
```
goos: darwin
goarch: amd64
pkg: main
cpu: Intel(R) Core(TM) i7-9750H CPU @ 2.60GHz
BenchmarkStrconv-12             217894554               31.76 ns/op
BenchmarkStrconv-12             217140132               31.45 ns/op
BenchmarkStrconv-12             219136828               31.79 ns/op
BenchmarkFmtSprint-12           70683580                89.53 ns/op
BenchmarkFmtSprint-12           63881758                82.51 ns/op
BenchmarkFmtSprint-12           64984329                82.04 ns/op
PASS
ok      main  54.296s
```
结果变化也不大，看来strconv是真的比fmt.Sprint快很多。那快是快，会不会内存分配上情况就相反呢？

通过**go test -bench . -benchmem**这个命令我们可以看到两个方法的内存分配情况，结果如下：
```
goos: darwin
goarch: amd64
pkg: main
cpu: Intel(R) Core(TM) i7-9750H CPU @ 2.60GHz
BenchmarkStrconv-12         43700922    27.46 ns/op    7 B/op   0 allocs/op
BenchmarkFmtSprint-12       143412      80.88 ns/op   16 B/op   2 allocs/op
PASS
ok      main  7.031s
```
可以看到strconv在内存分配上是0次，每次运行使用的内存是7字节，只是fmt.Sprint的43.8%，简直是全方面的优于fmt.Sprint啊。那究竟是为什么strconv比fmt.Sprint好这么多呢？

通过查看strconv的代码，我们发现，对于小于100的数字，strconv是直接通过digits和smallsString这两个常量进行转换的，而大于等于100的数字，则是通过不断除以100取余，然后再找到余数对应的字符串，把这些余数的结果拼起来进行转换的。
```
const digits = "0123456789abcdefghijklmnopqrstuvwxyz"
const smallsString = "00010203040506070809" +
  "10111213141516171819" +
  "20212223242526272829" +
  "30313233343536373839" +
  "40414243444546474849" +
  "50515253545556575859" +
  "60616263646566676869" +
  "70717273747576777879" +
  "80818283848586878889" +
  "90919293949596979899"
// small returns the string for an i with 0 <= i < nSmalls.
func small(i int) string {
  if i < 10 {
    return digits[i : i+1]
  }
  return smallsString[i*2 : i*2+2]
}
func formatBits(dst []byte, u uint64, base int, neg, append_ bool) (d []byte, s string) {
    ...
    for j := 4; j > 0; j-- {
        is := us % 100 * 2
        us /= 100
        i -= 2
        a[i+1] = smallsString[is+1]
        a[i+0] = smallsString[is+0]
    }
    ...
}
```
而fmt.Sprint则是通过反射来实现这一目的的，fmt.Sprint得先判断入参的类型，在知道参数是int型后，再调用fmt.fmtInteger方法把int转换成string，这多出来的步骤肯定没有直接把int转成string来的高效。
```
// fmtInteger formats signed and unsigned integers.
func (f *fmt) fmtInteger(u uint64, base int, isSigned bool, verb rune, digits string) {
    ...
    switch base {
  case 10:
    for u >= 10 {
      i--
      next := u / 10
      buf[i] = byte('0' + u - next*10)
      u = next
    }
    ...
}
```
benchmark还有很多实用的函数，比如ResetTimer可以重置启动时耗费的准备时间，StopTimer和StartTimer则可以暂停和启动计时，让测试结果更集中在核心逻辑上。
#### 2.2 pprof
##### 2.2.1 使用介绍
pprof是go语言官方提供的profile工具，支持可视化查看性能报告，功能十分强大。pprof基于定时器(10ms/次)对运行的go程序进行采样，搜集程序运行时的堆栈信息，包括CPU时间、内存分配等，最终生成性能报告。

pprof有两个标准库，使用的场景不同：
* runtime/pprof 通过在代码中显式的增加触发和结束埋点来收集指定代码块运行时数据生成性能报告。
* net/http/pprof 是对runtime/pprof的二次封装，基于web服务运行，通过访问链接触发，采集服务运行时的数据生成性能报告。

runtime/pprof的使用方法如下：
```
package main

import (
  "os"
  "runtime/pprof"
  "time"
)

func main() {
  w, _ := os.OpenFile("test_cpu", os.O_RDWR | os.O_CREATE | os.O_APPEND, 0644)
  pprof.StartCPUProfile(w)
  time.Sleep(time.Second)
  pprof.StopCPUProfile()
}
```
我们也可以使用另外一种方法，net/http/pprof：
```
package main

import (
    "net/http"
    _ "net/http/pprof"
)

func main() {
    err := http.ListenAndServe(":6060", nil)
    if err != nil {
        panic(err)
    }
}
```
将程序run起来后，我们通过访问http://127.0.0.1:6060/debug/pprof/就可以看到如下页面：
![image](https://github.com/weifansym/workDoc/assets/6757408/3252ca37-25a2-41b1-b902-da93fbd89b5c)

点击profile就可以下载cpu profile文件。那我们如何查看我们的性能报告呢？

pprof支持两种查看模式，终端和web界面，注意: **想要查看可视化界面需要提前安装graphviz。**

这里我们以web界面为例，在终端内我们输入如下命令：
> go tool pprof -http :6060 test_cpu

就会在浏览器里打开一个页面，内容如下：

![image](https://github.com/weifansym/workDoc/assets/6757408/90eebe14-0596-413b-a85e-1e9b2b1a6f75)

从界面左上方VIEW栏下，我们可以看到，pprof支持Flame Graph，dot Graph和Top等多种视图，下面我们将一一介绍如何阅读这些视图。

##### 2.2.1 火焰图 Flame Graph如何阅读
首先，**推荐直接阅读火焰图**，在查函数耗时场景，这个比较直观；

最简单的：**横条越长，资源消耗、占用越多**； 

注意：每一个function 的横条虽然很长，但可能是他的下层“子调用”耗时产生的，所以一定要关注“下一层子调用”各自的耗时分布；

每个横条支持点击下钻能力，可以更详细的分析子层的耗时占比。
![image](https://github.com/weifansym/workDoc/assets/6757408/a342a9c8-0e25-42f2-8a1e-2399643ad167)

##### 2.2.2 dot Graph 图如何阅读

英文原文在这里：https://github.com/google/pprof/blob/master/doc/README.md

![image](https://github.com/weifansym/workDoc/assets/6757408/b6031236-102a-44c5-a714-3f53d6de259e)

* **节点颜色**:
  * 红色表示耗时多的节点；
  * 绿色表示耗时少的节点；
  * 灰色表示耗时几乎可以忽略不计（接近零）；

* **节点字体大小**:
  * 字体越大，表示占“上层函数调用”比例越大；（其实上层函数自身也有耗时，没包含在此）
  * 字体越小，表示占“上层函数调用”比例越小；

* **线条（边）粗细**:
  * 线条越粗，表示消耗了更多的资源；
  * 反之，则越少；

* **线条（边）颜色**:
  * 颜色越红，表示性能消耗占比越高；
  * 颜色越绿，表示性能消耗占比越低；
  * 灰色，表示性能消耗几乎可以忽略不计；

* 虚线：表示中间有一些节点被“移除”或者忽略了；(一般是因为耗时较少所以忽略了) 

* 实线：表示节点之间直接调用

* 内联边标记：被调用函数已经被内联到调用函数中（对于一些代码行比较少的函数，编译器倾向于将它们在编译期展开从而消除函数调用，这种行为就是内联。）

##### 2.2.3 TOP 表如何阅读
* flat：当前函数，运行耗时（不包含内部调用其他函数的耗时）
* flat%：当前函数，占用的 CPU 运行耗时总比例（不包含外部调用函数）
* sum%：当前行的 flat% 与上面所有行的flat%总和。
* cum：当前函数加上它内部的调用的运行总耗时（包含内部调用其他函数的耗时）
* cum%：同上的 CPU 运行耗时总比例

![image](https://github.com/weifansym/workDoc/assets/6757408/e1b56b94-f1ae-448c-8e32-b087f25ba978)

#### 2.3 trace
pprof已经有了对内存和CPU的分析能力，那trace工具有什么不同呢？虽然pprof的CPU分析器，可以告诉你什么函数占用了最多的CPU时间，但它并不能帮助你定位到是什么阻止了goroutine运行，或者在可用的OS线程上如何调度goroutines。这正是trace真正起作用的地方。

我们需要更多关于Go应用中各个goroutine的执行情况的更为详细的信息，可以从P（goroutine调度器概念中的processor)和G（goroutine调度器概念中的goroutine）的视角完整的看到每个P和每个G在Tracer开启期间的全部“所作所为”，对Tracer输出数据中的每个P和G的行为分析并结合详细的event数据来辅助问题诊断的。

Tracer可以帮助我们记录的详细事件包含有：
* 与goroutine调度有关的事件信息：goroutine的创建、启动和结束；goroutine在同步原语（包括mutex、channel收发操作）上的阻塞与解锁。
* 与网络有关的事件：goroutine在网络I/O上的阻塞和解锁；
* 与系统调用有关的事件：goroutine进入系统调用与从系统调用返回；
* 与垃圾回收器有关的事件：GC的开始/停止，并发标记、清扫的开始/停止。

Tracer主要也是用于辅助诊断这三个场景下的具体问题的：
* 并行执行程度不足的问题：比如没有充分利用多核资源等；
* 因GC导致的延迟较大的问题；
* Goroutine执行情况分析，尝试发现goroutine因各种阻塞（锁竞争、系统调用、调度、辅助GC）而导致的有效运行时间较短或延迟的问题。

##### 2.3.1 trace性能报告
打开trace性能报告，首页信息包含了多维度数据，如下图：
![image](https://github.com/weifansym/workDoc/assets/6757408/4fb65267-8759-4a68-ae77-f9071690baa7)

* View trace：以图形页面的形式渲染和展示tracer的数据，这也是我们最为关注/最常用的功能
* Goroutine analysis：以表的形式记录执行同一个函数的多个goroutine的各项trace数据
* Network blocking profile：用pprof profile形式的调用关系图展示网络I/O阻塞的情况
* Synchronization blocking profile：用pprof profile形式的调用关系图展示同步阻塞耗时情况
* Syscall blocking profile：用pprof profile形式的调用关系图展示系统调用阻塞耗时情况
* Scheduler latency profile：用pprof profile形式的调用关系图展示调度器延迟情况
* User-defined tasks和User-defined regions：用户自定义trace的task和region
* Minimum mutator utilization：分析GC对应用延迟和吞吐影响情况的曲线图

通常我们最为关注的是View trace和Goroutine analysis，下面将详细说说这两项的用法。
##### 2.3.2 view trace
如果Tracer跟踪时间较长,trace会将View trace按时间段进行划分，避免触碰到trace-viewer的限制：

![image](https://github.com/weifansym/workDoc/assets/6757408/709c6df8-e6ae-42bb-8463-5e49d571eb96)

![image](https://github.com/weifansym/workDoc/assets/6757408/1892d813-7c4e-4055-9fc3-3b995e1b8f10)

View trace使用快捷键来缩放时间线标尺：w键用于放大（从秒向纳秒缩放），s键用于缩小标尺（从纳秒向秒缩放）。我们同样可以通过快捷键在时间线上左右移动：s键用于左移，d键用于右移。(游戏快捷键WASD)

##### 采样状态
这个区内展示了三个指标：Goroutines、Heap和Threads，某个时间点上的这三个指标的数据是这个时间点上的状态快照采样：Goroutines：某一时间点上应用中启动的goroutine的数量，当我们点击某个时间点上的goroutines采样状态区域时（我们可以用快捷键m来准确标记出那个时间点），事件详情区会显示当前的goroutines指标采样状态：
![image](https://github.com/weifansym/workDoc/assets/6757408/6b0d0277-c4f8-4f6e-8903-d91e8b53dd5d)

Heap指标则显示了某个时间点上Go应用heap分配情况（包括已经分配的Allocated和下一次GC的目标值NextGC）：
![image](https://github.com/weifansym/workDoc/assets/6757408/683fac2d-10be-4d64-8556-17f3245cfd54)

Threads指标显示了某个时间点上Go应用启动的线程数量情况，事件详情区将显示处于InSyscall（整阻塞在系统调用上）和Running两个状态的线程数量情况：
![image](https://github.com/weifansym/workDoc/assets/6757408/24ddab7e-d6ed-4dfa-b706-dd5d3fee093e)

##### P视角区
这里将View trace视图中最大的一块区域称为“P视角区”。这是因为在这个区域，我们能看到Go应用中每个P（Goroutine调度概念中的P）上发生的所有事件，包括：EventProcStart、EventProcStop、EventGoStart、EventGoStop、EventGoPreempt、Goroutine辅助GC的各种事件以及Goroutine的GC阻塞(STW)、系统调用阻塞、网络阻塞以及同步原语阻塞(mutex)等事件。除了每个P上发生的事件，我们还可以看到以单独行显示的GC过程中的所有事件。
![image](https://github.com/weifansym/workDoc/assets/6757408/4262b2e2-ec51-430e-8b64-2bca34588819)

##### 事件详情区
点选某个事件后，关于该事件的详细信息便会在这个区域显示出来，事件详情区可以看到关于该事件的详细信息：
![image](https://github.com/weifansym/workDoc/assets/6757408/6a32ee90-17e0-4729-83d7-d920d5bbf786)

* Title：事件的可读名称；
* Start：事件的开始时间，相对于时间线上的起始时间；
* Wall Duration：这个事件的持续时间，这里表示的是G1在P4上此次持续执行的时间；
* Start Stack Trace：当P4开始执行G1时G1的调用栈；
* End Stack Trace：当P4结束执行G1时G1的调用栈；从上面End Stack Trace栈顶的函数为runtime.asyncPreempt来看，该Goroutine G1是被强行抢占了，这样P4才结束了其运行；
* Incoming flow：触发P4执行G1的事件；
* Outgoing flow：触发G1结束在P4上执行的事件；
* Preceding events：与G1这个goroutine相关的之前的所有的事件；
* Follwing events：与G1这个goroutine相关的之后的所有的事件
* All connected：与G1这个goroutine相关的所有事件。

##### 2.3.3 Goroutine analysis
Goroutine analysis提供了从G视角看Go应用执行的图景。与View trace不同，这次页面中最广阔的区域提供的G视角视图，而不再是P视角视图。在这个视图中，每个G都会对应一个单独的条带（和P视角视图一样，每个条带都有两行），通过这一条带可以按时间线看到这个G的全部执行情况。通常仅需在goroutine analysis的表格页面找出执行最快和最慢的两个goroutine，在Go视角视图中沿着时间线对它们进行对比，以试图找出执行慢的goroutine究竟出了什么问题。
![image](https://github.com/weifansym/workDoc/assets/6757408/d3887c2d-b246-48ab-8ea5-494ac70a0374)

#### 2.4 后记
虽然pprof和trace有着非常强大的profile能力，但在使用过程中，仍存在以下痛点：
* 获取性能报告麻烦：一般大家做压测，为了更接近真实环境性能态，都使用生产环境/pre环境进行。而出于安全考虑，生产环境内网一般和PC办公内网是隔离不通的，需要单独配置通路才可以获得生产环境内网的profile 文件下载到PC办公电脑中，这也有一些额外的成本；

* 查看profile分析报告麻烦：之前大家在本地查看profile 分析报告，一般 go tool pprof -http=":8083" profile 命令在本地PC开启一个web service 查看，并且需要至少安装graphviz 等库。

* 查看trace分析同样麻烦：查看go trace 的profile 信息来分析routine 锁和生命周期时，也需要类似的方式在本地PC执行命令 go tool trace mytrace.profile 。

* 分享麻烦：如果我想把自己压测的性能结果内容，分享个另一位同学，那只能把1中获取的性能报告“profile文件”通过钉钉发给被分享人。然而有时候本地profile文件比较多，一不小心就发错了，还不如截图，但是截图又没有了交互放大、缩小、下钻等能力。处处不给力！

* 留存复盘麻烦：系统的性能分析就像一份病历，每每看到阶段性的压测报告，总结或者对照时，不禁要询问，做过了哪些优化和改造，病因病灶是什么，有没有共性，值不值得总结归纳，现在是不是又面临相似的性能问题？

那么能不能开发一个平台工具，解决以上的这些痛点呢？目前在阿里集团内部，高德的研发同学已经通过对go官方库的定制开发，实现了go语言性能平台，解决了以上这些痛点，并在内部进行了开源。该平台已面向阿里集团，累计实现性能场景快照数万条的获取和分析，解决了很多的线上服务性能调试和优化问题，这里暂时不展开，后续有机会可以单独分享。
### 三、性能调优-技巧篇
除了前面提到的尽量用strconv而不是fmt.Sprint进行数字到字符串的转化以外，我们还将介绍一些在实际开发中经常会用到的技巧，供各位参考。
####3.1 字符串拼接
拼接字符串为了书写方便快捷，最常用的两个方法是运算符 + 和 fmt.Sprintf()

运算符 + 只能简单地完成字符串之间的拼接，fmt.Sprintf() 其底层实现使用了反射，性能上会有所损耗。

从性能出发，兼顾易用可读，如果待拼接的变量不涉及类型转换且数量较少（<=5），拼接字符串推荐使用运算符 +，反之使用 fmt.Sprintf()。
```
// 推荐：用+进行字符串拼接
func BenchmarkPlus(b *testing.B) {
  for i := 0; i < b.N; i++ {
    s := "a" + "b"
    _ = s
  }
}
// 不推荐：用fmt.Sprintf进行字符串拼接
func BenchmarkFmt(b *testing.B) {
  for i := 0; i < b.N; i++ {
    s := fmt.Sprintf("%s%s", "a", "b")
    _ = s
  }
}

goos: darwin
goarch: amd64
pkg: main
cpu: Intel(R) Core(TM) i7-9750H CPU @ 2.60GHz
BenchmarkPlus-12        1000000000               0.2658 ns/op          0 B/op          0 allocs/op
BenchmarkFmt-12         16559949                70.83 ns/op            2 B/op          1 allocs/op
PASS
ok      main  5.908s
```
#### 3.2 提前指定容器容量
在初始化slice时，尽量指定容量，这是因为当添加元素时，如果容量的不足，slice会重新申请一个更大容量的容器，然后把原来的元素复制到新的容器中。
```
// 推荐：初始化时指定容量
func BenchmarkGenerateWithCap(b *testing.B) {

  nums := make([]int, 0, 10000)
  for n := 0; n < b.N; n++ {
    for i:=0; i < 10000; i++ {
      nums = append(nums, i)
    }
  }
}
// 不推荐：初始化时不指定容量
func BenchmarkGenerate(b *testing.B) {
  nums := make([]int, 0)
  for n := 0; n < b.N; n++ {
    for i:=0; i < 10000; i++ {
      nums = append(nums, i)
    }
  }
}

goos: darwin
goarch: amd64
pkg: main
cpu: Intel(R) Core(TM) i7-9750H CPU @ 2.60GHz
BenchmarkGenerateWithCap-12        23508            336485 ns/op          476667 B/op          0 allocs/op
BenchmarkGenerate-12               22620             68747 ns/op          426141 B/op          0 allocs/op
PASS
ok      main  16.628s
```
#### 3.3 遍历 []struct{} 使用下标而不是 range
常用的遍历方式有两种，一种是for循环下标遍历，一种是for循环range遍历，这两种遍历在性能上是否有差异呢？让我们来一探究竟。

针对[]int，我们来看看两种遍历有和差别吧
```
func getIntSlice() []int {
  nums := make([]int, 1024, 1024)
  for i := 0; i < 1024; i++ {
    nums[i] = i
  }
  return nums
}
// 用下标遍历[]int
func BenchmarkIndexIntSlice(b *testing.B) {
  nums := getIntSlice()
  b.ResetTimer()
  for i := 0; i < b.N; i++ {
    var tmp int
    for k := 0; k < len(nums); k++ {
      tmp = nums[k]
    }
    _ = tmp
  }
}
// 用range遍历[]int元素
func BenchmarkRangeIntSlice(b *testing.B) {
  nums := getIntSlice()
  b.ResetTimer()
  for i := 0; i < b.N; i++ {
    var tmp int
    for _, num := range nums {
      tmp = num
    }
    _ = tmp
  }
}

goos: darwin
goarch: amd64
pkg: demo/test
cpu: Intel(R) Core(TM) i7-9750H CPU @ 2.60GHz
BenchmarkIndexIntSlice-12        3923230               270.2 ns/op             0 B/op          0 allocs/op
BenchmarkRangeIntSlice-12        4518495               287.8 ns/op             0 B/op          0 allocs/op
PASS
ok      demo/test       3.303s
```
可以看到，在遍历[]int时，两种方式并无差别。

我们再看看遍历[]struct{}的情况
```
type Item struct {
  id  int
  val [1024]byte
}
// 推荐：用下标遍历[]struct{}
func BenchmarkIndexStructSlice(b *testing.B) {
  var items [1024]Item
  for i := 0; i < b.N; i++ {
    var tmp int
    for j := 0; j < len(items); j++ {
      tmp = items[j].id
    }
    _ = tmp
  }
}
// 推荐：用range的下标遍历[]struct{}
func BenchmarkRangeIndexStructSlice(b *testing.B) {
  var items [1024]Item
  for i := 0; i < b.N; i++ {
    var tmp int
    for k := range items {
      tmp = items[k].id
    }
    _ = tmp
  }
}
// 不推荐：用range遍历[]struct{}的元素
func BenchmarkRangeStructSlice(b *testing.B) {
  var items [1024]Item
  for i := 0; i < b.N; i++ {
    var tmp int
    for _, item := range items {
      tmp = item.id
    }
    _ = tmp
  }
}

goos: darwin
goarch: amd64
pkg: demo/test
cpu: Intel(R) Core(TM) i7-9750H CPU @ 2.60GHz
BenchmarkIndexStructSlice-12             4413182               266.7 ns/op             0 B/op          0 allocs/op
BenchmarkRangeIndexStructSlice-12        4545476               269.4 ns/op             0 B/op          0 allocs/op
BenchmarkRangeStructSlice-12               33300             35444 ns/op               0 B/op          0 allocs/op
PASS
ok      demo/test       5.282s
```
可以看到，用for循环下标的方式性能都差不多，但是用range遍历数组里的元素时，性能则相差很多，前面两种方法是第三种方法的130多倍。主要原因是通过for k, v := range获取到的元素v实际上是原始值的一个拷贝。所以在面对复杂的struct进行遍历的时候，推荐使用下标。但是当遍历对象是复杂结构体的指针([]*struct{})时，用下标还是用range迭代元素的性能就差不多了。

####3.4 利用unsafe包避开内存copy
unsafe包提供了任何类型的指针和 unsafe.Pointer 的相互转换及uintptr 类型和 unsafe.Pointer 可以相互转换，如下图
![image](https://github.com/weifansym/workDoc/assets/6757408/f576accd-263b-4d37-90b7-3d26da7de029)

依据上述转换关系，其实除了string和[]byte的转换，也可以用于slice、map等的求长度及一些结构体的偏移量获取等，但是这种黑科技在一些情况下会带来一些匪夷所思的诡异问题，官方也不建议用，所以还是慎用，除非你确实很理解各种机制了，这里给出项目中实际用到的常规string和[]byte之间的转换，如下：
```
func Str2bytes(s string) []byte {
   x := (*[2]uintptr)(unsafe.Pointer(&s))
   h := [3]uintptr{x[0], x[1], x[1]}
   return *(*[]byte)(unsafe.Pointer(&h))
}

func Bytes2str(b []byte) string {
   return *(*string)(unsafe.Pointer(&b))
}
```
我们通过benchmark来验证一下是否性能更优：
```
/ 推荐：用unsafe.Pointer实现string到bytes
func BenchmarkStr2bytes(b *testing.B) {
  s := "testString"
  var bs []byte
  for n := 0; n < b.N; n++ {
    bs = Str2bytes(s)
  }
  _ = bs
}
// 不推荐：用类型转换实现string到bytes
func BenchmarkStr2bytes2(b *testing.B) {
  s := "testString"
  var bs []byte
  for n := 0; n < b.N; n++ {
    bs = []byte(s)
  }
  _ = bs
}
// 推荐：用unsafe.Pointer实现bytes到string
func BenchmarkBytes2str(b *testing.B) {
  bs := Str2bytes("testString")
  var s string
  b.ResetTimer()
  for n := 0; n < b.N; n++ {
    s = Bytes2str(bs)
  }
  _ = s
}
// 不推荐：用类型转换实现bytes到string
func BenchmarkBytes2str2(b *testing.B) {
  bs := Str2bytes("testString")
  var s string
  b.ResetTimer()
  for n := 0; n < b.N; n++ {
    s = string(bs)
  }
  _ = s
}

goos: darwin
goarch: amd64
pkg: demo/test
cpu: Intel(R) Core(TM) i7-9750H CPU @ 2.60GHz
BenchmarkStr2bytes-12           1000000000               0.2938 ns/op          0 B/op          0 allocs/op
BenchmarkStr2bytes2-12          38193139                28.39 ns/op           16 B/op          1 allocs/op
BenchmarkBytes2str-12           1000000000               0.2552 ns/op          0 B/op          0 allocs/op
BenchmarkBytes2str2-12          60836140                19.60 ns/op           16 B/op          1 allocs/op
PASS
ok      demo/test       3.301s
```
可以看到使用unsafe.Pointer比强制类型转换性能是要高不少的，从内存分配上也可以看到完全没有新的内存被分配。
#### 3.5 协程池
go语言最大的特色就是很容易的创建协程，同时go语言的协程调度策略也让go程序可以最大化的利用cpu资源，减少线程切换。但是无限度的创建goroutine，仍然会带来问题。我们知道，一个go协程占用内存大小在2KB左右，无限度的创建协程除了会占用大量的内存空间，同时协程的切换也有不少开销，一次协程切换大概需要100ns，虽然相较于线程毫秒级的切换要优秀很多，但依然存在开销，而且这些协程最后还是需要GC来回收，过多的创建协程，对GC也是很大的压力。所以我们在使用协程时，可以通过协程池来限制goroutine数量，避免无限制的增长。

限制协程的方式有很多，比如可以用channel来限制：
```
var wg sync.WaitGroup
ch := make(chan struct{}, 3)
for i := 0; i < 10; i++ {
    ch <- struct{}{}
  wg.Add(1)
  go func(i int) {
        defer wg.Done()
        log.Println(i)
        time.Sleep(time.Second)
        <-ch
    }(i)
}
wg.Wait()
```
这里通过限制channel长度为3，可以实现最多只有3个协程被创建的效果。

当然也可以使用errgoup。使用方法如下：
```
func Test_ErrGroupRun(t *testing.T) {
  errgroup := WithTimeout(nil, 10*time.Second)
  errgroup.SetMaxProcs(4)
  for index := 0; index < 10; index++ {
    errgroup.Run(nil, index, "test", func(context *gin.Context, i interface{}) (interface{},
      error) {
      t.Logf("[%s]input:%+v, time:%s", "test", i, time.Now().Format("2006-01-02 15:04:05"))
      time.Sleep(2*time.Second)
      return i, nil
    })
  }
  errgroup.Wait()
}
```
输出结果如下：
```
=== RUN   Test_ErrGroupRun
    errgroup_test.go:23: [test]input:0, time:2022-12-04 17:31:29
    errgroup_test.go:23: [test]input:3, time:2022-12-04 17:31:29
    errgroup_test.go:23: [test]input:1, time:2022-12-04 17:31:29
    errgroup_test.go:23: [test]input:2, time:2022-12-04 17:31:29
    errgroup_test.go:23: [test]input:4, time:2022-12-04 17:31:31
    errgroup_test.go:23: [test]input:5, time:2022-12-04 17:31:31
    errgroup_test.go:23: [test]input:6, time:2022-12-04 17:31:31
    errgroup_test.go:23: [test]input:7, time:2022-12-04 17:31:31
    errgroup_test.go:23: [test]input:8, time:2022-12-04 17:31:33
    errgroup_test.go:23: [test]input:9, time:2022-12-04 17:31:33
--- PASS: Test_ErrGroupRun (6.00s)
PASS
```
errgroup可以通过SetMaxProcs设定协程池的大小，从上面的结果可以看到，最多就4个协程在运行。

#### 3.6 sync.Pool 对象复用
我们在代码中经常会用到json进行序列化和反序列化，举一个投放活动的例子，一个投放活动会有许多字段会转换为字节数组。
```
type ActTask struct {
  Id                 int64                `ddb:"id"`             // 主键id
  Status             common.Status        `ddb:"status"`         // 状态 0=初始 1=生效 2=失效 3=过期
  BizProd            common.BizProd       `ddb:"biz_prod"`       // 业务类型
  Name               string               `ddb:"name"`            // 活动名
  Adcode             string               `ddb:"adcode"`         // 城市
  RealTimeRuleByte   []byte               `ddb:"realtime_rule"`  // 实时规则json
  ...
}

type RealTimeRuleStruct struct {
  Filter []*struct {
    PropertyId   int64    `json:"property_id"`
    PropertyCode string   `json:"property_code"`
    Operator     string   `json:"operator"`
    Value        []string `json:"value"`
  } `json:"filter"`
  ExtData [1024]byte `json:"ext_data"`
}

func (at *ActTask) RealTimeRule() *form.RealTimeRule {
  if err := json.Unmarshal(at.RealTimeRuleByte, &at.RealTimeRuleStruct); err != nil {
    return nil
  }
  return at.RealTimeRuleStruct
}
```
以这里的实时投放规则为例，我们会将过滤规则反序列化为字节数组。每次json.Unmarshal都会申请一个临时的结构体对象，而这些对象都是分配在堆上的，会给 GC 造成很大压力，严重影响程序的性能。

对于需要频繁创建并回收的对象，我们可以使用对象池来提升性能。sync.Pool可以将暂时不用的对象缓存起来，待下次需要的时候直接使用，不用再次经过内存分配，复用对象的内存，减轻 GC 的压力，提升系统的性能。

sync.Pool的使用方法很简单，只需要实现 New 函数即可。对象池中没有对象时，将会调用 New 函数创建。
```
var realTimeRulePool = sync.Pool{
    New: func() interface{} { 
        return new(RealTimeRuleStruct) 
    },
}
```
然后调用 Pool 的 Get() 和 Put() 方法来获取和放回池子中。
```
rule := realTimeRulePool.Get().(*RealTimeRuleStruct)
json.Unmarshal(buf, rule)
realTimeRulePool.Put(rule)
```
* Get() 用于从对象池中获取对象，因为返回值是 interface{}，因此需要类型转换。
* Put() 则是在对象使用完毕后，放回到对象池。

接下来我们进行性能测试，看看性能如何：
```
var realTimeRule = []byte("{\\\"filter\\\":[{\\\"property_id\\\":2,\\\"property_code\\\":\\\"search_poiid_industry\\\",\\\"operator\\\":\\\"in\\\",\\\"value\\\":[\\\"yimei\\\"]},{\\\"property_id\\\":4,\\\"property_code\\\":\\\"request_page_id\\\",\\\"operator\\\":\\\"in\\\",\\\"value\\\":[\\\"all\\\"]}],\\\"white_list\\\":[{\\\"property_id\\\":1,\\\"property_code\\\":\\\"white_list_for_adiu\\\",\\\"operator\\\":\\\"in\\\",\\\"value\\\":[\\\"j838ef77bf227chcl89888f3fb0946\\\",\\\"lb89bea9af558589i55559764bc83e\\\"]}],\\\"ipc_user_tag\\\":[{\\\"property_id\\\":1,\\\"property_code\\\":\\\"ipc_crowd_tag\\\",\\\"operator\\\":\\\"in\\\",\\\"value\\\":[\\\"test_20227041152_mix_ipc_tag\\\"]}],\\\"relation_id\\\":0,\\\"is_copy\\\":true}")
// 推荐：复用一个对象，不用每次都生成新的
func BenchmarkUnmarshalWithPool(b *testing.B) {
  for n := 0; n < b.N; n++ {
    task := realTimeRulePool.Get().(*RealTimeRuleStruct)
    json.Unmarshal(realTimeRule, task)
    realTimeRulePool.Put(task)
  }
}
// 不推荐：每次都会生成一个新的临时对象
func BenchmarkUnmarshal(b *testing.B) {
  for n := 0; n < b.N; n++ {
    task := &RealTimeRuleStruct{}
    json.Unmarshal(realTimeRule, task)
  }
}

goos: darwin
goarch: amd64
pkg: demo/test
cpu: Intel(R) Core(TM) i7-9750H CPU @ 2.60GHz
BenchmarkUnmarshalWithPool-12    3627546     319.4 ns/op   312 B/op  7 allocs/op
BenchmarkUnmarshal-12            2342208     490.8 ns/op  1464 B/op  8 allocs/op
PASS
ok      demo/test       3.525s
```
可以看到，两种方法在时间消耗上差不太多，但是在内存分配上差距明显，使用sync.Pool后内存占用仅为不使用的1/5。
#### 3.7 避免系统调用
系统调用是一个很耗时的操作，在各种语言中都是，go也不例外，在go的GPM模型中，异步系统调用G会和MP分离，同步系统调用GM会和P分离，不管何种形式除了状态切换及内核态中执行操作耗时外，调度器本身的调度也耗时。所以在可以避免系统调用的地方尽量去避免。
```
// 推荐：不使用系统调用
func BenchmarkNoSytemcall(b *testing.B) {
   b.RunParallel(func(pb *testing.PB) {
      for pb.Next() {
         if configs.PUBLIC_KEY != nil {
         }
      }
   })
}
// 不推荐：使用系统调用
func BenchmarkSytemcall(b *testing.B) {
   b.RunParallel(func(pb *testing.PB) {
      for pb.Next() {
         if os.Getenv("PUBLIC_KEY") != "" {
         }
      }
   })
}
goos: darwin
goarch: amd64
pkg: demo/test
cpu: Intel(R) Core(TM) i7-9750H CPU @ 2.60GHz
BenchmarkNoSytemcall-12         1000000000              0.1495 ns/op          0 B/op          0 allocs/op
BenchmarkSytemcall-12           37224988                31.10 ns/op           0 B/op          0 allocs/op
PASS
ok      demo/test       1.877s
```
### 四、性能调优-实战篇
##### 案例1: go协程创建数据库连接不释放导致内存暴涨
应用背景

感谢@路现提供的案例。

遇到的问题及表象特征

线上机器偶尔出现内存使用率超过百分之九十报警。

分析思路及排查方向

在报警触发时，通过直接拉取线上应用的profile文件，查看内存分配情况，我们看到内存分配主要产生在本地缓存的组件上。
![image](https://github.com/weifansym/workDoc/assets/6757408/91053943-9c2c-47ab-8fb2-e72493ee42f2)

但是分析代码并没有发现存在内存泄露的情况，看着像是资源一直没有被释放，进一步分析goroutine的profile文件。
![image](https://github.com/weifansym/workDoc/assets/6757408/e64dcf59-b03f-4a2c-92ba-3035e7718aa8)

发现存在大量的goroutine未释放，表现在本地缓存击穿后回源数据库，对数据库的查询访问一直不释放。

调优手段与效果

最终通过排查，发现使用的数据库组件存在bug，在极端情况下会出现死锁的情况，导致数据库访问请求无法返回也无法释放。最终bug修复后升级数据库组件版本解决了问题。

##### 案例2: 优惠索引内存分配大，gc 耗时高
应用背景

感谢@梅东提供的案例。

遇到的问题及表象特征

接口tp99高，偶尔会有一些特别耗时的请求，导致用户的优惠信息展示不出来。

分析思路及排查方向

通过直接在平台上抓包观察，我们发现使用的分配索引这个方法占用的堆内存特别高，通过 top 可以看到是排在第一位的。
![image](https://github.com/weifansym/workDoc/assets/6757408/5160a249-5ae3-46b2-8120-2e9875936bcd)

![image](https://github.com/weifansym/workDoc/assets/6757408/8cebf3c5-37d9-4e5a-b84a-2c09ef5074bc)
我们分析代码，可以看到，获取城市索引的地方，每次都是重新申请了内存的，通过改动为返回指针，就不需要每次都单独申请内存了，核心代码改动：
![image](https://github.com/weifansym/workDoc/assets/6757408/6d6b690d-0d2e-4cfe-bc7f-9f29aedf9c4c)

调优手段与效果
修改后，上线观察，可以看到使用中的内存以及gc耗时都有了明显降低
![image](https://github.com/weifansym/workDoc/assets/6757408/843f0c02-9fed-4380-89ec-42f2af0c7451)

##### 案例3：流量上涨导致cpu异常飙升
应用背景

感谢@君度提供的案例。

遇到的问题及表象特征

能量站v2接口和task-home-page接口流量较大时，会造成ab实验策略匹配时cpu飙升

分析思路及排查方向
![image](https://github.com/weifansym/workDoc/assets/6757408/b0a42463-604b-4ce6-8b19-1be75c63ecb7)

调优手段与效果

主要优化点如下：

1、优化toEntity方法，简化为单独的ID()方法

2、优化数组、map初始化结构

3、优化adCode转换为string过程

4、关闭过多的match log打印

优化后profile：

![image](https://github.com/weifansym/workDoc/assets/6757408/6d2a2cda-b1f5-451c-b313-edd5197fe81a)

优化上线前后CPU的对比

![image](https://github.com/weifansym/workDoc/assets/6757408/bf459ecb-aa1d-473c-8e9a-3e81aa659b07)

##### 案例4：内存对象未释放导致内存泄漏
应用背景

感谢@淳深提供的案例，提供案例的服务，日常流量峰值在百万qps左右，是高德内部十分重要的服务。此前该服务是由java实现的，后来用go语言进行重构，在重构完成切全量后，有许多性能优化的优秀案例，这里选取内存泄漏相关的一个案例分享给大家，希望对大家在自己服务进行内存泄漏问题排查时能提供参考和帮助。

遇到的问题及表象特征

go语言版本全量切流后，每天会对服务各项指标进行详细review，发现每日内存上涨约0.4%，如下图
![image](https://github.com/weifansym/workDoc/assets/6757408/c49fde6f-9eb1-4534-8476-6cca9683f06a)

![image](https://github.com/weifansym/workDoc/assets/6757408/97b3c84b-2a32-4ec9-aef4-076acea1ad22)

在go版本服务切全量前，从第一张图可以看到整个内存使用是平稳的，无上涨趋势，但是切go版本后，从第二张图可以看到，整个内存曲线呈上升趋势，遂认定内存泄漏，开始排查内存泄漏的“罪魁祸首”。

分析思路及排查方向

我们先到线上机器抓取当前时间的heap文件，间隔一天后再次抓取heap文件，通过pprof diff对比，我们发现time.NewTicker的内存占用增长了几十MB(由于未保留当时的heap文件，此处没有截图)，通过调用栈信息，我们找到了问题的源头，来自中间件vipserver client的SrvHost方法，通过深扒vipserver client代码，我们发现，每个vipserver域名都会有一个对应的协程，这个协程每隔三秒钟就会新建一个ticker对象，且用过的ticker对象没有stop，也就不会释放相应的内存资源。
![image](https://github.com/weifansym/workDoc/assets/6757408/be39b0c6-ac2b-4aaa-84a5-b3a94e3ab2e6)

而这个time.NewTicker会创建一个timer对象，这个对象会占用72字节内存。
![image](https://github.com/weifansym/workDoc/assets/6757408/abe97606-4e31-45b3-9f9d-de923fac9833)

在服务运行一天的情况下，进过计算，该对象累计会在内存中占用约35.6MB，和上述内存每日增长0.4%正好能对上，我们就能断定这个内存泄漏来自这里。

调优手段与效果

知道是timer对象重复创建的问题后，只需要修改这部分的代码就好了，最新的vipserver client修改了此处的逻辑，如下
![image](https://github.com/weifansym/workDoc/assets/6757408/5d89e82d-2c50-4158-8ee4-396bad5c076d)

修改完后，运行一段时间，内存运行状态平稳，已无内存泄漏问题。

![image](https://github.com/weifansym/workDoc/assets/6757408/6b5526bb-40d8-41db-9498-3bd1f5f0c032)

### 结语

目前go语言不仅在阿里集团内部，在整个互联网行业内也越来越流行，希望本文能为正在使用go语言的同学在性能优化方面带来一些参考价值。在阿里集团内部，高德也是最早规模化使用go语言的团队之一，目前高德线上运行的go服务已经达到近百个，整体qps已突破百万量级。在使用go语言的同时，高德也为集团内go语言生态建设做出了许多贡献，包括开发支持阿里集团常见的中间件（比如配置中心-Diamond、分布式RPC服务框架-HSF、服务发现-Vipserver、消息队列-MetaQ、流量控制-Sentinel、日志追踪-Eagleeye等）go语言版本，并被阿里中间件团队官方收录。但是go语言生态建设仍然有很长的道路要走，希望能有更多对go感兴趣的同学能够加入我们，一起参与阿里的go生态建设，乃至为互联网业界的go生态发展添砖加瓦。

转自：https://mp.weixin.qq.com/s/UHaCLhiIyLYVrba-nEUONA






