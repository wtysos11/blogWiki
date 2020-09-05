# bigquery官方教程

标签：bigquery 数据库

[官方教程地址](https://cloud.google.com/bigquery/docs/how-to)

本文为对bigquery的学习总结与笔记，主要记录一些容易忘记的知识点。部分重点知识点会附带上原文记录。原文应该是需要翻墙的（或者说谷歌的网站有哪个是能直接连上的）

## 知识点

* [控制费用](https://cloud.google.com/bigquery/docs/best-practices-costs)，几乎是必看的文章，如果不想在bq上花过多的钱。简单来说，就是要认识到bq是列式存储，按检索到的列来收费。如果可能，尽量避免`SELECT *`这样的语句，因为它会检索所有的列。同时`LIMIT`对节省花费并没有任何帮助。
