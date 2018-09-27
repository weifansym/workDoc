## 关于CAS与ABA相关问题
CAS即（compare and swap）是原子操作的一种，可用于在多线程编程中实现不被打断的数据交换操作，从而避免多线程同时改写某一数据时由于执行顺序不确定性以及
中断的不可预知性产生的数据不一致问题。 该操作通过将内存中的值与指定数据进行比较，当数值一样时将内存中的数据替换为新的值。

一个CAS操作的过程可以用以下c代码表示: 
```
int cas(long *addr, long old, long new)
{
    /* Executes atomically. */
    if(*addr != old)
        return 0;
    *addr = new;
    return 1;
}
```
在使用上，通常会记录下某块内存中的旧值，通过对旧值进行一系列的操作后得到新值，然后通过CAS操作将新值与旧值进行交换。如果这块内存的值在这期间内没被修改过，
则旧值会与内存中的数据相同，这时CAS操作将会成功执行 使内存中的数据变为新值。如果内存中的值在这期间内被修改过，则一般[2]来说旧值会与内存中的数据不同，
这时CAS操作将会失败，新值将不会被写入内存。

### 应用
在应用中CAS可以用于实现无锁数据结构，常见的有无锁队列(先入先出)[3] 以及无锁堆(先入后出)。对于可在任意位置插入数据的链表以及双向链表，
实现无锁操作的难度较大。

### ABA问题
ABA问题是无锁结构实现中常见的一种问题，可基本表述为：

1. 进程P1读取了一个数值A
2. P1被挂起(时间片耗尽、中断等)，进程P2开始执行
3. P2修改数值A为数值B，然后又修改回A
4. P1被唤醒，比较后发现数值A没有变化，程序继续执行。

对于P1来说，数值A未发生过改变，但实际上A已经被变化过了，继续使用可能会出现问题。在CAS操作中，由于比较的多是指针，这个问题将会变得更加严重。

具体查看: [wikipedia关于CAS](https://zh.wikipedia.org/wiki/%E6%AF%94%E8%BE%83%E5%B9%B6%E4%BA%A4%E6%8D%A2#cite_ref-2)
### 解决ABA问题
通过给数据添加版本号或者计数器都可解决ABA问题。
### 参看
* https://coolshell.cn/articles/8239.html
* https://zhuanlan.zhihu.com/p/28049542
* https://www.w3cschool.cn/architectroad/architectroad-cas-optimization.html
* https://monkeysayhi.github.io/2018/01/02/CAS%E4%B8%AD%E7%9A%84ABA%E9%97%AE%E9%A2%98/
