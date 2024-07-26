- 部署
```shell
kubectl apply -f apps.yaml
```

- 滚动更新
```shell
kubectl rollout restart deployment openapi-business
```