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
这里的并发限制主要是指定创建promise的数量。例如如果并发是3，map中指定的执行函数并发执行的数量是3个，只要3个中有一个执行为resolve状态，才会有剩余的其他函数执行。

map方法在数组元素上并没有指定顺序，所以不保证顺序，要是需要顺序执行就需要使用Promise.mapSeries方法。
### Promise.reduce
```
Promise.reduce(
  Iterable<any>|Promise<Iterable<,
  function(any accumulator, any item, int index, int length) reducer,
  [any initialValue]
) -> Promise
```
为指定的Iterable数组或Iterablepromise，执行一个处理函数reducer，并返回一个经过Reduce处理后promise。
reducer函数会最终返回一个promise。数组中的所有promise处理成功后，会返回一个成功状态的promise，任何一个执行失败都会返回一个拒绝状态的promise。
如，使用Promise.reduce计算从三个文件中读取值的总和，每个文件中都有一个数字10：
```
Promise.reduce(["file1.txt", "file2.txt", "file3.txt"], function(total, fileName) {
  return fs.readFileAsync(fileName, "utf8").then(function(contents) {
    return total + parseInt(contents, 10);
  });
}, 0).then(function(total) {
  //Total is 30
});
```
### Promise.filter
```
Promise.filter(
  Iterable<any>|Promise<Iterable<,
  function(any item, int index, int length) filterer,
  [Object {concurrency: int=Infinity} options]
) -> Promise
```
为指定的Iterable数组或Iterablepromise，执行一个过滤函数filterer，并返回经过筛选后promises数组。
```
var Promise = require("bluebird");
var E = require("core-error-predicates");
var fs = Promise.promisifyAll(require("fs"));

fs.readdirAsync(process.cwd()).filter(function(fileName) {
  return fs.statAsync(fileName)
    .then(function(stat) {
      return stat.isDirectory();
    })
    .catch(E.FileAccessError, function() {
      return false;
    });
}).each(function(directoryName) {
  console.log(directoryName, " is an accessible directory");
});
```
### Promise.each
```
Promise.each(
  Iterable<any>|Promise<Iterable>,
  function(any item, int index, int length) iterator
) -> Promise
```
为指定的Iterable数组或Iterablepromise，执行一个函数iterator，该函数参数为(value, index, length)，value输入数组中promise的resolved值。
iterator函数会最终返回一个promise。数组中的所有promise处理成功后，会返回一个成功状态的promise，任何一个执行失败都会返回一个拒绝状态的promise。
### Promise.mapSeries
顺序执行的Map操作
```
Promise.mapSeries(
  Iterable<any>|Promise<Iterable<,
  function(any item, int index, int length) mapper
) -> Promise
```
详见Promise.map，但会按数组顺序依次执行mapper。

