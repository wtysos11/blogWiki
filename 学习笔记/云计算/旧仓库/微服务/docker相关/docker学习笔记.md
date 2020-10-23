# docker学习笔记

标签：docker 学习笔记

[学习对象](https://docs.docker.com/)
[dockerfile从入门到实践，中文gitbook](https://yeasy.gitbooks.io/docker_practice/introduction/what.html)

## 速记 

[docker cheat sheet](https://www.docker.com/sites/default/files/Docker_CheatSheet_08.09.2016_0.pdf)

## 概念

### Images and containers

Image(镜像) and containers(容器)

A container is launched by running an image. An image is an executable package that includes everything needed to run an application--the code, a runtime, libraries, environment variables, and configuration files.

A container is a runtime instance of an image--what the image becomes in memory when executed (that is, an image with state, or a user process). You can see a list of your running containers with the command, docker ps, just as you would in Linux.

## 特殊注意


### 镜像

除了能够选择现有镜像作为基础镜像外，Docker还存在一个特殊的镜像scratch。这个镜像是个虚拟的概念，并不实际存在，表示一个空白的镜像。

如果以scratch为基础镜像的话，意味着你不以任何镜像作为基础，接下来所写的指令将作为镜像第一层存在。

### Dockerfile

Dockerfile中每一个指令都会建立一层。比如RUN指令，与其一条条写，不如

```Dockerfile
FROM debian:stretch

RUN buildDeps='gcc libc6-dev make wget' \
    && apt-get update \
    && apt-get install -y $buildDeps \
    && wget -O redis.tar.gz "http://download.redis.io/releases/redis-5.0.3.tar.gz" \
    && mkdir -p /usr/src/redis \
    && tar -xzf redis.tar.gz -C /usr/src/redis --strip-components=1 \
    && make -C /usr/src/redis \
    && make -C /usr/src/redis install \
    && rm -rf /var/lib/apt/lists/* \
    && rm redis.tar.gz \
    && rm -r /usr/src/redis \
    && apt-get purge -y --auto-remove $buildDeps
```

一条条写的话会重复创建多条镜像。每一个条指令执行结束后，engine都会commit这一层的修改，从而构成新的镜像。

Union FS有最大层数限制，现在是不得超过127层。

#### ENTRYPOINT 与 CMD 的区别

如果需要往容器中的命令添加参数的时候，CMD需要输入完整的命令，但是ENTRYPOINT可以只添加需要的参数。[来源](https://yeasy.gitbooks.io/docker_practice/image/dockerfile/entrypoint.html)

#### CMD的执行

[来源](https://stackoverflow.com/questions/43237183/dockerfile-how-use-cmd-or-entrypoint-from-base-image)

如果本容器没有CMD，engine会自动在最后调用base image的entrypoint或cmd。