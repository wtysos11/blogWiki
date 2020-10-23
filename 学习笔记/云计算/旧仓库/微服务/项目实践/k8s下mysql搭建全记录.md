# k8s下mysql搭建全记录

标签：k8s docker mysql

## mysql的Dockerfile

上学期的时候我写过一个Dockerfile来作为镜像，其实效果还是不错的：

```Dockerfile
FROM guyton/centos6
RUN yum install -y mysql mysql-server
RUN /etc/init.d/mysqld start &&\
mysql -e "grant all privileges on *.* to 'root'@'%' identified by '123456';" &&\
mysql -e "grant all privileges on *.* to 'root'@'localhost' identified by '123456';"
EXPOSE 3306
CMD ["mysqld_safe"]
```

### 菜鸟教程的镜像

```Dockerfile
FROM debian:jessie

# add our user and group first to make sure their IDs get assigned consistently, regardless of whatever dependencies get added
RUN groupadd -r mysql && useradd -r -g mysql mysql

# add gosu for easy step-down from root
ENV GOSU_VERSION 1.7
RUN set -x \
    && apt-get update && apt-get install -y --no-install-recommends ca-certificates wget && rm -rf /var/lib/apt/lists/* \
    && wget -O /usr/local/bin/gosu "https://github.com/tianon/gosu/releases/download/$GOSU_VERSION/gosu-$(dpkg --print-architecture)" \
    && wget -O /usr/local/bin/gosu.asc "https://github.com/tianon/gosu/releases/download/$GOSU_VERSION/gosu-$(dpkg --print-architecture).asc" \
    && export GNUPGHOME="$(mktemp -d)" \
    && gpg --keyserver ha.pool.sks-keyservers.net --recv-keys B42F6819007F00F88E364FD4036A9C25BF357DD4 \
    && gpg --batch --verify /usr/local/bin/gosu.asc /usr/local/bin/gosu \
    && rm -r "$GNUPGHOME" /usr/local/bin/gosu.asc \
    && chmod +x /usr/local/bin/gosu \
    && gosu nobody true \
    && apt-get purge -y --auto-remove ca-certificates wget

RUN mkdir /docker-entrypoint-initdb.d

# FATAL ERROR: please install the following Perl modules before executing /usr/local/mysql/scripts/mysql_install_db:
# File::Basename
# File::Copy
# Sys::Hostname
# Data::Dumper
RUN apt-get update && apt-get install -y perl pwgen --no-install-recommends && rm -rf /var/lib/apt/lists/*

# gpg: key 5072E1F5: public key "MySQL Release Engineering <mysql-build@oss.oracle.com>" imported
RUN apt-key adv --keyserver ha.pool.sks-keyservers.net --recv-keys A4A9406876FCBD3C456770C88C718D3B5072E1F5

ENV MYSQL_MAJOR 5.6
ENV MYSQL_VERSION 5.6.31-1debian8

RUN echo "deb http://repo.mysql.com/apt/debian/ jessie mysql-${MYSQL_MAJOR}" > /etc/apt/sources.list.d/mysql.list

# the "/var/lib/mysql" stuff here is because the mysql-server postinst doesn't have an explicit way to disable the mysql_install_db codepath besides having a database already "configured" (ie, stuff in /var/lib/mysql/mysql)
# also, we set debconf keys to make APT a little quieter
RUN { \
        echo mysql-community-server mysql-community-server/data-dir select ''; \
        echo mysql-community-server mysql-community-server/root-pass password ''; \
        echo mysql-community-server mysql-community-server/re-root-pass password ''; \
        echo mysql-community-server mysql-community-server/remove-test-db select false; \
    } | debconf-set-selections \
    && apt-get update && apt-get install -y mysql-server="${MYSQL_VERSION}" && rm -rf /var/lib/apt/lists/* \
    && rm -rf /var/lib/mysql && mkdir -p /var/lib/mysql /var/run/mysqld \
    && chown -R mysql:mysql /var/lib/mysql /var/run/mysqld \
# ensure that /var/run/mysqld (used for socket and lock files) is writable regardless of the UID our mysqld instance ends up having at runtime
    && chmod 777 /var/run/mysqld

# comment out a few problematic configuration values
# don't reverse lookup hostnames, they are usually another container
RUN sed -Ei 's/^(bind-address|log)/#&/' /etc/mysql/my.cnf \
    && echo 'skip-host-cache\nskip-name-resolve' | awk '{ print } $1 == "[mysqld]" && c == 0 { c = 1; system("cat") }' /etc/mysql/my.cnf > /tmp/my.cnf \
    && mv /tmp/my.cnf /etc/mysql/my.cnf

VOLUME /var/lib/mysql

COPY docker-entrypoint.sh /usr/local/bin/
RUN ln -s usr/local/bin/docker-entrypoint.sh /entrypoint.sh # backwards compat
ENTRYPOINT ["docker-entrypoint.sh"]

EXPOSE 3306
CMD ["mysqld"]
```

