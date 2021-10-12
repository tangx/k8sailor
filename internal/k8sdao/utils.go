package k8sdao

import (
	"fmt"
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ConvertMapToSelector convert map to string, use comma connection: k1=v1,k2=v2
func ConvertMapToSelector(labels map[string]string) string {
	l := []string{}
	for k, v := range labels {
		l = append(l, fmt.Sprintf("%s=%s", k, v))
	}

	return strings.Join(l, ",")
}

func objectMeta(ns string, name string) metav1.ObjectMeta {
	return metav1.ObjectMeta{
		Namespace: ns,
		Name:      name,
		Labels:    map[string]string{"app": name},
	}
}
