## Airbnb JavaScript Style Guide() {

[Airbnb javascript style](https://github.com/airbnb/javascript#airbnb-javascript-style-guide-):用合理的方式编写JavaScript

## Table of Contents内容列表)
### Types
基本类型：当操作基本类型值的时候，是只对这个值的操作，基本类型包括：string，number，boolean，null，undefined，实例如下：

    const foo = 1;
    let bar = foo;
    bar = 9;
    console.log(foo, bar); // => 1, 9
复杂类型：当操作复杂类型的时候，是对这个值得引用操作，复杂类型包括：object，array，function
    
    const foo = [1, 2];
    const bar = foo;

    bar[0] = 9;
    console.log(foo[0], bar[0]); // => 9, 9
### References(变量)
使用cont声明所有的变量，避免是有var，为什么要这么做呢，这样可以保证不能对变量重新分配以及减少难以理解的代码，以及bug

    // bad
    var a = 1;
    var b = 2;

    // good
    const a = 1;
    const b = 2;
    
如果你必须对变量进行重分配，那就是使用let代替var进行声明，这是为什么呢？因为let属于块级作用域而var是则是函数作用域

    // bad
    var count = 1;
    if (true) {
     count += 1;
    }

    // good, use the let.
    let count = 1;
    if (true) {
      count += 1;
    }
注意let与const都是块级作用域

    // const and let only exist in the blocks they are defined in.
    {
     let a = 1;
     const b = 1;
    }
    console.log(a); // ReferenceError
    console.log(b); // ReferenceError    

### Objects(对像)
1、使用对象字面量进行创建对象

    // bad
    const item = new Object();

    // good
    const item = {};
2、在创建带有动态属性的对象的时候使用计算属性，因为他允许你在同一个地方定义对象的属性

    function getKey(k) {
      return `a key named ${k}`;
    }

    // bad
    const obj = {
      id: 5,
      name: 'San Francisco',
    };
    obj[getKey('enabled')] = true;

    // good
    const obj = {
      id: 5,
      name: 'San Francisco',
      [getKey('enabled')]: true,
    };
3、使用简写的对象方法

    // bad
    const atom = {
      value: 1,

      addValue: function (value) {
        return atom.value + value;
      },
    };

    // good
    const atom = {
      value: 1,

      addValue(value) {
        return atom.value + value;
      },
    };
4、使用属性的简写形式，这样在描述和写的时候都比较方便

    const lukeSkywalker = 'Luke Skywalker';

    // bad
    const obj = {
      lukeSkywalker: lukeSkywalker,
    };

    // good
    const obj = {
      lukeSkywalker,
    };
5、把简写形式的属性放在定义对象的开头，这样很容易区分什么样的形式才是速记形式

    const anakinSkywalker = 'Anakin Skywalker';
    const lukeSkywalker = 'Luke Skywalker';

    // bad
    const obj = {
      episodeOne: 1,
      twoJediWalkIntoACantina: 2,
      lukeSkywalker,
      episodeThree: 3,
      mayTheFourth: 4,
      anakinSkywalker,
    };

    // good
    const obj = {
      lukeSkywalker,
      anakinSkywalker,
      episodeOne: 1,
      twoJediWalkIntoACantina: 2,
      episodeThree: 3,
      mayTheFourth: 4,
    };
6、只引用无效标识符的属性，通常我们人为这种形式容易读，提供语法高亮，而且更容易在其他js引擎下优化

    // bad
    const bad = {
      'foo': 3,
      'bar': 4,
      'data-blah': 5,
    };

    // good
    const good = {
      foo: 3,
      bar: 4,
      'data-blah': 5,
    };

 7、不要直接调用**Object.prototype**方法 ，例如：hasOwnProperty, propertyIsEnumerable, and isPrototypeOf等，因为这些方法可能会被有关对象的属性所跟踪。想想下**{ hasOwnProperty: false }**或者对象可能是一个空对象**(Object.create(null))**
 
     // bad
     console.log(object.hasOwnProperty(key));

     // good
     console.log(Object.prototype.hasOwnProperty.call(object, key));

     // best
     const has = Object.prototype.hasOwnProperty; // cache the lookup once, in module scope.
     /* or */
     import has from 'has';
     // ...
     console.log(has.call(object, key)); 
8、再使用Object.assign对对象进行前拷贝的时候，更喜欢使用对象拓展操作符（rest operator）获得一个具有某些属性的新对象
 
    // very bad
    const original = { a: 1, b: 2 };
    const copy = Object.assign(original, { c: 3 }); // this mutates `original` ಠ_ಠ
    delete copy.a; // so does this

    // bad
    const original = { a: 1, b: 2 };
    const copy = Object.assign({}, original, { c: 3 }); // copy => { a: 1, b: 2, c: 3 }

    // good
    const original = { a: 1, b: 2 };
    const copy = { ...original, c: 3 }; // copy => { a: 1, b: 2, c: 3 }

    const { a, ...noA } = copy; // noA => { b: 2, c: 3 }
    
### 数组
1、使用字面值创建数组

    // bad
    const items = new Array();

    // good
    const items = [];
2、向数组添加元素时使用 Arrary#push 替代直接赋值。

    const someStack = [];
    // bad
    someStack[someStack.length] = 'abracadabra';
    // good
    someStack.push('abracadabra');
3、使用拓展运算符 ... 复制数组。

    // bad
    const len = items.length;
    const itemsCopy = [];
    let i;

    for (i = 0; i < len; i++) {
      itemsCopy[i] = items[i];
    }

    // good
    const itemsCopy = [...items];
   4、使用 Array#from 把一个类数组对象转换成数组。
   
       const foo = document.querySelectorAll('.foo');
       const nodes = Array.from(foo);
