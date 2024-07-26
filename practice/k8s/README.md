- 部署
```shell
kubectl apply -f apps.yaml
```

- 滚动更新
```shell
kubectl rollout restart deployment openapi-business
```

- job测试
```shell
kubectl apply -f job.yaml
kubectl get pods | grep echo
kubectl get jobs | grep echo
kubectl logs job/echo-job
kubectl delete job echo-job
```