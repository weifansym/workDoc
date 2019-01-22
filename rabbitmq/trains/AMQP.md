## 消息队列协议AMQP的设计原理
### AMQP是什么
AMQP（Advanced Message Queuing Protocol 高级消息队列协议）是一个消息协议，它支持符合标准的客户端请求程序与符合标准的消息中间件代理进行通信。
### 协议结构
![协议结构](https://github.com/weifansym/workDoc/blob/master/images/rabbitmq/rabbit1-1.png)
* Broker: 接收和分发消息的应用，RabbitMQ Server就是Message Broker。
* Virtual host: 出于多租户和安全因素设计的，把AMQP的基本组件划分到一个虚拟的分组中，类似于网络中的namespace概念。当多个不同的用户使用同一个中间件提供的服务时，可以划分出多个vhost，每个用户在自己的vhost创建exchange／queue等。
* Channel: 如果每一次访问都建立一个Connection，在消息量大的时候建立Connection的开销将是巨大的，效率也较低。Channel作为轻量级的Connection极大减少了操作系统建立connection的开销。线程池的思想。
* Exchange: message到达broker的第一站，根据分发规则，匹配查询表中的routing key，分发消息到queue中去。常用的类型有：direct (point-to-point), topic (publish-subscribe) and fanout (multicast)。
* Queue: 消息最终被送到这里等待consumer取走。一个message可以被同时拷贝到多个queue中。
* Binding: exchange和queue之间的虚拟连接，binding中可以包含routing key。Binding信息被保存到exchange中的查询表中，用于message的分发依据。其实可以理解为Exchange与Queue的关系对照表。

#### 工作流程：
1. 消息被发布到exchanges，通常可将exchanges比作邮局或者邮箱
2. exchanges按照bindings中的规则，将消息副本分发到queues
3. AMQP brokers传递消息给与queues关联的consumers，或者consumers按照需求从queues拉取信息
### 其他内容
#### exchange
Default exchange：

default exchange是一个没有名称的、被broker预先申明的direct exchange。它所拥有的一个特殊属性使它对于简单的应用程序很有作用：每个创建的queue会与它自动绑定，使用queue名称作为routing key。
举例说，当你申明一个名称为“search-indexing-online”的queue时，AMQP broker使用“search-indexing-online”作为routing key将它绑定到default exchange。因此，一条被发布到default exchange并且routing key为”search-indexing-online”将被路由到名称为”search-indexing-online”的queue。

Direct exchange：

direct exchange根据消息的routing key来传送消息。direct exchange是单一传播路由消息的最佳选择（尽管他们也可以用于多路传播路由），以下是它们的工作原理：

* 一个routing key为K的queue与exchange进行绑定
* 当一条新的routing key为R的消息到达direct exchange时，如果K=R,exchange 将它路由至该queue
![direct exchange](https://github.com/weifansym/workDoc/blob/master/images/rabbitmq/rabbit1-2.png)

Fanout exchange:

fanout exchange路由消息到所有的与其绑定的queue中，忽略routing key。如果N个queue被绑定到一个fanout exchange，当一条新消息被发布到exchange时，消息会被复制并且传送到这N个queue。fanout exchange是广播路由的最佳选择。
![fanout exchange](https://github.com/weifansym/workDoc/blob/master/images/rabbitmq/rabbit1-3.png)

Topic exchange：
根据routing key，通过表达式将消息分配到匹配的队列中，Topic exchange将分发到目标queue中。如， 包含分类与标签的新闻信息推送。

Headers exchange：
headers exchanges忽略routing key属性，相反用于路由的属性是从headers属性中获取的。如果消息头的值等于指定的绑定值，则认为消息是匹配的。

### 消息发送与接受
#### 发送者确认：
在消息生产者与broker传递消息的过程中，由于出现种种原因，消息无法传递到broker。此时，确认消息是否正确传递就变得至关重要，AMQP支持生产者对消息的确认，可以通过两种方式：

1. 事物方式
2. confirm模式

![](https://github.com/weifansym/workDoc/blob/master/images/rabbitmq/rabbit1-4.png)

#### 消费者应答: 
同样在broker与消费者传递消息的过程中，也有可能由于某种原因导致消息无法传递到消费者。这就会出现一个问题，什么时候讲队列中的消息删除？可以通过两种方式：

1. 投递的标识：Delivery Tags
2. 应答模式

![](https://github.com/weifansym/workDoc/blob/master/images/rabbitmq/rabbit1-5.png)
### 总图如下：
![](https://github.com/weifansym/workDoc/blob/master/images/rabbitmq/rabbit1-6.png)
