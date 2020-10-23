# kubernetes上构建外部能够访问到的应用

标签：kubernetes 服务发现

起因是我在安装dashboard的时候，外部无论如何也不能够访问到dashboard，后面改了service也只是提供了一个内网ip，而且登陆不上。

问题：如果有一个和k8s集群不在同一个子网的机器，它该如何访问k8s集群内部署的服务？

参考：

* [kubernetes如何设置external ip](http://dockone.io/question/903)
* [kubernetes部署ingress](https://www.cnblogs.com/dingbin/p/9754993.html)
* [从外部访问kubernetes中的Pod](https://jimmysong.io/posts/accessing-kubernetes-pods-from-outside-of-the-cluster/)

基本思想：如果云服务商不提供负载均衡器的话，外部想要访问只能够使用NodePort，但是一般需要nginx再进行一次反向代理。而ingress出现后可以用ingress controller读取集群中的ingress规则变化，动态生成nginx配置文件并注入到nginx的Pod中。

为什么要部署ingress：Service NodePort是4层调度，无法直接卸载https协议，而kubernetes ingress是7层调度，可以直接卸载https回话。

## 实验目标

1. 实现外网到内网服务器通过ServicePort的访问
2. 实现外网到内网服务器通过Ingress的访问
3. 实现其他方式的服务发现

## 开始实验

因为我安装了istio，先将istio注入关闭：`kubectl label namespace default istio-injection=disabled --overwrite`，希望我实验完之后记得把它弄回去。

首先配置nginx的Deployment,rc也可以

```yaml
apiVersion: extensions/v1beta1
kind: Deployment
metadata:
    name: test-nginx
    labels:
        app: nginx
spec:
    replicas: 1
    template:
        metadata:
            labels:
                app: nginx
                version: v1
        spec:
            containers:
            - name: nginx
              image: nginx:latest
              imagePullPolicy: IfNotPresent
              ports:
              - containerPort: 80

```

然后是service

```yaml
apiVersion: v1
kind: Service
metadata:
  name: nginx-service
spec:
  type: NodePort
  ports:
  - port: 80
    nodePort: 30010
  selector:
    app: nginx

```

上述节点打开container，也就是nginx的80端口，并且在service上做30010到80端口的映射。

这时候就可以在集群内部进行访问，集群内部所有节点都可以通过访问节点内网IP的30010端口访问到nginx的服务器。公网通过直接访问ip也可以访问到，但是经过测试只能够在nginx的Pod所运行的服务器上成立，也就是说其他机子是不能通过30010端口访问到nginx服务的。

这样子是不行的，因为公网可能会有多个出口节点，这时候如果只有一个节点能够访问通，显然不可以。虽然可以通过nginx解决，但是如果每加一条服务就加一条nginx配置的话，未免过于繁琐。

通过搜索我发现了kubernetes中新出了ingress的机制，希望它可以帮助到我。

## 安装ingress

遵循官方的[安装指南](https://kubernetes.github.io/ingress-nginx/deploy/)

首先执行命令`kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/master/deploy/mandatory.yaml`进行安装。

需要安装metalLB作为负载均衡器，运行`kubectl apply -f https://raw.githubusercontent.com/google/metallb/v0.7.3/manifests/metallb.yaml`，里面包含了`metallb-system/controller`的deployment和`metallb-system/speaker`的daemonset，以及一些service。在`kubectl get po -nmetallb-system`指令下确认所有的pods都成功运行后进行下一步。

根据[官方文档](https://kubernetes.github.io/ingress-nginx/deploy/baremetal/)的说法，或者用[metallb的说法](https://metallb.universe.tf/tutorial/layer2/)，使用configmap来配置IP池。我使用了自己的IP。（注：即使只有一个IP也要表示为范围的形式，不然会报错）

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  namespace: metallb-system
  name: config
data:
  config: |
    address-pools:
    - name: default
      protocol: layer2
      addresses:
      - 129.204.7.185-129.204.7.185

```

然后执行Bare-metal-using NodePort部分的命令`kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/master/deploy/provider/baremetal/service-nodeport.yaml`

运行`kubectl get pods --all-namespaces -l app.kubernetes.io/name=ingress-nginx`可以查看生成的pods，可以看到这时候pods成功运行。

### 查看运行是否成功

这时候，删除之前的Nginx，运行`kubectl apply -f https://raw.githubusercontent.com/google/metallb/v0.7.3/manifests/tutorial-2.yaml`，里面有官方配置好的nginx镜像和service。

运行`kubectl get sertvice nginx`，如果有出现ExternalIP，则任务完成。这时候，只要访问外部的IP的指定端口，就可以访问到内网的服务。

真是不容易啊，中间遇到了好多的问题，感谢互联网之神，今天又顺利的解决了呢。