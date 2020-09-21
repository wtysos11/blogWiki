# Predicting cloud performance for HPC applications before deployment

Future Generation Computer Systems，2018，CCF C类

大致内容是针对高性能集群HPC应用系统，提出了一种基于机器学习的方法，用于预测应用程序在云中部署前的运行性能

## 摘要

为了减轻投入，现在很多HPC用户都选择将应用迁移到云上。在应用上云的过程中，用户可能：a.不能了解应用与云系统软件层的交互。b.对云系统的一些硬件细节缺乏了解。c.无法了解云系统与他人的共享资源是如何带来性能损失的。这些误解可能导致用户选择了代价上或性能上的次优解。

在本文，我们采用了一种机器学习的方式来支持用户选择针对指定目标流量在云上部署时的最优配置。这使用户可以在面对迁移和分析云中的应用程序的成本之前决定是否购买以及购买什么。我们在云服务提供商这一角度建立了云服务性能预测模型（cloud performance prediction model,CP）以及对应的用户端下一个硬件独立的profile prediction model(PP)。PP能够捕捉应用级别的调度行为，用户在一个小的数据集上进行profile，然后在一个大的数据集上使用机器学习算法建立PP。CP的产生是由云服务提供商进行，来学习硬件独立的Profile和云性能直接的联系。

## 3 Proposed methodology

* CP, cloud-performance-prediction model，是云服务提供商方面产生的一个模型。CP是产生来推导在执行指定流量时计算集群的细节。
* PP，profile-prediction model，是从用户角度产生的模型。PP是依据硬件独立的Profile建立的模型。

将这两个模型联结起来可以在应用软件执行之前就能够预测性能的可用性。

### Provider side

首先，搜集所有的要训练的应用和它们的数据集中的cloud configuration。然后对于每一个要训练的应用a和数据集d，沃恩会搜集一个与硬件无关的profile,p(a,d)。然后，通过在数据集合d，配置c下运行需要测试的应用a，可以得到性能指标$\zeta$。

最后，使用一个机器学习方法来学习一个近似函数$\hat{\zeta}(p,c) \sim \zeta(a,d,c)$。本文选用的机器学习方法是随机森林，因为它可以很简单地被使用，并证明了在特定问题中具有高准确度。(G. Mariani, A. Anghel, R. Jongerius, G. Dittmann, Predicting cloud performance for HPC applications: A user-oriented approach, in: Proceedings of the 17th IEEE/ACM International Symposium on Cluster, Cloud and Grid Computing, CCGrid ’17, IEEE Press, Piscataway, NJ, USA, 2017, pp. 524–533. http://dx.doi. org/10.1109/CCGRID.2017.11. https://doi.org/10.1109/CCGRID.2017.11. )

### User side

PP在用户访问云基础设施之前产生，使用机器学习方法来近似这个训练的profile p(d)，使用一个analytical function $\hat{p}(d)$。在本文中我们使用LR和ANN来实现PP。

### Coupling CP 和 PP

CP需要捕获云系统的行为，因此它需要符合实际训练应用在云系统中的执行性能指标。CP的输入是一个cloud configuration和硬件独立的profile。

PP学习用户指定目标的伸缩行为，这个应用还没有在CP的训练中出现过，它被用来在部署前于用户端训练PP。所有的训练数据需要去产生PP，用户可以使用PP来预测未知的数据。

最终的目标是用户可以选择一个最佳的云配置来在云中处理目标数据集。

构建应用profile所需要的feature:
* Operation mix：各种操作的比例(integer,floating point,memory read/write等)
* ILP：instruction level parallelism对于一个理想的机器执行LLVM IR指令。每一个操作类型都会进行计算。
* Reuse distance：对于不同的距离我们会计算概率，用来作为内存复用的估算。
* Library calls：对于external library中每个函数调用量的估计。
* communication requirements：对于communication library中的每个函数调用量的估计。