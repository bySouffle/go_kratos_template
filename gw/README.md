# Apinto

## 1. 安装apinto网关

* 配置文件

```/data/apinto/config-dir/config.yml
version: 2
#certificate: # 证书存放根目录
#  dir: /etc/apinto/cert
client:
  #advertise_urls: # open api 服务的广播地址
  #- http://127.0.0.1:9400
  listen_urls: # open api 服务的监听地址
    - http://0.0.0.0:9400
  #certificate:  # 证书配置，允许使用ip的自签证书
  #  - cert: server.pem
  #    key: server.key
gateway:
  #advertise_urls: # 转发服务的广播地址
  #- http://127.0.0.1:9400
  listen_urls: # 转发服务的监听地址
    - https://0.0.0.0:8099
    - http://0.0.0.0:8099
peer: # 集群间节点通信配置信息
  listen_urls: # 节点监听地址
    - http://0.0.0.0:9401
  #advertise_urls: # 节点通信广播地址
  # - http://127.0.0.1:9400
  #certificate:  # 证书配置，允许使用ip的自签证书
  #  - cert: server.pem
  #    key: server.key
```

* 启动docker

```shell
docker pull eolinker/apinto-gateway

docker run -td  -p 8099:8099 -p 9400:9400 \
-v /data/apinto/data-dir:/var/lib/apinto \
-v /data/apinto/log-dir:/var/log/apinto \
-v /data/apinto/config-dir/config.yml:/etc/apinto/config.yml \
--name=apinto_node  eolinker/apinto-gateway:latest

docker run -td --restart=always --net=host -v /data/apinto/data-dir:/var/lib/apinto -v /data/apinto/log-dir:/var/log/apinto -v /data/apinto/config-dir/config.yml:/etc/apinto/config.yml -v /data/tls:/data/tls --name=apinto_node  eolinker/apinto-gateway:latest
```

## 2. 安装apinto-dashboard

```shell
# 新建docker网段
docker network create --driver bridge --subnet=172.100.0.0/24 --gateway=172.100.0.1 apinto

# 安装redis_cluster
docker run -dt --name redis_cluster --restart=always --net=host -v /var/lib/apinto/redis-cluster/data:/usr/local/cluster_redis/data -e REDIS_PWD=root  -e PORT=6380 eolinker/cluster-redis:6.2.7

docker run -dt --name apinto-dashboard --restart=always \
-p 18080:8080 -v /var/log/apinto/apinto-dashboard/work:/apinto-dashboard/work \
--network=apinto --privileged=true \
-e MYSQL_USER_NAME=root -e MYSQL_IP=apinto_mysql \
-e MYSQL_PWD=root -e MYSQL_PORT=3306 -e MYSQL_DB=apinto \
-e REDIS_ADDR=127.0.0.1:6380 \
-e REDIS_PWD=root eolinker/apinto-dashboard

# EX
docker run -dt --name apinto-dashboard --restart=always --net=host -v /var/log/apinto/apinto-dashboard/work:/apinto-dashboard/work --privileged=true -e MYSQL_USER_NAME=root -e MYSQL_IP=127.0.0.1 -e MYSQL_PWD=root -e MYSQL_PORT=3306 -e MYSQL_DB=apinto -e REDIS_ADDR=127.0.0.1:6380 -e REDIS_PWD=root eolinker/apinto-dashboard
```

## 3. 登录 http://127.0.0.1:8080/
