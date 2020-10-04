# Auto-Scaling Web Applications in Clouds : A Taxonomy and Survey

标签：

## 前言

云计算时一种新的提供计算资源的模式，特点：

* subscription-oriented services
* pay-as-you-go basis

## 2 Problem Definition and challenges

问题定义和挑战

定义：在单一云上，弹性调度问题可以被定义为如何自动和动态地调整一系列资源来满足一个波动的应用负载，从而在满足资源消耗最小的情况下满足SLA或SLO。

是一个经典的自动控制问题，需要系统自动决定分配资源的类型和数量，从而达到一定的性能要求。

经典模型：MAPE控制循环

## 4 Application Architectures

可以参考一下多层应用，不过由于现阶段大部分多层应用都是基于排队网络来做的，因此没有什么好的想法。

SOA架构与微服务架构其实是很类似的，可以参考一下，虽然肯定是不一样的。但在网络架构上还是很类似地，即不想多层应用一样是线性的连接（可以参考腾讯的DAGOR）

## 5 Session Stickness

我们目前还没有考虑到Session的问题，可能之后考虑到链接持续性的时候会用到这一部分。

## 6 Adaptivity

自动调度系统是否是自适应的，即可以自行调节参数

## 7 Scaling Indicators

将调度的触发器分为底层和高层指标。底层指标包括CPU利用率、内存利用率等，是通用的；高层指标是服务请求时间和请求混合值等普通的VM很难获得的指标。

目前的工作对这方面的探索可能并不多，如果有一些能够很好地获取高层指标的方法也是可以借鉴一下的。

## 8 Resource estimation

资源估计的目的就是决定当前流量下需要的最小计算资源。

目前来说有六个组：rule-based, fuzzy inference, application profiling, analytical modeling, machine learning, and hybrid approaches

## 9 Osillcation mitigation

Oscillation为波动、震荡，即自动调度器反复重复增加/删除应用。这个不在我们的考虑范围之内。

## 10 Scaling timing

决定调度的时间点，这个可以考虑一下，因为它一般与预测式调度相关

## 11 Scaling methods

是使用垂直调度还是水平调度，一般来说仅考虑水平调度

## 13 额外考虑

我个人是认为SOA中的很多知识与微服务是相同的，在不考虑架构细节的情况下，都是多入口多出口的模式。

一些对隐藏参数的监控工具也是值得探索的，比如如果突然间有一种工具能够探究应用的服务时间，那么弹性调度器一下子就可以工作的更加顺畅。

资源预测模型很多都是有一定历史的，存在很多的误差，因此如果能够提出更精确的资源预测模型也是能发一篇文章的。


