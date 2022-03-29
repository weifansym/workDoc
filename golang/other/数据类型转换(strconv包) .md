## 数据类型转换(strconv包) 
Go不会对数据进行隐式的类型转换，只能手动去执行转换操作。
#### 简单的转换操作
转换数据类型的方式很简单。
```
valueOfTypeB = typeB(valueOfTypeA)
```
例如：
```
// 浮点数
a := 5.0

// 转换为int类型
b := int(a)
```
Go允许在底层结构相同的两个类型之间互转。例如：
```
// IT类型的底层是int类型
type IT int

// a的类型为IT，底层是int
var a IT = 5

// 将a(IT)转换为int，b现在是int类型
b := int(5)

// 将b(int)转换为IT，c现在是IT类型
c := IT(b)
```
但注意：
1. 不是所有数据类型都能转换的，例如字母格式的string类型"abcd"转换为int肯定会失败
2. 低精度转换为高精度时是安全的，高精度的值转换为低精度时会丢失精度。例如int32转换为int16，float32转换为int
3. 这种简单的转换方式不能对int(float)和string进行互转，要跨大类型转换，可以使用**strconv**包提供的函数

#### strconv
strconv包提供了字符串与简单数据类型之间的类型转换功能。可以将简单类型转换为字符串，也可以将字符串转换为其它简单类型。

这个包里提供了很多函数，大概分为几类：
* 字符串转int：Atoi()
* int转字符串: Itoa()
* ParseTP类函数将string转换为TP类型：ParseBool()、ParseFloat()、ParseInt()、ParseUint()。因为string转其它类型可能会失败，所以这些函数都有第二个返回值表示是否转换成功
* FormatTP类函数将其它类型转string：FormatBool()、FormatFloat()、FormatInt()、FormatUint()
* AppendTP类函数用于将TP转换成字符串后append到一个slice中：AppendBool()、AppendFloat()、AppendInt()、AppendUint()

还有其他一些基本用不上的函数，见官方手册：**go doc strconv**或者https://golang.org/pkg/strconv/。

当有些类型无法转换时，将报错，返回的错误是strconv包中自行定义的error类型。有两种错误：
```
var ErrRange = errors.New("value out of range")
var ErrSyntax = errors.New("invalid syntax")
```
例如，使用Atoi("a")将"a"转换为int类型，自然是不成功的。如果print输出err信息，将显示：
```
strconv.Atoi: parsing "a": invalid syntax
```
#### string和int的转换
最常见的是字符串和int之间的转换：
##### 1.int转换为字符串：Itoa()
```
// Itoa(): int -> string
println("a" + strconv.Itoa(32))  // a32
```
##### 2.string转换为int：Atoi()
```
func Atoi(s string) (int, error)
```
由于string可能无法转换为int，所以这个函数有两个返回值：第一个返回值是转换成int的值，第二个返回值判断是否转换成功。
```
// Atoi(): string -> int
i,_ := strconv.Atoi("3")
println(3 + i)   // 6

// Atoi()转换失败
i,err := strconv.Atoi("a")
if err != nil {
    println("converted failed")
}
```
#### Parse类函数
**Parse类函数用于转换字符串为给定类型的值**：ParseBool()、ParseFloat()、ParseInt()、ParseUint()。
由于字符串转换为其它类型可能会失败，所以这些函数都有两个返回值，第一个返回值保存转换后的值，第二个返回值判断是否转换成功。
```
b, err := strconv.ParseBool("true")
f, err := strconv.ParseFloat("3.1415", 64)
i, err := strconv.ParseInt("-42", 10, 64)
u, err := strconv.ParseUint("42", 10, 64)
```
ParseFloat()只能接收float64类型的浮点数。

ParseInt()和ParseUint()有3个参数：
```
func ParseInt(s string, base int, bitSize int) (i int64, err error)
func ParseUint(s string, base int, bitSize int) (uint64, error)
```
**bitSize**: 参数表示转换为什么位的int/uint，有效值为0、8、16、32、64。当bitSize=0的时候，表示转换为int或uint类型。例如bitSize=8表示转换后的值的类型为int8或uint8。
**base**: 参数表示以什么进制的方式去解析给定的字符串，有效值为0、2-36。当base=0的时候，表示根据string的前缀来判断以什么进制去解析：
0x开头的以16进制的方式去解析，0开头的以8进制方式去解析，其它的以10进制方式解析。

以10进制方式解析"-42"，保存为int64类型：
```
i, _ := strconv.ParseInt("-42", 10, 64)
```
以5进制方式解析"23"，保存为int64类型：
```
i, _ := strconv.ParseInt("23", 5, 64)
println(i)    // 13
```
因为5进制的时候，23表示进位了2次，再加3，所以对应的十进制数为5*2+3=13。

以16进制解析23，保存为int64类型：
```
i, _ := strconv.ParseInt("23", 16, 64)
println(i)    // 35
```
因为16进制的时候，23表示进位了2次，再加3，所以对应的十进制数为16*2+3=35。

以15进制解析23，保存为int64类型：
```
i, _ := strconv.ParseInt("23", 15, 64)
println(i)    // 33
```
因为15进制的时候，23表示进位了2次，再加3，所以对应的十进制数为15*2+3=33。
#### Format类函数
将给定类型格式化为string类型：FormatBool()、FormatFloat()、FormatInt()、FormatUint()。
```
s := strconv.FormatBool(true)
s := strconv.FormatFloat(3.1415, 'E', -1, 64)
s := strconv.FormatInt(-42, 16)
s := strconv.FormatUint(42, 16)
```
FormatInt()和FormatUint()有两个参数：
```
func FormatInt(i int64, base int) string
func FormatUint(i uint64, base int) string
```
第二个参数base指定将第一个参数转换为多少进制，有效值为2<=base<=36。当指定的进制位大于10的时候，超出10的数值以a-z字母表示。例如16进制时，10-15的数字分别使用a-f表示，
17进制时，10-16的数值分别使用a-g表示。

例如：FormatInt(-42, 16)表示将-42转换为16进制数，转换的结果为-2a。

FormatFloat()参数众多：
```
func FormatFloat(f float64, fmt byte, prec, bitSize int) string
```
bitSize表示f的来源类型（32：float32、64：float64），会据此进行舍入。

fmt表示格式：'f'（-ddd.dddd）、'b'（-ddddp±ddd，指数为二进制）、'e'（-d.dddde±dd，十进制指数）、'E'（-d.ddddE±dd，十进制指数）、'g'（指数很大时用'e'格式，否则'f'格式）、
'G'（指数很大时用'E'格式，否则'f'格式）。

prec控制精度（排除指数部分）：对'f'、'e'、'E'，它表示小数点后的数字个数；对'g'、'G'，它控制总的数字个数。如果prec 为-1，则代表使用最少数量的、但又必需的数字来表示f。
#### Append类函数
AppendTP类函数用于将TP转换成字符串后append到一个slice中：AppendBool()、AppendFloat()、AppendInt()、AppendUint()。

Append类的函数和Format类的函数工作方式类似，只不过是将转换后的结果追加到一个slice中。
```
package main

import (
	"fmt"
	"strconv"
)

func main() {
    // 声明一个slice
	b10 := []byte("int (base 10):")
    
    // 将转换为10进制的string，追加到slice中
	b10 = strconv.AppendInt(b10, -42, 10)
	fmt.Println(string(b10))

	b16 := []byte("int (base 16):")
	b16 = strconv.AppendInt(b16, -42, 16)
	fmt.Println(string(b16))
}
```
输出结果：
```
int (base 10):-42
int (base 16):-2a
```

