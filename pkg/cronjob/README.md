## 1. Asynq

## [server.go](server.go)

    周期任务服务，用于注册触发的回调函数

## [periodic.go](periodic.go)

    可修改周期的定时器，用于注册回调周期

## [register.go](register.go)

    不可修改周期的定时器，用于注册回调周期

## 2. Dash

https://github.com/hibiken/asynqmon#readme

```shell
docker run --rm     --name asynqmon   --net=host  -e REDIS_DB=1 -e REDIS_PASSWORD=passwd hibiken/asynqmon
```

```shell
Flag	                        Env	                      Description	                                                        Default
--port(int)                     PORT                      port number to use for web ui server	                                8080
---redis-url(string)	        REDIS_URL                 URL to redis or sentinel server. See godoc for supported format       ""
--redis-addr(string)	        REDIS_ADDR                address of redis server to connect to	                                "127.0.0.1:6379"
--redis-db(int)	                REDIS_DB                  redis database number	                                                0
--redis-password(string)        REDIS_PASSWORD            password to use when connecting to redis server                       ""
--redis-cluster-nodes(string)   REDIS_CLUSTER_NODES       comma separated list of host:port addresses of cluster nodes          ""
--redis-tls(string)             REDIS_TLS                 server name for TLS validation used when connecting to redis server	""
--redis-insecure-tls(bool)      REDIS_INSECURE_TLS        disable TLS certificate host checks	                                false
--enable-metrics-exporter(bool) ENABLE_METRICS_EXPORTER   enable prometheus metrics exporter to expose queue metrics	        false
--prometheus-addr(string)       PROMETHEUS_ADDR           address of prometheus server to query time series                     ""
--read-only(bool)               READ_ONLY                 use web UI in read-only mode	                                        false
```