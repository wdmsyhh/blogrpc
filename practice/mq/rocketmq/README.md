# RocketMQ

## 问题

在运行的过程中总是发现有这样的一个问题：rocketmq-broker exited with code 253，也没日志打印。这里可能是挂载路径没有权限的问题。加上权限即可。

容器内是rocketmq用户运行，对应得gid:uid就是3000:3000

```shell
chmod -R 755 /opt/rocketmq && chown -R 3000:3000 /opt/rocketmq
```