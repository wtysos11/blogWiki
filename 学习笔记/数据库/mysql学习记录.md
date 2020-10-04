# mysql学习记录

标签：mysql 文章 记录

速查速记，不涉及原理

## 快速QA

## 命令行操作记录

## 知识点记录

### 1. 什么时候用\`，什么时候用'。

参考：[SO](https://stackoverflow.com/questions/42319049/where-do-we-use-backticks-and-quotes-in-sql)。重音号用于包含数据库、表名、列名。除非需要使用复杂字符，不然一般可以不用重音号。


        SELECT * FROM `database`.`table` WHERE `column` = "value";`

    而引号一般用于表示字符串。

### 2. 外联结


