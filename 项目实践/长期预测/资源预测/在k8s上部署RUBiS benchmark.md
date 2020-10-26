# 在k8s上部署RUBiS benchmark(php version)

因为时间有限，所以我选择不做镜像，直接将容器当作虚拟机来进行后续的操作。

RUBiS 2002年的实现实际上是java servelet、php和EJB三种方式实现的，因此我只用选择一种就好了。这里因为找到了别人配置的过程，所以选择了php（我记得有一篇论文也选用的php版本，在调研中发现与Apache关联会更紧密，配置会容易些）

参考文献：
* [RUBiS php配置](https://www.cnblogs.com/damn-chris/archive/2012/03/06/2382146.html)、[RUBiS的安装与配置](https://www.dazhuanlan.com/2020/02/13/5e44e2d334638/)
* [sguazt/RUBiS](https://github.com/sguazt/RUBiS)/[uillianluiz/RUBiS](https://github.com/uillianluiz/RUBiS)。这两个是目前能够找到的github上对RUBiS的修改后结果。我个人选择的是第一个，因为第一个版本比较老，而且没有做太多的改动。说实话如果改动太多我还不如用其他语言自己再写一个，毕竟servelet和php我一个都不会。

## 安装步骤

下面简要介绍安装步骤。虽然计划使用prometheus，但是保险起见还是安装上之前的监控软件。官方的Dockerhub有`CentOS5`的镜像，但是国内很多镜像站都不支持了，预计安装起来会很困难。经过分析我觉得很多软件的版本应该不重要（最关键的php版本已经被github大佬修改过了），因此我冒险选用了`Centos7`。使用`docker pull centos:7`拉去镜像。

计划：
1. 安装数据库，暴露3306端口
2. 安装php版本的RUBiS
3. 使用client来初始化，要修改配置文件中的`database_server`来指定mysql数据库所在的IP，而且必须要为3306端口。本项与下一项都会在`edu/rice/rubis/client/RUBiSProperties.java`中被导入到client的设置中。
4. 访问要修改配置文件中的`httpd_hostname`为网络服务器所在的节点IP来进行。由于client需要远程登陆所有的服务器并访问sysstat服务，相当困难，因此原始的client脚本应该是用不了了，大部分都是java硬编码，很难通过istio来进行访问。这个可以之后再做进一步尝试。


### 数据库

数据库使用的是任意版本的mysql，基本上包含以下几个步骤：
1. 创建用户`rubis`和`rubis@localhost`,例`CREATE USER 'rubis'@'%' IDENTIFIED BY 'password';`
2. 给1中创建的两个用户赋予权限，例`GRANT ALL PRIVILEGES ON *.* TO 'rubis'@'%';`
3. 检查3306端口是否访问（在容器上需要提前暴露出3306端口）
4. 调用sql文件初始化rubis所需要的数据库
   ```bash
    cd /home/RUBiS/database
    mysql -urubis -ppassword < rubis.sql
    mysql -urubis -ppassword rubis < regions.sql
    mysql -urubis -ppassword rubis < categories.sql
   ```

需要构建一个mysql镜像，目前先考虑手动输入，后续再考虑做一个自动脚本。比较麻烦的是这个Dockerfile需要挂载到k8s上，这和挂载到Docker上又是有些不一样的，因此测试的时候用外部的mysql仓库代替，实际使用是直接用k8s上的mysql。

#### k8s-mysql

参考：
* [官方](https://kubernetes.io/zh/docs/tasks/run-application/run-single-instance-stateful-application/)

配置文件（出于方便直接在yaml内配置账号和密码）

service and deployment: 

```yaml
apiVersion: v1
kind: Service
metadata:
  name: mysql
spec:
  type: NodePort
  selector:
    app: mysql
  ports:
  - port: 3306
    targetPort: 3306
    nodePort: 31107
---
apiVersion: apps/v1 # for versions before 1.9.0 use apps/v1beta2
kind: Deployment
metadata:
  name: mysql
spec:
  selector:
    matchLabels:
      app: mysql
  strategy:
    type: Recreate
  template:
    metadata:
      labels:
        app: mysql
    spec:
      containers:
      - image: mysql:5.7
        name: mysql
        env:
          # Use secret in real usage
        - name: MYSQL_ROOT_PASSWORD
          value: password
        ports:
        - containerPort: 3306
          name: mysql
        volumeMounts:
        - name: mysql-persistent-storage
          mountPath: /var/lib/mysql
      volumes:
      - name: mysql-persistent-storage
        persistentVolumeClaim:
          claimName: mysql-pv-claim

```

pv与pvc

```yaml
apiVersion: v1
kind: PersistentVolume
metadata:
  name: mysql-pv-volume
  labels:
    type: local
spec:
  storageClassName: manual
  capacity:
    storage: 10Gi
  accessModes:
    - ReadWriteOnce
  hostPath:
    path: "/mnt/data"
---
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: mysql-pv-claim
spec:
  storageClassName: manual
  accessModes:
    - ReadWriteOnce
  resources:
    requests:
      storage: 10Gi
```

然后先将pv安装上去，再将deployment和svc放上去即可。

访问的话因为本地没有mysql客户端，我使用了mysql镜像来进行。原文为：
```
This image can also be used as a client for non-Docker or remote instances:

$ docker run -it --rm mysql mysql -hsome.mysql.host -usome-mysql-user -p

k8s:
kubectl run -it --rm --image=mysql:5.6 --restart=Never mysql-client -- mysql -h mysql -ppassword
```

#### RUBiS-client

需要有一个外部服务器能够执行sql脚本，具体来说，则是运行命令：`docker run -it --rm mysql:5.7 mysql -h139.9.74.70 -uroot -P31107 -p`来访问mysql，不过还需要带上那些.sql文件。

因此，我们需要有一个client能够挂载到服务器上，并且携带指定的脚本。构建需要官方包`/RUBiS/database`下的所有文件。原本想要使用官方镜像，但是发现很麻烦，还不如直接用回`ubuntu:18.04`

```Dockerfile
FROM ubuntu:18.04
MAINTAINER wtysos11 "wtysos11@gmail.com"

RUN mkdir test
COPY sources.list ./test/
ADD database ./database
RUN rm /etc/apt/sources.list\
  && mv /test/sources.list /etc/apt/sources.list\
  && apt-get update \ 
  && apt-get install -y mysql-client vim 
  #&& apt-get install gatling

CMD ["sh","-c","sleep infinity"]
```

编译`docker build . -t wty-mysqltest:0.0.1`,运行镜像`docker run --name test -d wty-mysqltest:0.0.1`，然后使用bash命令进入`docker exec -it test bash`。

在进入镜像内部后，进入到`database`文件夹中，执行下列命令

```bash
 cd /home/RUBiS/database
 mysql -h139.9.74.70 -urubis -ppassword -P31107 < rubis.sql
 mysql -h139.9.74.70 -urubis -ppassword -P31107 rubis < regions.sql
 mysql -h139.9.74.70 -urubis -ppassword -P31107 rubis < categories.sql
```

如此，数据库创建完毕，所有的表都已经被导入。

#### Dockerfile及其他文件

考虑到国内很多软件下载会很慢（而且兼容性应该没有问题），因此使用`yum`下载依赖库，并用[中科大镜像](http://mirrors.ustc.edu.cn/help/centos.html)来加速。

```Dockerfile
FROM centos:7
MAINTAINER wtysos11 "wtysos11@gmail.com"

RUN sudo sed -e 's|^mirrorlist=|#mirrorlist=|g' \
         -e 's|^#baseurl=http://mirror.centos.org/centos|baseurl=https://mirrors.ustc.edu.cn/centos|g' \
         -i.bak \
         /etc/yum.repos.d/CentOS-Base.repo \
    && sudo yum makecache \ 
    && sudo yum install -y ...

CMD ["sleep infinity"]
```

### php配置

参考：
* [apache & php on docker](https://writing.pupius.co.uk/apache-and-php-on-docker-44faef716150)，该文章介绍了如何在docker上部署Apache+php。
* [centos-apahce-php](https://hub.docker.com/r/naqoda/centos-apache-php/)，支持多个版本的apache和php，可以参考一下Dockerfile。
* [官方php:version-apache](https://hub.docker.com/_/php/)，官方提供了相关的镜像，我之前居然没有发现。这个镜像可以自行添加Configuration File来实现一定的修改。

我个人还是倾向于直接使用别人的镜像，而不是自己来手动做。毕竟构建镜像是一项比较麻烦的事情，而且我对于运行环境的要求也不高，有一定的冗余也是可以接受的。

选用php版本很重要的原因就是php脚本不用编译。RUBiS使用Apache来提供服务，我们使用`$RUBIS_HOME`来指代当前版本RUBiS所在的目录。

对于Apache来说，一般文档的根目录是在`/var/www`下，因此我们只需要把`$RUBIS_HOME/php`的内容拷贝到`/var/www`即可

```bash
mkdir -p /var/www/rubis
cp -r $RUBIS_HOME/php /var/www/rubis/PHP
```

这一步唯一的难点在于Apache的docker镜像该如何使用，以及如何让php脚本访问到指定的数据库（可能需要改写原脚本中的地址）

从几篇文档的信息来看，主要是修改`RUBiS/php/PHPPrinter.php`文件中的`getDatabaseLink`函数，需要将`mysql_pconnect()`里面的三个参数进行修改。

#### 最终结果

经过了多次尝试，最终还是选择了官方镜像，Dockerfile如下：

```Dockerfile
FROM php:7.2-apache
ENV HOSTNAME="10.247.188.145"
ENV SQLUSER="rubis"
ENV SQLPASSWD="password"
ENV SQLPORT="31107"

COPY --from=mlocati/php-extension-installer /usr/bin/install-php-extensions /usr/bin/

RUN install-php-extensions mysqli

RUN sed -i 's#http://deb.debian.org#https://mirrors.163.com#g' /etc/apt/sources.list\
    && sed -i 's#http://security.debian.org#https://mirrors.163.com#g' /etc/apt/sources.list\
    && mv "$PHP_INI_DIR/php.ini-production" "$PHP_INI_DIR/php.ini"\
    && sed -i '901s/;//' /usr/local/etc/php/php.ini \
    && mkdir -p /var/www/html/PHP\
    && apt-get update\
    && apt-get install -y vim

EXPOSE 80
COPY php/ /var/www/html/PHP/
```

其中`php`为原来的php脚本，这里选择了[github](https://github.com/sguazt/RUBiS)代码的php文件夹，就可以不用费事修改了（还是要修改的，因为php7不能用mysql插件而只能用mysqli插件，所以要替换所有文件中所有对应的方法）。如果是使用原来RUBiS 1.4.3版本的代码，需要按照参考文章中所说，修改两个过期的方法来使得新文件适应php7。

最后还需要再加上一个环境变量，方便在php中配置数据库连接。这里使用的方法为在Dockerfile中加上`ENV Test="http://192.168.0.173:19001"`这样的字段（因为在启动后再加上后php读不到），然后使用php的`getenv`方法，如

```php
<?php
$ip = getenv('Test');
echo $ip;
?>
```

就可以获得结果了。

很重要的一点是该镜像中并没有开启`mysql`的extension，需要访问`/usr/local/etc/php/php.ini`，去掉第901行`extension=mysqli`前面的分号。

对`RUBiS/PHP/PHPprinter.php`中的前面进行如下修改，其中`mysql_pconnect`在PHP7中被废除了，使用的是另外一个`$con=mysqli_connect("localhost","my_user","my_password","my_db");`

```php
function getDatabaseLink(&$link)
{
  $hostname = getenv("HOSTNAME");
  $user = getenv("SQLUSER");
  $passwd = getenv("SQLPASSWD");
  $port = getenv("SQLPORT");
  $link = mysqli_connect($hostname,$user,$passwd,'rubis',$port);
//...

```

#### 一些坑

1. 在进入apahce-php镜像中修改php代码的时候，新建的php代码必须要`chmod 777`，不然apache服务器将没有访问权限（因为在镜像中权限一般都是root级别）
2. php7官方的包中默认是不带`mysqli`的，这就很难受了。直接改写`/usr/local/etc/php/php.ini`中第901行的分号是没有用的，因为它就没带这个。需要按照[插件安装](https://github.com/mlocati/docker-php-extension-installer)里面的步骤修改Dockerfile来把它装上。使用`php -m`可以列出所有能用的插件。
3. 一开始的时候mysql在k8s中的访问方式是NodePort，但是在php服务器中它居然无法访问到外部的这个IP（挺奇怪的，但是可以访问到istio的ingressgateway），可能是DNS配置问题。我后来改用ClusterIP之后就可以了。

### 客户端

客户端应该还是要配置的，有两个作用，一个是初始化服务器，另一个是模拟应用请求对rubis进行测试。

客户端是使用java来实现的，而且依赖的版本很老（`jdk 1.3.1`），但是跑了一下代码，虽然有些特性deprecated了，但是应该还是能跑起来的。

步骤：
1. 修改Client文件夹中的`rubis.properties`文件，主要修改以下几个条目：apache服务器IP/域名`httpd_hostname`、Apache服务器端口`httpd_port`、`httpd_use_version`（使用的服务器版本，目前为PHP）