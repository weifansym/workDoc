## Gin框架(十一):验证器Validator使用
### 介绍
[validator](https://github.com/go-playground/validator)是一个开源的验证器包，可以快速校验输入信息是否符合自定规则。目前Star 7.8k,源码地址: https://github.com/go-playground/validator
#### 1.1 安装
```
go get github.com/go-playground/validator
```
#### 1.2 引用
```
import "github.com/go-playground/validator"
```
#### 1.3 示例
代码:
```
package main
import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)
// 定义一个添加用户参数结构体
type AddUserPost struct {
	Name  string `json:"name" validate:"required"`        //必填
	Email string `json:"email" validate:"required,email"` // 必填，并且格式是email
	Age   uint8    `json:"age" validate:"gte=18,lte=30"` // 年龄范围
}
// 简单示例
func main() {
	engine := gin.Default()
	engine.POST("/valid", func(context *gin.Context) {
		var adduserPost  AddUserPost
		// 接收参数
		err := context.ShouldBindJSON(&adduserPost)
		if err != nil {
			context.JSON(500,gin.H{"msg":err})
			return
		}
		fmt.Printf("adduserPost: %+v\n",adduserPost)
		// 使用Validate验证
		validate := validator.New()
		err = validate.Struct(adduserPost)
		if err != nil {
			fmt.Println(err)
			context.JSON(500,gin.H{"msg":err.Error()})
			return
		}
		context.JSON(200,gin.H{"msg":"success"})
	})
	_ = engine.Run()
}
```
请求:
```
# email不合法时
➜ curl -X POST http://127.0.0.1:8080/valid -d '{"name":"张三","email":"123","age":21}'
{"msg":"Key: 'AddUserPost.Email' Error:Field validation for 'Email' failed on the 'email' tag"}
# age 不在指定范围时
➜ curl -X POST http://127.0.0.1:8080/valid -d '{"name":"张三","email":"123@163.com","age":17}'
{"msg":"Key: 'AddUserPost.Age' Error:Field validation for 'Age' failed on the 'gte' tag"}
# 姓名不填时
➜ curl -X POST http://127.0.0.1:8080/valid -d '{"name":"","email":"123@163.com","age":20}'
{"msg":"Key: 'AddUserPost.Name' Error:Field validation for 'Name' failed on the 'required' tag"}
```
### 2.改成中文
#### 2.1 代码
```
package main
import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zhs "github.com/go-playground/validator/v10/translations/zh"
)
var (
	validate = validator.New()          // 实例化验证器
	chinese  = zh.New()                 // 获取中文翻译器
	uni      = ut.New(chinese, chinese) // 设置成中文翻译器
	trans, _ = uni.GetTranslator("zh")  // 获取翻译字典
)
type User struct {
	Name  string `form:"name" validate:"required,min=3,max=5"`
	Email string `form:"email" validate:"email"`
	Age   int8   `form:"age" validate:"gte=18,lte=20"`
}

func main() {
	engine := gin.Default()
	engine.GET("/language", func(context *gin.Context) {
		var user User
		err := context.ShouldBindQuery(&user)
		if err != nil {
			context.JSON(500, gin.H{"msg": err})
			return
		}
		// 注册翻译器
		_ = zhs.RegisterDefaultTranslations(validate, trans)
		// 使用验证器验证
		err = validate.Struct(user)
		if err != nil {
			if errors, ok := err.(validator.ValidationErrors); ok {
				// 翻译，并返回
				context.JSON(500, gin.H{
					"翻译前": errors.Error(),
					"翻译后": errors.Translate(trans),
				})
				return
			}
		}
		context.JSON(200,gin.H{"msg":"success"})
	})
	_ = engine.Run()
}
```
#### 2.2 请求
```
# 不传参数
➜ curl -X GET http://127.0.0.1:8080/language
{
    "翻译前":"Key: 'User.Name' Error:Field validation for 'Name' failed on the 'required' tag
Key: 'User.Email' Error:Field validation for 'Email' failed on the 'email' tag
Key: 'User.Age' Error:Field validation for 'Age' failed on the 'gte' tag",
    "翻译后":{
        "User.Age":"Age必须大于或等于18",
        "User.Email":"Email必须是一个有效的邮箱",
        "User.Name":"Name为必填字段"
    }
}
```
#### 3.校验规则
整理一些常用的规则


参见：http://liuqh.icu/2021/05/30/go/gin/11-validate/



