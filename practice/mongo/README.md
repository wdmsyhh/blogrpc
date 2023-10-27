# 搭建本地测试环境

## 创建单节点 mongo 集群

- 使用 docker compose 启动 mongo 容器

```yaml
version: '3'
services:
  mongo:
    image: mongo:4.4
    ports:
      - "27017:27017"
    command: "--replSet rs0"

networks:
  default:
    external:
      name: my_default
```

- 进入容器初始化副本集

```shell
# 执行
mongo
# 初始化副本集
rs.initiate({_id: "rs0", members: [{_id: 0, host: "mongo:27017"}]})
# 查看副本集状态
rs.status()
```