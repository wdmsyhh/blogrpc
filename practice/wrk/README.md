## 压测工具

安装使用：

https://www.cnblogs.com/quanxiaoha/p/10661650.html

https://segmentfault.com/a/1190000023212126

```shell
# 100线程 101http连接 1s压测时间
wrk -t100 -c101 -d1s http://127.0.0.1:8888/api/ping/baidu
```