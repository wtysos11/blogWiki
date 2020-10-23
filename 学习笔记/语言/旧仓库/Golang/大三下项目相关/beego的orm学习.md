# beego中的orm学习

标签：beego orm

因为我没怎么用过orm，一般就直接用着sql驱动直接上了。数据库上的也有一段时间了，对于各种关系忘得也是差不多了。

这里主要是官方写的不是特别好，比如官方中文文档在[这里](https://beego.me/docs/mvc/model/orm.md)给出的models.go的例子，让人看的很迷茫。

```golang
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

这……我知道那些标签能够用反射获得，可是orm后面的那些什么一对一、一对多关系都是什么东西？下面将记录一些我在orm学习的一些坑点，以及数据库知识与orm知识的回顾与结合。

## 使用orm进行原始的SQL操作

主要包括三个：直接用Raw函数，使用Exec执行SQL指令与事务提交。一般而言orm就不要直接用SQL了，不然和直接使用SQL驱动有什么区别。

### QueryTable

传入表名，或model对象，返回一个[QuerySeter](https://beego.me/docs/mvc/model/query.md)。ORM以QuerySeter来组织查询，每个返回QuerySeter的方法都会获得一个新的QuerySeter对象。

```golang
o := orm.NewOrm()

// 获取 QuerySeter 对象，user 为表名
qs := o.QueryTable("user")

// 也可以直接使用对象作为表名
user := new(User)
qs = o.QueryTable(user) // 返回 QuerySeter
```

这里是高级查询的具体实现部分。

#### 接口实现

QuerySeter可以使用Filter方法，接受参数并进行数组上的判断。多个Filter之间使用And进行连接。

### Raw

使用sql语句直接进行操作。使用Raw函数，会直接返回一个[RawSeter](https://beego.me/docs/mvc/model/rawsql.md)用以对设置的sql语句和参数进行操作。

```golang
o := orm.NewOrm()
var r RawSeter
r = o.Raw("UPDATE user SET name = ? WHERE name = ?", "testing", "slene")
```

可以创建一个新的RawSeter对象。

## 设置参数

终于找到了，在官方文档的[模型定义](https://beego.me/docs/mvc/model/models.md)这里。模型定义功能主要用于数据库数据转换和自动建表，对于我这种比较喜欢把所有相关的语句放在一个地方，同时又不想写SQL语句的人而言，这个功能还是值得深入研究的。（毕竟能够在本语言内解决的东西最好还是不要再加技术栈）

设置参数：`orm:"null;rel(fk)"`，多个设置之间使用`;`进行分隔。设置的值如果是多个，使用`,`进行分隔。

### 忽略字段

设置`-`可忽略struct中的字段，比如

```golang
type User struct {
...
    AnyField string `orm:"-"`
...
}
```

### auto

当Field类型为int等类型的时候，可以设置字段为自增键(extra:auto_increment)。当模型定义中没有主键时，符合上述类型且名称为`Id`的Field将被视为自增键。

### pk

primary key，设置为逐渐。适用于定义其他名称与类型为主键。

### null

数据库表默认为`NOT NULL`，设置null代表允许，即`ALLOW NULL`，比如`Name string orm:"null"`

### index

为单个字段增加索引（模型定义的开头有提到复合索引）

### unique

为单个字段增加unique键

### column

为字段设置db字段的名称（不设置将会自动转换，转换规则也是在模型定义页面的开头）。

        Name string `orm:"column(user_name)"`

### size

string类型默认为varchar(255)，设置size之后会变为varchar(size)

        Title string `orm:"size(60)"`

### digits/decimals

设置float32,float64类型的浮点精度

        Money float64 `orm:"digits(12);decimals(4)"`

即总长度12，小数点后4位

### auto_now/auto_now_add

        Created time.Time `orm:"auto_now_add;type(datetime)"`
        Updated time.Time `orm:"auto_now;type(datetime)"`

* auto_now每次model保存时都会对时间自动更新
* auto_now_add第一次保存时才设置时间

对于批量的update此设置是不生效的。

### type

设置为date时，time.Time字段的对应db类型使用date。设置为datetime时，time.Time字段的对应db类型使用datetime

### default

为字段设置默认值，类型必须符合（目前仅用于级联删除时的默认值）

### Comment

        type User struct {
            ...
            Status int `orm:"default(1)" description:(这是状态字段)`
            ...
        }

为字段添加注释，注释中禁止包含引号

## 表关系设置

rel/reverse

关系参考了[Django ORM 一对一、一对多、多对多详解](https://www.cnblogs.com/pythonxiaohu/p/5814247.html)，例子使用的是官方中文教程的例子。

### 一对一关系

子表从母表中选出一条数据一一对应，母表中选出来一条就少一条，子表中不可以再选择母表中已经选择的数据。

从数学的角度上来说，就是将子表作为定义域，母表作为值域的单射（不确定是否是满射，好像可以不是的样子）

RelOneToOne:

```golang
type User struct {
    ...
    Profile *Profile `orm:"null;rel(one);on_delete(set_null)"`
    ...
}
```

对应的反向关系RelReverseOne

```golang
type Profile struct {
    ...
    User *User `orm:"reverse(one)"`
    ...
}
```

实验打出来的表格信息：

```mysql
create table `user`
    -- --------------------------------------------------
    --  Table Structure for `main.User`
    -- --------------------------------------------------
    CREATE TABLE IF NOT EXISTS `user` (
        `id` integer AUTO_INCREMENT NOT NULL PRIMARY KEY,
        `profile_id` integer UNIQUE
    ) ENGINE=InnoDB;

create table `profile`
    -- --------------------------------------------------
    --  Table Structure for `main.Profile`
    -- --------------------------------------------------
    CREATE TABLE IF NOT EXISTS `profile` (
        `id` integer AUTO_INCREMENT NOT NULL PRIMARY KEY
    ) ENGINE=InnoDB;
```

可以看到实现上是作为Unique的对象，并没有实现类似外键之类的东西。实际在添加的时候也是并没有说我加了User，然后Profile就不用再添加了。实际上还是要一起添加的，不会去检查存在性与否，在数据库上的实现真的就是多了一个一对一的id作为Unique的值而已，并没有什么其他的区别。

感觉还是可以酌情添加吧，场景上上面那篇文章说是可以用于某张表的补充，比如用户信息与用户账号密码的分离，我个人的想法是可以将不必要的信息放在其他表里面减少查询所需要的消耗。感觉就是将一张原始的大表分为子表和母表时用到的。

### 外键关系，一对多关系

从子表从母表中选出一条数据一一对应，但母表的这条数据还可以被其他子表数据选择。从数据库实现上这个是真的有外键了。

RelForeignKey:

```golang
type Post struct {
    ...
    User *User `orm:"rel(fk)"` // RelForeignKey relation
    ...
}
```

对应的反向关系RelReverseMany:

```golang
type User struct {
    ...
    Posts []*Post `orm:"reverse(many)"` // fk 的反向关系
    ...
}
```

上面的Post中，每个Post对应一个User，但是User可以对应多个Post。数据库的实现如下：

```mysql
create table `user`
    -- --------------------------------------------------
    --  Table Structure for `main.User`
    -- --------------------------------------------------
    CREATE TABLE IF NOT EXISTS `user` (
        `id` integer AUTO_INCREMENT NOT NULL PRIMARY KEY
    ) ENGINE=InnoDB;

create table `post`
    -- --------------------------------------------------
    --  Table Structure for `main.Post`
    -- --------------------------------------------------
    CREATE TABLE IF NOT EXISTS `post` (
        `id` integer AUTO_INCREMENT NOT NULL PRIMARY KEY,
        `user_id` integer NOT NULL
    ) ENGINE=InnoDB;
```

测试代码：

```golang
package main

import (
    "fmt"
    "github.com/astaxie/beego/orm"
    _ "github.com/go-sql-driver/mysql" // import your used driver
)

// Model Struct
type Post struct {
	Id int
    User *User `orm:"rel(fk)"` // RelForeignKey relation
}

type User struct {
	Id    int
	Posts []*Post `orm:"reverse(many)"` // fk 的反向关系
}


func init() {
    // set default database
	orm.RegisterDataBase("default", "mysql", "wty:97112500@tcp(127.0.0.1:3306)/test?charset=utf8")

    // register model
    orm.RegisterModel(new(User),new(Post))

    // create table
    orm.RunSyncdb("default", false, true)
}

func main() {
	
	o := orm.NewOrm()
	user := User{Id:1}
	post1 := Post{Id:1,User:&user}
	post2 := Post{Id:2,User:&user}
	storage := []*Post{&post1,&post2}
	user.Posts = storage
	id,err := o.Insert(&post1)
	fmt.Println(id,err)
	id,err = o.Insert(&post2)
	fmt.Println(id,err)
	id,err = o.Insert(&user)
	fmt.Println(id,err)
	fmt.Println("Finish")
}

```

实现场景上也并没有外键，与上面类似，仅仅是实现了对应关系的主键。

### 多对多关系

RelManyToMany

```golang
type Post struct {
    ...
    Tags []*Tag `orm:"rel(m2m)"` // ManyToMany relation
    ...
}
```

对应的反向关系RelReverseMany：

```golang
type Tag struct {
    ...
    Posts []*Post `orm:"reverse(many)"`
    ...
}
```

这里我就不实验了，感觉上，表关系设置在数据库层面并没有太多的支持，可能在实际orm操作上有类似级联删除的机制，这个我之后在继续进行测试。

更多的详细设置还是要参考文档。