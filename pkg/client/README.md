###  

## 1. EMQX 安装

```shell
docker pull emqx/emqx:5.3.2
docker run -d --name emqx -p 1883:1883 -p 8083:8083 -p 8084:8084 -p 8883:8883 -p 18083:18083  emqx:5.3.2
```

## 2. Docker部署注意事项

1. 如果您想持久化 EMQX Docker 容器中生成的数据，则需要保留以下目录，这样即使容器不再存在，数据也能持久保存。

```shell
/opt/emqx/data
/opt/emqx/log
```

2. 启动容器并挂载目录：

```shell
docker run -d --name emqx \
  -p 1883:1883 -p 8083:8083 \
  -p 8084:8084 -p 8883:8883 \
  -p 18083:18083 \
  -v $PWD/data:/opt/emqx/data \
  -v $PWD/log:/opt/emqx/log \
  emqx/emqx:5.3.2
```

3. 在 Docker 环境中，localhost或127.0.0.1指的是容器自己的内部网络接口，而不是主机的内部网络接口。要访问主机上运行的服务，请使用主机的
   IP 地址或使用主机网络设置。如果您使用的是 Docker for Mac 或 Docker for Windows，则可以用作host.docker.internal主机地址。
4. EMQX 使用data/mnesia/<node_name>目录来存储数据。选择一个稳定的标识符（例如主机名或完全限定域名 (FQDN)
   ）作为节点名称至关重要。这种做法可以避免因节点名称更改而导致的数据丢失