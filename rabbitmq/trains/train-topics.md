## Topics(主题/订阅消息)
### 主题模式（Topics）
广播（fanout）交换机仅能够广播消息，使用direct交换机，可以对消息进行筛选过滤。尽管使用direct交换机改进了日志系统，但它还是有所限制，
它不能使用更复杂的路由规则。

在日志系统中，可能不仅仅订阅基于日志级别的筛选消息，日志可能来自不同的源，也要加以区分。像unix的syslog日志，它即可以区分日志级别（info/warn/crit…），
也可以使用其它更灵活的区分机制（auth/cron/kern…）。如果想监听特殊的错误日志、cron产生的日志以及kern产生的所有日志，这将会复杂一些。为了实现此功能，
需要了解主题交换机（topic）。

### 主题交换机（Topic exchange）
被发送至主题交换机的消息不可以是任意的路由关键字，它必须是以点号分隔的多个单词串。这些单词可以任意，通常是和消息的特性有关，如下一些有效的路由关键字：
”stock.usd.nyse”, “nyse.vmw”, “quick.orange.rabbit”。路由关键字的上限是255 bytes。

和路由关键字相对应，绑定关键字应该有相似的格式。主题交换机和direct交换机的机制类似，都是根据绑定关键字和路由关键字的匹配来决定消息被传送至哪些消息队列。
关于绑定关键字有两点比较重要：

* 匹配一个单词。
# 匹配一个或多个单词。
使用下图的例子来说明：
![Topics](https://www.rabbitmq.com/img/tutorials/python-five.png)
在这个例子中，将发送一些用来描述动物的消息。这些消息的路由关键字均包含三个单词（即含有两个点号），第一个单词描述速度，第二个描述颜色，第三个是种类：
”<speed>.<colour>.<species>”。此处有三个绑定规则：Q1使用绑定关键字”*.orange.*”，Q2使用”*.*.rabbit”和”lazy.#”，这些绑定可解释如下：

* Q1关注的是所有颜色为橙色（orange）的动物。
* Q2希望收到所有兔子（rabbit）的消息以及比较懒惰（lazy）的动物。

路由关键字为”quick.orange.rabbit”或”lazy.orange.elephant”的消息将会传送至所有的队列。而”quick.orange.fox”将会传送至Q1，
”lazy.brown.fox”传送至Q2。”lazy.pink.rabbit”虽然可以匹配Q2的两个绑定关键字，但它仅会被传送至Q2一次。”quick.brown.fox”
将不会匹配任何绑定关键字，将会被丢弃。

如果使用1个或4个单词的路由关键字（如：”orange”或”quick.orange.male.rabbit”）将会发生什么呢？当然，它们匹配不上任何规则，将会被丢弃。
”lazy.orange.male.rabbit”即使有4个单词，也会匹配lazy.#规则，而被传送至Q2。

> 注意：主题交换机功能强大，它可以模拟其它的交换机功能。
如果队列使用”#”绑定关键字，它将匹配所有的路由关键字而接收所有的消息，正如fanout交换机的功能所示。
如果队列的绑定关键字未使用字符”#”和”*”，那么主题关键字将会表现出direct交换机的特性。

### Putting it all together
#### 生产者
```
package com.zenfery.example.rabbitmq;
 
import java.io.IOException;
import java.util.Scanner;
 
import com.rabbitmq.client.Channel;
import com.rabbitmq.client.Connection;
import com.rabbitmq.client.ConnectionFactory;
import com.rabbitmq.client.ConsumerCancelledException;
import com.rabbitmq.client.ShutdownSignalException;
 
//生产者
public class EmitLogTopic {
 
  private static final String EXCHANGE_NAME = "topic_logs"; // 交换机名称
   
  public static void main(String[] args) throws IOException
  , ShutdownSignalException, ConsumerCancelledException
  , InterruptedException {
    ConnectionFactory factory = new ConnectionFactory();
    factory.setHost("localhost");
    factory.setPort(5672);
    //创建连接
    Connection connection = factory.newConnection();
    Channel channel = connection.createChannel();
    //定义类型为topic的交换机
    channel.exchangeDeclare(EXCHANGE_NAME, "topic");
     
    //发送消息，接收控制台输入，按行输入，路由关键字和日志内容中间使用空格分隔
    Scanner scanner = new Scanner(System.in);
    while(scanner.hasNextLine()){
      String line = scanner.nextLine();
      String key = line.split(" ")[0];
      String message = line.split(" ")[1];
      //发送消息
      channel.basicPublish(EXCHANGE_NAME, key, null, message.getBytes());
      System.out.println("  >>>发送：["+key+"] "+message+"");
    }
  }
}
```
### 消费者
```
package com.zenfery.example.rabbitmq;
 
import java.io.IOException;
 
import com.rabbitmq.client.Channel;
import com.rabbitmq.client.Connection;
import com.rabbitmq.client.ConnectionFactory;
import com.rabbitmq.client.ConsumerCancelledException;
import com.rabbitmq.client.QueueingConsumer;
import com.rabbitmq.client.ShutdownSignalException;
 
//消费者
public class ReceiveLogsTopic {
 
  private static final String EXCHANGE_NAME = "topic_logs"; // 交换机名称
   
  public static void main(String[] args) throws IOException
  , ShutdownSignalException, ConsumerCancelledException
  , InterruptedException {
    ConnectionFactory factory = new ConnectionFactory();
    factory.setHost("localhost");
    factory.setPort(5672);
    //创建连接
    Connection connection = factory.newConnection();
    Channel channel = connection.createChannel();
    //定义类型为topic的交换机
    channel.exchangeDeclare(EXCHANGE_NAME, "topic");
     
    //自动生成队列，根据main方法传入的关注的主题规则并绑定队列至交换机
    String queueName = channel.queueDeclare().getQueue();
    for(String key: args){
      channel.queueBind(queueName, EXCHANGE_NAME, key);
    }
     
    //创建消费者对象
    QueueingConsumer consumer = new QueueingConsumer(channel);
    channel.basicConsume(queueName, true, consumer);
     
    while (true) {
          QueueingConsumer.Delivery delivery = consumer.nextDelivery();
          String message = new String(delivery.getBody());
          String key = delivery.getEnvelope().getRoutingKey();
          System.out.println(" >>>接收消息：["+key+"]" + message);
    }
  }
}
```
#### 演示执行
传入绑定关键字”#”启动消费者客户端，传入绑定关键字”kern.*”启动第二个消费者客户端。

启动生产者客户端依次发送消息：
```
other otherlogs
kern kernlogs
kern.log kernfulllogs
```
第一个消费者输出:
```
>>>接收消息：[other]otherlogs
>>>接收消息：[kern]kernlogs
>>>接收消息：[kern.log]kernfulllogs
```
第二个消费者输出:
```
>>>接收消息：[kern.log]kernfulllogs
```
> 注意:
“#.*” 将会匹配”..”或”.”路由关键字，它还可以匹配单个单词的路由关键字。
“*” 将不匹配单个路由关键字，它也不能匹配空字符串的路由关键字。

参考：[主题/订阅消息](https://blog.zenfery.cc/archives/111.html)
