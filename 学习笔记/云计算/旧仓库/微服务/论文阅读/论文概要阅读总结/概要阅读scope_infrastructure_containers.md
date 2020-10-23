# 概要阅读scope_infrastructure_containers

标签：论文阅读 概要阅读 containers

来自综述：elasticity in cloud computing: state of the art and research challenges

目标是概要阅读大量的摘要来大概了解该领域文章的内容，同时选出几篇来进行精读。

## auto-scaling of micro-services using containerization

### 基本信息

* 会议/期刊：International Journal of Science and Research IJSR
* 作者：Priyanka P. Kukade, Prof. Geetanjali Kale
* 机构：Pune Institute of Computer technology, Dhankawadi, Pune - 411043, India

### 摘要内容

使用了一种方法，将传统的平台服务拆分成微服务。这些微服务可以松耦合，并且独立地进行部署。

## DoCloud: An Elastic Cloud Platform for WEb Applications Based on Docker

看过了，杨泓臻上次讲的论文。

## Model-Driven Management of Docker containers

### 基本信息

* 会议/期刊：2016 IEEE 9th International Conference on Cloud Computing
* 作者：Fawaz Paraiso, Stephanie Challita, Yahya Al-Dhuraibi, Philippe Merle
* 机构：University of Lille & Inria Lille - Nord Europe CRIStAL UMR CNRS 9189,France

### 摘要

docker缺少设计时部署检查工具(Docker lacks of deployability verification tool for containers at design time)，docker也没有提供一种工具来进行容器资源的按需变更。而且容器比虚拟机更加底层，需要用户关注更底层的问题。

关注容器云的管理问题，开发了一种方法来制作docker容器模型，并开发了一种工具来确保容器的部署和管理。

### 没看前的问题

这东西不是kubernetes要做的事情吗

## MultiBox: Lightweight Containers for Vendor-Independent Multi-cloud Deployments

### 基本信息

* Communications in Computer and Information Science 系列
* 作者：James Hadley, Yehia Elkhatib, Gordon Blair, and Utz Roedig
* 机构： Lancaster University, Lancaster LA1 4WA, UK

### 摘要

云上虚拟机间容器的迁移问题。（kubernetes干的）