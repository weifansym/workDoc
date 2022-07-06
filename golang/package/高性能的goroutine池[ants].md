## 高性能的goroutine池[ants]
### 介绍
处理大量并发是 Go 语言的一大优势。语言内置了方便的并发语法，可以非常方便的创建很多个轻量级的 goroutine 并发处理任务。相比于创建多个线程，goroutine 更轻量、资源占用更少、
切换速度更快、无线程上下文切换开销更少。但是受限于资源总量，系统中能够创建的 goroutine 数量也是受限的。默认每个 goroutine 占用 8KB 内存，一台 8GB 内存的机器满打满算也
只能创建 8GB/8KB = 1000000 个 goroutine，更何况系统还需要保留一部分内存运行日常管理任务，go 运行时需要内存运行 gc、处理 goroutine 切换等。使用的内存超过机器内存容量，
系统会使用交换区（swap），导致性能急速下降。
我们可以简单验证一下创建过多 goroutine 会发生什么：
```
func main() {
  var wg sync.WaitGroup
  wg.Add(10000000)
  for i := 0; i < 10000000; i++ {
    go func() {
      time.Sleep(1 * time.Minute)
    }()
  }
  wg.Wait()
}
```
在我的机器上（8G内存）运行上面的程序会报**errno 1455**，即Out of Memory错误，这很好理解。谨慎运行。

另一方面，goroutine 的管理也是一个问题。goroutine 只能自己运行结束，外部没有任何手段可以强制j结束一个 goroutine。如果一个 goroutine 因为某种原因没有自行结束，
就会出现 goroutine 泄露。此外，频繁创建 goroutine 也是一个开销。

鉴于上述原因，自然出现了与线程池一样的需求，即 goroutine 池。一般的 goroutine 池自动管理 goroutine 的生命周期，可以按需创建，动态缩容。向 goroutine 池提交一个任务，
goroutine 池会自动安排某个 goroutine 来处理。

**ants**就是其中一个实现 goroutine 池的库。
### 快速使用
本文代码使用 Go Modules。

创建目录并初始化：
```
$ mkdir ants && cd ants
$ go mod init github.com/darjun/go-daily-lib/ants
```
安装ants库，使用v2版本：
```
$ go get -u github.com/panjf2000/ants/v2
```
我们接下来要实现一个计算大量整数和的程序。首先创建基础的任务结构，并实现其执行任务方法：
```
type Task struct {
  index int
  nums  []int
  sum   int
  wg    *sync.WaitGroup
}

func (t *Task) Do() {
  for _, num := range t.nums {
    t.sum += num
  }

  t.wg.Done()
}
```
很简单，就是将一个切片中的所有整数相加。

然后我们创建 goroutine 池，注意池使用完后需要手动关闭，这里使用defer关闭：
```
p, _ := ants.NewPoolWithFunc(10, taskFunc)
defer p.Release()

func taskFunc(data interface{}) {
  task := data.(*Task)
  task.Do()
  fmt.Printf("task:%d sum:%d\n", task.index, task.sum)
}
```
上面调用了**ants.NewPoolWithFunc()** 创建了一个 goroutine 池。第一个参数是池容量，即池中最多有 10 个 goroutine。第二个参数为每次执行任务的函数。
当我们调用p.Invoke(data)的时候，ants池会在其管理的 goroutine 中找出一个空闲的，让它执行函数**taskFunc**，并将data作为参数。

接着，我们模拟数据，做数据切分，生成任务，交给 ants 处理：
```
const (
  DataSize    = 10000
  DataPerTask = 100
)

nums := make([]int, DataSize, DataSize)
for i := range nums {
  nums[i] = rand.Intn(1000)
}

var wg sync.WaitGroup
wg.Add(DataSize / DataPerTask)
tasks := make([]*Task, 0, DataSize/DataPerTask)
for i := 0; i < DataSize/DataPerTask; i++ {
  task := &Task{
    index: i + 1,
    nums:  nums[i*DataPerTask : (i+1)*DataPerTask],
    wg:    &wg,
  }

  tasks = append(tasks, task)
  p.Invoke(task)
}

wg.Wait()
fmt.Printf("running goroutines: %d\n", ants.Running())
```
随机生成 10000 个整数，将这些整数分为 100 份，每份 100 个，生成Task结构，调用p.Invoke(task)处理。wg.Wait()等待处理完成，然后输出ants正在运行的 goroutine 数量，这时应该是 0。

