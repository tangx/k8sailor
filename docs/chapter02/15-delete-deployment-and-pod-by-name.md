# 根据名字删除 deployment 和 pod

> tag: https://github.com/tangx/k8sailor/tree/feat/15-delete-deployment-and-pod-by-name

调用 k8s api 没什么好说的。

**k8sdao**

```go

func DeleteDeploymentByName(ctx context.Context, namespace string, name string) error {
	opts := metav1.DeleteOptions{}
	return clientset.AppsV1().Deployments(namespace).Delete(ctx, name, opts)
}

```


**biz**

```go

type DeleteDeploymentByNameInput struct {
	Name      string `uri:"name"`
	Namespace string `query:"namespace"`
}

// DeleteDeploymentByName 根据名字删除 deployment
func DeleteDeploymentByName(ctx context.Context, input DeleteDeploymentByNameInput) error {
	err := k8sdao.DeleteDeploymentByName(ctx, input.Namespace, input.Name)
	if err != nil {
		return fmt.Errorf("k8s internal error: %w", err)
	}
	return nil
}
```

**api**

```go
func handlerDeleteDeploymentByName(c *gin.Context) {
	input := deployment.DeleteDeploymentByNameInput{}
	if err := ginbinder.ShouldBindRequest(c, &input); err != nil {
		httpresponse.Error(c, http.StatusBadRequest, err)
		return
	}

	if err := deployment.DeleteDeploymentByName(c, input); err != nil {
		httpresponse.Error(c, http.StatusInternalServerError, err)
		return
	}

	httpresponse.OK(c, true)
}
```