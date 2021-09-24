# 获取所有 deployments 数据

```
apis ->
  biz ->
  dao ->
```

## 三层结构
1. `apis` 只用于

## 重新调整目录结构

### internal

### confk8s

### 删除 cobra flags


## 获取 deployments 信息

1. api 处理用户请求参数， 请求 biz Operator 方法
2. biz Operator， 请求 k8sdao Operator， 并 **处理/过滤** 原始数据
3. k8sdao 与 cluster 交互， 返回原始数据。

> 有点问题， 三个模块， 三次同名方法。 有点麻烦。

