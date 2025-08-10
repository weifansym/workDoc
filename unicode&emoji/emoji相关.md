## emoji的字符的unicode 范围，以及占用的字节数
关于 Emoji 字符的 Unicode 范围及其在 UTF-8 编码中的字节数，整理如下：
### 1. Emoji 的 Unicode 范围
Emoji 并不集中在单一的 Unicode 区块，而是分布在多个区块中，主要包括：
<img width="659" height="447" alt="截屏2025-08-11 00 02 44" src="https://github.com/user-attachments/assets/65ee9495-e2bb-4b05-a4bc-527477820f56" />

此外，Emoji 还可能由多个 Unicode 码点组合而成（如肤色修饰符、性别修饰符、ZWJ 连接符等），形成复合 Emoji。

#### 2. Emoji 在 UTF-8 中的字节数
Unicode 码点对应的 UTF-8 编码字节数规则：

<img width="623" height="269" alt="截屏2025-08-11 00 03 19" src="https://github.com/user-attachments/assets/70940a6f-ba17-4c84-bede-54114b97a424" />

##### **具体到 Emoji**：
* 大部分 Emoji 码点位于辅助平面（U+1F300 及以上），因此单个 Emoji 码点通常占用 **4 字节 UTF-8**。
* 一些较早的符号（如 Dingbats U+2700~U+27BF）位于 BMP 范围，占用 **3 字节 UTF-8**。
* 复合 Emoji 由多个码点组成，整体字节数是各码点字节数之和。
 
##### 3. 示例

<img width="673" height="522" alt="截屏2025-08-11 00 05 42" src="https://github.com/user-attachments/assets/631d4207-e996-4ca5-8680-71a82ede8778" />

## U+0800 – U+FFFF 编码 中包含了哪些emoji
关于 Unicode 范围 U+0800 – U+FFFF 中包含的 Emoji，以下是详细说明：
#### 1. 范围说明
* U+0800 – U+FFFF 是 Unicode 的基本多文种平面（BMP，Basic Multilingual Plane）的大部分区间。
* 这个范围内包含了大量字符，包括各种语言文字、符号、标点等。
* Emoji 在这个范围内主要分布在以下区块：

<img width="543" height="228" alt="截屏2025-08-11 00 08 20" src="https://github.com/user-attachments/assets/3c16b45b-bb83-4aea-8b46-d1070bb25036" />
  
#### 2. 具体包含的 Emoji 主要区块
##### 2.1 Miscellaneous Symbols (U+2600 – U+26FF)
* 包含许多常用符号和部分 Emoji，如天气符号、星星、符号等。
* 典型 Emoji 示例：

<img width="223" height="265" alt="截屏2025-08-11 00 09 31" src="https://github.com/user-attachments/assets/d7bc9b5f-b574-4ff0-a244-120ce6660393" />
 
##### 2.2 Dingbats (U+2700 – U+27BF)
* 包含装饰符号和部分 Emoji。
* 典型 Emoji 示例：

<img width="223" height="266" alt="截屏2025-08-11 00 10 34" src="https://github.com/user-attachments/assets/ab937403-036a-4c9f-9920-8d318b68e957" />

##### 2.3 Enclosed Alphanumeric Supplement (U+1F100 – U+1F1FF)
* 这个区块位于 BMP 末尾，包含封闭的字母数字符号，部分被用作 Emoji。
* 例如：
 
<img width="245" height="121" alt="截屏2025-08-11 00 11 26" src="https://github.com/user-attachments/assets/888af617-e717-4163-832a-f78703fe4793" />

##### 3. 其他 BMP 区域
* 大部分 Emoji 码点位于辅助平面（U+10000 以上），如表情符号、交通工具、动物等。
* BMP 中的 CJK 统一汉字区（U+4E00 – U+9FFF）不包含 Emoji。
* 其他符号区块（如数学符号、货币符号）也不属于 Emoji。

##### 4. 总结

<img width="503" height="206" alt="截屏2025-08-11 00 12 36" src="https://github.com/user-attachments/assets/46d3cb63-032e-4cde-8eea-386c4934bb73" />

如果你需要，我可以帮你列出该范围内所有 Emoji 的完整列表，或者帮你写代码检测某字符是否属于 BMP 范围内的 Emoji。

## 请使用golang代码，写出检测某字符是否属于 BMP 范围内的 Emoji。
下面是一个用 Go 语言编写的示例代码，用于检测某个字符（rune）是否属于 BMP 范围内的 Emoji。代码基于 Unicode BMP 范围内常见的 Emoji 区块（主要是 U+2600–U+26FF 和 U+2700–U+27BF）进行判断。
```
package main

import (
    "fmt"
)

// isBMPEmoji 判断一个 rune 是否属于 BMP 范围内的 Emoji
func isBMPEmoji(r rune) bool {
    // Miscellaneous Symbols: U+2600 – U+26FF
    if r >= 0x2600 && r <= 0x26FF {
        return true
    }
    // Dingbats: U+2700 – U+27BF
    if r >= 0x2700 && r <= 0x27BF {
        return true
    }
    // Enclosed Alphanumeric Supplement: U+1F100 – U+1F1FF
    // 该区块部分码点在辅助平面，不属于 BMP，故不判断
    // 这里只判断 BMP 范围内的 Emoji

    return false
}

func main() {
    testRunes := []rune{
        '☀', // U+2600 太阳
        '☁', // U+2601 云
        '✈', // U+2708 飞机
        '✌', // U+270C 胜利手势
        'A',  // 普通字母
        '😀', // U+1F600 辅助平面 Emoji，不属于 BMP
    }

    for _, r := range testRunes {
        fmt.Printf("字符: %c, Unicode: U+%04X, 是否 BMP Emoji: %v\n", r, r, isBMPEmoji(r))
    }
}
```
##### 运行结果示例：
```
字符: ☀, Unicode: U+2600, 是否 BMP Emoji: true
字符: ☁, Unicode: U+2601, 是否 BMP Emoji: true
字符: ✈, Unicode: U+2708, 是否 BMP Emoji: true
字符: ✌, Unicode: U+270C, 是否 BMP Emoji: true
字符: A, Unicode: U+0041, 是否 BMP Emoji: false
字符: 😀, Unicode: U+1F600, 是否 BMP Emoji: false
```
##### 说明
* 该函数只检测 BMP 范围内的 Emoji，主要覆盖了 Miscellaneous Symbols 和 Dingbats 区块。
* 辅助平面（如 U+1F600 及以上）的 Emoji 不属于 BMP，函数返回 false。
* 如果需要检测更全面的 Emoji，建议使用第三方库或维护完整的 Emoji 码点列表。








