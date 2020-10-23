# Python中调用C

标签：python 链接

感谢大佬们在stackoverflow上的[问题](https://stackoverflow.com/questions/145270/calling-c-c-from-python)中的回复。

## 官方支持

[官方文档](https://docs.python.org/3.7/extending/extending.html)

作用是通过特别定义的C文件，调用`Python.h`，得以在C中使用Python的数据或调用Python的函数，并在Python中直接调用C的函数或使用C的数据。实现方式应该是动态链接，将编译产生的\*.o文件利用上。

但是过程相当的繁琐而且底层，而且说实话我真的没看懂。

## ctypes

[官方文档](https://docs.python.org/3.7/library/ctypes.html)

挺有名的一种调用C代码的方式，是标准库里面提供的。需要先将C编译成动态链接库的形式，即windows下的.dll文件，以及linux下的.so文件

ctypes作为python和C联系的桥梁，定义了一系列的专用数据类型来连接这两种编程语言，详见ctypes type，其中值得注意的是c_long和c_int在long和int字长相同的时候是可以互换的。这些属性的具体值可以通过`value`属性进行访问。如果需要可变的内存块，可以使用`create_string_buffer()`。

* 动态链接库中的函数会作为python对象对应的属性进行调用。如果函数名字不是有效的python属性名，像是`??2@YAPAXI@Z`，则可以使用`getattr`方法来获取，比如`getattr(cdll.msvcrt, "??2@YAPAXI@Z")`

总体而言，ctypes虽然比上一个简单，但还是相当的繁琐。

## swig

[中文文档](http://www.swig.org/translations/chinese/index.html)，版本有点老，最近的一次更新是2011年的时候，不过应该还可以。

## Boost.Python

## Cython

[官方文档](https://cython.org/)，优点是很多科学计算包会带上它，像是我用的anaconda。
[中文介绍](https://www.jianshu.com/p/fc5025094912)


Cython代码包括.pyx文件和.c文件。其中.pyx文件会被Cython编译成.c文件（内含Python扩展模块），.c文件会被C编译⑦；编译成.so文件（windows上为.pyd文件），可以直接被python所调用。