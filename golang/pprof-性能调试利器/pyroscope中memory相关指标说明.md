### Pyroscope 在内存分析的时候,alloc_space,alloc_objects,inuse_objects,inuse_space这几个指标都有什么区别？那个对查找内存泄漏帮助最大
非常好的问题！在使用 Pyroscope（或 go tool pprof）进行 内存分析时，常见的四个指标：
```
- alloc_space
- alloc_objects
- inuse_space
- inuse_objects
```
它们各自有不同的意义与用途，下面为你详细对比，并重点说明**哪个对内存泄漏排查最有用**。
<img width="816" height="380" alt="截屏2025-07-16 17 03 00" src="https://github.com/user-attachments/assets/361b892b-c247-43ff-a132-82cbfda3c3b1" />

### 哪个对查找内存泄漏帮助最大？
**优先使用 inuse_space**

> 它可以直接反映 当前哪些函数/调用栈的内存仍然未释放（可能泄漏）。

### 举个例子
```
type User struct {
    ID   int
    Name string
}
```
如果你不断创建 User 对象并加入一个 []User 切片中，但从不释放它：
* alloc_space 会一直增长（累计内存分配）；
* inuse_space 也会增长（因为你没有释放）；
* inuse_space 可告诉你内存一直挂在哪些函数/变量中；
* 但如果你已经清理引用，GC 回收了，那么：
  * alloc_space 仍然高；
  * inuse_space 会回落。

<img width="668" height="272" alt="截屏2025-07-16 17 08 17" src="https://github.com/user-attachments/assets/0ea41773-6f69-480d-b63b-aa22d9f7e08a" />

### 实战技巧
* 内存泄漏常见来源：
  * 全局变量引用了大量数据结构；
  * goroutine 泄漏，携带上下文对象未退出；
  * sync.Map 或缓存没有及时清理；
* 用 inuse_space 火焰图找到堆顶函数（如 handler.HandleConn → doSomething → []*User）即可定位泄漏点。

### 总结
<img width="664" height="276" alt="截屏2025-07-16 17 10 44" src="https://github.com/user-attachments/assets/a72eeff5-96d9-4dc0-b6f8-8c47290274ec" />



