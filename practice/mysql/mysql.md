# Mysql

## 运行 mysql 容器

```shell
docker run -itd --name mysql-test -p 3306:3306 --network my_default --network-alias mysql -e MYSQL_ROOT_PASSWORD=root123 mysql:5.7

docker exec -it mysql-test bash

mysql -u root -p

show databases;

# mysql 数据库不能加中划线，这里用下划线
CREATE DATABASE portal_master;
# 选择数据库 
use portal_master;
# 查看当前使用的哪个库（3种方式）
show tables;
select database();
status;
```