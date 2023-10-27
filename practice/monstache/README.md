# monstache

## 安装 monstache

[https://github.com/rwynn/monstache](https://github.com/rwynn/monstache)

[https://rwynn.github.io/monstache-site/](https://rwynn.github.io/monstache-site/)

```shell
git clone git@github.com:rwynn/monstache.git

go install
```

## 测试 monstache

```shell
monstache -f config.toml
```

然后创建 testdb 数据库和 member 集合并添加一条数据，数据会同步到 es 中，在 es 中的索引名称为 testdb.member

## 容器方式启动 monstache

- 参考 [https://github.com/rwynn/monstache/blob/rel6/docker/test/docker-compose.test.yml](https://github.com/rwynn/monstache/blob/rel6/docker/test/docker-compose.test.yml)