## Js在使用async/await时的并发情况
### 0. 前言
Js的异步一直都是很难处理的一个东西，即便到了最新的ES7，正式引入了async/await能将异步语法糖封装到类同步的写法，里面的很多细节及原理仍旧需要我们非常小心地学习和处理。这次遇到了一个比较傻的问题，而我居然都楞了一下，事后觉得有必要记录下来。

问题如下：

一个代码片逻辑中，如果使用await来等待异步的返回，代码逻辑是继续执行下去了（异步）还是等待在这个await语句上没有执行下去？
如果是等待的话，那么当前的node进程还能不能处理其他任务了？是整个进程阻塞了？还是当前的业务阻塞了？
### 1. 问题1
问题1的答案显而易见，也很容易理解，async/await的出现本来就是为了解决异步语法的理解问题的，因此在使用await进行返回等待的时候当前的逻辑是停止等待状态。

查看2个比对例子：

await 例子：
```
const waitForResponse = function() {
  return new Promise((resolve) => {
    setTimeout(() => {
      resolve('Done');
    }, 2000);
  });
};

const main = async function() {
  console.log('Before');
  console.log(await waitForResponse());
  console.log('After');
};

main().then(_ => _);

// Result:
// Before
// Done
// After
```
promise 例子：
```
const waitForResponse = function() {
  return new Promise((resolve) => {
    setTimeout(() => {
      resolve('Done');
    }, 2000);
  });
};

const main = async function() {
  console.log('Before');
  waitForResponse().then((response) => {
    console.log(response);
  });
  console.log('After');
};

main().then(_ => _);

// Result:
// Before
// After
// Done
```
### 2. 问题2
问题2的结论也非常简单，任何并发请求都会进入到事件循环中，并不会受到await语句的阻塞。await语句的停顿只会在逻辑层面，而不会在进程层面。
并发请求不会受到任何影响，这是node赖以生存的关键核心，不可能受到干扰。

例子：
```
// server.js
const Koa = require('koa');
const app = new Koa();
require('koa-qs')(app);

async function response(index) {
  return new Promise((resolve, reject) => {
    if (parseInt(index) === 1 || parseInt(index) === 2) { // magic number
      setTimeout(() => {
        resolve(`Response ${index} waited ${index}s`);
      }, index * 1000); // wait Xs
    } else {
      resolve(`Response ${index}`);
    }
  });
}

app.use(async ctx => {
  ctx.body = await response(ctx.query.index);
});

app.listen(3000);

// client.js
const fetch = require('node-fetch');

function sendRequest(index) {
  fetch(`http://localhost:3000/?index=${index}`)
    .then(function(res) {
      return res.text();
    }).then(function(body) {
      process.send(`Child response of ${index}: ${body}`);
    });
}

process.on('message', function(m) {
  process.send(`Child got msg: ${m}`);
  sendRequest(m);
});
process.send('Child initialized!');

// cluster.js
const cp = require('child_process');
const fork = cp.fork;

const clients = [];
for (let i = 0; i < 5; i++) {
  let client = fork('./client.js');
  client.on('message', (msg) => {
    console.log(`Client[${i}]: ${msg}`);
  });
  clients.push(client);
}

setTimeout(() => {
  console.log('Clients initialized!');
  clients.forEach((client, i) => {
    client.send(i);
  });
}, 2000); // wait 2s for clients to be initialized

// Execution:
// node server.js // 启动服务器
// node cluster.js // 查看结果

// Result:
// Client[1]: Child initialized!
// Client[0]: Child initialized!
// Client[2]: Child initialized!
// Client[3]: Child initialized!
// Client[4]: Child initialized!
// Clients initialized!
// Client[2]: Child got msg: 2
// Client[1]: Child got msg: 1
// Client[3]: Child got msg: 3
// Client[4]: Child got msg: 4
// Client[0]: Child got msg: 0
// Client[3]: Child response of 3: Response 3
// Client[0]: Child response of 0: Response 0
// Client[4]: Child response of 4: Response 4
// Client[1]: Child response of 1: Response 1 waited 1s
// Client[2]: Child response of 2: Response 2 waited 2s 
```
可以看到 0、3、4的返回结果几乎是同一个时间节点出来的，而1和2则各等待了1秒和2秒
