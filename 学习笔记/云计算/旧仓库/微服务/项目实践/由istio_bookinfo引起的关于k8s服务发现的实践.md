# k8s服务发现

标签：kubernetes 服务发现 istio

参考资料：

* 阿里云写的这一篇[文章](https://yq.aliyun.com/articles/636511)，这里将istio的gateway和kubernetes ingress controller进行了对比，但是并没有提到怎么样在外网上访问，不过对比还是很有价值的，展示了istio的gateway的一些特性。

## bookinfo实例

实例来自[istio中文网](https://preliminary.istio.io/zh/docs/examples/bookinfo/)

执行命令：`kubectl apply -f samples/bookinfo/platform/kube/bookinfo.yaml`，如果此时所有的svc和pods都跑起来了，可以通过运行`kubectl exec -it $(kubectl get pod -l app=ratings -o jsonpath='{.items[0].metadata.name}') -c ratings -- curl productpage:9080/productpage | grep -o "<title>.*</title>"`来检测到结果。

这个结果颠覆了我对所谓的集群内的认识，我一直以为在跑k8s的集群的虚拟机内叫集群内，其实是k8s虚拟机的pod中才算是集群内，k8s外统一为集群外。所以要想在clusterIP下访问服务，只有在Pod内才能够做到。

其中`kubectl get pod -l app=ratings -o jsonpath='{.items[0].metadata.name}'`返回的是rating的pod名称。整个命令的意思是进入到rating所在pod中，执行`curl productpage:9080/productpage | grep -o "<title>.*</title>"`这条命令。因为`productpage`的ClusterIP开放的端口是9080，所以可以访问成功。

然后确定Ingress的IP和端口，执行`kubectl apply -f samples/bookinfo/networking/bookinfo-gateway.yaml`定义入口网关。

执行`kubectl apply -f destination-rule-all.yaml`定义destinationRule

确定端口：

```bash
export INGRESS_PORT=$(kubectl -n istio-system get service istio-ingressgateway -o jsonpath='{.spec.ports[?(@.name=="http2")].nodePort}')
export SECURE_INGRESS_PORT=$(kubectl -n istio-system get service istio-ingressgateway -o jsonpath='{.spec.ports[?(@.name=="https")].nodePort}')
```

确定IP：

```bash
export INGRESS_HOST=$(kubectl get po -l istio=ingressgateway -n istio-system -o 'jsonpath={.items[0].status.hostIP}')
```

设置GATEWAY_URL：`export GATEWAY_URL=$INGRESS_HOST:$INGRESS_PORT`，这时候直接`curl $GATEWAY_URL/productpage`就可以访问到目标页面了。

## 外网访问

现在bookinfo成功搭建起来了，但是服务器内网可以通过内网IP访问到，但是外网该如何访问呢？

现在的IP情况如下：

```
[root@k8s-master install]# kubectl get svc
NAME          TYPE        CLUSTER-IP      EXTERNAL-IP   PORT(S)    AGE
details       ClusterIP   10.100.82.10    <none>        9080/TCP   34m
kubernetes    ClusterIP   10.96.0.1       <none>        443/TCP    47d
productpage   ClusterIP   10.104.214.24   <none>        9080/TCP   34m
ratings       ClusterIP   10.109.4.97     <none>        9080/TCP   34m
reviews       ClusterIP   10.105.125.8    <none>        9080/TCP   34m
```

其中`$GATEWAY_URL=10.104.223.158:31380`，且能够正常访问到目标地址。但是腾讯云主机的外部地址是`129.204.7.185`等地址，显然，这些外部地址是无法和内部地址相对应的。

我记得是宋净超大佬的[从外部访问kubernetes中的Pod](https://jimmysong.io/posts/accessing-kubernetes-pods-from-outside-of-the-cluster/)一文中有介绍各种操作，但是我不是很确定istio中是否能够使用

下面参考了b站上华为云发布的[istio教程](https://www.bilibili.com/video/av40758305?from=search&seid=2717202429341329089)的第一篇，感觉很多概念一下子清楚了很多，推荐新人朋友看一下。

首先要认识到，所有的istio流量都是走istio-ingressgateway进入istio集群里面的，这个服务如下：

```
NAME                   TYPE           CLUSTER-IP       EXTERNAL-IP     PORT(S)                                                                                                                                      AGE
istio-ingressgateway   LoadBalancer   10.109.139.108   129.204.7.185   15020:30937/TCP,80:31380/TCP,443:31390/TCP,31400:31400/TCP,15029:30721/TCP,15030:31100/TCP,15031:31773/TCP,15032:30877/TCP,15443:32541/TCP   82m
```

它的yaml如下图(这个LoadBalancer的IP是由[这篇文章](https://blog.csdn.net/u012837895/article/details/89052503)里的ingress-controller进行的分配)，只要在外网访问这个IP加上对应的NodePort的端口，就可以访问到服务。

```yaml
apiVersion: v1
kind: Service
metadata:
  annotations:
    kubectl.kubernetes.io/last-applied-configuration: |
      {"apiVersion":"v1","kind":"Service","metadata":{"annotations":{},"labels":{"app":"istio-ingressgateway","chart":"gateways","heritage":"Tiller","istio":"ingressgateway","release":"istio"},"name":"istio-ingressgateway","namespace":"istio-system"},"spec":{"ports":[{"name":"status-port","port":15020,"targetPort":15020},{"name":"http2","nodePort":31380,"port":80,"targetPort":80},{"name":"https","nodePort":31390,"port":443},{"name":"tcp","nodePort":31400,"port":31400},{"name":"https-kiali","port":15029,"targetPort":15029},{"name":"https-prometheus","port":15030,"targetPort":15030},{"name":"https-grafana","port":15031,"targetPort":15031},{"name":"https-tracing","port":15032,"targetPort":15032},{"name":"tls","port":15443,"targetPort":15443}],"selector":{"app":"istio-ingressgateway","istio":"ingressgateway","release":"istio"},"type":"LoadBalancer"}}
  creationTimestamp: "2019-04-23T01:30:47Z"
  labels:
    app: istio-ingressgateway
    chart: gateways
    heritage: Tiller
    istio: ingressgateway
    release: istio
  name: istio-ingressgateway
  namespace: istio-system
  resourceVersion: "6828596"
  selfLink: /api/v1/namespaces/istio-system/services/istio-ingressgateway
  uid: 6bf0a4b3-6567-11e9-96eb-5254001218ec
spec:
  clusterIP: 10.109.139.108
  externalTrafficPolicy: Cluster
  ports:
  - name: status-port
    nodePort: 30937
    port: 15020
    protocol: TCP
    targetPort: 15020
  - name: http2
    nodePort: 31380
    port: 80
    protocol: TCP
    targetPort: 80
  - name: https
    nodePort: 31390
    port: 443
    protocol: TCP
    targetPort: 443
  - name: tcp
    nodePort: 31400
    port: 31400
    protocol: TCP
    targetPort: 31400
  - name: https-kiali
    nodePort: 30721
    port: 15029
    protocol: TCP
    targetPort: 15029
  - name: https-prometheus
    nodePort: 31100
    port: 15030
    protocol: TCP
    targetPort: 15030
  - name: https-grafana
    nodePort: 31773
    port: 15031
    protocol: TCP
    targetPort: 15031
  - name: https-tracing
    nodePort: 30877
    port: 15032
    protocol: TCP
    targetPort: 15032
  - name: tls
    nodePort: 32541
    port: 15443
    protocol: TCP
    targetPort: 15443
  selector:
    app: istio-ingressgateway
    istio: ingressgateway
    release: istio
  sessionAffinity: None
  type: LoadBalancer
status:
  loadBalancer:
    ingress:
    - ip: 129.204.7.185

```

而bookinfo-gateway.yaml就是sample/bookinfo/networking中的文件，描述了bookinfo定义了什么样的网关（也即之前操作所定义的入口网关）。内部定义了一个gateway和一个VirtualService。这个gateway和上面提到的gateway绑定，暴露的端口是80，使用的协议是http

```yaml
apiVersion: networking.istio.io/v1alpha3
kind: Gateway
metadata:
  name: bookinfo-gateway
spec:
  selector:
    istio: ingressgateway # use istio default controller
  servers:
  - port:
      number: 80
      name: http
      protocol: HTTP
    hosts:
    - "*"
---
apiVersion: networking.istio.io/v1alpha3
kind: VirtualService
metadata:
  name: bookinfo
spec:
  hosts:
  - "*"
  gateways:
  - bookinfo-gateway
  http:
  - match:
    - uri:
        exact: /productpage
    - uri:
        exact: /login
    - uri:
        exact: /logout
    - uri:
        prefix: /api/v1/products
    route:
    - destination:
        host: productpage
        port:
          number: 9080

```

这样子访问`http://129.204.7.185:31380/productpage`就可以访问到bookinfo的服务了，至于为什么是31380，如果输入`kubectl get svc -nistio-system | grep 31380`可以看到在istio-ingressgateway中它会将80端口映射到外网的31380端口，也即上面的name:http2的配置。

这样子，bookinfo的例子就暂时告一段落，后面istio的特性将继续进行学习。