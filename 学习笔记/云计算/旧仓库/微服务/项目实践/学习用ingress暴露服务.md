# 学习用ingress暴露服务

标签：kubernetes ingress 服务发现

之前我一直使用的是NodePort和LoadBalancer，但是今天在暴露prometheus时发现这两种都不行。有一篇[博客](https://www.qikqiak.com/post/kubernetes-monitor-prometheus-grafana/)提到了用ingress暴露http服务，我之前没有试过ingress，先做个实验吧。

参考文章：

* [Kubernetes Nginx Ingress教程](https://mritd.me/2017/03/04/how-to-use-nginx-ingress/)
* [kubernetes ingress 官方文档](https://kubernetes.io/docs/concepts/services-networking/ingress/)

首先创建nginx基础服务

```yaml
apiVersion: apps/v1beta2
kind: Deployment
metadata:
  name: nginx
spec:
  selector:
    matchLabels:
      app: nginx
  template:
    metadata:
      labels:
        app: nginx
    spec:
      containers:
      - name: nginx
        image: nginx:1
        ports:
        - name: http
          containerPort: 80

---
apiVersion: v1
kind: Service
metadata:
  name: nginx
spec:
  ports:
  - name: http
    port: 80
    protocol: TCP
    targetPort: 80
  selector:
    app: nginx
  type: ClusterIP

```


然后使用最简单的Single Service Ingress来配置。

```yaml
apiVersion: extensions/v1beta1
kind: Ingress
metadata:
  name: test-ingress
spec:
  backend:
    serviceName: nginx
    servicePort: 80
```

这里我有个疑问，我该怎么样访问到这里服务呢？毕竟集群有三个服务器，每个服务器都有一个外网IP和一个内网IP。然后我就没有疑惑了：

```
[root@k8s-master code]# kubectl get ingress
NAME           HOSTS   ADDRESS   PORTS   AGE
test-ingress   *                 80      29s
```

根本就没有分配ip嘛。

## 尝试1

根据[stackoverflow](https://stackoverflow.com/questions/49845021/getting-an-kubernetes-ingress-endpoint-ip-address)上的一篇回答，我认为可能是原始的ingress-nginx采用的NodePort的原因，下面将改写为LoadBalancer，看看是否有变化。

修改前

```yaml
apiVersion: v1
kind: Service
metadata:
  annotations:
    kubectl.kubernetes.io/last-applied-configuration: |
      {"apiVersion":"v1","kind":"Service","metadata":{"annotations":{},"labels":{"app.kubernetes.io/name":"ingress-nginx","app.kubernetes.io/part-of":"ingress-nginx"},"name":"ingress-nginx","namespace":"ingress-nginx"},"spec":{"ports":[{"name":"http","port":80,"protocol":"TCP","targetPort":80},{"name":"https","port":443,"protocol":"TCP","targetPort":443}],"selector":{"app.kubernetes.io/name":"ingress-nginx","app.kubernetes.io/part-of":"ingress-nginx"},"type":"NodePort"}}
  creationTimestamp: "2019-04-06T02:01:43Z"
  labels:
    app.kubernetes.io/name: ingress-nginx
    app.kubernetes.io/part-of: ingress-nginx
  name: ingress-nginx
  namespace: ingress-nginx
  resourceVersion: "4194999"
  selfLink: /api/v1/namespaces/ingress-nginx/services/ingress-nginx
  uid: ece1c1cf-580f-11e9-96eb-5254001218ec
spec:
  clusterIP: 10.101.73.96
  externalTrafficPolicy: Cluster
  ports:
  - name: http
    nodePort: 32623
    port: 80
    protocol: TCP
    targetPort: 80
  - name: https
    nodePort: 32538
    port: 443
    protocol: TCP
    targetPort: 443
  selector:
    app.kubernetes.io/name: ingress-nginx
    app.kubernetes.io/part-of: ingress-nginx
  sessionAffinity: None
  type: NodePort
status:
  loadBalancer: {}
```

修改后（我没有改下面的status，它是自己加上去的）

```yaml
apiVersion: v1
kind: Service
metadata:
  annotations:
    kubectl.kubernetes.io/last-applied-configuration: |
      {"apiVersion":"v1","kind":"Service","metadata":{"annotations":{},"labels":{"app.kubernetes.io/name":"ingress-nginx","app.kubernetes.io/part-of":"ingress-nginx"},"name":"ingress-nginx","namespace":"ingress-nginx"},"spec":{"ports":[{"name":"http","port":80,"protocol":"TCP","targetPort":80},{"name":"https","port":443,"protocol":"TCP","targetPort":443}],"selector":{"app.kubernetes.io/name":"ingress-nginx","app.kubernetes.io/part-of":"ingress-nginx"},"type":"NodePort"}}
  creationTimestamp: "2019-04-06T02:01:43Z"
  labels:
    app.kubernetes.io/name: ingress-nginx
    app.kubernetes.io/part-of: ingress-nginx
  name: ingress-nginx
  namespace: ingress-nginx
  resourceVersion: "6879819"
  selfLink: /api/v1/namespaces/ingress-nginx/services/ingress-nginx
  uid: ece1c1cf-580f-11e9-96eb-5254001218ec
spec:
  clusterIP: 10.101.73.96
  externalIPs:
  - 129.204.90.87
  - 203.195.219.185
  externalTrafficPolicy: Cluster
  ports:
  - name: http
    nodePort: 30828
    port: 80
    protocol: TCP
    targetPort: 80
  - name: https
    nodePort: 31802
    port: 443
    protocol: TCP
    targetPort: 443
  selector:
    app.kubernetes.io/name: ingress-nginx
    app.kubernetes.io/part-of: ingress-nginx
  sessionAffinity: None
  type: LoadBalancer
status:
  loadBalancer:
    ingress:
    - ip: 203.195.219.185
```

修改后查看

```
NAME            TYPE           CLUSTER-IP     EXTERNAL-IP                                     PORT(S)                      AGE
ingress-nginx   LoadBalancer   10.101.73.96   203.195.219.185,129.204.90.87,203.195.219.185   80:30828/TCP,443:31802/TCP   17d
```

可以看到现在ingress-nginx变为了负载均衡器，再次查看原始的Ingress，可以看到确实有分配ip

```
[root@k8s-master code]# kubectl get ingress
NAME           HOSTS   ADDRESS                         PORTS   AGE
test-ingress   *       129.204.90.87,203.195.219.185   80      19m
```

可还是访问不了服务，nginx根本没有反应。修改了原始的svc为ClusterIP，取消了NodePort

测试内网IP`curl 10.101.73.96`，发现可以连接上nginx。在集群内部其他机器上测试`curl 203.195.219.185`，发现也能够连接上nginx。应该是中间转发有些问题。