# kubernetes虚拟机多级环境部署与Istio的安装

标签：项目实践 kubernetes istio

计划在三台2C4G的机器上安装kubernetes。其中master和slave的系统均为Ubuntu18.04桌面版。（原本想用CentOS的，下载的时候就是这样了，费事改了）

目前需要在实验室的服务器上进行部署，得了，最终结果也是一样的。

分为以下几个步骤：
1. 安装基本软件
2. 访问谷歌镜像仓库gcr.io
3. 完成三台机器关于安装kubeadm的相关工作
4. 部署相关应用（prometheus、grafana）
5. 使用helm安装Istio

参考文献：


## 安装基本软件

### 配置docker

### 配置kubectl、kubeadm、kubelet

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

在执行完第一步和第二步之后，可以直接使用kubeadm进行安装：`kubeadm init ---apiserver-advertise-address=192.168.1.1 --pod-network-cidr=192.168.0.0/24`

命令参考：[官方](https://kubernetes.io/zh/docs/reference/setup-tools/kubeadm/kubeadm-init/)，也可以在命令行下查看。

* apiserver-advertise-address：这个命令指定了监听的API地址。
* pod-network-cidr：规定了pod能够使用的IP地址段。我之前用的是16位子网掩码，但是现在给的子网就是24位掩码，我也不确定使用其他子网能不能行……先保险起见吧。
* kubernetes-version：指定kubeadm安装的kubernetes版本。这个是很重要的，因为默认情况下kubeadm会安装与它版本相同的kubernetes版本，而由于国内的网络问题，每次都需要重新下载一遍镜像，非常的麻烦。如果之后版本使用这个脚本，可以加上`--kubernetes-version=v1.19.2`