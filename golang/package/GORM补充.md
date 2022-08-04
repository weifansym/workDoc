### GORM补充
### DBResolver
DBResolver 为 GORM 提供了多个数据库支持，支持以下功能：
* 支持多个 sources、replicas
* 读写分离
* 根据工作表、struct 自动切换连接
* 手动切换连接
* Sources/Replicas 负载均衡
* 适用于原生 SQL
* Transaction

https://github.com/go-gorm/dbresolver

转自：https://kuriyama.net.cn/blog/212
