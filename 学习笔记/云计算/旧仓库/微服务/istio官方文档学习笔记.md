# istio官方文档学习笔记

标签：kubernetes istio 学习笔记

## 架构

Istio架构上分为数据平面和控制平面。

* 数据平面由一组以sidecar方式部署的智能代理(Envoy)组成。这些代理用以调节和控制微服务以及Mixer之间所有的网络通信。
* 控制平面负责管理和配置代理来路由流量。以此控制平面配置Mixer以实施策略和收集遥测数据。

一些概念：

* Envoy，C++开发的高性能代理，以sidecar的方式和同一个服务被部署在同一个k8s Pod中。
* Mixer，一个独立于平台的组件，负责在服务网格上执行访问控制和使用策略，并从Envoy代理和其他服务收集遥测数据。在高层次上，Mixer好；会将istio的后端服务器全部封装起来，并将自己作为一个中间件使用。
* Pilot，为Envoy sidecar提供服务发现功能，为智能路由和弹性提供流量管理功能。它将控制流量行为的高级路由规则转换为特定于Envoy的配置，并在运行时把它们传播到sidecar。
* Citadel，通过内置身份和凭证管理提供强大的服务间和最终用户身份验证。

## 流量管理

Istio流量管理的核心组件是Pilot，它管理和配置部署在特定Istio服务网格中的所有Envoy代理实例，允许指定在Envoy代理之间使用的是什么样的路由流量规则。

Istio的流量管理模型，本质上是将流量与基础设施解耦，让运维人员通过Pilot指定流量遵循什么样的规则，而不是执行哪些pod/VM应该接受流量。

### 服务发现和负载均衡

* 服务注册：使用k8s的服务注册表实现服务发现，服务的新实例会自动注册到服务注册表，并且不健康的实例将会被自动删除。
* 服务发现：Pilot使用来自服务注册的信息，并提供与平台无关的服务发现接口。

Istio仅支持三种负载平衡模式：轮询、随机和带权重的最少请求。

### 故障处理

Envoy提供了一套开箱即用的，可选的故障恢复功能：

1. 超时。
2. 具备超时预算，能够在重试之间进行可变抖动的有限重试功能。
3. 并发连接数和上游服务请求数限制
4. 对负载均衡池中的每个成员进行主动运行健康检查
5. 细粒度熔断器（被动健康检查）

### VirtualService

VirtualService表示如何对服务的请求进行路由控制。

#### 规则的目标描述

路由规则对应着一个或多个VirtualService配置指定的请求目的主机，称为`hosts`条目。条目内可以是内部的名称，或是域名。

hosts字段用显示或隐式的方式定义了一个或多个完全限定名（FQDN）。

#### 根据来源或Header制定规则

1. 根据制定用户进行限定。可以制定一个规则，只针对来自reviews服务的Pod生效：

```yaml
apiVersion: networking.istio.io/v1alpha3
kind: VirtualService
metadata:
  name: ratings
spec:
  hosts:
  - ratings
  http:
  - match:
      sourceLabels:
        app: reviews
    ...
```

sourceLabels依赖于服务的实现，在k8s中，与Pod的选择标签是一致的。

2. 根据调用方的特定版本进行限定。

```yaml
apiVersion: networking.istio.io/v1alpha3
kind: VirtualService
metadata:
  name: ratings
spec:
  hosts:
  - ratings
  http:
  - match:
    - sourceLabels:
        app: reviews
        version: v2
    ...
```

3. 根据HTTP Header选择规则。下面的规则只会对包含了end-user头，且值为jason的请求生效。

```yaml
apiVersion: networking.istio.io/v1alpha3
kind: VirtualService
metadata:
  name: reviews
spec:
  hosts:
    - reviews
  http:
  - match:
    - headers:
        end-user:
          exact: jason
    ...
```

多个header之间是与的关系。

#### 目标规则

在请求被VirtualService路由之后，Destination配置的一系列策略就生效了。

#### Gateway

Gateway为HTTP/TCP流量配置了一个负载均衡，多数情况下在网格边缘进行操作，用于启动一个服务的ingress流量。