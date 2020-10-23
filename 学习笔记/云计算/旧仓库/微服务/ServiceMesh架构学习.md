# Service Mesh架构学习

标签：ServiceMesh 微服务 Kubernetes

[资料来源](http://www.servicemesher.com/blog/service-mesh-architectures/


如何使用Service Mesh

* 在软件中引用相关的库。库方法不需要与底层基础架构进行太多的合作，但如果要实现多语言支持，就必须用不同的语言去重复实现多次。
* 节点代理：在此架构中，每个节点上运行着一个单独的代理，为异构的服务提供负载。它不关心应用程序的语言，可以为许多不同的微服务租户提供服务。Linkerd使用的是这种方案。
* sidecar。这是Istio与Envoy使用的模型，为每个应用的容器部署一个伴生容器，对于Service Mesh，sidecar接管进出应用程序容器的所有网络流量。这种方法介于库和节点代理模型之间。Sidecar利于工作审计，特别是一些与安全相关的方面。