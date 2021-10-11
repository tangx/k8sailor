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
	typ, servicePorts, _ := parsePorts(ports)
	spec := corev1.ServiceSpec{
		Selector: map[string]string{
			"app": name,
		},
		Type:  typ,
		Ports: servicePorts,
	}

	return spec
}

type Port struct {
	NodePort   string
	Port       string
	TargetPort string
}

// parsePorts convert Simpl Port Phrase to ServicePort
// port:= !nodePort:port:targetPort
// port:= !nodePort:port
// port:= port:targetPort
func parsePorts(ports []string) (corev1.ServiceType, []corev1.ServicePort, error) {

	v1ServicePorts := []corev1.ServicePort{}

	isNodePort := false
	typ := corev1.ServiceTypeClusterIP
	for _, port := range ports {
		if port[0] == '!' {
			typ = corev1.ServiceTypeNodePort
			port = port[1:]
			isNodePort = true
		}

		parts := strings.Split(port, ":")
		if len(parts) < 2 {
			continue
		}

		port := corev1.ServicePort{}
		if len(parts) == 2 {
			port.Port = str2Int32(parts[0])
			port.TargetPort = intstr.Parse(parts[1])

			if isNodePort {
				port.NodePort = str2Int32(parts[0])
				port.Port = str2Int32(parts[1])
				port.TargetPort = intstr.Parse(parts[1])
			}
		}

		if len(parts) > 2 {
			port.NodePort = str2Int32(parts[0])
			port.Port = str2Int32(parts[1])
			port.TargetPort = intstr.Parse(parts[2])
		}

		port.Name = fmt.Sprintf("%d-%s", port.Port, port.TargetPort.String())
		v1ServicePorts = append(v1ServicePorts, port)
	}

	return typ, v1ServicePorts, nil
}

func str2Int32(s string) int32 {
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}

	return int32(i)
}
