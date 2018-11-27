## node.js socket
废话不多说，先谈谈优、缺点
* 优点
高性能
实时性
加密方便
* 缺点
需要对数据处理
对开发水平要求较高
### 简单介绍下优缺点
* 高性能: 传输多次数据时，只需一次链接。你以为ajax很快吗，跟socket比，慢的一批
* 实时性: 服务端也能主动推送消息给客户端，你ajax试试？
* 加密方便: 字符串传输效率多低，我就不说了。转buffer或者aes等其他加密，可以让你节省70%以上带宽，谁用谁知道

### 开发过程中，可以参考以下几点经验
1. 要严格制定数据格式，特别在跨语言传输时
2. 可以把数据转为buffer或可逆加密传输节省服务器资源
3. 需要处理粘包、缓冲区数据截断问题
### demo的数据包结构为：
长度包头(Buffer(4) Int32LE) + 数据包(Buffer)

每个包的头4个字节，为数据包的字节长度 这么做的原因是：socket原生开发中，取缓冲区数据是一次性全部取出，但可能缓冲区数据不完整或有多个包，
所以包头返回这个包的数据长度，也就是上面说的经验3。

该示例代码逻辑为：客户端连接后，就准备一次数据，然后发送给服务端，服务端收到数据后，解析完原样返回给客户端

直接上代码(只贴客户端，服务端逻辑相同，自己改改)
```
var net = require('net');
var _maps = require("./maps");
//_maps这里是自己导入的一个非常非常长的数据，用来测试数据截断使用

//你要连接的服务端ip跟端口，看不懂的自杀
var HOST = 'localhost';
var PORT = 8888;

var client = new net.Socket();
client.connect(PORT, HOST, function () {

    console.log('CONNECTED TO: ' + HOST + ':' + PORT);

    var _recData = new Buffer("");

    // 建立连接后立即向服务器发送数据，服务器将收到这些数据
    var _data = {
        userList: [{
            userName: "111111",
            brokeId: "11"
        }, {
            brokeId: "222222",
            userName: "22"
        }],
        product: [1111, 2222, 3333, 4444],
        price: 333.01,
        maps: _maps
    };

    sendData(_data);

    client.on("end", function () {
        console.log('disconnected from server');
    });

    client.on('data', function (data) {
        //收到消息时，处理数据
        _recData = Buffer.concat([_recData, data]);
        getData();
    });

    client.on("error", function (err) {
        console.log("connect error", err);
    });

    function sendData(_data) {
        var _dataString = JSON.stringify(_data);
        var _buffer = new Buffer(_dataString);
        var _fBuffer = new Buffer(4);
        _fBuffer.writeInt32LE(_buffer.length, 0, 4);
        var _lastBuffer = Buffer.concat([_fBuffer, _buffer]);

        console.log("send length:", _lastBuffer.length, "buffer length", _buffer.length, "head length:", _fBuffer);

        //发送数据
        client.write(_lastBuffer, "", function () {
            // client.destroy();
            //是否关闭连接
        });

    }

    function getData() {
        if (_recData.length >= 4) {
            //缓存包长度大于4字节
            var _dataLength = _recData.readInt32LE(0, 4).toString() * 1;
            //取包头长度
            if (_recData.length - 4 >= _dataLength) {
                //判断包完整性
                var _data = _recData.toString("utf8", 4, 4 + _dataLength);
                _recData = _recData.slice(4 + _dataLength);
                handleData(_data);
                if (_recData.length > 4) {
                    //如果缓冲区数据可能还有包，递归处理
                    getData();
                }
            }
        }
    }

    function handleData(data) {
        var _data = JSON.parse(data);
        //这里_data就是解析后的json数据
        console.log("handleData:", _data.maps.length, "cache length:", _recData.length);
    }
});
```
以上为本人开发经验，如有错误或建议，欢迎交流

转自：https://www.w3cvip.org/topics/129
