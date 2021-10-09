# 创建 deployment

> tag: https://github.com/tangx/k8sailor/tree/feat/16-create-deployment

使用 kubectl 命令创建如下
```bash
kubectl create deployment my-nginx-5 --image=nginx:alpine --replicas=3 --port=80
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
  selector:
    matchLabels:
      app: my-nginx-5
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
# ... 省略 ...
```

1. 在使用命令行传递参数的时候只传递了 `name, image, replicas, pods`

2. kubectl 根据传递的信息， 自动补全了一些信息
    + 以 name 补全了 `labels`： `app: my-nginx-5`
    + 以 image 生成了 `container name`； 如果传递多个 **相同名称** image 将会报错 **container name** 冲突。

接下来的工作就真没什么技术含量了， 就是最简单最无脑的拼凑字段。

在 `Annotations` 字段中， 也可以夹带很多 **私货**。 例如， 在 CI 中可以加入很多关于 **commit** 的信息， 例如 提交人， 提交信息 等。 一切想夹带的都可以放进去。

```go
type CreateDeploymentInput struct {
	Name     string
	Replicas *int32
	Images   []string
}

func CreateDeployment(ctx context.Context, namespace string, input CreateDeploymentInput) (*appsv1.Deployment, error) {
	labels := map[string]string{
		"app": input.Name,
	}
	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      input.Name,
			Namespace: namespace,
			Labels:    labels,
			// 在 CI 的时候， 可以在这里加上关键的 commit 信息。
			Annotations: map[string]string{
				"manager": "k8sailor",
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: input.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					Containers: containers(input.Images),
				},
			},
		},
	}
	opts := metav1.CreateOptions{}

	return clientset.AppsV1().Deployments(namespace).Create(ctx, dep, opts)

}

func containers(images []string) []corev1.Container {
	containers := make([]corev1.Container, len(images))
	for i, image := range images {
		container := corev1.Container{
			Image: image,
			Name:  imageName(i, image),
		}

		containers[i] = container
	}

	return containers
}
```

一个标准的 image 地址类似 `docker.io/libray/nginx:1.13-alpine`。 其中出现了多种特殊符号。 而 container name 只允许 **小写字母， 数字，`-`**， 所以需要改造一下。 为了解决相同镜像冲突的问题， 还在末尾加上了 container 的 id

```go
// image: docker.io/libray/nginx:1.13-alpine
func imageName(i int, image string) string {
	for _, char := range []string{"/", ":", "."} {
		image = strings.ReplaceAll(image, char, "-")
	}
	return fmt.Sprintf("%s-%d", image, i)
}
```

**写下来也发现了一个问题** ， 即使要拼凑字段， 也不应该写在一个结构体中。 

1. 本身 k8s 提供的诸多 Workloads 就是对 Pods 的 **调度管理** 的一层抽象，从而应对不同的场景。
2. Pod 与 Containers 之间是一个组合关系。 Container 和 InitContainer 本身是一样的， 只是在不同阶段运行的区分。

因此， 更应该将三个 Spec 分开处理， 这样就能更好、更方便的进行组合复用。

1. 向 Container 的构造函数中传入关键参数, 创建 Container 对象: `NewContainerSpec(params)`
2. 向 Pod 的构造函数中传入 **关键参数和 container 对象**， 创建 Pod 对象: `NewPodSpec(params, containers)`
3. 向 Deployment(或其他 Workload) 的构造函数中传入 **关键参数和 pod 对象**， 创建 `Deployment` 对象， `NewDeploymentSpec(params, pod)`
