## 
由于nodejs是单线程模式，一旦程序遇到未捕获的异常后，整个程序就崩溃了。

* 有两个解决办法，用其他守护进程模块来启动nodejs项目，例如forever,pm2
具体是使用请看相关包的是使用。
* 还有一种是用 worker 进程来启动server。master进程负责维护和管理 worker进程。这样master进程就相当于守护进程了。具体查看官方案例
[官方代码](https://nodejs.org/api/cluster.html)如下:
```
const cluster = require('cluster');
const http = require('http');
const numCPUs = require('os').cpus().length;

if (cluster.isMaster) {
  console.log(`Master ${process.pid} is running`);

  // Fork workers.
  for (let i = 0; i < numCPUs; i++) {
    cluster.fork();
  }

  cluster.on('exit', (worker, code, signal) => {
    console.log(`worker ${worker.process.pid} died`);
  });
} else {
  // Workers can share any TCP connection
  // In this case it is an HTTP server
  http.createServer((req, res) => {
    res.writeHead(200);
    res.end('hello world\n');
  }).listen(8000);

  console.log(`Worker ${process.pid} started`);
}

```
这种方式是多进程模式，也是多实例模式。如果启动了4个worker进程，相当于运行了4个server，而这4个进程都监听同一个端口，请求会分配到任意worker进程中，
如果有一个进程挂了，那么还有其他进程继续工作。master主进程则再次启动新的worker进程，保证workder进程数量

使用多进程模式，会是程序变得更为复杂（比如进程间数据传输，资源共享问题，并发）如果只是小项目可以只启动一个worker进程
