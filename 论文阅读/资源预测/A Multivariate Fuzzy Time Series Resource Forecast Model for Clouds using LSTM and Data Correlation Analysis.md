# A Multivariate Fuzzy Time Series Resource Forecast Model for Clouds using LSTM and Data Correlation Analysis

标签：资源预测 云计算 

2018 Procedia Computer Science （CCF未列入），期刊影响因子2.5，不属于SCI。

针对的是虚拟机环境，以资源使用率为指标的固定阈值法的问题。比如设定内存利用率80%，但是调度需要时间，在完成调度的时候内存利用率可能仍然大于80%或小于80%，造成资源的浪费与无效操作。

使用的数据集：[google workload](https://www.pdl.cmu.edu/PDL-FTP/CloudComputing/googletrace-socc2012.pdf)，并没有直接给出数据，而是给出了一篇论文。

做法：
1. 使用fuzzy-technology来预处理数据
2. 考虑到相关性来选择指标。
3. 