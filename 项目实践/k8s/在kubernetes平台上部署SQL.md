# 在kubernetes平台上部署SQL

在华为云的云容器引擎上部署SQL服务器，供之后的学习。

参考文献：
* [运行一个单实例有状态应用](https://kubernetes.io/zh/docs/tasks/run-application/run-single-instance-stateful-application/)
* [PV和PVC](https://www.cnblogs.com/along21/p/10342788.html)

## 相关知识

在有状态数据库的部署中，有几个问题，第一：如何挂在硬盘；第二：如何挂载Secret(不使用直接写入yaml参数的方式)；第三：如何暴露使其能够被访问。

* PV：PersistentVolume，指的是与pod相独立的存储资源，可以理解为容量插件。
* PVC:PersistentVolume Claim，是用户进行存储的请求，类似于pod。Pod消耗节点资源，PVC消耗PV资源。PVC和PV是一一对应的。

在生产中一般使用Storage Class动态存储，不过这里就用简单的就行了。

SELECT *  FROM(
    SELECT  * FROM `bigquery-public-data.wikipedia.pageviews_2015` UNION ALL
    SELECT  * FROM `bigquery-public-data.wikipedia.pageviews_2016` UNION ALL
)
WHERE views=0 LIMIT 10