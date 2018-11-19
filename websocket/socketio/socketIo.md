## socketIo相关内容
### 概念
* github源码地址：[socket.io](https://github.com/socketio/socket.io)
* api: [api](https://socket.io/docs/server-api/)

关于socket.io我们从源码对应的[官网文档](https://socket.io/docs/)中可以清晰的了解到他的基本使用,下面我们来对官网文档做一个详细的解读：

先来看一个简单的demo吧，服务端代码如下，这段代码应该会与client建立一条websocket连接并传递基本信息：
```
var io = require('socket.io').listen(8080);

io.sockets.on('connection', function (socket) {
    socket.emit('news', { hello: 'world' });
    socket.on('my other event', function (data) {
        console.log(data);
    });
});
```
client端：
```
<script src="/socket.io/socket.io.js"></script>
<script>
    var socket = io.connect('http://localhost:8080');
    socket.on('news', function (data) {
        console.log(data);
        socket.emit('my other event', { my: 'data' });
    });
</script>
```
这里我们主要讲下面几个部分：
### socket请求的header，querystring
### namespace（命名空间）与room(channel)房间
socket.io有连个非常重要的概念namespace与room，他们两个存在的原因就是**分组**，把不同的信息转到不同的分组中。
#### namespace
首先来看下namespace的用法，从上面的例子中我们可以知道，我们都是用io来进行所有的操作，假设我们想要发送数据到所有的socket,我们只需要通过如下方式：
```
io.emit('news', 'Hi I am Mark')
```
未指定namespace的啥时候，会使用/作为默认的namespace。接下来我们可以使用room来将io里面的socket分类多不同的房间中，所以当我们要传递指定的信息到某个房间中的时候，我们可以如下操作：
```
io.to('room001').emit('news', 'Hi I am Mark')
```
那如果我想建立另一个io呢？比如我想建立一个专门处理股票的io,二另一个是专门处理其他的io。这时我们就可以使用namespaces里面的of方法来处理了。
```
var stock_io = io.of('/stock');
var futurs_io = io.of('/futures);
```
然后我们就可以使用stock_io进行我们上面说的所有操作了。你可以把namespaces看成是子io，他可以做io做的所有事情。
#### room
接下来我们就来看一下room，这个room和namespace很像，但是要记住的是rooms 是在 namespaces 底下存在的。

建设你有一个需求，有三个用户在使用你的股票报价信息，其中两个人再看1001这只股票，其外一个再看1002这只股票，这个使用就需要注意了，不能把1002的数据推到1001哪里去。

room这个功能其实就能解决上面的问题，room的意思其实就是说你可以发送数据到你需要的房间，所以上面的两只股票可以看成不同的房间，所以如果有1002的数据要
推送，此时我么只需要推送1002的数据就好了。

怎么加入房间呢？如下：
```
io.on('connection', function(socket){
  socket.join('1002');
});
```
如何传送信息到指定的房间呢？这里需要注意的地方就是，我们这里是使用io来发送信息。上面的简单实例是用socket.emit来发送信息的，因为socket.emit就是针对本socket发送的信息。这里我们要发送信息是针对1002这个房间中的所有socket，所以要用io来发。
```
io.to('1101').emit('xxxx股價')
```
下面总结一下这两者之间的差别：
> 你可以把 namespace 的功能想成可以建立多个io ，不同的子io 可以处理自已的事情，也代表有自已的rooms ，io1与io2两个如果都有room为movie的时候，也不会又是什么影响，因为他们分属不同的io。
### 在 Socket.io 中使用 middleware
在我们平时的web开发中经常会有这种需求：每一个http请求到来之前需要先检查登录状态。通常这个时候我们就会增加一个middleware来进行统一处理。middleware是中间件，用来对http请求进行统一处理的方式。

在socket.io中我们同样会有这种需求的。例如，每当要建立webSocket链接时，我们都有一些需要预处理的事情，其实就可以通过middleware来处理。

socket.io提供了use方法来让我们建立middleware，实例代码如下：
```
var srv = require('http').createServer();
var io = require('socket.io')(srv);
var run = 0;
io.use(function(socket, next){
  run++; // 0 -> 1
  next();
});

var socket = require('socket.io-client')();
socket.on('connect', function(){
  // run == 2 at this time
});
```
上面的代码中，use方法就是一个middleware，在每次webSocket链接建立前都会先执行这个方法。所以假设我们需要先处理日志，加载配置等，我们可以这么处理：
```
io.use(logMiddleWare);
io.use(cacheMiddleWare);

function logMiddleWare(socket, next){
    //  处理 log ~~~
    next();
}

function cacheMiddleWare(socket, next){
    //  获取缓存的配置信息
    next();
}
```
### Sending volatile messages
### acknowledgements
### Broadcasting messages
```
socket.broadcast.to('room').emit('event_name',data)//emit to 'room' except this socket／＊发送消息给room所有的socket client端，除了发送者自己＊／

socket.broadcast.emit('event_name',data)//emit to all sockets except this one／＊发送信息给所有连接到server的client端＊／

io.sockets.in('room').emit('event_name',data)//emit to all clients in a particular room／＊发送消息给room所有的socket client端＊／

io.sockets.emit('event_name',data) //emit an event to all clients／＊发送信息给所有连接到server的client端＊／

io.of('namespace').in('room').emit();// emit an event to all clients in a namespace of a particular room

```


参考：
* http://marklin-blog.logdown.com/posts/2907665-socket-io-talking-island
* http://marklin-blog.logdown.com/posts/2906519
