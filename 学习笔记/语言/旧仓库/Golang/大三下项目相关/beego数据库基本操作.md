# beego数据库基本操作实践

标签：beego 数据库 实践

beego中对数据库的操作比较多，本文主要实践除了orm以外的几种对数据库的访问方式。而且实际上，如果你经常使用数据库而不怎么用orm的话，直接使用下面的方法比除了CURD四种情况外使用orm要方便很多。

其实原始的中文文档讲的就很清除，这里只是做一个自己看的总结而已，这样我查东西的时候就不用忍受抽风的延迟了。

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

下面是一段例子，及其简单的测试了一下QueryTable的具体使用方法

```golang
package main

import (
	"fmt"
	_ "apiproject/routers"
    "github.com/astaxie/beego/orm"
    _ "github.com/go-sql-driver/mysql"
	"github.com/astaxie/beego"
)

type Test struct {
    Id   int
    Name string `orm:"size(100)"`
}

func init() {
    orm.RegisterDriver("mysql", orm.DRMySQL)

	orm.RegisterDataBase("default", "mysql", "wty:97112500@tcp(127.0.0.1:3306)/test?charset=utf8")
	orm.RegisterModel(new(Test))
	orm.RunSyncdb("default", false, true)
}



func main() {
	o := orm.NewOrm()
	test := Test{Name:"test2"}
	
	//insert
	id, err := o.Insert(&test)
	fmt.Printf("ID:%d, ERR:%v\n",id,err)

	qs := o.QueryTable("test")
	qs.Filter("Name","test2")

	var users []*Test

	num,err := o.QueryTable("test").Filter("Name","test2").All(&users)
	fmt.Println("Returned Rows Num:",num,err)

	for i:=int64(0) ; i<num ; i++{
		fmt.Println(users[i].Id)
	}

	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}
	beego.Run()
}

```

这个例子会返回表test中所有名字为test2的对象到Tset数组users中，并打印它们的id

总结

#### 接口实现

前排提醒，所有的QuerySeter的方法都会返回一个新的QuerySeter对象，因此可以自由的使用链式方法（应该是叫这个名字）

##### 读取

QuerySeter可以使用Filter方法，接受参数并进行数组上的判断。多个Filter之间使用And进行连接。

Filter方法的第一个参数为描述字段的名称与sql操作符的符合，第二个参数为具体的值。

此外还有用来排除的Exclude，自定义的条件表达式SetCond，聚类GroupBy，参数使用expr是用的排序表示OrderBy，对应sql重复的Distinct，参数使用expr是使用的关系查询RelatedSel，返回结果行数Count

限制最大行数的Limit方法与设置偏移行数的Offset方法。

##### 更新

甚至还有依据当前查询条件进行批量更新操作

```golang
num, err := o.QueryTable("user").Filter("name", "slene").Update(orm.Params{
    "name": "astaxie",
})
fmt.Printf("Affected Num: %s, %s", num, err)
// SET name = "astaixe" WHERE name = "slene"
```

以及原子操作增加字段值

```golang
// 假设 user struct 里有一个 nums int 字段
num, err := o.QueryTable("user").Update(orm.Params{
    "nums": orm.ColValue(orm.ColAdd, 100),
})
// SET nums = nums + 100
```

##### 删除

```golang
num, err := o.QueryTable("user").Filter("name", "slene").Delete()
fmt.Printf("Affected Num: %s, %s", num, err)
// DELETE FROM user WHERE name = "slene"
```

##### 插入

用一次prepare多次insert插入，以提高批量插入的速度。具体速度还没有测试，感觉应该会挺快的。

```golang
var users []*User
...
qs := o.QueryTable("user")
i, _ := qs.PrepareInsert()
for _, user := range users {
    id, err := i.Insert(user)
    if err == nil {
        ...
    }
}
// PREPARE INSERT INTO user (`name`, ...) VALUES (?, ...)
// EXECUTE INSERT INTO user (`name`, ...) VALUES ("slene", ...)
// EXECUTE ...
// ...
i.Close() // 别忘记关闭 statement
```

#### 返回结果

##### 返回所有结果

1. All方法

All方法返回所有结果到指定的对象数组中，收到Limit的限制（似乎可以自己改），默认最大行数为1000

```golang
var users []*User
num, err := o.QueryTable("user").Filter("name", "slene").All(&users)
fmt.Printf("Returned Rows Num: %s, %s", num, err)
```

