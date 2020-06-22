## python下ARIMA方法学习

标签：python ARIMA

### 简介

使用python学习ARIMA，参考数据：[result.csv](https://github.com/wtysos11/WikipediaDataPredict/blob/master/result.csv)维基百科总访问量数据，小时粒度。

参考资料：
* 教程[How to Create an ARIMA Model for Time Series Forecasting in Python](https://machinelearningmastery.com/arima-for-time-series-forecasting-with-python/)
* [中文教程](https://blog.csdn.net/shine19930820/article/details/72667656)
* 使用的[statsmodel-arima](https://www.statsmodels.org/stable/generated/statsmodels.tsa.arima_model.ARIMA.html)

目标：
1. 阅读后能够掌握使用python进行ARIMA预测的代码能力，最好能够有cheat sheet
2. 可以考虑对ARIMA的原理进行描述，在超参数的选择上要有所侧重

### cheat sheet



### 原理介绍

ARIMA的全称

* AR
* I
* MA

参数：
* p
* d
* q

### 例子

确定参数

单步预测

滚动前向预测



师兄，今天我讲的时候突然发现我好像把指标算错了。我脑子有点晕你帮我看一下。

我的逻辑是可以按照真实流量T_real的值分配资源，使得响应时间一定不会超过上界，比如每个实例50req/s时保证不会超时。这样如果按照预测流量值乘以某个倍数可以保证满足响应时间。