# kubernetes虚拟机多级环境部署与Istio的安装

标签：项目实践 kubernetes istio

计划在三台2C4G的机器上安装kubernetes。其中master和slave的系统均为Ubuntu18.04桌面版。（原本想用CentOS的，下载的时候就是这样了，费事改了）

目前需要在实验室的服务器上进行部署，得了，最终结果也是一样的。

分为以下几个步骤：
1. 安装基本软件
2. 访问谷歌镜像仓库gcr.io
3. 完成三台机器关于安装kubeadm的相关工作
4. 安装Istio
5. 部署相关应用（prometheus、grafana）

参考文献：
* [掘金](https://juejin.im/post/6844903908326932493)
* [cnblog](https://www.cnblogs.com/tylerzhou/p/10971336.html)，这一篇是我后来看到的，写的更好一些。

## 安装基本软件

如果使用包管理软件，请务必检查自身的包管理软件为最新。

以Ubuntu的apt为例，首先请务必使用最新的镜像仓库（最好手动弄一下），然后再执行一遍`sudo apt-get update`来更新。

### 配置docker

略

### 配置kubectl、kubeadm、kubelet

在执行命令之前，请首先使用检查是否能够访问到指定的包。以apt为例，使用`apt-cache policy <package>`可以检查远程仓库中包的版本。

* `apt-get install kubectl kubeadm kubelet`

### 系统配置

* 关闭防火墙
* 关闭SELinux
* 关闭swap

## 下载kubeadm所需镜像

### 从阿里云下载

首先使用`kubeadm config images list`列出kubeadm所需要的所有镜像

```
k8s.gcr.io/kube-apiserver:v1.19.2
k8s.gcr.io/kube-controller-manager:v1.19.2
k8s.gcr.io/kube-scheduler:v1.19.2
k8s.gcr.io/kube-proxy:v1.19.2
k8s.gcr.io/pause:3.2
k8s.gcr.io/etcd:3.4.13-0
k8s.gcr.io/coredns:1.7.0
```

即所有的系统镜像为v1.19.2版本，pause3.2,etcd3.4.13,coredns1.7.0

[这篇博客](https://juejin.im/post/6844903908326932493)提供了一个脚本

```bash
#########################################################################
# File Name: pull_master_image.sh
# Description: pull_master_image.sh
# Author: zhangyi
# mail: 450575982@qq.com
# Created Time: 2019-07-31 21:38:14
#########################################################################
#!/bin/bash
kube_version=:v1.19.2
kube_images=(kube-proxy kube-scheduler kube-controller-manager kube-apiserver)
addon_images=(etcd:3.4.13-0 coredns:1.7.0 pause:3.2)

for imageName in ${kube_images[@]} ; do
  docker pull registry.cn-hangzhou.aliyuncs.com/google_containers/$imageName-amd64$kube_version
  docker image tag registry.cn-hangzhou.aliyuncs.com/google_containers/$imageName-amd64$kube_version k8s.gcr.io/$imageName$kube_version
  docker image rm registry.cn-hangzhou.aliyuncs.com/google_containers/$imageName-amd64$kube_version
done

for imageName in ${addon_images[@]} ; do
  docker pull registry.cn-hangzhou.aliyuncs.com/google_containers/$imageName
  docker image tag registry.cn-hangzhou.aliyuncs.com/google_containers/$imageName k8s.gcr.io/$imageName
  docker image rm registry.cn-hangzhou.aliyuncs.com/google_containers/$imageName
done
```

这个脚本应该是没有问题的，有可能因为阿里云那边的容器镜像发生改变导致脚本需要变化，但大体思路是没错的。

可能的错误信息：
* 如果使用`sh`命令运行此脚本，可能会遇到`pullImage.sh: 3: pullImage.sh: Syntax error: "(" unexpected`。参考[这篇回答](https://unix.stackexchange.com/questions/253892/syntax-error-unexpected-when-creating-an-array)，该问题的原因是因为使用`sh xx.sh`的时候是使用sh shell，是没有数组的。正确的执行方式是`bash xx.sh`或者如回答中一样使用`./xx.sh`(需要注意权限)

## 使用kubeadm安装kubernetes集群

### 安装集群主体

在执行完第一步和第二步之后，可以直接使用kubeadm进行安装：`sudo kubeadm init --apiserver-advertise-address=192.168.1.1 --pod-network-cidr=10.244.0.0/16 --token-ttl=0`（与之相反，执行`kubeadm reset`就可以尽量还原`kubeadm init`或者`kubeadm join`所带来的影响）

命令参考：[官方](https://kubernetes.io/zh/docs/reference/setup-tools/kubeadm/kubeadm-init/)，也可以在命令行下查看。

* apiserver-advertise-address：这个参数指定了监听的API地址。若没有设置，则使用默认网络接口。
* apiserver-bind-port：这个参数指定了API服务器暴露出的端口号，默认是6443。
* pod-network-cidr：规定了pod能够使用的IP地址段。我之前用的是16位子网掩码，但是现在给的子网就是24位掩码，我也不确定使用其他子网能不能行……先保险起见吧。
* kubernetes-version：指定kubeadm安装的kubernetes版本。这个是很重要的，因为默认情况下kubeadm会安装与它版本相同的kubernetes版本，而由于国内的网络问题，每次都需要重新下载一遍镜像，非常的麻烦。如果之后版本使用这个脚本，可以加上`--kubernetes-version=v1.19.2`
* image-repository：默认是"k8s.gcr.io"。我觉得如果修改这个可以不用像之前那样从阿里云下载下来后手动tag，不过没有尝试。
* token-ttl：令牌被删除前的时间，默认是24h。kubeadm初始化完毕后会生成一个令牌，让其他节点能够加入集群，过时之后这个令牌会自动删除。如果设置为0之后令牌就永不过期。

这一步的难点在于如何设置pod-network-cidr，参数的[作用](https://blog.csdn.net/shida_csdn/article/details/104334372)。根据[官方教程](https://kubernetes.io/zh/docs/setup/production-environment/tools/kubeadm/create-cluster-kubeadm/)的说法，Pod网络与任何主机网络不得有重叠。但是目前我看到的很多教程都是在主机局域网络下构建的。比如说三台主机都在`192.168.1.1/16`子网，而pod网络也在同样的子网下。

这个参数的设置似乎与所使用的CNI有关系：
* [flannel](https://coreos.com/flannel/docs/latest/kubernetes.html)，要求的参数为`--pod-network-cidr=10.244.0.0/16`
* [calico](https://docs.projectcalico.org/getting-started/kubernetes/quickstart)，要求的参数为`--pod-network-cidr=192.168.0.0/16`

本文采用flannel，一个很重要的原因是因为服务器的子网与pod的子网部分重叠，可能存在风险，以及flannel似乎更容易部署一些。

产生的日志：

```bash
stack@workflow-1:~/k8sTest$ sudo kubeadm init --apiserver-advertise-address=192.168.1.1 --pod-network-cidr=10.244.0.0/16 --token-ttl=0
W1012 19:18:16.455833   86405 configset.go:348] WARNING: kubeadm cannot validate component configs for API groups [kubelet.config.k8s.io kubeproxy.config.k8s.io]
[init] Using Kubernetes version: v1.19.2
[preflight] Running pre-flight checks
[preflight] Pulling images required for setting up a Kubernetes cluster
[preflight] This might take a minute or two, depending on the speed of your internet connection
[preflight] You can also perform this action in beforehand using 'kubeadm config images pull'
[certs] Using certificateDir folder "/etc/kubernetes/pki"
[certs] Generating "ca" certificate and key
[certs] Generating "apiserver" certificate and key
[certs] apiserver serving cert is signed for DNS names [kubernetes kubernetes.default kubernetes.default.svc kubernetes.default.svc.cluster.local workflow-1] and IPs [10.96.0.1 192.168.1.1]
[certs] Generating "apiserver-kubelet-client" certificate and key
[certs] Generating "front-proxy-ca" certificate and key
[certs] Generating "front-proxy-client" certificate and key
[certs] Generating "etcd/ca" certificate and key
[certs] Generating "etcd/server" certificate and key
[certs] etcd/server serving cert is signed for DNS names [localhost workflow-1] and IPs [192.168.1.1 127.0.0.1 ::1]
[certs] Generating "etcd/peer" certificate and key
[certs] etcd/peer serving cert is signed for DNS names [localhost workflow-1] and IPs [192.168.1.1 127.0.0.1 ::1]
[certs] Generating "etcd/healthcheck-client" certificate and key
[certs] Generating "apiserver-etcd-client" certificate and key
[certs] Generating "sa" key and public key
[kubeconfig] Using kubeconfig folder "/etc/kubernetes"
[kubeconfig] Writing "admin.conf" kubeconfig file
[kubeconfig] Writing "kubelet.conf" kubeconfig file
[kubeconfig] Writing "controller-manager.conf" kubeconfig file
[kubeconfig] Writing "scheduler.conf" kubeconfig file
[kubelet-start] Writing kubelet environment file with flags to file "/var/lib/kubelet/kubeadm-flags.env"
[kubelet-start] Writing kubelet configuration to file "/var/lib/kubelet/config.yaml"
[kubelet-start] Starting the kubelet
[control-plane] Using manifest folder "/etc/kubernetes/manifests"
[control-plane] Creating static Pod manifest for "kube-apiserver"
[control-plane] Creating static Pod manifest for "kube-controller-manager"
[control-plane] Creating static Pod manifest for "kube-scheduler"
[etcd] Creating static Pod manifest for local etcd in "/etc/kubernetes/manifests"
[wait-control-plane] Waiting for the kubelet to boot up the control plane as static Pods from directory "/etc/kubernetes/manifests". This can take up to 4m0s
[apiclient] All control plane components are healthy after 26.502830 seconds
[upload-config] Storing the configuration used in ConfigMap "kubeadm-config" in the "kube-system" Namespace
[kubelet] Creating a ConfigMap "kubelet-config-1.19" in namespace kube-system with the configuration for the kubelets in the cluster
[upload-certs] Skipping phase. Please see --upload-certs
[mark-control-plane] Marking the node workflow-1 as control-plane by adding the label "node-role.kubernetes.io/master=''"
[mark-control-plane] Marking the node workflow-1 as control-plane by adding the taints [node-role.kubernetes.io/master:NoSchedule]
[bootstrap-token] Using token: a8svth.m8zu4mdj60m6zjzd
[bootstrap-token] Configuring bootstrap tokens, cluster-info ConfigMap, RBAC Roles
[bootstrap-token] configured RBAC rules to allow Node Bootstrap tokens to get nodes
[bootstrap-token] configured RBAC rules to allow Node Bootstrap tokens to post CSRs in order for nodes to get long term certificate credentials
[bootstrap-token] configured RBAC rules to allow the csrapprover controller automatically approve CSRs from a Node Bootstrap Token
[bootstrap-token] configured RBAC rules to allow certificate rotation for all node client certificates in the cluster
[bootstrap-token] Creating the "cluster-info" ConfigMap in the "kube-public" namespace
[kubelet-finalize] Updating "/etc/kubernetes/kubelet.conf" to point to a rotatable kubelet client certificate and key
[addons] Applied essential addon: CoreDNS
[addons] Applied essential addon: kube-proxy

Your Kubernetes control-plane has initialized successfully!

To start using your cluster, you need to run the following as a regular user:

  mkdir -p $HOME/.kube
  sudo cp -i /etc/kubernetes/admin.conf $HOME/.kube/config
  sudo chown $(id -u):$(id -g) $HOME/.kube/config

You should now deploy a pod network to the cluster.
Run "kubectl apply -f [podnetwork].yaml" with one of the options listed at:
  https://kubernetes.io/docs/concepts/cluster-administration/addons/

Then you can join any number of worker nodes by running the following on each as root:

kubeadm join 192.168.1.1:6443 --token a8svth.m8zu4mdj60m6zjzd \
    --discovery-token-ca-cert-hash sha256:a98361e5d879fe14734236285b8ad28d0a4d3d1470bd424194011bee41ee8c9e 
```

在看到如上的命令之后，需要按照它的说明执行这三条命令，其作用是将kubectl所需要的配置文件拉到用户目录下并设置访问权限。（不执行的话kubectl是无法正常工作的）

```bash
mkdir -p $HOME/.kube
sudo cp -i /etc/kubernetes/admin.conf $HOME/.kube/config
sudo chown $(id -u):$(id -g) $HOME/.kube/config
```

之后，执行`kubectl get pods --all-namespaces`(或者`-n=kube-system`)，就可以看到kubernetes的系统节点了。

```
NAME                                 READY   STATUS    RESTARTS   AGE
coredns-f9fd979d6-7m7bb              1/1     Running   0          2m26s
coredns-f9fd979d6-9xr4b              1/1     Running   0          2m26s
etcd-workflow-1                      1/1     Running   0          2m30s
kube-apiserver-workflow-1            1/1     Running   0          2m30s
kube-controller-manager-workflow-1   1/1     Running   0          2m30s
kube-proxy-pdpkp                     1/1     Running   0          2m26s
kube-scheduler-workflow-1            1/1     Running   0          2m30s
```

这里的CoreDNS需要在安装完CNI插件之后才会正常运行，Pending是正常情况。这里是Running，我觉得是因为我第一次的时候没有删除干净。

### 安装flannel

参考[flannel官方教程](https://github.com/coreos/flannel/blob/master/Documentation/kubernetes.md)，对于1.7版本以上的kubernetes，可以直接用`kubectl apply -f https://raw.githubusercontent.com/coreos/flannel/master/Documentation/kube-flannel.yml`。当然，国内网络，懂得都懂，我就基本上没有直连成功过。直接手动弄下来搞也可以。

日志记录

```bash
stack@workflow-1:~/k8sTest$ kubectl apply -f kube-flannel.yml 
podsecuritypolicy.policy/psp.flannel.unprivileged created
Warning: rbac.authorization.k8s.io/v1beta1 ClusterRole is deprecated in v1.17+, unavailable in v1.22+; use rbac.authorization.k8s.io/v1 ClusterRole
clusterrole.rbac.authorization.k8s.io/flannel created
Warning: rbac.authorization.k8s.io/v1beta1 ClusterRoleBinding is deprecated in v1.17+, unavailable in v1.22+; use rbac.authorization.k8s.io/v1 ClusterRoleBinding
clusterrolebinding.rbac.authorization.k8s.io/flannel created
serviceaccount/flannel created
configmap/kube-flannel-cfg created
daemonset.apps/kube-flannel-ds created
```

可能出现的问题：
* flannel需要镜像`quay.io/coreos/flannel:v0.13.0-rc2`，之前的时候我的虚拟机是无法直接从`quay.io`下拉镜像，因此有可能需要曲线救国。

### 连接其他节点

将主节点装完之后就可以安装其他节点了（也可以在之前一起安装）。其他节点的前置条件需要像主节点一样，包括系统配置、软件安装等等。

比较麻烦的一点是从节点也需要安装那一些镜像，如果我没有理解错的话。这就需要很久了……

命令就是之前弹出来的那个

```
sudo kubeadm join 192.168.1.1:6443 --token a8svth.m8zu4mdj60m6zjzd \
    --discovery-token-ca-cert-hash sha256:a98361e5d879fe14734236285b8ad28d0a4d3d1470bd424194011bee41ee8c9e 
```

可能出现的问题：
* 如果安装过比较多次的话，会发现一些奇怪的报错信息，比如端口被占用之类的。请记得执行`kubeadm reset`

## 安装istio

原本计划是用helm来安装istio，但是最新版本的[官方文档](https://istio.io/latest/zh/docs/setup/install/helm/)弃用了helm。

### 下载

首先到官方给的[发布页面](https://github.com/istio/istio/releases/tag/1.7.3)下载（这个链接是1.7.3版本）。

然后使用`tar -zxvf`解压

### 安装

按照官方说明[文档](https://istio.io/latest/zh/docs/setup/getting-started/#download)

* 把文件夹下的`istioctl`所在目录加入PATH变量之中。（我没有把istioctl丢到`/usr/bin`下，因为它这个安装过程有可能要用到其他目录中的文件）
* 执行`istioctl manifest install --set profile=demo`

```
stack@workflow-1:~/k8sTest/istio-1.7.3/bin$ istioctl manifest install --set profile=demo
Detected that your cluster does not support third party JWT authentication. Falling back to less secure first party JWT. See https://istio.io/docs/ops/best-practices/security/#configure-third-party-service-account-tokens for details.
✔ Istio core installed                                                               
✔ Istiod installed                                                                   
✔ Egress gateways installed                                                          
✔ Ingress gateways installed                                                         
✔ Installation complete 
```

结果：
```
stack@workflow-1:~/k8sTest/istio-1.7.3/bin$ kubectl get pods -n=istio-system
NAME                                    READY   STATUS    RESTARTS   AGE
istio-egressgateway-8556f8c8dc-rj922    1/1     Running   0          2m42s
istio-ingressgateway-589d868684-pf7g4   1/1     Running   0          2m42s
istiod-86d65b6959-lch8x                 1/1     Running   0          3m54s
stack@workflow-1:~/k8sTest/istio-1.7.3/bin$ kubectl get svc -n istio-system
NAME                   TYPE           CLUSTER-IP      EXTERNAL-IP   PORT(S)                                                                      AGE
istio-egressgateway    ClusterIP      10.108.56.199   <none>        80/TCP,443/TCP,15443/TCP                                                     3m7s
istio-ingressgateway   LoadBalancer   10.96.27.251    <pending>     15021:32083/TCP,80:31588/TCP,443:31926/TCP,31400:30951/TCP,15443:31996/TCP   3m7s
istiod                 ClusterIP      10.111.3.121    <none>        15010/TCP,15012/TCP,443/TCP,15014/TCP,853/TCP                                4m20s
```

这个步骤真的是比我当时（0.2.x）要简单多了。我记得我试了好多次都没有成功，都有心理阴影了。才过了一年多就有了这么大的变化，软件行业真的是日新月异啊。

## 部署prometheus和grafana

参考资料：
* [简书-k8s安装prometheus+grafana](https://www.jianshu.com/p/ac8853927528)
* [prometheus-github](https://github.com/prometheus/prometheus)
* [掘金-k8s监控 安装prometheus](https://juejin.im/post/6844903908251451406) 
* [promethues官方-配置文件规范](https://prometheus.io/docs/prometheus/latest/configuration/configuration/)

istio之前有一套prometheus的系统，因此我一开始并没有考虑安装。但是在[最新版本1.7](https://istio.io/latest/docs/ops/integrations/prometheus/)中它并没有自带了，只是在addon里面加上了一个prometheus的样本用于验证（不过似乎也够用？）。

[istio/addon](https://github.com/istio/istio/tree/master/samples/addons)

可以直接使用`kubectl apply -f samples/addons`进行安装。当然如果之前是自定义安装istioctl的话那这一步就需要看一下官方文档怎么说了。

可能出现的问题：
* 在K8s 1.19版本情况下可能会出现如[issues](https://github.com/istio/istio/issues/27417)所示的情况，因为某些未知的原因k8s无法识别出在同一个文件内配置的配置项，因此必须要执行两次`kubectl apply -f kiali.yaml`

### 暴露服务

此时prometheus等服务还是只能够在内网进行访问，需要在外部网络暴露服务，不然就只能做反向代理了。

参考：
* [Kubernetes Ingress实战(二)：使用Ingress将第一个HTTP服务暴露到集群外部](https://blog.frognew.com/2018/06/kubernetes-ingress-2.html)
* [istio-ingress control](https://istio.io/latest/zh/docs/tasks/traffic-management/ingress/ingress-control/)

不过istio自己有一个ingressgateway，可能又与普通的K8s集群有所区别。需要等我把组会需要的东西做完了，才能把这个坑给填上。
