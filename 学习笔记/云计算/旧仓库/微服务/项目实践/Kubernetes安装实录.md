# Kubernetes安装实录

标签：Kubernetes 安装记录

## 其他博客

在学习过程中陆续搜集到了一些不错的文章

* [Ubuntu16.04下用kubeadm安装kubernetes，英文](https://medium.com/@SystemMining/setup-kubenetes-cluster-on-ubuntu-16-04-with-kubeadm-336f4061d929)

## 原因

上学期的时候我在上潘茂林老师的服务计算课程的时候也接触到了容器化技术和Kubernetes，但是很可惜那时候因为计算机视觉的期末大作业过于难搞，因此没有学习。但是据说Kubernetes的安装极其艰难，因此先写下这篇安装实录，以备不时之需。安装环境是位于虚拟机中的Centos7，希望不要出什么大的问题。

因为不知道下一次看到这篇文章是什么时候，我将尽量与版本解耦合，如果到时候不存在这个工具了，那我也没什么办法。

## 文章安装

[修改主机名](https://www.jianshu.com/p/39d7000dfa47)方便后续的实验，区分主机。然后手动修改`/etc/hosts`中的127.0.0.1等原主机名字。

[Virtualbox配置centos7网络](https://www.jianshu.com/p/044fc0b85521)，经过实践，这篇文章讲的是没有问题的。配置ip除了可以按照文章上说的以外，也可以使用nmtui工具进行配置。

[从零开始搭建Kubernetes集群](https://www.jianshu.com/p/78a5afd0c597)主要参考这篇文章。此外还有使用rancher来进行安装的，也十分简单，这里不赘述。建议安装图形界面，然后使用ssh登陆进行端口映射登陆（这个链接里有提到怎么设置），这样子比较好安装增强功能（记得提前更新系统内核，可以先按照[官方docker安装](https://yeasy.gitbooks.io/docker_practice/install/centos.html)上装好依赖，装完docker之后装增强功能更容易，我本人成功率为100%。docker的版本需要注意[安装指定版本的docker](https://www.jianshu.com/p/d9dbf7e23722)）。

因为`kubeadm init`操作需要从谷歌的网站拉镜像，国内基本是不用想的了。[kubeadm搭建kubernetes](https://mritd.me/2016/10/29/set-up-kubernetes-cluster-by-kubeadm/)这里提供了两种方法，都可以。也可以使用[azure的国内gcr proxy](http://mirror.azure.cn/help/gcr-proxy-cache.html)用前一种方法提前拉取镜像，然后再修改对应的镜像名就好了(比如`docker tag imageid name:tag`)。

`kubectl taint nodes --all node-role.kubernetes.io/master-`可以允许Master节点作为工作节点，这样可以不使用minikube而创建一个单节点的K8S集群。

kubeadm init --pod-network-cidr=10.104.0.0/16 --apiserver-advertise-address=10.104.223.158

## 安装镜像

```

quay.io 中科大镜像 quay.mirrors.ustc.edu.cn，比如quay.mirrors.ustc.edu.cn/calico/cni
# quay.io/calico/cni 替换为
quay.mirrors.ustc.edu.cn/calico/cni

# gcr.io/namespace/image_name:image_tag 
# 替换为
# gcr.mirrors.ustc.edu.cn/namespace/image_name:image_tag 
gcr.mirrors.ustc.edu.cn/kubernetes-helm/tiller
# Kubernetes官方教程经常用到k8s.gcr.io, 
# 相应的 k8s.gcr.io 等同于 gcr.io/google-containers/
# 因此 k8s.gcr.io/busybox 等价于
gcr.mirrors.ustc.edu.cn/google-containers/busybox



docker pull gcr.azk8s.cn/google_containers/kube-apiserver:v1.13.4
docker pull gcr.azk8s.cn/google_containers/kube-controller-manager:v1.13.4
docker pull gcr.azk8s.cn/google_containers/kube-scheduler:v1.13.4
docker pull gcr.azk8s.cn/google_containers/kube-proxy:v1.13.4
docker pull gcr.azk8s.cn/google_containers/pause:3.1
docker pull gcr.azk8s.cn/google_containers/etcd:3.2.24
docker pull gcr.azk8s.cn/google_containers/coredns:1.2.6

docker pull gcr.azk8s.cn/google_containers/kubernetes-dashboard-amd64:v1.10.1

k8s.gcr.io/kubernetes-dashboard-amd64:v1.10.1

k8s.gcr.io/kube-apiserver:v1.13.4
k8s.gcr.io/kube-controller-manager:v1.13.4
k8s.gcr.io/kube-scheduler:v1.13.4
k8s.gcr.io/kube-proxy:v1.13.4
k8s.gcr.io/pause:3.1
k8s.gcr.io/etcd:3.2.24
k8s.gcr.io/coredns:1.2.6
```

## 常见问题解决

执行完`kubeadm init`后忘记token，或者token24小时过期后还想再加入新的节点，可以参考[stackoverflow](https://stackoverflow.com/questions/40009831/cant-find-kubeadm-token-after-initializing-master)上的这篇文章，根据官方文档来进行操作（第二个）

执行`kubeadm init`的时候，可能会报错` /proc/sys/net/bridge/bridge-nf-call-iptables contents are not set to 1`，这里[这篇文章](https://segmentfault.com/a/1190000009374196)给出了很好的解决方法：

`echo 1 > /proc/sys/net/bridge/bridge-nf-call-iptables`与`echo 1 > /proc/sys/net/bridge/bridge-nf-call-ip6tables`

因故或者意外删除节点的时候需要执行`kubeadm reset`，如果出现Port 10250被占用的情况也需要进行处理，如[这个](https://stackoverflow.com/questions/41732265/how-to-use-kubeadm-to-create-kubernetes-cluster)

如何debug，当节点出现异常的时候，[stackoverflow问题](https://stackoverflow.com/questions/47107117/how-to-debug-when-kubernetes-nodes-are-in-not-ready-state)

问题：出现了下面的情况，且无法连接服务器

```
Unable to connect to the server: x509: certificate signed by unknown authority (possibly because of "crypto/rsa: verification error" while trying to verify candidate authority certificate "kubernetes")
```

解决方案：依据[官方文档](https://k8smeetup.github.io/docs/setup/independent/troubleshooting-kubeadm/)，执行下面的命令来覆盖"admin"用户的默认kubeconfig文件（很多情况都可以由这个解决）

```
mv  $HOME/.kube $HOME/.kube.bak
mkdir -p $HOME/.kube
sudo cp -i /etc/kubernetes/admin.conf $HOME/.kube/config
sudo chown $(id -u):$(id -g) $HOME/.kube/config
```

问题：kubectl指令出现`The connection to the server localhost:8080 was refused`的情况：

```
sudo cp /etc/kubernetes/admin.conf $HOME/
sudo chown $(id -u):$(id -g) $HOME/admin.conf
export KUBECONFIG=$HOME/admin.conf
```

问题：coredns频繁掉线，且出现`HINFO: unreachable backend: read udp 192.168.1.2:53484->10.8.4.4:53: i/o timeout`，后面的地址是配置的DNS nameserver，见`etc/resolv.conf`。详见[stackoverflow](https://stackoverflow.com/questions/52837574/coredns-couldnt-reach-to-host-nameserver)

问题：重启服务器之后出现`The connection to the server 10.104.223.158:6443 was refused - did you specify the right host or port?`

## 配置kubernetes-dashboard

也可以参考[这个](https://jimmysong.io/posts/kubernetes-dashboard-upgrade/)（虽然我用的不是nodePort而是kube proxy，这一条`kubectl proxy --address='0.0.0.0' --accept-hosts='^*$'`）

或是配置service nodeport，端口31064.

首先是镜像问题，因为在国内无法访问到

访问`http://203.195.219.185:8001/api/v1/namespaces/kube-system/services/https:kubernetes-dashboard:/proxy/#!/login`登陆，登陆前要先按照[kubernetes Dashboards中的身份认证详解](https://jimmysong.io/posts/kubernetes-dashboard-upgrade/#%E4%BD%BF%E7%94%A8-kubeconfig)中所说的，制作token。

得到token，并获取token。

将token经过base64解码后作为token登陆，既可以拥有管理员权限操作整个k8s集群中的对象。或者将这串token进行base64解码后，加到admin用户的Kubeconfig文件中。

token例子

```
eyJhbGciOiJSUzI1NiIsImtpZCI6IiJ9.eyJpc3MiOiJrdWJlcm5ldGVzL3NlcnZpY2VhY2NvdW50Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9uYW1lc3BhY2UiOiJrdWJlLXN5c3RlbSIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VjcmV0Lm5hbWUiOiJhZG1pbi10b2tlbi10NjhtOSIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VydmljZS1hY2NvdW50Lm5hbWUiOiJhZG1pbiIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VydmljZS1hY2NvdW50LnVpZCI6IjU5NTRjZTcxLTQwMDktMTFlOS04ZDZiLTUyNTQwMDEyMThlYyIsInN1YiI6InN5c3RlbTpzZXJ2aWNlYWNjb3VudDprdWJlLXN5c3RlbTphZG1pbiJ9.aCaemFbyBCm_uluT5MIiME_tAo4XE_S6d9WNj4bQIxvx_XyXALnAhZ13ysUBONx2NddoxxRAe_mF6GsNOZ9_RDdU65veJwkI8sp3AoJn5gtBysILiW1Sv1O4BV4Nnk-GVrBJysW5wcHFMN0cgC00F1iI3sry-KQRm8lV49Z3eiEn4AbuT4GqcprNGt0dN_WkiOWar-iJduwX6_xXsQ8Zf1iBmFWrZY-TAiok0nKtCn6QPex5FoZKJ5fWtJPJHvnaioKmY3lrYTHXrbz_ql08N9EHuSLUdhLdkZI-F2RzN3fY47a437h8mw3hJQhUyHUzgrq8iR8f6bOWyBjxmNW3SA
```

### 遇到的问题

kubernetes安装dashboard后输入token没反应。后面才发现没有改官网的yaml文件。像[这里](https://blog.csdn.net/spring_trees/article/details/80825473)

## 安装helm

helm client的安装就不说了，很简单。

安装完client后开始安装server: helm tiller。使用`helm init --upgrade`命令即可开始安装，但因为安装的时候默认会向拉取gcr.io/kubernetes-helm/tiller镜像，国内怎么可能拉得了，所以可以使用阿里云的`helm init --upgrade -i registry.cn-hangzhou.aliyuncs.com/google_containers/tiller:v2.13.0 --stable-repo-url https://kubernetes.oss-cn-hangzhou.aliyuncs.com/charts`，其中这个tiller的版本可以自行更改

在过程中可能会遇到`Error: configmaps is forbidden: User "system:serviceaccount:kube-system:default" cannot list resource "configmaps" in API group "" in the namespace "kube-system"`这样的错误信息，这时可以参考这篇issue[helm的权限问题](https://github.com/helm/helm/issues/3130)中noprom提到的，重新安装helm tiller，建立高权限用户来解决。（虽然他执行到这一步之后还是有问题，但是我自己试过一遍之后是可以的）

## 使用helm安装istio

参考[官方教程](https://istio.io/zh/docs/setup/kubernetes/helm-install/)，使用Helm的helm template来安装。

PS：官方教程有点问题，可以参考这篇[文章](https://blog.csdn.net/c5113620/article/details/82765686)，注意修改istio的版本。

在下载好的压缩包内执行`kubectl apply -f install/kubernetes/helm/istio/templates/crds.yaml`，安装istio定义的资源。然后执行`kubectl apply -f install/kubernetes/istio-demo-auth.yaml`配置用户。

因为使用的不是docker的官方镜像，所以比较慢，需要耐心等待（10分钟内基本完成），或者自己先行下载好镜像。

## 新的节点如何加入到集群中

因为网络环境复杂，很容易出现像是虚拟机系统错误、崩溃，云服务器被人入侵等情况，导致机器不可用甚至重装系统的情况。这时候如何快速恢复可用。

1. 更新系统软件，`yum update`
2. 重新安装docker指定版本，首先添加docker仓库源，推荐使用官方给的`sudo yum-config-manager --add-repo https://mirrors.ustc.edu.cn/docker-ce/linux/centos/docker-ce.repo`。然后用`yum list docker-ce --showduplicates|sort -r  `显示所有版本后直接找到对应版本进行安装，比如`yum install docker-ce-17.09.0.ce -y`。

3. 修改`/etc/docker/daemon.json`，添加下面的信息来配置镜像。

```json
{
  "registry-mirrors": [
    "https://registry.docker-cn.com",
    "https://docker.mirrors.ustc.edu.cn/"
  ]
}
```

再然后，启动docker服务：`systemctl start docker & systemctl enable docker`，检查是否成功安装。如果成功，进入下一步。

4. 关闭swap、防火墙（centos7下防火墙使用`firewall-cmd --state`检查状态，用` systemctl stop firewalld & systemctl disable firewalld`进行关闭），关闭selinux：`setenforce 0`。
5. 配置k8s的yum源。

```bash
cat <<EOF > /etc/yum.repos.d/kubernetes.repo

[kubernetes]

name=Kubernetes

baseurl=http://mirrors.aliyun.com/kubernetes/yum/repos/kubernetes-el7-x86_64

enabled=1

gpgcheck=0

repo_gpgcheck=0

gpgkey=http://mirrors.aliyun.com/kubernetes/yum/doc/yum-key.gpg

        http://mirrors.aliyun.com/kubernetes/yum/doc/rpm-package-key.gpg

EOF
```

安装k8s组件`yum install -y kubelet kubeadm kubectl`

6. 启动kubelet，`systemctl enable kubelet && systemctl start kubelet`
7. 下载gcr镜像，详见上面的说明。然后修改下载好地镜像的tag。
8. 生成token，详见[stackoverflow](https://stackoverflow.com/questions/40009831/cant-find-kubeadm-token-after-initializing-master)和[这里](https://blog.csdn.net/mailjoin/article/details/79686934)

openssl x509 -pubkey -in /etc/kubernetes/pki/ca.crt | openssl rsa -pubin -outform der 2>/dev/null | openssl dgst -sha256 -hex | sed 's/^.* //'

c8710bca46713969d4fd099050e80151d7b90bc172ee3649e6d79da02c2d61c1

kubeadm join --token ku38an.j2brd6sshybrfqg0 --discovery-token-ca-cert-hash sha256:c8710bca46713969d4fd099050e80151d7b90bc172ee3649e6d79da02c2d61c1  10.104.235.215:6443 --skip-preflight-checks