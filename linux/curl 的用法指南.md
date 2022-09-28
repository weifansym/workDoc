## curl 的用法指南
curl 是常用的命令行工具，用来请求 Web 服务器。它的名字就是客户端（client）的 URL 工具的意思。

它的功能非常强大，命令行参数多达几十种。如果熟练的话，完全可以取代 Postman 这一类的图形界面工具。

本文介绍它的主要命令行参数，作为日常的参考，方便查阅。内容主要翻译自[《curl cookbook》](https://catonmat.net/cookbooks/curl)。为了节约篇幅，下面的例子不包括运行时的输出，初学者可以先看我以前写的[《curl 初学者教程》](https://www.ruanyifeng.com/blog/2011/09/curl.html)。

不带有任何参数时，curl 就是发出 GET 请求。
```
$ curl https://www.example.com
```
上面命令向www.example.com发出 GET 请求，服务器返回的内容会在命令行输出。

#### -A
-A参数指定客户端的用户代理标头，即User-Agent。curl 的默认用户代理字符串是curl/[version]。
```
$ curl -A 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/76.0.3809.100 Safari/537.36' https://google.com
```
上面命令将User-Agent改成 Chrome 浏览器。
```
$ curl -A '' https://google.com
```
上面命令会移除User-Agent标头。

也可以通过-H参数直接指定标头，更改User-Agent。


转自：https://www.ruanyifeng.com/blog/2019/09/curl-reference.html
