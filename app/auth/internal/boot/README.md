# boot

## 启动配置载入

1. config.go 加载配置文件
2. bootstrap 加载配置
    1. Load()    加载配置
    2. Setting() 自定义配置Options模式
    3. Run() 启动需要提前boot的服务，如trace
    4. Close() defer时需要释放的资源