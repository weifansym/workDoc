## 关于Node.js 包 amqplib
* api地址：[amqplib](http://www.squaremobius.net/amqp.node/channel_api.html#confirmchannel)

### Channel与ConfirmChannel

### Channel#assertQueue
在声明一个队列的时候除了队列名之外还有其他可选参数：
* exclusive(独占)：如果设置为true,则这个队列会独占这个链接（默认为false）
* durable（持久化）:  如果为true时，当rabbitmq的服务重启后队列依然存在。在exclusive 和 autoDelete模式下也适用，默认是true
* autoDelete: 如果设置为true时，在消费者不存在的时候队列将会删除，默认是false
* arguments: 附加参数，通常是某种特定于mq服务的扩展的参数，例如高可用性，TTL，rabbitmq的扩展也可以作为option参数。具体看文档
* messageTtl: 队列中的消息多久过期
* expires: 队列不是使用的情况下，将会在n毫秒后被销毁，这里的使用是有消费者
* deadLetterExchange: 一个交换器，队列中丢弃的消息，将被重新发送。
* maxLength: 设置队列存储消息的最大数，老的消息将会被丢弃
* maxPriority: 把队列设置成一个优先级队列

### Channel#assertExchange
声明exchange的时候必须指定一个字符串，切不为空，可选参数如下：
* durable: 是否持久化，在服务重启后依然存在，默认为true
* internal: 如果为true消息不直接发送到exchange上，默认为false
* autoDelete: 如果为true，当exchange没有都应的绑定，exchange就会销毁，默认为false
* alternateExchange: 备用exchange，当消息发到exchange上，但是没有路由到任何队列的时候
* arguments: exchange类型对应的其他可选参数

### Channel#bindExchange
```
#bindExchange(destination, source, pattern, [args])
```
绑定一个exchange到其他的exchange，destination这个exchange会接收到source这个exchange的消息，根据source的类型以及给定的pattern。
例如，一个direct类型的exchange将会传递消息到一个到一个路由key和pattern相同的exchange。
### Channel#publish
```
#publish(exchange, routingKey, content, [options])
```
传递消息到exchange，必填参数如下：
* exchange和routingKey：exchange和路由key,决定了消息的流向。一个特殊的例子是exchange为空字符串，这样消息将会直接发送到路由key对应的队列里。**equivalent**方法和这个特殊情况相同。如果指定的exchange不存在，则对应的channel将会关闭。
* content：buffer类型的消息内容

余下的参数在options中指定，这些参数分为对RabbitMQ有一定意义的那些以及RabbitMQ会忽略但传递给消费者的那些。这些参数都可以忽略，将会提供默认值。
option中的参数可以是BasicDeliver（用来传递消息的方法），BasicProperties（消息header frame）,以及RabbitMQ的扩展。

其中在RabbitMQ中起作用且发送给消费者的：
* expiration：如果存在，则队列中的消息在制定时间之后就会销毁。
* userId：如果提供，Rabbitmq将会把他和在打开连接时使用的username相比较，只有匹配的消息才会传递。
* CC：路邮key数组，消息将会路由到除了routingKey之外的其他路由中。
* priority：消息的优先级，如果队列不是优先级队列则不使用。
* persistent：如果为true,消息以及对应的队列持久化，即：Rabbimq服务重启，消息和他对应的队列都存在。
* deliveryMode：

在RabbitMQ中器作用且不会发送给消费者的：
* mandatory：如果是true,则如果消息没有路由到队列则返回。
* BCC

不作用在RabbitMQ上，且不发送给消费者的：
* immediate

RabbitMQ会忽略的但程序中可能会用到的
* contentType
* contentEncoding
* headers
* correlationId
* replyTo
* messageId
* timestamp
### Channel#consume
```
#consume(queue, function(msg) {...}, [options])
```
可选参数：Options如下：
* consumerTag
* noLocal: 理论上如果为真，如果在同一个链接上有publish消息，则rabbitmq服务将不会传递消息给他。默认为false
* noAck: 如果为true,服务器将不会等待对发送给消费者的消息确认。消息一旦被发送就会将消息出列。 默认为false
* exclusive：如果设置为true,服务不会让任何人消费这个队列。
* priority：给消费者设置优先级，优先级较高的消费者优先于低优先级消费者获取消息。
* arguments：

服务端的返回值包含一个consumerTag字段，最好保留下这个值，因为如果你想取消这个消费的操作（停止获取消息），
第二个回调的参数如下形式：
```
{
  content: Buffer,
  fields: Object,
  properties: Object
}
```
### Channel#ack
确认给定的消息。

如果#consume or #get带有参数noAck（默认值为false）。服务器会在删除消息之前会等待消息确认。如果没有进行确认，一旦channel关闭，消息将会重新入队。

如果allUpTo设置为true：在给定消息之前和之后的所有未完成消息都应被视为已确认。 如果为false或省略，则仅确认提供的消息

提供不需要确认或已经确认的消息是错误的。 这样做会使channel出错。如果您要确认所有消息，并且没有特定的消息，使用#ackAll

### Channel#ackAll
对channel上所有已出队列的消息进行确认。这是一个安全的操作，即使在没这个消息的情况下也不会导致错误发生。
