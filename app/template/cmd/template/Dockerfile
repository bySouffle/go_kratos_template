FROM base_image_go:latest AS builder
MAINTAINER bySouffle <bysouffle@gmail.com>
LABEL authors="bys"

ENV DEBIAN_FRONTEND=noninteractive
WORKDIR /app

COPY ./template /app

RUN gosu nobody true \
&& useradd -s /bin/bash -m user

ENTRYPOINT ["/usr/local/bin/docker-entrypoint.sh"]

VOLUME /data/conf
VOLUME /data/tls

CMD ["./template", "-conf", "/data/conf"]


#docker build -t go_kratos_template:v0.0.1 -f Dockerfile .
#docker run -d --restart=always -e LOCAL_USER_ID=`id -u $USER` --name=go_kratos_template -v /data/conf:/data/conf -v /data/tls:/data/tls -v /etc/localtime:/etc/localtime --network host --privileged  --cap-add=SYS_ADMIN --cap-add=IPC_LOCK  go_kratos_template:v0.0.1


#
#COPY . /src
#WORKDIR /src
#
#RUN GOPROXY=https://goproxy.cn make build
#
#FROM debian:stable-slim
#
#RUN apt-get update && apt-get install -y --no-install-recommends \
#		ca-certificates  \
#        netbase \
#        && rm -rf /var/lib/apt/lists/ \
#        && apt-get autoremove -y && apt-get autoclean -y
#
#COPY --from=builder /src/bin /app
#
#WORKDIR /app
#
#EXPOSE 8000
#EXPOSE 9000
#VOLUME /data/conf
#
#CMD ["./server", "-conf", "/data/conf"]
