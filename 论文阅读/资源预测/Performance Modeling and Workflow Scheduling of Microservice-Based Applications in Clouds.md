# Performance Modeling and Workflow Scheduling of Microservice-Based Applications in Clouds

IEEE Transactions on Parallel and Distributed Systems,2019, CCF A类

第一眼看上去差点被摘要吓到，完全是我的想法的加强升级版。读完之后想了一下，好像跟我的方向也没什么太大关系。

本文要解决的问题是所谓的MAWS-BC问题，全称Microservice-based Application Workflow Scheduling problem for minimum end-to-end delay under a user-specified Budget Constraint。这个问题更多和工作流相关而不是和调度算法或者资源估计相关。

## 2 Related Work

本文的相关工作写得不错，其中的Performance Modeling and Prediction这里有很多论文我可能需要去看一下。从它写的相关工作可以看到，它引用的很多相关工作都是做execution time，也就是和本文相关的内容。

对于Microservice Architecture的一些内容也不错，我入门微服务也比较短，对于其中的一些基础知识也不太清楚，比如service compositions我也只是明白个大概。其中很多文章都出自[microservices.io](https://microservices.io/)这个网站，我觉得我有空应该多去阅读一下。

## 3 A System Overview

大概介绍了本文的建模对象，即一个泛用的工作流系统（见图3）。这个系统有一个统一的输入队列（消息队列），有一个业务层处理器，有一个事务层处理器（连接DB），并且使用多线程即使回复。

## 4 Performance Modeling and Prediction of a microservice

对Section 3提到的整体系统进行建模。这个模型其实还是挺泛用的，但就是太泛用了才让我感觉到不够实用，这种级别的建模注定无法落地于实际，我还是比较喜欢黑盒模型。

回到原文，作者对process time(下简称PT)进行了建模。PT说是处理时间，实际指的是请求从用户端产生服务端收到到从服务端发回的时间（就是服务端处理一个请求的总时间）

$PT=f_1(R)+f_2(C)+f_3(D)+\gamma$

* $f_1(R)$指的是请求缓存时间(request caching time)，即一个请求在输入队列中等待的时间，可以用排队论来建模。
* $f_2(C)$指的是处理时间（processing time），在文中指在business processing中处理的时间。在本文认为业务工作流的处理时间与输入参数的数量呈正相关，因此用输入参数的大小来估算处理时间的长度。
* $f_3(D)$指事务处理时间(transaction processing time)。本文认为IO事务的时间=逻辑IO时间（CPU调用请求）+DB的IO时间

### 4.2 Microservice Performance Prediction

作者将特征分为四个部分，即硬件特征、请求特征、业务特征和数据库特征（见Table 1）。目标是使用这些特征来估计微服务组件的执行时间。

## 5 Analytical models and problem formulation

使用了三层模型，见原图Fig.5.

* Application Model，使用工作流方式进行描述的DAG，每个结点为一个function，每个function由microservice组成。
* System Model，构建了具体实现是，每个function对应的microservice是如何被部署在VM中的。VM之间是全连接的。
* Performance Model，构建模型的最终目的是为了调整微服务，达到minimum end to end delay(MED)的目的。为此，需要明确每一个组件的execution time和monetary cost。

后面就没什么好说的了，它缩小端到端时延的方式就是缩小关键路径的时延，找到合适的微服务来实现功能，至少我的理解是这样的。这是一个服务选择的问题，还是个NP-hard的问题，因此用的是启发式算法。

## 7 Implementation and experimental results

结果说实话挺一般的。ROC曲线倒是画的不错，Fig7可以参考一下。但是结果真的是……我第一次看到响应时间是按ns来计算的，换算成秒的MSE能有10+，这误差也太大了吧。