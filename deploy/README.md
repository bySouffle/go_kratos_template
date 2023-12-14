# 部署
## 容器内已非root权限部署
1. 构建基础镜像
```shell
docker build -t base_image_go:v0.0.1 -f base_image_go.Dockerfile .
```

2. 将编译完成的命名为`go_kratos_template`APP拷贝到当前目录并构建
```shell

```
