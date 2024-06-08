### 1.介绍
Protobuf是Protocol Buffers的简称，它是谷歌公司开发的一种数据描述语言，并于2008年开源。Protobuf刚开源时的定位类似于XML、JSON等数据描述语言，通过附带工具生成代码并实现将结构化数据序列化的功能。

在序列化结构化数据的机制中，ProtoBuf是灵活、高效、自动化的，相对常见的XML、JSON，描述同样的信息，ProtoBuf序列化后数据量更小、序列化/反序列化速度更快、更简单。
### 2. XML、JSON、Protobuf对比
* json: 一般的web项目中，最流行的主要还是json。因为浏览器对于json数据支持非常好，有很多内建的函数支持。
* xml: 在webservice中应用最为广泛，但是相比于json，它的数据更加冗余，因为需要成对的闭合标签。json使用了键值对的方式，不仅压缩了一定的数据空间，同时也具有可读性。
* protobuf:是后起之秀，是谷歌开源的一种数据格式，适合高性能，对响应速度有要求的数据传输场景。因为profobuf是二进制数据格式，需要编码和解码。数据本身不具有可读性。因此只能反序列化之后得到真正可读的数据。
<img width="1024" alt="截屏2024-06-08 20 38 10" src="https://github.com/weifansym/workDoc/assets/6757408/3bf3d3c1-25b4-4072-a132-6f9644c4e512">

### 3. 定义语法
#### 3.1 示例
```
// 指明当前使用proto3语法，如果不指定，编译器会使用proto2
syntax = "proto3";
// package声明符，用来防止消息类型有命名冲突
package msg;
// 选项信息，对应go的包路径
option go_package = "server/msg";
// message关键字，像go中的结构体
message FirstMsg {
  // 类型 字段名 标识号
  int32 id = 1;
  string name=2;
  string age=3;
}
```
#### 3.2 示例说明
* syntax: 用来标记当前使用proto的哪个版本。
* package: 包名，用来防止消息类型命名冲突。
* option go_package: 选项信息，代表生成后的go代码包路径。
* message: 声明消息的关键字，类似Go语言中的struct。
* 定义字段语法: 类型 字段名 标识号

> 标识号说明
> * 每个字段都有唯一的一个数字标识符，一旦开始使用就不能够再改变。
> * [1, 15]之内的标识号在编码的时候会占用一个字节。[16, 2047]之内的标识号则占用2个字节。
> * 最小的标识号可以从1开始，最大到2^29 - 1, or 536,870,911。不可以使用其中的[19000－19999],因为是预留信息，如果使用，编译时会报警。
#### 3.3 简单类型与Go对应
<img width="1009" alt="截屏2024-06-08 20 43 22" src="https://github.com/weifansym/workDoc/assets/6757408/1a4c7e33-f16e-486a-8835-f1dfea034f61">

#### 4. 保留标识符(reserved)
> 什么是保留标示符？reserved 标记的标识号、字段名，都不能在当前消息中使用。

```
syntax = "proto3";
package demo;

// 在这个消息中标记
message DemoMsg {
  // 标示号：1，2，10，11，12，13 都不能用
  reserved 1,2, 10 to 13;
  // 字段名 test、name 不能用
  reserved "test","name";
  // 不能使用字段名，提示:Field name 'name' is reserved
  string name = 3;
  // 不能使用标示号,提示:Field 'id' uses reserved number 11
  int32 id = 11;
}

// 另外一个消息还是可以正常使用
message Demo2Msg {
  // 标示号可以正常使用
  int32 id = 1;
  // 字段名可以正常使用
  string name = 2;
}
```
#### 5. 枚举类型
```
syntax = "proto3";
package demo;
// 声明生成Go代码，包路径
option go_package ="server/demo";
// 枚举消息
message DemoEnumMsg {
  enum Gender{
    // 枚举字段标识符,必须从0开始
    UnKnown = 0;
    Body = 1;
    Girl = 2;
  }
  // 使用自定义的枚举类型
  Gender sex = 2;
}
// 在枚举信息中，重复使用标识符
message DemoTwoMsg{
  enum Animal {
    // 开启允许重复使用 标示符
    option allow_alias = true;
    Other = 0;
    Cat = 1;
    Dog = 2;
    // 白猫也是毛，标示符也用1
    // 不开启allow_alias，会报错： Enum value number 1 has already been used by value 'Cat'
    WhiteCat = 1;
  }
}
```
> 每个枚举类型必须将其第一个类型映射为0, 原因有两个：1.必须有个默认值为0； 2.为了兼容proto2语法，枚举类的第一个值总是默认值.