最后我们将结果汇总，并验证一下结果，与直接相加得到的结果做一个比较：
```
var sum int
for _, task := range tasks {
  sum += task.sum
}

var expect int
for _, num := range nums {
  expect += num
}

fmt.Printf("finish all tasks, result is %d expect:%d\n", sum, expect)
```
运行：
```
$ go run main.go
...
task:96 sum:53275
task:88 sum:50090
task:62 sum:57114
task:45 sum:48041
task:82 sum:45269
running goroutines: 0
finish all tasks, result is 5010172 expect:5010172
```
确实，任务完成之后，正在运行的 goroutine 数量变为 0。而且我们验证了，结果没有偏差。另外需要注意，**goroutine池中任务的执行顺序是随机的，与提交任务的先后没有关系。**
由上面运行打印的任务标识我们也能发现这一点。
### 函数作为任务
**ants**支持将一个不接受任何参数的函数作为任务提交给 goroutine 运行。由于不接受参数，我们提交的函数要么不需要外部数据，只需要处理自身逻辑，否则就必须用某种方式将需要的数据传递进去，例如闭包。

提交函数作为任务的 goroutine 池使用**ants.NewPool()** 创建，它只接受一个参数表示池子的容量。调用池子对象的**Submit()** 方法来提交任务，将一个不接受任何参数的函数传入。

最开始的例子可以改写一下。增加一个任务包装函数，将任务需要的参数作为包装函数的参数。包装函数返回实际的任务函数，该任务函数就可以通过闭包访问它需要的数据了：
```
type taskFunc func()

func taskFuncWrapper(nums []int, i int, sum *int, wg *sync.WaitGroup) taskFunc {
  return func() {
    for _, num := range nums[i*DataPerTask : (i+1)*DataPerTask] {
      *sum += num
    }

    fmt.Printf("task:%d sum:%d\n", i+1, *sum)
    wg.Done()
  }
}
```
调用ants.NewPool(10)创建 goroutine 池，同样池子用完需要释放，这里使用defer：
```
p, _ := ants.NewPool(10)
defer p.Release()
```
生成模拟数据，切分任务。提交任务给ants池执行，这里使用taskFuncWrapper()包装函数生成具体的任务，然后调用p.Submit()提交：
```
nums := make([]int, DataSize, DataSize)
for i := range nums {
  nums[i] = rand.Intn(1000)
}

var wg sync.WaitGroup
wg.Add(DataSize / DataPerTask)
partSums := make([]int, DataSize/DataPerTask, DataSize/DataPerTask)
for i := 0; i < DataSize/DataPerTask; i++ {
  p.Submit(taskFuncWrapper(nums, i, &partSums[i], &wg))
}
wg.Wait()
```
汇总结果，验证：
```
var sum int
for _, partSum := range partSums {
  sum += partSum
}

var expect int
for _, num := range nums {
  expect += num
}
fmt.Printf("running goroutines: %d\n", ants.Running())
fmt.Printf("finish all tasks, result is %d expect is %d\n", sum, expect)
```
这个程序的功能与最开始的完全相同。

### 执行流程
GitHub 仓库中有个执行流程图，我重新绘制了一下：
![image](https://user-images.githubusercontent.com/6757408/177616580-51b9e53d-32fe-4e82-b5ac-362820c059df.png)
执行流程如下：
* 初始化 goroutine 池；
* 提交任务给 goroutine 池，检查是否有空闲的 goroutine：
  * 有，获取空闲 goroutine
  * 无，检查池中的 goroutine 数量是否已到池容量上限：
    * 已到上限，检查 goroutine 池是否是非阻塞的：
      * 非阻塞，直接返回nil表示执行失败
      * 阻塞，等待 goroutine 空闲
    * 未到上限，创建一个新的 goroutine 处理任务
* 任务处理完成，将 goroutine 交还给池，以待处理下一个任务
### 选项
ants提供了一些选项可以定制 goroutine 池的行为。选项使用Options结构定义：
```
// src/github.com/panjf2000/ants/options.go
type Options struct {
  ExpiryDuration time.Duration
  PreAlloc bool
  MaxBlockingTasks int
  Nonblocking bool
  PanicHandler func(interface{})
  Logger Logger
}
```



转自：https://darjun.github.io/2021/06/03/godailylib/ants/










