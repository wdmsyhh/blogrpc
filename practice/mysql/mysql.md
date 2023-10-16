# Mysql

## 运行 mysql 容器

```shell
docker run -itd --name mysql-test -p 3306:3306 -e MYSQL_ROOT_PASSWORD=root123 mysql:5.7
```