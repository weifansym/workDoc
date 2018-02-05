## API接口定义规范
本规范主要是基于RESTful风格
### 协议
API与用户的通信协议，总是使用HTTPs协议，确保交互数据的传输安全。

### 域名
应该尽量将API部署在专用域名之下。
```
https://api.example.com
```
如果确定API很简单，不会有进一步扩展，可以考虑放在主域名下。
```
https://example.org/api/
```
### api版本控制
应该将API的版本号放入URL。
https://api.example.com/v{n}/
另一种做法是，将版本号放在HTTP头信息中，但不如放入URL方便和直观。Github采用这种做法。
```
采用多版本并存，增量发布的方式
v{n} n代表版本号,分为整形和浮点型
整形的版本号: 大功能版本发布形式；具有当前版本状态下的所有API接口 ,例如：v1,v2
浮点型：为小版本号，只具备补充api的功能，其他api都默认调用对应大版本号的api 例如：v1.1 v2.2
```

### API 路径规则
路径又称"终点"（endpoint），表示API的具体网址。
在RESTful架构中，每个网址代表一种资源（resource），所以网址中不能有动词，只能有名词，而且所用的名词往往与数据库的表格名对应。一般来说，数据库中的表都是同种记录的"集合"（collection），所以API中的名词也应该使用复数。
举例来说，有一个API提供动物园（zoo）的信息，还包括各种动物和雇员的信息，则它的路径应该设计成下面这样。
```
https://api.example.com/v1/products
https://api.example.com/v1/users
https://api.example.com/v1/employees
```
### HTTP请求方式
对于资源的具体操作类型，由HTTP动词表示。
常用的HTTP动词有下面四个（括号里是对应的SQL命令）。
```
GET（SELECT）：从服务器取出资源（一项或多项）。
POST（CREATE）：在服务器新建一个资源。
PUT（UPDATE）：在服务器更新资源（客户端提供改变后的完整资源）。
DELETE（DELETE）：从服务器删除资源。
```
下面是一些例子。
```
GET /products：列出所有商品
POST /products：新建一个商品
GET /products/ID：获取某个指定商品的信息
PUT /products/ID：更新某个指定商品的信息
DELETE /products/ID：删除某个商品
GET /products/ID/purchases ：列出某个指定商品的所有投资者
get /products/ID/purchases/ID：获取某个指定商品的指定投资者信息
```
### 过滤信息
如果记录数量很多，服务器不可能都将它们返回给用户。API应该提供参数，过滤返回结果。
```
下面是一些常见的参数。
?limit=10：指定返回记录的数量
?offset=10：指定返回记录的开始位置。
?page=2&per_page=100：指定第几页，以及每页的记录数。
?sortby=name&order=asc：指定返回结果按照哪个属性排序，以及排序顺序。
?producy_type=1：指定筛选条件
```
### API 传入参数
参入参数分为4种类型：
```
地址栏参数
* restful 地址栏参数 /api/v1/product/122 122为产品编号，获取产品为122的信息
* get方式的查询字串 见过滤信息小节
请求body数据
cookie
request header
```
cookie和header 一般都是用于OAuth认证的2种途径

### 返回数据
只要api接口成功接到请求，就不能返回200以外的HTTP状态。
为了保障前后端的数据交互的顺畅，建议规范数据的返回，并采用固定的数据格式封装。

接口返回模板：
```
{
status:0,
data:{}||[],
msg:’’
}
```
status: 接口的执行的状态
```
=0表示成功
<0 表示有异常
》0 表示接口有部分执行失败
```
Data 接口的主数据：
```
可以根据实际返回数组或JSON对象
```
Msg
```
当status!=0 都应该有错误信息
```
### 非Restful Api的需求
```
由于实际业务开展过程中，可能会出现各种的api不是简单的restful 规范能实现的，因此，需要有一些api突破restful规范原则。
特别是移动互联网的api设计，更需要有一些特定的api来优化数据请求的交互。
```
### 页面级的api
把当前页面中需要用到的所有数据通过一个接口一次性返回全部数据

举例
api/v1/get-home-data 返回首页用到的所有数据
```
这类API有一个非常不好的地址，只要业务需求变动，这个api就需要跟着变更。
```
### 自定义组合api
```
把当前用户需要在第一时间内容加载的多个接口合并成一个请求发送到服务端，
服务端根据请求内容，一次性把所有数据合并返回,相比于页面级api，具备更高的灵活性，同时又能很容易的实现页面级的api功能。
```
规范
地址：api/v1/batApi

传入参数：
```
data:[
{url:'api1',type:'get',data:{...}},
{url:'api2',type:'get',data:{...}},
{url:'api3',type:'get',data:{...}},
{url:'api4',type:'get',data:{...}}
]

返回数据
{status:0,msg:'',
data:[
{status:0,msg:'',data:[]},
{status:-1,msg:'',data:{}},
{status:1,msg:'',data:{}},
{status:0,msg:'',data:[]},
]
}
```
### Api共建平台
RAP是一个GUI的WEB接口管理工具。在RAP中，您可定义接口的URL、请求&响应细节格式等等。通过分析这些数据，RAP提供MOCK服务、测试服务等自动化工具。
RAP同时提供大量企业级功能，帮助企业和团队高效的工作。

什么是RAP?
在前后端分离的开发模式下，我们通常需要定义一份接口文档来规范接口的具体信息。如一个请求的地址、有几个参数、参数名称及类型含义等等。
RAP 首先方便团队录入、查看和管理这些接口文档，并通过分析结构化的文档数据，重复利用并生成自测数据、提供自测控制台等等... 大幅度提升开发效率。

RAP的特色
强大的GUI工具 给力的用户体验，你将会爱上使用RAP来管理您的API文档。
完善的MOCK服务 文档定义好的瞬间，所有接口已经准备就绪。有了MockJS，无论您的业务模型有多复杂，它都能很好的满足。
庞大的用户群 RAP在阿里巴巴有200多个大型项目在使用，也有许多著名的公司、开源人士在使用。RAP跟随这些业务的成行而成长，专注细节，把握质量，经得住考验。
免费 + 专业的技术支持 RAP是免费的，而且你的技术咨询都将在24小时内得到答复。大多数情况，在1小时内会得到答复。
RAP是一个可视化接口管理工具 通过分析接口结构，动态生成模拟数据，校验真实接口正确性， 围绕接口定义，通过一系列自动化工具提升我们的协作效率。
我们的口号：提高效率，回家吃晚饭！

上手视频
片源：淘宝视频 优酷

RAP学习中心
定期更新的分步骤视频教程

http://thx.github.io/RAP/study.html

参考链接
http://www.ruanyifeng.com/blog/2014/05/restful_api.html
https://github.com/thx/RAP/wiki/about_cn
http://www.vinaysahni.com/best-practices-for-a-pragmatic-restful-api

转自：https://github.com/mishe/blog/issues/129
