# Pod 的阶段(phase)与状态(status)

Pod 的生命周期 https://kubernetes.io/zh/docs/concepts/workloads/pods/pod-lifecycle/

![pod lifecycle](https://d33wubrfki0l68.cloudfront.net/aecab1f649bc640ebef1f05581bfcc91a48038c4/728d6/images/docs/pod.svg)



Pod 的 Status 不是 Phase。

Pod 的 Status 需要根据  Pod 中的 `ContainerStatuses` 进行计算得到。

```go
// extractPod 转换成业务本身的 Pod
func extractPod(item corev1.Pod) *Pod {

	reason := ""
	message := ""
    
    // 计算 Pod 在 Phase Running 时候的真实 Status
	for _, status := range item.Status.ContainerStatuses {
		if !status.Ready && status.State.Waiting != nil {
			reason = status.State.Waiting.Reason
			message = status.State.Waiting.Message
			break
		}
	}

	return &Pod{
		Name:       item.Name,
		Namespace:  item.Namespace,
		Images:     PodImages(item.Spec),
		NodeName:   item.Spec.NodeName,
		NodeIp:     item.Status.HostIP,
		CreateTime: item.CreationTimestamp.Time,
		PodIP:      item.Status.PodIP,
		Status: PodStatus{
			Phase:   item.Status.Phase,
			Message: message,
			Reason:  reason,
		},
		Labels: item.Labels,
	}
}
```