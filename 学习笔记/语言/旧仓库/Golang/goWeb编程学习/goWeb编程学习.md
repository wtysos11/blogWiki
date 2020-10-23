# goWeb编程学习

标签：golang web 学习笔记

## 来源

[gitBook](https://github.com/astaxie/build-web-application-with-golang/tree/master/zh)
很不错的一本书，介绍了基础、http解析，表单处理，数据库，session和数据存储

[目录](https://github.com/astaxie/build-web-application-with-golang/blob/master/zh/preface.md)

事实上，原书已经足够的简介了，为了不将原书再抄一遍，我只会列出我觉得重点的部分和章节概要。

## 实验安排

* 实验1：http表单解析实验。目的是了解表单解析的过程
* 实验2：http表单解析与模版使用。了解如何使用模板，以及表单解析`r.parseForm()`的必要性。
* 实验3：http表单提交文件

## 第二章 基础


基础部分跳过，没有看基础的可以看一下，写的还是很不错的。

[2.5 method](https://github.com/astaxie/build-web-application-with-golang/blob/master/zh/02.5.md)讲的挺好的，提到了golang中如何实现类似于面向对象的继承与方法重写（虽然我的老师和许多网上的资料都说GO本身不是OO的，但是OO的很多特性还是十分方便且符合人类习惯的）

[2.6 interface](https://github.com/astaxie/build-web-application-with-golang/blob/master/zh/02.6.md)
类型断言还是很重要的。

[Go官方的反射教学](https://blog.golang.org/laws-of-reflection)讲的很好

## 第三章 Web基础

[Go Http包执行流程](https://github.com/astaxie/build-web-application-with-golang/blob/master/zh/03.4.md)对于阅读HTTP包的源很有帮助。

## 第四章 表单

[向服务器上传文件](https://github.com/astaxie/build-web-application-with-golang/blob/master/zh/04.5.md)，没有怎么实践过这东西，找天仔细看看。

## 第五章 数据库

例子很多，挺不错的。

## 第六章 session和cookie

[go中如何使用session](https://github.com/astaxie/build-web-application-with-golang/blob/master/zh/06.2.md)

## 第七章 文本处理

JSON、XML以及正则表达式，字符串的处理

## 第八章 WEB

挺重要的一章

[RPC](https://github.com/astaxie/build-web-application-with-golang/blob/master/zh/08.4.md)，有所耳闻，但确实没有在go上做过RPC。值得一提的是，golang内部提供了以gob编码的RPC只能在以go语言作为服务端和客户端的机器上相互交流。不过也提供了以json编码的RPC，具体还需要进一步的实验。

## 第九章 安全

主要是几种基本的网络攻击方式的预防，go中加密与解密方法的学习。

## 第十章 国际化与本地化

国际化与本地化是十分重要的，这里介绍了如何实现类似的包，但是实际上一般直接使用第三方包也足够用了。

## 第十一章 错误处理，调试和测试

非常重要，介绍了错误类型，gdb调试以及测试样例的书写。

## 第十二章 部署与维护

