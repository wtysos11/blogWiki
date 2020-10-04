# Research on Auto-Scaling of Web Applications in Cloud: Survey, Trends and Future Directions

标签：计算机文献 微服务 综述

2019 Scalable Computing: Practice and Experience，不在CCF和中科院分区上的垃圾期刊

综述做的也很一般

## 摘要

Cloud computing emerging environment attracts many applications providers to deploy web applications on cloud data centers. The primary area of attraction is elasticity, which allows to auto-scale the resources on-demand. However, web applications usually have dynamic workload and hard to predict. Cloud service providers and researchers are working to reduce the cost while maintaining the Quality of Service (QoS). One of the key challenges for web application in cloud computing is auto-scaling. The auto-scaling in cloud computing is still in infancy and required detail investigation of taxonomy, approach and types of resources mapped to the current research. In this article, we presented the literature survey for auto-scaling techniques of web applications in cloud computing. This survey supports the research community to find the requirements in auto-scaling techniques. We present a taxonomy of reviewed articles with parameters such as auto-scaling techniques, approach, resources, monitoring tool, experiment, workload, and metric, etc. Based on the analysis, we proposed the new areas of research in this direction.

本文主要面对的是多层结构的网络应用的弹性扩缩。

## 2 Auto scaling

IBM提出的MAPE-K循环。

特殊状态:
* Over-provision，即分配了过多的资源
* Oscillation，抖动。

## 3 Taxonomy of Auto-scaling

弹性调度的分类

## 4 Elastic Application

大部分云应用都是天然有着elastic特性的。

从架构上来说，本文主要关注于multi-tier应用。

## Survey on Auto-scaling techniques

分为了七类：
1. Application Profiling
2. threshold-based rules
3. fuzzy rule
4. Control theory
5. queuing theory
6. machine learning
7. time series analysis

### 5.1 Application profiling

Application profiling是一个找到资源利用率顶点的过程。用来测试这个点的流量可以是真实的或者模拟的，可以在线或是离线进行测试（我理解的是生产环境），是最简单的在不同时间点进行资源分配的一种方式。

有方法可以动态地测试出应用的类别进行分类（采用决策树方式）。

或者是对多层应用的每一层进行测试，从而分配充足的资源。

可以在不加入虚拟机的情况下离线对每一层应用性能进行测试，从而实现快速分配。

* Self-managing cloud-native applications: Design, implementation, and experience. 2017, Future Generation Computer Systems.
* Pascal: An architecture for proactive auto-scaling for distributed services. 2019, Future Generation Computer Systems.

Off-line application profiling，这种方式可以在每次应用更新后离线进行，包括；
* Integer Linear Programming，数学最优化问题，所有的变量都是整数，约束和目标函数都是线性的。
* Workload Profiling Technique，这种技术会收集负载的信息，对负载进行建模

Online Application Profiling，动态进行来满足需求，用的比较多
* Application Signatures：确定一个小的测试流量集合，对应用进行分类。
* Elastic distributed resource scaling(AGILE)：预测式方法，在产生SLA违约前进行调度。
* Rapid Profiling，对虚拟机实例进行测试，判断需要被分到哪一层

### 5.2 Threshold-based rules

最经典的方法

### 5.3 Fuzzy Rules

一类rule-based的方式，被单独归类。rules被定义为if-else条件，优点时可以使用语言条件，比如low,medium,high来控制，用control-system来作为根据指定资源在指定流量下性能的估算。

### 4.3 Control Theory

排队论部分倒是写得挺好的。
