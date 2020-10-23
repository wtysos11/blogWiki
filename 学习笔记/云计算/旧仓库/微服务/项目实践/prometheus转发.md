# promethus转发

标签：promethus istio

promethus是非常重要的流量监控工具，但是它的浏览器界面只能够在内网下打开，如何在外网访问这个好用的流量监控工具呢

根据[bookinfo遥测教学](https://istio.io/zh/docs/tasks/telemetry/metrics/collecting-metrics/)，在本机上访问[链接地址](http://localhost:9090/graph#%5B%7B%22range_input%22%3A%221h%22%2C%22expr%22%3A%22istio_double_request_count%22%2C%22tab%22%3A1%7D%5D)即可打开prometheus界面，我们的目标就是让外网也能够打开这个界面，这就需要建立外网到内网9090端口的映射。

## 直接修改svc/istio-ingressgateway

在ports中加入自己写的一段，将外部端口31399映射到内部的9090

```yaml
  - name: prometheus
    nodePort: 31399
    port: 9090
    protocol: TCP
    targetPort: 9090
```

结果自然是被拒绝了，果然没有这么简单。

## 加入自己的入口网关

遥测教学的做法是先在后台运行这个`kubectl -n istio-system port-forward $(kubectl -n istio-system get pod -l app=prometheus -o jsonpath='{.items[0].metadata.name}') 9090:9090`，将本地的9090端口与目标容器的9090端口作端口转发，然后访问本地的9090端口来访问。好处是避免了pod的ip变动造成的失效，但是外部很显然访问不到这个服务。

自己仿照bookinfo的例子写了个网关，参考了[API document of VirtualService](https://istio.io/docs/reference/config/networking/v1alpha3/virtual-service/)，目标是将外部端口9090访问转发到内部prometheus的9090端口上。

```yaml
apiVersion: networking.istio.io/v1alpha3
kind: Gateway
metadata:
  name: prometheus-gateway
spec:
  selector:
    istio: ingressgateway # use istio default controller
  servers:
  - port:
      number: 90
      name: http
      protocol: HTTP
    hosts:
    - "*"
---
apiVersion: networking.istio.io/v1alpha3
kind: VirtualService
metadata:
  name: prometheus
spec:
  hosts:
  - "*"
  gateways:
  - prometheus-gateway
  http:
  - match:
    - uri:
        regex: "*"
    route:
    - destination:
        host: prometheus
        port:
          number: 9090

```

看了下istio的入口网关，发现内部的9090对应的是外部的31399。而且看了一下之前博客里面记录下来的数据，这个是实时添加上去的。很好奇它的实现方式，找时间拆掉。

```
istio-ingressgateway     LoadBalancer   10.109.139.108   129.204.7.185   15020:30937/TCP,80:31380/TCP,9090:31399/TCP,443:31390/TCP,31400:31400/TCP,15029:30721/TCP,15030:31100/TCP,15031:31773/TCP,15032:30877/TCP,15443:32541/TCP   6h41m
```

理论上来说访问`129.204.7.185:31399/productpage`就可以访问到prometheus的服务了，大概。

但实际上试了好久好像都不行。


## 改变prometheus的服务

在网上看到了一个issue，说是在istio环境下如何访问prometheus时，[有人](https://github.com/istio/istio/issues/6652)提到可以将prometheus的svc换成LoadBalancer，然后从外网ip访问

修改如下：
1. 将type从ClusterIP换成了LoadBalancer
2. 修改了最后的status，原本时空对象`{}`，现在为

```
  loadBalancer:
    ingress:
    - ip: 203.195.219.185
```

结果发现配置时成功了，但是LoadBalancer的external-IP pending，说明没有落实。应该按照之前[博客的文章](https://blog.csdn.net/u012837895/article/details/89052503)对ingress controller进行修改。

使用`kubectl get configmaps -nmetallb-system`找到了原来配置的那个configmap，接下来就是将ip加入ip池中。

改完之后大概长这个样子

```yaml
# Please edit the object below. Lines beginning with a '#' will be ignored,
# and an empty file will abort the edit. If an error occurs while saving this file will be
# reopened with the relevant failures.
#
apiVersion: v1
data:
  config: |
    address-pools:
    - name: default
      protocol: layer2
      addresses:
      - 129.204.7.185-129.204.7.185
      - 203.195.219.185-203.195.219.185
kind: ConfigMap
metadata:
  annotations:
    kubectl.kubernetes.io/last-applied-configuration: |
      {"apiVersion":"v1","data":{"config":"address-pools:\n- name: default\n  protocol: layer2\n  addresses:\n  - 129.204.7.185-129.204.7.185\n"},"kind":"ConfigMap","metadata":{"annotations":{},"name":"config","namespace":"metallb-system"}}
  creationTimestamp: "2019-04-06T02:01:03Z"
  name: config
  namespace: metallb-system
  resourceVersion: "6869478"
  selfLink: /api/v1/namespaces/metallb-system/configmaps/config
  uid: d526d2a2-580f-11e9-96eb-5254001218ec
```

改完之后发现externalIP果然显示出来了，效果下：

```
prometheus               LoadBalancer   10.102.237.167   203.195.219.185   9090:31588/TCP                                                                                                                                              7h8m
```

理论上来说，访问外网ip`203.195.219.185:31588`就可以访问到prometheus了。

按照官方文档的说法，可以访问`203.195.219.185:31588/graph`访问图像，访问`203.195.219.185:31588/metrics`访问流量数据。然而并没有成功。

## 参考文章再次改写ingreess

参考了[这篇文章](https://www.qikqiak.com/post/kubernetes-monitor-prometheus-grafana/)，作者使用ingress将prometheus暴露到外网。不过原文似乎不是在istio环境下跑的，我应该考虑把自动注入关掉。

首先将自动注入关闭`kubectl label namespace default istio-injection=disabled --overwrite`，希望我记得还有这件事，把它改回来。

