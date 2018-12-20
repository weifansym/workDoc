## Routing(路由消息)
### Routing(路由)
之前创建过了一个简单的日志系统。可以将日志信息广播至许多接收者。在本节中，会将日志系统增加一个特性：仅订阅日志消息的一个子集。
例如，仅仅将关键的错误日志消息写入日志文件（保存在磁盘），同时还能够将所有的日志消息打印至控制台。

### Bindings(绑定)
在之前的例子中，已经使用过绑定了，调用代码如下：
```
channel.queueBind(queueName, EXCHANGE_NAME, "");
```
绑定即交换机和队列之间的一个关系。可以简单理解为：该队列关注此交换机中的消息。绑定方法可以传入一个routingKey参数，为了避免和basic_publish参数混淆，
将它叫做绑定关键字（binding key）。下面是如何使用关键字来进行绑定：
```
channel.queueBind(queueName, EXCHANGE_NAME, "black");
```
绑定关键字的意义视交换机类型而定。先前使用的广播交换机（fanout）将会忽略该关键字。
### Direct exchange
之前的日志系统将所有的消息广播至所有的消费者客户端。我们希望在此基础上根据日志的级别来过滤消息。例如，我们可能仅需要将关键的错误日志写入磁盘中，
从而不会在warning或info日志消息上浪费磁盘空间。

广播交换机不能提供复杂的特性，它仅能实现简单的广播机制。下面使用direct交换机来代替广播（fanout）交换机，direct交换机的路由算法相对简单：
只有队列的绑定关键字和消息的路由关键字完全匹配时，消息才能够发送至队列。如下图的路由机制：
![Direct](https://www.rabbitmq.com/img/tutorials/direct-exchange.png)
可以看出在此路由机制下，有两个队列绑定在同一个direct类型的交换机上。第一个使用orange关键字进行绑定，第二个队列有两个绑定关键字，black和green。
在此机制下，使用路由关键字orange发布的消息将会被发布至队列Q1，使用路由关键字black或green将会至Q2，其它所有的消息将会被丢弃。
### Multiple bindings
![Multiple bindings](https://www.rabbitmq.com/img/tutorials/direct-exchange-multiple.png)
多个队列使用相同的绑定关键字是非常合法的，在上图所示的例子中，在例子中，可以为Q1增加一个绑定关键字black绑定至交换机X。在这种情况下，
direct交换机就像fanout交换机一样，将会广播消息至所有匹配的队列中。使用路由关键字black的消息将会传送至Q1和Q2。

### Emitting logs
在日志系统中使用这种模型，使用direct交换机代替fanout交换机来发送消息。使用日志级别做为路由关键字，接收程序可以选择它想要接收的级别的日志。

和以前一样，首先需要先创建交换机：
```
channel.exchangeDeclare(EXCHANGE_NAME, "direct");
```
发送消息：
```
channel.basicPublish(EXCHANGE_NAME, severity, null, message.getBytes());
```
在上面的代码中，可以假定变量severity可取值为：info、warning、error。
### Subscribing
订阅消息很简单，仅需要将队列关注的级别绑定至相应的交换机即可：
```
String queueName = channel.queueDeclare().getQueue();
channel.queueBind(queueName, EXCHANGE_NAME, severity);
```
### Putting it all together
![together](https://www.rabbitmq.com/img/tutorials/python-four.png)
#### 生产者
接收控制台的输入，每行输入以空格分隔，第一位表示日志级别。详细代码如下：
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
public class EmitLogDirect {
 
  private static final String EXCHANGE_NAME = "direct_logs"; // 交换机名称
   
  public static void main(String[] args) throws IOException
  , ShutdownSignalException, ConsumerCancelledException
  , InterruptedException {
    ConnectionFactory factory = new ConnectionFactory();
    factory.setHost("localhost");
    factory.setPort(5672);
    //创建连接
    Connection connection = factory.newConnection();
    Channel channel = connection.createChannel();
    //定义类型为direct的交换机
    channel.exchangeDeclare(EXCHANGE_NAME, "direct");
     
    //发送消息，接收控制台输入，按行输入，日志级别和日志内容中间使用空格分隔
    Scanner scanner = new Scanner(System.in);
    while(scanner.hasNextLine()){
      String line = scanner.nextLine();
      String severity = line.split(" ")[0];
      String message = line;
      //发送消息
      channel.basicPublish(EXCHANGE_NAME, severity, null, message.getBytes());
      System.out.println("  >>>发送：["+severity+"] "+message+"");
    }
  }
}
```
#### 消费者
启动的消费者根据传入的日志级别参数列表，来决定监听哪些日志级别的日志，其余的将会被忽略：
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
public class ReceiveLogsDirect {
 
  private static final String EXCHANGE_NAME = "direct_logs"; // 交换机名称
   
  public static void main(String[] args) throws IOException
  , ShutdownSignalException, ConsumerCancelledException
  , InterruptedException {
    ConnectionFactory factory = new ConnectionFactory();
    factory.setHost("localhost");
    factory.setPort(5672);
    //创建连接
    Connection connection = factory.newConnection();
    Channel channel = connection.createChannel();
    //定义类型为fanout的交换机
    channel.exchangeDeclare(EXCHANGE_NAME, "direct");
     
    //自动生成队列，根据main方法传入的关注的日志级别并绑定队列至交换机
    String queueName = channel.queueDeclare().getQueue();
    for(String severity: args){
      channel.queueBind(queueName, EXCHANGE_NAME, severity);
    }
     
    //创建消费者对象
    QueueingConsumer consumer = new QueueingConsumer(channel);
    channel.basicConsume(queueName, true, consumer);
     
    while (true) {
          QueueingConsumer.Delivery delivery = consumer.nextDelivery();
          String message = new String(delivery.getBody());
          System.out.println(" >>>接收消息：" + message);
    }
  }
 
}
```
#### 执行演示
传入参数“error”，启动第一个消费者。传入参数“error info warning”启动第二个消费者。

启动第一个消费者，依次输入:
```
info logs123.
error logs456.
warning logs789.
other logs890.
```
第一个消费者输出：
```
>>>接收消息：error logs456.
```
第二个消费者输出：
```
>>>接收消息：info logs123.
>>>接收消息：error logs456.
>>>接收消息：warning logs789.
```
