# calico安装文件阅读

标签：calico kubernetes 源码阅读

[文件地址](https://docs.projectcalico.org/v3.6/getting-started/kubernetes/installation/hosted/kubernetes-datastore/calico-networking/1.7/calico.yaml)

注：比较特殊的是`POD_CIDR`这个属性，指代的是网络的运行子网，如果与`192.168.0.0/16`不相同需要自行修改。

[api object参考文档](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.10/)

## calico-config

这个文件的类型是`ConfigMap`，比较重要的字段就是`data`了。在data中值以key-value的形式存放，并且给予容器作为配置信息。

值得一提的是yaml的文件格式，yaml支持三种数据结构：对象、数组、纯量。对象是被花括号括起来的值，也可以作为键值对的值，称为行内对象。数值是一组以`-`开头的行，构成一个数组。

特殊记号`|`表示保留换行符，也可以使用`>`进行折叠换行。`+`表示保留文字块末尾的换行，`-`表示删除字符串末尾的换行。也就是说，`|`和`>`后的都是多行连续字符串。

锚点`&`和别名`*`，可以用来引用。

### API对象与资源对象：CustomResourceDefinition

kind:[CustomResourceDefinition](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.10/#customresourcedefinition-v1beta1-apiextensions-k8s-io)

crd对象代表了API服务器暴露出的一种资源类型，它的名称必须形如`<.spec.name>.<.spec.group>`,比如`virtualservices.networking.istio.io`

spec字段指定了用户想要这个对象展现出的信息，这里指`CustomResourceDefinitionSpec`，包含group、names(CustomResourceDefinitionNames，表示对象的各种状态时对应的名字)、scope（作用域是全局的还是命名空间内）、subresources、validation、version等。

比较重要且容器弄错的name字段中的属性是categories，在其他API中这个字段是一个数组(string array)，表示这个资源属于哪些资源。

## calico-kube-controllers

source:calico/templates/rbac.yaml，属于calico的rbac系统的组成部分之一，kind: ClusterRole

### API对象与资源对象：ClusterRole

clusterrole是在cluster层面上，PolicyRules的逻辑分组。可以被RoleBinding或ClusterRoleBinding引用。

重要的字段是rules（一个PolicyRule array），保留这个ClusterRole的所有PolicyRules。此外还有可选的aggregationRule，描述了如何为这个ClusterRole构建Rules。

### API对象：PolicyRule

[reference](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.10/#policyrule-v1beta1-rbac-authorization-k8s-io)

PolicyRule并不是资源对象，不能够在kubernetes服务器上找到。它保有着能够描述一个policy rule的信息，但是并不保有这个rule适用于哪个对象或是哪个作用域。

在本对象中所使用的字段：

* apiGroups(string array)：apiGroup的集合，每个apiGroup都是一群资源的定义（实际测试过，与CustomResourceDefinition中spec/group字段的值是一样的。也就是说，只带的是一个资源群）。原文的说法是如果指定了多个apiGroups，那么policyRule的动作作用于任何一个apiGroups所枚举出的资源类型都是被允许的。
* resources(string array)：代表的是动作所对应的资源类型
* verbs(string array)：一系列作用于资源或属性的动作

### API对象与资源对象：ClusterRoleBinding

ClusterRoleBinding指向一个ClusterRole，但是并不包含它。

对象所含有的字段：

* roleRef(RoleRef):RoleRef只能够引用全局作用域的一个ClusterRole。
* roleRef->apiGroup(string)：被引用资源所属的apiGroup（经检验，与apiVersion斜杠前的部分完全一致）
* roleRef->kind：被引用资源的类型
* roleRef->name：被引用资源的名称
* subject(Subject array)：保留了role所申请的对象的引用(hold references to the objects the role applied to)

