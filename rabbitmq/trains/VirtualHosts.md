## Virtual Host与权限管理

官方文档：[Virtual Host](https://www.rabbitmq.com/vhosts.html)

我们知道在RabbitMQ服务刚起动的时候会默认生成一个名字是"/"的virtual host,分配一个用户名“guest”以及初始密码“guest”。下面来看一下	virtual host。

RabbitMQ是一个多租户系统，链接，交换机，队列，绑定，用户权限,协议，以及一些其他的东西都是属于Virtual Host的，他是一个完整的逻辑组合。Virtual hosts
提供的是一个逻辑集合来对资源进行分离。virtual hosts是物理分离。因此在说明用户权限以及连接的时候说明对应的vhost是很有必要的。
每个VirtualHost相当月一个相对独立的RabbitMQ服务器，每个VirtualHost之间是相互隔离的。exchange、queue、message不能互通。
其实也不要想得太复杂，我们知道mysql有数据库的概念。其实Virtual Host类似于mysql中的某个库的概念。我们可以对这个库以及表等设定操作权限。

在RabbitMQ中无法通过AMQP创建VirtualHost，可以通过以下命令来创建。
```
rabbitmqctl add_vhost [vhostname]
```
当然也可以通过WEB管理插件来创建。 
![vhost](https://github.com/weifansym/workDoc/blob/master/images/rabbitmq/vhost.png)
如上图在创建完vhost后可以在All Virtual Host标签看到新建的VirtualHost。

## 用户权限管理
通常在权限管理中主要包含三步：
* 新建用户
* 配置权限
* 配置角色
### 新建用户

```
rabbitmqctl add_user superrd superrd
```

### 配置权限

```
set_permissions [-p <vhostpath>] <user> <conf> <write> <read>1
```

其中， 的位置分别用正则表达式来匹配特定的资源，如

> ‘^(amq.gen.*|amq.default)$’

可以匹配server生成的和默认的exchange，’^$’不匹配任何资源

* exchange和queue的declare与delete分别需要exchange和queue上的配置权限
* exchange的bind与unbind需要exchange的读写权限
* queue的bind与unbind需要queue写权限exchange的读权限 发消息(publish)需exchange的写权限
* 获取或清除(get、consume、purge)消息需queue的读权限

示例：我们赋予superrd在“/”下面的全部资源的配置和读写权限。

```
rabbitmqctl set_permissions -p / superrd ".*" ".*" ".*"
```

> 注意”/”代表virtual host为“/”这个“/”和linux里的根目录是有区别的并不是virtual host为“/”可以访问所以的virtual host，把这个“/”理解成字符串就行。

  ### 配置角色

```
rabbitmqctl set_user_tags [user] [role]1
```

RabbitMQ中的角色分为如下五类：none、management、policymaker、monitoring、administrator

官方解释如下：

```
management 
User can access the management plugin 
policymaker 
User can access the management plugin and manage policies and parameters for the vhosts they have access to. 
monitoring 
User can access the management plugin and see all connections and channels as well as node-related information. 
administrator 
User can do everything monitoring can do, manage users, vhosts and permissions, close other user’s connections, and manage policies and parameters for all vhosts.
```

* none : 不能访问 management plugin

* management : 用户可以通过AMQP做的任何事外加： 列出自己可以通过AMQP登入的virtual hosts ,查看自己的virtual hosts中的queues, exchanges 和 bindings ,查看和关闭自己的channels 和 connections ,查看有关自己的virtual hosts的“全局”的统计信息，包含其他用户在这些virtual hosts中的活动。

* policymaker : management可以做的任何事外加： 查看、创建和删除自己的virtual hosts所属的policies和parameters

* monitoring : management可以做的任何事外加： 列出所有virtual hosts，包括他们不能登录的virtual hosts ,查看其他用户的connections和channels ,查看节点级别的数据如clustering和memory使用情况 ,查看真正的关于所有virtual hosts的全局的统计信息

* administrator : policymaker和monitoring可以做的任何事外加: 创建和删除virtual hosts ,查看、创建和删除users ,查看创建和删除permissions ,关闭其他用户的connections

  

  如下示例将superrd设置成administrator角色。

  ```
  rabbitmqctl set_user_tags superrd administrator
  ```
