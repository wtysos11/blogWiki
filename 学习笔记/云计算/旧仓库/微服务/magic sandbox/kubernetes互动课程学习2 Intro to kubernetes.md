# kubernetes互动课程学习之2

标签：kubernetes 实践

因为内容很多[前文](https://github.com/wtysos11/NoteBook/blob/master/%E5%BE%AE%E6%9C%8D%E5%8A%A1/magic%20sandbox/kubernetes%E4%BA%92%E5%8A%A8%E8%AF%BE%E7%A8%8B%E5%AD%A6%E4%B9%A0.md)都学过了一遍，而且概念性的东西都比较简单，就不赘述了。

## 2.1 Kubernetes简介

### k8s特性

这是从整体的角度来分析k8s的优点，或者说，特性。

Kubernetes是在容器层进行操作而不是硬件层。因此只要能够在容器上部署的应用，就可以在Kubernetes上部署。

Kubernetes提供了很多PaaS都有的特性，如部署、弹性伸缩、负载均衡、日志和监控，但这些并不是统一的，而是可选的插件。

* automatic binpacking：自动将容器根据资源消耗和其他限制条件，在不牺牲可用性的前提下部署到节点上。尽可能提高资源的使用率。
* self-healing：在容器失败的时候重启、重新调度。
* horizontal scaling：水平拉伸。拉伸分为垂直和水平两种，垂直指的是分配给实例的资源增加，水平指的是分配更多的实例。
* service discovery and load balancing：服务发现和负载均衡。如[前文](https://github.com/wtysos11/NoteBook/blob/master/%E5%BE%AE%E6%9C%8D%E5%8A%A1/magic%20sandbox/kubernetes%E4%BA%92%E5%8A%A8%E8%AF%BE%E7%A8%8B%E5%AD%A6%E4%B9%A0.md)提到的，k8s中的服务发现主要是service进行的，而负载均衡是由kube-proxy完成的。
* automated rollouts and rollbacks：k8s会渐进地roll out来更新应用，如果发生错误，会roll back，以避免服务掉线。
* secret and configuration management。可以部署和升级Secret对象来应对服务器中需要加密的数据，而不需要重构镜像，或者说在配置文件中暴露。
* storage orchestration：自动挂载存储系统，支持公有云服务商的文件存储系统，或者网络存储系统。
* batch execution：除了服务以外，k8s还可以管理批处理脚本和持续集成(CI,continuous integration)的负载压力，替换失败的容器。

### k8s做不到什么

* 自动部署源代码和编译，CI/CD需要用户自行完成。
* 不支持应用级的服务，比如中间件、数据处理框架、数据库等，需要用户自己封装完成。
* 不提供日志、监控或者警告模块的解决方案，虽然有很多的第三方插件。
* 不提供configuration language，提供了声明式的API。

## 2.2 Kubernetes Architecture

### Kubernetes的组件

K8s可以分为3部分：Master组件、Node组件、Addons

Master组件提供了cluster的控制平面，进行关于cluster的全局决策（比如调度），检测cluster的活动并作出回应。Master组件可以在cluster的任意机器上运行，但为了简便处理，一般只在一台机器上运行，并且不再这台机器上分配容器，管这台机器叫做Master节点（虽然实际上也可以让它分配Pod）

Master组件包括：

* kube-apiserver：暴露k8s的API，控制平面的前端。
* etcd：key-value型的存储器，使用raft一致性算法，存放着所有cluster中的数据。（记得备份）
* kube-scheduler：负责调度Pod到节点上，调度操作的执行者。具体可以见[这篇文章](https://github.com/wtysos11/NoteBook/blob/master/%E5%BE%AE%E6%9C%8D%E5%8A%A1/Kubernetes%E5%AD%A6%E4%B9%A0%E7%AC%94%E8%AE%B0.md)
* kube-controller-manager：负责运行controller的组件。Controller会从api server监听容器的当前状态，并试图将它变为用户的期望状态。原则上，每个controller是分离的进程，但为了减少复杂性，它们被编译进了一个二进制文件中，并在一个进程上运行。Controller包括Node controller(负责提示和报告Node下线)、Replication Controller(副本管理器，负责维护Pod的副本数量)、Endpoints Controller（产生Endpoint对象，连接Service和Pod）、Service Account & Token Controller（为新的命名空间创建账户和API访问令牌）
* cloud-controller-manager：负责运行和云服务商相关的controller，通过调用它们提供的API来完成行为。

Node组件在每一个节点上运行，维护Pod的运行和k8s的运行时环境。Node可以是虚拟机或是物理机。

* kubelet：在cluster上每个节点都运行着的客户端，负责维护Pod中容器的运行。kubelet拿到一个PodSpecs的集合，确保这个集合中的容器健康运行。
* kube-proxy：kube-proxy具体实现service的功能，详见[前文](https://github.com/wtysos11/NoteBook/blob/master/%E5%BE%AE%E6%9C%8D%E5%8A%A1/magic%20sandbox/kubernetes%E4%BA%92%E5%8A%A8%E8%AF%BE%E7%A8%8B%E5%AD%A6%E4%B9%A0.md)，通过维护转发表来实现，iptables在大流量下表现不好。
* container runtime：容器运行时，负责运行容器的软件，目前支持:docker,rkt,runc和所有符合OCI接口的实现。

Addons是实现其他非必须cluster特性的pod和service，例如DNS服务、互联网看板、资源监控和日志框架等等。

## 2.3 interacting with kubernetes

一般来说，用户通过使用k8s的API对象来描述他所希望的集群状态。

k8s通过使用k8s object来完成系统状态的抽象。一个k8s对象是目的的记录，一旦你创建了一个对象，k8s系统会尽量去维护它。通过创建这个对象，你告诉k8s集群你所希望的cluster状态。

教程中使用了RESTful API来与api server进行交互，但是我几次尝试都失败了，不知道原因。

实际上，官方提供的kubectl就是通过RESTful API和api server进行交互的，就像docker也是通过RESTful API和docker守护进程进行交互一样。

* `kubectl apply`，通过`-f`参数可以执行指定配置文件来配置资源，支持JSON和yaml.
* `kubectl get`，通过kubectl get resources name拿到资源列表，了解资源的状态
* `kubectl describe`，通过kubectl describe name拿到详细信息。
* `kubectl cluster-info`，查看集群信息
* 

get和describe支持的参数基本一致，包括：

* -n,--namespace，指定命名空间。
* `resource type/name`，通过下划线加名字可以指定特定名字的资源
* -l参数，指明label，后跟一串key、value作为筛选。

## 2.4 kubectl command

[kubectl官方文档](https://kubernetes.io/docs/reference/generated/kubectl/kubectl-commands)

[kubectl操作简介](https://kubernetes.io/docs/reference/kubectl/overview/#operations)

通用命令格式：`kubectl [command] [TYPE] [NAME] [flags]`

* command指想要进行的状态，如get、describe、delete、create等。
* type和name指定了操作的目标对象。其中type为资源类型，比如pod、service等，name为实例的名字属性。type是大小写不敏感的，而且容错率较高。[支持的type类型](https://kubernetes.io/docs/reference/kubectl/overview/#resource-types)。如果忽略name的话，会显示type的所有实例。
* 如果想要同时显示多个对象，可以使用`TYPE1 name1 name2 name<N>`，即在后面罗列姓名。如果类型不一样，可以使用斜杠，如`TYPE1/name1 TYPE1/name2 TYPE<#>/name<#>`

## 2.5 YAML file structure

参考：[kubectl API reference doc](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.11/)

* apiVersion：Kubernetes会使用多个版本的API，但一般都使用v1，如果不确定可以看[这里](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.11/)
* kind:指定创建的对象类型
* metadata：包含了可以唯一识别一个对象的数据，比如名字和命名空间。比较重要的属性还有labels，一般用来帮助给对象归类。
* spec：对每个kind的对象来说这里都不一样，指定了创建的细节，

可以将多个对象放入同一个yaml文件中，每个对象之间换行，使用`---`进行分隔。

## 2.6 kubernetes self-healing mechanisms

复习一下Deployment的定义：为Pods和ReplicaSets提供声明式升级服务的controller。

ReplicationController负责根据Deployment的定义维护Pod的副本数，在副本数量不对的时候进行调整，同时提供滚动升级服务以及回滚。

比较重要的概念：系统间的连接是并行进行的，label是唯一确定连接的重要机制。