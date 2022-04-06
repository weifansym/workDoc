在go中，可以用strings包里的Replace方法来做字符串替换
### 函数申明
```
package strings

// Replace returns a copy of the string s with the first n
// non-overlapping instances of old replaced by new.
// If old is empty, it matches at the beginning of the string
// and after each UTF-8 sequence, yielding up to k+1 replacements
// for a k-rune string.
// If n < 0, there is no limit on the number of replacements.
func Replace(s, old, new string, n int) string{
    ...
}
```
返回将s中前n个不重叠old子串都替换为new的新字符串，
如果old为空，则向前插入n个new
如果n<0会替换所有old子串。
```
package log

import (
    "fmt"
    "strings"
)

func main(){
        s := "123abcd123abcd123abcd123abcd123abcd"
    old := "123"
    new := "3915"
    // n < 0 ,用 new 替换所有匹配上的 old；n=-1:  3915abcd3915abcd3915abcd3915abcd3915abcd
    fmt.Println("n=-1: ", strings.Replace(s, old, new, -1))

    // n = 0 ,不替换任何匹配的 old; n=0: 123abcd123abcd123abcd123abcd123abcd
    fmt.Println("n=0: ", strings.Replace(s, old, new, 0))

    // n = 1 ,用 new 替换第一个匹配的 old；n=-1:  3915abcd123abcd123abcd123abcd123abcd
    fmt.Println("n=1: ", strings.Replace(s, old, new, 1))

    // n = 2 ,用 new 替换第二个匹配的 old；n=-1:  3915abcd3915abcd123abcd123abcd123abcd
    fmt.Println("n=0: ", strings.Replace(s, old, new, 2))

    // n = 2,old="" 在最前面插入二个new；n=2:  39151391523abcd123abcd123abcd123abcd123abcd
    fmt.Println("n=2: ", strings.Replace(s, "", new, 2))
}
```
