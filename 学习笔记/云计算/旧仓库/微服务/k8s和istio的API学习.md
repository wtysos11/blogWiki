# k8s和istio的API学习

标签：kubernetes istio API

按照师兄的要求，开始学习k8s和istio的API，目标是能够通过API获取监控数据。

参考资料：

* [k8s api doc](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.10/)
* [k8s官方API概览](https://kubernetes.io/docs/reference/using-api/api-overview/)与[概览中的API](https://kubernetes.io/docs/concepts/overview/kubernetes-api/)
* [k8s官方API文档 using api](https://kubernetes.io/docs/reference/using-api/api-concepts/)
* [k8s官方API文档 accessing api](https://kubernetes.io/docs/reference/access-authn-authz/controlling-access/)
* [istio/api的github](https://github.com/istio/api)
* [kubeapi golang](https://github.com/kubernetes/client-go/)

https相关：

* [关于curl访问https的若干问题](https://blog.csdn.net/q553716434/article/details/39500467)
* [使用https访问kubernetes](https://www.jianshu.com/p/af4f9349c3c2)

中文资料：

* [掘金-大家都较熟悉之kubernetes api解析](https://juejin.im/post/5a11131ff265da4315237fcb)
* [kubernetes源码之旅：从kubectl到API server](https://www.kubernetes.org.cn/2324.html)
* [kubernetes API使用指南](https://www.jianshu.com/p/91d0330f4bf1)
* [kubernetes restful api](https://blog.csdn.net/jinzhencs/article/details/51452208)
* [kubernetes指南 kube-api server，写的很好](https://kubernetes.feisky.xyz/he-xin-yuan-li/index-1/apiserver)

## API的构成

来自[k8s官方API概览](https://kubernetes.io/docs/reference/using-api/api-overview/)

需要注意到，所有k8s中的组件都被称为API对象，并在API中有与之对应的元素。所有操作和通信都可以通过与API服务器的交互完成，即使是kubectl这个命令行工具也是实现了自己的API。

### API版本

有许多的API版本，比如`/api/v1`和`/apis/extensions/v1beta1`等，这些不同的API版本用来代替作用并重组资源表示。

放在API层而不是资源层或作用域层的原因是：

* 确保API代表了清楚和连续系统资源和行为
* 能够取得EOF和实验性API的控制

### API groups

API group通过REST路径指定，名称形如`/apis/$GROUP_NAME/$VERSION`。可以通过添加CRD或是对现有API进行聚合(aggregator)来扩展API

### OpenAPI

[OpenAPI](https://www.openapis.org/),，API的细节用OpenAPI进行管理（documented）。[k8s openapi仓库](https://github.com/kubernetes/kube-openapi)，允许使用openapi进行访问。

[Openapi是什么](https://swagger.io/docs/specification/about/)

## 使用API

[k8s官方教程](https://kubernetes.io/docs/reference/using-api/api-concepts/)

### API的专有名词

* resource type：指在URL中被使用的名字。
* kind：所有的resource type都有用JSON格式的具体表示
* collection：一种resource的资源组成的列表
* resource：单个resource资源的实例

资源被cluster限制`/apis/GROUP/VERSION/*`或是被命名空间限制`/apis/GROUP/VERSION/namespaces/NAMESPACE/*`

## 用API进行交互

如何是使用kubeadm进行的kubernetes配置且没有修改`--apiserver-bind-port`，那么默认的监听端口是6443，使用https协议。一般来说，使用`$USER/.kube/config`作为用户登陆用的证书是足够的。（或是环境变量配置`$KUBECONFIG`）

在kubectl下执行命令，如`kubectl --v=8 get pods`，可以输出其API指令与返回的HTTP信息。可以通过启用`kubectl proxy --port=8080`将api服务器端口转接到本地的8080端口，这样子可以通过`http://127.0.0.1:8080/api/...`来访问，通过合适的设置甚至可以实现外网进行访问。

如果想要证书的话，kubeadm的证书默认是存放在`--cert-dir="/etc/kubernetes/pki"`中的，实测里面有apiserver.crt/key，apiserver-etcd-client.crt/key，apiserver-kubelet-client.crt/key，ca.crt/key，front-proxy-ca.crt/key，front-proxy-client.crt/key,sa.pub/key，至于具体怎么用还不是很清楚。

### 使用curl

[参考](https://kubernetes.feisky.xyz/he-xin-yuan-li/index-1/apiserver)

Pod内部

```bash
# In Pods with service account.
$ TOKEN=$(cat /run/secrets/kubernetes.io/serviceaccount/token)
$ CACERT=/run/secrets/kubernetes.io/serviceaccount/ca.crt
$ curl --cacert $CACERT --header "Authorization: Bearer $TOKEN"  https://$KUBERNETES_SERVICE_HOST:$KUBERNETES_SERVICE_PORT/api
{
  "kind": "APIVersions",
  "versions": [
    "v1"
  ],
  "serverAddressByClientCIDRs": [
    {
      "clientCIDR": "0.0.0.0/0",
      "serverAddress": "10.0.1.149:443"
    }
  ]
}
```

Pod外部，亲测可用

```bash
# Outside of Pods.
$ APISERVER=$(kubectl config view | grep server | cut -f 2- -d ":" | tr -d " ")
$ TOKEN=$(kubectl describe secret $(kubectl get secrets | grep default | cut -f1 -d ' ') | grep -E '^token'| cut -f2 -d':'| tr -d '\t')
$ curl $APISERVER/api --header "Authorization: Bearer $TOKEN" --insecure
{
  "kind": "APIVersions",
  "versions": [
    "v1"
  ],
  "serverAddressByClientCIDRs": [
    {
      "clientCIDR": "0.0.0.0/0",
      "serverAddress": "10.0.1.149:443"
    }
  ]
}
```