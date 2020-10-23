# prometheus学习笔记

标签：prometheus

## 资料

* [从零开始：Prometheus](https://blog.csdn.net/qq_37843943/article/details/80510976)非常的一般。
* [实战Prometheus搭建监控系统](https://www.aneasystone.com/archives/2018/11/prometheus-in-action.html)，十分的不错。

## 概念

### [Metric Type](https://prometheus.io/docs/concepts/metric_types/)

Counter, Gauge, Histogram, Summary

## [Querying](https://prometheus.io/docs/prometheus/latest/querying/basics/)

Prometheus提供了一套名为PromQL的查询语言，从而使得用户可以获得并聚合实时的数据。

符号：

* `=`等于某个标签
* `!=`不等于某个标签
* `=~`使用正则匹配
* `!~`使用正则不匹配

Selectors:

* Instant vector selectors:`{}`符号，用于根据label选择指定的时间序列
* Range Vector Selectors:`[]`符号，用于配合s、m等时间符号选定时间范围。可以使用`:`来选择精度。

### 操作符

Aggregation operators:

* sum：求和
* min：最小值
* max：最大值
* avg，计算平均值
* stddev，standard deviation
* stdvar，standard variance
* count，计数
* count_values，计算同样值的元素数量
* bottomk，最小的k个值
* topk，最大的k个值
* quantile，这个是真的不知道

`<aggr-op>([parameter,] <vector expression>) [without|by (<label list>)]`

### Functions

常见的内置函数：

* rate：`rate(v range-vector)`
* abs：`abs(v instant-vector)`绝对值
* ceil：`ceil(v instant-vector)`求底
* floor：求顶

### [Example](https://prometheus.io/docs/prometheus/latest/querying/examples/)

从学习的角度来说例子会更容易懂。

### [HTTP API](https://prometheus.io/docs/prometheus/latest/querying/api/)



## istio下的流量监控

istio请求入口应用product_page的请求数，可以看到source基本为空，source_app="unknown"，而destination_app="productpage"，同时destination_workload="productpage-v1"。

`istio_requests_total{connection_security_policy="none",destination_app="productpage",destination_principal="unknown",destination_service="productpage.default.svc.cluster.local",destination_service_name="productpage",destination_service_namespace="default",destination_version="v1",destination_workload="productpage-v1",destination_workload_namespace="default",instance="172.16.0.158:42422",job="istio-mesh",reporter="destination",request_protocol="http",response_code="503",source_app="unknown",source_principal="unknown",source_version="unknown",source_workload="client",source_workload_namespace="unknown"}`

使用`http_request_duration_microseconds`可以请求到请求延迟，但是这个延迟是跟着节点走的，每个机器三条，分别代表着50%,90%和99%。

## 实现prometheus请求响应时间

英文：

* [stackoverflow](https://stackoverflow.com/questions/47305424/measure-service-latency-with-prometheus)上的问题
* [使用golang实现Prometheus自定义监控例子](https://medium.com/@zhimin.wen/custom-prometheus-metrics-for-apps-running-in-kubernetes-498d69ada7aa)
* [prometheus golang API](https://godoc.org/github.com/prometheus/client_golang/prometheus)
* [istio - Collecting Metrics](https://istio.io/docs/tasks/telemetry/metrics/collecting-metrics/)，非常重要，如何使用istio上的prometheus收集自定义数据。
* [Application metrics in istio](https://meteatamel.wordpress.com/2019/01/07/application-metrics-in-istio/)，提出了如果使用自定义metrics的话，istio环境下需要进行什么样的配置的问题。做法是在pod的spec的template下的annotations中添加：`prometheus.io/scrape:"true"`与`prometheus.io/port="8080"`两行

中文资料：

* [使用golang编写Prometheus Exporter](https://blog.csdn.net/u014029783/article/details/80001251)，不仅仅是怎样在web service上使用prometheus，而且还可以将其注册到正在运行的prometheus服务器上。

### 一些常见的prometheus的endpoint

* /config: to see the current configuration of Prometheus.
* /metrics: to see the scraped metrics.
* /targets: to see the targets that’s being scraped and their status.


华为云yaml备份
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  annotations:
    description: 'test for prometheus'
  labels:
    app.kubernetes.io/name: 'wtytest'
    version: 'v1'
  name: 'wtytest'
  namespace: default
spec:
  replicas: 1
  selector:
    matchLabels:
      app: 'wtytest'
  template:
    metadata:
      annotations:
        metrics.alpha.kubernetes.io/custom-endpoints: '[{"api":"","path":"","port":"","names":""}]'
        prometheus.io/scrape: 'true'
        prometheus.io/port: '8080'
      labels:
        app: 'wtytest'
        version: 'v1'
    spec:
      containers: 
        - name: wtytest
          image: 'wtysos11/prometheus'
          ports:
            - containerPort: 8080
              protocol: TCP
      terminationGracePeriodSeconds: 30
      dnsConfig:
        options:
          - name: single-request-reopen
      affinity: {}
      tolerations:
        - key: node.kubernetes.io/not-ready
          operator: Exists
          effect: NoExecute
          tolerationSeconds: 300
        - key: node.kubernetes.io/unreachable
          operator: Exists
          effect: NoExecute
          tolerationSeconds: 300
  minReadySeconds: 0
  strategy:
    type: RollingUpdate
    rollingUpdate:
      maxSurge: 0
      maxUnavailable: 1
```

第二个想法参见[属性参考](https://istio.io/docs/reference/config/policy-and-telemetry/attribute-vocabulary/)中的connection.duration和response.duration，以及[Collecting Metrics](https://istio.io/docs/tasks/telemetry/metrics/collecting-metrics/)，通过添加自定义的规则来取得响应时间的数据。

华为云上默认的rule如下：

```yaml
 Please edit the object below. Lines beginning with a '#' will be ignored,
# and an empty file will abort the edit. If an error occurs while saving this file will be
# reopened with the relevant failures.
#
apiVersion: config.istio.io/v1alpha2
kind: rule
metadata:
  annotations:
    kubectl.kubernetes.io/last-applied-configuration: |
      {"apiVersion":"config.istio.io/v1alpha2","kind":"rule","metadata":{"annotations":{},"name":"promhttp","namespace":"istio-system"},"spec":{"actions":[{"handler":"handler.prometheus","instances":["r
  creationTimestamp: 2019-03-20T12:18:17Z
  generation: 1
  name: promhttp
  namespace: istio-system
  resourceVersion: "1158765"
  selfLink: /apis/config.istio.io/v1alpha2/namespaces/istio-system/rules/promhttp
  uid: 3de81518-4b0a-11e9-bfe4-fa163ebf73fb
spec:
  actions:
  - handler: handler.prometheus
    instances:
    - requestcount.metric
    - requestduration.metric
    - requestsize.metric
    - responsesize.metric
  match: context.protocol == "http" || context.protocol == "grpc"
```

由上可知官方的handler为handler.prometheus，该rule实现了requestcount.metric等监控instance

通过`kubectl edit prometheus/handler -n=istio-system`可以访问prometheus的配置文件。

配置文件如图：

```yaml
apiVersion: config.istio.io/v1alpha2
kind: prometheus
metadata:
  annotations:
    kubectl.kubernetes.io/last-applied-configuration: |
{  
   "apiVersion":"config.istio.io/v1alpha2",
   "kind":"prometheus",
   "metadata":{  
      "annotations":{  

      },
      "name":"handler",
      "namespace":"istio-system"
   },
   "spec":{  
      "metrics":[  
         {  
            "instance_name":"requestcount.metric.istio-system",
            "kind":"COUNTER",
            "label_names":[  
               "reporter",
               "source_app",
               "source_principal",
               "source_workload",
               "source_workload_namespace",
               "source_version",
               "destination_app",
               "destination_principal",
               "destination_workload",
               "destination_workload_namespace",
               "destination_version",
               "destination_service",
               "destination_service_name",
               "destination_service_namespace",
               "request_protocol",
               "response_code",
               "connection_security_policy"
            ],
            "name":"requests_total"
         },
         {  
            "buckets":{  
               "explicit_buckets":{  
                  "bounds":[  
                     0.005,
                     0.01,
                     0.025,
                     0.05,
                     0.1,
                     0.25,
                     0.5,
                     1,
                     2.5,
                     5,
                     10
                  ]
               }
            },
            "instance_name":"requestduration.metric.istio-system",
            "kind":"DISTRIBUTION",
            "label_names":[  
               "reporter",
               "source_app",
               "source_principal",
               "source_workload",
               "source_workload_namespace",
               "source_version",
               "destination_app",
               "destination_principal",
               "destination_workload",
               "destination_workload_namespace",
               "destination_version",
               "destination_service",
               "destination_service_name",
               "destination_service_namespace",
               "request_protocol",
               "response_code",
               "connection_security_policy"
            ],
            "name":"request_duration_seconds"
         },
         {  
            "buckets":{  
               "exponentialBuckets":{  
                  "growthFactor":10,
                  "numFiniteBuckets":8,
                  "scale":1
               }
            },
            "instance_name":"requestsize.metric.istio-system",
            "kind":"DISTRIBUTION",
            "label_names":[  
               "reporter",
               "source_app",
               "source_principal",
               "source_workload",
               "source_workload_namespace",
               "source_version",
               "destination_app",
               "destination_principal",
               "destination_workload",
               "destination_workload_namespace",
               "destination_version",
               "destination_service",
               "destination_service_name",
               "destination_service_namespace",
               "request_protocol",
               "response_code",
               "connection_security_policy"
            ],
            "name":"request_bytes"
         },
         {  
            "buckets":{  
               "exponentialBuckets":{  
                  "growthFactor":10,
                  "numFiniteBuckets":8,
                  "scale":1
               }
            },
            "instance_name":"responsesize.metric.istio-system",
            "kind":"DISTRIBUTION",
            "label_names":[  
               "reporter",
               "source_app",
               "source_principal",
               "source_workload",
               "source_workload_namespace",
               "source_version",
               "destination_app",
               "destination_principal",
               "destination_workload",
               "destination_workload_namespace",
               "destination_version",
               "destination_service",
               "destination_service_name",
               "destination_service_namespace",
               "request_protocol",
               "response_code",
               "connection_security_policy"
            ],
            "name":"response_bytes"
         },
         {  
            "instance_name":"tcpbytesent.metric.istio-system",
            "kind":"COUNTER",
            "label_names":[  
               "reporter",
               "source_app",
               "source_principal",
               "source_workload",
               "source_workload_namespace",
               "source_version",
               "destination_app",
               "destination_principal",
               "destination_workload",
               "destination_workload_namespace",
               "destination_version",
               "destination_service",
               "destination_service_name",
               "destination_service_namespace",
               "connection_security_policy"
            ],
            "name":"tcp_sent_bytes_total"
         },
         {  
            "instance_name":"tcpbytereceived.metric.istio-system",
            "kind":"COUNTER",
            "label_names":[  
               "reporter",
               "source_app",
               "source_principal",
               "source_workload",
               "source_workload_namespace",
               "source_version",
               "destination_app",
               "destination_principal",
               "destination_workload",
               "destination_workload_namespace",
               "destination_version",
               "destination_service",
               "destination_service_name",
               "destination_service_namespace",
               "connection_security_policy"
            ],
            "name":"tcp_received_bytes_total"
         }
      ]
   }
}
  creationTimestamp: 2019-03-20T12:18:17Z
  generation: 1
  name: handler
  namespace: istio-system
  resourceVersion: "1158764"
  selfLink: /apis/config.istio.io/v1alpha2/namespaces/istio-system/prometheuses/handler
  uid: 3de671d4-4b0a-11e9-bfe4-fa163ebf73fb           
spec:                                                 
  metrics:    

  - instance_name: requestcount.metric.istio-system   
    kind: COUNTER                                     
    label_names:                                      
    - reporter                                        
    - source_app                                      
    - source_principal
    - source_workload                                 
    - source_workload_namespace                       
    - source_version                                  
    - destination_app                                 
    - destination_principal                           
    - destination_workload                            
    - destination_workload_namespace                  
    - destination_version                             
    - destination_service                             
    - destination_service_name                        
    - destination_service_namespace                   
    - request_protocol                                
    - response_code                                   
    - connection_security_policy                      

    name: requests_total 
  - buckets:
      explicit_bu
        bounds:
        - 0.005
        - 0.01
        - 0.025
        - 0.05
        - 0.1
        - 0.25
        - 0.5       
        - 1         
        - 2.5       
        - 5         
        - 10        
    instance_name   
    kind: DISTRIBUTION                                
    label_names:                                      
    - reporter                                        
    - source_app                                      
    - source_principal                                
    - source_workload                                 
    - source_workload_namespace                       
    - source_version                                  
    - destination_app                                 
    - destination_principal                           
    - destination_workload                            
    - destination_workload_namespace                  
    - destination_version                             
    - destination_service                             
    - destination_service_name                        
    - destination_service_namespace                   
    - request_protocol                                
    - response_code                                   
    - connection_security_policy                      
    name: request_duration_seconds                    
  - buckets:                                          
      exponentialBuckets:                             
        growthFactor: 10                              
        numFiniteBuckets: 8                           
        scale: 1                                      

    instance_name: requestsize.metric.istio-system    
    kind: DISTRIBUTION                                
    label_names: 
    - reporter
    - source_app
    - source_principal
    - source_workload
    - source_workload_namespace
    - source_version
    - destination_app
    - destination_principal
    - destination_workload
    - destination_workload_namespace
    - destination_version
    - destination_service
    - destination_service_name
    - destination_service_namespace                   
    - request_protocol                                
    - response_code                                   
    - connection_security_policy                      
    name: request_bytes                               
  - buckets:                                          
      exponentialBuckets:                             
        growthFactor: 10                              
        numFiniteBuckets: 8                           
        scale: 1          

    instance_name: responsesize.metric.istio-system   
    kind: DISTRIBUTION                                
    label_names:                                      
    - reporter                                        
    - source_app                                      
    - source_principal                                
    - source_workload                                 
    - source_workload_namespace                       
    - source_version                                  
    - destination_app                                 
    - destination_principal                           
    - destination_workload                            
    - destination_workload_namespace                  
    - destination_version                             
    - destination_service                             
    - destination_service_name                        
    - destination_service_namespace                   
    - request_protocol                                
    - response_code                                   
    - connection_security_policy
    name: response_bytes

  - instance_name: tcpbytesent.metric.istio-system
    kind: COUNTER
    label_names:
    - reporter
    - source_app
    - source_principal
    - source_workload
    - source_workload_namespace
    - source_version                                  
    - destination_app                                 
    - destination_principal                           
    - destination_workload                            
    - destination_workload_namespace                  
    - destination_version                             
    - destination_service                             
    - destination_service_name                        
    - destination_service_namespace                   
    - connection_security_policy                      
    name: tcp_sent_bytes_total                        

  - instance_name: tcpbytereceived.metric.istio-system
    kind: COUNTER                                     
    label_names:                                      
    - reporter                                        
    - source_app                                      
    - source_principal                                
    - source_workload                                 
    - source_workload_namespace                       
    - source_version                                  
    - destination_app                                 
    - destination_principal                           
    - destination_workload                            
    - destination_workload_namespace                  
    - destination_version                             
    - destination_service                             
    - destination_service_name                        
    - destination_service_namespace                   
    - connection_security_policy                      
    name: tcp_received_bytes_total                    
```

metrics:responsesize

```yaml
apiVersion: config.istio.io/v1alpha2
kind: metric
metadata:
  annotations:
    kubectl.kubernetes.io/last-applied-configuration: |
      {"apiVersion":"config.istio.io/v1alpha2","kind":"metric","metadata":{"annotations":{},"name":"responsesize","namespace":"istio-system"},"spec":{"dimensions":{"connection_security_policy":"conditio
  creationTimestamp: 2019-03-20T12:18:17Z
  generation: 1
  name: responsesize
  namespace: istio-system
  resourceVersion: "1158761"
  selfLink: /apis/config.istio.io/v1alpha2/namespaces/istio-system/metrics/responsesize
  uid: 3de13280-4b0a-11e9-bfe4-fa163ebf73fb
spec:
  dimensions:
    connection_security_policy: conditional((context.reporter.kind | "inbound") ==
      "outbound", "unknown", conditional(connection.mtls | false, "mutual_tls", "none"))
    destination_app: destination.labels["app"] | "unknown"
    destination_principal: destination.principal | "unknown"
    destination_service: destination.service.host | "unknown"
    destination_service_name: destination.service.name | "unknown"
    destination_service_namespace: destination.service.namespace | "unknown"
    destination_version: destination.labels["version"] | "unknown"
    destination_workload: destination.workload.name | "unknown"
    destination_workload_namespace: destination.workload.namespace | "unknown"
    reporter: conditional((context.reporter.kind | "inbound") == "outbound", "source",
      "destination")
    request_protocol: api.protocol | context.protocol | "unknown"
    response_code: response.code | 200
    source_app: source.labels["app"] | "unknown"
    source_principal: source.principal | "unknown"
    source_version: source.labels["version"] | "unknown"
    source_workload: source.workload.name | "unknown"
    source_workload_namespace: source.workload.namespace | "unknown"
  monitored_resource_type: '"UNSPECIFIED"'
  value: response.size | 0
```

requestsize

```yaml
apiVersion: config.istio.io/v1alpha2
kind: metric
metadata:
  annotations:
    kubectl.kubernetes.io/last-applied-configuration: |
      {"apiVersion":"config.istio.io/v1alpha2","kind":"metric","metadata":{"annotations":{},"name":"requestsize","namespace":"istio-system"},"spec":{"dimensions":{"connection_security_policy":"condition
  creationTimestamp: 2019-03-20T12:18:17Z
  generation: 1
  name: requestsize
  namespace: istio-system
  resourceVersion: "1158760"
  selfLink: /apis/config.istio.io/v1alpha2/namespaces/istio-system/metrics/requestsize
  uid: 3ddf9d1c-4b0a-11e9-bfe4-fa163ebf73fb
spec:
  dimensions:
    connection_security_policy: conditional((context.reporter.kind | "inbound") ==
      "outbound", "unknown", conditional(connection.mtls | false, "mutual_tls", "none"))
    destination_app: destination.labels["app"] | "unknown"
    destination_principal: destination.principal | "unknown"
    destination_service: destination.service.host | "unknown"
    destination_service_name: destination.service.name | "unknown"
    destination_service_namespace: destination.service.namespace | "unknown"
    destination_version: destination.labels["version"] | "unknown"
    destination_workload: destination.workload.name | "unknown"
    destination_workload_namespace: destination.workload.namespace | "unknown"
    reporter: conditional((context.reporter.kind | "inbound") == "outbound", "source",
      "destination")
    request_protocol: api.protocol | context.protocol | "unknown"
    response_code: response.code | 200
    source_app: source.labels["app"] | "unknown"
    source_principal: source.principal | "unknown"
    source_version: source.labels["version"] | "unknown"
    source_workload: source.workload.name | "unknown"
    source_workload_namespace: source.workload.namespace | "unknown"
  monitored_resource_type: '"UNSPECIFIED"'
  value: request.size | 0
```

requestcount

```yaml
apiVersion: config.istio.io/v1alpha2
kind: metric
metadata:
  annotations:
    kubectl.kubernetes.io/last-applied-configuration: |
      {"apiVersion":"config.istio.io/v1alpha2","kind":"metric","metadata":{"annotations":{},"name":"requestcount","namespace":"istio-system"},"spec":{"dimensions":{"connection_security_policy":"conditio
  creationTimestamp: 2019-03-20T12:18:17Z
  generation: 1
  name: requestcount
  namespace: istio-system
  resourceVersion: "1158758"
  selfLink: /apis/config.istio.io/v1alpha2/namespaces/istio-system/metrics/requestcount
  uid: 3ddc6325-4b0a-11e9-bfe4-fa163ebf73fb
spec:
  dimensions:
    connection_security_policy: conditional((context.reporter.kind | "inbound") ==
      "outbound", "unknown", conditional(connection.mtls | false, "mutual_tls", "none"))
    destination_app: destination.labels["app"] | "unknown"
    destination_principal: destination.principal | "unknown"
    destination_service: destination.service.host | "unknown"
    destination_service_name: destination.service.name | "unknown"
    destination_service_namespace: destination.service.namespace | "unknown"
    destination_version: destination.labels["version"] | "unknown"
    destination_workload: destination.workload.name | "unknown"
    destination_workload_namespace: destination.workload.namespace | "unknown"
    reporter: conditional((context.reporter.kind | "inbound") == "outbound", "source",
      "destination")
    request_protocol: api.protocol | context.protocol | "unknown"
    response_code: response.code | 200
    source_app: source.labels["app"] | "unknown"
    source_principal: source.principal | "unknown"
    source_version: source.labels["version"] | "unknown"
    source_workload: source.workload.name | "unknown"
    source_workload_namespace: source.workload.namespace | "unknown"
  monitored_resource_type: '"UNSPECIFIED"'
  value: "1"
```

最终解决方式是补上了requestDuration的yaml

```yaml
apiVersion: config.istio.io/v1alpha2
kind: metric
metadata:
  annotations:
    kubectl.kubernetes.io/last-applied-configuration: |
      {"apiVersion":"config.istio.io/v1alpha2","kind":"metric","metadata":{"annotations":{},"name":"requestduration","namespace":"istio-system"},"spec":{"dimensions":{"connection_security_policy":"conditional((context.reporter.kind | \"inbound\") == \"outbound\", \"unknown\", conditional(connection.mtls | false, \"mutual_tls\", \"none\"))","destination_app":"destination.labels[\"app\"] | \"unknown\"","destination_principal":"destination.principal | \"unknown\"","destination_service":"destination.service.host | \"unknown\"","destination_service_name":"destination.service.name | \"unknown\"","destination_service_namespace":"destination.service.namespace | \"unknown\"","destination_version":"destination.labels[\"version\"] | \"unknown\"","destination_workload":"destination.workload.name | \"unknown\"","destination_workload_namespace":"destination.workload.namespace | \"unknown\"","reporter":"conditional((context.reporter.kind | \"inbound\") == \"outbound\", \"source\", \"destination\")","request_protocol":"api.protocol | context.protocol | \"unknown\"","response_code":"response.code | 200","source_app":"source.labels[\"app\"] | \"unknown\"","source_principal":"source.principal | \"unknown\"","source_version":"source.labels[\"version\"] | \"unknown\"","source_workload":"source.workload.name | \"unknown\"","source_workload_namespace":"source.workload.namespace | \"unknown\""},"monitored_resource_type":"\"UNSPECIFIED\"","value":"response.duration | \"0ms\""}}
  creationTimestamp: "2019-06-22T12:57:59Z"
  generation: 1
  name: requestduration
  namespace: istio-system
  resourceVersion: "4756465"
  selfLink: /apis/config.istio.io/v1alpha2/namespaces/istio-system/metrics/requestduration
  uid: 5c795634-94ed-11e9-a62c-5254001ba903
spec:
  dimensions:
    connection_security_policy: conditional((context.reporter.kind | "inbound") ==
      "outbound", "unknown", conditional(connection.mtls | false, "mutual_tls", "none"))
    destination_app: destination.labels["app"] | "unknown"
    destination_principal: destination.principal | "unknown"
    destination_service: destination.service.host | "unknown"
    destination_service_name: destination.service.name | "unknown"
    destination_service_namespace: destination.service.namespace | "unknown"
    destination_version: destination.labels["version"] | "unknown"
    destination_workload: destination.workload.name | "unknown"
    destination_workload_namespace: destination.workload.namespace | "unknown"
    reporter: conditional((context.reporter.kind | "inbound") == "outbound", "source",
      "destination")
    request_protocol: api.protocol | context.protocol | "unknown"
    response_code: response.code | 200
    source_app: source.labels["app"] | "unknown"
    source_principal: source.principal | "unknown"
    source_version: source.labels["version"] | "unknown"
    source_workload: source.workload.name | "unknown"
    source_workload_namespace: source.workload.namespace | "unknown"
  monitored_resource_type: '"UNSPECIFIED"'
  value: response.duration | "0ms"

```