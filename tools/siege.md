## Mac安装压测工具siege
### 修改Mac文件描述符限制
在压测开始前，你需要确保你的open files足够大，否则会报TOO MANY FILES OPEN错误，可以通过ulimit -a查看
```
$ ulimit -a
core file size          (blocks, -c) 0
data seg size           (kbytes, -d) unlimited
file size               (blocks, -f) unlimited
max locked memory       (kbytes, -l) unlimited
max memory size         (kbytes, -m) unlimited
open files                      (-n) 256
pipe size            (512 bytes, -p) 1
stack size              (kbytes, -s) 8192
cpu time               (seconds, -t) unlimited
max user processes              (-u) 709
virtual memory          (kbytes, -v) unlimited
```
使用ulimit -n 10000可以修改该值。不过这种修改并不是永久的，关闭终端会话，又会恢复回来。

### 在Mac上安装siege
这里有两种方式
#### 第一种
```
wget http://download.joedog.org/siege/siege-latest.tar.gz
tar -xvf siege-latest.tar.gz
cd siege-4.0.2/
./configure 
make
make install
```
默认安装在/usr/local/bin/，并自动添加到系统环境变量中，在终端输入siege看查看命令命令。下面介绍一些常用命令。
#### 第二种
mac 安装siege很简单，brew install siege 即可。
### siege压测命令
* -C, --config  在屏幕上打印显示出当前的配置,配置是包括在他的配置文件$HOME/.siegerc中,可以编辑里面的参数,这样每次siege 都会按照它运行.
* -v, --verbose 运行时能看到详细的运行信息.
* -c, --concurrent=NUM 模拟有n个用户在同时访问,n不要设得太大,因为越大,siege消耗本地机器的资源越多.
* -r, --reps=NUM 重复运行测试n次,不能与-t同时存在
* -t, --time=NUMm 持续运行siege ‘n’秒(如10S),分钟(10M),小时(10H)
* -d, --delay=NUM 每个url之间的延迟,在0-n之间.
* -b, --benchmark 请求无需等待 delay=0.
* -i, --internet  随机访问urls.txt中的url列表项.
* -f, --file=FILE 指定用特定的urls文件运行 ,默认为urls.txt,位于siege安装目录下的etc/urls.txt
* -R, --rc=FILE   指定用特定的siege 配置文件来运行,默认的为$HOME/.siegerc
* -l, --log[=FILE] 运行结束,将统计数据保存到日志文件中siege .log,一般位于/usr/local/var/siege .log中,也可在.siegerc中自定义
### siege压测结果
```
//并发10个,发生5次,共50个请求
siege -c 10 -r 5 http://www.baidu.com

Transactions: 350 hits //总共测试次数
Availability: 100.00 % //成功次数百分比
Elapsed time: 4.27 secs //总共耗时多少秒
Data transferred: 7.08 MB //总共数据传输
Response time: 0.07 secs //等到响应耗时
Transaction rate: 81.97 trans/sec //平均每秒处理请求数
Throughput: 1.66 MB/sec //吞吐率
Concurrency: 6.06 //最高并发
Successful transactions: 350 //成功的请求数
Failed transactions: 0 //失败的请求数
Longest transaction: 0.24 //每次传输所花最长时间
Shortest transaction: 0.01 //每次传输所花最短时间
```
### 常用命令
```
# 200个并发对http://www.google.com发送请求100次
siege -c 200 -r 100 http://www.google.com

# 在urls.txt中列出所有的网址
siege -c 200 -r 100 -f urls.txt

# 随机选取urls.txt中列出所有的网址
siege -c 200 -r 100 -f urls.txt -i

# delay=0，更准确的压力测试，而不是功能测试
siege -c 200 -r 100 -f urls.txt -i -b

# 指定http请求头 文档类型
siege -H "Content-Type:application/json" -c 200 -r 100 -f urls.txt -i -b
```
### 注意事项
发送post请求时，url格式为 http://www.xxxx.com/ POST p1=v1&p2=v22
如果url中含有空格和中文，要先进行url编码，否则siege发送的请求url不准确。

参考：https://my.oschina.net/OSrainn/blog/1628722

