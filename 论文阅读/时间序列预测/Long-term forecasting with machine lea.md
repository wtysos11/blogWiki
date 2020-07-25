# Long-term forecasting with machine learning models

标签：计算机文章 时间序列预测

[原文](https://thuijskens.github.io/2016/08/03/time-series-forecasting/)，是2016年的一篇网络博文。

不使用原始算法的原因是因为数据量大（>500K）
1. 传统的ARIMA对于短期预测表现很好，但长期预测时AR部分难以收敛。
2. MCMC sampling algorithm对于一些贝叶斯状态空间模型来说计算代价很大。

因此本文采用机器学习方法。然而机器学习方法大多是为了independent and identically distributted(IID)数据设计的，因此将其用在非IID数据上是很有趣的。

## 预测策略

本文遵循non-linear autoregressive representation (NAR) assumption

对于$y_{t+1}=f(y_t,...,y_{t-n+1}+\epsilno_t$，即对过去的连续n个时间点，建立一个未知函数模型f，直接预测未来t->t+H的时间点的数据(而不只是第t+H个点的数据)。

当H=1的时候是单步预测，此时很容易进行。当H>1的时候有三种方法：
* iterated one-step ahead forecasting
* direct H-step ahead forecasting
* multiple input multiple output models

### iterated forecasting

在这种模式下，我们先基于one-step ahead criterion建立一个模型，然后在预测的时候迭代向前推进。

这个模式的问题在于，在后面的估计中实际上使用估计值作为预测的输入，因此会很容易受到预测误差的影响，十分的不稳定。

### Direct H-step ahead forecasting

直接预测H步以后的时间点，因此不会受到累计误差的影响。然而并不一定存在这样的关系，因此最终训练出来的模型质量是无法保证的。

本文主要引用了A review and comparison of strategies for multi-step ahead time series forecasting based on the NN5 forecasting competition此文的想法。里面提到的模型选择和模型平衡问题十分值得注意。