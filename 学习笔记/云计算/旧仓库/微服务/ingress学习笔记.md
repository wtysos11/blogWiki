# ingress学习笔记

标签：kubernetes ingress 学习笔记

ingress是什么：ingress是kubernetes支持的API对象，负责管理集群对外部的连接，特别是HTTP服务。

## 配置模版

实例：

```yaml
apiVersion: extensions/v1beta1
kind: Ingress
metadata:
  name: test-ingress
  annotations:
    nginx.ingress.kubernetes.io/rewrite-target: /
spec:
  rules:
  - http:
      paths:
      - path: /testpath
        backend:
          serviceName: test
          servicePort: 80
```

### rule条目

每个http rule包含下面的信息

* 可选的host：在本例子中，没有指定host，所以rule被用在所有发送到主机的HTTP流量上。如果指定了host，则只针对URL为这个host的主机。
* paths条目，每个paths条目中跟着一个path词条，里面有一个backend属性。backend被serviceName和servicePort两个属性所定义。host和path必须与到来的request相对应。
* backend是service和端口名称的组合。满足请求的HTTP request会被发往指定的backend。

一般而言，会有一个默认的backend来配置ingress controller，处理所有不满足任何path的请求。

### ingress类型

#### Single Service Ingress

K8s允许你只暴露一个Service，你可以使用默认backend的ingress来实现它。

```yaml
apiVersion: extensions/v1beta1
kind: Ingress
metadata:
  name: test-ingress
spec:
  backend:
    serviceName: testsvc
    servicePort: 80
```

#### simple fanout

指定host，解析过程如下：

```
foo.bar.com -> 178.91.123.132 -> / foo    service1:4200
                                 / bar    service2:8080
```

```yaml
apiVersion: extensions/v1beta1
kind: Ingress
metadata:
  name: simple-fanout-example
  annotations:
    nginx.ingress.kubernetes.io/rewrite-target: /
spec:
  rules:
  - host: foo.bar.com
    http:
      paths:
      - path: /foo
        backend:
          serviceName: service1
          servicePort: 4200
      - path: /bar
        backend:
          serviceName: service2
          servicePort: 8080
```

ingress controller保证特定的loadbalancer的实现来满足ingress的要求，只要服务s1和s2存在。

#### name based virtual hosting

基于名称的虚拟host，支持多域名的HTTP路由。

```
foo.bar.com --|                 |-> foo.bar.com s1:80
              | 178.91.123.132  |
bar.foo.com --|                 |-> bar.foo.com s2:80
```