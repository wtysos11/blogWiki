# kubernetes互动课程学习

标签：kubernetes 实践

在师兄的介绍下发现了神奇的[awesome-kubernetes](https://github.com/ramitsurana/awesome-kubernetes)，里面有着很多的资源，其中之一就是非常有趣的interactive learning environments。其实官方英文教程里面也有一些简单的H5互动模块，但是真的太简单了。

k8s的搭建非常的麻烦，因为我个人的技术问题，不会使用国内的容器云，所以用的是VPN。搭建的这么麻烦，自然不想因为自己的原因导致故障或者重新搭建之类的事情发生，所以我还是比较喜欢这种在线的环境，至少出了问题可以随意的重启。

之前写了很长的一篇[kubernetes学习笔记](https://github.com/wtysos11/NoteBook/blob/master/%E5%BE%AE%E6%9C%8D%E5%8A%A1/Kubernetes%E5%AD%A6%E4%B9%A0%E7%AC%94%E8%AE%B0.md)，但是结构太散了，或者说它不是我的结构，而不是别人的结构。这篇文章的纵向线是magic sandbox的组织安排，但是横向线我会根据自己的学习内容进行扩展，所以看起来会比较杂，但我写起来、复习起来会很舒服。

本次笔记使用的是[Magic Sandbox](https://magicsandbox.com/)，主要的问题是延迟比较高，开香港节点的小飞机延迟在350ms左右，终端界面有明显的卡顿。

## 一些资源

由于主要涉及的应该是kubectl的使用，因此官方的[kubectl cheat sheet](https://kubernetes.io/docs/reference/kubectl/cheatsheet/)还是很值得使用的。

## magic sandbox课程

### 第1课 Introduction to magic sandbox

介绍界面，使用说明

[general yaml file format](https://learnxinyminutes.com/docs/yaml/)，yaml是kubernetes用来定义对象的配置文件所用的语言。

#### 一些yaml记号和规范

* literal block `|`，可以开启多行。
* folded block `>`，同上。

例子

```yaml
literal_block: |
    This entire block of text will be the value of the 'literal_block' key,
    with line breaks being preserved.

    The literal continues until de-dented, and the leading indentation is
    stripped.

        Any lines that are 'more-indented' keep the rest of their indentation -
        these lines will be indented by 4 spaces.
folded_style: >
    This entire block of text will be the value of 'folded_style', but this
    time, all newlines will be replaced with a single space.

    Blank lines, like above, are converted to a newline character.

        'More-indented' lines keep their newlines, too -
        this text will appear over two lines.
```

yaml的缩进使用空格而不是literal tab

#### Introduction to service : 疑问，kubernetes中的端口

[资料来源](https://matthewpalmer.net/kubernetes-app-developer/articles/kubernetes-ports-targetport-nodeport-service.html)

这里的疑问来源主要是Service的yaml文件中，有着port和targetPort两个条目，有点疑问。

一个是Pod的port list，在pod.spec.containers[].ports中定义，给出容器暴露的端口列表。这个列表是不用用户去指定的，容器自己会在对应的端口上进行监听。

一个是Service的port list，定义在service.spec.ports中，指定service的哪个端口映射到pod的哪个端口上。一个成功的请求可以从cluster外部发到节点的IP地址和Service的nodePort中，然后转发到service的port中，然后被pod的targetPod接收。

例子：

```yaml
kind: Service
apiVersion: v1
metadata:
  name: port-example-svc
spec:
  # Make the service externally visible via the node
  type: NodePort 

  ports:
    # Which port on the node is the service available through?
    - nodePort: 31234

    # Inside the cluster, what port does the service expose?
    - port: 8080

    # Which port do pods selected by this service expose?
    - targetPort: 

  selector:
    # ...
```

其中nodePort使得service在kubernetes cluster集群外可见，通过node的IP地址和该属性指定的nodePort来与该服务进行通信。

port在cluster内部暴露指定service。或者说，service在这个端口下变得可见，向这个端口发送信息会被转发到对应的pod上。

targetPort，是service向Pod发送的端口号，用户的应用需要在Pod中监听这个端口以实现与外部的通信。

这里就涉及到了service的工作原理，service解决的是逻辑上的一组Pod，外部如何访问它们的问题，因为Pod是可以随时撤销的，所以Service需要肩负起映射表的功效，通常是通过Label Selector进行管理

一些相关的名词

* endpoint：endpoint是K8s集群中的一个资源对象，存储在etcd中，用来记录service对应的所有pod的访问地址。service配置selector后，endpoint controller才会创建对应的endpoint对象，否则不创建。没有selector的service可以手动绑定endpoint。这也实现了kubernetes的服务发现机制。
* kube-proxy，负责service的实现，即实现k8s内部从pod到service和外部从node port到service的访问。实现了kubernetes的负载均衡机制。它的具体实现是基于iptables，也就是转发表的。主要职责是两大块，一块是侦听service更新事件，更新规则；一块是侦听endpoint更新事件，并更新规则，然后将包请求转发入endpoint对应的Pod中。

![kube-proxy-structure](https://github.com/wtysos11/NoteBook/blob/master/%E5%BE%AE%E6%9C%8D%E5%8A%A1/img/kube-proxy-structure.jpg?raw=true)

kube-proxy部分详情见这里[浅谈kubernetes service那些事](https://zhuanlan.zhihu.com/p/39909011)，是网易云在知乎上发布的，写的不错。官方的[service说明](https://kubernetes.io/zh/docs/concepts/services-networking/service/)也还可以。

#### Introduction to Deployment

Deployment通过RepolicaSets管理多个Pod。通过Deployment对象，用户描述了自己期望的理想状态，而控制器就会尽量让实际状态向理想状态靠拢。

```yaml
apiVersion: apps/v1beta2
kind: Deployment
metadata:
  name: nginx-deployment
  labels:
    app: nginx
spec:
  replicas: 3 #定义了期望3个Pod实例
  selector: #selector字段定义了什么样的Pod会在deployment的管理范围内
    matchLabels:
      app: nginx
  template:  #定义了创建的Pod的模版
    metadata:
      labels:
        app: nginx
    spec:
      containers:
      - name: nginx
        image: nginx:1.7.9
        ports:
        - containerPort: 80
```

#### updating deployment

如果想要更新deployment，可以通过kubectl来进行。即重复运行`kubectl apply -f deployment.yaml`，就可以更新设置。