甚至可以指定返回的字段

```golang
type Post struct {
    Id      int
    Title   string
    Content string
    Status  int
}

// 只返回 Id 和 Title
var posts []Post
o.QueryTable("post").Filter("Status", 1).All(&posts, "Id", "Title")
```

2. Values方法

Values方法返回结果集的key=>value值，key为Model里的Field name，value为interface{}类型，因此要使用断言来获取真实值，比如`Name : m["Name"].(string)`

下面为返回指定的Field数据（Id、Name）

```golang
var maps []orm.Params
num, err := o.QueryTable("user").Values(&maps)
if err == nil {
    fmt.Printf("Result Nums: %d\n", num)
    for _, m := range maps {
        fmt.Println(m["Id"], m["Name"])
    }
}
```

可以使用expr级联返回需要的数据

```golang
var maps []orm.Params
num, err := o.QueryTable("user").Values(&maps, "id", "name", "profile", "profile__age")
if err == nil {
    fmt.Printf("Result Nums: %d\n", num)
    for _, m := range maps {
        fmt.Println(m["Id"], m["Name"], m["Profile"], m["Profile__Age"])
        // map 中的数据都是展开的，没有复杂的嵌套
    }
}
```

3. ValuesList方法

与Value类似，不过返回的结果集以slice的形式存储。结果的排列与Model中定义的Field顺序是一致的。返回的每个元素值以string形式保存。

```golang
var lists []orm.ParamsList
num, err := o.QueryTable("user").ValuesList(&lists)
if err == nil {
    fmt.Printf("Result Nums: %d\n", num)
    for _, row := range lists {
        fmt.Println(row)
    }
}
```

##### 返回单条记录

One方法尝试返回单条记录，也可以返回单条记录，返回单条记录的时候形式与上面完全一样，这里就不了列举了。

```golang
var user User
err := o.QueryTable("user").Filter("name", "slene").One(&user)
if err == orm.ErrMultiRows {
    // 多条的时候报错
    fmt.Printf("Returned Multi Rows Not One")
}
if err == orm.ErrNoRows {
    // 没有找到记录
    fmt.Printf("Not row found")
}
```

### 关系查询

好吧，我承认我失误了，我没有看这里

#### OneToOne关系

##### 根据User拿Profile

已经取得了User对象，查询Profile

```golang
user := &User{Id: 1}
o.Read(user)
if user.Profile != nil {
    o.Read(user.Profile)
}
```

完整代码

```golang
package main

import (
    "fmt"
    _ "apiproject/routers"
    "github.com/astaxie/beego/orm"
    _ "github.com/go-sql-driver/mysql"
    "github.com/astaxie/beego"
)

type User struct {
    Id   int
    Name string `orm:"size(100)"`
	Profile *Profile `orm:"null;rel(one);on_delete(set_null)"`
}

type Profile struct{
	Id 	int
	Name	string
	User	*User	`orm:"reverse(one)"`
}

func init() {
    orm.RegisterDriver("mysql", orm.DRMySQL)

	orm.RegisterDataBase("default", "mysql", "wty:97112500@tcp(127.0.0.1:3306)/test?charset=utf8")
	orm.RegisterModel(new(User),new(Profile))
	orm.RunSyncdb("default", false, true)
}



func main() {
	o := orm.NewOrm()
	u := User{Name:"wty"}
	p := Profile{Name:"wtyOrigin"}
	u.Profile = &p
	p.User = &u
	fmt.Println(u.Profile)
	o.Insert(&p)
	o.Insert(&u)
	

	user := &User{Id:5}//id必须要为已经有的id
	o.Read(user)
	fmt.Println("Get user",user)
	fmt.Println(user.Profile)//&{5  <nil>}
	if user.Profile!=nil{
		o.Read(user.Profile)//注意，user.Profile已经是指针了，不需要加&
	}
	fmt.Println(user.Profile)//&{5 wtyOrigin <nil>}

	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}
	beego.Run()
}

```

##### 直接查询

直接关联查询（我以前以为只能够用这种方式，感觉好蠢）

