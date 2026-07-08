# python语言中with as的用法使用详解

**With语句是什么？**

有一些任务，可能事先需要设置，事后做清理工作。对于这种场景，Python的with语句提供了一种非常方便的处理方式。一个很好的例子是文件处理，你需要获取一个文件句柄，从文件中读取数据，然后关闭文件句柄。  
如果不用with语句，代码如下：

```csharp
file = open("/tmp/foo.txt")
data = file.read()
file.close()
```

这里有两个问题。一是可能忘记关闭文件句柄；二是文件读取数据发生异常，没有进行任何处理。下面是处理异常的加强版本：

```csharp
file = open("/tmp/foo.txt")
try:
    data = file.read()
finally:
    file.close()
```

虽然这段代码运行良好，但是太冗长了。这时候就是with一展身手的时候了。除了有更优雅的语法，with还可以很好的处理上下文环境产生的异常。下面是with版本的代码：

```csharp
with open("/tmp/foo.txt") as file:
    data = file.read()
```

**with如何工作？**

这看起来充满魔法，但不仅仅是魔法，Python对with的处理还很聪明。基本思想是with所求值的对象必须有一个`__enter__()`方法，一个`__exit__()`方法。

紧跟with后面的语句被求值后，返回对象的`__enter__()`方法被调用，这个方法的返回值将被赋值给as后面的变量。当with后面的代码块全部被执行完之后，将调用前面返回对象的`__exit__()`方法。

下面例子可以具体说明with如何工作：

```python
#!/usr/bin/env python
# with_example01.py

class Sample:
    def __enter__(self):
        print "In __enter__()"
        return "Foo"

    def __exit__(self, type, value, trace):
        print "In __exit__()"

def get_sample():
    return Sample()

with get_sample() as sample:
    print "sample:", sample
```

运行代码，输出如下

```scss
In __enter__()
sample: Foo
In __exit__()
```

正如你看到的，

1.  `__enter__()`方法被执行
2.  `__enter__()`方法返回的值 - 这个例子中是"Foo"，赋值给变量'sample'
3.  执行代码块，打印变量"sample"的值为 "Foo"
4.  `__exit__()`方法被调用

with真正强大之处是它可以处理异常。可能你已经注意到Sample类的`__exit__`方法有**三个参数- val, type 和 trace。 这些参数在异常处理中相当有用**。我们来改一下代码，看看具体如何工作的。

```python
#!/usr/bin/env python
# with_example02.py
class Sample:
    def __enter__(self):
        return self

    def __exit__(self, type, value, trace):
        print "type:", type
        print "value:", value
        print "trace:", trace

    def do_something(self):
        bar = 1/0
        return bar + 10

with Sample() as sample:
    sample.do_something()
```

这个例子中，with后面的get\_sample()变成了Sample()。这没有任何关系，只要紧跟with后面的语句所返回的对象有`__enter__()`和`__exit__()`方法即可。此例中，Sample()的`__enter__()`方法返回新创建的Sample对象，并赋值给变量sample。

代码执行后：

```vbnet
bash-3.2$ ./with_example02.py\
type: <type 'exceptions.ZeroDivisionError'>\
value: integer division or modulo by zero\
trace: <traceback object at 0x1004a8128>\
Traceback (most recent call last):\
  File "./with_example02.py", line 19, in <module>\
    sample.do_something()\
  File "./with_example02.py", line 15, in do_something\
    bar = 1/0\
ZeroDivisionError: integer division or modulo by zero
```

实际上，在with后面的代码块抛出任何异常时，**exit**()方法被执行。正如例子所示，异常抛出时，与之关联的type，value和stack trace传给`__exit__()`方法，因此抛出的ZeroDivisionError异常被打印出来了。**开发库时，清理资源，关闭文件**等等操作，都可以放在`__exit__`方法当中。

因此，Python的with语句是提供一个有效的机制，让代码更简练，同时在异常产生时，清理工作更简单。

