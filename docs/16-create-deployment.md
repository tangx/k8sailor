# 创建 deployment

## 单 pod 单 container

使用 kubectl 命令创建如下
```bash
kubectl create deployment my-nginx-5 --image=nginx:alpine --replicas 3 --port 80
```

创建成功后查看结果， 大部分参数为默认参数。

```yaml
# kgd -o yaml my-nginx-5
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    app: my-nginx-5 # 根据 deployment 自动匹配名字自动生成
  name: my-nginx-5  # 用户指定
  namespace: default # 用户选择，默认为当前 namespace
spec:
  progressDeadlineSeconds: 600  # 默认值
  replicas: 3   # --replicas
  revisionHistoryLimit: 10
  selector:
    matchLabels:
      app: my-nginx-5
  strategy:
    rollingUpdate:
      maxSurge: 25%
      maxUnavailable: 25%
    type: RollingUpdate
  template:
    metadata:
      labels:
        app: my-nginx-5
    spec:
      containers:
      - image: nginx:alpine     # --image
        imagePullPolicy: IfNotPresent
        name: nginx # 根据镜像名字获取
        ports:      # --port
        - containerPort: 80
          protocol: TCP
        resources: {}
        terminationMessagePath: /dev/termination-log
        terminationMessagePolicy: File
      dnsPolicy: ClusterFirst
      restartPolicy: Always
      schedulerName: default-scheduler
      securityContext: {}
      terminationGracePeriodSeconds: 30
```


## 多 pod

```
kubectl create deployment my-nginx-7 --image=nginx:alpine --image=uyinn28/tools --replicas 1 --port 80
```