然后使用`docker build -t mysql .`来创建一个镜像。

这里的一个小问题是这个不是mysql的官方镜像，原因是当时我不知道官方镜像怎么用。一个大问题是它实际上是现场下载mysql，其实还是有很大的延迟的。因此我想要探索一下如果使用mysql作为Dockerfile的基础镜像该怎么搭建。

### mysql的镜像

首先使用`docker pull mysql`拉取镜像。

然后使用`docker run -p 3306:3306 --name mymysql -v $PWD/conf:/etc/mysql/conf.d -v $PWD/logs:/logs -v $PWD/data:/var/lib/mysql -e MYSQL_ROOT_PASSWORD=123456 -d mysql:5.6`来运行

可以使用mysql的控制台远程登陆指定的3306端口进行访问。

来自菜鸟教程的例子：

```bash
# docker 中下载 mysql
docker pull mysql

#启动
docker run --name mysql -p 3306:3306 -e MYSQL_ROOT_PASSWORD=Lzslov123! -d mysql

#进入容器
docker exec -it mysql bash

#登录mysql
mysql -u root -p
ALTER USER 'root'@'localhost' IDENTIFIED BY 'Lzslov123!';

#添加远程登录用户
CREATE USER 'liaozesong'@'%' IDENTIFIED WITH mysql_native_password BY 'Lzslov123!';
GRANT ALL PRIVILEGES ON *.* TO 'liaozesong'@'%';
```

## k8s下的StatefulSet如何写

参考资料：

