## 如何告诉gorm将缺少的time.Time字段保存为NULL，而不是'0000-00-00'？
我为用户提供了以下自定义类型：
```
type User struct {
    UserID    uint64 `gorm:"primaryKey"`
    CreatedAt time.Time
    UpdatedAt time.Time
    LastLogin time.Time
}
```
当传递给gorm的db.Create()方法时，用户初始化如下：
```
return User{
    UserID:    userID,
    AccountID: accountID,
    UserType:  userType,
    Status:    status,
    UserInfo:  &userInfo,
    CreatedAt: now,
    UpdatedAt: now,
}
```
因为LastLogin在MySQL中是一个可以为空的timestamp列，所以我没有在这里初始化它的值。

现在gorm将在SQL语句中将未设置的值解析为'0000-00-00 00:00:00.000000'，并导致以下错误。

> Error 2021-01-29 15:36:13,388 v1(7) error_logger.go:14 192.168.10.100 - - default - 0 Error 1292: Incorrect datetime value: '0000-00-00' for column 'last_login' at row 1

虽然我理解MySQL在不更改某些模式的情况下不允许时间戳值为零的原因，但我可以很容易地将time.Time字段初始化为一些较远的日期，例如2038年左右。如何告诉gorm将零时间字段作为NULL传递到SQML中？

## 回答
所以你有几个选择。您可以将LastLogin设为指针，这意味着它可以是一个nil值：
```
type User struct {
    ID        uint64 `gorm:"primaryKey"`
    CreatedAt time.Time
    LastLogin *time.Time
}
```
或者像@aureliar提到的那样，你可以使用sql.NullTime类型
```
type User struct {
    ID        uint64 `gorm:"primaryKey"`
    CreatedAt time.Time
    LastLogin sql.NullTime
}
```
现在，当您在数据库中创建该对象时，如果没有设置LastLogin，它将在数据库中保存为NULL。

https://gorm.io/docs/models.html

值得注意的是，如果您使用sql.NullTime，在结构中您将看到一个默认时间戳，而不是nil值

转自：https://cloud.tencent.com/developer/ask/sof/478567
