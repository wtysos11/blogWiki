# Web Application Resource Requirements Estimation based on the Workload Latent Features

IEEE transactions on services computing, 2019 ,CCF B类

我个人是觉得这篇文章是非常重要的。

本文提出过去的预测式调度方法通常都只基于流量的输入值，而在相同流量输入值的情况下，CPU、内存占用、响应时间和网络带宽的值都可能是不一样的。因此本文更多关注于流量的隐特征，在本文中，流量的隐特征指的是请求的概率分布。

## 3 Proposed System methodology

在本阶段完成的是URI空间的映射。本文在前面分析了输入流量不同对于资源的压力是不一样的（见原文Fig1）。

### 3.1 URI-Space Partitioning

进行URI空间切分的目的是为了将需求相似资源的URI放在一类。本文将service time和document size作为聚类的重要指标，最终会分为k类（指定k）。分类时假定需要的资源在每一类中都呈正态分布（Assuming that the required resources, including CPU, memory, network and I/O are sampled randomly from a normal distribution for every URI in the subspace）

本文用来一个代价函数，聚类的目的是最小化这个代价函数，其实本质上还是一个k-means的过程。k从2逐渐往上增加，直到代价函数（描述error rate）变得最小。

这个函数太复杂了，我就不记录了，原文式(1)

### 3.2 Workload Latent Features Identification

上面一节将URI空间划分为K个子空间，下面会将读取历史上的日志，并将一个时间间隔内的所有日志进行聚合，将其到K个子空间的概率计算出来，所得到的概率即为隐特征（看原文Fig2会更清楚一点）

## 4 Resource Demand estimation

对比算法列了很多，Kriging Model(我有点没看明白，这不是个插值算法吗)、Ridge Regression、LASSO Regression

本文主打的算法是MLP，结构见Fig3。其实是很基本的想法，和我与师兄做的那个类似。每个ANN为两个隐层，首先各用一个ANN来预测CPU、内存和带宽，然后用三个+隐向量来预测响应时间。作者认为这样的好处就是训练的比较彻底，实际上应该差不多，因为ANN的梯度消失还是比较明显的，如果不这样做的话需要的数据量会大很多。

（但是我们之前实践的时候ANN是很不准确的，我们还是只用了单个URI做测试，最后是用强化学习强行把误差也考虑了进去，效果最后倒是还不错）

## 5 Experiment setup and design

测试用了两个benchmark：RUBIS和Acme Air（这两个都挺常用的……后面那个我见过微服务版本的，不知道前面的有没有）

实验中除了用隐特征向量之外，还使用了单纯的请求到达率作为baseline（这个就有点欺负人了）

从结果上来看，隐特征向量对于几乎所有的方法提升都很明显，并且在大部分情况下MLP效果都不错。

虽然图看着挺好的，但是实际指标上预测误差，特别是响应时间的预测误差还是很大的。而且对比指标中用的是MSE，而不是NRMSE这种比例性质的指标，因此不同数值之间也是有影响的，我觉得不够直观。

## 总结

总体而言是很不错的一篇文章，这个隐特征向量真的是有点东西。但是基于文档大小和响应时间来进行聚类有些不妥当，原因讲得不够让我信服。尽管效果很好，但我比较怀疑这种聚类方式的泛用性。