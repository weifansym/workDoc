## docker部署简单的web项目
今天就简单的使用docekr来部署一个很简单web应用，这个web营部不需要和其他应用交互。具体过程如下：

### 创建node项目
```
mkdir nodedemo
cd nodedemo/
npm init -y
```
构建node项目需要的package.json，然后创建一个文件
```
touch index.js
```
在创建的index.js中创建一个最简单的http服务
```
const http = require('http')
const port = 8888

const requestHandler = (request, response) => {
  console.log(request.url)
  response.end('Hello Node.js Server!')
}

const server = http.createServer(requestHandler)

server.listen(port, (err) => {
  if (err) {
    return console.log('something bad happened', err)
  }

  console.log(`server is listening on ${port}`)
})
```
修改package.json在 script中添加下面语句 ，便于通过npm启动项目
```
"start": "node index.js"
```
这样一个简单的node项目就创建完成了。

### Dockerfile文件
首先，在node工程的根目录创建Dockerfile文件，该文件是node工程中对docker的配置文件。
* 创建Dcokerfile文件
```
vi Dockerfile
```
* 输入如下内容：
```
#node镜像版本
FROM node:10-alpine
#声明作者
MAINTAINER LI
#在image中创建文件夹
RUN mkdir -p /home/node-app
#将该文件夹作为工作目录
WORKDIR /home/node-app

# 将node工程下所有文件拷贝到Image下的文件夹中
COPY . /home/node-app

#使用RUN命令执行npm install安装工程依赖库
RUN npm install

#暴露给主机的端口号
EXPOSE 8888
#执行npm start命令，启动Node工程
CMD [ "npm", "start" ]
```
### 构建image

执行命令docker build -t node-app:v1 . 需要注意v1后面还有一个.

