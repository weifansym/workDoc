## Netty精粹之TCP粘包拆包问题
粘包拆包问题是处于网络比较底层的问题，在数据链路层、网络层以及传输层都有可能发生。我们日常的网络应用开发大都在传输层进行，由于UDP有消息保护边界，
不会发生这个问题，因此这篇文章只讨论发生在传输层的TCP粘包拆包问题。

### 什么是粘包、拆包？

对于什么是粘包、拆包问题，我想先举两个简单的应用场景：

客户端和服务器建立一个连接，客户端发送一条消息，客户端关闭与服务端的连接。

客户端和服务器简历一个连接，客户端连续发送两条消息，客户端关闭与服务端的连接。

对于第一种情况，服务端的处理流程可以是这样的：当客户端与服务端的连接建立成功之后，服务端不断读取客户端发送过来的数据，当客户端与服务端连接断开之后，
服务端知道已经读完了一条消息，然后进行解码和后续处理...。对于第二种情况，如果按照上面相同的处理逻辑来处理，那就有问题了，我们来看看第二种情况下客户端
发送的两条消息递交到服务端有可能出现的情况：

#### 第一种情况：
服务端一共读到两个数据包，第一个包包含客户端发出的第一条消息的完整信息，第二个包包含客户端发出的第二条消息，那这种情况比较好处理，服务器只需要简单的从
网络缓冲区去读就好了，第一次读到第一条消息的完整信息，消费完再从网络缓冲区将第二条完整消息读出来消费。