#### 6.嵌套消息
```
syntax = "proto3";
option go_package = "server/nested";
// 学员信息
message UserInfo {
  int32 userId = 1;
  string userName = 2;
}
message Common {
  // 班级信息
  message CLassInfo{
    int32 classId = 1;
    string className = 2;
  }
}
// 嵌套信息
message NestedDemoMsg {
  // 学员信息 (直接使用消息类型)
  UserInfo userInfo = 1;
  // 班级信息 (通过Parent.Type，调某个消息类型的子类型)
  Common.CLassInfo classInfo =2;
}
```
#### 7.map类型消息
##### 7.1 protobuf源码
```
syntax = "proto3";
option go_package = "server/demo";

// map消息
message DemoMapMsg {
  int32 userId = 1;
  map<string,string> like =2;
}
```
##### 7.2 生成Go代码
```
// map消息
type DemoMapMsg struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserId int32             `protobuf:"varint,1,opt,name=userId,proto3" json:"userId,omitempty"`
	Like   map[string]string `protobuf:"bytes,2,rep,name=like,proto3" json:"like,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}
```
#### 8.切片(数组)类型消息
##### 8.1 protobuf源码
```
syntax = "proto3";
option go_package = "server/demo";

// repeated允许字段重复，对于Go语言来说，它会编译成数组(slice of type)类型的格式
message DemoSliceMsg {
  // 会生成 []int32
  repeated int32 id = 1;
  // 会生成 []string
  repeated string name = 2;
  // 会生成 []float32
  repeated float price = 3;
  // 会生成 []float64
  repeated double money = 4;
}
```
##### 8.2 生成Go代码
```
...
// repeated允许字段重复，对于Go语言来说，它会编译成数组(slice of type)类型的格式
type DemoSliceMsg struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// 会生成 []int32
	Id []int32 `protobuf:"varint,1,rep,packed,name=id,proto3" json:"id,omitempty"`
	// 会生成 []string
	Name []string `protobuf:"bytes,2,rep,name=name,proto3" json:"name,omitempty"`
	// 会生成 []float32
	Price []float32 `protobuf:"fixed32,3,rep,packed,name=price,proto3" json:"price,omitempty"`
	Money []float64 `protobuf:"fixed64,4,rep,packed,name=money,proto3" json:"money,omitempty"`
}
...
```
#### 9. 引入其他proto文件
##### 9.1 被引入文件class.proto
文件位置:proto/class.proto
```
syntax="proto3";
// 包名
package dto;
// 生成go后的文件路径
option go_package = "grpc/server/dto";

message ClassMsg {
  int32  classId = 1;
  string className = 2;
}
```
##### 9.2 使用引入文件user.proto
文件位置:proto/user.proto
```
syntax = "proto3";

// 导入其他proto文件
import "proto/class.proto";

option go_package="grpc/server/dto";

package dto;

// 用户信息
message UserDetail{
  int32 id = 1;
  string name = 2;
  string address = 3;
  repeated string likes = 4;
  // 所属班级
  ClassMsg classInfo = 5;
}
```
##### 9.3 Idea提示:Cannot resolve import..
![image](https://github.com/weifansym/workDoc/assets/6757408/0c6d710a-6bf3-4c1f-b87f-7d0cbcaf56b4)
### 10.定义服务(Service)
> 上面学习的都是怎么定义一个消息类型，如果想要将消息类型用在RPC(远程方法调用)系统中，需要使用关键字(service)定义一个RPC服务接口，使用rpc定义具体方法，而消息类型则充当方法的参数和返回值。

#### 10.1 编写service
文件位置:proto/hello_service.proto
```
syntax = "proto3";

option go_package = "grpc/server";

// 定义入参消息
message HelloParam{
  string name = 1;
  string context = 2;
}

// 定义出参消息
message HelloResult {
  string result = 1;
}

