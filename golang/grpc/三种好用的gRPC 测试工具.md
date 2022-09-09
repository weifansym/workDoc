## 三种好用的gRPC 测试工具
![image](https://user-images.githubusercontent.com/6757408/189290554-78fb8e7c-426a-4179-bc2a-98bdc2085b49.png)

最近在用Go 语言实作微服务，沟通的接口采用[gRPC](https://grpc.io/)，除了可以透过gRPC支援的[第三方语言](https://grpc.io/docs/languages/)来写客户端的测试之外，有没有一些好用的工具来验证检查gRPC 实现的接口。
刚好今年看到Postman宣布开始[支援gRPC](https://blog.postman.com/postman-now-supports-grpc/)，相信大家对于Postman 工具并不会太陌生，毕竟测试[Websocket](https://blog.postman.com/postman-supports-websocket-apis/)或RESTful API 都是靠这工具呢。本篇除了介绍Postman 之外，
还有一套CLI 工具[grpcurl](https://github.com/fullstorydev/grpcurl)及一套GUI 工具[grpcui](https://github.com/fullstorydev/grpcui)也是不错用，后面这两套都是由同一家公司[FullStory](https://www.fullstory.com/blog/tag/engineering/)开源出来的专案，底下就来一一介绍。
### gRPC 服务范本
用Go 语言写好一个[测试范例版本](https://github.com/go-training/proto-go-sample)，gRPC 定义的proto 档案可以[从这边查看](https://github.com/go-training/proto-def-demo)
```
syntax = "proto3";

package gitea.v1;

message GiteaRequest {
  string name = 1;
}

message GiteaResponse {
  string giteaing = 1;
}

message IntroduceRequest {
  string name = 1;
}

message IntroduceResponse {
  string sentence = 1;
}

service GiteaService {
  rpc Gitea(GiteaRequest) returns (GiteaResponse) {}
  rpc Introduce(IntroduceRequest) returns (stream IntroduceResponse) {}
}
```
### grpcurl
相信大家都有用过强大的[curl](https://curl.se/)工具，而[grpcurl](https://github.com/fullstorydev/grpcurl)可以想像成curl 的gRPC 版本，安装方式非常简单，可以到[Relase](https://github.com/fullstorydev/grpcurl/releases) 页面找你要的OS 版本，如果本身是Go 开发者可以透过go install 安装
```
go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest
```
MacOS 环境可以透过Homebrew 安装
```
brew install grpcurl
```
或者透过Docker 来执行也可以的
```
# Download image
docker pull fullstorydev/grpcurl:latest
# Run the tool
docker run fullstorydev/grpcurl api.grpc.me:443 list
```
先把测试环境搭建起来，就可以在本地端测试8080 连接埠
```
$ grpcurl -plaintext localhost:8080 list
gitea.v1.GiteaService
grpc.health.v1.Health
ping.v1.PingService
```
由于本地端没有SSL 凭证，所以请加上 -plaintext 参数才可以使用，详细参数可透过 -h 查看，接下来打第一个测试接口
```
grpcurl \
  -plaintext \
  -d '{"name": "foobar"}' \
  localhost:8080 \
  gitea.v1.GiteaService/Gitea
```
会拿到底下结果
```
{
  "giteaing": "Hello, foobar!"
}
```
如果服务有加上Health Check 可以直接使用底下指令
```
grpcurl \
  -plaintext \
  -d '{"service": "gitea.v1.GiteaService"}' \
  localhost:8080 \
  grpc.health.v1.Health/Check
```
可以拿到底下结果
```
{
  "status": "SERVING"
}
```
服务如果有支援Server Streaming RPC 也是同样用法
```
grpcurl \
  -plaintext \
  -d '{"name": "foobar"}' \
  localhost:8080 \
  gitea.v1.GiteaService/Introduce
```
可以看到Server 会回应两个讯息
```
{
  "sentence": "foobar, How are you feeling today 01 ?"
}
{
  "sentence": "foobar, How are you feeling today 02 ?"
}
```
### grpcui
除了上面grpcurl 外，同一个团队也推出Web UI 的gRPC 测试工具[grpcui](https://github.com/fullstorydev/grpcui)，安装方式跟grpcurl 一样，这边就不多做说明了，grpcui 也支援全部RPC 功能，包含streaming 等。
底下一行指令就可以启动GUI 画面了
```
grpcui -plaintext localhost:8080
```
![image](https://user-images.githubusercontent.com/6757408/189292174-1af9938d-c095-49ff-bd62-4bdbd7e32e18.png)

此页面已经帮忙把Service 及可以用的Method 都准备完毕了，你也不用知道任何服务提供哪些Method，相当方便，选择不同的Service 及Method，画面上的Request Data 都会随着变动

![image](https://user-images.githubusercontent.com/6757408/189292241-e879f9d4-b948-4427-b621-cd6d7e475ac0.png)

填写完Request Form 在切换到Raw Request (JSON) 就可以看到JSON Format 的资料

![image](https://user-images.githubusercontent.com/6757408/189292341-63980564-ef3d-4813-a897-27fa05a184fc.png)

按下invoke 后，可以看到结果

![image](https://user-images.githubusercontent.com/6757408/189292403-9023b90e-312e-4650-9fe2-b889e0956413.png)

测试Streaming 结果

![image](https://user-images.githubusercontent.com/6757408/189292992-d0fdf70c-7289-452c-9a09-ce5aaef807f2.png)

### Postman
相信大家最熟悉的还是Postman，测试任何服务都离不开此工具，也很高兴今年一月看到支援了gRPC Beta 版本，打开软体后，按下左上角New 就可以看到底下画面，选择gRPC Request
![image](https://user-images.githubusercontent.com/6757408/189293113-1a09bf8d-deb5-4c82-ac55-4900123d1d78.png)

左边API 可以把所有服务的proto 资料写进去

![image](https://user-images.githubusercontent.com/6757408/189293190-b899ecfd-a451-4216-8d52-8c892bbb19e5.png)

接着就可以透过New Request 来选择了，并且执行invoke

![image](https://user-images.githubusercontent.com/6757408/189293269-4722650e-3a8f-4c28-95ab-ed609e13724b.png)

Stream Testing 如下

![image](https://user-images.githubusercontent.com/6757408/189293349-08fab62b-6bcc-4dc3-bd89-1f113af4b075.png)

### 心得
grpcurl CLI 工具可以在基本没有桌面的Linux 环境中使用，所以这边算蛮推荐的，但是Postman 在使用前都需要手动将Proto 的资料写进去，才可以测试使用，所以我更强烈推荐使用grpcui 
可以直接动态读取proto 资讯，将资讯转成Request Form，减少查询资料属性的时间。






