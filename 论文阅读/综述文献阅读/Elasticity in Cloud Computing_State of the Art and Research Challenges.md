# Elasticity in Cloud Computing : State of the Art and Research Challenges

标签：计算机文献 微服务 综述

引用格式：
Al-Dhuraibi Y, Paraiso F, Djarallah N, et al. (2017). Elasticity in Cloud Computing : State of the Art and Research Challenges. IEEE Transactions on Services Computing, 11(2), 430–447. https://doi.org/10.1109/TSC.2017.2711009

下面分为几个部分：
* 摘要：摘录原文的摘要并进行翻译，大致了解这篇文章是讲什么的。
* 概念：理清原文所介绍的各种概念，这也是阅读综述类文献的意义。
* 总结：进行总结

## 摘要

Elasticity is a fundamental property in cloud computing that has recently witnessed major developments. This article reviews both classical and recent elasticity solutions and provides an overview of containerization, a new technological trend in lightweight virtualization. It also discusses major issues and research challenges related to elasticity in cloud computing. We comprehensively review and analyze the proposals developed in this field. We provide a taxonomy of elasticity mechanisms according to the identified works and key properties. Compared to other works in literature, this article presents a broader and detailed analysis of elasticity approaches and is considered as the first survey addressing the elasticity of containers.

## 概念

### Elasticity

Elasticity is one of the key features in cloud computing that dynamically adjusts the amount of allocated resources to meet changes in workload demands.

是从资源角度出发定义的。弹性指的是调度器能够动态调整资源使其满足负载需求的能力。

We define elasticity as the ability of a system to add and remove resources (such as CPU cores, memory, VM and container instances) "on the fly" to adapt to the load variation in real time.

可以分为水平弹性和垂直弹性两种。

相近的概念有scalability和efficiency。Elasticity=scalability+automation+optimization，scalability指的是系统能够根据增长的流量动态分配额外资源的能力。

efficiency值的是云服务提供商能如何有效地利用资源，以使用的资源总量作为衡量标准（the amount of resources consumed for processing a given amount of work）。一般来说，在完成相同任务时，消耗的资源越少，则有效性越高。

### 状态

* Over-Provision：分配的资源超出需要的资源，尽管QoS被满足了，但是会导致额外和不必要的资源使用。
* Under-Provision：分配的资源小于需要的资源，会导致服务退化和SLA违约。
* just-in-need：代表一种平衡状态，负载可以被适当地处理，QoS可以被保证。

### 模式

手工调度(manual mode)

自动调度(all actions are done automatically)

* 响应式调度(reactive mode)：弹性动作被基于确定阈值或规则触发
* 预测式调度(proactive mode)：这种方法的实现依靠预测技术，需要明确未来的需求。

此外还有将两者结合起来的混合式调度方法。

### 架构

中心化：只有一个调度器
去中心化：有多个调度器组成。

Multi-Agent System：分布式调度系统。基于此的系统可以称为agent-based computing，更加轻量化、更加自动化。

### Virtualization

虚拟化是云计算的基石，云服务提供商通过虚拟化技术构建环境。虚拟化技术包括hypervisor技术和container-based virtualization。

### 

## 总结
