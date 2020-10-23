# delve快速入门

标签：golang delve 调试 工具学习

安装：`go get -u github.com/go-delve/delve/cmd/dlv`

[官方仓库](https://github.com/go-delve/delve)

命令主体：`dlv debug xxx`

输入后即可进入调试界面

## gdb

其实gdb就可以做到很多强大的功能，推荐[这篇文章](https://www.oschina.net/translate/using-gdb-debugger-with-go)

## 调试界面命令

### 设置断点

比如使用`break main.main`可以在package main的main函数头设置断点

可以使用"文件：行号"来打断点，`b /home/goworkspace/src/github.com/mytest/main.go:20`，也可以使用相对位置。

### 运行到断点

`continue`，或者简写`c`可以运行到断点。

### 单步执行

`step`，或者简写`s`，不过它会直接进入到函数中。

### 单步执行

`next`，或者简写`n`，下一行，不进入函数
