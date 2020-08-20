# SQL必知必会

关于SQL必知必会这本书的学习笔记。

## 第一章

主键 primary key：一列（或一组列），用来唯一标识表中的每一行。

## 第二章 检索数据

### 2.1 SELECT语句

例子：
```sql
SELECT product_name
FROM Products;
```

SELECT+列名+FROM+表名，这样的语句返回的是未排序的数据。

需要注意：在处理SQL语句的时候，所有的空格都可以被忽略。将SQL语句分成多行可以更便于进行调试。

### 2.5 检索不同的值

使用DISTINCT关键字可以让结果中只出现一个值，比如例子：
```sql
SELECT DISTINCT vend_id
FROM Products;
```

需要注意的是，如果选择了多余一列，则当且仅当两行所有的列的值相同的时候才会只显示其中一行。即不能部分使用DISTINCT。

### 2.6 限制结果

SELECT语句返回的是指定表中所有匹配的行，如果仅需要第一行或一定数量的行，则需要使用其他方法。

下文仅列出mysql的例子：
1. LIMIT关键字。

```sql
SELECT pod_name
FROM Products
LIMIT 5;
```

表示选择前5个。

2. OFFSET关键字，用来指定开始检索的行数（一般第一行的行数为0）

```sql
SELECT pod_name
FROM Products
LIMIT 5 OFFSET 5;
```

PS：MYSQL支持简化版本的LIMIT 5 OFFSET 4，即LIMIT 5,4。

### 2.7 使用注释

单行注释包括：
* 使用--（两个连字符），之后的文本即为注释。连字符前必须加上空格。
* 使用#，这种方法不一定得到支持。

多行注释为/**/，可以任意使用

## 第三章 排序检索数据

### 3.1 排序数据

使用SELECT检索出的数据并非随机显示的，一般以数据存在表中的顺序显示，也有可能是数据最初添加到表中的顺序。

1. **ORDER BY**关键字。例子：
```sql
SELECT prod_name
FROM Products
ORDER BY prod_name;
```

默认为升序，可以按照多个列顺序排序。如果对一个列名加上DESC关键字，则表示该列为降序排序。

排序的关键字必须在SELECT中出现。

```sql
SELECT prod_id,prod_prices,prod_name
FROM Products
ORDER BY prod_price DESC, prod_name;
```

## 第四章 过滤数据

### 4.1 使用WHERE字句

一般来说只需要获得部分数据，此时可以使用WHERE进行过滤

```sql
SELECT prod_name,prod_price
FROM Products
WHERE prod_price=3.18;
```

### 4.2 WHERE子句操作符
关键在于不等号：`<>`与`!=`，以及对空值NULL的判断。

* 范围值检查：`WHERE prod_price BETWEEN 5 AND 10`
* 空值检查：`WHERE prod_price IS NULL`

其中任何操作符的比较都不能列出空值。

## 第五章 高级数据过滤

### 5.1 组合WHERE子句

逻辑操作符：AND、OR，其中如果不用圆括号改变求值顺序，则默认先执行AND再执行OR。

### 5.2 IN操作符

表示数值是否为某几个值中的一个，例如：

```sql
SELECT prod_name,prod_price
FROM Products
WHERE vend_id IN ('DLL01','BRS01')
ORDER BY prod_name;
```

IN操作符与OR操作符比较类似，优点主要在于在多个条件同时存在的时候简单直观，同时不会改变求值顺序。

### 5.3 NOT操作符

## 第六章 用通配符过滤

### 6.1 LIKE操作符

在搜索子句中使用通配符必须使用LIKE操作符。此时使用通配符或搜索模式而不是简单匹配。

* 百分号%通配符：表示任何字符出现任意次数。比如以Fish开头可以表示为'Fish%'（在有些数据库中使用`*`）
* 下划线_通配符：与%一样，但是只能匹配单个字符。或者说，只能表示任意字符出现一次的情况。
* 方括号通配符([])，通常用来指定一个字符集，它必须匹配指定位置的一个字符。（数据库并不总是支持集合）

### 6.2 通配符的技巧

尽量避免将通配符置于开始，同时尽量减少使用通配符

## 第七章 创建计算字段

### 7.1 计算字段

存储在数据库中的字段一般不是程序所需要的数据，这时候需要直接从数据库中检索出转换、计算或格式化后的数据，而不是单纯检索出数据。

### 7.2 拼接字段

拼接，即将值联接在一起构成单个值。

操作符可以使用+或者||

例如：
```sql
SELECT vend_name+'('+ vend_country +')'
FROM Vendors
ORDER BY vend_name;
```

就可以得到一个新的列，按照规则组织拼接。
