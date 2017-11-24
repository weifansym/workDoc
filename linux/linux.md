##  linux操作
## curl发送POST请求
curl发送POST请求 今天写Gitlab的一个merge request hook,使用curl来简化测试请求.
简单备忘一下,如何使用curl发送POST请求.以下为使用curl发送一个携带json数据的POST请求.

命令介绍：
```
-H, —header LINE Custom header to pass to server (H)
-d, —data DATA HTTP POST data (H)
```
实例命令：
```
curl -H "Content-Type:application/json" -X POST --data '{"userName": "test11"}' http://127.0.0.1:9001/v1/message/leaveMessage
```
