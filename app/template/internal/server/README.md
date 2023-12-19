# Server

为http和grpc实例的创建和配置，以及注册对应的service

## HTTP

## GRPC

## MQTT

## Register

## Cron

### 定时任务服务

1. 在 `app/template/internal/task` 中定义定时任务
2. 在 `app/template/internal/server/cron.go` 中注册定时任务，其中
    1. `RegisterCronHandler`，注册定时任务的回调函数
    2. `RegisterCronJob`，注册定时任务的执行周期
    3. 任务以topic完成执行函数和周期的匹配
