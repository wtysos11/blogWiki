# 如何进入k8s平台下正在运行的容器的终端+精确测算k8s+istio下pod的CPU占用率

参考资料:
* [官方文档](https://kubernetes.io/zh/docs/tasks/debug-application-cluster/get-shell-running-container/)

## 解决的问题

* 如何进入运行中的容器？或者说，如何打开运行中容器的bash界面。
* 如何让容器长期运行？或者说，如何让ubuntu镜像长期运行？

## 实验目标

要在istio网络下建立一个基于ubuntu的服务器，然后通过像这个服务器发送流量的方式来弄清楚prometheus与华为云提供的监控哪一个关于CPU占用率是更准确地。

为了实现这个目的，需要在k8s下进入正在运行的容器的终端

## 知识

可以使用`kubectl exec -it shell-demo -- /bin/bash`来进入名为`shell-demo`的容器中执行命令`/bin/bash`

双破折号`--`用于将要传递给命令的参数与kubectl的参数分开。

在普通的命令窗口（不是shell）可以执行单个命令。比如`kubectl exec shell-demo env`相当于在容器内执行`env`

在Pod包含多个容器时打开shell，可以使用`--container`或者`-c`在`kubectl exec`命令中指定容器。

## 构造镜像

我们基于bookinfo应用的productpage来构造镜像。[source](https://github.com/istio/istio/tree/master/samples/bookinfo/src/productpage)

### 镜像选择

目前计划使用ubuntu来作为基础镜像。在部署并测试的时候发现了自己犯了一个愚蠢的错误——我直接使用了官方的ubuntu镜像而没有构造自己的镜像，这导致镜像在部署上k8s后就一直在重复启动，因为它根本就没有其他的命令……

对于docker而言，它并不是一个虚拟机，而更应该被看作是一个被隔离的镜像。其生命周期完全取决于CMD命令，一旦CMD命令被执行完毕，则整个容器就会进入结束状态。

对于`ubuntu:18.04`这样的镜像而言，其命令为：

```
"Cmd": [
    "/bin/sh",
    "-c",
    "#(nop) ",
    "CMD [\"/bin/bash\"]"
],
或者
"Cmd": ["/bin/bash"]
```

上述信息可以使用`docker inspect ubuntu:18.04 | grep Cmd -C 5`得到。事实上，在线下可以直接使用`docker run -it ubuntu:18.04 /bin/bash`进入，但是一旦退出容器也将终止，因为`/bin/bash`命令已经执行完毕了。

如何让其长期运行？这个问题在[SO](https://stackoverflow.com/questions/43419500/how-do-you-start-a-docker-ubuntu-container-into-bash)上也被问过挺多次的。

我使用ubuntu镜像的主要原因是为了验证Dockerfile中所写的命令执行起来会不会有问题，因此也不用特意去细究Dockerfile该怎么写。这个问题中提到的解决方案很好，就是在正常的ubuntu镜像下使用`docker run -d  stens_ubuntu  sleep infinity`，或者将`sleep infinity`作为CMD的参数，让容器永久停留。这样就可以直接登陆容器内部测试参数/命令，同时退出之后也不用担心容器会直接停止了。

### 构造镜像

测试镜像的Dockerfile：

```Dockerfile
FROM ubuntu:18.04
MAINTAINER wtysos11 "wtysos11@gmail.com"

COPY sources.list ./
COPY pip.conf ./
RUN rm /etc/apt/sources.list\
  && mv sources.list /etc/apt/sources.list\
  && mkdir -p ~/.pip \
  && mv pip.conf ~/.pip/pip.conf \
  && apt-get update \
  && apt-get install -y python3-pip python3-dev \
  && cd /usr/local/bin \
  && ln -s /usr/bin/python3 python \
  && ln -s /usr/bin/pip3 pip \
  && python -m pip install --upgrade setuptools pip wheel 

CMD ["/bin/bash"]
```

其中sources.list使用的是[清华镜像站](https://mirrors.tuna.tsinghua.edu.cn/help/ubuntu/)18.04版本，注意需要使用http，不然会报证书错误。pip镜像采用的是清华镜像。

执行`docker build . -t ubuntu-test:0.0.5`构造镜像，然后执行`docker run --name test -d  ubuntu-test:0.0.5  sleep infinity`。再使用`docker exec -it test bash`

经过测试后，原来的命令执行是没有问题的，并且可以执行`top`命令。

### 进入容器开始测试

pod名称为`productpage-test-66b95cdbd8-phrw8`，istio注入pod。

使用`kubectl exec -it productpage-test-66b95cdbd8-phrw8 --container container-0 -n=wtytest /bin/bash`即可进入在线容器。

开始的时候比较菜，只知道使用`top`命令，然后发现top刷新起来不好记录。随后发现python有个很好用的库`psutil`，比如可以使用下面的代码实现类似top的效果：

```python
for x in range(10):
    print(psutil.cpu_percent(interval=1, percpu=True))
```

### 如何计算CPU占用率

这里我们遇到了两个问题，一个是如何在k8s的平台下计算pod的CPU占用率；另外一个是是否要将istio注入到每个pod中的container: istio-proxy的CPU占用率纳入考量。

#### kubernetes中的CPU测算（基于prometheus）

这个一直以来我都没有太弄明白。根据[issue](https://github.com/google/cadvisor/issues/2026)中某位大佬的说法，如果使用kubernetes来运行容器的话，`container_spec_cpu_shares`为容器的请求量，而`container_spec_cpu_quota`为容器的限制量。从这个角度来说，应该是采用`container_cpu_usage_seconds_total/container_spec_cpu_quota * constant`的方式来计算（虽然我有些没有弄明白，不是应该计算当前申请的CPU中有多少被消耗吗）。

之后这位大佬又用另外一个想法，即`rate(container_cpu_usage_seconds[10m])/(container_spec_cpu_quota / container_spec_cpu_period)`类似这种式子来进行计算，赞同的人很多。

同时有人展示了自己所用的grafana绘图公式`sum(rate(container_cpu_usage_seconds_total{name!~".*prometheus.*", image!="", container_name!="POD"}[5m])) by (pod_name, container_name) /sum(container_spec_cpu_quota{name!~".*prometheus.*", image!="", container_name!="POD"}/container_spec_cpu_period{name!~".*prometheus.*", image!="", container_name!="POD"}) by (pod_name, container_name)`

或者我可以使用另外一种想法，即不统计CPU占用率（因为这个CPU占用率变动会很大），而选择直接预测实际使用的CPU时间，或者说，等效vCPU的数量。但是这样获得的前提是能够确定拿到的数据是准确的，因此回过头来还是要验证prometheus拿到的数据与本地拿到的数据在一定时间内至少相差不能太大。

[prometheus中文文档](https://yunlzheng.gitbook.io/prometheus-book/part-ii-prometheus-jin-jie/exporter/commonly-eporter-usage/use-prometheus-monitor-container)里列了cAdvisor搜集的CPU相关指标，包括：
* container_cpu_load_average_10s，gauge型，过去10s容器CPU的平均负载（负载代表队伍中等待的进程数量，详见[High CPU utilization but low load average](https://serverfault.com/questions/667078/high-cpu-utilization-but-low-load-average)
* container_cpu_usage_seconds_total，counter，容器CPU累积占用时间
* container_cpu_system_seconds_total，System CPU累积占用时间
* container_cpu_user_seconds_total，User CPU累积占用时间

此外[docs of cAdvisor](https://docs.signalfx.com/en/latest/integrations/agent/monitors/cadvisor.html)更详尽地介绍了cAdvisor能够抓到的指标，包括：
* container_spec_cpu_period ,Gauge型，The number of microseconds that the CFS scheduler uses as a window when limiting container processes。
* container_spec_cpu_quota，Gauge型，In CPU quota for the CFS process scheduler. In K8s this is equal to the containers’s CPU limit as a fraction of 1 core and multiplied by the container_spec_cpu_period. So if the CPU limit is 500m (500 millicores) for a container and the container_spec_cpu_period is set to 100,000, this value will be 50,000. 这个值在kubernetes中应该与容器的CPU限制量是相关的（PS：CPU限制中的m表示千分之一个CPU）
* container_spec_cpu_shares，Gauge型，CPU share of the container，应该是实际分配给容器的CPU数量。


其中[CFS Scheduler](https://www.kernel.org/doc/Documentation/scheduler/sched-bwc.txt)应该是关于CPU限额的调度。The bandwidth allowed for a group is specified using a quota and period. Within
each given "period" (microseconds), a group is allowed to consume only up to
"quota" microseconds of CPU time. 这也是container_spec_cpu_quota这个值为什么是这么计算的原因。