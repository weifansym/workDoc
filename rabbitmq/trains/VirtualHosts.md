## Virtual Host与权限管理
RabbitMQ是一个多租户系统，链接，交换机，队列，绑定，用户权限,协议，以及一些其他的东西都是属于Virtual Host的，他是一个完整的逻辑组合。Virtual hosts
提供的是一个逻辑集合来对资源进行分离。virtual hosts是物理分离。因此在说明用户权限以及连接的时候说明对应的vhost是很有必要的。
每个VirtualHost相当月一个相对独立的RabbitMQ服务器，每个VirtualHost之间是相互隔离的。exchange、queue、message不能互通。
其实也不要想得太复杂，我们知道mysql有数据库的概念。其实Virtual Host类似于mysql中的某个库的概念。我们可以对这个库以及表等设定操作权限。
 
在RabbitMQ中无法通过AMQP创建VirtualHost，可以通过以下命令来创建。
```
rabbitmqctl add_vhost [vhostname]
```
当然也可以通过WEB管理插件来创建。 

如上图在创建完vhost后可以在All Virtual Host标签看到新建的VirtualHost。

## 用户权限管理
通常在权限管理中主要包含三步：
* 新建用户
* 配置权限
* 配置角色
### 新建用户