**with-as语句**

从python2.6开始，with就成为默认关键字了。With是一个控制流语句，跟if for while try之类的是一类，with可以用来简化try-finally代码，看起来比try finally更清晰，所以说with用很优雅的方式处理上下文环境产生的异常。with关键字的用法如下：

```csharp
with expression as variable:
    with block
```

该代码快的执行过程是：

1.先执行expression，然后执行该表达式返回的对象实例的`__enter__`函数，然后将该函数的返回值赋给as后面的变量。(注意，是将\_\_enter\_\_函数的返回值赋给变量)  
2.然后执行with block代码块，不论成功，错误，异常，在with block执行结束后，会执行第一步中的实例的`__exit__`函数

**with-as语句使用举例**

（1）打开文件的例子

with-as语句最常见的一个用法是打开文件的操作，如下：

```csharp
with open("decorator.py") as file:
    print file.readlines()
```

（2）自定义  
with语句后面的对象必须要有`__enter__`和`__exit__`方法，如下是一个自定义的例子：

```python
class WithTest():
    def __init__(self,name):
        self.name = name
        pass

    def __enter__(self):
        print "This is enter function"
        return self 

    def __exit__(self,e_t,e_v,t_b):
        print "Now, you are exit"

    def playNow(self):
        print "Now, I am playing"

print "**********"
with WithTest("coolboy") as test:
    print type(test)
    test.playNow() 
    print test.name
print "**********"
```

上述代码运行的结果如下：

```bash
**********
This is enter function
<type 'instance'>
Now, I am playing
coolboy
Now, you are exit
**********
```

分析以上代码： 一二行，执行open函数，该函数返回一个文件对象的实例，然后执行了该实例的`__enter__`函数，该函数返回此实例本身，最后赋值给file变量。从456句可以印证。

自定义的类WithTest，重载了`__enter__`和`__exit__`函数，就可以实现with这样的语法了，注意在`__enter__`函数中，返回了self，在`__exit__`函数中，可以通过`__exit__`的返回值来指示with-block部分发生的异常是否需要reraise，如果返回false，则会reraise with block异常，如果返回ture，则就像什么也没发生。

**上下文管理器contextlib模块对with-as的支持**

contextlib 模块提供了3个对象：装饰器 contextmanager、函数 nested 和上下文管理器 closing。使用这些对象，可以对已有的生成器函数或者对象进行包装，加入对上下文管理协议的支持，避免了专门编写上下文管理器来支持 with 语句。

以contextlib的closing来说，closing帮助实现了`__enter__`和`__exit__`方法，用户不需要自己再实现这两个方法，但是被closing分装的对象必须提供close方法。contextlib.closing类的实现代码如下：

```ruby
class closing(object):
    # help doc here
    def __init__(self, thing):
        self.thing = thing
    def __enter__(self):
        return self.thing
    def __exit__(self, *exc_info):
        self.thing.close()
```

下面是一个使用contextlib.closing的例子：

```javascript
import contextlib
request_url = ('http://www.sina.com.cn/')
with contextlib.closing(urlopen(request_url)) as response:
    return response.read().decode('utf-8')
```

> **针对上面的with的方法进行总结如下几点：**  
> 1.with 后面的方法会得到一个对象，这个对象来自于`__enter__`方法，向上面我们的"Foo"字符串对象、还有return self中得到的WithTest的对象，都是一个实例instance，最终通过as 赋值给后面的变量，在with代码中进行使用  
> 2.with的代码块中最后一行执行结束会执行`__exit__`方法，一般在这个方法中执行文件关闭、异常处理等  
> 3.当我们的with代码块中的代码执行出错的时候，代码执行会直接跳到`__exit__`中执行相对应的代码，最后退出with代码块
> 4.注意上面的例子中WithTest的`__enter__`方法中return self的使用、我们在构造上下文管理器的时候，也可以使用这样的方式。
> 
