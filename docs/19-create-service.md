# 为 Deployment 创建 Service

```bash
kubectl create service clusterip nginx-web --clusterip="port:targetPort"
kubectl create service clusterip nginx-web --clusterip="8082:80"
kubectl create service nodeport  nginx-web --clusterip="8081:80"
```

需要注意, 使用 `kubectl get service` 查看到的 `Ports` 的展示结果为 `port:nodePort`， 而 `targetPort` 不展示。

```bash
# kubectl get service
NAME                    TYPE        CLUSTER-IP      EXTERNAL-IP   PORT(S)         AGE
demo-nginx-nodeport-3   NodePort    10.43.181.29    <none>        80:32425/TCP    4s
```

## port, targetPort, nodePort

![nodeport-port-targetport](./assets/img/19/nodeport-port-targetport.png)

端口映射中的四个 **比较关键** 的要素:

1. name: 避免端口相同时，默认名字冲突
2. port: service 对外提供服务的端口
3. targetPort: service 指向的 pod 的端口， 即 pod 对外服务的端口。
4. nodePort: node 对外提供服务的端口， 通过 kube-proxy 修改 iptables 将流量转发到 service 上。

其中， targetPort 可以是 `string / int32` 的复合类型， 定义如下。

```go
// k8s.io/api@v0.21.4/core/v1/types.go

// ServicePort contains information on service's port.
type ServicePort struct {
	// Number or name of the port to access on the pods targeted by the service.
	// Number must be in the range 1 to 65535. Name must be an IANA_SVC_NAME.
	// If this is a string, it will be looked up as a named port in the
	// target Pod's container ports. If this is not specified, the value
	// of the 'port' field is used (an identity map).
	// This field is ignored for services with clusterIP=None, and should be
	// omitted or set equal to the 'port' field.
	// More info: https://kubernetes.io/docs/concepts/services-networking/service/#defining-a-service
	// +optional
	TargetPort intstr.IntOrString `json:"targetPort,omitempty" protobuf:"bytes,4,opt,name=targetPort"`
}
```

在 **apimachinery`(k8s.io/apimachinery@v0.21.4/pkg/util/)`** 包中， 提供了很多针对 k8s 对象的常用方法。 

### 解析 port

为了能简化 api 请求的参数， 因此对 clusterIp 和 nodePort 使用了符号表示： 默认为 clusterIp， 如果有 **叹号`!`** 则为 nodePort。

```go
port    // clusterIp, port 与 targetPort 一致
port:targetPort // clusterIp, port 与 targetPort 可能不一致

!port:targetPort // nodeport, 端口号随机: port 与 targetPort 可能不一致
!nodePort:port:targetPort // nodeport, 指定端口号, 端口可能因为被使用而创建失败
```


## headless service

对 statefuleset 有效， 对 deployment 无效

> https://kubernetes.io/zh/docs/concepts/workloads/controllers/statefulset/

![service-headless](./assets/img/19/service-headless.png)

使用 `headless` 之后， k8s 将不再创建 `service` 进行 pod 的负载均衡。 取而代之的是 **DNS** 将每个 pod 直接解析暴露， 域名规则 `podName.serviceName.namespace.Cluster`


### 解析 Headless Port

```bash
kubectl create service clusterip my-nginx-web  --clusterip="None" --tcp=8088:80
```

从 kubectl 的命令中可以看到， `Headless Service` 可以被认为 `clusterip` 的一个子类， 其特殊之处就是 `ClusterIp: None`

鉴于此， 对之前的 Port 规则进行了一定扩展。 

规则基本类似， 采用了的新的符号 `#` 表示 Headless 服务。 `#` 在很多地方表示注释， 注释对外看不见，因此用以表示 `Headless`。

```go
#port:targetPort // headless 
```

本身 **NodePort 和 Headless** 就是不兼容的， 由于 `#, !` 在同一个位置， 也一定程度上避免了 **误写**


## external name

