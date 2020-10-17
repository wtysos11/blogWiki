# 在k8s中上线gatling镜像并在内网发送流量

很多时候我们会面临一个问题，即外网的带宽是有限的，虽然未来有扩容的可能，但是短时间内也不能直接扩容，而测试本身是无限的。因此，如果不能够在内网下直接发包进行测试，那由于带宽限制打不到较大的压力，对于一些容器的测试很可能就达不到效果。

因此我们需要在内网有一个能够配置的压力测试容器，目前选定了gatling，因为其功能比较强大，而且安装很方便。

## 镜像构造

### 初始镜像构造

虽然之前已经做了`ubuntu`的镜像，并且可以使用`apt-get install`来安装gatling，但是这种方式安装后有些不太会用，似乎更多是作为一个插件存在而不是独立存在的软件。

我还是选择了自己最熟悉的方式，直接从官网上下载了开源版本的standalone gatling.zip，解压后将目录重命名为gatling，Dockerfile如下:

```Dockerfile
FROM ubuntu:18.04
MAINTAINER wtysos11 "wtysos11@gmail.com"

COPY sources.list ./
ADD gatling ./gatling

RUN rm /etc/apt/sources.list\
  && mv sources.list /etc/apt/sources.list\
  && apt-get update \ 
  && apt-get install -y openjdk-8-jdk scala 
  #&& apt-get install -y gatling

CMD ["/bin/bash"]
```

sources.list为清华的apt镜像，为了加速；gatling可以在java8下运行，必须要安装scala（其实我个人觉得只安装scala就够了，保险起见）

操作完之后执行`docker build . -t ubuntu-wtynettest:0.0.2`构造镜像，然后执行`docker run --name test -d  ubuntu-wtynettest:0.0.2  sleep infinity`。再使用`docker exec -it test bash`

经过测试，gatling软件能够正常运行并且访问外界指定端口。如果我没有记错，k8s中的pod暴露端口主要是为了转发流量，那容器自己往外发流量应该是不用暴露端口的，因此直接上k8s是没有问题的。

### 进阶流量压力测试镜像构造

下面的任务为：
1. 在镜像文件中配置环境变量，该变量最好是能够在`docker build`的时候修改而不是要手动改写Docerfile，这样后续写bash脚本之类的会比较方便。（如果能够实时传入就更好了，不过这要将gatling作为插件实现，改写太多了，计划放在第三步）
2. 删除原有系统中的脚本文件，并上传指定的脚本文件`test.scala`。这个脚本文件要能够读取环境变量来替换指定的值。
3. 命令直接设为指定的发包命令。

#### 环境变量配置

我看了一下，使用`docker build`加参数的方式似乎并不常见，而且其他方式也挺麻烦的。

因此我直接使用了`ENV Key=value`的形式（如果value中间有空格，两边要加上双引号）

#### 脚本读取环境变量

scala脚本是可以读取到环境变量的，方法挺多的。目前选择的是直接使用`sys.env["EnvVar"]`，此时需要环境中能够读取到`$EnvVar`，不需要引入任何库。

这个方法的缺点是如果环境中没有设置环境变量会报错，不过这也不是什么大问题，毕竟在docker内部。

接下来就很简单了，将这个值作为方法的参数进行传递，然后把脚本送到指定的位置。

#### 命令配置

由于standalone版本的gatling是使用`gatling.sh`进行执行的，因此我预先写了一个输入文件进行重定向（其实就是一个只有`1`+回车的文件）。如此，容器的命令配置完毕。

最终的Dockerfile：
```Dockerfile
FROM ubuntu:18.04
MAINTAINER wtysos11 "wtysos11@gmail.com"

ENV Test="http://192.168.0.173:19001"

COPY sources.list ./
COPY nettest.scala ./
COPY command.txt ./
ADD gatling ./gatling

RUN rm /etc/apt/sources.list\
  && mv sources.list /etc/apt/sources.list\
  && rm /gatling/user-files/simulations/computerdatabase/ -R \
  && mv nettest.scala /gatling/user-files/simulations/nettest.scala \
  && apt-get update \ 
  && apt-get install -y openjdk-8-jdk 
  #&& apt-get install gatling

#CMD ["sh","-c","/gatling/bin/gatling.sh", "<","/test/command.txt"]
CMD ["sleep infinity"]
```

我失败了好多次……这里的命令行有一个exec form和shell form的[区别](https://stackoverflow.com/questions/42805750/dockerfile-cmd-shell-versus-exec-form)。我之前使用`CMD ["/gatling/bin/gatling.sh < /test/command.txt"]`的时候系统经常提示我找不到后面那个输入文件，真是让人摸不着头脑。然后换用exec form之后gatling根本输入

其中移除gatling内系统自带脚本的目的是为了让用户脚本一定排在第一位。由于版本不同，系统自带脚本可能有所区别，需要注意。

下面的文件：
* command.txt，内含`1`+空格，表示输入给`gatling.sh`的内容
* nettest.scala，一个可以读取`$Test`作为目标地址的gatling脚本
* gatling，解压官方包gatling.zip后的文件夹

执行测试部分命令：
* `docker build . -t ubuntu-wtynettest:0.0.3`
* `docker run --name test -d ubuntu-wtynettest:0.0.3 sleep infinity`
* `docker exec -it test bash`

### 动态挂载

上面的实现方案还是有一个问题，即没有办法灵活控制gatling，只能够每次生成一个实例在挂载到k8s上，非常麻烦。

我在思考有没有一种方式，能够将一个gatling程序传到k8s集群中，只需要通过网络端口向其上传配置文件、发送命令就可以调用指定的压力测试脚本。

我的实现思路需要用scala做一个简易的服务器，而网上的思路似乎有些不太一样。

* [gatling docker image](https://github.com/denvazh/gatling)是github上一个gatling的docker镜像，通过挂载配置文件能够在本地的docker上进行压力测试。
* [Distributed load testing with Gatling and Kubernetes](https://movile.blog/distributed-load-testing-with-gatling-and-kubernetes/)这个是用gatling docker kubernetes关键词搜索出来的文章，似乎和我的思路比较类似。

不过有没有必要作出这个项目也是一个问题，毕竟gatling中仍然存在一些问题没有弄清楚，比如atOnceUser和constantUser等测试方式之间的选择等。