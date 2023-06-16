gorm在插入完成数据后,想要获取插入的自增id ,可以使用Create()方法执行插入,在结构体里直接就能获取到ID

例如下面这个User 
```
type User struct {
  Model
  Name string `json:"name"`
  Password string `json:"password"`
  Nickname string `json:"nickname"`
  Avator string `json:"avator"`
  RoleName string `json:"role_name" sql:"-"`
}
func CreateUser(name string,password string,avator string,nickname string)uint{
  user:=&User{
    Name:name,
    Password: password,
    Avator:avator,
    Nickname: nickname,
  }
  DB.Create(user)
  return user.ID
}
```
当RoleName这个成语不想映射到字段里的时候 `sql:"-"`
```
RoleName string `json:"role_name" sql:"-"`
```
测试效果可以直接点击与我交流

代码地址:

https://github.com/taoshihan1991/go-fly
