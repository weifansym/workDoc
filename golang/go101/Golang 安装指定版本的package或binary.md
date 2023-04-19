## Golang 安装指定版本的package或binary
Golang要安装指定版本的package或可执行档binary的方式如下。
### 安装install
在命令列输入go get <package_path>@<versoin>即可下载指定版本(tag)的package，<package_path>为要下载的package资源路径，<version>为指定版本。

例如下载swag套件版本v1.8.1则输入go get -u github.com/swaggo/swag@v1.8.1。
```
$ go get -u github.com/swaggo/swag@v1.8.1
go: downloading golang.org/x/tools v0.3.0
go: downloading golang.org/x/sync v0.1.0
go: downloading golang.org/x/net v0.2.0
go: downloading golang.org/x/mod v0.7.0

```  
在命令列输入go install <binary_path>@<versoin>即可下载指定版本(tag)的binary执行档。<binary_path>为要下载的binary资源路径。

例如下载swag的cmd/swag binary版本v1.8.3则输入go install github.com/swaggo/swag/cmd/swag@v1.8.3。
```
$ go install github.com/swaggo/swag/cmd/swag@v1.8.3
go: downloading golang.org/x/tools v0.1.10
go: downloading golang.org/x/net v0.0.0-20220425223048-2871e0cb64e4
go: downloading golang.org/x/sys v0.0.0-20211216021012-1d35b9e2eb4e
```
若要删除安装的package或bin参考「Golang 删除安装的package或binary档 」。
