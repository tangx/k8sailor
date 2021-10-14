# 创建 ingress

> https://kubernetes.io/zh/docs/concepts/services-networking/ingress/

```bash
  # Create an ingress with a default backend
  kubectl create ingress ingdefault --class=default \
  --default-backend=defaultsvc:http \
  --rule="foo.com/*=svc:8080,tls=secret1" --dry-run -o yaml
```


```yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  creationTimestamp: null
  name: ingdefault
spec:
  defaultBackend:
    service:
      name: defaultsvc
      port:
        name: http
  ingressClassName: default
  rules:
  - host: foo.com
    http:
      paths:
      - backend:
          service:
            name: svc
            port:
              number: 8080
        path: /
        pathType: Prefix   # 匹配方式
  tls:
  - hosts:
    - foo.com
    secretName: secret1
status:
  loadBalancer: {}
```


### 路径类型 

Ingress 中的每个路径都需要有对应的路径类型（Path Type）。未明确设置 pathType 的路径无法通过合法性检查。当前支持的路径类型有三种：

`ImplementationSpecific` ：对于这种路径类型，匹配方法取决于 IngressClass。 具体实现可以将其作为单独的 pathType 处理或者与 Prefix 或 Exact 类型作相同处理。

`Exact` ：精确匹配 URL 路径，且区分大小写。

`Prefix` ：基于以 / 分隔的 URL 路径前缀匹配。匹配区分大小写，并且对路径中的元素逐个完成。 路径元素指的是由 / 分隔符分隔的路径中的标签列表。 如果每个 p 都是请求路径 p 的元素前缀，则请求与路径 p 匹配。

说明： 如果路径的最后一个元素是请求路径中最后一个元素的子字符串，则不会匹配 （例如：/foo/bar 匹配 /foo/bar/baz, 但不匹配 /foo/barbaz）。


> 不建议使用 `ImplementationSpecific`，因为这个参数实现的功能取决于 **IngressClass / IngressController** ， 对用户 **1. 或不可控**， **2. 或认知盲点**。 如果后期换了一个 Controller 或换了一个人维护则可能出现异常。

事实上， 在 `kubectl create ingress` 的时候， 也是通过判断 rule 中的 uri 最后一个字符是否为 `*` 号确定使用使用哪种规则。 有 `*` 使用 prefix， 没有使用 exact 。

```yaml

# http://www.baidu.com/v0/api$   # ImplementationSpecific  精确匹配
http://www.baidu.com/v0/api   # Exact  精确匹配

http://www.baidu.com/v1/api*    # Prefix 前缀匹配
```

### 