## RPC编程(三):gRPC快速入门
### 1.什么是gRPC
[gRPC](https://grpc.io/)是一个高性能、开源、通用的RPC框架，由Google推出，基于[HTTP2](https://http2.github.io/)协议标准设计开发，默认采用[Protocol Buffers](https://developers.google.com/protocol-buffers/)数据序列化协议，支持多种开发语言。
gRPC提供了一种简单的方法来精确的定义服务，并且为客户端和服务端自动生成可靠的功能库。
### grpc技术栈
Go 语言的 gRPC 技术栈

![image](https://user-images.githubusercontent.com/6757408/189268332-dd8b80f3-a6b2-427a-bf1a-18c4cd1bb554.png)

> 最底层为TCP或Unix套接字协议，在此之上是HTTP/2协议的实现，然后在HTTP/2协议之上又构建了针对Go语言的gRPC核心库（gRPC内核+解释器）。应用程序通过gRPC插件生成的Stub代码和gRPC核心库通信，也可以直接和gRPC核心库通信。


转自：
* http://liuqh.icu/2022/01/20/go/rpc/03-grpc-ru-men/
* https://chai2010.cn/advanced-go-programming-book/ch4-rpc/ch4-04-grpc.html
