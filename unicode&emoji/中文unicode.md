## 中文在unicode编码下都是3个字节码？有没一些占4个字节的
关于中文字符在 Unicode 编码下的字节数，具体说明如下：
#### 1. 中文字符的 Unicode 范围
* 绝大多数常用中文汉字位于 基本多文种平面（BMP），范围大致是：
  * CJK Unified Ideographs：U+4E00 – U+9FFF
  * CJK Unified Ideographs Extension A：U+3400 – U+4DBF
* 这些汉字都属于 BMP 范围（U+0000 – U+FFFF）。

#### 2. UTF-8 编码字节数规则回顾

<img width="436" height="195" alt="截屏2025-08-11 00 17 22" src="https://github.com/user-attachments/assets/97fae4b0-845c-414d-a175-402c1487bf84" />

#### 3. 中文字符的 UTF-8 字节数
* **大部分中文汉字（BMP 范围）：**
  * Unicode 码点在 U+0800 – U+FFFF 之间，
  * 在 UTF-8 编码下占用 3 个字节。
* 部分扩展汉字（辅助平面）：
  * Unicode 还定义了多个扩展区块（Extension B、C、D、E、F、G），范围在 U+20000 及以上，属于辅助平面。
  * 这些扩展汉字在 UTF-8 编码下占用 4 个字节。
  * 例如：
    * U+20000（𠀀）
    * U+2A6D6（𪛖）
    * U+2B740（𫝀）
   
#### 4. 总结

<img width="616" height="122" alt="截屏2025-08-11 00 19 39" src="https://github.com/user-attachments/assets/1617cecd-ae7b-4554-978b-9d9983003cb3" />

##### 5. 示例代码（Go）
下面示例演示如何判断一个中文字符的 UTF-8 字节数：
```
package main

import (
    "fmt"
    "unicode/utf8"
)

func main() {
    runes := []rune{'汉', '𠀀'} // '汉' U+6C49，'𠀀' U+20000

    for _, r := range runes {
        buf := make([]byte, 4)
        n := utf8.EncodeRune(buf, r)
        fmt.Printf("字符: %c, Unicode: U+%X, UTF-8 字节数: %d\n", r, r, n)
    }
}
```
输出：
```
字符: 汉, Unicode: U+6C49, UTF-8 字节数: 3
字符: 𠀀, Unicode: U+20000, UTF-8 字节数: 4
```
如果你需要，我可以帮你列出辅助平面扩展汉字的具体范围，或者帮你写代码判断一个字符是否属于扩展汉字。


  
