# Survey on prediction models of applications for resources provisioning in the cloud
标签：计算机论文 资源阅读

发表于Journal of Network and Computer Application 2017,CCF 计算机网络领域C类期刊

内容：总结了目前最新的application prediction methods

## 2 Application Prediction

* 定义：在多角度预测未来一段时间应用的行为，比如它所承受的负载以及表现。
* 意义：应用预测是在云上进行有效资源分配的必要一步

不同类型的应用对QoS的要求不同，文中表1给出。对于网站应用而言，QoS要求为存储可靠性、高网络带宽，以及高可靠性。（高可靠性可以理解为低响应时间）

### 2.1 Different dimensions of prediction

通常来说，应用预测的目标是为了从不同角度描述一段时间后应用的行为，基于已经收集的数据。

在虚拟机层次上，大部分文章着眼于预测流量、表现和SLA参数。

workload在不同文章中表现不同：
* The number of requests
* resource demand
* resource utilization

performance:
* throughput
* response time

### 2.2 Characteristics, challenges, needs and evaluation metrics

#### 2.2.1 The needs for application prediction

进行应用预测的原因：
1. Application Management：在云环境中，应用之间互相有资源上的共享，因此他们的资源需求需要被合适地预测，即分配给应用的资源应当于其流量成比例。这里会尽量避免QoS下降
2. Resource/Cost Management：为了最大化资源利用率。

#### 2.2.2 Prediction characteristics

对预测的要求：
* Accurate：预测模型可以用预测结果的准确性来评估。
* Adaptability(online learning)：云环境经常在变动，如果预测模型应该能够很好地适应这种变动。随着时间的推移，模型应该能够增进对应用行为的了解并降低预测误差。
* proactive：VM的创建和迁移是需要大量时间的，因此预测需要提前进行，在流量发生变动前即完成预测。
* Historic Data：不同种类的数据被分配给了云服务，要利用这些不同的数据。

#### 2.2.3 Prediction challenges

* Complexity：预测模型所消耗的资源不能显著高于应用资源，不然得不偿失。
* Data granularity：初始阶段要决定去监控哪一个资源，下一步要决定监控的区间。区间不能过小也不能过大。
* Pattern Length：在大多数预测模型中，模式长度是固定的。该如何确定这个擦汗高难度

### 2.3 Evaluation metrics

* Cost：由SLA违约或资源浪费所带来的误差
* Success：成功指标决定预测模型能够多精准地预测应用未来的行为
* Profit：用于计算云服务提供商的profit-rate
* Error：用来评估真实行为和预测行为的差别。

## 3 Modeling approaches for the application prediction 

### 3.1 Table driven methods

在本方法中，应用的行为被记录在表中，用不同的负载密度值和不同的资源数来代表。

很老的方法，可扩展性很低。

### 3.2 Control theory

在控制论中，目标是控制云应用之间的资源分配。

* 如果模型控制一种资源，则是单进单出SISO
* 如果模型控制多种资源，则是多进多出MIMO

### 3.3 Queuing theory

排队论使用排队网络来预测应用的性能，通过对负载和性能之间关系的建模来进行。在排队论中，每个服务器的应用都在一个排队系统内，任务从一个队列流转到另一个队列。

### 3.4 Machine Learning techniques

最新的方法是机器学习方法，可以在多个维度预测应用的行为。

可以用来预测资源行为、SLA违约、应用表现和任务的执行时间。

机器学习方法将应用的行为看作是一个时间序列，绝大多数方法都是基于一个长度为m的滑动窗口进行预测。

强化学习是另外一种更新的方法，可以从动态地环境中学习最佳的策略。同时机器学习和统计预测非常接近，一般看作是相同方法。

## 4 Overview of prediction methods

* 大部分方法提供一步预测，部分方法可以进行多步预测。多步预测中，预测准确性会随着步长增加而下降。
* 有一些方法会采用聚类算法，根据资源种类、任务和VM类型进行聚类。将相似对象聚为一类能够利用更多的数据，提取出对象之间的相关性，并得到更精确的结果。类似的技术还有主成分分析可以得到时间和空间上的缩减。
* 有一些方法能够平滑应用的时间序列，提供更精确的结果。
* 马尔科夫模型等能够基于系统之前的行为学习到趋势的模型，预测结果可以提供更多的信息。
* 模糊算法提供了一种新的预测方法。进化算法能够提供相似的解。

### 4.1 Methods based on machine learning and statistical techniques

#### 4.1.1 Regression and moving average

AR、ARMA算法

模型简单，可靠性基于对流量的简单假设。这些方法基于过去信息进行训练，因此无法动态捕获应用行为的变化。因此在误差增加的时候模型需要被重新训练。

#### 4.1.2 神经网络

神经网络不需要对应用流量进行假设，可以对非线性行为建模。超参数较多，很复杂，而且是黑盒模型，不能提供关于行为模式的解释。

#### 4.1.3 Marrkov model, clustering, dimension reduction and fuzzy logic

Markov model可以快速从历史数据中计算得到。然而，其成立的假设非常严格，不一定适用于所有的流量。无法用于长期预测，需要经常根据最新的观测值重建模型。

聚类算法常用于资源、任务和VM的种类，基于聚类的预测方法会更加节省时间和空间，同时可以捕获物体与自身的相关性。

fuzzy logic模型可以产生不确定的数据，模拟了人类认知，因此结果更具有解释性。但是需要先验知识。

#### 4.1.4 Hurst exponent and bayesian theory

Hurst exponent非常简单，决定了应用行为是否可预测。可以基于历史数据来提供未来一段时间的趋势。但不是在所有应用下都是可计算的，可以与其他方法一起来提高预测精度。

Bayesian theory很简单，可以将初始信息作为先验知识导入预测器中。如果资源管理者无法提供足够的先验知识，它会用过去的行为进行估计。这要求描述行为的特征需要是独立的，而且先验概率经常需要被重新计算。

#### Histogram, probability distribution and benchmarking

尽管很简单，且计算代价固定，但是结果不一定可靠。

#### Mathematical derivation and string matching

灰度预测模型、多项式预测、卡尔曼滤波器、KMP等

优点是需要很少甚至不需要训练数据，计算速度快，可以根据应用负载快速变化。但参数选择比较具有挑战性。

#### 4.1.7 Combination of support vector machine, neural network and regression

相结合的方法可以提供多步预测，但是结果会更好。缺点仍然是黑盒模型，而且需要经常重复训练，SVM训练需要较长时间。

#### 4.1.8 Filtering and signal processing

用滤波器平滑时间序列，但是可能会消减掉有用的信息。

信号处理方法，包括FFT和小波变换，可以提取周期性信息，没有对负载有严格的前置假设，但都会丢失细节信息。

#### 4.1.9 Workload generator and workload factoring

