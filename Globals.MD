## node 之Globals
**Global object**这个对象在所有的其他模块中都可以访问,需要注意的是：他们中的一些属性其实在全局作用域中是不能够访问的，
但是在模块作用域中是可以访问的。
### Class: Buffer
* [Function](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/Function)

用来处理二进制数据，详细信息查看[buffer section](http://nodejs.cn/doc/node/buffer.html)
### __dirname
* string
当前运行脚本的目录，例如，运行**node example.js**在**/Users/mjr**目录下
 
   console.log(__dirname);
    // /Users/mjr
