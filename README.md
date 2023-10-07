# blogrpc

## grpc-gateway 和 gin 结合适用的 demo

- 执行 `./proto/gen_stub.sh` proto 生成对应的 go 文件
  - 在 blogrpc 目录下执行 `./build.sh` 构建镜像
    - 注意：如果打包成镜像，代码中所有访问另一个容器的地方，比如：localhost:8081 需要改成 {容器别名或容器名}:{容器内部端口}
    - 不打包成镜像直接代码调试的时候需要在 /etc/hosts 文件中配置 {容器别名或容器名} 对应的容器 IP，但是代码中访问的应该是 {容器别名或容器名}:{容器内部端口}
    - 启动容器：
    ```shell
    docker network create my_default
    # 容器别名是 business
    docker run -it --rm --name businesstest -p 8080:8080 --network my_default --network-alias business business:v1
    # 容器别名是 hello
    docker run -it --rm --name hellotest -p 8081:8081 --network my_default --network-alias hello hello:v1
    # 容器别名是 member
    docker run -it --rm --name membertest -p 8082:8082 --network my_default --network-alias member member:v1
  
    # 启动 MongoDB 容器，设置用户名密码，容器别名是 mongo
    docker run -itd --name mongotest -p 27012:27017 --network my_default --network-alias mongo mongo:4.2 --auth
    docker exec -it mongotest mongo admin
    db.createUser({ user:'root',pwd:'root',roles:[ { role:'userAdminAnyDatabase', db: 'admin'},"readWriteAnyDatabase"]});
    db.auth('root', 'root');
    # 查看所有数据库
    show dbs;
    # 查看当前使用的库
    db;
    # 查看当前库中的表
    show collections;
    # 切换数据库
    use portal-master;
    # 添加一条数据,如果是容器内部需要使用别名加内部端口
    db.accountDBConfig.insert(
        {
        "title" : "",
        "accountId" : ObjectId("5e7873c4c3307000f272c9e2"),
        "dsn" : "mongodb://root:root@127.0.0.1:27012/portal-tenants-shared?authSource=admin",
        "options" : {
            "replicaSet" : "rs0"
        },
        "createdAt" : ISODate("2023-06-07T08:14:11.583Z")
        }
    )
  
    # 启动 Redis 容器
    docker pull redis:latest
    docker run -itd --name redis -p 6379:6379 redis
    docker exec -it redis /bin/bash
    redis-cli
    set key1 value1
    ```
  
## test

```shell
docker run -it --rm -p 9091:9091 --network my_default registry.ap-southeast-1.aliyuncs.com/yhhnamespace/blogrpc-openapi-business:local

docker run -it --rm -p 1701:1701 --name blogrpc-hello --network my_default registry.ap-southeast-1.aliyuncs.com/yhhnamespace/blogrpc-hello:local
```

## 问题

- server selection timeout

需要查看数据库 host 和端口是否正确