```golang
user := &User{}
o.QueryTable("user").Filter("Id", 1).RelatedSel().One(user)
// 自动查询到 Profile
fmt.Println(user.Profile)
// 因为在 Profile 里定义了反向关系的 User，所以 Profile 里的 User 也是自动赋值过的，可以直接取用。
fmt.Println(user.Profile.User)

// [SELECT T0.`id`, T0.`name`, T0.`profile_id`, T1.`id`, T1.`age` FROM `user` T0 INNER JOIN `profile` T1 ON T1.`id` = T0.`profile_id` WHERE T0.`id` = ? LIMIT 1000] - `1`
```

##### 通过UserId拿Profile

通过User反向查询Profile

```golang
var profile Profile
err := o.QueryTable("profile").Filter("User__Id", 1).One(&profile)
if err == nil {
    fmt.Println(profile)
}
```

#### ManyToOne关系

Post和User是ManyToOne关系，也就是说ForeignKey为User

```golang
type Post struct {
    Id    int
    Title string
    User  *User  `orm:"rel(fk)"`
    Tags  []*Tag `orm:"rel(m2m)"`
}
var posts []*Post
num, err := o.QueryTable("post").Filter("User", 1).RelatedSel().All(&posts)
if err == nil {
    fmt.Printf("%d posts read\n", num)
    for _, post := range posts {
        fmt.Printf("Id: %d, UserName: %d, Title: %s\n", post.Id, post.User.UserName, post.Title)
    }
}
// [SELECT T0.`id`, T0.`title`, T0.`user_id`, T1.`id`, T1.`name`, T1.`profile_id`, T2.`id`, T2.`age` FROM `post` T0 INNER JOIN `user` T1 ON T1.`id` = T0.`user_id` INNER JOIN `profile` T2 ON T2.`id` = T1.`profile_id` WHERE T0.`user_id` = ? LIMIT 1000] - `1`
```

完整版。（但我平时就是这么写的……）。但这样子有一个好处就是类与类之间的联系更加的紧密，post可以直接拿到与其相关的user

```golang
package main

import (
	"fmt"
	_ "apiproject/routers"
    "github.com/astaxie/beego/orm"
    _ "github.com/go-sql-driver/mysql"
	"github.com/astaxie/beego"
)

type User struct {
    Id   int
    Name string `orm:"size(100)"`
	Post []*Post `orm:"reverse(many)"`
}

type Post struct{
	Id 	int
	Title	string
	User	*User	`orm:"rel(fk)"`
}

func init() {
    orm.RegisterDriver("mysql", orm.DRMySQL)

	orm.RegisterDataBase("default", "mysql", "wty:97112500@tcp(127.0.0.1:3306)/test?charset=utf8")
	orm.RegisterModel(new(User),new(Post))
	orm.RunSyncdb("default", false, true)
}



func main() {
	o := orm.NewOrm()

	user1 := User{Name:"wty"}
	post1 := Post{Title:"post1"}
	post2 := Post{Title:"post2"}
	post1.User = &user1
	post2.User = &user1



	user1.Post = []*Post{&post1,&post2}
	o.Insert(&user1)
	o.Insert(&post1)
	o.Insert(&post2)
	newUser := User{Id:2}
	o.Read(&newUser)
	fmt.Println("NewUser",newUser)

	var posts []*Post
	num, err := o.QueryTable("post").Filter("user_id", 2).RelatedSel().All(&posts)
	if err == nil{
		fmt.Printf("%d posts read\n",num)
		for _,post := range posts{
			fmt.Println(post.Id,post.Title,post.User.Id,post.User.Name)
		}
	}

	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}
	beego.Run()
}

```

根据Post.Title查询对应的User。因为RegisterModel时，ORM会自动建立User中Post的反向关系，所以可以直接进行查询

```golang
var user User
err := o.QueryTable("user").Filter("Post__Title", "The Title").Limit(1).One(&user)
if err == nil {
    fmt.Printf(user)
}
```

多对多关系因为用不上，所以我就不在这里列举了。

#### 载入关系字段

LoadRelated用于载入模型的关系字段，包括所有的rel/reverse-one/many关系

User是Post的ForeignKey，对应的ReverseMany关系载入

```golang
type User struct {
    Id    int
    Name  string
    Posts []*Post `orm:"reverse(many)"`
}

user := User{Id: 1}
err := dORM.Read(&user)
num, err := dORM.LoadRelated(&user, "Posts")
for _, post := range user.Posts {
    //...
}
```

完整版示例

