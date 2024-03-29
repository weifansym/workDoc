## Gin框架(二):服务启动源码分析
### 1.启动服务
#### 1.1 服务源码
```
package main
// 引入gin框架
import "github.com/gin-gonic/gin"
func main() {
	// 创建一个默认的路由引擎
	engine := gin.Default()
	// 注册Get路由
	engine.GET("/", func(ctx *gin.Context) {
		ctx.JSON(200,gin.H{
			"msg":"请求成功",
		})
	})
  // 默认监听的是 0.0.0.0:8080
	_ = engine.Run()
}
```
#### 1.2 启动输出
```
➜ go run main.go 
[GIN-debug] [WARNING] Creating an Engine instance with the Logger and Recovery middleware already attached.

[GIN-debug] [WARNING] Running in "debug" mode. Switch to "release" mode in production.
 - using env:   export GIN_MODE=release
 - using code:  gin.SetMode(gin.ReleaseMode)

[GIN-debug] GET    /                         --> main.main.func1 (3 handlers)
[GIN-debug] Environment variable PORT is undefined. Using port :8080 by default
[GIN-debug] Listening and serving HTTP on :8080
```
### 2.代码分析
#### 2.1 import “github.com/gin-gonic/gin”
在引入Gin框架时，包内相关的init()方法都会被执行；经查找发现下面两个init方法;

a. 第一个init方法

方法位置: github.com/gin-gonic/gin/context_appengine.go
```
package gin
func init() {
  // 设置AppEngine = true
	defaultAppEngine = true
}
```
b. 第二个init方法

方法位置: github.com/gin-gonic/gin/mode.go
```
func init() {
  // 设置服务的运行模式，默认是DebugMode
  // 分别有三种模式:DebugMode=0(开发模式)、releaseCode=1(生产模式)、testCode=2(测试模式)
	mode := os.Getenv(EnvGinMode)
	SetMode(mode)
}
```
#### 2.2 gin.Default()
gin.Default源码如下:
```
func Default() *Engine {
  // 打印gin-debug信息
	debugPrintWARNINGDefault()
  // 新建一个无路由无中间的引擎
	engine := New()
  // 注册全局日志和异常捕获中间件
	engine.Use(Logger(), Recovery())
	return engine
}
```
> 注意: Gin框架中注册中间件是通过 engine.Use(xx)的方式。 

2.3 engine.GET(“/“,…)

1.源码
```
// 注册一个匹配路径(relativePath)的Get请求路由
// handlers是对应的处理逻辑
func (group *RouterGroup) GET(relativePath string, handlers ...HandlerFunc) IRoutes {
	return group.handle(http.MethodGet, relativePath, handlers)
}
```
在上述示例代码中我们注册了一个匹配根目录("/")的路由，处理handlers是匿名函数,直接调用ctx.JSON返回json格式的数据；
```
// 注册Get路由
engine.GET("/", func(ctx *gin.Context) {
  ctx.JSON(200,gin.H{
    "msg":"请求成功",
  })
})
```
2. 多种返回格式

在gin.Context支持多种返回格式,整理常用的返回格式如下:
| 方法名  | 描述 |
| ------------- | ------------- |
| ctx.XML(code int, obj interface{})  | 返回xml  |
| ctx.AsciiJSON(code int, obj interface{})	  | 返回json,将使特殊字符编码  |
| ctx.PureJSON(code int, obj interface{})	  | 返回json，有html的不转义。  |

3.Json、AsciiJSON、PureJSON 对比
```
package main
import "github.com/gin-gonic/gin"
func main() {
	// 创建一个默认的路由引擎
	engine := gin.Default()
	// 注册Get路由
	engine.GET("/", func(ctx *gin.Context) {
		key, _ := ctx.GetQuery("key")
		msgBody := gin.H{
			"msg": "请求成功",
			"html":"<span>我是一段html代码</span>",
		}
		switch key {
		case "1":
			msgBody["method"] = "ctx.JSON"
			ctx.JSON(200, msgBody)
		case "2":
			msgBody["method"] = "ctx.PureJSON"
			ctx.PureJSON(200, msgBody)
		case "3":
			msgBody["method"] = "ctx.AsciiJSON"
			ctx.AsciiJSON(200, msgBody)
		default:
			ctx.JSON(500, gin.H{
				"msg": "请求失败",
			})
		}
		return
	})
	_ = engine.Run()
}

```
请求返回:
```
➜  ~ curl http://127.0.0.1:8080/\?key\=1
{"html":"\u003cspan\u003e我是一段html代码\u003c/span\u003e","method":"ctx.JSON","msg":"请求成功"}%

➜  ~ curl http://127.0.0.1:8080/\?key\=2
{"html":"<span>我是一段html代码</span>","method":"ctx.PureJSON","msg":"请求成功"}

➜  ~ curl http://127.0.0.1:8080/\?key\=3
{"html":"\u003cspan\u003e\u6211\u662f\u4e00\u6bb5html\u4ee3\u7801\u003c/span\u003e","method":"ctx.AsciiJSON","msg":"\u8bf7\u6c42\u6210\u529f"}%
```
总结
| 方法名  | 现象 |
| ------------- | ------------- |
| ctx.JSON  | 默认会把html转成unicode字符,对汉字不做额外处理  |
| ctx.PureJSON	  | 会把html原样返回，,对汉字不做额外处理  |
| ctx.AsciiJSON  | 会对汉字和html都做处理。  |

2.4 engine.Run()
1.Run源码如下
```
func (engine *Engine) Run(addr ...string) (err error) {
  // 延迟关闭输出ERROR类型的日志信息
	defer func() { debugPrintError(err) }()
  // 设置CIDR（无类型域间路由)信息，默认返回: 0.0.0.0/0
	trustedCIDRs, err := engine.prepareTrustedCIDRs()
	if err != nil {
		return err
	}
	engine.trustedCIDRs = trustedCIDRs
  // 设置监听IP和端口信息，默认是":8080"
	address := resolveAddress(addr)
	debugPrint("Listening and serving HTTP on %s\n", address)
  // 启动服务
	err = http.ListenAndServe(address, engine)
	return
}
```
2.为什么默认监听是”:8080”
在Run方法中调用 resolveAddress(addr),该方法源码如下:
```
// 接收一个字符串切片参数
func resolveAddress(addr []string) string {
  // 如果参数长度为0，默认监听8080
	switch len(addr) {
	case 0:
		if port := os.Getenv("PORT"); port != "" {
			debugPrint("Environment variable PORT=\"%s\"", port)
			return ":" + port
		}
		debugPrint("Environment variable PORT is undefined. Using port :8080 by default")
		return ":8080"
	case 1:
      // 如果参数长度为1，监听IP和端口
		return addr[0]
	default:
     //  如果参数长度大于1，则报错
		panic("too many parameters")
	}
}
```
