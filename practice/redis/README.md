
[Docker 启动 Redis 并添加密码](https://jueee.github.io/2021/03/2021-03-14-Docker%E5%90%AF%E5%8A%A8Redis%E5%B9%B6%E6%B7%BB%E5%8A%A0%E5%AF%86%E7%A0%81/)


```shell
docker run -itd --name redis-test -p 6379:6379 --network my_default --network-alias redis redis:5.0.12 --requirepass "root123"

docker exec -it redis-test bash

redis-cli

auth root123

set key1 value1

get key1
```