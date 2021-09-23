# 搭建 k3s 集群

安装过程参考

> https://tangx.in/2021/06/07/k3s-architecture-single-server/


## 安装

k3s 集群版本为 v1.21.4。 因此 k8s client-go sdk 的版本也需要安装对应版本

```bash

# curl -sfL http://rancher-mirror.cnrancher.com/k3s/k3s-install.sh | INSTALL_K3S_MIRROR=cn sh -

[INFO]  Finding release for channel stable
[INFO]  Using v1.21.4+k3s1 as release
[INFO]  Downloading hash http://rancher-mirror.cnrancher.com/k3s/v1.21.4-k3s1/sha256sum-amd64.txt
[INFO]  Downloading binary http://rancher-mirror.cnrancher.com/k3s/v1.21.4-k3s1/k3s
[INFO]  Verifying binary download
[INFO]  Installing k3s to /usr/local/bin/k3s

... 省略
```

## 初始化环境

通过命令创建一些工作负载， 以便后续 k8s api 调用查看

这里创建了两个 deployment: 
+ my-nginx-1 : 1 个 pod
+ my-nginx-2 : 2 个 pod

```bash
# kubectl create deployment my-nginx-1 --image=nginx:alpine
deployment.apps/my-nginx-1 created

# kubectl create deployment my-nginx-2 --image=nginx:alpine --replicas=2
deployment.apps/my-nginx-2 created
```

通过 kubectl 命令查看结果

```bash
# kubectl get pod

NAME                          READY   STATUS    RESTARTS   AGE
my-nginx-1-6d9577949b-98hzv   1/1     Running   0          105s
my-nginx-2-cd544c6f7-sf68x    1/1     Running   0          91s
my-nginx-2-cd544c6f7-zm974    1/1     Running   0          91s
```