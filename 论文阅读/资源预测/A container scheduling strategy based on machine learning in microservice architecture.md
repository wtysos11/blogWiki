# A container scheduling strategy based on machine learning in microservice architecture

IEEE Conference on Services Computing， 2019， CCF C类

居然是我们学校的……

作者运用total access, concurrency, average response time of users, error rate采用随机森林回归的方式来预测微服务架构下所需要的容器数量。比较吸引我的是微服务架构，虽然我也没有看出它哪里微服务了（指不存在调用关系）

尽管作者自己给了三个不同的流量曲线，并且说明是使用JMeter进行测试流量的发送，但是也没有说明是对哪一个应用进行预测（而且我估计这个应用只有一个容器）【我问了一下，是istio的测试应用bookinfo，将整个应用作为单个应用进行扩缩。居然还能做到这么高的准确率】

没什么好说的，比较恐怖的一点是用随机森林的准确性实在高的有点离谱，我回去再看一下。