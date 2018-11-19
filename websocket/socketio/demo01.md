## namespace和room相关demo
### namespace的demo
app.js
```
var io = require('socket.io')(8888);
var chat = io
  .of('/chat')
  .on('connection', function (socket) {
    console.log('server chat connection');
    socket.emit('a message', {
      that: 'only',
      '/chat': 'will get'
    });
    /*//  广播给每一个用户
    chat.emit('a message', {
      everyone: 'in'
      , '/chat': 'will get'
    });*/
    socket.broadcast.emit('a message', {
      that: 'broadcast'
    })
    socket.on('hi', (data) => {
      console.log('hi data: ', data);
    })
  });

var news = io
  .of('/news')
  .on('connection', function (socket) {
    console.log('server news connection');
    const num = Math.floor(10000 * Math.random());
    socket.emit('item', { news: 'item' + num });
  });
```
client.js
```
var io = require('socket.io-client');
var chat = io('http://localhost:8888/chat');
let news = io('http://localhost:8888/news');

const num = Math.floor(10000 * Math.random());

chat.on('connect', function () {
  console.log('client chat connect');
  chat.emit('hi', 'name: ' + num);
});

chat.on('a message', function (msg) {
  console.log('a meesage: ', msg);
})

news.on('connect', function () {
  console.log('client news connect');
  news.on('item', function (item) {
    console.log('news item: ', item);
  });
  // news.emit('woot');
});
```

### room的demo
app.js
```
const io = require('socket.io')(8889);

const chat = io.of('/chat');

var count = 0;

/*chat.on('connect', (socket) => {
  console.log('serRoom chat start');

  var roomNum = 'room' + (++count);//自增实现socket进入不同房间

  socket.join(roomNum, function () {
    console.log('has join: ', socket.rooms);
  }); //加入房间后，打印出socket和room的信息

  socket.on('chat message', function (msg) {
    var room = Object.keys(socket.rooms)[1]; //这是当前socket的房间，这个对象设置得有点怪，但是事实如此。
    chat.to(room).emit('chat message', room);
    console.log('message: ', msg);
    console.log(room);//打印出房间。
  });
});*/

const room = 'testRoom';

chat.on('connect', (socket) => {
  console.log('serRoom chat start');
  const room = 'room' + (++count);
  socket.join(room, (err) => {
    console.log(socket.id, 'join room: ', socket.rooms);
  });

  socket.on('chat message', (msg) => {
    var room = Object.keys(socket.rooms)[1]; //这是当前socket的房间，这个对象设置得有点怪，但是事实如此。
    chat.to(room).emit('chat message', room);   //  推送数据给房间中(只在此房间中)的所有用户，也包括自己
    // socket.to(room).emit(('chat message'), room);   //  广播房间(其他房间的用户无法收到)里的其他用户可以收到，除了自己
    // socket.broadcast.emit('chat message', room);  //  基于命名空间的广播,命名空间中的其他房间用户也可以收到消息
    console.log('message: ', msg);
    console.log(room);//打印出房间。
  });

})
```
client.js
```
const io = require('socket.io-client');
const chat = io('http://localhost:8889/chat');

chat.on('connect', () => {
  console.log('cli room has starts');
  chat.emit('chat message', 'cliRoom send message!!!');
  /*chat.on('chat message', (msg) => {
    console.log('server msg: ', msg);
  });*/
});

chat.on('chat message', (msg) => {
  console.log('server msg: ', msg);
});
```
