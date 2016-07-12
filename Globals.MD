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
**__dirname**不是一个全局目录但是对于每一个模块来说是一个本地目录。  
例如：有两个模块**a**和**b**,其中**b**是**a**的依赖，他们具有如下的目录结构：
* /Users/mjr/app/a.js
* /Users/mjr/app/node_modules/b/b.js

在**b.js**文件中引用**__dirname**，将会返回**/Users/mjr/app/node_modules/b**。在**a**中引用**__dirname**将会返回**/Users/mjr/app**

