## Go HTTP 服务超时控制
系统对外提供 HTTP 服务的时候一定要控制好各个环节的超时时间，不然很容易受到 DDos 攻击。我们部门使用的业务框架是基于 Go 语言的 net/http 标准库二次开发的。在当年开发框架的时候，
我对 Go 语言 HTTP 服务器的超时控制理解并不深刻。当时觉着只要在最外层加一个 http.TimeoutHandler 就足够了。系统上线后也一直没有出这方面的问题，还自我感觉良好。其实是因为我们运维
在最外层的 nginx 设置了各项超时控制，没有把系统的问题暴露出来。等我们在 AWS 上运行另一套业务系统的时候，因为 AWS 的 ALB 配置跟原来的 nginx 不同，我们发现只用 
http.TimeoutHandler 居然在特殊场景中会产生「死锁」！我当场就阵亡了，赶紧排查原因。大致看了一遍 Go 语言 HTTP 服务的源码，找到了死锁的原因。今天就把相关经验分享给大家。我在看代码
的时候发现 Go HTTP 服务器在读到完整的 HTTP 请求后会再起一个协程，该协程会试着再读一个字节的内容。非常奇怪🤔，也一并研究了一下，最终找到了相关的 issue 和提交记录，并发现了 Go 
语言的一个缺陷，今天也一并分享给大家。
#### Go HTTP 服务器对应的结构是 net/http.Server，跟超时相关的配置有四个：
* ReadTimeout
* ReadHeaderTimeout
* WriteTimeout
* IdleTimeout

除了这四个配置外，还可以使用 TimeoutHandler，但这个需要调用net/http.TimeoutHandler() 生成，函数签名如下：
```
package http
func TimeoutHandler(h Handler, dt time.Duration, msg string) Handler
```
以上配置和 TimeHandler 对应的作用过程如下图（来自 Cloudflare）：


转自：https://taoshu.in/go/go-http-server-timeout.html
