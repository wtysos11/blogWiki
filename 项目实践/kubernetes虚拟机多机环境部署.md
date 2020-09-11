# kubernetes虚拟机多级环境部署与Istio的安装

标签：项目实践 kubernetes istio

计划在三台2C4G的机器上安装kubernetes。其中master和slave的系统均为Ubuntu18.04桌面版。（原本想用CentOS的，下载的时候就是这样了，费事改了）

分为以下几个步骤：
1. 访问谷歌镜像仓库gcr.io
2. 完成三台机器关于安装kubeadm的相关工作
3. 初始化k8s集群，验证通过。备份
4. 安装Istio

## 
