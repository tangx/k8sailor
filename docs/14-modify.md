# modify

### 将 LabelSelector 转换为 Selector

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
