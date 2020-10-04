# PASCAL_An architecture for proactive auto-scaling of distributed services

摘要：在分布式系统中使用的调度方法，但总体而言还是十分接近我们的调度目标。它会有一个响应时间阈值上界的要求。最终的目标就是在响应时间没有超过上界的情况下，最小化自身的资源数。

总体而言分为两个阶段，第一个阶段是profiling phase，包括完成workload model和performance model，然后在运行时使用，预测未来的输入流量（workload model），并且计算所需要的最小资源配置(performance model)，并触发调度。