其中 -t node-app:v1 为构建的镜像名称及标签
```
weifandeMacBook-Pro:nodedemo weifan$ docker build -t node-app:v1 .
Sending build context to Docker daemon  4.096kB
Step 1/8 : FROM node:10-alpine
10-alpine: Pulling from library/node
89d9c30c1d48: Pull complete
0eaf5bd7a6e1: Pull complete
6fb5c3a20092: Pull complete
004e30fa1cb9: Pull complete
Digest: sha256:da8161962573bd6ab16b54a9bfa81a263458e5199074d0678d0556376b22bd22
Status: Downloaded newer image for node:10-alpine
 ---> a0708430821e
Step 2/8 : MAINTAINER LI
 ---> Running in aabbd55ed517
Removing intermediate container aabbd55ed517
 ---> a10f70f38593
Step 3/8 : RUN mkdir -p /home/node-app
 ---> Running in 7a6b8e161391
Removing intermediate container 7a6b8e161391
 ---> 920e3bc9c472
Step 4/8 : WORKDIR /home/node-app
 ---> Running in 6e2f5757f693
Removing intermediate container 6e2f5757f693
 ---> a96bbf9470e3
Step 5/8 : COPY . /home/node-app
 ---> 42649c246eaf
Step 6/8 : RUN npm install
 ---> Running in 06557e5d84b2
npm notice created a lockfile as package-lock.json. You should commit this file.
npm WARN nodedemo@1.0.0 No description
npm WARN nodedemo@1.0.0 No repository field.

up to date in 0.77s
found 0 vulnerabilities

Removing intermediate container 06557e5d84b2
 ---> 841ec25f4818
Step 7/8 : EXPOSE 8888
 ---> Running in 7099d3bbc129
Removing intermediate container 7099d3bbc129
 ---> c1a0c4d7ebbc
Step 8/8 : CMD [ "npm", "start" ]
 ---> Running in 05ad9d5f9fea
Removing intermediate container 05ad9d5f9fea
 ---> 0ed562689cbc
Successfully built 0ed562689cbc
Successfully tagged node-app:v1
weifandeMacBook-Pro:nodedemo weifan$
```
查看生成的image: docker images命令
```
weifandeMacBook-Pro:nodedemo weifan$ docker images
REPOSITORY          TAG                 IMAGE ID            CREATED             SIZE
node-app            v1                  0ed562689cbc        5 minutes ago       75.4MB
nginx               v3                  bc0ced96c716        5 hours ago         126MB
redis               5.0.7               dcf9ec9265e0        2 weeks ago         98.2MB
nginx               latest              231d40e811cd        2 weeks ago         126MB
node                10-alpine           a0708430821e        3 weeks ago         75.4MB
```
### 运行container
> 执行命令 docker run -d -p 8888:8888 0ed5
其中， -d表示在容器后台运行，-p表示端口映射，将本机的8888端口映射到container的8888端口，外网访问本机的8888端口即可访问container。0ed5为生成的IMAGE的ID,只需要写入对应ID的前几位系统能辨识出对应的image即可。
```
weifandeMacBook-Pro:nodedemo weifan$ docker run -d -p 8888:8888 0ed562689cbc
b3d374e5f2c92d581c8151a2cb963b1a3c66e1f52c3afd750cc8f5903897018b
```
> 执行命令docker ps查看container是否运行
```
weifandeMacBook-Pro:nodedemo weifan$ docker ps
CONTAINER ID        IMAGE               COMMAND                  CREATED             STATUS              PORTS      NAMES
b3d374e5f2c9        0ed562689cbc        "docker-entrypoint.s…"   9 seconds ago       Up 8 seconds     0.0.0.0:8888->8888/tcp   upbeat_bhabha
```
> 通过命令docker logs b3d374e5f2c9 还可查看container的日志
```
weifandeMacBook-Pro:nodedemo weifan$ docker logs b3d374e5f2c9

> nodedemo@1.0.0 start /home/node-app
> node index.js

server is listening on 3000
```
此时服务器已经正常启动了。
在浏览器中访问：http://localhost:8888/，你会看到Hello Node.js Server!
### 进入容器
为了方便查看容器内部文件和调试，可以通过命令进入容器中。容器内部就像一个小型的linux系统一样。命令为docker exec -it b3d374e5f2c9 /bin/sh
```
weifandeMacBook-Pro:nodedemo weifan$ docker exec -it b3d374e5f2c9 /bin/sh
/home/node-app # ls
Dockerfile         index.js           package-lock.json  package.json
```
### 日志
* docker镜像中node工程会有打印日志功能，因为docker容器一旦挂掉，容器中的文件也会访问不了，所以日志必须要放在docker镜像外的文件路径下。此时，必须要将centos系统中的日志文件目录挂在到docker容器中，在容器启动时开启数据卷，实现日志采集。
* 在启动容器时，使用命令docker run -d -p 8888:8888 -v /home/logs:/data/logs 190f即可。/home/logs为centos系统中日志文件目录，data/logs为docker容器中node工程写入日志路径。
* 如果docker容器中工程需要写入文件，则在启动时要加上--privileged=true才可以。

### 打包与解压
> 如果没有私有仓库，则可以通过save和load命令来打包和解压。这样方便我们备份镜像或者给被人使用， save将docker镜像压缩为tar文件，load为将tar文件解压生成镜像。
1. 打包镜像
```
weifandeMacBook-Pro:nodedemo weifan$ docker save 0ed562689cbc -o node-app-v1.tar
weifandeMacBook-Pro:nodedemo weifan$ ls
Dockerfile	index.js	node-app-v1.tar	package.json
weifandeMacBook-Pro:nodedemo weifan$
```
2. 解压载入镜像
```
weifandeMacBook-Pro:nodedemo weifan$ docker load < node-app-v1.tar
Loaded image ID: sha256:0ed562689cbceb490fbe264e908c045ce737381f7d9b0d70e63b285344af3099
```

参考：
* https://juejin.im/post/5b82613f6fb9a019ce1490fe
* https://www.cnblogs.com/linjiqin/p/8604756.html
* https://github.com/nodejs/docker-node
