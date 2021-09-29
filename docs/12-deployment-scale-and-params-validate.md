# deployment 副本数量设置 与 参数的有效性验证

## deployment scale

```bash
kubectl scale deployment my-nginx-1 --replicas 1
```

在 client-go sdk 中， scale 参数是一个对象， 因此不能直接传入 **一个数字**。

1. 需要通过 `GetScale()` 方法获取到 `*autoscalingv1.Scale` 对象。
2. 修改 `Scale` 对象中的 `Replicas` 数值。
3. 使用 `UpdateScale()` 方法更新设置。

[SetDeploymentReplicas](/internal/k8sdao/deployment.go#L25)

## params validtor

参数验证在任何情况下都不能放松警惕， 尤其是 **边界验证** 和 **0值混淆** 。

对于参数的验证， 可以自己在业务代码中实现， 也可以使用已有的公共库。 gin 默认使用的是 https://github.com/go-playground/validator


```go
// SetDeploymentReplicasInput 调整 deployment pod 数量参数
// Replicas 为了避免 **0值** 影响。
//   1. 使用为 *int 指针对象， 自行在业务逻辑中进行校验
//   2. 另外也可以使用， `binding` tag， 由 gin 框架的 valicator 帮忙校验。 https://github.com/go-playground/validator
// Namespace 设置了默认值， 如果请求不提供将由 gin 框架自己填充。
type SetDeploymentReplicasInput struct {
	Namespace string `query:"namespace,default=default"`
	Name      string `uri:"name"`
	Replicas  *int   `query:"replicas" binding:"required"`
}
```

## vue 