## pod

- 使用k8s pod执行echo命令
```yaml
apiVersion: v1
kind: Pod
metadata:
  name: echo-pod
spec:
  containers:
    - name: echo-container
      image: ubuntu:jammy
      command: ["/bin/bash", "-c"]
      args: ["echo '你好$小明'"]
  restartPolicy: Never
```

```shell
kubectl apply -f echo.yaml
kubectl get pods | grep echo
kubectl logs -f echo-pod  # 发现输出乱码  你好$å��明
kubectl delete pod echo-pod
```