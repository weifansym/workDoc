### Gin框架(八):中间件
#### 1.介绍
在Gin框架中,中间件本质上是gin.HandlerFunc 函数，如果我们要自定义中间件，只需要返回类型是gin.HandlerFunc 即可。

在Gin框架中,使用中间件可分为以下几种场景:
* 全局使用
* 单个路由使用
* 路由组使用

#### 2.使用
##### 2.1 全局使用
在gin.Default()函数中，默认注册全局中间件Logger、Recovery，具体代码如下：
```
func Default() *Engine {
  // 打印debug信息
	debugPrintWARNINGDefault()
  // 创建引擎
	engine := New()
  // 注册全局中间件
	engine.Use(Logger(), Recovery())
	return engine
}
```
##### 2.2 单个路由使用
```
func main()  {
	engine := gin.New()
	// 添加一个中间件: Logger()
	engine.GET("/route",gin.Logger(), func(context *gin.Context) {
		context.JSON(200,gin.H{"msg":"针对单个路由"})
	})
	// 添加多个中间件: Logger(),Recovery()
	engine.GET("/route2",gin.Logger(),gin.Recovery(), func(context *gin.Context) {
		context.JSON(200,gin.H{"msg":"针对单个路由添加多个中间件"})
	})
	_ = engine.Run()
}
```
##### 2.3 路由组使用
```
func main()  {
	engine := gin.New()
  // 声明路由后，通过Use注册中间件
	v1Group := engine.Group("/v1").Use(gin.Logger())
	{
		v1Group.GET("/middleGroup", func(context *gin.Context) {
			 context.JSON(200,gin.H{"msg":"succss"})
		})
	}
	_ = engine.Run()
}
```
#### 3.自定义中间件
我们可以创建自己的中间件，在中间件中需要继续执行时,使用c.Next(),而要终止执行时则需调用c.Abort()终止请求。
##### 3.1 语法结构
```
func MiddleName() gin.HandlerFunc {
	return func(context *gin.Context) {
		// 请求前业务逻辑...
    // 继续往下执行
		context.Next()
    // context.Abort() 停止往下执行
		// 请求后业务逻辑...
		end := time.Now().Unix()
		fmt.Printf("接口耗时:%v 秒 \n", end - t)
	}
}
```
代码说明:
* MiddleName: 自定义中间件名称。
* gin.HandlerFunc: 中间件始终返回这个类型。
* func(context *gin.Context): 匿名函数,参数是上下文(context *gin.Context)
* context.Next(): 继续执行调用这个函数。
* context.Abort()：终止执行调用这个函数

