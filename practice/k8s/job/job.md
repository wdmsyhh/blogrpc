## job

- job测试
```shell
kubectl apply -f job.yaml
kubectl get pods | grep echo
kubectl get jobs | grep echo
kubectl logs job/echo-job  #输出 Hello, Kubernetes! 你好$å��明
kubectl delete job echo-job
```