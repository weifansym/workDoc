## Linux常用操作
### 1、运行.sh文件
第一种方法：

首先你要打开一个终端。
然后输入sudo su
随后输入密码。这样就取得了root用户权限。
然后找到那个文件
首先看这个文件是不是有执行权限，没执行权限的话就要更改文件的执行权限然后在进行余下的操作
执行./sh文件名字
这样.sh就运行了。
第二种方法：
```
sh xx.sh
```

### 2、查看文件状态，例如修改时间，创建时间，文件的大小等
例如我有一个文件名为test.log，使用stat来查看文件的状态
```
stat test.log
```
### 3、查看某个文件夹的大小
du -h --max-depth=1


### 4、Linux查找某个文件夹下是否包含某个字符串
grep -rn "6402105992922202358" *         //  执行的命令 6402105992922202358 是要查找的字符串
例如：



### 5、使用SSH链接远程ip
Windows下我习惯用Xshell来ssh登录，Mac直接使用Terminal即可。

ssh的一些常用命令：

ssh root@ip
使用root账号登录指定ip的服务器。下面需要把ip换成你自己服务器的ip。如果服务器使用的不是标准端口，比如是2345端口，则是：

ssh root@ip -p 2345
退出当前登录的服务器：

exit
### 6、SCP命令的使用
下面是一些简单例子：

copy本地文件到服务器的命令如下：

scp <local file> <remote user>@<remote machine>:<remote path>
 

上传文件：

[root@test test]# scp ./mytest/password.php 172.30.4.42:/tmp/test2
将当前目录中的mytest目录下的password.php上传到172.30.4.42服务器/tmp/test2目录下面。

 

上传目录：

[root@test test]# scp -r ./mytest 172.30.4.42:/tmp/test2
将当前目录中的mytest目录上传到172.30.4.42服务器/tmp/test2目录下面。

 

如果想Copy远程文件到本地，则是：

scp <remote user>@<remote machine>:<remote path> <local file>
下载文件

[root@test test]# scp 172.30.4.42:/tmp/test2/aaa.php ./
将172.30.4.42linux系统中/tmp/test2/aaa.php文件copy到当前目录下面

下载目录

[root@test test]# scp -r root@172.30.4.42:/tmp/test2 ./
将172.30.4.42linux系统中/tmp/test2目录copy到当前目录下面，在这172.30.4.42前面加了root@,提示输入密码，如果不加呢，会提示你输入用户名和密码

具体查看scp命令 ：scp --help

### 7、zip压缩与unzip解压
把某个目录压缩，命令如下：

zip -r  file.zip FolderName
其中-r 表示对文件夹进行压缩（即循环处理文件），file.zip表示要压缩后生成的文件名，FolderName表示要压缩的目录或文件夹名

例如：

zip -r Projects.zip Projects/
如果是压缩某个文件，去掉-r参数即可。

解压缩：

unzip file.zip
例如：

unzip Projects.zip
### 8、服务器与本地文件传输
对于经常使用Linux系统的人员来说，少不了将本地的文件上传到服务器或者从服务器上下载文件到本地，rz / sz命令很方便的帮我们实现了这个功能，但是很多Linux系统初始并没有这两个命令。今天，我们就简单的讲解一下如何安装和使用rz、sz命令。yum安装：
yum install -y lrzsz 
 使用如下：  sz命令发送文件到本地： 
# sz filename
rz命令本地上传文件到服务器：
# rz
执行该命令后，在弹出框中选择要上传的文件即可。
说明：打开SecureCRT软件 -> Options -> session options -> X/Y/Zmodem 下可以设置上传和下载的目录。
 
### 9、创建用户与切换用户
首先用adduser命令添加一个普通用户，命令如下：
```
#adduser tommy  //添加一个名为tommy的用户
#passwd tommy   //修改密码
Changing password for user tommy.
New UNIX password:     //在这里输入新密码
Retype new UNIX password:  //再次输入新密码
passwd: all authentication tokens updated successfully.
```
删除用户
```
userdel test
```
将test用户删除

切换用户：
```
可以使用su命令来切换用户，su是switch user切换用户的缩写。可以是从普通用户切换到root用户，也可以是从root用户切换到普通用户。从普通用户切换到root用户需要输入密码，从root用户切换到普通用户不需要输入密码。
命令格式：su [参数] [-] [用户名]
用户名的默认值为root。
用法示例：
su zhidao #切换到zhidao用户
su #切换到root用户

```


 

 

 
 
