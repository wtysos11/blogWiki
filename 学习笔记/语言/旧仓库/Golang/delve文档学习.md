# delve文档学习

标签：golang delve 文档学习

[官方文档](https://github.com/go-delve/delve/tree/master/Documentation)

## usage

调试主体使用

### 大纲

delve是源代码级别的调试程序，旨在提供一个调试的接口。

如果想要传递参数给想要调试的程序，可以使用`dlv exec ./hello -- server --config conf/config.toml`

### dlv debug

## cli

[位置](https://github.com/go-delve/delve/tree/master/Documentation/cli)

### break

设置断点，`break [name] <linespec>`

其中[linespec的格式](https://github.com/go-delve/delve/blob/master/Documentation/cli/locspec.md)，主要使用`<filename>:<line>`的形式，指定文件中某一行为断点，这比较符合常理。

其他形式：

* `*<address>`，指定内存地址，可以使用十进制、八进制或者十六进制。
* `<line>`，指定当前文件的指定行
* `+<offset>`，指定当前行的后若干行
* `-<offset>`，指定当前行的前若干行
* `<function>[:line]`，直接指定函数名，其中函数名前要带上包名。函数名后的行号可以省略。

甚至可以使用正则表达式。

### breakpoints

简写`bp`，可以打印出所有的活跃断点的信息

### call

重新开始进程，插入一个函数调用（实验性质）

`call [-unsafe] <function call expression>`

刚做出来的样子，似乎基本不会用到

### clear

清除breakpoint

`clear <breakpoint name or id>`

### clearall

`clearall [<linespec>]`

### condition

设置条件断点，别名`cond`

`condition <breakpoint name or id> <boolean expression>.`

### continue

执行直到遇到断点或者程序终止，别名`c`

### deferred

在当前上下文中异步调用命令，比如print、args、local等，使其在第n次deferred call的语境下被调用

`deferred <n> <command>`

### disassemble

反汇编工具

`[goroutine <n>] [frame <m>] disassemble [-a <start> <end>] [-l <locspec>]`

如果没有指定参数，那么会指定正在调用栈中执行的函数

```bash
-a <start> <end>	disassembles the specified address range
-l <locspec>		disassembles the specified function
```

### down

move the current down

格式

```bash
down [<m>]
down [<m>] <command>
```

实在不知道什么意思

### frame

设置当前的frame，或者在其他frame下执行命令。

```bash
frame <m>
frame <m> <command>
```

第一种形式是用来制作frame，这个frame可以被print或set所使用。第二种形式在指定的frame下运行命令。

### funcs

输出函数列表

`funcs [<regex>]`，如果存在正则表达式，程序只会返回符合正则表达式的函数名

### goroutine

展示或改变当前的goroutine

```bash
goroutine
goroutine <id>
goroutine <id> <command>
```

不使用参数的话会显示当前的goroutine。

如果使用一个参数的话会切换到指定的goroutine中

如果使用更多的参数则会在指定的goroutine中执行命令。

### goroutines

列出程序的goroutines

`goroutines [-u (default: user location)|-r (runtime location)|-g (go statement location)|-s (start location)] [ -t (stack trace)]`

打印所有goroutine的信息，flag具体信息如下：

```bash
-u	displays location of topmost stackframe in user code
-r	displays location of topmost stackframe (including frames inside private runtime functions)
-g	displays location of go instruction that created the goroutine
-s	displays location of the start function
-t	displays stack trace of goroutine
```

如果没有flag，默认使用-u

### libraries

列出加载的动态运行库

### list

显示源代码

`[goroutine <n>] [frame <m>] list [<linespec>]`

别名`ls`或`l`，显示指定点或行的代码

### locals

打印局部变量

`[goroutine <n>] [frame <m>] locals [-v] [<regex>]`

当前作用于内不可视的对象会用圆括号包裹。

如果指定了正则表达式，那么只有名称匹配正则表达式的局部变量

### next

跳到源代码的下一行（不进入函数中）

### on

在遇到断点的时候执行命令

`on <breakpoint name or id> <command>.`，支持的命令有print,stack和goroutine

### print

执行表达式`[goroutine <n>] [frame <m>] print <expression>`，支持的[表达式](https://github.com/go-delve/delve/tree/master/Documentation/cli/expr.md)

### regs

打印CPU寄存器的值，`regs [-a]`

### restart

在断点或事件处重新启动进程，别名`r`

`restart [event number or checkpoint id]`

### set

改变变量的值

`[goroutine <n>] [frame <m>] set <variable> = <value>`

支持的[表达式](https://github.com/go-delve/delve/blob/master/Documentation/cli/expr.md)参考之前的。

### sources

打印源文件列表，如果使用了正则表达式则进行正则表达式匹配。

`sources [<regex>]`

### stack

打印调用栈

```bash
[goroutine <n>] [frame <m>] stack [<depth>] [-full] [-offsets] [-defer] [-a <n>] [-adepth <depth>]

-full		every stackframe is decorated with the value of its local variables and arguments.
-offsets	prints frame offset of each frame.
-defer		prints deferred function call stack for each frame.
-a <n>		prints stacktrace of n ancestors of the selected goroutine (target process must have tracebackancestors enabled)
-adepth <depth>	configures depth of ancestor stacktrace
```

### step

单步执行，会进入函数中。别名`s`

### step-instruction

单步执行单条CPU指令。别名`si`

### stepout

跳出当前函数

### thread

切换到指定进程，别名`tr`

`thread <id>`

### threads

打印所有的调用进程

### trace

设置tracepoint，`trace [name] <linespec>`，别名`t`

tracepoint是一种breakpoint，但是它不会停止程序的执行，而是产生一个提示信息。[linespec](https://github.com/go-delve/delve/tree/master/Documentation/cli/locspec.md)

### types

显示所有的类型`types [<regex>]`

### up

将当前frame向上移

```
up [<m>]
up [<m>] <command>
```

### vars

打印包的变量

`vars [-v] [<regex>]`