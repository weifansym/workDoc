## Module ngx_http_core_module
[ngx_http_core_module](http://nginx.org/en/docs/http/ngx_http_core_module.html#root)

### alias
```
Syntax:	alias path;
Default:	—
Context:	location
```
定义一个指定位置的替代项，例如下面的配置
```
location /i/ {
    alias /data/w3/images/;
}
```
当请求**“/i/top.gif”**的时候，文件**/data/w3/images/top.gif**将会被发送。

**path**的值可以包含变量，除了：$document_root 和 $realpath_root之外。

如果**alias**在location中使用正则表达式定义，正则表达式应该包含捕获，alias应该指向这个捕获，例如：
```
location ~ ^/users/(.+\.(?:gif|jpe?g|png))$ {
    alias /data/w3/images/$1;
}
```
当location匹配指令值最后的部分
```
location /images/ {
    alias /data/w3/images/;
}
```
此时应该最好使用root指令值进行替代
```
location /images/ {
    root /data/w3;
}
```

### root
```
Syntax:	root path;
Default:	
root html;
Context:	http, server, location, if in location
```
为请求设置根目录，例如下面的配置：
```
location /i/ {
    root /data/w3;
}
```
当请求**/i/top.gif**的时候/data/w3/i/top.gif文件将会被发送。

**path**的值可以包含变量，除了：$document_root 和 $realpath_root之外。
一个文件的路径被创建通过仅仅在root指令中添加url，如果url必须修改则要使用alias指令


### alias与root的区别
```
root    实际访问文件路径会拼接URL中的路径
alias   实际访问文件路径不会拼接URL中的路径
```
示例如下：
```
location ^~ /sta/ {  
   alias /usr/local/nginx/html/static/;  
}
```
请求：http://test.com/sta/sta1.html
实际访问：/usr/local/nginx/html/static/sta1.html 文件

```
location ^~ /tea/ {  
   root /usr/local/nginx/html/;  
}
```
请求：http://test.com/tea/tea1.html
实际访问：/usr/local/nginx/html/tea/tea1.html 文件


