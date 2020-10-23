# istio Ingressgateway配置

已经在之前完成了Istio的部署工作，现在的目标是配置Ingressgateway，并且学会注入sidecar，完成指定pod的流量收集工作。

工作计划：
1. 配置命名空间，在该空间下注入istio-sidecar
2. 暴露ingressgateway
3. 在k8s集群上上载之前的gatling镜像与productpage镜像
4. 进行内网和外网的压力测试

与之前相比，Istio的中文文档已经很完善了，很多东西查起来更方便了。感谢各位大佬们的贡献。

## 配置命名空间

Istio的功能实现依赖于istio-sidecar。根据[官方文档](https://preliminary.istio.io/latest/zh/docs/setup/additional-setup/sidecar-injection/#automatic-sidecar-injection),可以使用`kubectl label namespace default istio-injection=enabled`给`default`命名空间进行标记。这样子所有在默认命名空间下的应用都将被注入sidecar。

当然，这样可能会导致一些问题，所以我选择新建一个命名空间来进行之后的操作。

```yaml
apiVersion: v1
kind: Namespace
metadata:
    name: test
    labels: 
        istio-injection: enabled
```

效果如下：

```
Name:         test
Labels:       istio-injection=enabled
Annotations:  <none>
Status:       Active

No resource quota.

No LimitRange resource.
```

或者使用命令`kubectl get namespaces --show-labels`查看。这样，只要将新的pod在test命名空间下上传给k8s，就可以接入到Istio中。

Istio官方在其文件夹下提供了`samples/sleep/sleep.yaml`作为自动测试注入功能的测试样例，还是很好用的。可以在istio文件夹下执行`kubectl apply -f samples/sleep/sleep.yaml --namespace=test`。

PS：如何设置Deployment的命名空间：
    一般来说是通过设置metadata字段中的namespace字段来实现指定namespace的功能。详见[API](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.19/#objectmeta-v1-meta)(需要往下翻，metadata是属于ObjectMeta类型的)

PPS：很值得注意的是，官方在文档里指导我们如何使用sleep这个比较简单的包完成内网测试。(能够完成测试的原因很简单，因为这个镜像是用curl镜像+sleep命令来做的。所以也不用想着登陆命令行这种操作，做不到的)
```bash
export SLEEP_POD=$(kubectl get pod -l app=sleep -o jsonpath={.items..metadata.name})#这里要自行加上-n=test
kubectl exec -it $SLEEP_POD -c sleep curl http://ratings.default.svc.cluster.local:9080/ratings/1
{"id":1,"ratings":{"Reviewer1":5,"Reviewer2":4}}
```

官方在httpbin里也展示了如何通过`curl`镜像来进行内网测试(同样，注意命名空间)
```bash
kubectl run -i --rm --restart=Never dummy --image=dockerqa/curl:ubuntu-trusty --command -- curl --silent httpbin:8000/html
kubectl run -i --rm --restart=Never dummy --image=dockerqa/curl:ubuntu-trusty --command -- curl --silent --head httpbin:8000/status/500
time kubectl run -i --rm --restart=Never dummy --image=dockerqa/curl:ubuntu-trusty --command -- curl --silent httpbin:8000/delay/5
```

## 内网暴露ingressgateway

此处参考[官方教程](https://istio.io/latest/zh/docs/tasks/traffic-management/ingress/ingress-control/)。

官方推荐使用httpbin来进行网络功能测试。可以使用`kubectl apply -f samples/httpbin/httpbin.yaml --namespace=test`来安装，详情见[istio-httpbin](https://github.com/istio/istio/tree/release-1.7/samples/httpbin)，还是很方便的。

有两种情况，一种是存在外网负载均衡器，一种是不存在外网的LB。我所在的环境是不存在LB的，所以我按照官方的操作执行，设置端口和IP，这样就只能在内网访问了，除非做反向代理。

### 确定ingreegateway的外网IP

```bash
export INGRESS_PORT=$(kubectl -n istio-system get service istio-ingressgateway -o jsonpath='{.spec.ports[?(@.name=="http2")].nodePort}') 
export SECURE_INGRESS_PORT=$(kubectl -n istio-system get service istio-ingressgateway -o jsonpath='{.spec.ports[?(@.name=="https")].nodePort}')
export INGRESS_HOST=$(kubectl get po -l istio=ingressgateway -n istio-system -o jsonpath='{.items[0].status.hostIP}')
```

通过`kubectl edit svc/istio-ingressgateway -n istio-system`命令我们可以查看service所设定的规则。实际上INGREE_PORT端口指的就是配置的http2端口（对应端口80），SECURE_INGRESS_PORT端口指的就是配置的https端口（对应端口443）。INGRESS_HOST实际上就是ingressgateway所在的host节点的IP（所以如果后面ingressgateway发生了迁移，这之后应该是要重新跑一遍）

下面配置网关和路由。这一步在本文中需要注意，因为我的服务是在test命名空间下配置的，因此下面的配套设施也需要在test命名空间下。我之前就是忘了这一步，导致卡了好久。

执行下面两个语句，可以构建网关和路由。
```bash
kubectl apply -f - <<EOF
apiVersion: networking.istio.io/v1alpha3
kind: Gateway
metadata:
  name: httpbin-gateway
spec:
  selector:
    istio: ingressgateway # use Istio default gateway implementation
  servers:
  - port:
      number: 80
      name: http
      protocol: HTTP
    hosts:
    - "httpbin.example.com"
EOF

kubectl apply -f - <<EOF
apiVersion: networking.istio.io/v1alpha3
kind: VirtualService
metadata:
  name: httpbin
spec:
  hosts:
  - "httpbin.example.com"
  gateways:
  - httpbin-gateway
  http:
  - match:
    - uri:
        prefix: /status
    - uri:
        prefix: /delay
    route:
    - destination:
        port:
          number: 8000
        host: httpbin
EOF
```

再执行`curl -s -I -HHost:httpbin.example.com "http://$INGRESS_HOST:$INGRESS_PORT/status/200"`，即可通过ingressgateway访问到指定的Pod。

上述两个文件中大部分都是比较好懂的。第一个文件实现了一个网关，这个网关暴露了自身的80端口用于通信，并且能够处理`httpbin.example.com`这个域名。而且这个80端口是集群内的端口，因为集群外我访问的是31941端口，而这个端口在负载均衡器中是转发到80端口的。

第二个文件实现了一个VirtualService，解析`httpbin.example.com`这个域名，并且只允许通过`/status`和`/delay`这两个路由（我在想如果match这个字段不配置的话是不是可以允许全部通过）。重点是route字段，destination中的port指定了应用的服务所暴露出的端口（类型为ClusterIP），因此这个host很可能指的是服务的名字。

但是我还是有几点疑问

一个是spec->servers->hosts数组中的那个网址，我的理解是这个网址可以被其他节点所解析，其他节点通过访问这个网址即可以访问到80端口。这个网址在第二个文件的spec->hosts中也有所体现。在官方文档的第二个改进中这个hosts被改成了通配符`*`，那原来那个是什么？

第二个是VirtualService的http中的match，是不是没有配置的路由都是不会被访问的？我觉得应该缺色hi是这样子。最终的route->destination->port指向了8000，这个8000应该是httpbin暴露出的端口（80是httpbin暴露的端口，8000是httpbin对应的服务暴露的端口）。但是这个host又很迷，难道是app的名字吗（似乎是的）？

详见[gateway-configuration](https://istio.io/latest/docs/reference/config/networking/gateway/)

## 部署镜像

### 导入镜像

主要是gatling镜像与productpage镜像。productpage镜像还好，直接用`samples`里的版本就可以了。但是`gatling`镜像是我自己拿ubuntu改的，因此就需要docker镜像的导入与导出了。

参考[这篇文章](https://dockerlabs.collabnix.com/beginners/saving-images-as-tar/)来进行镜像导入与导出。

使用`dokcer export CONTAINER_ID > file.tar`导出镜像，并用`docker import - IMAGE_NAME < file.tar`导入镜像。

或者使用`docker save -o file.tar IMAGE_NAME`，然后用`docker load < file.tar`导入。

两者的区别：
* `docker save`用来将一个或多个镜像打包。虽然命令行要求指定image，但实际测试中container也是可以打包的（打包将保存对应的image）。其应用场景是帮助不能链接外网的客户服务器部署容器，因此会打包所有的数据。
* `docker export`，从帮助中是将container的文件系统进行打包，参数必须指定为container。其对应操作`docker import`是可以指定镜像名和标签。其应用场景是用于制作基础镜像，然后分发给其他人开发使用，作为基础的开发环境。

`docker save`保存的是镜像包，`docker export`保存的是容器包，两者的结果是不能混用的，最终都会导致容器不能正常启动。容器包打包的是原有容器的文件系统，而镜像包打包的是构建镜像每一次commit堆叠起来的所有层，是两个完全不同的概念（尽管最终都表现为是压缩包）。

### 调整网关

#### 该如何写网关

将ingressgateway暴露出来并不是一件容易的事情，特别是需要做到特定的要求。之前实验室在华为云上时，可以通过Istio的负载均衡器自动配置，实现访问`外网IP：指定端口`的形式就可以访问到对应的服务。因此我也需要在新的集群上实现同样的功能。

以productpage-test应用为例，通过访问`外网IP：19002`即可直接访问到该服务，并且路由全部正常转发。

其ingressgateway为
```yaml
apiVersion: networking.istio.io/v1alpha3
kind: Gateway
metadata:
  annotations:
    creator: asm
  creationTimestamp: 2020-10-09T04:27:30Z
  generation: 1
  name: productpage-test0-wtytest-gateway
  namespace: istio-system
  resourceVersion: "121292656"
  selfLink: /apis/networking.istio.io/v1alpha3/namespaces/istio-system/gateways/productpage-test0-wtytest-gateway
   uid: be5a3234-09e7-11eb-b6b2-fa163ebf73fb
 spec:
   selector:
     istio: ingressgateway
   servers:
   - hosts:
     - 139.9.57.167
     port:
       name: http-be57584a
       number: 19002
       protocol: http 
```

有一个对应的LB型的service，nodePort为30609，port和targetPort均为19002。

对应的服务为9080的ClusterIP型，忽略。

VirtualService：
```yaml
spec:
   gateways:
   - productpage-test0-wtytest-gateway.istio-system.svc.cluster.local
   - mesh
   hosts:
   - productpage-test0
   - 外网IP
   http:
   - match:
     - gateways:
       - productpage-test0-wtytest-gateway.istio-system.svc.cluster.local
       port: 19002
       uri: {}
     route:
     - destination:
         host: productpage-test0
         port:
           number: 9080
         subset: v1
   - match:
     - gateways:
       - mesh 
       port: 9080
     route:
     - destination:
         host: productpage-test0
         port:
           number: 9080
         subset: v1  
```

这里需要弄清楚的一点是，19002与9080应该如何对应上。从现在的情况来看，19002应该配置在gateway这边，负责转接外网的流量。9080应该配置在VirtualService这边，负责将外网流量对应到内网服务上。

感觉这个更加复杂了，因此还是先参考官方的版本。

#### 进行网关编写

由于productpage是部署在istio上的，需要设置入口网关

Service使用服务自带的Service，其他的也使用自带的部分（刨除了其他服务，因为我只需要一个productpage进行测试）
```yaml
apiVersion: v1
kind: Service
metadata:
  name: productpage
  labels:
    app: productpage
    service: productpage
spec:
  ports:
  - port: 9080
    name: http
  selector:
    app: productpage
```

Gateway:
```yaml
apiVersion: networking.istio.io/v1alpha3
kind: Gateway
metadata:
  name: productpage-gateway
spec:
  selector:
    istio: ingressgateway # use Istio default gateway implementation
  servers:
  - port:
      number: 80
      name: http
      protocol: HTTP
    hosts:
    - "*"
```

VirtualService:
```yaml
apiVersion: networking.istio.io/v1alpha3
kind: VirtualService
metadata:
  name: productpage
spec:
  hosts:
  - "*"
  gateways:
  - productpage-gateway
  http:
  - match:
    - uri:
        prefix: /productpage
    - uri:
        prefix: /
    route:
    - destination:
        port:
          number: 9080
        host: productpage
```

#### 测试访问

```bash
export INGRESS_PORT=$(kubectl -n istio-system get service istio-ingressgateway -o jsonpath='{.spec.ports[?(@.name=="http2")].nodePort}') 
export INGRESS_HOST=$(kubectl get po -l istio=ingressgateway -n istio-system -o jsonpath='{.items[0].status.hostIP}')
```

访问`curl -s -I "http://$INGRESS_HOST:$INGRESS_PORT"`或者加上`/productpage`都可以。问题也很明显，全部都会被转接到这一个服务上，我暂时还没有想到怎么做NodePort，之前试了一下发现转不过去……

### 调整节点亲和

为了充分利用节点资源，需要将两个发送端实例和接收端实例放在不同的物理机节点上。

下面为作为发送端的gatling脚本，设置了affinity进行节点亲和，同时设置了toleration使其能够部署在该节点上。

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: gatling
  labels:
    app: gatling
    version: v1
spec:
  replicas: 1
  selector:
    matchLabels:
      app: gatling
      version: v1
  template:
    metadata:
      labels:
        app: gatling
        version: v1
    spec:
      containers:
      - name: gatling
        image: ubuntu-wtynettest:0.0.1
        imagePullPolicy: Never
      affinity:
        nodeAffinity:
          requiredDuringSchedulingIgnoredDuringExecution:
            nodeSelectorTerms:
              - matchExpressions:
                  - key: kubernetes.io/hostname
                    operator: In
                    values:
                    - workflow-1
      tolerations:
        - key: node-role.kubernetes.io/master
          operator: Exists
          effect: NoSchedule
```

## 内网压力测试

测试反而是最好做的，`kubectl exec -it gatling-5998854d89-mkfwq bash`，登陆后运行gatling，完成。