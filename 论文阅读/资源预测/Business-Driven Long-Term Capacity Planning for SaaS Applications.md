# Business-Driven Long-Term Capacity Planning for SaaS Applications

目标：使得SaaS provider的利润最大化

设想的场景：一个SaaS服务商向IaaS服务商拿取资源。云服务有按需付费（on-demand）和预付费(reservation)两种，前者昂贵、不可获得但是灵活，后者便宜、保证获得但是只能按时购买。因此如何选择预付费资源的购买数量和时间是很重要的。

## 4 Capacity Planning Heuristics

采取了两种容量计划的启发式方法:
* 一种基于instance utilization(UT)
* 一种基于Queueing Theory

两种都接受一个未来时间段D中流量的预测结果。D即为计划的时间（长时间，比如一年）