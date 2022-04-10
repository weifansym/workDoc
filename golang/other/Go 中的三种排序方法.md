### Go 中的三种排序方法
排序操作是很多程序经常使用的操作。尽管一个简短的快排程序只要二三十行代码就可以搞定，但是一个健壮的实现需要更多的代码，并且我们不希望每次我们需要的时候都重写或者拷贝这些代码。幸运的是，
Go 内置的 sort 包中提供了根据一些排序函数来对任何序列进行排序的功能。
#### 排序整数、浮点数和字符串切片
对于 **[]int,[]float,[]string** 这种元素类型是基础类型的切片使用 sort 包提供的下面几个函数进行排序。
* sort.Ints
* sort.Floats
* sort.Strings
```
s := []int{4, 2, 3, 1}
sort.Ints(s)
fmt.Println(s) // 输出[1 2 3 4]
```
#### 使用自定义比较器排序
* 使用 sort.Slice 函数排序，它使用一个用户提供的函数来对序列进行排序，函数类型为 func(i, j int) bool，其中参数 i, j 是序列中的索引。
* sort.SliceStable 在排序切片时会保留相等元素的原始顺序。
* 上面两个函数让我们可以排序结构体切片 (order by struct field value)。
```
family := []struct {
    Name string
    Age  int
}{
    {"Alice", 23},
    {"David", 2},
    {"Eve", 2},
    {"Bob", 25},
}

// 用 age 排序，年龄相等的元素保持原始顺序
sort.SliceStable(family, func(i, j int) bool {
    return family[i].Age < family[j].Age
})
fmt.Println(family) // [{David 2} {Eve 2} {Alice 23} {Bob 25}]

 //下面实现排序order by age asc, name desc，如果 age 和 name 都相等则保持原始排序
sort.SliceStable(family, func(i, j int) bool {
    if family[i].Age != family[j].Age {
        return family[i].Age < family[j].Age
    }
    return strings.Compare(family[i].Name, family[j].Name) == 1
})

fmt.Println(family) // [{Eve 2} {David 2} {Alice 23} {Bob 25}]
```
#### 排序任意数据结构
* 使用 sort.Sort 或者 sort.Stable 函数。
* 他们可以排序实现了 sort.Interface 接口的任意类型

一个内置的排序算法需要知道三个东西：序列的长度，表示两个元素比较的结果，一种交换两个元素的方式；这就是 sort.Interface 的三个方法：
```
type Interface interface {
    Len() int
    Less(i, j int) bool // i, j 是元素的索引
    Swap(i, j int)
}
```
还是以上面的结构体切片为例子，我们为切片类型自定义一个类型名，然后在自定义的类型上实现 srot.Interface 接口
```
type Person struct {
    Name string
    Age  int
}

// ByAge 通过对age排序实现了sort.Interface接口
type ByAge []Person

func (a ByAge) Len() int           { return len(a) }
func (a ByAge) Less(i, j int) bool { return a[i].Age < a[j].Age }
func (a ByAge) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

func main() {
    family := []Person{
        {"David", 2},
        {"Alice", 23},
        {"Eve", 2},
        {"Bob", 25},
    }
    sort.Sort(ByAge(family))
    fmt.Println(family) // [{David, 2} {Eve 2} {Alice 23} {Bob 25}]

    sort.
}
```
实现了 sort.Interface 的具体类型不一定是切片类型；下面的 customSort 是一个结构体类型。
```
type customSort struct {
    p    []*Person
    less func(x, y *Person) bool
}

func (x customSort) Len() int {len(x.p)}
func (x customSort) Less(i, j int) bool { return x.less(x.p[i], x.p[j]) }
func (x customSort) Swap(i, j int)      { x.p[i], x.p[j] = x.p[j], x.p[i] }
```
让我们定义一个根据多字段排序的函数，它主要的排序键是 Age，Age 相同了再按 Name 进行倒序排序。下面是该排序的调用，其中这个排序使用了匿名排序函数：
```
sort.Sort(customSort{persons, func(x, y *Person) bool {
    if x.Age != y.Age {
        return x.Age < y.Age
    }
    if x.Name != y.Name {
        return x.Name > y.Name
    }
    return false
}})
```
#### 排序具体的算法和复杂度
Go 的 sort 包中所有的排序算法在最坏的情况下会做 n log n 次 比较，n 是被排序序列的长度，所以排序的时间复杂度是 O(n log n*)。其大多数的函数都是用改良后的快速排序算法实现的。