##### 3.2 示例
```
// 声明自定义中间件
func MyMiddleware() gin.HandlerFunc {
	return func(context *gin.Context) {
		// 请求前
		fmt.Println("中间件-- 请求前")
		t := time.Now().Unix()
    // 继续往下执行
		context.Next()
    // context.Abort() 停止往下执行
		// 请求后
		end := time.Now().Unix()
		fmt.Printf("接口耗时:%v 秒 \n", end - t)
	}
}
// 自定义中间件使用
func RunWithMyMiddle()  {
	engine := gin.New()
	// 注册自定义中间件
	engine.GET("/route",MyMiddleware(), func(context *gin.Context) {
		time.Sleep(time.Second * 3)
		context.JSON(200,gin.H{"msg":"自定义路由"})
	})
	_ = engine.Run()
}
```
##### 4.多个中间件
如果同时注册多个中间件，那么他们之间的执行顺序是什么样呢? 下面是GinContext的两个执行模型。分别对应中间件中使用Next()和不使用Next()。
![image](https://user-images.githubusercontent.com/6757408/187849758-6ba0357c-0b2a-4c05-a4e2-dc8b78aafeb0.png)

##### 4.2 代码示例
```
// 自定义中间件A
func MyMiddleware() gin.HandlerFunc {
	return func(context *gin.Context) {
		// 请求前
		fmt.Println("中间件1-- 请求前")
		context.Next()
		// 请求后
		fmt.Println("中间件1-- 请求后")
	}
}
// 自定义中间件B
func MyMiddleware2() gin.HandlerFunc {
	return func(context *gin.Context) {
		// 请求前
		fmt.Println("中间件2-- 请求前")
		context.Next()
		// 请求后
		fmt.Println("中间件2-- 请求后")
	}
}
// 启动服务
func main()  {
	engine := gin.New()
	// 使用多个中间件
	engine.GET("/route",MyMiddleware(), MyMiddleware2(),func(context *gin.Context) {
		time.Sleep(time.Second * 3)
		context.JSON(200,gin.H{"msg":"自定义路由"})
	})
	_ = engine.Run()
}
```
请求控制台输出:
```
中间件1-- 请求前
中间件2-- 请求前
中间件2-- 请求后
中间件1-- 请求后
```
#### 5.实践
创建自定义中间件用来判断Token是否有效(是否是登录状态)。
##### 5.1 源码
./mian.go代码:
```
package main
import "go-use/practise"
func main() {
	practise.RunServeWithCheckToken()
}
```
./practise/check_token.go 代码:
```
package practise
import (
	"github.com/gin-gonic/gin"
)
// 检测token
func CheckTokenMiddle() gin.HandlerFunc {
	return func(context *gin.Context) {
		// 获取token
		token := context.DefaultQuery("token","")
		// 检测token
		if token != "abcd" {
			context.JSON(500,gin.H{"msg":"请先登录!"})
			// 终止请求
			context.Abort()
		}
		// token合法则继续往下执行
		context.Next()
	}
}
// 启动服务
func RunServeWithCheckToken()  {
	engine := gin.Default()
	// 登录接口不需要添加检测Token中间件
	engine.GET("/user/login", func(context *gin.Context) {
		context.JSON(200,gin.H{"msg":"登录成功!","token":"abcd"})
	})
	// 需要验证token的路由
	user := engine.Group("/user").Use(CheckTokenMiddle())
	{
		// 用户基本信息
		user.GET("/info", func(context *gin.Context) {
			data := map[string]interface{} {
				"name": "张三",
				"age": 18,
				"likes": []string{"打游戏","旅游"},
			}
			context.JSON(200,gin.H{"msg":"请求成功","data":data})
		})
		// 更新用户
		user.GET("/update", func(context *gin.Context) {
			context.JSON(200,gin.H{"msg":"请求成功"})
		})
	}
	// 启动服务
	_ = engine.Run()
}
```
##### 5.2 请求
```
# 不需要验证token的接口
➜ curl -X GET http://127.0.0.1:8080/user/login
{"msg":"登录成功!","token":"abcd"}
# 需要验证Token，但是没传时
➜ curl -X GET http://127.0.0.1:8080/user/info
{"msg":"请先登录!"}
# 需要验证Token，并且已传时
➜  curl -X GET http://127.0.0.1:8080/user/info?token=abcd
{"data":{"age":18,"likes":["打游戏","旅游"],"name":"张三"},"msg":"请求成功"}
```
#### 6.使用Goroutine
当在中间件或 handler 中启动新的 Goroutine 时，不能使用原始的上下文，必须使用只读副本。
##### 6.1 代码示例
```
// 在中间件中使用Goroutine
func main()  {
	engine := gin.Default()
	engine.Use(MyMiddleWithGoRoutine())
	engine.GET("/useGo", func(context *gin.Context) {
		context.JSON(200,gin.H{"msg":"success"})
	})
	_ = engine.Run()
}
// 在中间件中使用Goroutine
func MyMiddleWithGoRoutine() gin.HandlerFunc {
	return func(context *gin.Context) {
		// 复制上下文
		cpContext := context.Copy()
		go func() {
			time.Sleep(3 * time.Second)
			fmt.Println(cpContext.Request.URL.Path)
		}()
		context.Next()
	}
}
```
#### 7.流行的中间件
在[contrib](https://github.com/gin-gonic/contrib)库中整理很多常用的中间件，可以直接使用

转自：http://liuqh.icu/2021/05/15/go/gin/8-middleware/


