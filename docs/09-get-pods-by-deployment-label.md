# 通过 deployment label 获取 pod 信息

有了之前结构铺垫， 获取 Pod 还是很简单简单的。 其中需要注意的是 `ListOptions` 中的 `LabelSelector` 是一个字符串， 多组 `key=value` 之间使用 **逗号 `,`** 进行连接。

```go
labelSelector := `key1=value1,key2=value2,...`
```

而通过 client-go API 获取的 Deployment, Pod 等信息中的 `MatchLabel` 字段是一个 `map[string]string` 的 map。

因此， 在使用 `k8s client` 查询的时候， 需要对进行一些传参转换。

```go
// convertMapToSelector convert map to string, use comma connection: k1=v1,k2=v2
func convertMapToSelector(labels map[string]string) string {
	l := []string{}
	for k, v := range labels {
		l = append(l, fmt.Sprintf("%s=%s", k, v))
	}

	return strings.Join(l, ",")
}
```

## 获取 Pod

Pod 本身是 k8s 的一个最核心的概念， 独立于其他 Workloads ， 这点从 API 上也可以看出来。
Pod 的 API 是 `core v1` 而 Deployment 是 `apps v1`。

可以直接通过 Label 获取 Pod 信息

[/internal/k8sdao/pod.go](/internal/k8sdao/pod.go)

```go
func GetPodByLabels(ctx context.Context, namespace string, labels map[string]string) (*corev1.PodList, error) {

	opts := metav1.ListOptions{
		LabelSelector: convertMapToSelector(labels),
	}

	return clientset.CoreV1().Pods(namespace).List(ctx, opts)
}
```

## 通过 Deployment 获取 Pod

Pod 与其他 Workloads 之间的关联是 **弱关联 / 间接关联**， 以 Deployment 为例。  **Deployment 创建 ReplicaSet,  ReplicaSet 创建 Pod**

首先， 通过 `clientset` 的 `Get` 方法根据 **Name** 获取到 Deployment 对象， 在通过 Deployment 中的 Label 信息获取对应的 Pod 对象。 这里需要注意的是上述所讲的的 **Pod 与 Deployment 之间的弱关联关系**， 因为是通过标签匹配的，所以结果可能根本与 `Deployment` 无关。

假如现在有两个 Deployment 如下

```yaml
# dep1
kind: Deployment
metadata:
  creationTimestamp: null
  labels:   # 只有一组标签
    app: my-nginx-1
  name: my-nginx-1
# ... 省略
---

# dep2
kind: Deployment
metadata:
  creationTimestamp: null
  labels:   # 有两组标签， 其中 app 组与 dep1 相同
    app: my-nginx-1
    srv: my-nginx-2
  name: my-nginx-2
## ... 省略
```

如果单独的通过 `app=my-nginx-1` 标签来匹配，还会得到 dep2 的 Pod

```bash
kubectl get pod -l app=my-nginx-1
```

### 获取 ReplicaSet 再获取 Pod

[/internal/k8sdao/replicaset.go](/internal/k8sdao/replicaset.go)

因此， 在获取获取 Pod 信息之前， 应该先获取 `ReplicaSet`， 再获取 `Pod`

通过 deployment 的 label 获取 ReplicaSet

```bash
# kubectl get rs -l app=my-nginx-1

NAME                    DESIRED   CURRENT   READY   AGE
my-nginx-1-6d9577949b   1         1         1       4d14h
```

得到 rs 详细信息如下

```yaml
# kubectl get rs -o yaml my-nginx-1-6d9577949b
apiVersion: apps/v1
kind: ReplicaSet
metadata:
  # ... 省略
  labels:
    app: my-nginx-1
    pod-template-hash: 6d9577949b
  # ... 省略
```

通过带有 rs 的 label 进行查询

```bash
# kubectl pod -l app=my-nginx-1,pod-template-hash=6d9577949b

NAME                          READY   STATUS    RESTARTS   AGE
my-nginx-1-6d9577949b-bhwlp   1/1     Running   0          4d14h
```

查询出来的 Pod 结果符合预期

```yaml
# kubectl get pod my-nginx-1-6d9577949b-bhwlp -o yaml
apiVersion: v1
kind: Pod
metadata:
  creationTimestamp: "2021-09-23T16:13:51Z"
  generateName: my-nginx-1-6d9577949b-
  labels:
    app: my-nginx-1
    pod-template-hash: 6d9577949b
# ... 省略
```