// 定义service
service HelloService{
  // 定义方法 
  rpc SayHello(HelloParam) returns (HelloResult);
}
```
#### 10.2 生成Go代码
> 使用命令:**protoc --go_out=. --go-grpc_out=. proto/hello_service.proto**生成代码。

上述命令会生成很多代码,我们的工作主要是要实现SayHello方法,下面是生成的部分代码。
##### 1.客户端部分代码
```
// 客户端接口
type HelloServiceClient interface {
	SayHello(ctx context.Context, in *HelloParam, opts ...grpc.CallOption) (*HelloResult, error)
}
// 客户端实现调用SayHello方法
func (c *helloServiceClient) SayHello(ctx context.Context, in *HelloParam, opts ...grpc.CallOption) (*HelloResult, error) {
	out := new(HelloResult)
	err := c.cc.Invoke(ctx, "/HelloService/SayHello", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}
```
##### 2. 服务端部分代码
```
// 生成的接口
type HelloServiceServer interface {
	SayHello(context.Context, *HelloParam) (*HelloResult, error)
	mustEmbedUnimplementedHelloServiceServer()
}
// 需要实现SayHello方法
func (UnimplementedHelloServiceServer) SayHello(context.Context, *HelloParam) (*HelloResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SayHello not implemented")
}
```
#### 10.3 实现服务端(SayHello)方法
```
func (UnimplementedHelloServiceServer) SayHello(ctx context.Context, p *HelloParam) (*HelloResult, error) {
	//return nil, status.Errorf(codes.Unimplemented, "method SayHello not implemented")
	return &HelloResult{Result: fmt.Sprintf("%s say %s",p.GetName(),p.GetContext())},nil
}
```
#### 10.4 运行服务
##### 1. 编写服务端
```
package main

import (
	"52lu/go-rpc/grpc/server"
	"fmt"
	"google.golang.org/grpc"
	"net"
)
// 服务端代码
func main() {
	// 创建grpc服务
	rpcServer := grpc.NewServer()
	// 注册服务
	server.RegisterHelloServiceServer(rpcServer, new(server.UnimplementedHelloServiceServer))
	// 监听端口
	listen, err := net.Listen("tcp", ":1234")
	if err != nil {
		fmt.Println("服务启动失败", err)
		return
	}
	fmt.Println("服务启动成功!")
	rpcServer.Serve(listen)
}
```
##### 2. 编写客户端
```
package main

import (
	"52lu/go-rpc/grpc/server"
	"context"
	"fmt"
	"google.golang.org/grpc"
)

// 客户端代码
func main() {
	// 建立链接
	dial, err := grpc.Dial("127.0.0.1:1234", grpc.WithInsecure())
	if err != nil {
		fmt.Println("Dial Error ", err)
		return
	}
	// 延迟关闭链接
	defer dial.Close()
	// 实例化客户端
	client := server.NewHelloServiceClient(dial)
	// 发起请求
	result, err := client.SayHello(context.TODO(), &server.HelloParam{
		Name:    "张三",
		Context: "hello word!",
	})
	if err != nil {
		fmt.Println("请求失败:", err)
		return
	}
	// 打印返回信息
	fmt.Printf("%+v\n", result)
}
```
##### 3.启动运行
```
# 启动服务端
➜ go run server.go
服务启动成功!

# 启动客户端
➜  go run client.go
result:"张三 say hello word!"
```
#### 11. oneof(只能选择一个)
##### 11.1 编写proto
修改上面的服务中的请求参数:HelloParam,
```
// 定义入参消息
message HelloParam{
  string name = 1;
  string context = 2;
  // oneof 最多只能设置其中一个字段
  oneof option {
    int32 age= 3;
    string gender= 4;
  }
}
```
#### 11.2 使用
1.客户端传参
生成Go代码后，入参只能设置其中一个值，如下
```
...省略
	// 实例化客户端
	client := server.NewHelloServiceClient(dial)
	// 定义参数
	reqParam := &server.HelloParam{
		Name:    "张三",
		Context: "hello word!",
	}
	// 只能设置其中一个
	reqParam.Option = &server.HelloParam_Age{Age: 19}
	// 这个会替代上一个值
	//reqParam.Option = &server.HelloParam_Gender{Gender: "男"}
	// 发起请求
	result, err := client.SayHello(context.TODO(), reqParam)
...
```
2. 服务端接收参数
```
func (UnimplementedHelloServiceServer) SayHello(ctx context.Context, p *HelloParam) (*HelloResult, error) {
	//return nil, status.Errorf(codes.Unimplemented, "method SayHello not implemented")
	res := fmt.Sprintf("name:%s| gender:%s| age:%s | context:%s",p.GetName(),p.GetGender(),p.GetAge(),p.GetContext())
	return &HelloResult{Result:res,},nil
}
```
#### 12. Any
##### 12.1 编写proto   
```
syntax = "proto3";

option go_package = "grpc/server";

// 使用any类型，需要导入这个
import "google/protobuf/any.proto";

// 定义准备传的消息
message Context {
  int32 id = 1;
  string title = 2;
}
// 定义入参消息
message HelloParam{
  // any，代表可以是任何类型
  google.protobuf.Any data = 1;
}

// 定义出参消息
message HelloResult {
  string result = 1;
}
```
#### 12.2 使用
1. 客户端传参
```
//...省略...
	// 使用Any参数
	content := &server.Context{
		Id:    100,
		Title: "Hello Word",
	}
	// 序列化Any类型参数
	any, err := anypb.New(content)
	if err != nil {
		fmt.Println("any 类型参数序列化失败")
		return
	}
  // 注意这里，一开始在proto文件中，没有定义使用消息类型Context，
  // 现在通过any类型，同样可以使用
	reqParam := &server.HelloParam{Data: any}
	// 发起请求
	result, err := client.SayHello(context.TODO(), reqParam)
	if err != nil {
		fmt.Println("请求失败:", err)
		return
	}
// ....省略...
```
2. 服务端接收参数
```
func (UnimplementedHelloServiceServer) SayHello(ctx context.Context, p *HelloParam) (*HelloResult, error) {
	var context Context
  // 反序列化参数
	p.Data.UnmarshalTo(&context)
	res := fmt.Sprintf("%v|%v",context.GetId(),context.GetTitle())
	return &HelloResult{Result: res},nil
	//return nil, status.Errorf(codes.Unimplemented, "method SayHello not implemented")
}
```
   
