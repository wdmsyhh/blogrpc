## consul 的安装和配置

```shell
docker run -d --name consul_1 -p 8500:8500 -p 8300:8300 -p 8301:8301 -p 8302:8302 -p 8600:8600/udp consul consul agent -dev -client=0.0.0.0

# 每次启动电脑时启动 consul
docker container update --restart=always 容器名字
```

- 相关 api: https://developer.hashicorp.com/consul/api-docs/agent/service