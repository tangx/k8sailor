# k8sailor

k8s manage dashbooard

golang 练习的前后端

![](/docs/assets/img/gopher-pirate-medium.jpg)

## Summary

### 环境准备

1. [搭建 k3s 集群](/docs/01-install-k3s-cluster.md)
2. [使用 cobra 管理命令与参数](/docs/02-design-cobra-command.md)
3. [连接 k3s 集群并获取 deployments 信息](/docs/03-connect-cluster.md)
4. [使用 gin 初始化一个 API server](/docs/04-init-httpserver.md)
5. [设计 RESTful API 和响应请求](/docs/05-design-restful-api-and-response-data.md)

### 数据获取

1. [使用 api/biz/dao 分层结构管理数据请求，获取 deployment 数据](/docs/06-get-all-deployments.md)
2. [vue3 - 初始化 vue3 + vite2](/docs/07-initial-vue3-vite2.md)
3. [vue3 - 获取并展示 deployments 信息](/docs/08-fetch-and-display-deployments.md)
4. [通过 deployment label 获取 pod 信息](/docs/09-get-pods-by-deployment-label.md)
5. [vue3 - 使用 vue-router 和 less 优化展示页面](/docs/10-vue-router-and-less.md)
6. [vue3 - 展示 deployment 详情页](/docs/11-display-deployment-detail.md)
7. [deployment 副本数量设置 与 参数的有效性验证](/docs/12-deployment-scale-and-params-validate.md)
8. [使用 informer 监听变化并在本地缓存数据](./docs/13-k8s-informer.md)
9. [一些优化](./docs/14-some-optimize.md)
    + 将 LabelSelector 转换为 Selector
    + 自动刷新前端数据
    + 使用 informer 订阅 k8s event
    + defineProps 传入自定义类型

## feed

![wx-qrcode](https://tangx.in/assets/images/wx-qrcode.png)

