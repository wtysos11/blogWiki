# istio源码阅读

标签：istio kubernetes 源码阅读

使用istio1.0.6-release的源代码，具体见github。

## crd.yaml

具体位置：`install/kubernetes/helm/istio/templates/crds.yaml`。

crds是需要第一个采用的配置，不然后面的istio配置文件会因为缺少crd而失败（亲测，它真的不会自动加载这个文件，不知道是不是我的问题）

