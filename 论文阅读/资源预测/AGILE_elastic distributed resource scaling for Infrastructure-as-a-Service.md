# AGILE:elastic distributed resource scaling for Infrastructure-as-a-Service

2013 IEEE ICAC

## 简介

本文是虚拟机环境下的调度算法，在multi-tier架构下（一般e-commerce），用wavelet decomposition来实现中程预测，并使用RUBiS和Google cluster traces来验证预测式算法。

本文的一大特色是用户可以指定SLO违约率，从而对调度结果产生影响。

AGILE可以使用online profiling和polynomial curve fitting技术，来提供一个黑盒性能模型，根据给定的资源压力预测未来的SLO违约速率，同时根据环境变化实时更新。

重点贡献；
* resource pressure model，来根据给定的SLO违约率确定资源数量。
* 模型建立是依赖于RUBiS系统的，有可能在其他系统上不能够正常运行。

对于超时情况的预测，有3.42倍的正例和0.34倍的负例，更加精确。

## 2 AGILE system design

### 2.1 基于小波变换的中期资源预测

系统采用的是预测式的设计，中等长度，指的是两分钟（即VM的启动时间，60个时间点，每个时间点占两秒）

采用滑动窗口法，滑动窗口大小D=6000s，并且不需要离线的application profiling以及一些白盒/灰盒应用建模方法。

思路：
1. 采用小波变换进行时间序列分解(Haar小波)
2. 对每个分解的时间序列进行预测
3. 组合预测数据，估算资源消耗。

与FFT的异同：FFT只能用正弦系列函数，因此只对存在周期性的数据表现良好。Wavelet则可以更好地处理acyclic pattern。

小波变化需要配置的参数有二，其一为选择使用的小波函数，其二为scale的数量。因为模拟信号会导致更好的数值被预测，我们需要最大化approximation signal和original signal之间的相似性，因此选择使得两者欧氏距离最小的配置组合。

预测方式采用Markov model based prediction，可能其他方式会更好。

### 2.2 Online resource pressure modeling

AGILE的目的就是为了管理资源分配使得能够满足应用的SLO需求。一种方式是预测未来的输入流量，并通过建立模型来推演出所需要的资源数量，从而完成输入流量到所需要资源的映射从而满足SLO。然而这种方式需要对应用非常详尽的了解（比如白盒级别），在IaaS级别很难进行精确预测，而且很难泛化。

因此我们所使用的方式是预测应用的资源使用率，使用application-agnostic resource pressure model来将应用的SLO违约率目标（比如小于5%）映射到需要维持的最大资源压力。其中，资源压力(Resource pressure)是被分配的资源所使用的比例。注意我们显然需要分配更多一点的资源来应对workload spke以及预留一些空间。

这个映射模型使用online profiling进行来对应用的每一层生成一个资源压力模型。（资源压力->SLO违约率，目前难以理解的一点是SLO违约率是怎么计算出来的）

显然，这个resource pressure model是应用特化的（甚至每一层都不同），而且因为workload mix的改变，可能会在运行时发生变化。（比如在RUBiS中，写需求更多的时候需求更多的CPU）

为了应对这个，采用了online profilng和curve fitting：
1. 为了构建新的函数，第一步肯定是要收集数据。通过调整应用资源的分配量，如果是多层应用则在每一层分别动态调整（其他层要分配足够的资源来确保不是瓶颈层）。如果SLO被多种资源所影响，则一个个来测试。然后对同层的所有应用结果取平均。（profiling interval一般取一分钟）
2. 然后使用不同阶的polynomial来拟合profiling data的曲线，并选出least-square error的那一条曲线，最大阶数为16来避免过拟合。当运行时函数发生巨大改变且拟合误差超出预定阈值时(5%)，则AGILE会用新的模型来替代旧的模型。同时为了避免频繁地重新训练，AGILE会对模型进行缓存，在对于有着不同流量模式的应用来说比较合适。

## 3 Exp

评价指标使用RUBiS online auction benchmark在Apache Xassandra key-value store上进行测试，同时使用Google cluster data验证预测算法。

在overload prediction上，使用的指标为准确率和召回率。