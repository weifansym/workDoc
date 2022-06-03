## Go语言中for-range使用踩坑指南
golang常用的遍历方式，有两种： for 和 for-range。
而for-range使用中有些坑常会遇到，今天我们一起来捋一捋。
### 指针数据坑: 遍历取不到所有元素指针?
如下代码想从数组遍历获取一个指针元素切片集合
```
arr := [2]int{1, 2}
res := []*int{}
for _, v := range arr {
    res = append(res, &v)
}
//expect: 1 2
fmt.Println(*res[0],*res[1]) 
//but output: 2 2
```
答案是【取不到】
同样代码对切片[]int{1, 2}或map[int]int{1:1, 2:2}遍历也不符合预期。
问题出在哪里？

通过查看[go编译源码](https://github.com/golang/gofrontend/blob/e387439bfd24d5e142874b8e68e7039f74c744d7/go/statements.cc#L5501)可以了解到, for-range其实是语法糖，内部调用还是for循环，初始化会拷贝带遍历的列表（如array，slice，map），变量v是用来保存迭代切片所得的值，因为v只被声明了一次，每次迭代的值都是赋值给v，该变量的内存地址始终未变，这样讲他的地址追加到新的切片中，该切片保存的都是同一个地址，这肯定无法达到预期效果的。对应伪代码如下：
```
// len_temp := len(range)
// range_temp := range
// for index_temp = 0; index_temp < len_temp; index_temp++ {
//     value_temp = range_temp[index_temp]
//     index = index_temp
//     value = value_temp
//     original body
//   }
```
那么怎么改？
#### 有两种
* 使用局部变量拷贝v
```
for _, v := range arr {
    //局部变量v替换了v，也可用别的局部变量名
    v := v 
    res = append(res, &v)
}
```
* 直接索引获取原来的元素
```
//这种其实退化为for循环的简写
for k := range arr {
    res = append(res, &arr[k])
}
```
理顺了这个问题后边的坑基本都好发现了，来迅速过一遍
### 迭代修改变量问题
看下这个例子，现在我们要对切片中保存的每个用户的年龄进行修改，因为我们都是永远18岁，嘎嘎嘎～～～。
```
package main

import (
	"fmt"
)

type user struct {
	name string
	age uint64
}

func main()  {
	u := []user{
		{"asong",23},
		{"song",19},
		{"asong2020",18},
	}
	for _,v := range u{
		if v.age != 18{
			v.age = 20
		}
	}
	fmt.Println(u)
}
```
来看一下运行结果：
```
[{asong 23} {song 19} {asong2020 18}]
```
哎呀，怎么回事。怎么没有更改呢。其实道理都是一样，还记得，我在上文说的一个知识点嘛。对，就是这个，想起来了吧。v变量是拷贝切片中的数据，修改拷贝数据怎么会对原切片有影响呢，还是这个问题，copy这个知识点很重要，一不注意，就会出现问题。知道问题了，我们现在来把这个问题解决吧。
```
package main

import (
	"fmt"
)

type user struct {
	name string
	age uint64
}

func main()  {
	u := []user{
		{"asong",23},
		{"song",19},
		{"asong2020",18},
	}
	for k,v := range u{
		if v.age != 18{
			u[k].age = 18
		}
	}
	fmt.Println(u)
}
```
可以看到，我们直接对切片的值进行修改，这样就修改成功了。所以这里还是要注意一下的，防止以后出现bug。
### 遍历会停止么？
```
v := []int{1, 2, 3}
for i := range v {
    v = append(v, i)
}
```
答案是【会】，因为遍历前对v做了拷贝，所以期间对原来v的修改不会反映到遍历中

### 对大数组这样遍历有啥问题？
```
//假设值都为1，这里只赋值3个
var arr = [102400]int{1, 1, 1} 
for i, n := range arr {
    //just ignore i and n for simplify the example
    _ = i 
    _ = n 
}
```
答案是【有问题】！遍历前的拷贝对内存是极大浪费啊
怎么优化？有两种

* 对数组取地址遍历
```
for i, n := range &arr
```
* 对数组做切片引用
```
for i, n := range arr[:]
```
反思题：对大量元素的slice和map遍历为啥不会有内存浪费问题？
（提示，底层数据结构是否被拷贝）
### 对大数组这样重置效率高么？
```
//假设值都为1，这里只赋值3个
var arr = [102400]int{1, 1, 1} 
for i, _ := range &arr {
    arr[i] = 0
}
```
答案是【高】，这个要理解得知道go对这种重置元素值为默认值的遍历是有优化的, 详见[go源码：memclrrange](https://github.com/golang/go/blob/ea020ff3de9482726ce7019ac43c1d301ce5e3de/src/cmd/compile/internal/gc/range.go#L363)
```
// Lower n into runtime·memclr if possible, for
// fast zeroing of slices and arrays (issue 5373).
// Look for instances of
//
// for i := range a {
// 	a[i] = zero
// }
//
// in which the evaluation of a is side-effect-free.
```
### 对map遍历时删除元素能遍历到么？
```
var m = map[int]int{1: 1, 2: 2, 3: 3}
//only del key once, and not del the current iteration key
var o sync.Once 
for i := range m {
    o.Do(func() {
        for _, key := range []int{1, 2, 3} {
            if key != i {
                fmt.Printf("when iteration key %d, del key %d\n", i, key)
                delete(m, key)
                break
            }
        }
    })
    fmt.Printf("%d%d ", i, m[i])
}
```
答案是【不会】
map内部实现是一个链式hash表，为保证每次无序，初始化时会随机一个遍历开始的位置,
这样，如果删除的元素开始没被遍历到（上边once.Do函数内保证第一次执行时删除未遍历的一个元素），那就后边就不会出现。

### 对map遍历时新增元素能遍历到么？
```
var m = map[int]int{1:1, 2:2, 3:3}
for i, _ := range m {
    m[4] = 4
    fmt.Printf("%d%d ", i, m[i])
}
```
答案是【可能会】，输出中可能会有44。原因同上一个, 可以用以下代码验证
```
var createElemDuringIterMap = func() {
    var m = map[int]int{1: 1, 2: 2, 3: 3}
    for i := range m {
        m[4] = 4
        fmt.Printf("%d%d ", i, m[i])
    }
}
for i := 0; i < 50; i++ {
    //some line will not show 44, some line will
    createElemDuringIterMap()
    fmt.Println()
}
```
### 这样遍历中起goroutine可以么？
```
var m = []int{1, 2, 3}
for i := range m {
    go func() {
        fmt.Print(i)
    }()
}
//block main 1ms to wait goroutine finished
time.Sleep(time.Millisecond)
```
答案是【不可以】。预期输出0,1,2的某个组合，如012，210..
结果是222. 同样是拷贝的问题
怎么解决

* 以参数方式传入
```
for i := range m {
    go func(i int) {
        fmt.Print(i)
    }(i)
}
```
* 使用局部变量拷贝
```
for i := range m {
    i := i
    go func() {
        fmt.Print(i)
    }()
}
```
发现没，一个简单的for-range，仔细剖析下来也是有不少有趣的地方。

参考：
* https://asong.cloud/go%E8%AF%AD%E8%A8%80range%E8%B8%A9%E5%9D%91%E6%8C%87%E5%8D%97/
* http://blog.newbmiao.com/2020/01/03/dig101-golang-for-range.html

