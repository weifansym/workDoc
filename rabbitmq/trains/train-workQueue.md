## 工作队列(Work Queues)

官网地址: [Work Queues](https://www.rabbitmq.com/tutorials/tutorial-two-python.html)

工作队列（Work Queue）是为了使用多个work进程来处理分布式耗时任务。

工作队列（任务队列）的设计目的避免即时执行计算密集型任务而不得不长时间等待执行完成。取而代之，可以延时执行任务。我们将任务封装成一个消息并将其发送至队列中，运行在后台的work进程就会从队列中取出任务并最终执行它。可以使用多个work进程来分摊任务。业务模型如下：



在web项目中，不可能在一个Http连接（通常http连接时长都比较短）下去执行一个复杂的任务，而工作队列的特性用来解决这种问题最合适不过了。

### Java实现

为了模拟复杂的计算任务，约定根据消息中出现的”.“的个数来决定任务的执行时长，每一个点号将会执行一秒，比如消息”Hello…“将会执行3秒。

#### 2.1、生产者

由于要模拟许多耗时的任务，在此，生产者将不再由键盘输入消息，使用for循环来生成（让任务的发送几乎是在同一时间）。将新的生产者命名为Sender4Work.java。具体产生消息的逻辑如下：

```
//发送消息，随机生成消息
    for(int i=0; i<10; i++){
      //生成消息
      String message = ""+i;
      int r = (int)(Math.random()*10);
      for(int j=0; j<r; j++)
        message +=".";
      //退出
      if(message != null && "quit".equals(message)){
        channel.close();
          connection.close();
          break;
      }
      //发送消息
      else{
        channel.basicPublish("", QUEUE_NAME, null, message.getBytes());
        System.out.println("  >>>发送："+message+"");
      }
    }
```

#### 2.2、消费者

消费者相比上节中的HelloWorld示例将增加逻辑：根据消息的点数，来执行相应时间的任务。将其复制命名为Receiver4Work.java。具体如下：

处理消息逻辑

```
while (true) {
  QueueingConsumer.Delivery delivery = consumer.nextDelivery();
  String message = new String(delivery.getBody());
  System.out.println(" >>>接收消息：" + message);
  if(message!=null && "ok".equals(message)){
    channel.close();
    connection.close();
        break;
  }else{
    doWork(message);
  }
}
```

实际的work逻辑 

```
public static void doWork(String message) throws InterruptedException{
  for(char c : message.toCharArray()){
    if(c == '.'){
      System.out.println("["+message+"] 在"
        +System.currentTimeMillis()+"执行了。");
      Thread.sleep(1000);
    }
  }
}
```

### Round-robin dispatching(轮巡分发机制)

使用工作队列的主要优势是很容易实现并行任务处理。很容易增加work来实现。

默认情况下，RabbitMQ将会有序地将消息分发给消费者，平均每个消费者将会收到相同数目的消息，这种机制叫做**轮巡制**。从上面的输出可以看出，无论某个任务执行了多久，它都采用的是轮巡制（不会因某个消费者上的任务执行时间过长，导致消息会传递到其它消费者上处理）。

### Message acknowledgment(消息的应答机制)

执行任务可能会执行数秒。可能会担心一个消费者客户端在执行一个耗时较长的任务时，只执行了一部分便宕掉了。使用之前的代码，RabbitMQ在将消息传送给消费者后便会从内存中移除它，在这种情况下，一旦杀掉一个worker进程，将会丢掉它正在处理的任务，更可怕的是，我们还会丢失所有已经传送给这个worker且还未来得及处理的任务。

如果worker死掉，我们更希望任务会被传送至其它worker并进行处理，从而达到不丢失任务的目的。

为了确保任务不会被丢失，RabbitMQ支持消息应答机制。当worker接收了任务并且处理完成后，将会给RabbitMQ Server发送一个ack，告知RabbitMQ可以安全删除它了。

当一个消费者客户端挂掉后，RabbitMQ就会意识到任务并没有完全执行成功，并且会把它重新传递给其它消费者。这种办法可以保证消息不会被丢失，即使在消费者偶断的情况下。

RabbitMQ仅会在worker连接中断的情况下才会重新将消息发送给其它消费者客户端，它没有消息超时机制。如果任务执行非常非常长的时间也是没有问题的。

消息应答机制默认是开启的。之前的例子我们通过autoAck=true显式地关闭了，开启它是在消费者客户端开启的。如下：

```
QueueingConsumer consumer = new QueueingConsumer(channel);
    boolean autoAck = false;//启动消息应答机制
    channel.basicConsume(QUEUE_NAME, autoAck, consumer);
 
    while (true) {
          QueueingConsumer.Delivery delivery = consumer.nextDelivery();
          String message = new String(delivery.getBody());
          System.out.println(" >>>接收消息：" + message);
          if(message!=null && "ok".equals(message)){
            //在关闭前也需要，发送消息应答
            channel.basicAck(delivery.getEnvelope().getDeliveryTag(), false);
 
            channel.close();
            connection.close();
            break;
          }else{
            doWork(message);
            //发送消息应答
            channel.basicAck(delivery.getEnvelope().getDeliveryTag(), false);
          }
 
    }
```

> 如果在需要ack开启的时候，在处理完相关逻辑后忘记了ack，此时当客户端退出的时候，消息将会被再次投递。此时可能会导致消息在RabbitMQ中堆积，消息将会占用大量的内存，得不到释放。
>
> 为了调试这种类型的错误，你可以使用rabbitmqctl，来打印messages_unacknowledged字段：
>
> ```
> sudo rabbitmqctl list_queues name messages_ready messages_unacknowledged
> ```

 ### Message durability（）

前面已经了解了如何在消费者客户端挂掉后保证任务不会被丢失。但是如果RabbitMQ Server如果宕掉的话，任务还是会被丢失。

如果不进行设置，当RabbitMQ退出或者挂掉后，所有的队列和消息将会丢失。要保证消息不会真正被丢失，需要做两件事：将队列和消息标记为持久化的。

为了确保队列不丢失，需要将队列声明为持久化的。如下：

```
boolean durable = true;
channel.queueDeclare("hello", durable, false, false, null);
```

尽管语法没有问题，但是不会起作用，因为之前已经创建了一个叫hello的非持久化的队列。RabbitMQ不允许创建名字相同参数不同的队列，且会返回一个错误。修改个名字就好了。

queueDeclare()方法的改变必须应用到所有的生产者和消费者中。



 MessageProperties.PERSISTENT_TEXT_PLAIN可以保证消息也能够持久化。如下：

```
channel.basicPublish("", "task_queue",
            MessageProperties.PERSISTENT_TEXT_PLAIN,
            message.getBytes());
```

> 注意：将消息标记为持久化的并不意味着消息完全不会丢失。尽管告诉RabbitMQ将消息保存至硬盘，但RabbitMQ可能接收到消息还没有来得及保存至硬盘中，即RabbitMQ刚将消息保存至内存中还没有写进硬盘。尽管消息持久化机制并不能完全保证，但对于简单的消息队列，这已经足够了。如果需要更强的保证，可以使用**发布确认（publisher confirms）**机制。

### Fair dispatch（公平分发机制）

但消息分发有时候还是不能达到预期。比如在两个work的场景下，当所有奇数序号的消息均为重量级任务，而偶数序号的消息均为轻量级的，那么就会有一个队列非常繁忙，而另一个几乎什么都不做。RabbitMQ将不会知道这些，它还是会均匀地分发消息。

这种现象的产生是由于当一个消息进入队列后，RabbitMQ就会立即分发这个消息，而不会关心从消费者客户端返回了多少应答消息。它已经提前将第n个消息分发给了对应的第n个消费者。

为了改善这种问题，可以使用basicQos()方法来设置，传入参数为1。它告知RabbitMQ在同一时间最多只能给一个worker分发一条消息，换句话说就是在worker处理完并返回一个应答之前，不要再分发一个消息给它，结果RabbitMQ会将消息分发给那些空闲的worker。

```
int prefetchCount = 1;
channel.basicQos(prefetchCount);
```

> 注意：如果所有的worker都非常繁忙，队列可能会被填满。需要时常关注它，可能需要增加worker来解决它，或者使用其它策略。



参考：[工作队列(Work Queues)](https://blog.zenfery.cc/archives/85.html)
