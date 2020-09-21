# Service agreement trifecta: Backup resources, price and penalty in the availability-aware cloud

2018 information systems research，这个期刊很垃圾，中科院管理类三区。

从云数据中心角度出发，号称是第一个建立资源分配和付费的统一模型。一般吧，好像没什么用。

我没有理解错的话，作者的想法应该是最小化备用资源。里面涉及到了一些contract design的工作，即用户与云服务提供商之间如何达成SLA协议。一般来说，contract design literature工作主要关注IT business value和underlying risks之间的平衡，大部分时候是要对用户的不确定需求进行量化。

作者提到ITIC 2017的调查显示，80%的受访者目前愿意接受的最小商业服务器可用性目前是99.99%。显然，对SLA更加严格的要求使得在给定资源数的情况下SLA违约的可能性变大了，这就令对SLA要求严格的用户需要支付更多的成本，因为云服务提供商需要准备更多的后备资源来应对可能出现的服务失败情况；而SLA要求更宽泛的用户需要支付更少的成本。

用户和云服务提供商之间有着它们各自的cost-benefit模型来进行谈判，并且有着不同的决策目标。因此SLA中具体值的确定（违约代价、具体指标）的上界和下界都是动态博弈的结果，上界取决于双方的议价能力，下界取决于双方的成本。

## 相关工作

1. contract design literature，关注于如何量化用户的不确定需求，从而平衡IT商业花费与潜在风险，协助SLA的建立与合同的签署。
2. Optimal resource allocation，根据SLA或QoS的要求来分配最佳资源，一般是cost-aware的策略或分析。