* [k8s容器编排之Stateful持久化部署MySQL](https://www.jianshu.com/p/1fcaf702d104)：详细介绍了Mysql的Deployment和PV的写法，以及Secret，但是不是StatefulSet，还是有些不一样。问题在于挂载上后提示endpoints not found，似乎不能够用。
* [k8s官方中文翻译-基于Persistent Volumes搭建WordPress和MySQL应用](https://kubernetes.io/zh/docs/tutorials/stateful-application/mysql-wordpress-persistent-volume/)
* [Example: Deploying WordPress and MySQL with Persistent Volumes](https://kubernetes.io/docs/tutorials/stateful-application/mysql-wordpress-persistent-volume/)，就是上面那篇的英文版，挺不错的。
* [k8s中文社区-k8s部署高可用的mysql](https://www.kubernetes.org.cn/3985.html)，没有提到persistent volume，我当时也不知道，看的有点一知半解的。
* [官方-run a replicated stateful application](https://kubernetes.io/docs/tasks/run-application/run-replicated-stateful-application/)
* [stackoverflow-kubernetes mysql statefulset with root password](https://stackoverflow.com/questions/51905060/kubernetes-mysql-statefulset-with-root-password)


比较好的学习资料：

* [stackoverflow-kubernetes mountPath vs hostPath](https://stackoverflow.com/questions/51107390/kubernetes-mountpath-vs-hostpath)，volumes.hostPath和volumeMounts.mountPath想必很多人第一眼看过去都会有所疑问。注意，hostPath为主机上用作映射的文件或目录，（应该要事先创建，没有验证）
* [官方文档-如何在Pod上使用pv](https://kubernetes.io/docs/tasks/configure-pod-container/configure-persistent-volume-storage/)，首先，管理员创建PV，并且没有将这个PV与任何的Pod相联系。用户创建PVC，联系上这个PV，Pod使用这个PVC作为存储。

### 前置条件


提供容量大于10g的，读取类型为`ReadWriteOnce`的持久化储存卷。

PV的示例：

```yaml
apiVersion: v1
kind: PersistentVolume
metadata:
  name: pv0003
spec:
  capacity:
    storage: 10Gi
  volumeMode: Filesystem
  accessModes:
    - ReadWriteOnce
  persistentVolumeReclaimPolicy: Retain
  mountOptions:
    - hard
    - nfsvers=4.1
  nfs:
    path: /tmp
    server: 172.17.0.2
```

我所用的PV：

```yaml
kind: PersistentVolume
apiVersion: v1
metadata:
  name: mysql-pv-volume
  labels:
    type: local
spec:
  capacity:
    storage: 10Gi
  accessModes:
    - ReadWriteOnce
  hostPath:
    path: "/mnt/data"

```

注：

1. 我事先在节点上创建好了`/mnt/data`的文件夹，不知道有没有影响
2. 最好不要加上**storageName**，按原理来讲，PVC没有加上storageName是不会去请求有storageName的PV的。但就算加上了storageName，请求居然也失败了，不清楚具体的原因。没有加上的时候就对了。

### 开始布置StatefulSet

主要参考的是[这篇](https://stackoverflow.com/questions/51905060/kubernetes-mysql-statefulset-with-root-password)

国内需要提前下载`gcr.io/google-samples/xtrabackup:1.0`镜像，我是使用的中科大的镜像`gcr.mirrors.ustc.edu.cn/google-samples/xtrabackup:1.0`，不过速度相当的慢。还有`gcr.azk8s.cn/namespace/image_name:image_tag `的镜像，这个速度快很多。之后使用`docker tag hash newname`修改名字就好了，比如`docker tag c415dbd7af07 gcr.io/google-samples/xtrabackup:1.0`。此外还可以参考[这篇文章](https://blog.csdn.net/weixin_39961559/article/details/80739352)，用阿里的镜像站来做代理。

执行`kubectl create secret generic mysql-pass --from-literal=MYSQL_ROOT_PASSWORD=wu97112500`创建mysql密钥

mysql Configmap.yaml

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: mysql
  labels:
    app: mysql
data:
  master.cnf: |
    # Apply this config only on the master.
    [mysqld]
    log-bin
  slave.cnf: |
    # Apply this config only on slaves.
    [mysqld]
    super-read-only
```

mysql Servcie.yaml

```yaml
# Headless service for stable DNS entries of StatefulSet members.
apiVersion: v1
kind: Service
metadata:
  name: mysql
  labels:
    app: mysql
spec:
  ports:
  - name: mysql
    port: 3306
  clusterIP: None
  selector:
    app: mysql
---
# Client service for connecting to any MySQL instance for reads.
# For writes, you must instead connect to the master: mysql-0.mysql.
apiVersion: v1
kind: Service
metadata:
  name: mysql-read
  labels:
    app: mysql
spec:
  ports:
  - name: mysql
    port: 3306
  selector:
    app: mysql
```

mysql.yaml

```yaml
apiVersion: apps/v1beta1
kind: StatefulSet
metadata:
  name: mysql
spec:
  selector:
    matchLabels:
      app: mysql
  serviceName: mysql
  replicas: 1
  template:
    metadata:
      labels:
        app: mysql
    spec:
      initContainers:
      - name: init-mysql
        image: mysql:5.7
        command:
        - bash
        - "-c"
        - |
          set -ex
          # Generate mysql server-id from pod ordinal index.
          [[ `hostname` =~ -([0-9]+)$ ]] || exit 1
          ordinal=${BASH_REMATCH[1]}
          echo [mysqld] > /mnt/conf.d/server-id.cnf
          # Add an offset to avoid reserved server-id=0 value.
          echo server-id=$((100 + $ordinal)) >> /mnt/conf.d/server-id.cnf
          # Copy appropriate conf.d files from config-map to emptyDir.
          if [[ $ordinal -eq 0 ]]; then
            cp /mnt/config-map/master.cnf /mnt/conf.d/
          else
            cp /mnt/config-map/slave.cnf /mnt/conf.d/
          fi
        volumeMounts:
        - name: conf
          mountPath: /mnt/conf.d
        - name: config-map
          mountPath: /mnt/config-map
      - name: clone-mysql
        image: gcr.io/google-samples/xtrabackup:1.0
        command:
        - bash
        - "-c"
        - |
          set -ex
          # Skip the clone if data already exists.
          [[ -d /var/lib/mysql/mysql ]] && exit 0
          # Skip the clone on master (ordinal index 0).
          [[ `hostname` =~ -([0-9]+)$ ]] || exit 1
          ordinal=${BASH_REMATCH[1]}
          [[ $ordinal -eq 0 ]] && exit 0
          # Clone data from previous peer.
          ncat --recv-only mysql-$(($ordinal-1)).mysql 3307 | xbstream -x -C /var/lib/mysql
          # Prepare the backup.
          xtrabackup --prepare --target-dir=/var/lib/mysql
        volumeMounts:
        - name: data
          mountPath: /var/lib/mysql
          subPath: mysql
        - name: conf
          mountPath: /etc/mysql/conf.d
      containers:
      - name: mysql
        image: mysql:5.7
        env:
        - name: MYSQL_ROOT_PASSWORD
          valueFrom:
            secretKeyRef:
              name: mysql-pass
              key: MYSQL_ROOT_PASSWORD
        ports:
        - name: mysql
          containerPort: 3306
        volumeMounts:
        - name: data
          mountPath: /var/lib/mysql
          subPath: mysql
        - name: conf
          mountPath: /etc/mysql/conf.d
        resources:
          requests:
            cpu: 500m
            memory: 1Gi
        livenessProbe:
          exec:
            command:
            - bash
            - "-c"
            - |
              set -ex
              mysqladmin -uroot -p$MYSQL_ROOT_PASSWORD ping &> /dev/null
          initialDelaySeconds: 30
          periodSeconds: 10
          timeoutSeconds: 5
        readinessProbe:
          exec:
            # Check we can execute queries over TCP (skip-networking is off).
            command:
            - bash
            - "-c"
            - |
              set -ex
              mysql -h 127.0.0.1 -uroot -p$MYSQL_ROOT_PASSWORD -e "SELECT 1" &> /dev/null
          initialDelaySeconds: 5
          periodSeconds: 2
          timeoutSeconds: 1
      - name: xtrabackup
        image: gcr.io/google-samples/xtrabackup:1.0
        ports:
        - name: xtrabackup
          containerPort: 3307
        command:
        - bash
        - "-c"
        - |
          set -ex
          cd /var/lib/mysql

          # Determine binlog position of cloned data, if any.
          if [[ -f xtrabackup_slave_info ]]; then
            # XtraBackup already generated a partial "CHANGE MASTER TO" query
            # because we're cloning from an existing slave.
            mv xtrabackup_slave_info change_master_to.sql.in
            # Ignore xtrabackup_binlog_info in this case (it's useless).
            rm -f xtrabackup_binlog_info
          elif [[ -f xtrabackup_binlog_info ]]; then
            # We're cloning directly from master. Parse binlog position.
            [[ `cat xtrabackup_binlog_info` =~ ^(.*?)[[:space:]]+(.*?)$ ]] || exit 1
            rm xtrabackup_binlog_info
            echo "CHANGE MASTER TO MASTER_LOG_FILE='${BASH_REMATCH[1]}',\
                  MASTER_LOG_POS=${BASH_REMATCH[2]}" > change_master_to.sql.in
          fi

          # Check if we need to complete a clone by starting replication.
          if [[ -f change_master_to.sql.in ]]; then
            echo "Waiting for mysqld to be ready (accepting connections)"
            until mysql -h 127.0.0.1 -uroot -p$MYSQL_ROOT_PASSWORD -e "SELECT 1"; do sleep 1; done

            echo "Initializing replication from clone position"
            # In case of container restart, attempt this at-most-once.
            mv change_master_to.sql.in change_master_to.sql.orig
            mysql -h 127.0.0.1 -uroot -p$MYSQL_ROOT_PASSWORD <<EOF
          $(<change_master_to.sql.orig),
            MASTER_HOST='mysql-0.mysql',
            MASTER_USER='root',
            MASTER_PASSWORD=$MYSQL_ROOT_PASSWORD,
            MASTER_CONNECT_RETRY=10;
          START SLAVE;
          EOF
          fi

          # Start a server to send backups when requested by peers.
          exec ncat --listen --keep-open --send-only --max-conns=1 3307 -c \
            "xtrabackup --backup --slave-info --stream=xbstream --host=127.0.0.1 --user=root"
        volumeMounts:
        - name: data
          mountPath: /var/lib/mysql
          subPath: mysql
        - name: conf
          mountPath: /etc/mysql/conf.d
        resources:
          requests:
            cpu: 100m
            memory: 100Mi
      volumes:
      - name: conf
        emptyDir: {}
      - name: config-map
        configMap:
          name: mysql
  volumeClaimTemplates:
    - metadata:
        name: data
      spec:
        accessModes: ["ReadWriteOnce"]
        resources:
          requests:
            storage: 10Gi
```

### 使用测试

由于我有两个节点，使用官方的方法可能会部署在另外一个节点上，这样就无法进行测试了。所以我直接使用service上的NodePort进行登陆。

首先查看svc

```
[root@k8s-master ~]# kubectl get svc
NAME         TYPE        CLUSTER-IP      EXTERNAL-IP   PORT(S)    AGE
kubernetes   ClusterIP   10.96.0.1       <none>        443/TCP    6d18h
mysql        ClusterIP   None            <none>        3306/TCP   6d17h
mysql-read   ClusterIP   10.102.14.245   <none>        3306/TCP   6d17h
```

进行登陆`mysql -h 10.102.14.245 -uroot -pwu97112500`（比较奇怪的是居然非得要用户名和密码都加上才能够登陆，只用用户名居然没有提示密码输入的选项）

执行以下命令：

```mysql
CREATE DATABASE demo; 
CREATE TABLE demo.messages (message VARCHAR(250)); 
INSERT INTO demo.messages VALUES ('hello');
USE demo;
SELECT * FROM messages;
```

如果一切正常的话，则表示集群内部登陆正常。

接下来准备集群外部的访问，其实直接使用NodePort好像就可以了，我需要进行一下实验。

### 集群外部访问mysql

因为我重新装了集群，没有了强大的负载均衡器，所以简单起见，我将使用NodePort作为访问方法。

执行命令`kubectl edit svc/mysql-read`进入修改界面，修改type为`NodePort`，同时给`spec/ports`中加入`nodePort:30000`条目。

再次执行`kubectl get svc`，可以看到`mysql-read   NodePort    10.102.14.245   <none>        3306:30000/TCP   6d18h`，说明映射完毕。

在本地电脑执行登陆选项。因为我用的是mysql shell，所以在界面内执行`\connect mysql://root:wu97112500@129.204.7.185:30000`，成功了。

创建数据库的时候要指明UTF-8

```mysql
CREATE DATABASE demo
  DEFAULT CHARACTER SET utf8
  DEFAULT COLLATE utf8_general_ci;
```

