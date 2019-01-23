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


