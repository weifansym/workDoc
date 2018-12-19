## bluebird 常用文档
官网文档：[bluebird](http://bluebirdjs.com/docs/api-reference.html)
* Promise.join: 用法类似于Promise.all，用来控制并发
* Promise.props: Promise.props，只不过参数是对象，而不是数组
* Promise.some: 参数是一个可枚举的类型，和一个count，这个count就是枚举类型执行为fulfilled状态的数量，具体请看官方文档。
* Promise.any: 类似于Promise.some，只是把Promise.some的count变成了1，然后返回值是一个值，而不是一个是数组。
### Promise.map
```
Promise.map(
  Iterable<any>|Promise<Iterable<,
  function(any item, int index, int length) mapper,
  [Object {concurrency: int=Infinity} options]
 ) -> Promise
```
为指定的Iterable数组或Iterablepromise，执行一个处理函数mapper并返回执行后的数组或Iterablepromise。

Promises会等待mapper全部执行完成后返回，如果数组中的promise执行全部分成功则Promises中是执行成功值。如果任何一个promise执行失败，Promises对应的也是拒绝值。

Promise.map可以用于替数组.push+Promise.all方法：
```
// 对于如下一个操作：
var promises = [];
for (var i = 0; i < fileNames.length; ++i) {
    promises.push(fs.readFileAsync(fileNames[i]));
}
Promise.all(promises).then(function() {
    console.log("done");
});

// 使用 Promise.map处理如下：
Promise.map(fileNames, function(fileName) {
    // Promise.map 等待操作成功后返回
    return fs.readFileAsync(fileName);
}).then(function() {
    console.log("done");
});
```
#### Map Option: concurrency
你可以指定一个可选的并发限制：
```
...map(..., {concurrency: 3});
```
这里的并发限制主要是指定创建promise的数量。例如如果并发是3，map中指定的执行函数并发执行的数量是3个，
只要3各种