![https://github.com/weifansym/workDoc/blob/master/images/20181128001.png](https://github.com/weifansym/workDoc/blob/master/images/20181128001.png)
                                                没有发生粘包、拆包示意图
#### 第二种情况：
服务端一共就读到一个数据包，这个数据包包含客户端发出的两条消息的完整信息，这个时候基于之前逻辑实现的服务端就蒙了，因为服务端不知道第一条消息从哪儿结束
和第二条消息从哪儿开始，这种情况其实是发生了TCP粘包。

![https://github.com/weifansym/workDoc/blob/master/images/20181128002.png](https://github.com/weifansym/workDoc/blob/master/images/20181128002.png)
TCP粘包示意图
#### 第三种情况：
服务端一共收到了两个数据包，第一个数据包只包含了第一条消息的一部分，第一条消息的后半部分和第二条消息都在第二个数据包中，或者是第一个数据包包含了第一条
消息的完整信息和第二条消息的一部分信息，第二个数据包包含了第二条消息的剩下部分，这种情况其实是发送了TCP拆，因为发生了一条消息被拆分在两个包里面发送了，
同样上面的服务器逻辑对于这种情况是不好处理的。

![https://github.com/weifansym/workDoc/blob/master/images/20181128003.png](https://github.com/weifansym/workDoc/blob/master/images/20181128003.png)
TCP拆包示意图

### 为什么会发生TCP粘包、拆包呢？

发生TCP粘包、拆包主要是由于下面一些原因：

* 应用程序写入的数据大于套接字缓冲区大小，这将会发生拆包。
* 应用程序写入数据小于套接字缓冲区大小，网卡将应用多次写入的数据发送到网络上，这将会发生粘包。
* 进行MSS（最大报文长度）大小的TCP分段，当TCP报文长度-TCP头部长度>MSS的时候将发生拆包。
* 接收方法不及时读取套接字缓冲区数据，这将发生粘包。
* ……

### 如何处理粘包、拆包问题？
知道了粘包、拆包问题及根源，那么如何处理粘包、拆包问题呢？TCP本身是面向流的，作为网络服务器，如何从这源源不断涌来的数据流中拆分出或者合并出
有意义的信息呢？通常会有以下一些常用的方法：

* 使用带消息头的协议、消息头存储消息开始标识及消息长度信息，服务端获取消息头的时候解析出消息长度，然后向后读取该长度的内容。
* 设置定长消息，服务端每次读取既定长度的内容作为一条完整消息。
* 设置消息边界，服务端从网络流中按消息编辑分离出消息内容。
* ……
### 如何基于Netty处理粘包、拆包问题？

我的上一篇文章的ChannelPipeline部分大概讲了Netty网络层数据的流向以及ChannelHandler组件对网络数据的处理，这一小节也会涉及到相关重要组件：

1. ByteToMessageDecoder
2. MessageToMessageDecoder

这两个组件都实现了ChannelInboundHandler接口，这说明这两个组件都是用来解码网络上过来的数据的。而他们的顺序一般是ByteToMessageDecoder
位于head channel handler的后面，MessageToMessageDecoder位于ByteToMessageDecoder的后面。Netty中，涉及到粘包、拆包的逻辑主要
在ByteToMessageDecoder及其实现中。

### ByteToMessageDecoder
顾名思义、ByteToMessageDecoder是用来将从网络缓冲区读取的字节转换成有意义的消息对象的，对于源码层面指的说明的一段是下面这部分：
```
protected void callDecode(ChannelHandlerContext ctx, ByteBuf in, List<Object> out) {
    try {
        while (in.isReadable()) {
            int outSize = out.size();

            if (outSize > 0) {
                fireChannelRead(ctx, out, outSize);
                out.clear();
                
                if (ctx.isRemoved()) {
                    break;
                }
                outSize = 0;
            }

            int oldInputLength = in.readableBytes();
            decode(ctx, in, out);

            if (ctx.isRemoved()) {
                break;
            }

            if (outSize == out.size()) {
                if (oldInputLength == in.readableBytes()) {
                    break;
                } else {
                    continue;
                }
            }

            if (oldInputLength == in.readableBytes()) {
                throw new DecoderException(
                        StringUtil.simpleClassName(getClass()) +
                        ".decode() did not read anything but decoded a message.");
            }

            if (isSingleDecode()) {
                break;
            }
        }
    } catch (DecoderException e) {
        throw e;
    } catch (Throwable cause) {
        throw new DecoderException(cause);
    }
}
```
为了节省篇幅，我把注释删除掉了，当上面一个channel handler传入的ByteBuf有数据的时候，这里我们可以把in参数看成网络流，这里有不断的数据流入，
而我们要做的就是从这个byte流中分离出message，然后把message添加给out。分开将一下代码逻辑：

1. 当out中有Message的时候，直接将out中的内容交给后面的channel handler去处理。

2. 当用户逻辑把当前channel handler移除的时候，立即停止对网络数据的处理。

3. 记录当前in中可读字节数。

4. decode是抽象方法，交给子类具体实现。

5. 同样判断当前channel handler移除的时候，立即停止对网络数据的处理。

6. 如果子类实现没有分理出任何message的时候，且子类实现也没有动bytebuf中的数据的时候，这里直接跳出，等待后续有数据来了再进行处理。

7. 如果子类实现没有分理出任何message的时候，且子类实现动了bytebuf中的数据，则继续循环，直到解析出message或者不在对bytebuf中数据进行处理为止。

8. 如果子类实现解析出了message但是又没有动bytebuf中的数据，那么是有问题的，抛出异常。

9. 如果标志位只解码一次，则退出。

可以知道，如果要实现具有处理粘包、拆包功能的子类，及decode实现，必须要遵守上面的规则，我们以实现处理第一部分的第二种粘包情况和第三种情况拆包情况的
服务器逻辑来举例：

对于粘包情况的decode需要实现的逻辑对应于将客户端发送的两条消息都解析出来分为两个message加入out，这样的话callDecode只需要调用一次decode即可。

对于拆包情况的decode需要实现的逻辑主要对应于处理第一个数据包的时候第一次调用decode的时候out的size不变，从continue跳出并且由于不满足继续可读
而退出循环，处理第二个数据包的时候，对于decode的调用将会产生两个message放入out，其中两次进入callDecode上下文中的数据流将会合并为一个bytebuf
和当前channel handler实例关联，两次处理完毕即清空这个bytebuf。

当然，尽管介绍了ByteToMessageDecoder，用户自己去实现处理粘包、拆包的逻辑还是有一定难度的，Netty已经提供了一些基于不同处理粘包、拆包规则的实现：
如DelimiterBasedFrameDecoder、FixedLengthFrameDecoder、LengthFieldBasedFrameDecoder和LineBasedFrameDecoder等等。其中：

* DelimiterBasedFrameDecoder: 是基于消息边界方式进行粘包拆包处理的。
* FixedLengthFrameDecoder: 是基于固定长度消息进行粘包拆包处理的。
* LengthFieldBasedFrameDecoder: 是基于消息头指定消息长度进行粘包拆包处理的。
* LineBasedFrameDecoder: 是基于行来进行消息粘包拆包处理的。

用户可以自行选择规则然后使用Netty提供的对应的Decoder来进行具有粘包、拆包处理功能的网络应用开发。

最后

在通常的高性能网络应用中，客户端通常以长连接的方式和服务端相连，因为每次建立网络连接是一个很耗时的操作。比如在RPC调用中，如果一个客户端远程调用的过程中，
连续发起了多次调用，而如果这些调用对应于同一个连接的时候，那么就会出现服务器需要对于这些多次调用消息的粘包拆包问题的处理。如果是你，你会选择哪种策略呢？

转自：https://my.oschina.net/andylucc/blog/625315
