## Publish/Subscribe(发布/订阅消息)
### 发布/订阅消息
之前创建的是一个工作队列。工作队列的设计思想是：每个任务仅能由一个worker消费。接下来做一些复杂点的东西：将一个消息传送至多个消费者客户端。这种模式称为“发布/订阅”。

创建一个简单的日志系统来演示这种模式，该系统包含两个简单的程序：一个是产生日志消息，一个接收消息并打印它们。在此日志系统中，所有启动的接收者都将接收这些消息。

实质上，发布的消息将会被广播至所有的消费者。

### 交换机（Exchanges）
之前只对一个队列发送和接收消息。下面介绍RabbbitMQ中所有的消息组件。先看下之前已经介绍过的：

* 生产者（Producer）：产生消息。
* 队列（Queue）：存储消息的缓存区。
* 消费者（Consumer）：接收消息。

RabbitMQ消息组件的核心设计架构是生产者从来都不会将消息直接发送至队列中，实际上，生产者甚至一点都不了解消息是否被传送至队列中。

生产者仅可以将消息发送至一个交换机（Exchange）。一个交换机是非常简单的东西。一方面是从生产者接收消息，另一方面是将消息发送至队列中。交换机必须知道如何处理它接收到的消息：发送至一个队列，发送至多个队列，或者丢弃。这些规则由交换机类型来定义。
![pub-sub.png](https://github.com/weifansym/workDoc/blob/master/images/rabbitmq/pub-sub.png)

可用的Exchange类型：direct、topic、headers、fanout。主要关注广播（fanout）类型。创建一个名为”logs“的广播交换机：
```
channel.exchangeDeclare("logs", "fanout");
```
广播交换机非常简单。从名字就能理解它，它会将所有接收到的消息发送至和它关联的队列。这个特性正是日志系统所需要的。

在RabbitMQ中默认会创建一些名字为amq.*和默认（未命名的）交换机（可以在Web UI界面上查看到），现在应该不太可能会用到它们。

之前，还没有接触到交换机，但还是能够将消息发送至队列。因为使用的是默认的交换机，它使用空字符串（””）来标识。再看一下先前用到的发布消息的代码：

```
channel.basicPublish("", "hello", null, message.getBytes());
```
第一个参数指的就是交换机的名称，空字符串意味着使用默认的交换机，通过指定的routingKey，消息就可以被路由至相应的队列。现在可以将消息发送至指定的交换机：
```
channel.basicPublish( "logs", "", null, message.getBytes());
```
### Temporary queues(临时队列)
之前，定义了如hello、task_queue的队列，因为需要将多个worker绑定至相同的队列，所以指定队列名称是非常必要的。但在此处的日志系统却不需要这样做，每个消费者都需要接收所有的日志消息，且每个消费者连接上时仅需要接收连接点之后的消息。为了达到此目的，需要做两件事情：

首先，无论何时连接至RabbitMQ都将会是一个全新的、空的队列。为此，需要队列名称为自动生成，甚至可以把自动生成队列名称的事情交给RabbitMQ Server。

其次，一旦消费者断开连接，其绑定的队列将会自动删除。

在使用Java客户端时，只需要不给方法queueDeclare()传递参数，即会创建一个名字为自动生成的、非持久化的、独一无二的、自动清除的队列：

```
String queueName = channel.queueDeclare().getQueue();
```
生成的随机名称可能为这种格式：amq.gen-JzTY20BRgKO-HjmUJj0wLg。

而且，一旦消费者断开连接，队列应该被删除，这里有一个exclusive标记来处理他。
```
result = channel.queue_declare(exclusive=True)
```
### Bindings(绑定)
![Bindings](https://www.rabbitmq.com/img/tutorials/bindings.png)
到此为止，已经创建了广播交换机和队列。下面需要告知交换机将消息发送给所有的队列。队列和交换机之间的这种关系叫做绑定。
```
channel.queueBind(queueName, "logs", "");
```
现在logs交换机则会将消息传送至所有的队列。可以使用命令rabbitmqctl list_bindings来查看所有的绑定关系。
### Putting it all together(汇总代码)
![together](https://www.rabbitmq.com/img/tutorials/python-three-overall.png)
产生消息的生产者程序，看起来和之前的程序并没有太大的差别，最大的差别是：之前发布消息是发至未命名的交换机，现在将消息发布至logs交换机。在此处指定routingKey对广播交换机来说是无用的。下面是生产者EmitLog.java程序：
```
package com.zenfery.example.rabbitmq;
import java.io.IOException;
import java.util.Scanner;
import com.rabbitmq.client.Channel;
import com.rabbitmq.client.Connection;
import com.rabbitmq.client.ConnectionFactory;
import com.rabbitmq.client.ConsumerCancelledException;
import com.rabbitmq.client.QueueingConsumer;
import com.rabbitmq.client.ShutdownSignalException;
//生产者
public class EmitLog {
  private static final String EXCHANGE_NAME = "logs"; // 交换机名称
   
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
    channel.exchangeDeclare(EXCHANGE_NAME, "fanout");
     
    //发送消息，接收控制台输入，并将其发送
    Scanner scanner = new Scanner(System.in);
    while(scanner.hasNextLine()){
      String message = scanner.nextLine();
      //退出
      if(message != null && "quit".equals(message)){
          channel.close();
          connection.close();
          break;
      }
      //发送消息
      else{
        channel.basicPublish(EXCHANGE_NAME, "", null, message.getBytes());
        System.out.println("  >>>发送："+message+"");
      }
    }
  }
}
```
在建立连接之后，紧接着声明了交换机，这一步是必需的，因为禁止向未知的交换机发送消息。在没有队列绑定至此交换机时，消息将会全部丢失，在日志系统中，这样是没有问题的。如果没有消费者进行监听，丢掉消息是安全的。ReceiveLogs.java的具体代码如下：
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
public class ReceiveLogs {
  private static final String EXCHANGE_NAME = "logs"; // 交换机名称
   
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
    channel.exchangeDeclare(EXCHANGE_NAME, "fanout");
     
    //自动生成队列，并绑定队列至交换机
    String queueName = channel.queueDeclare().getQueue();
    channel.queueBind(queueName, EXCHANGE_NAME, "");
     
    //创建消费者对象
    QueueingConsumer consumer = new QueueingConsumer(channel);
    channel.basicConsume(queueName, true, consumer);
     
    while (true) {
          QueueingConsumer.Delivery delivery = consumer.nextDelivery();
          String message = new String(delivery.getBody());
          System.out.println(" >>>接收消息：" + message);
          if(message!=null && "ok".equals(message)){
            channel.close();
            connection.close();
            break;
          }
    }
  }
}
```
下面演示执行：

* 启动EmitLog，再启动每一个ReceiveLogs。并在EmitLog控制台输入”first message.“。
* 启动第二个ReceiveLogs。并在EmitLog控制台输入”second message.“。

结果：

第一个ReceiveLogs的输出：
```
>>>接收消息：first message.
>>>接收消息：second message.
```
第二个ReceiveLogs的输出：
```
>>>接收消息：second message.
```

参考：[发布/订阅消息](https://blog.zenfery.cc/archives/90.html)



