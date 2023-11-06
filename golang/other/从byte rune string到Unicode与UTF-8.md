## 从byte rune string到Unicode与UTF-8
Go语言使用UTF-8编码，因此任何字符都可以用Unicode表示。为此，Go在代码中引入了一个新术语，称为 rune。

rune是int32的类型别名:
```
// rune is an alias for int32 and is equivalent to int32 in all ways. It is
// used, by convention, to distinguish character values from integer values.
type rune = int32
```
另外，字符串经常被转换为[]byte使用，要详细说清楚rune、byte、字符串之间的关系，必须得从人和宇宙的关系说起，呸！从字符编码说起。

#### 1. ASCII码
通过数字电路的知识，我们知道使用二进制对信息进行编码与度量。最初现代计算机由美国人发明使用，自然而然就考虑把英语进行编码，所以ASCII码就是英语字符对应的二进制位，而且一直沿用至今，ASCII码占用1个字节，最高位统一规定为0，所以只使用了7位，一共可以表示27=128个字符，包括32个不能打印的字符。

#### 2.Unicode
现代计算机早已不是美国一家独大，互联网更是让世界互联互通。但是文字确实多种多样，各个国家拥有一套编码规则，同一个二进制数会被不同编码解释为不同符号。如果每次不把编码方式勾兑清楚，谁也不知道该怎么解码。有没有不需要勾兑的方式？有，就是抛开各个国家独有的编码方式，统一使用一个编码方式：Unicode

#### 3.UTF-8
Unicode规定了字符的二进制代码，但是却没有规定如何存储。而且，各个字符占的字节是可能不同的，比如汉字很多都有10几位二进制，可能需要2个字节，3个字节，甚至4个字节。虽有unicode对应，肯定是该多少字节就存多少字节，而不是每个字符都存相同大小字节，毕竟unicode有100多万，全存相同大小字节，肯定浪费空间。但是就有了最终要解决的问题：什么时候该读3个字节以表示1个字符，什么时候该读1个字节以表示字符？

UTF-8就是存储Unicode的方式，但不是唯一的，其他utf-16,utf-32交给童鞋们自己探索，我们主要深究一下utf-8。来看下UTF-8是如何解决上面的问题：

#### 什么时候读1个字节的字符？
字节的第一位为0，后面7位为符号的unicode码。所以这样看，英语字母的utf-8和ascii一致。

#### 什么时候读多个字节的字符？
对于有n个字节的字符，（n>1）….其中第一个字节的高n位就为1，换句话说：
* 第一个字节读到0，那就是读1个字节
* 第一个字节读到n个1，就要读n个字节

然后第一个字高n位后1位设为0，**后续其他字节前两位都设为10**
```
0xxxxxxx # 读1个字节
110xxxxx 10xxxxxx # 读两个字节
1110xxxx 10xxxxxx 10xxxxxx #读3个字节
11110xxx 10xxxxxx 10xxxxxx 10xxxxxx #读4个字节

Unicode符号范围     |        UTF-8编码方式
(十六进制)        |              （二进制）
----------------------+---------------------------------------------
0000 0000-0000 007F | 0xxxxxxx
0000 0080-0000 07FF | 110xxxxx 10xxxxxx
0000 0800-0000 FFFF | 1110xxxx 10xxxxxx 10xxxxxx
0001 0000-0010 FFFF | 11110xxx 10xxxxxx 10xxxxxx 10xxxxxx
```
#### 怎样完成UTF-8最终编码？
解决了读几个字节的问题，还有一个问题：Unicode怎么填充UTF-8的各个字节？

比如 张 字，unicode编码5F20，对应的十六进制处于0000 0800-0000 FFFF中，也就是3个字节。
* 1110xxxx 10xxxxxx 10xxxxxx
* 张的unicode对应的二进制：101 111100 100000
* 从后向前填充，高位不够的补0
  * 010000 填充至第三个字节 10xxxxxx → 10100000
  * 111100 填充至第二个字节 10xxxxxx → 10111100
  * 101 填充至第一个字节 1110xxxx → 1110x101
  * 高位补0 1110x101 → 11100101
  * 最终结果：11100101 10111100 10100000 16进制 E5BCA0

#### 4.go语言的字符串
字符串是Go 语言中最常用的基础数据类型之一，实际上字符串是一块连续的内存空间，一个由字符组成的数组，既然作为数组来说，它会占用一片连续的内存空间，这片连续的内存空间就存储了多个字节，整个字节数组组成了字符串。

