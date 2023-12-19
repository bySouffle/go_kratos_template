# conf

## 项目配置文件

### 1. Bootstrap:

程序启动的必要配置

* APPInfo
    * ID
    * Name 服务名称
    * Version 版本号
    * Environment 环境配置 local
* General
    * SSL 开启TLS加密
    * Mode 工作模式 debug release
    * CSRF
    * Debug
    * Cert TLS
    * Key TLS
* Security
    * CookieName 加密字段名
    * JwtSecret 加密密钥
    * JwtTimeout 失效时间
* Experiment
    * EnableTrace 开启追踪
    * EnablePprof 开启性能分析
    * Trace
        * ServiceName
        * LocalAgentHostPort
        * CollectorEndpoint 将span发往jaeger-collector的服务地址 eg "http://localhost:14268/api/traces"
        * Exporter 输出模式 stdout | file | jaeger
        * TraceFilePath 输出到file时的路径
        * Token
* Server
    * HTTP
        * network 协议类型
        * addr 地址端口
        * timeout 超时时间
    * GRPC
        * network 协议类型
        * addr 地址端口
        * timeout 超时时间
    * Cron
        * network 协议类型
        * addr 地址端口
        * db redis数据库编号
        * password redis密码
        * dial_timeout 连接超时 默认 0.5s
        * read_timeout 默认 0.5s
        * write_timeout 默认 0.5s
        * MinIdleConn 默认 200
        * PoolSize 默认 100
        * PoolTimeout 默认 240s
        * Concurrency 并发数 默认 10
    * MQTT
        * network 协议类型
        * addr 地址端口
        * ClientID 客户端名
        * Username
        * Password
        * AutoReconnect 是否自动重连
        * MaxReconnectInterval 最大重连间隔 默认 5s
* Data
    * Database
        * driver 数据库驱动名 mysql sqlite
        * source root:root@tcp(127.0.0.1:3306)/test
        * Name 数据库名称
        * Addr 如果是 docker,可以替换为对应的服务名称，eg: db:3306
        * UserName
        * Password
        * ShowLog 是否打印所有SQL日志
        * MaxIdleConn 最大闲置的连接数，0意味着使用默认的大小2，小于0表示不使用连接池 默认 10
        * MaxOpenConn 最大打开的连接数, 需要小于数据库配置中的max_connections数 默认 60
        * ConnMaxLifeTime 单个连接最大存活时间，建议设置比数据库超时时长(wait_timeout)稍小一些 默认 14400s
        * SlowThreshold 慢查询阈值，设置后只打印慢查询日志，默认为200ms 默认 0.5s
    * Redis
        * network 协议类型
        * addr 地址端口
        * db redis数据库编号
        * password redis密码
        * dial_timeout 连接超时 默认 0.5s
        * read_timeout 默认 0.5s
        * write_timeout 默认 0.5s

### 2. Registry

* Registry
    * Consul
    * address 注册中心地址
    * scheme 协议
    * health_check
    * EndPoint
    * address
    * scheme
    * discovery 服务名 eg: bys.core.service

### 3. Logger

* Logger
    * Development 是否是开发环境，可选值：true 和 false，默认为false
    * DisableCaller 是否打印日志的文件调用文件和行号，即日志文件里的 caller 字段
    * DisableStacktrace
    * Encoding 打印的日志格式，默认为 json, 也可以修改为存文本格式，可选值为：json 和 console
    * Level 配置的日志级别，本地和测试可以开启为 debug，生产环境可以配置为 warn
    * Name 服务名，对应到日志里就是 app_id 字段
    * Writers 是日志需要输出到的位置 file console 可以两者同时选择
    * LoggerDir 保存目录 /tmp/logs/
    * LogRollingPolicy 日志切割方式: daily/hourly，默认按天(daily)进行切割
    * LogRotateDate
    * LogRotateSize
    * LogBackupCount 日志备份数

### 4. Folder

一些文件保存路径

### 5. Device

一些设备的配置

## 生成配置

```shell
# 在服务路径下
cd app/template/
protoc --proto_path=./internal \
--proto_path=../../third_party \
--go_out=paths=source_relative:./internal \
internal/conf/conf.proto
```