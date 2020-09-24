# Enterprise applications cloud rightsizing through a joint becnmarking and optimization approach

Future Generation Computer Systems，2018，CCF C类

平台：基于VM的公有云平台。

解决的问题：对于指定的流量和指定的应用，选用哪一种类型的VM更合适（不同平台/不同规格）

作者在实验中提到，对于某一个应用在某一个流量级别上的表现而言，虚拟机越小其SE越大。这个原因是因为该流量级别不足以充分利用虚拟机的性能，所以应该选择容量最小的虚拟机。

## 摘要

本文所提出的背景是，将应用迁移到云平台需要考虑花费、性能和QoS的问题。即我们需要在满足用户QoS的同时找到花费最小的云配置方案。在这里的一个核心问题是需要知道应用在云平台上的行为（对各种资源的需求程度）来选择最佳的运行环境。

本文的工作包括：
* 云服务的性能评价
* 分类器/评测器机制，来确定一个随机服务的特性(footprint)，并根据此特性为随机服务提供其最适合的运行环境。（性能和代价）
* 一个设计空间探索工具，可以探索不同的云配置。

## 3 Cloud benchmarking

文章提出了一个SE(Service Efficiency)指标来描述benchmark的结果：

$SE=\frac{Clients}{w_1\times delay + w_2\times Cost}$

SE值越高越好

文章提到，为了更好地测量服务的性能方面，并且选择最适合该应用的服务，我们需要使用benchmarking的技术。

1. 第一步，基于不同的服务种类定义应用的性能基准(performance stereotype)。这些基准的目的是为了提取出性能指标中满足QoS所需要的特性。
2. 第二步，根据现实的流量模式，将应用的表现与这些特征相对应。

文章会用不同的benchmark数据集来与不同的应用类型相对应，从而得到应用关于CPU、IO、网络和数据相关性等特性。

## 4 Profiling and classification mechanism

对应用进行分类：
1. 第一步是通过相关的基准测试程序来测试出每个典型应用的性能最佳，然后使用性能分析工具对其进行度量。
2. 第二步是在VM下的相同环境运行这个程序，并测得23个特征作为footprint，用这些特征来训练classfication Tool。这样分类工具就可以根据特征来使用KNN算法，决定每个应用最适合的基准测试类型。

### 4.1 Profiling architecture

图3所示的Profiling Tool是在Linux系统下工作的。通过Pidstat和Tshark外部工具，可以监控虚拟环境下的进程并进行性能分析。（前者监控任务的资源消耗，后者通过IP监控网络）

Profiling工具有两种：Application Profiling和benchmark Profiling，这两种的区别还是很明显的。

### 4.2 Classification architecture

一旦application和benchmark的profiling都完成了，应用的开发者会使用Classification Tool来将应用组件映射到预定义的benchmark目录上，并且根据SE指标提供适合的云配置。

分类工具由下面三种组成：
* GUI，提供用户界面
* 分类器，使用KNN算法
* 控制器，通用的进程管理工具，负责与DB进行交互，并根据SE指标的计算来选择最佳的实例类型。

如此，用户可以完成随机应用的分类

## 5 QoS assessment and optimization

Design-time的探索是由SPACE4CLOUD工具支持的，它提供了对云应用的QoS性质的探索和评估。该应用允许云应用遵循模型驱动的开发方法来描述、分析和最优化自己。

该软件可以对应用的运行代价进行评估，通过一个给定的代价模型，从而达到最小化运行代价的目的。通过使用Layered Queuing Networks(LQNs)进行建模，系统可以对多线程环境下的云应用进行QoS分析。

软件的核心是一个Mixed Integer Linear Program，可以产生最初的部署。软件采用本地搜索方式（禁忌搜索或迭代本地搜索）来使得解空间向期望方向移动。最终结果是确定应用每一层的初始分布，并得到最佳分布。

之所以使用LQN的原因，是因为它能够很好地描述复杂系统（多层系统）。采用了24个LQN来对每天24小时进行分别建模。

## 6 Experimental analysis

