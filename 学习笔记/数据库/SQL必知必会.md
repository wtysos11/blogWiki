# SQL必知必会

关于SQL必知必会这本书的学习笔记。

## 读后感

挺好的一本入门书籍，使用例子提供了对SQL的感性认识，与菜鸟教程配合用来极速入门非常的不错。

本文主要是对重点内容的摘录与一些笔记记录，而且只是会用级别，具体到性能优化需要再阅读其他书籍。

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
相等操作为`=`而非`==`

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

如SELECT X, Y AS B这样可以为计算列命一个别名

## 第8课

### 8.1 函数

[Mysql函数](https://www.runoob.com/mysql/mysql-functions.html)

如果选择使用函数，至少做好注释

### 8.2 使用函数

* 文本处理函数：比较常用的有去括号类、返回指定位置类、转换类、插入类等
* 日期和时间处理函数
* 数值处理函数

## 第九章 汇总数据

### 9.1 聚集函数

经常需要汇总数据而不是将其列举出来，比如：统计行数、统计某些行的和、找出列的最大/最小/平均值。

* AVG()，返回某列的平均值。AVG只能够用于得到一个列的平均值，如果需要多个列的数据需要调用多次。（会忽略NULL行）
* COUNT()，返回某列的行数。可以使用COUNT(*)对表中行的数目进行计算，不管是否为NULL；COUNT(COLUMN)可以对特定列中具有值的行进行计数，忽略NULL值。
* MAX()，返回某列的最大值。忽略NULL
* MIN()，返回某列的最小值。忽略NULL
* SUM()，返回某列之和。忽略NULL

### 9.2 聚集不同值

默认对所有行进行操作，即ALL。如果需要对不同数值进行操作需要加上DISTINCT

PS:
* DISTINCT不能用于COUNT(*)

### 9.3 组织聚集函数

可以同时使用多个聚集函数。

## 第十章 分组数据

GROUP BY子句与HAVING子句

### 10.1 数据分组

例子：需要对每一个供应商返回各自的产品

### 10.2 创建分组

```sql
SELECT vend_id,COUNT(*) AS num_prods
FROM Products
GROUP BY vend_id
```

上述可以统计不同供应商的供应商ID与供应商品数量。

GROUP BY 子句可以包含任意的列。如果进行了嵌套，将在最后一个分组上进行汇总。
位置上，GROUP BY子句必须出现在WHERE子句之后，ORDER BY子句之前。

### 10.3 过滤分组

SQL允许规定包含哪些分组，例如
```sql
SELECT cust_id,COUNT(*) AS orders
FROM Orders
GROUP BY cust_id
HAVING COUNT(*) >= 2;
```
最后一句话过滤哪些COUNT(*) >=2 的分组。可以看到，WHERE子句在这里不起作用。

PS：WHERE与HAVING的差别。WHERE在数据分组之前进行过滤，HAVING在数据分组之后进行过滤。WHERE过滤的行很可能会影响之后分组的结果。

因此在使用HAVING时应该与GROUP BY共用，并且要与WHERE区分开来。

### 10.4 分组和排序

GROUP BY和ORDER BY的区别

ORDER BY会对任意输出进行排序；GROUP BY只是对行分组，但输出可能并不是分组的顺序
ORDER BY可以对任意列（甚至非选择列）来使用；GROUP BY 只可能用于选择列或表达式列，并且必须使用每个表达式列。
ORDER BY如果不需要对数据排序时，不一定需要；但是GROUP BY如果与聚集函数一起使用的时候是一定要用的。

## 第十一章 使用子查询

目前的所有SELECT语句都是简单查询，即从单个数据库表中检索数据的单条语句。SQL允许创建子查询，即嵌套在其他查询中的查询。

### 11.2 利用子查询进行过滤

例如
```sql
SELECT cust_id
FROM Orders
WHERE order_num IN (SELECT order_num
    FROM OrderItems
    WHERE prod_id = 'RAGN01');
```

如此嵌套可以灵活实现很多功能，虽然受限于性能限制一般不会嵌套很多层。

并且，作为子查询的SELECT语句只能够查询单个列，如果试图查询多个列会返回错误。

### 11.3 作为计算字段使用子查询

使用子查询的另一个方式是使用计算字段，或者说这样子才比较常用。

PS：完全限定列名，即`WHERE ORders.cust_id = Customers.cust_id`，因为一个同样的列名很可能会出现在多个表中。

## 第十二章 联结表

### 12.1 联结

联接：join，利用SQL的SELECT能执行的最重要的操作（最影响性能）

为什么使用联接：当数据存在多个表中时，要正确取出数据，只能够考虑多重嵌套SELECT或join。而多重嵌套会带来很大的性能和编程负担，而且每次嵌套只能够向上返回一个列，不能适用很多情况，因此联结会更好。

### 12.2 创建联结

指定要联结的所有表与关联它们的方式即可

```sql
SELECT vend_name,prod_name,prod_price
FROM Vendors,Products
WHERE Vendors.vend_id = Products.vend_id;
```

此处语句将Vendors表与Products表中vend_id相同的行链接在一起，此处应该是保证两个表中vend_id各自唯一。

由没有联结条件的表关系返回的结果是笛卡尔积，因此一定要记得使用WHERE子句。

#### 内链接：INNER JOIN

```sql
SELECT vend_name,prod_name,prod_price
FROM Vendors INNER JOIN Products
ON Vendors.vend_id = Products.vend_id;
```

两个表之间的关系是INNER JOIN指定的部分FROM子句，联结条件用特定的ON而非WHERE给出（虽然实际上内容是相同的）

#### 联结多个表

SWL不限制一个SELECT语句中可以联结的表的个数。

```sql
SELECT prod_name,vend_name,prod_price,quantity
FROM Orderitems,Products,Vendors
WHERE Products.vend_id = Vendors.vend_id
AND OrderItems.prod_id = Products.prod_id
AND order_num = 20007;
```

上述例子显示的是订单号为20007中的物品，也可以使用嵌套查询。

## 第十三章 创建高级联结

### 13.1 使用表别名

SQL除了可以对列名和计算字段使用别名，还允许给表名起别名，理由：
1. 缩短SQL语句
2. 允许在一条SELECT语句中多次使用相同的表。

例子：
``` sql
SELECT cust_name,cust_contact
FROM Customers AS C, Orders AS O, OrderItems AS OI
WHERE C.cust_id = O.cust_id
AND OI.order_num = O.order_num
AND prod_id = 'RGAN01';
```

### 13.2 使用不同类型的联结

前面提到的是内联接或者等值联结的简单联结，下面介绍三种：自联结self-join、自然联结natural join和外联结outer join

#### 13.2.1 自联结

例子：查询要求找到Jim Jones工作的公司，然后找出在该公司工作的顾客的信息。

第一种想法是嵌套子查询

```sql
SELECT cust_id,cust_name,cust_contact
FROM Customers
WHERE cust_name = (SELECT cust_name
    FROM Customers
    WHERE cust_contact = 'Jim Jones');
```

第二种想法是使用自联结
```sql
SELECT c1.cust_id,c1.cust_name,c1.cust_contact
FROM Customers AS c1, Customers AS c2
WHERE c1.cust_name = c2.cust_name
AND c2.cust_contact = 'Jim Jones';
```

PS：一般而言，DBMS处理自联结的速度要快于嵌套子句，实际中需要具体计算来考虑到底使用哪一种。

#### 13.2.2 自然联结

自然联结可以排除两个相同的列，使每个列只出现一次。实际上就是一种特殊的内联结，要求比较中的分量必须是相同的属性组，并将重复的属性组去掉。系统不自动完成，需要用户来指定，如

```sql
SELECT C.*,O.order_num,O.order_data,
    OI.prod_id,OI.quantity,OI.item_price
FROM Customers AS C, Orders AS O, OrderItems AS OI
WHERE C.cust_id = O.cust_id
AND OI.order_num = O.order_num
AND prod_id = 'RGAN01';
```

因为通配符只对第一个表使用，所以其他列明确列出，没有多余的列被检索出来。目前的每个内联结都是自然联结，很可能永远也不会用到不是自然联结的内联结。

#### 13.2.3 外联结

许多联结需要将一个表中的行与另一个表中的行进行关联，但有时候需要包含没有关联行的哪些行，比如

要检索没有订单顾客在内的所有顾客

```sql
SELECT Customers.cust_id,Orders.order_num
FROM Customers LEFT OUTER JOIN Orders
ON Customers.cust_id = Orders.cust_id;
```

外联结是一个比较复杂的内容，这里不进行细究。详情可以查看[mysql学习记录](/学习笔记/数据库/mysql学习记录.md)中知识点记录的第二个知识点：外联结