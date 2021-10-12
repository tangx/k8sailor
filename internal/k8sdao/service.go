package k8sdao

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func CreateServcieByName(ctx context.Context, namespace string, name string, ports []string) (*corev1.Service, error) {
	opts := metav1.CreateOptions{}
	svc := &corev1.Service{
		ObjectMeta: objectMeta(namespace, name),
		Spec:       serviceSpec(name, ports),
	}
	svc, err := clientset.CoreV1().Services(namespace).Create(ctx, svc, opts)
	if err != nil {
		return nil, fmt.Errorf("Create Service Failed: %w", err)
	}

	return svc, nil
}

func GetServiceByName(ctx context.Context, namespace string, name string) (*corev1.Service, error) {
	opts := metav1.GetOptions{}
	svc, err := clientset.CoreV1().Services(namespace).Get(ctx, name, opts)
	if err != nil {
		return nil, fmt.Errorf("Get Service Failed: %w", err)
	}

	return svc, nil
}

func serviceSpec(name string, ports []string) corev1.ServiceSpec {

	spec := corev1.ServiceSpec{
		Selector: map[string]string{
			"app": name,
		},
	}

	parsePorts(ports, &spec)
	return spec
}

type Port struct {
	NodePort   string
	Port       string
	TargetPort string
}

// parsePorts convert Simpl Port Phrase to ServicePort
// port:= !nodePort:port:targetPort
// port:= !port:targetPort
// port:= port:targetPort,None
func parsePorts(ports []string, spec *corev1.ServiceSpec) {

	typ := corev1.ServiceTypeClusterIP
	v1Ports := []corev1.ServicePort{}

	for _, port := range ports {
		symbal := false

		if len(port) == 0 {
			continue
		}

		// 0. external service
		if external, symbal := isExternalName(port); symbal {
			spec.ClusterIP = ""
			spec.Type = corev1.ServiceTypeExternalName
			spec.ExternalName = external

			return
		}

		// 1. headless service
		if port, symbal = isHeadless(port); symbal {
			spec.ClusterIP = "None"
		}

		// 2. nodePort:port:targetPort
		if port, symbal = isNodePort(port); symbal {
			typ = corev1.ServiceTypeNodePort
		}

		// 3. handle ports mapping
		parts := strings.Split(port, ":")
		v1Port := corev1.ServicePort{}
		if len(parts) == 1 {
			v1Port.Port = str2Int32(parts[0])
			v1Port.TargetPort = intstr.Parse(parts[0])
		}
		if len(parts) == 2 {
			v1Port.Port = str2Int32(parts[0])
			v1Port.TargetPort = intstr.Parse(parts[1])
		}
		if len(parts) > 2 {
			v1Port.NodePort = str2Int32(parts[0])
			v1Port.Port = str2Int32(parts[1])
			v1Port.TargetPort = intstr.Parse(parts[2])
		}

		v1Port.Name = fmt.Sprintf("%d-%s", v1Port.Port, v1Port.TargetPort.String())
		v1Ports = append(v1Ports, v1Port)
	}

	spec.Type = typ
	spec.Ports = v1Ports
}

func str2Int32(s string) int32 {
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}

	return int32(i)
}

// 是否为 Headless 服务
// 以 # 开头: #port:targetPort
func isHeadless(port string) (_port string, ok bool) {
	if len(port) == 0 {
		return "", false
	}
	if port[0] == '#' {
		return port[1:], true
	}

	return port, false
}

// isNodePort 是否暴露 NodePort
// 以 ! 开头: !nodePort:port:targetPort
func isNodePort(port string) (_port string, ok bool) {
	if len(port) == 0 {
		return "", false
	}
	if port[0] == '!' {
		return port[1:], true
	}

	return port, false
}

// isExternalName 是否为 外部域名服务
// 以 @ 开头 :  @external.com
func isExternalName(external string) (_external string, ok bool) {
	if len(external) == 0 {
		return "", false
	}

	if external[0] == '@' {
		return external[1:], true
	}

	return external, false
}
