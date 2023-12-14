##  定义接口
### 1. 添加接口
```shell
kratos proto add api/template/v1/template.proto
```

### 2. 定义HTTP接口
```shell
service Template {
  rpc CreateTemplate (CreateTemplateRequest) returns (CreateTemplateReply){
    option (google.api.http) = {
      // 定义一个 GET 接口，并且把 name 映射到 CreateTemplate
      get: "/v1/template/create/{name}",
      // 可以添加附加接口
      additional_bindings {
        // 定义一个 POST 接口，并且把 body 映射到 CreateTemplate
        post: "/v1/template/create",
        body: "*",
      }
    };
  };
}
```
### 3. 定义错误
```shell
# 安装
go install github.com/go-kratos/kratos/cmd/protoc-gen-go-errors/v2@latest
```
```shell
syntax = "proto3";

package template.v1;

option go_package = "go_kratos_template/api/template/v1;v1";
option java_multiple_files = true;
option java_package = "template.v1";
option objc_class_prefix = "APITemplateV1";

enum ErrorReason {
  TEMPLATE_UNSPECIFIED = 0;
  USER_NOT_FOUND = 1;
}

```



### 2. 生成API
```shell
# 1. 生成所有 client API
make api
# 2. 生成指定API
# （1）生成 client 源码
kratos proto client api/template/v1/template.proto
# （2）生成 server 源码
kratos proto server api/template/v1/template.proto -t app/template/internal/service
```