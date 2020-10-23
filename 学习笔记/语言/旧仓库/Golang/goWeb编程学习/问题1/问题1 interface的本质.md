# 问题1 interface的本质是什么

在复习到interface的时候我有一个疑问：比如有一个struct Student，它实现了Sing和SayHi方法。Interface Men要求这两个方法。如果有一个`var i men; i = Student{xxx}`，那么其中发生的是复制行为还是引用行为。

## 实验

这个很简单，只要在赋值之后修改原struct的值查看interface是否发生更改即可。

```golang
package main

import (
    "fmt"
)

type Human struct{
    Name string
    age int
    phone string
}

type Student struct{
    Human
    school string
    loan float32
}

func (h Human) SayHi(){
    fmt.Printf("Hi, I am %s you can call me on %s \n",h.Name,h.phone)
}

func (h Human) Sing(lyrics string){
    fmt.Println("La la la la ...",lyrics)
}

type Men interface{
    SayHi()
    Sing(lyrics string)
}

func main(){
    mike := Student{Human{"Mike", 25, "222-222-XXX"}, "MIT", 0.00}

    var i Men

    i = mike
    fmt.Println(mike.Human.Name)
    fmt.Println(i) //{Human{"Mike", 25, "222-222-XXX"}, "MIT", 0.00}
    mike.Human.Name = "test"
    fmt.Println(i) //{Human{"Mike", 25, "222-222-XXX"}, "MIT", 0.00}
}
```

结果是不发生改变，这应该是一个复制行为。

## 本质

因此interface实际上就是一组抽象方法的集合。

那么interface在golang中是怎么进行实现的呢，如何访问到interface中的具体值呢？

访问具体值的方法我以前学过，可以通过类型判断或是reflect包，这里不再赘述。

Interface的实现是这样子的：接口对象由接口表指针和数据指针组成（来源与Go源码阅读）

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

reflect包就是基于这样的原理对interface进行读取，从而拿到存储于其中的值。