```golang
package main

import (
	"fmt"
	_ "apiproject/routers"
    "github.com/astaxie/beego/orm"
    _ "github.com/go-sql-driver/mysql"
	"github.com/astaxie/beego"
)

type User struct {
    Id   int
    Name string `orm:"size(100)"`
	Post []*Post `orm:"reverse(many)"`
}

type Post struct{
	Id 	int
	Title	string
	User	*User	`orm:"rel(fk)"`
}

func init() {
    orm.RegisterDriver("mysql", orm.DRMySQL)

	orm.RegisterDataBase("default", "mysql", "wty:97112500@tcp(127.0.0.1:3306)/test?charset=utf8")
	orm.RegisterModel(new(User),new(Post))
	orm.RunSyncdb("default", false, true)
}



func main() {
	o := orm.NewOrm()

	user1 := User{Name:"wty"}
	post1 := Post{Title:"post1"}
	post2 := Post{Title:"post2"}
	post1.User = &user1
	post2.User = &user1



	user1.Post = []*Post{&post1,&post2}
	o.Insert(&user1)
	o.Insert(&post1)
	o.Insert(&post2)
	newUser := User{Id:2}
	o.Read(&newUser)
	fmt.Println("NewUser",newUser)

	num,err := o.LoadRelated(&newUser,"Post")
	fmt.Println("Load Related from user",num,err)
	for _,post := range newUser.Post{
		fmt.Println("show post details",post.Id,post.Title,post.User.Id,post.User.Name)
	}

	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}
	beego.Run()
}

```

### Raw

使用sql语句直接进行操作。使用Raw函数，会直接返回一个[RawSeter](https://beego.me/docs/mvc/model/rawsql.md)用以对设置的sql语句和参数进行操作。

```golang
o := orm.NewOrm()
var r RawSeter
r = o.Raw("UPDATE user SET name = ? WHERE name = ?", "testing", "slene")
```

可以创建一个新的RawSeter对象。`?`占位符号自动进行转换，参数可以使用字符串，Model Struct和Slice,Array

#### Exec

执行sql语句，返回sql.Result对象，其中有方法RowsAffected()，可以返回受到影响的行的数目。

#### QueryRow

提供sql mapper功能，可以直接将结果绑定到指定的对象中（感觉需要orm级别的名称对应，尚未进行实验）

```golang
type User struct {
    Id       int
    UserName string
}

var user User
err := o.Raw("SELECT id, user_name FROM user WHERE id = ?", 1).QueryRow(&user)
```

此外还有QueryRows，用于多个user的数组的参数获取

#### Values

返回结果集的key=>value值

```golang
var maps []orm.Params
num, err := o.Raw("SELECT user_name FROM user WHERE status = ?", 1).Values(&maps)
if err == nil && num > 0 {
    fmt.Println(maps[0]["user_name"]) // slene
}
```

此外还有ValuesList和ValuesFlat等

#### RowsToMap

将查询的结果映射到map中。

```golang
res := make(orm.Params)
nums, err := o.Raw("SELECT name, value FROM options_table").RowsToMap(&res, "name", "value")
// res is a map[string]interface{}{
//  "total": 100,
//  "found": 200,
// }
```

此外还有RowsToStruct（这个要名称对应，snake->camel）

### 事物处理

orm可以进行简单的事物操作（比如回滚）

```golang
o := NewOrm()
err := o.Begin()
// 事务处理过程
...
...
// 此过程中的所有使用 o Ormer 对象的查询都在事务处理范围内
if SomeError {
    err = o.Rollback()
} else {
    err = o.Commit()
}
```

### 构造查询：QueryBuilder

QueryBuilder与ORM上重合，但是各有利弊。ORM更加适合简单的CRUDE操作，而QueryBuilder更适合复杂的查询，例如包含子查询和多重联结。

```golang
// User 包装了下面的查询结果
type User struct {
    Name string
    Age  int
}
var users []User

// 获取 QueryBuilder 对象. 需要指定数据库驱动参数。
// 第二个返回值是错误对象，在这里略过
qb, _ := orm.NewQueryBuilder("mysql")

// 构建查询对象
qb.Select("user.name",
    "profile.age").
    From("user").
    InnerJoin("profile").On("user.id_user = profile.fk_user").
    Where("age > ?").
    OrderBy("name").Desc().
    Limit(10).Offset(0)

// 导出 SQL 语句
sql := qb.String()

// 执行 SQL 语句
o := orm.NewOrm()
o.Raw(sql, 20).QueryRows(&users)
```