package k8sdao

import (
	"fmt"
	"strings"
)

// ConvertMapToSelector convert map to string, use comma connection: k1=v1,k2=v2
func ConvertMapToSelector(labels map[string]string) string {
	l := []string{}
	for k, v := range labels {
		l = append(l, fmt.Sprintf("%s=%s", k, v))
	}

	return strings.Join(l, ",")
}
