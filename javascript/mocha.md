## mocha
mocha官网地址如下：https://mochajs.org/#getting-started
## 关于单元测试的想法
对于一些比较重要的项目，每次更新代码之后总是要自己测好久，担心一旦上线出了问题影响的服务太多，此时就希望能有一个比较规范的测试流程。在github上看到牛逼的javascript开源项目，也都是有测试代码的，看来业界大牛们都比较注重单元测试这块。
就我自己的理解而言：

* 涉及到大量业务逻辑的代码，可能我没有精力去给每个函数都写上单元测试的代码，功能细节的测试应该交由测试的同事去完成，但是对会直接影响项目正常运行的重要的数据接口，还是可以看情况写上几个单元测试用例的，每一次修改之后跑一跑用例测试一下。
* 重要的框架底层模块，任何地方出一个小问题，都可能影响到很多服务。对于这种模块，最好是每个函数、每种接口都写上单元测试代码，不然一出问题就是一个大坑啊。
*  开放出去的公共模块，可以针对主要的函数和接口写上单元测试代码，这样可以确保模块代码比较健壮，看起来也专业一些：）。

基于以上几个想法，我决定学习一款Javascript单元测试框架，并试试去使用它写一些单元测试的代码。
看了很多技术站点和博客的文章，参考了一部分开源项目的测试代码，大致观望下风向，决定学习一下mocha.js这款单元测试框架。
别人的文章都是别人自己学习、咀嚼理解出来的内容，想学的透彻一点，还是自己学习并翻译一遍原版官方的文档比较好。

## mocha guide
mocha是一款功能丰富的javascript单元测试框架，它既可以运行在nodejs环境中，也可以运行在浏览器环境中。
javascript是一门单线程语言，最显著的特点就是有很多异步执行。同步代码的测试比较简单，直接判断函数的返回值是否符合预期就行了，而异步的函数，就需要测试框架支持回调、promise或其他的方式来判断测试结果的正确性了。mocha可以良好的支持javascript异步的单元测试。
mocha会串行地执行我们编写的测试用例，可以在将未捕获异常指向对应用例的同时，保证输出灵活准确的测试结果报告。

# INSTALLATION(安装)
* 使用**npm**进行安装:
* 全局安装：

    $ npm install --global mocha
* 作为你项目的开发依赖安装：
   
    $ npm install --save-dev mocha
* 使用**yarn**安装：

    $ yarn add --save-dev mocha
## GETTING STARTED（开始使用）
先看一个小栗子：
 
    $ npm install mocha
    $ mkdir test
    $ $EDITOR test/test.js # or open with your favorite editor
在你的编辑器中输入：

    var assert = require('assert');
    describe('Array', function() {
      describe('#indexOf()', function() {
        it('should return -1 when the value is not present', function() {
          assert.equal(-1, [1,2,3].indexOf(4));
        });
      });
    });
回到控制台

    $ ./node_modules/mocha/bin/mocha

      Array
        #indexOf()
          ✓ should return -1 when the value is not present

      1 passing (9ms)
在你的package.json中设置测试脚本：

    "scripts": {
        "test": "mocha"
      }
 运行测试：
 
     $ yarn test
     
 ## ASSERTIONS(断言)    
 
