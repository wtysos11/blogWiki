# golang微信服务端编写

标签：golang 服务端 微信

## 参考资料

* [GoWechat：微信各大平台通用golang API](https://yaotian.github.io/gowechat/)
* [silenceper/wechat:微信SDK](https://github.com/silenceper/wechat)
* [medivhzhan/weapp:go微信小程序SDK](https://github.com/medivhzhan/weapp)
* [golang微信支付服务端，似乎没什么用](https://studygolang.com/articles/5636)
* [studygolang上的问题，推荐用beego](https://studygolang.com/topics/2386)
* [微信小程序go语言后端趟坑小笔记](http://teawater.github.io/docs/wechatgo.html)
* [微信小程序demo之猜拳游戏：go+websocket+redis+mysql](http://w ww.wxapp-union.com/thread-4873-1-1.html)
* [gorilla/websocket：websocket协议的实现](https://github.com/gorilla/websocket)

### 微信需要什么

[微信网络要求](https://developers.weixin.qq.com/miniprogram/dev/framework/ability/network.html)，只能跟指定域名进行网络通信，包括普通的HTTPS请求、上传文件、下载文件和WebSocket通信

普通https请求`wx.request`，[API位置](https://developers.weixin.qq.com/miniprogram/dev/api/wx.request.html)，可以配置向指定端口，但是不能更改。默认端口只能够向指定端口通信。

WebSocket通信`wx.connectSocket`，[API位置](https://developers.weixin.qq.com/miniprogram/dev/api/wx.connectSocket.html)，与https的十分类似。

### TLS证书相关

我用的是腾讯云的证书，不过如果是自己手写证书也许也不错。腾讯云证书申请是免费的，选择域名型免费版即可，[地址](https://cloud.tencent.com/product/ssl)。

* [Go代码打通https](https://segmentfault.com/a/1190000013287122)

### golang怎么实现https

* [golang的https服务器](https://studygolang.com/articles/378)
* [golang https/tls simple sample](https://github.com/denji/golang-tls)

### golang怎么实现wss

* [go websocket wss的设置方法](http://www.zrray.com/art/300)
* [go中使用websocket](http://www.lijiaocn.com/%E7%BC%96%E7%A8%8B/2017/11/03/golang-websocket.html)，对websocket协议的细节进行了一些讲解，并提供了demo的范本。

### golang使用beego

我研究了一下，推荐使用beego，原因有二：1.容易实现restful API，支持MVC模式。2.支持https，或者可以使用nginx作转发，但是也可以自己实现。3.beego的模块是高度解耦合的，即使将来不使用beego的http逻辑，也可以独立的使用这些模块。

学长推荐使用[postman](https://www.getpostman.com/)，一款很好地API开发工具，能够比较方便地测试API。

* [beego搭建api服务](https://studygolang.com/articles/7115)，这是go语言中文网的，例子很不错，但是不够清楚。
* [beego+swagger快速上手](https://juejin.im/post/5a90bec3f265da4e9957a282)，非常好的教程，很实用，在10分钟之内绝对可以完成一个简单的demo。作者前面还写过一篇[swagger上手](https://www.jianshu.com/p/06b7b752a983)，可以看一下。关键点在于[beego的API自动化文档](https://beego.me/docs/advantage/docs.md)