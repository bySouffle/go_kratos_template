# Go Kratos Template
## 简介
本项目对Kratos进行了常用功能集成，包含发现注册服务，链路追踪，定时任务，日志集成，ws等，从而更方便的理解和使用kratos框架。

## 安裝环境
### 1. consul
[管理地址](http://127.0.0.1:8500)
```shell
docker pull hashicorp/consul
docker run -d --restart=always -p 8500:8500 -e CONSUL_BIND_INTER --name=consul hashicorp/consul agent -server -bootstrap -ui -client='0.0.0.0'
```


### 2. jaeger
[管理地址](http://127.0.0.1:16686)

```shell    
docker pull jaegertracing/all-in-one:latest
docker run -d --restart=always --name=jaeger -p 6831:6831/udp -p 16686:16686 -p 14268:14268 jaegertracing/all-in-one:latest 
```

### 3. EMQX
[管理地址](http://127.0.0.1:18083)

```shell   
docker pull emqx/emqx:5.3.2
docker run -d --name emqx -p 1883:1883 -p 8083:8083 -p 8084:8084 -p 8883:8883 -p 18083:18083  emqx:5.3.2
```

### 4. 运行
```shell
cd app/template/cmd/template/
go mod tidy
go build -ldflags "-X main.Version=x.y.z" .
./template -conf ../../configs/
```

### 5. docker 运行
```shell
# 构建基础镜像
cd deploy/
docker build -t base_image_go:v0.0.1A -f base_image_go.Dockerfile .

# 构建镜像网络采用主机模式以非root运行
cd ../app/template/cmd/template/
docker build -t go_kratos_template:v0.0.1 -f Dockerfile .
docker run -d --restart=always -e LOCAL_USER_ID=`id -u $USER` --name=go_kratos_template -v /data/conf:/data/conf -v /data/tls:/data/tls -v /etc/localtime:/etc/localtime --network host --privileged  --cap-add=SYS_ADMIN --cap-add=IPC_LOCK  go_kratos_template:v0.0.1
```

### 6. GRPCS测试
```shell
grpcui -insecure=false -plaintext=false  myregistry.domain.com:19000
```

## Install Kratos
```
go install github.com/go-kratos/kratos/cmd/kratos/v2@latest
```
## Create a service
```
# Create a template project
kratos new server

cd server
# Add a proto template
kratos proto add api/server/server.proto
# Generate the proto code
kratos proto client api/server/server.proto
# Generate the source code of service by proto file
kratos proto server api/server/server.proto -t internal/service

go generate ./...
go build -o ./bin/ ./...
./bin/server -conf ./configs
```
## Generate other auxiliary files by Makefile
```
# Download and update dependencies
make init
# Generate API files (include: pb.go, http, grpc, validate, swagger) by proto file
make api
# Generate all files
make all
```
## Automated Initialization (wire)
```
# install wire
go get github.com/google/wire/cmd/wire

# generate wire
cd cmd/server
wire
```

## Docker
```bash
# build
docker build -t <your-docker-image-name> .

# run
docker run --rm -p 8000:8000 -p 9000:9000 -v </path/to/your/configs>:/data/conf <your-docker-image-name>
```

