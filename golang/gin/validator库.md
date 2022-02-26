## validator库使用
在做API部分开发时，需要对请求参数的校验，防止用户的恶意请求。例如日期格式，用户年龄，性别等必须是正常的值，不能随意设置。
其实这种场景，我们其实不需要自己编写代码来进行参数校验的，应该会有对应的package进行处理的，这里就介绍一个常用的package，同时这个包还是GIN框架默认的参数校验包。
### 快速安装
使用之前，我们先要获取validator这个库。
```
# 第一次安装使用如下命令
$ go get github.com/go-playground/validator/v10
# 项目中引入包
import "github.com/go-playground/validator/v10"
```
### 简单示例
安装还是很简单的，下面我先来一个官方样例，看看是怎么使用的，然后展开分析。
```
package main

import (
 "fmt"
 "net/http"

 "github.com/gin-gonic/gin"
)

type RegisterRequest struct {
 Username string `json:"username" binding:"required"`
 Nickname string `json:"nickname" binding:"required"`
 Email    string `json:"email" binding:"required,email"`
 Password string `json:"password" binding:"required"`
 Age      uint8  `json:"age" binding:"gte=1,lte=120"`
}

func main() {

 router := gin.Default()

 router.POST("register", Register)

 router.Run(":9999")
}

func Register(c *gin.Context) {
 var r RegisterRequest
 err := c.ShouldBindJSON(&r)
 if err != nil {
  fmt.Println("register failed")
  c.JSON(http.StatusOK, gin.H{"msg": err.Error()})
  return
 }
 //验证 存储操作省略.....
 fmt.Println("register success")
 c.JSON(http.StatusOK, "successful")
}
```
具体内容可以参加：https://mp.weixin.qq.com/s?__biz=MzkyNzI1NzM5NQ==&mid=2247484752&idx=1&sn=24a691b9305df828c24b5d9a56f25b46&scene=21#wechat_redirect