#### 5.rune与byte的使用
#### Ascii码字符
```
package main

import (
	"fmt"
	"unsafe"
)

func main() {
	s := 'a'       //rune
	fmt.Println(s) // 97
	t := unsafe.Sizeof(s) 
	fmt.Println(t) // 4
}
```
a是Ascii码字符，单引号' ‘包裹的字符，go语言会将其视为rune类型，rune类型为int32，所以占4个字节。

#### 全为Ascii码的字符串
```
package main

import (
	"fmt"
	"unsafe"
)

func main() {
	b := "golang"
	fmt.Println(b)
	s_rune := []rune(b)
	s_byte := []byte(b)
	fmt.Println(s_byte) // [103 111 76 97 110 103]
	fmt.Println(s_rune) // [103 111 76 97 110 103]
}
```
* []rune()将字符串转换为rune切片
* []byte()将字符串转换为byte切片
* 由于都是Ascii码字符串，所以输出的整数都一致
#### 包含非ascii码的字符串
```
package main

import (
	"fmt"
	"unicode/utf8"
	"unsafe"
)

func main() {
	c := "go语言"
	s_rune_c := []rune(c)
	s_byte_c := []byte(c)
	fmt.Println(s_rune_c) // [103 111 35821 35328]  
	fmt.Println(s_byte_c) // [103 111 232 175 173 232 168 128]
	fmt.Println(utf8.RuneCountInString(c)) 	//4
  fmt.Println(len(c))   					//8
	fmt.Println(len(s_rune_c)) 				//4
}
```
* 汉字占3个字节，所以转换的[]byte长度为8
* 由于已经转换为[]rune，所以长度为4
* utf8.RuneCountInString()获取UTF-8编码字符串的长度，所以跟[]rune一致
#### 6.汉字的输出详解
```
package main

import (
	"fmt"
	"unsafe"
)

func main() {
	f := "张"
	s_byte_f := []byte(f)
	s_rune_f := []rune(f)
	t := unsafe.Sizeof(s_byte_f) 
	fmt.Println(s_byte_f)	// [299 188 160]
	t = unsafe.Sizeof(s_rune_f) 
	fmt.Println(s_rune_f) // [24352]
  e := '张'
	s_byte_e := byte(e)
  t = unsafe.Sizeof(s_byte_e) 
	fmt.Println(t) // 1
  fmt.Println(s_byte_e) // 张32?
}
```
24352？[299 188 160] ? 32???

* 张 输出的值24352是unicode
  * 十六进制 5F20
  * 十进制 24352
  * 二进制101111100100000
* 存储方式是utf-8
  * uft-8编码:11100101 10111100 10100000
  * 11100101 - 299
  * 10111100 - 188
  * 10100000 - 160
* 这就解释了为什么转换后的[]byte是[299 188 160]

在go语言中，byte其实是uint8的别名，byte和 uint8 之间可以直接进行互转，只能将0~255范围的int转成byte。超出这个范围，go在转换的时候，就会把多出来数据砍掉；**但是rune转byte，又有些不同：会先把rune从UTF-8转换为Unicode，由于Unicode依然超出了byte表示范围，所以取低8位，其余的全部扔掉 101111100100000，就能解释为什么是输出32**（这里有专门的[汉字对应表](http://www.chi2ko.com/tool/CJK.htm)，可以用其他做验证。）
#### 7.总结
* Go 语言中的字符串是一个只读的字节切片
* 声明的任何单个字符，go语言都会视其为rune类型
* []rune()可以把字符串转换为一个rune数组(即unicode数组)
  * 一个rune就表示一个Unicode字符
  * 每个Unicode字符，在内存中是以utf-8的形式存储
  * Unicode字符，输出[]rune，会把每个UTF-8转换为Unicode后再输出
* []byte()可以把字符串转换为一个byte数组
  * 输出[]byte,会按字符串在内存中实际存储形式(UTF-8)输出
    * Unicode字符，按[]byte输出，就会把UTF-8的每个字节单个输出
* 而Unicode字符做强制转换时，会优先计算出Unicode值，再做转换
* 对于Ascii码字符，rune与byte值是一样的
  * 这是因为Ascii码字符的Unicode也只需要1个字节，且一致

转自：http://www.randyfield.cn/post/2022-01-14-rune-unicode-utf8/
