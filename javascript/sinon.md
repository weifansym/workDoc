# Sinon指南: 使用Mocks, Spies 和 Stubs编写JavaScript测试

> 本文转载自：[众成翻译](http://www.zcfy.cc)
> 译者：[loveky](http://www.zcfy.cc/@loveky)
> 链接：[http://www.zcfy.cc/article/422](http://www.zcfy.cc/article/422)
> 原文：[https://www.sitepoint.com/sinon-tutorial-javascript-testing-mocks-spies-stubs/](https://www.sitepoint.com/sinon-tutorial-javascript-testing-mocks-spies-stubs/)

编写单元测试时遇到的最大的困难之一就是如何测试复杂代码。

在真实的项目中，各种各样的功能会使代码很难被测试。不论是浏览器中的Ajax请求，定时器，日期时间等功能，还是Node.js中的数据库，网络，文件操作。

这些功能之所以很难被测试是因为你不能在代码中控制它们。如果你使用了Ajax，为了让你的代码通过测试，你需要有一个服务器来响应请求。如果你使用了`setTimeout`，你的测试代码将不得不等待定时器被触发。如果你访问了数据库或是网络，也是同样的道理，你需要一个含有正确数据的数据库或是一个网络服务器。

现实生活并不总像有的测试指南中描述的那样简单。但是你是否知道针对这个问题有一个解决方案呢？

**通过使用Sinon，我们可以让测试复杂代码变得不复杂！**

让我们看看该怎么做把。

## 是什么让Sinon这么重要和有用呢？

简单来说，[Sinon](http://sinonjs.org/)允许你把代码中难以被测试的部分替换为更容易测试的内容。

测试一段代码时，你不希望它被测试以外的部分影响。如果有任何外部因素可以影响测试，那么这个测试就会变得更复杂并且很容易失败。

如果你想测试一段发送Ajax的代码，该怎么做呢？你需要运行一个服务器并且确保它返回了你测试需要的数据。这种方式导致准备测试环境变得很复杂，同时也给编写和执行单元测试带来的很大的不便。

如果你的代码依赖时间又会怎么样呢？假设一段代码在执行操作之前要等待1秒钟，该怎么办？你可能会通过`setTimeout`来把测试代码的执行延迟1秒钟，但这样会导致测试变慢。如果这个等待间隔比1秒钟更长呢，比如5分钟。我猜你一定不想在每次执行测试代码前都等上5分钟。

通过使用Sinon，我们可以解决以上这些（还有很多其它的）问题，并降低测试的复杂性。

## Sinon是如何工作的呢?

Sinon通过帮助你很容易的创建所谓的*“测试替身”*来消除复杂性。

像它名字暗示的一样，测试替身用来替换测试中的部分代码。让我们回顾一下前文提到的Ajax例子，我们可以使用一个测试替身取代这个Ajax请求，而不是为这个请求搭建一台服务器。在定时器的例子中，我们可以利用测试替身“在时间上前行”。

这听上去有些不可思议，但是基本概念很简单。由于JavaScript是非常灵活的，我们可以把任何方法替换成其它内容。测试替身只不过是把这个想法更进一步罢了。使用Sinon，我们可以把任何JavaScript函数替换成一个测试替身。通过配置，测试替身可以完成各种各样的任务来让测试复杂代码变得简单。

Sinon将测试替身分为3种类型：

*   [Spies](http://sinonjs.org/docs/#spies), 可以提供函数调用的信息，但不会改变函数的行为

*   [Stubs](http://sinonjs.org/docs/#stubs), 与spies类似，但是会完全替换目标函数。这使得一个被stubbed的函数可以做任何你想要的 —— 例如抛出一个异常，返回某个特定值等等。

*   [Mocks](http://sinonjs.org/docs/#mocks), 通过组合spies和stubs，使替换一个完整对象更容易。

此外，Sinon还提供了其它的辅助方法，虽然它们已经超出了本文的讨论范围：

*   [Fake timers](http://sinonjs.org/docs/#clock), 可以用来穿越时间，例如触发一个`setTimeout`

*   [Fake XMLHttpRequest and server](http://sinonjs.org/docs/#server), 用来伪造Ajax请求和响应

有了这些功能，Sinon就可以帮你解决外部依赖在测试时带来的难题。如果掌握了有效利用Sinon的技巧，你就不再需要任何其它工具了。

## 安装Sinon

首先我们需要安装Sinon。

**Node.js测试:**

1.  通过npm安装Sinon `npm install sinon`

2.  在测试代码中引入Sinon `var sinon = require('sinon');`

**基于浏览器环境的测试:**

1.  你可以通过三种途径获取Sinon： 通过npm `npm install sinon`, 使用 [CDN](https://cdnjs.com/libraries/sinon.js), 或是从 [Sinon官网](http://sinonjs.org/download/)下载。

2.  在测试运行页面引入 `sinon.js`。

## 入门

Sinon有很多功能，但是大部分都是建立在它自身之上。在掌握了一部分之后，你就自然而然的了解下一部分。因此当你学习了Sinon的基础知识并了解各个组件的功能之后，使用Sinon就会变得很容易。

我们通常会在某个被调用的函数给我们测试带来麻烦时使用Sinon。

当使用Ajax时，它可能是`$.get`或是`XMHttpRequest`。当使用定时器时，它可能是`setTimeout`。当使用数据库时，它可能是`mongodb.findOne`。

为了让关于这些被调用函数的讨论变得更简单，我会称它们为*依赖*。我们要测试的方法*依赖*于另一个方法的返回值。

我们可以说，使用Sinon的基本模式就是**使用测试替身替换掉不确定的依赖**。

*   当测试Ajax时，我们把`XMLHttpRequest`替换为一个模拟发送Ajax请求的测试替身。

*   当测试定时器时，我们把`setTimeout`替换为一个伪定时器。

*   当测试数据库访问时，我们把`mongodb.findOne`替换为一个可以立即返回伪数据的测试替身。

让我们看看如何在实践中使用这些模式吧。

## Spies

Spies是Sinon中最简单的功能，其它功能都是建立在它之上。

Spies的主要用途是收集函数调用的信息。你也可以用它验证诸如某个函数是否被调用过之类的断言。

```
var spy = sinon.spy();

// 我们可以像调用函数一样调用一个spy
spy('Hello', 'World');

// 现在我们可以获取关于这次调用的信息
console.log(spy.firstCall.args); //output: ['Hello', 'World'] 
```

`sinon.spy`返回一个`Spy`对象。该对象不仅可以像函数一样被调用，还可以收集每次被调用时的信息。在上边的例子中，`firstCall`属性包含了关于第一次调用的信息，比如`firstCall.args`包含了这次调用传递的参数。

虽然你可以向上例中一样利用`sinon.spy`创建一个匿名spy，但更常见的做法是把一个现有函数替换成一个spy。

```
var user = {
  ...
  setName: function(name){
    this.name = name;
  }
}

// 为setName方法创建一个spy
var setNameSpy = sinon.spy(user, 'setName');

// 现在开始，每次调用这个方法时，相关信息都会被记录下来
user.setName('Darth Vader');

// 通过spy对象可以查看这些记录的信息
console.log(setNameSpy.callCount); //output: 1

// 重要的最后一步，移除spy
setNameSpy.restore();
```

把一个现有函数替换成一个spy与前一个例子相比(译者注：这里指匿名spy的例子)并没有什么特殊之处，除了一个关键步骤：当spy使用完成后，切记把它恢复成原始函数，就像上边例子中最后一步那样。如果不这样做，你的测试可能会出现不可预知的结果。

Spies有很多不同的属性提供了关于调用的各方面信息。[Sinon’s spy documentation](http://sinonjs.org/docs/#spies-api)有一份关于所有可用属性的详细介绍。

在实践中，你可能不会经常用到spy。你往往更多的用到stub。但当你需要验证某个函数是否被调用过时，spy还是很方便的：

```
function myFunction(condition, callback){
  if(condition){
    callback();
  }
}

describe('myFunction', function() {
  it('should call the callback function', function() {
    var callback = sinon.spy();

    myFunction(true, callback);

    assert(callback.calledOnce);
  });
});
```

在本例中，我使用了[Mocha](https://mochajs.org/)作为测试框架，[Chai](http://chaijs.com/)作为断言库。如果你想了解关于它们的更多信息，请参考我之前的文章：[使用Mocha和Chai进行JavaScript单元测试](http://www.sitepoint.com/unit-test-javascript-mocha-chai/)。

[查看CodePen示例](http://codepen.io/SitePoint/pen/XdbYzb/)

## Sinon的断言

在继续往前讨论stubs之前，让我们快速看一看[Sinon的断言](http://sinonjs.org/docs/#assertions)。

在使用spies(和stubs)的大多数情境中，你需要通过某种方式来验证结果。

你可以使用任何类型的断言来验证。在上一个关于callback的例子中，我们使用了Chai提供的`assert`方法来验证值是否为真。

```
assert(callback.calledOnce);
```

这种断言方式的问题在于测试失败时的错误信息不够明确。你只会得到一条类似“false不是true”这样的信息。你可能已经想到了，这样的信息对于确定测试为何会失败并没有什么帮助，你还是不得不查看测试代码来找到哪里出错了。这可不好玩。

为了解决这个问题，我们可以在断言中加入一条自定义的错误信息。

```
assert(callback.calledOnce, 'Callback was not called once');
```

但是为何不用*Sinon提供的断言呢*？

```
describe('myFunction', function() {
  it('should call the callback function', function() {
    var callback = sinon.spy();

    myFunction(true, callback);

    sinon.assert.calledOnce(callback);
  });
});
```

像这样使用Sinon的断言可以为我们提供一种更加友好的错误信息。当你需要验证更复杂的情况，比如某个函数的调用参数时，这将变得非常有用。

以下是另外一些Sinon提供的实用断言：

*   `sinon.assert.calledWith` 可以用来验证某个函数被调用时是否传入了特定的参数(这很可能是我最常用的了)

*   `sinon.assert.callOrder` 用来验证函数是否按照一定顺序被调用

和spies一样，[Sinon的断言文档](http://sinonjs.org/docs/#assertions)列出了所有可用的选项。如果你习惯使用Chai，那么有一个[sinon-chai插件](https://github.com/domenic/sinon-chai)可供选择，它可以让你通过Chai的`expect`和`should`接口使用Sinon的断言。

## Stubs

由于其灵活和方便，stubs成为了Sinon中最常用的测试替身类型。它拥有spies提供的所有功能，区别在于它会完全替换掉目标函数，而不只是记录函数的调用信息。换句话说，当使用spy时，原函数还会继续执行，但使用stub时就不会。

这使得stubs非常适用于以下场景：

*   替换掉那些使测试变慢或是难以测试的外部调用

*   根据函数返回值来触发不同的代码执行路径

*   测试异常情况，例如代码抛出了一个异常

我们可以用类似创建spies的方法创建stubs：

```
var stub = sinon.stub();

stub('hello');

console.log(stub.firstCall.args); // 输出: ['hello']
```

我们可以创建匿名stubs，就和使用spies时一样，但只有当你用stubs替换一个现有函数时它才开始真正的发挥作用。

举例来说，如果你有一段代码使用了jQuery的Ajax功能，那这段代码就会很难被测试。这段代码会向某台服务器发送请求，你不得不保证测试期间该服务器的可用性。或者你可能会想到在代码里增加一段特殊逻辑以便在测试环境下不会真正的发送请求 —— 这可犯了大忌。在绝大多数情况下你应该保证代码中不会出现针对测试环境的特殊逻辑。

我们可以通过Sinon把Ajax功能替换为一个stub，而不是寻求其它糟糕的实践方式。这会使得测试变得很简单。

以下是一个我们要测试的函数。它接受一个对象作为参数，并通过Ajax把该对象发送给某个预定的URL。

```
function saveUser(user, callback) {
  $.post('/users', {
    first: user.firstname,
    last: user.lastname
  }, callback);
}
```

通常情况下，由于涉及到Ajax调用和某个特定的URL，对它测试是比较困难的。但如果我们使用了stub，这就变得很简单。

比方说我们要确保传给`saveUser`的回调函数在请求结束后被正确执行。

```
describe('saveUser', function() {
  it('should call callback after saving', function() {

    // 我们会stub $.post，这样就不用真正的发送请求
    var post = sinon.stub($, 'post');
    post.yields();

    // 针对回调函数使用一个spy
    var callback = sinon.spy();

    saveUser({ firstname: 'Han', lastname: 'Solo' }, callback);

    post.restore();
    sinon.assert.calledOnce(callback);
  });
});
```

[查看CodePen示例](http://codepen.io/SitePoint/pen/vGOrwj/)

这里我们把Ajax方法替换成了一个stub。这意味着代码里并不会真的发出请求，因此也就不需要相应的服务器了，这样我们就对测试代码里的逻辑取得了完全控制。

由于我们要确保传给`saveUser`的回调函数被执行了，我们指示stub要*yield*。这意味着stub会自动执行作为参数传入的第一个函数。这就模拟了`$.post`的行为 —— 请求一旦完成就执行回调函数。

除了stub以外，我们还在测试中创建了一个spy。我们也可以使用一个普通函数作为回调，但是使用了spy后利用Sinon提供的`sinon.assert.calledOnce`断言可以很容易的验证结果。

在使用stub的大多数情况中，你都可以遵循以下模式：

*   找到导致问题的函数，比如`$.post`

*   观察它是如何工作的以便在测试中模拟它的行为

*   创建一个stub

*   配置stub以便按照你期望的方式工作

Stub不必模拟目标对象的所有行为。只要模拟在测试中用到的行为就够了，其它的都可以忽略。

Stub的另一个常见使用场景是验证某个函数被调用时传入了正确的参数。

例如，针对Ajax的功能，我们想验证发送的数据正确与否。因此，我们可以这样：

```
describe('saveUser', function() {
  it('should send correct parameters to the expected URL', function() {

    // 像之前一样为$.post设置stub
    var post = sinon.stub($, 'post');

    // 创建变量保存我们期望看到的结果
    var expectedUrl = '/users';
    var expectedParams = {
      first: 'Expected first name',
      last: 'Expected last name'
    };

    // 创建将要作为参数的数据
    var user = {
      firstname: expectedParams.first,
      lastname: expectedParams.last
    }

    saveUser(user, function(){} );
    post.restore();

    sinon.assert.calledWith(post, expectedUrl, expectedParams);
  });
});
```
[查看CodePen示例](http://codepen.io/SitePoint/pen/eZNKqZ/)

同样的，我们又为`$.post()`创建了一个stub，但这次我们没有设置它为yield。这是因为此次的测试我们并不关心回调函数，因此设置yield就没有意义了。

我们创建了一些变量用来保存期望得到的数据 —— URL和参数。创建这样的变量是一种不错的实践，因为这样就可以很容易看出这个测试要测哪些数据。我们还可以利用这些值创建`user`变量从而避免重复输入。

这次我们使用了`sinon.assert.calledWith()`断言。我们把stub作为第一个参数传入，因为我们要验证这个stub被调用时是否传入了正确的参数。

还有另外一种使用Sinon测试Ajax请求的方法。那就是使用Sinon提供的伪XMLHttpRequest功能。我们不会在本文中详细阐述，如果你想了解更多详细信息，请参考[使用Sinon的伪XMLHttpRequest测试Ajax](http://codeutopia.net/blog/2015/03/21/unit-testing-ajax-requests-with-mocha/)。

## Mocks

Mocks是使用stub的另一种途径。如果你曾经听过“mock对象”这种说法，这其实是一码事 —— Sinon的mock可以用来替换整个对象以改变其行为，就像函数stub一样。

基本上只有需要针对一个对象的多个方法进行stub时才需要使用mock。如果只需要替换一个方法，使用stub更简单。

使用mock时你要很小心。由于mock强大的功能，它很容易导致你的测试过于具体 —— 测试了太多，太细节的内容 —— 这很容易在不经意间导致你的测试变得脆弱。

与spy和stub不同的是，mock有内置的断言。你需要预先定义好mock对象期望的行为并在测试结束前执行验证函数。

比方说我们代码中使用了[store.js](https://github.com/marcuswestin/store.js)向localStorage中写入数据，我们希望测试一个与这部分内容相关的函数。我们可以使用一个mock来协助测试：

```
describe('incrementStoredData', function() {
  it('should increment stored value by one', function() {
    var storeMock = sinon.mock(store);
    storeMock.expects('get').withArgs('data').returns(0);
    storeMock.expects('set').once().withArgs('data', 1);

    incrementStoredData();

    storeMock.restore();
    storeMock.verify();
  });
});
```

[查看CodePen示例](http://codepen.io/SitePoint/pen/EKjpYW/)

使用mock时，我们使用链式调用的方式定义一系列方法以及相应的返回值。除了预先定义好行为并在测试结束前调用`storeMock.verify()`来验证结果以外，这和使用断言验证测试结果没什么两样。

在Sinon的mock对象术语中，执行`mock.expects('something')`创建了一个*预期*。例如，函数`mock.something()`期望被调用。每一个预期除了mock特殊的功能外，还支持spy和stub的功能。

你可能会发现大多数时候使用stub比使用mock简单的多，这很正常。Mock应该被小心的使用。

如果想要了解mock支持的所有方法，请参考[Sinon’s mock documentation](http://sinonjs.org/docs/#mocks)。

## 重要的最佳实践: 使用sinon.test()

无论何时使用spy，stub或是mock，都有一条重要的最佳实践需要牢记。

如果你使用测试替身替换了一个现有函数，记得使用`sinon.test()`。

在前面的示例中，我们使用了`stub.restore()`或`mock.restore()`来执行清理操作。这个操作是必要的，否则测试替身会一直存在并给其它测试带来负面影响或是导致错误。

但是直接使用`restore()`方法是有问题的。因为有可能在`restore()`执行之前测试代码就因为错误提前结束执行了。

有两种方法可以解决这个问题：我们可以把所有的代码放在一个`try catch`块中，这样就可以在`finally`块中执行`restore()`而不用担心测试代码是否报错。

还有一种更好的方式，就是把测试代码包裹在`sinon.test()`中：

```
it('should do something with stubs', sinon.test(function() {
  var stub = this.stub($, 'post');

  doSomething();

  sinon.assert.calledOnce(stub);
});
```

在上边的示例中，要注意的是传递给`it()`的第二个参数被包裹在`sinon.test()`中。另一个要注意的点是我们使用的是`this.stub()`而不是`sinon.stub()`。

把测试代码包裹在`sinon.test()`之中后，我们就可以使用Sinon的*沙盒*特性了。它允许我们通过`this.spy()`，`this.stub()`和`this.mock()`来创建spy，stub和mock。任何使用沙盒特性创建的测试替身都会被自动清理。

注意上边的例子中没有`stub.restore()`操作 —— 因为在沙盒特性下的测试里它变得不必要了。

如果在所有地方都使用了`sinon.test()`，那么你就可以避免由于某个测试未能清理它内部的测试替身而导致后续测试随机失败的情况。

## Sinon没有魔法

Sinon功能强大，可能看上去很难理解它是如何工作的。为了更好的理解Sinon的工作原理，让我们看一些和它工作原理有关的例子。这将有利于我们更好的理解Sinon究竟做了哪些工作并在不同的场景中更好的利用它。

我们也可以手工创建spy，stub或是mock。使用Sinon的原因在于它使得这个过程更简单了 —— 手工创建通常比较复杂。不过为了理解Sinon，还是让我们看看如何手工创建吧。

首先，spy在本质上就是一个函数包装器：

```
// 一个简单的spy辅助函数
function createSpy(targetFunc) {
  var spy = function() {
    spy.args = arguments;
    spy.returnValue = targetFunc.apply(this, arguments);
    return spy.returnValue;
  };

  return spy;
}

// 基于一个函数创建spy:
function sum(a, b) { return a + b; }

var spiedSum = createSpy(sum);

spiedSum(10, 5);

console.log(spiedSum.args); // 输出: [10, 5]
console.log(spiedSum.returnValue); // 输出: 15 
```

我们可以使用一个像这样的方法很容易的创建spy。但是要明白Sinon的spy提供了包括断言在内的丰富得多的功能，这使得使用Sinon相当容易。

### 那么Stub呢?

要创建一个简单的stub，只需把一个函数替换成另一个：

```
var stub = function() { };

var original = thing.otherFunction;
thing.otherFunction = stub;

// 现在开始，所有对thing.otherFunction的调用都会被stub的调用所取代
```

但同样需要指出的是，Sinon的stub有若干优势：

*   它们包含了spy的所有功能

*   你可以使用`stub.restore()`轻松的恢复原始函数

*   你可以针对Sinon stub使用断言

Mock只不过是把spy和stub组合在一起，使得可以灵活使用它们的功能。

虽然Sinon某些时候看起来使用了很多“魔法”，但大多数情况下，你都可以使用自己的代码实现相同的功能。与自己开发一个库比起来，使用Sinon只不过是更方便罢了。

## 总结

对现实中的代码进行测试有时看起来会非常复杂，你很容易就彻底放弃了。但是有了Sinon的帮助以后，对几乎所有类型代码的测试都变得很简单。

你只需要记住主要原则：如果一个函数让你的测试难于编写，尝试把它替换成一个测试替身。这条原则对任何函数都适用。

还想了解更多关于如何在你的代码中使用Sinon的信息吗？请访问我的网站，我会给你发送免费的[现实中的Sinon指南](http://codeutopia.net/)。它包括了Sinon最佳实践以及三个如何在真实场景中使用Sinon的实例！
