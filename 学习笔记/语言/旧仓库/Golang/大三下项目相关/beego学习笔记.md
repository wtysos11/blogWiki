# beego学习笔记

标签：beego 学习笔记

参考资料：

* [beego官方中文文档](https://beego.me/docs/intro/)

推荐使用[postman](https://www.getpostman.com/)，一款很好地API开发工具，能够比较方便地测试API。

* [beego搭建api服务](https://studygolang.com/articles/7115)，这是go语言中文网的，例子很不错，但是不够清楚。
* [beego+swagger快速上手](https://juejin.im/post/5a90bec3f265da4e9957a282)，非常好的教程，很实用，在10分钟之内绝对可以完成一个简单的demo。作者前面还写过一篇[swagger上手](https://www.jianshu.com/p/06b7b752a983)，可以看一下。关键点在于[beego的API自动化文档](https://beego.me/docs/advantage/docs.md)

## 快速入门

### 架构概述

beego是MVC结构，其中model主要负责数据库（逐层抽象），view层为模版，controller层作路由。

### 极速上手demo

要求：提前安装好了go（或者docker跑golang镜像）

安装工具：`go get github.com/astaxie/beego`和`go get github.com/beego/bee`

创建api项目：`bee api apiTest`

自动下载swagger文件，自动化文档，即可在本地浏览默认API：`http://localhost:8080/swagger/`。`bee run -gendoc=true -downdoc=true`

API目录

```bash
├── conf
│   └── app.conf
├── controllers
│   └── object.go
│   └── user.go
├── docs
│   └── doc.go
├── main.go
├── models
│   └── object.go
│   └── user.go
├── routers
│   └── router.go
└── tests
    └── default_test.go
```

剩下的就是阅读controllers、models文件下的go文件源代码

最终效果

![result](https://user-gold-cdn.xitu.io/2018/2/24/161c56909f0f5f56?imageView2/0/w/1280/h/960/format/webp/ignore-error/1)

让beego运行Https的配置

```conf
appname = apiproject
httpport = 8080
runmode = dev
autorender = false
copyrequestbody = true
EnableDocs = true
EnableHTTPS = true
HTTPSCertFile = "cert.crt"
HTTPSKeyFile = "key.key"
sqlconn = 
```

根目录为main.go所在目录，把CA证书换成自己的即可使用。具体的配置信息参考[这里](https://beego.me/docs/mvc/controller/config.md)

## 核心部分

### controller设计

#### 路由设置

beego存在三种方式的路由：固定路由、正则路由和自动路由。

* 固定路由：就是全匹配的路由，一个路由对应一个控制器，然后根据用户请求方法不同去请求控制器中的对应的方法。
* 正则路由：为了用户更方便的设置路由，beego参考了sinatra的路由实现，支持多种方式的路由。
* 自动路由：首先需要把路由的控制器注册到自动路由中，beego会通过反射获取该结构体中所有的实现方法，然后可以通过api对应自动调用方法(`/:controller/:method`匹配，剩余的url beego会自动化解析为参数，保存在`this.Ctx.Input.Params`中)。
* 注解路由：用户无需在router中注册路由，只需要include相应的controller，然后在controller的method方法上写上router注释(//@router)就可以了。beego会自动进行源码分析，注意，只会在dev模式下进行生成，生成的路由会放在`/routers/commentsRouter.go`文件中。

##### namespace

通过namespace可以支持更多层次的api，可以前置过滤、条件判断、无限嵌套namespace。

### 控制器函数

基于beego的Controller设计可以参考controller文件夹下的文件，只要匿名组合beego.Controller就可以了，如下。而beego.Controller实现了接口beego.ControllerInterface以及一系列的函数。

```golang
type xxxController struct {
    beego.Controller
}
```

提前终止运行：比如用户登陆逻辑没有通过，就直接终止进程。使用`StopRun`可以做到。

```golang
type RController struct {
    beego.Controller
}

func (this *RController) Prepare() {
    this.Data["json"] = map[string]interface{}{"name": "astaxie"}
    this.ServeJSON()
    this.StopRun()
}
```

### model设计

安装ORM：`go get github.com/astaxie/beego/orm`

简单例子 models.go：

```go
package main

import (
    "github.com/astaxie/beego/orm"
)

type User struct {
    Id          int
    Name        string
    Profile     *Profile   `orm:"rel(one)"` // OneToOne relation
    Post        []*Post `orm:"reverse(many)"` // 设置一对多的反向关系
}

type Profile struct {
    Id          int
    Age         int16
    User        *User   `orm:"reverse(one)"` // 设置一对一反向关系(可选)
}

type Post struct {
    Id    int
    Title string
    User  *User  `orm:"rel(fk)"`    //设置一对多关系
    Tags  []*Tag `orm:"rel(m2m)"`
}

type Tag struct {
    Id    int
    Name  string
    Posts []*Post `orm:"reverse(many)"`
}

func init() {
    // 需要在init中注册定义的model
    orm.RegisterModel(new(User), new(Post), new(Profile), new(Tag))
}
```

main.go

```go
package main

import (
    "fmt"
    "github.com/astaxie/beego/orm"
    _ "github.com/go-sql-driver/mysql"
)

func init() {
    orm.RegisterDriver("mysql", orm.DRMySQL)

    orm.RegisterDataBase("default", "mysql", "root:root@/orm_test?charset=utf8")
}

func main() {
    o := orm.NewOrm()
    o.Using("default") // 默认使用 default，你可以指定为其他数据库

    profile := new(Profile)
    profile.Age = 30

    user := new(User)
    user.Profile = profile
    user.Name = "slene"

    fmt.Println(o.Insert(profile))
    fmt.Println(o.Insert(user))
}
```

#### 注册模型

如果使用`orm.QuerySeter`进行高级查询的话，这个是必须的。反之，如果只使用Raw查询和map struct，是无需这一步的。

首先，注册定义的Model，最佳的设计是有单独的models.go文件，在它的init函数内完成注册，比如：

迷你的models.go

```go
package main

import "github.com/astaxie/beego/orm"

type User struct {
    Id   int
    Name string
}

func init(){
    orm.RegisterModel(new(User))
}
```

可以如前文一样同时注册多个model。可以使用`RegisterModelWithPrefix`，这样创建的表名会自动带上前缀。

#### 高级查询

ORM以QuerySeter来组织查询，每个返回QuerySeter的方法都会获得一个新的QuerySeter对象。

QuerySeter使用`Filter`组织查询，字段组合的前后顺序依照表的关系。字段的分隔符使用`__`双下划线。除了描述字段，expr的尾部可以增加操作符以执行对应的sql操作，比如`Profile__Age__gt`代表`Profile.Age>18`的条件查询。

## mysql中shell常用命令

因为使用的是mysql，而且经常要用到一些命令，所以在这里记录一下

登陆`mysql -u root -p`，或者在mysql shell中输入`\connect [scheme://][user[:password]@]<host[:port]|socket>[/schema][?option=value&option=value...]`，例如`\connect mysql://wty:97112500@127.0.0.1:3306`来登陆

使用`\sql`可以切换到sql模式，这样子就可以使用常用的sql命令来操作。

## 遇到的问题

### beego中orm怎么使用

在哪里定义，是在main里定义后面使用呢，还是在