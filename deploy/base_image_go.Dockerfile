#   Go 基础开发镜像 包含
#   系统权限 gosu
FROM golang:1.20.4 AS builder
MAINTAINER bySouffle <bysouffle@gmail.com>
LABEL authors="bys"

ENV DEBIAN_FRONTEND=noninteractive
ENV GOPROXY https://proxy.golang.com.cn,direct

RUN       apt-get update \
      &&  apt-get install -y --no-install-recommends \
          ca-certificates  \
          netbase \
          gosu \
          && rm -rf /var/lib/apt/lists/ \
          && apt-get autoremove -y && apt-get autoclean -y

COPY ./docker-entrypoint.sh /usr/local/bin/docker-entrypoint.sh
RUN chmod a+x /usr/local/bin/docker-entrypoint.sh

#docker build -t base_image_go:v0.0.1A -f base_image_go.Dockerfile .
#docker run -d --restart=always -e LOCAL_USER_ID=`id -u $USER` --name=base_image_go -v /opt:/opt -v /etc/localtime:/etc/localtime --network host --privileged  --cap-add=SYS_ADMIN --cap-add=IPC_LOCK  base_image_go:v0.0.1A
