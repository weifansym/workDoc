## fmt 源码分析——fmt 如何进行格式化？
![image](https://user-images.githubusercontent.com/6757408/162599376-55f625c6-5fea-47c7-be9f-f68f2e7041b7.png)
本文将介绍 fmt 包格式化的一些原理，以及 Formatter、State和Stringer这几个接口的作用。
### format
fmt包虽然不建议用来打印日志，但是格式化字符串确实是必不可少的，比如打印日志的时候。先详细介绍一下格式化的格式format。
format由百分号%开始，后面的部分可以分为四部分：
#### verb 占位符。
完整的格式可以参考Go 文档，下面我大概列几个：
```
%v  通过默认格式打印
%t  用于布尔类型，打印true或者false
%d  以10进制格式打印数字
%c  将数据转换成 Unicode 里面的字符打印
%x  以16进制格式打印数字
%e  科学计数法表示
%f  以10进制表示浮点数
%s  字符串
%p  指针，以0x开头的16进制地址
```
还有Go语言自己定义的类型：
```
%#v
```
#### 宽度
比如%3c，c是占位符，表示把整数转成 Unicode 字符展示，而前面的3就是宽度了。

源码如下，看到num = num*10 + int(s[newi]-'0')很熟悉有没有，就是一个把字符转成整形的方法。那么%3c返回的num就是3了。
```
func parsenum(s string, start, end int) (num int, isnum bool, newi int) {
	if start >= end {
		return 0, false, end
	}
	for newi = start; newi < end && '0' <= s[newi] && s[newi] <= '9'; newi++ {
		if tooLarge(num) {
			return 0, false, end // Overflow; crazy long number most likely.
		}
		num = num*10 + int(s[newi]-'0')
		isnum = true
	}
	return
}
```
打印下面的语句：
```
fmt.Printf("%3c\n", 'a')
```
控制输出3位，a输出占用一位，前面需要补两个0。在源码中，是通过pad实现：
```
func (f *fmt) pad(b []byte) {
	...
	width := f.wid - utf8.RuneCount(b)
	if !f.minus {
		// left padding
		f.writePadding(width)
		f.buf.Write(b)
	} else {
		// right padding
		f.buf.Write(b)
		f.writePadding(width)
	}
}
```
上面的代码width就是2，调用writePadding会打印对应宽度的空格。

**如果打印的内容很长，比如有10位，而宽度只设置了3位，会展示完整的数字还是只显示3位呢？**

答案是完整展示，当n小于等于0，直接返回，而打印内容方面，不受影响：
```
func (f *fmt) writePadding(n int) {
	if n <= 0 { // No padding bytes needed.
		return
	}
	...
```
#### 精度
比如%3.2f，f是占位符，表示浮点数展示。3表示宽度，而小数点后面的2则是精度。精度在浮点数的格式化中会用到。精度的控制是通过strconv包的字符串转换函数来实现的：
```
num := strconv.AppendFloat(f.intbuf[:1], v, byte(verb), prec, size)
```
#### 标记
除了宽度和精度，还有标记可以用来控制输出。
```
+      总打印数值的正负号；对于%q（%+q）保证只输出ASCII编码的字符。 
-      在右侧而非左侧填充空格（左对齐该区域）
#      备用格式：为八进制添加前导 0（%#o）      Printf("%#U", '中')      U+4E2D
       为十六进制添加前导 0x（%#x）或 0X（%#X），为 %p（%#p）去掉前导 0x；
       如果可能的话，%q（%#q）会打印原始 （即反引号围绕的）字符串；
       如果是可打印字符，%U（%#U）会写出该字符的
       Unicode 编码形式（如字符 x 会被打印成 U+0078 'x'）。
' '    (空格)为数值中省略的正负号留出空白（% d）；
       以十六进制（% x, % X）打印字符串或切片时，在字节之间用空格隔开
0      填充前导的0而非空格；对于数字，这会将填充移到正负号之后
```
#### fmt.State 和 fmt.Formatter
上面提到的占位符、宽度、精度和标记，除了占位符，剩下的3个在解析后被保存到了接口fmt.State里面。这个接口还增加了一个函数Write用于写入数据。
```
type State interface {
	// Write is the function to call to emit formatted output to be printed.
	Write(b []byte) (n int, err error)
	// Width returns the value of the width option and whether it has been set.
	Width() (wid int, ok bool)
	// Precision returns the value of the precision option and whether it has been set.
	Precision() (prec int, ok bool)

	// Flag reports whether the flag c, a character, has been set.
	Flag(c int) bool
}
```
它会在Formatter接口中被用到。参数c就是占位符，这些终于都凑齐了。这个接口用来自定义格式化方法，你可以在自己的结构体中实现Format函数来实现自动调用解析。
```
type Formatter interface {
	Format(f State, c rune)
}
```
### 常见类型的格式化方法
func (p *pp) printArg(arg interface{}, verb rune)是底层真正进行转换的函数。

#### 指针 %p，类型 %T
```
func (p *pp) printArg(arg interface{}, verb rune) {
	...
	// Special processing considerations.
	// %T (the value's type) and %p (its address) are special; we always do them first.
	switch verb {
	case 'T':
		p.fmt.fmtS(reflect.TypeOf(arg).String())
		return
	case 'p':
		p.fmtPointer(reflect.ValueOf(arg), 'p')
		return
	}
	...
```
对于类型和指针的转换，有现成的方法调用，而这两个转换都是通过反射实现。

**这里并没有判断是否调用用户自定义的 Format 函数，说明所有类型打印内存地址和类型都只能通过上面的代码实现，不能自定义。**

#### 数字
数字支持多种进制，16进制、8进制、4进制、2进制、10进制。在fmtInteger中通过求余法实现。
```
switch base {
	case 10:
		for u >= 10 {
			i--
			next := u / 10
			buf[i] = byte('0' + u - next*10)
			u = next
		}
	...
```
#### 万能通用格式，%v
万能格式其实也有映射关系：
```
int, int8 etc.:          %d
uint, uint8 etc.:        %d, %#x if printed with %#v
float32, complex64, etc: %g
string:                  %s
chan:                    %p
pointer:                 %p
```
一般结构体会用到这种打印方式。如果是结构体：
```
if p.fmt.sharpV {
	p.buf.WriteString(f.Type().String())
}
p.buf.WriteByte('{')
for i := 0; i < f.NumField(); i++ {
	if i > 0 {
		if p.fmt.sharpV {
			p.buf.WriteString(commaSpaceString)
		} else {
			p.buf.WriteByte(' ')
		}
	}
	if p.fmt.plusV || p.fmt.sharpV {
		if name := f.Type().Field(i).Name; name != "" {
			p.buf.WriteString(name)
			p.buf.WriteByte(':')
		}
	}
	p.printValue(getField(f, i), verb, depth+1)
}
p.buf.WriteByte('}')
```
通过反射拿到字段 Field 和内容，如果格式是%+v，也就是p.fmt.plusV是true，这样会打印字段名称。

#### 异常
转换的时候还会有异常捕获，这个在 Go 源码中不多见：
```
defer p.catchPanic(p.arg, verb)
p.fmtString(v.String(), verb)
```
如果在转换的时候发生异常panic，并不会发生异常，转换后的结果会是这个样子：
```
type data struct {
	A string
	B int
}

func (d *data) String() string {
	panic("implement me")
}

func main() {
	d := &data{"1", 2}
	fmt.Printf("%s\n", d) // prints: %!s(PANIC=implement me)
}
```
结果是%!s(PANIC=implement me)，会有 PANIC 的字样。还有一个地方很有趣，String()方法并没有按要求返回字符串，只有一个panic，这样可以编译过。

#### fmt.Stringer
顺道介绍一下Stringer接口，上面的data对象就实现了这个方法。如果是通过%s打印，或者直接调用的Println，这时候会判断这个对象是否实现了Stringer接口，如果实现了，就调用对象的String方法，
上一节的data就是这个例子。
```
type Stringer interface {
	String() string
}
```
#### 一个 fmt.Formatter 例子
还是针对上面的data类型，我实现了Formatter接口：
```
func (d *data) Format(f fmt.State, c rune) {
	switch c {
	case 'v': // &{1 2}
		buf, err := json.Marshal(d)
		if err != nil {
			panic(err)
		}
		f.Write(buf)
	case 's':
		f.Write([]byte(d.String()))
	case 'x', 'X':
		//case 'p':
		v := reflect.ValueOf(d)
		f.Write([]byte{'('})
		f.Write([]byte(v.Type().String()))
		f.Write([]byte{')', '('})
		u := v.Pointer()
		f.Write([]byte(strconv.FormatUint(uint64(u), 16)))
		f.Write([]byte{')'})
	default:
		f.Write([]byte("http://cyeam.com"))
	}
}

d := &data{"1", 2}
fmt.Printf("v %v\n", d)
fmt.Printf("s %s\n", d)
fmt.Printf("p %p\n", d)
fmt.Printf("T %T\n", d)
fmt.Printf("b %b\n", d)
fmt.Printf("o %o\n", d)
fmt.Printf("x %x\n", d)
fmt.Printf("d %d\n", d)
```
结果如下：
```
v {“A”:”1”,”B”:2} s {“A”:”1”,”B”:2} p 0xc00006c020 T main.data b http://cyeam.com o http://cyeam.com x (main.data)(c00006c020) d http://cyeam.com
```
* b、o、d我没有实现，所以返回的是一个默认值；
* v是返回的json编码；
* p和T在前面也介绍了，它并不会调用Format，所以虽然我并没有实现这两个占位符，但是结果是对的；
* x手写一个基于反射的实现，能返回变量名称和地址。

### 完整流程
格式解析，把fmt.State接口要用到的数据解析完成
```
func (p *pp) doPrintf(format string, a []interface{}) {
	...
	// Do we have flags?
	// 解析格式串中的标记
	for ; i < end; i++ {
		c := format[i]
		switch c {
		case '#':
			p.fmt.sharp = true
		case '0':
			p.fmt.zero = !p.fmt.minus // Only allow zero padding to the left.
		case '+':
			p.fmt.plus = true
		case '-':
			p.fmt.minus = true
			p.fmt.zero = false // Do not pad with zeros to the right.
		case ' ':
			p.fmt.space = true
		default:
	}
	...
	// Do we have width?
	if i < end && format[i] == '*' {
	...
	} else {
		// 解析了格式串中的宽度内容
		p.fmt.wid, p.fmt.widPresent, i = parsenum(format, i, end)
		if afterIndex && p.fmt.widPresent { // "%[3]2d"
			p.goodArgNum = false			
		}
	}
	...
	// Do we have precision?
	if i+1 < end && format[i] == '.' {
		...
		// 解析了格式串中的精度内容
		p.fmt.prec, p.fmt.precPresent, i = parsenum(format, i, end)
		...
}

```
在func (p *pp) printArg(arg interface{}, verb rune)中进行格式化转换编码；
如果对象值是空，直接打印
```
if arg == nil {
	switch verb {
	case 'T', 'v':
		p.fmt.padString(nilAngleString)
	default:
		p.badVerb(verb)
	}
	return
}
```
如果是指针或者类型格式化，调用反射实现
```
switch verb {
case 'T':
	p.fmt.fmtS(reflect.TypeOf(arg).String())
	return
case 'p':
	p.fmtPointer(reflect.ValueOf(arg), 'p')
	return
}
```
格式化数据
```
switch f := arg.(type) {
	case bool:
		p.fmtBool(f, verb)
	case float32:
		p.fmtFloat(float64(f), 32, verb)
	case float64:
	...
	default:
		if !p.handleMethods(verb) {
			// Need to use reflection, since the type had no
			// interface methods that could be used for formatting.
			p.printValue(reflect.ValueOf(f), verb, 0)
		}
}
```
* 每种内置类型都有自己的格式化实现，这样就避免了反射；
* 如果不是内置类型，判断是否实现了Formatter接口，如果实现了调用此接口；
* 如果需要转成字符串，而对象实现了Stringer接口，调用其String方法转换；
* 上面两个逻辑在函数func (p *pp) handleMethods(verb rune) (handled bool)中，如果能通过接口实现转换，返回true并格式化数据，否则返回false；(其实还有一些细节的逻辑，
例如GoStringer，我就不展开细说了)
* 如果通过上面的转换失败，则需要使用默认转换策略。

默认转换策略 p.printValue(reflect.ValueOf(f), verb, 0)
默认转换就是通过反射实现，以结构体为例，如果反射出来是结构体，那就遍历所有字段打印，逻辑和上面提到的万能转换里提到的差不多。

### 总结
从格式化的完整流程中可以发现，底层格式化算法是有对性能优化的，那就是通过对每种内置对象单独编写格式化实现来规避反射来提高性能。

实际工作中经常需要对系统内复杂结构进行格式化，那么为这些对象实现Formatter接口也算是一种提升性能的有效方式。

本文涉及的完整代码请看这里。

转自：https://blog.cyeam.com/golang/2018/09/10/fmt
