# golang学习笔记

标签：golang 读书笔记 语言基础
参考资料：
基础部分参见go tour

* [go源码剖析与学习笔记](https://github.com/qyuhen/book)

过于基础的部分跳过，这里具体记录golang语言的特点

非常基础但是容易遗漏的常识：

* golang中没有分号
* 所有的左括号需要与当前行在一起，不能另起一行，也就是google开源项目风格。
* golang是静态变量语言，没有OO，但是有接口。
* 特殊的只写变量符`_`，用于占位，无其他意义。
* 定义后未使用的局部变量会直接报错，全局变量没有问题。同时注意go编译器会在编译过程后自动标准化源代码。
* 风格上，一般函数的第一个返回值作为错误，如果第一个返回值为nil，则视为没有错误。
* nil的判断为直接使用布尔表达式判断。

## 符号

### 变量

* 符号`:=`可以简要声明，不用给出具体的变量类型。
* 可以像python一样同时给多个变量按顺序赋值（python中的元组赋值，大概）。比如`data,i := [3]int{0,1,2},0`，其中多变量赋值的时候会先计算完所有相关值，再从左到右依次赋值。这里`i,data[i] = 2,100`时，data数组指向的是0号元素

### 常量

* 使用const声明常量
* 常量值除了是表达式外，还可以是len、cap、unsafe.Sizeof等编译器可以确定结果的函数返回值。
* 枚举用的关键字iota，可以从0开始计数的自增枚举值:
```golang
const (Sunday = iota //0
    Monday      //1
    Tuesday     //2
    Wednesday   //3
    Thursday    //4
    Friday      //5
    Saturday    //6
)
```
后面同样可以使用iota的表达式，表达式中iota的值同样会自增。同时在同一常量组中可以提供多个iota，它们各自增长。iota的增长被打断需要进行显式地恢复（再次用iota声明）

### 引用类型

引用类型包括slice、map和channel。它们有复杂的内部结构，除了申请内存外，还需要初始化相关属性。

使用内置函数new计算类型大小，分配零值内存，并返回指针。make则会被编译器翻译成具体的创建函数，由其分配内存和初始化成员结构，返回对象而非指针。

make函数第一个参数接收初始化的数据类型，第二个参数接收零值空间的个数。比如`make([]int,3)`，就会返回一个有3个零值空间的一维数组。

new与make均为golang中用于内存分配的原语。new与make的区别：

* new是用来分配内存的内建函数，但是并不初始化内存，只是将其置零，并返回它的地址，即一个\*T的类型的值。它是指向它所对应的空间的地址，值为0。new是得到指针的重要方法。除此之外还能够得到指针的方法只有使用地址引用符`&`。
* make只用来创建slice、map和channel，并返回一个初始化（不是置零）的，类型为T的指

### 类型转换

go语言不支持隐式的类型转换，只能够使用显式的类型转换。

### 字符串

go语言中的字符串是不可变的，使用索引可以访问到字符对应的编码制（默认为ASCII码），但是不能够修改，同时不能够使用索引取得指针。（不像C风格数组，尾部不包含NULL）一般而言使用bytes数组进行修改。`[]byte(str)`进行初始化，`string([]bytes)`将其转换回去。

可以使用\`来定义不作转义的原始字符串，支持换行。跨行连接字符串的时候，加号要在上一行的行尾。

utf8的使用可以参考标准库中的utf8包。

### 指针

支持指针，指针的指针以及带包前缀的指针。没有->操作符，统一使用.来访问属性成员。不能对指针进行加减法操作。可以在unsafe.Pointer和任意指针类型之间进行转换。

uintptr是个特殊的类型，它是golang的内置类型，是能存储指针的整形。在64位平台的时候底层的数据类型是uint64。虽然unsafe.Pointer指针可以被转换为uintptr类型，然后保存到指针型数值变量中用以进行必要的指针数值运算，但是这可能是非法的内存地址。

总而言之，unsafe包在不同的版本会返回不同的值，而且编译器难以有提示，所以尽量避免使用。

## 表达式

### 布尔运算符

大部分与C类似。

其中`^`作为单目运算符的时候是取反，作为双目运算符的时候是按位异或。

比较特殊的是`&^`清除标志位，比如`0110 &^ 1011 = 0100`，清除了第二位的标志位。

### 初始化

初始化复合对象必须使用类型标签，使用花括号，比如`var b []int = {1,2,3}`

### range

类似迭代器的操作，返回(索引，值)或(键，值)。对于byte数组，如果内部有utf8编码，逐个索引会将其作为ASCII解码，但是range处理的话可以正常识别多字节utf8编码。

注意：range会复制对象。对于range给出的迭代器的修改是不会反应到原始对象中的。（所以循环中如果是数组的话需要注意，slice作为引用对象数值是会被修改的，同样要注意。）

### switch

跟C有很大的不同，不会自动执行后面的语句，除非加上关键字`fallthrough`，同样可以使用break。可以省略表达式当作if语句使用。

### break和continue

与python类似，配合标签，可以跳出多层循环。

## 函数

不支持嵌套、重载和默认参数。

可以使用匿名函数，可以将函数作为参数传递。

### 变参

变参实质上是一个slice，只能有一个，而且必须是最后一个。
```golang
func test(s string, n ...int) string{

}
```

其中的n便是int型的slice，作为变参。

使用slice对象作为变参的时候必须展开
```golang
s:=[]int{1,2,3}
test("string",s...)
```

### 返回值

不能够使用容器对象接受多个返回值，只能用多个变量或者用`_`忽略。

然而多返回值可以直接作为其他函数的调用参数（也可作为变参）
```golang
func test() (int,int){
    return 1,2
}
func add(x,y int) int{
    return x+y
}

func main(){
    fmt.Println(add(test()))
}
```

命名返回参数可以看作与形参类似的局部变量，最后由return隐式返回。（如果有同名的局部变量，则需要显式返回）

值得一提的是命名返回参数允许defer延迟调用通过闭包来读取和修改
```golang
func add(x,y int) (z int){
    defer func(){
        z += 100
    }()

    z = x + y
    return
}

add(1,2) //103
```

这里的结果是103，因为在执行完`z=x+y`并return之后，defer中的语句通过闭包访问到了z并修改了它的值。

### 延迟调用

关键字defer用于注册延迟调用。这些调用直到ret前才被执行，通常用于释放资源或者错误处理。多个defer注册的时候会按照FILO顺序执行。哪怕函数或某个延迟调用发生错误，这些调用依旧会被执行。

延迟调用的函数参数在注册的时候求值或进行复制，也可以通过指针或闭包“延迟”读取。

滥用defer，特别是在大循环中，会导致性能问题。

### 错误处理

golang中没有结构化异常，使用panic来抛出错误，recover捕获错误。由于panic、recover参数类型为interface{}，因此可以抛出任何类型的对象。

延迟调用中引发的错误，可以被后续延迟调用捕获，但只能捕获最后一个。

捕获函数recover只有在延迟调用内直接调用才会终止错误，否则总是返回nil。任何未捕获的错误都会沿调用堆栈向外传递。
```golang
func except(){
    recover()
}

func test(){
    defer except()
    panic("test panic")
}
```

除了panic引发的中断性错误外，还可以返回error类型错误对象表示函数的调用状态。可以使用`errors.New`或是`fmt.Errorf`函数用于创建实现error接口的错误对象，通过判断错误对象实例来确定具体的错误类型。

如果导出关键流程出现不可修复性错误时应该使用panic终止程序，其他时候使用error

## 数据

### 数组
* Array不是引用类型，作为参数传递的时候必定进行值拷贝行为。由于值拷贝可能造成性能问题，一般建议使用slice或数组指针。
* slice是通过内部指针和相关属性引用数组片段，以实现变长方案。slice可以通过`silce := array[low:high:max]`进行创建

slice在ruhntime.h中的定义
```golang
struct Slice
{
    byte* array;
    uintgo len;
    uintgo cap;
};
```

可以使用make动态创建slice，避免了数组必须使用常量做长度的麻烦，还可以用指针来直接访问底层数组。

slice相关行为

* append，向slice尾部添加数据，返回新的slice对象。简单来说，就是在`array[slice.high]`中写数据。一旦超出原slice.cap限制，就会重新分配底层数组，即使原数组并未填满。（因此大批量添加数据时，建议一次性分配足够大的空间，以减少内存分配和数据复制开销。或初始化足够长的len属性。及早释放不再使用的slice对象，避免持有过期数组，造成GC无法回收）
* copy。函数copy在两个slice间复制数据，复制长度以len晓的为准。两个slice可指向同一底层数组，允许元素区间重叠。`copy(dst,src)`

### Map

引用类型，哈希表。键必须是支持相等运算符的类型，值可以是任意类型，没有限制。
预先给make参数一个合理元素数量参数有助于提升性能，因为事先申请一大块内存可以避免后续操作时频繁扩张。`m:=make(map[string]int,1000)`

从map中取回的是一个value的临时复制品，对其成员的修改是没有任何意义的。如果试图修改，正确的做法是完整替换value或者使用指针。（当map因扩张而重新哈希的时候，各键值项的存储位置都会发生变化，因此map被设计成not addressable，类似m\[1\].name这种期望透过原value指针修改成员的行为自然会被禁止）

* 索引操作。`v,ok := m[key]`，其中ok是个布尔值，为true表示该键存在，反之则表示不存在。如果key不存在，则会返回`\0`而不会出错。
* `delete(map,key)`删除。如果key不存在，不会出错。
* 迭代。可以使用range。不能保证迭代返回次序，通常是随机结果。

可以在迭代时安全删除键值，但是如果期间有新增操作，可能会出现意外。

### 结构体

结构体成员支持匿名结构，但是初始化的时候不能直接对匿名结构进行初始化

```golang
type File struct{
    name string
    size int
    attr struct{
        perm int
        owner int
    }
}

func main(){
    f:=File{
        name :"test.txt",
        size:1025,
        // attr :{0755,1}, //Error : missing type in composite literal
    }
    f.attr.owner = 1
    f.attr.perm = 2
}   
```

可以用\`定义字段标签，用反射读取。标签是类型的组成部分。

空结构节省内存，比如用来实现set数据结构，或者实现没有状态只有方法的静态类。

```golang
var null struct{}

set := make(map[string]struct{})
set["a"] = null
```

#### 匿名字段

匿名字段是一种语法糖。从根本上来说，就是一个与成员类型同名（不含包名）的字段。被匿名嵌入的可以是任何类型，包括指针。

```golang
type User struct{
    name string
}
type Manager struct{
    User
    title string
}

m := Manager{
    User: User{"Tom"},//匿名字段的显式字段名和类型名相同。
    title:"Administrator",
}
```

可以像普通字段那样访问匿名字段成员，编译器从外向内逐级查找所有层次的匿名字段，直到发现目标或是出错。

不能同时嵌入某一类型和其指针类型，因为它们名字相同。

结构体的内存布局上和C的相同，都是根据声明顺序顺序摆放，没有额外的object信息。

## 方法

方法总是绑定对象实例，并隐式地将实例作为第一实参只能为当前包内命名的类型定义方法。参数receiver类型可以是T或\*T，基类型T不能是接口或指针。

receiver T 和\*T的差别：方法不过是特殊的函数，receiver类型是否为指针并不影响该变量在方法内的使用情况，只是影响方法内对该变量的修改是否会影响到方法外，以及方法外的调用对象是否为指针。

从1.4开始，不再支持多级指针查找方法成员。

### 通过匿名字段实现继承

通过匿名字段，可以获得和继承类似的复用能力。依据编译器查找次序，只需要在外层定义同名方法，就可以实现重写。

```golang
type User struct{
    id int
    name string
}

type Manager struct{
    User
    title string
}

func (self *User) ToString() string{
    return fmt.Sprintf("User:%p,%v",self,self)
}

func (self *Manager) ToString() string{
    return fmt.Sprintf("Manager:%p,%v",self,self)
}

func main(){
    m := Manager{User{1,"Tom"},"Administrator"}
    fmt.Println(m.ToString())
    fmt.Println(m.User.ToString())
}
```

### 表达式

根据调用者不同，方法分为两种表现形式：

1. `instance.method(args...)`，称为method value，绑定实例。
2. `<type>.func(instance,args...)`，称为method expression，需要显式传参。

```golang
type User struct{
    id int
    name string
}
func (self *User) Test(){
    fmt.Printf("%p,%v\n",self,self)
}
func main(){
    u := User{1,"Tom"}
    u.Test()
    
    mValue := u.Test
    mValue()//隐式传递receiver

    mExpression := (*User).Test
    mExpression(&u) //显式传递receiver
}
```

在汇编层面上，method value和闭包的实现方式相同，实际返回FuncVal类型对象。`FuncVal{method_address,receiver_copy}`

## 接口

接口是一个或多个方法签名的集合，任何类型的方法集中只要拥有与之对应的全部方法，就表示它“实现”了该接口，无须在该类型上显式添加接口声明。

对应方法，指有相同名称、参数列表以及返回值。

空接口interface{}没有任何方法签名，也就意味着任何类型都实现了空接口。类似于面向对象语言中的根对象object。

### 执行机制

接口对象由接口表指针和数据指针组成。

``` C++
struct Iface
{
    Itab* tab;
    void* data;
};

struct Itab
{
    InterfaceType* inter;
    Type* type;
    void (*func[])(void);
};
```

接口表存储元数据信息，包括接口类型、动态类型，以及实现接口的方法指针。
数据指针持有的是目标对象的只读复制品，复制完整对象或指针。

接口转型：`var.(struct)`，可以将接口var转换为实际类型struct输出值。
接口转型返回临时对象，只有使用指针才能够修改其状态。

```golang
type User struct{
    id int
    name string
}
func main(){
    u := User{1,"Tom"}
    var vi,pi interface{} = u,&u

    //vi.(User).name = "Jack" //Error:cannot assign to vi.(User).name
    pi.(*User).name = "Jack"
}
```

显然，只有tab和data都为nil时，接口才为nil。

类型推断可以判断接口对象是否为某个具体的接口或类型。`value,ok := interface.(ConcreteStruct)`，如果接口对象属于某个具体的接口或类型，ok为true。还可以使用switch做批量的类型判断，不支持fallthrough。

超集接口对象可以转换成子集接口，但是反之会出错。

### 技巧

让编译器来检查，以确保某个类型实现接口`var _ fmt.Stringer = (*Data)(nil)`

## 并发

go在语言层面上对并发变成提供支持，一种类似协程，被称为gorouine的机制。另外还有与之配套的channel类型，用以实现“以通讯来共享内存”的CSP模式。

默认情况下，进程启动后仅允许一个系统线程服务于goroutine。可使用环境变量或标准库函数`runtime.GOMAXPROCS`修改，让调度器用多个线程实现多核并行，而不仅仅是并发。

### Channel

channel是CSP模式的具体实现，用于多个goroutine通讯，内部实现了同步以确保并发安全。使用make函数进行初始化,可以接受第二个参数，表示缓存区的大小。如果没有规定缓存区，如果一方阻塞，则整条线程阻塞。阻塞可能造成死锁，详见操作系统课。

可以使用range channel从队列迭代接受数据，直到close。队列的关闭并不是必须的，按照网上的说法，关闭是一种状态而不是信息，关闭并不是为了垃圾回收服务的，而是给接收方提供终止发送的信号，表示数据已经发送完毕。向close channel发送数据引发panic错误，接收则立即返回零值。而nil channel，不论收发都会被阻塞。

channel是第一类对象，可以传参（内部实现为指针）或者作为结构成员。

#### 单向

可以将channel隐式转换为单向队列，只收或只发，进行部分的封装。但是不能将单向队列转为普通队列。

疑问：如何规定谁收，谁发？
```golang
c := make(chan int,3)

var send chan<- int = c //send-only
var recv <-chan int = c //receive-only

send<-1
<-recv
```

#### 选择

如果需要同时处理多个channel，可以使用select语句。它随机选择一个可用的channel进行收发操作而不会阻塞，或者执行default case。

#### sync.WaitGroup

sync.WaitGroup提供了三个方法Add(),Done()和Wait()，可以帮助主进程在goroutine执行完毕前阻塞。Add(x int)表示阻塞的线程数量，Done()等价于Add(-1)，Wait()在主线程进行阻塞直到线程完成。一般在goroutine中`defer wg.Done()`。

#### 超时控制

使用select+time.After的方式实现超时控制。

```golang
func main(){
    w := make(chan bool)
    c := make(chan int,2)

    go func(){
        select{
            case v:= <-c : fmt.Println(v)
            case <-time.After(time.Second*3): fmt.Println("timeout.")
        }

        w<-true
    }

    //c<-1 //case time out
    <-w
}
```

## 包

### 工作空间

工作空间必须由bin、pkg、src三个目录组成。其中bin为go install安装目录，pkg为go build生成静态库.a存放目录，src为源代码目录。

可以在GOPATH环境变量列表中添加多个工作空间，但是不能和GOROOT相同。

### 包结构

源文件头部以`package <name>`声明包名称。包名称类似于namespace，与包所在目录名、编译文件名无关。可执行文件必须包含`package main`以及入口函数`main`。

包中成员如果首字母为大写，则可以被包外成员访问，否则只能被包内成员访问。

### 导入

```golang
import   "yuhen/test" //默认模式：test.A
import M "yuhen/test" //包重命名：M.A
import . "yuhen/test" //简便模式；A
import _ "yuhen/test" //非导入模式：仅让该包执行初始化函数。
```

对于当前目录下的子包，除使用默认完整导入路径外，还可使用local方法：`import "./test"`，此为本地模式，仅对于`go run main.go`有效


### 初始化

每个源文件可以定义一个或多个初始化函数，但是编译器并不保证多个初始化函数的执行次序。初始化函数为`func init()`，它无法被直接调用，只能在单一线程被调用且仅执行一次，会在包所有全局变量初始化后执行。在所有初始化函数结束后才执行`main.main`。

### 文档

使用扩展工具godoc自动提取注释生成文档

包内可以使用专门的doc.go保存package帮助信息，包文档第一整句（中英文句号结束）被当作packages列表说明。

#### bug

以`BUG(author)`开头的注释，将被识别为已知bug，显示在bugs区域。

## 进阶

### 指针陷阱

对象内存分配会受编译参数的影响。举个例子，当函数返回对象指针时，必然在堆上分配。可如果该函数被内联，那么这个指针就不会跨栈帧使用，就有可能在栈上直接分配，以实现代码优化的目的。

指针除了正常指针外，还有unsafe.Pointer和uintptr。其中uintptr被GC当作普通整数对象，它不能阻止所“引用”对象被回收。

合法的unsafe.Pointer会被当作正常指针，同样能确保对象不被回收。

### cgo

通过cgo，可在Go和C/C++代码间相互调用，受CGO_ENABLED参数限制。如下

```golang
package main
/*
    #include <stdio.h>
    #include <stdlib.h>

    void hello(){
        printf("Hello,world!\n");
    }
*/
import "C"

func main(){
    C.hello()
}
```

调试cgo代码是很麻烦的事情，建议单独保存在.c文件中，将其当作独立的C程序进行调试。

test.h

```c
#ifndef __TEST_H__
#define __TEST_H__
void hello();
#endif
```

test.c

```c
#include <stdio.h>
#include "test.h"

void hello(){
    printf("Hello,world\n");
}

#ifdef __TEST__ //避免和Go bootstrap main冲突。
int main(int argc,char* argv[]){
    hello();
    return 0;
}
#endif
```

main.go

```go
package main
/*
    #include "test.h"
*/

import "C"

func main(){
    C.hello()
}
```

编译和调试C，只需要在命令行提供宏定义即可。`gcc -g -D__TEST__ -o test test.c`

由于cgo仅扫描当前目录，如果需要包含其他C项目，可以在当前项目新建一个C文件，然后用`#include`指令将所需的.h、.c都包含近来。有时候还需要在CFLAGS中使用-l参数指定路径，或者指定-std参数。

#### 参数

可以使用`#cgo`命令定义CFLAGS、LDFLAGS等参数，自动合并多个设置。
```golang
/*
    #cgo CFLAGS: -g
    #cgo CFLAGS: -I./lib -D__VER__=1
    #cgo LDFLAGS: -lpthread

    #include "test.h"
*/
import "C"
```

可以设置GOOS、GOARCH编译条件，其中空格表示OR，逗号AND，感叹号NOT。

#### 字符串

```golang
/*
    #include <stdio.h>
    #include <stdlib.h>

    void test(char *s){
        printf("%s\n",s);
    }

    char* cstr(){
        return "abcde";
    }
*/
import "C"

func main(){
    s := "Hello,World!"

    cs := C.CString(s) //该函数在C heap分配内存，需要调用free释放。
    defer C.free(unsafe.Pointer(cs)) //#include<stdlib.h>

    C.test(cs)
    cs = C.cstr()

    fmt.Println(C.GoString(cs))//abcde
    fmt.Println(C.GoStringN(cs,2))//ab
    fmt.Println(C.GoBytes(unsafe.Pointer(cs),2))//[97 98]
}
```

#### struct/Enum/Union

对struct、enum支持良好，union会被转换成字节数组。如果没有使用typedef定义，那么必须添加struct_、enum_、union_前缀。

#### Export

导出Go函数给C调用，需要使用`//export`标记，建议在独立头文件中声明函数原型，避免"duplicate symbol"错误。

#### shared library

在cgo中使用C共享库

## Reflect

go语言不存在运行期类型对戏那个，实例也不会使用附加字段来表明身份。只有在转换成接口的时候，才会在其itab内部存储与类型有关的信息，Reflect所有操作都依赖于此。

导出类型的全部成员字段信息，包括非导出和匿名字段

```golang
type User struct{
    Username string
}

type Admin struct{
    User
    title string
}

func main(){
    var u Admin
    t := reflect.TypeOf(u)
    for i,n := 0; t.NumField();i<n;i++{
        f := t.Field(i)
        fmt.Println(f.Name,f.Type)
    }
}
```

如果是指针，应该先使用Elem方法获取目标类型。指针本身是没有字段成员的。

同时，接口的receiver是否为指针也会影响其方法集。下面代码显示所有的方法集

```golang
type User struct{
}

type Admin struct{
    User
}

func (*User) ToString(){}

func (Admin) test(){}

func main(){
    var u Admin
    methods := func(t reflect.Type){
        for i,n :=0,t.NumMethod(); i<n;i++{
            m = t.Method(i)
            fmt.Println(m.Name)
        }
    }
    fmt.Println("value interface")
    methods(reflect.TypeOf(u)) //test
    fmt.Println("pointer interface")
    methods(reflect.TypeOf(&u))//ToString test
}
```

方法Implements判断是否实现了某个具体的接口，AssignableTo、ConvertibleTo用于赋值和转换判断。

Value和Type使用方法类似，包括使用Elem获取指针目标对象。除了Int、String等转换方法，还可以返回interface{}。只是非导出字段无法使用，需要用CanInterface判断一下。

将对象转换为接口，会发生复制行为。该复制品只读，无法被修改。所以要通过接口改变目标对象的状态，必须是pointer-interface。

使用Method方法获取方法的参数、返回值类型等信息

```golang
func info(m reflect.Method){
    t := m.Type
    fmt.Println(m.Name)
    for i,n := 0,t.NumIn();i<n;i++{
        fmt.Println("  in[%d] %v\n",i,t.In(i))
    }
    for i,n := 0,t.NumOut();i<n;i++{
        fmt.Println("  out[%d] %v\n",i,t.Out(i))
    }
}
```

### 动态调用

动态调用的方法很简单，就是按照In列表准备好所需的参数即可。（不包括receiver）。非导出方法无法调用，甚至无法透过指针操作，因为接口类型信息中没有该方法地址。

```golang
func main(){
    d := new(Data)
    v := reflect.ValueOf(d)

    exec := func(name string, in []reflect.Value){
        m = v.MethodByName(name)
        out := m.Call(in)
        // out := m.CallSlice(in) //get slice

        for _,v := range out{
            fmt.Println(v.Interface())

        }
    }

    exec("Test",[]reflect.Value{
        reflect.ValueOf(1),
        reflect.ValueOf(2),
    })
}
```

### 实现泛型

利用Make、New等函数，可以实现近似泛型操作。下面为动态创建一个slice的泛型实现

```golang
var(
    Int = reflect.TypeOf(0)
    String = reflect.Typeof("")
)

func Make(T reflect.Type, fptr interface{}){
    //实际创建slice的包装函数
    swap := func(in []reflect.Value) []reflect.Value{
        //省略算法部分
        //返回和类型匹配的slice对象
        return []reflect.Value{
            reflect.MakeSlice(
                reflect.SliceOf(T), //slice type
                int(in[0].Int()),   //len
                int(in[1].Int())    //cap
            ),
        }
    }

//将函数变量指针指向swap函数
    fn := reflect.ValueOf(ptr).Elem()

    //获取指针类型，生成所需的swap function value
    v := reflect.MakeFunc(fn.Type(),swap)

    //修改函数指针实际指向，也就是swap。
    fn.Set(v)
}

func main(){
    var makeints func(int,int) []int
    var makestrings func(int,int) []string

    //用相同的算法生成不同类型的创建函数
    Make(Int,&makeints)
    Make(String,&makestrings)

    //按照实际类型使用
    x := makeints(5,10)
    fmt.Println(x)

    s := makestrings(3,10)
    fmt.Println(s)
}
```