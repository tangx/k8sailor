# modify

## 将 LabelSelector 转换为 Selector

`client-go` 提供了一个方法， 可以将 **Resource 中的 LabelSelector** 转换为 **Selector**, 并且 Selector 结构提供了一些常用的方法。 如 `String`

```go
import (
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func() {
    selector, _ := metav1.LabelSelectorAsSelector(dep.Spec.Selector)
    x := selector.String()
    fmt.Println(x)
}()
```

因此在使用 `GetXXXByLabels` 时， api 层 **可以考虑** 接收 `map[string]string` 类型的参数。 而在 biz 层应该将 **不同类型** 的参数 **统一** 转换为格式为 `key1=value1,key2=value2` 的 `string` 类型参数。 在 `dao` 层只接收 `string` 类型的 `string`。 这样就实现了前后一致性的问题。


## 自动刷新前端数据

在 vue 中， 如果数据是 **响应式** 数据， 那么当数据发生变化后， vue 会自动对页面进行刷新。

因此为了实现页面自动刷新， 需要保障：

1. 循环获取新数据
2. 解决新数据与老数据 **内容一致，顺序不一致** 导致的页面刷新。

### 循环任务

使用 `while` 循环请求接口不断请求数据， 需要注意的是

1. 循环一定要设置 **间隔时间**， typescript 中没有 `sleep` 函数， 可以使用 `Promise` 替代实现
2. 一定要设置 **循环开关**， 否则循环代码将一直在浏览器中的 **后台任务** 执行。 并且刷新一次就开启一个后台任务， 如果不加以限制， 机器风扇呜呜呜的转。 琢磨琢磨。
3. **循环开关** 可以放到页面的 `onMount() / onUnMount()` 两个 **生命周期钩子** 中实现。

```ts

let onOff = reactive({
  loop: false
})

const getAllByNamespaceLoop = async function () {
  while (onOff.loop) {
    let f = getAllByNamespace("default")

    // 间隔时间， ts 中没有 sleep 函数， 所以使用 Promise 实现
    await new Promise(f => setTimeout(f, 2000));
  }
}

onMounted(() => {
  onOff.loop = true
  console.log("onMounted: onOFF.loop", onOff.loop);

  getAllByNamespaceLoop()
})

onUnmounted(() => {
  onOff.loop = false

  console.log("onUnmounted: onOFF.loop", onOff.loop);
})

```

### 前端数据排序

typescript 中， 数组 `Array` 有一个方法 `sort( fn(n1,n2):number )`， 接收一个 **排序函数** 作为传参。

该 **排序函数** 接收 **两个参数** 表示元素， 返回一个 **数字** 类型表示是否交换位置。 

```ts
  // 对数组进行排序， 避免返回结果数据相同但顺序不同时， vue 不断重新渲染。
  let _items = resp.data.sort(
    (n1: Deployment, n2: Deployment) => {
      if (n1.name >= n2.name) {
        return 1
      }
      return -1
    }
  )
```

> https://stackoverflow.com/a/21689268


## 事件

使用 infromer 订阅 `Core/V1` 的 event 事件， 与 `EventsV1` 的 event 事件略有区别， 大体一致。

```go
	events, err := clientset.EventsV1().Events("default").List(ctx, v1.ListOptions{})
	events2, err := clientset.CoreV1().Events("default").List(ctx, v1.ListOptions{})
```

提取 event 事件的如下信息

```json
  "involvedObject": {
    "kind": "Pod",
    "namespace": "default",
    "name": "failed-nginx-6df5766f6d-vjn9n",
    "uid": "8726d44b-06b1-4d1c-9bad-efebf3fbb556",
    "apiVersion": "v1",
    "resourceVersion": "685855",
    "fieldPath": "spec.containers{nginx}"
  },
  "reason": "BackOff",
  "message": "Back-off pulling image \"nginx:alpine-11\"",
  "source": {
    "component": "kubelet",
    "host": "tangxin-test"
  },
```

并封装成一个 `map[string]Message` 的格式

```go
PodEvent["pod-namesapce-podname"] = Message{
    Reason: "BackOff",
    Message: "Back-off pulling image \"nginx:alpine-11\"",
}
```


