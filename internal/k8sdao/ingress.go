package k8sdao

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	netv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func GetIngressByName(ctx context.Context, namespace string, name string) (*netv1.Ingress, error) {
	opts := metav1.GetOptions{}

	ing, err := clientset.NetworkingV1().Ingresses(namespace).Get(ctx, name, opts)
	if err != nil {
		return nil, fmt.Errorf("Get Ingress failed: %w", err)
	}

	return ing, nil
}

func CreateIngressByName(ctx context.Context, namespace string, name string, backends []string) (*netv1.Ingress, error) {
	opts := metav1.CreateOptions{}

	spec, err := genIngressSpec(name, backends)
	if err != nil {
		return nil, err
	}
	meta := genIngressMeta(namespace, name)

	ing := &netv1.Ingress{
		ObjectMeta: meta,
		Spec:       *spec,
	}

	ing, err = clientset.NetworkingV1().Ingresses(namespace).Create(ctx, ing, opts)
	if err != nil {
		return nil, fmt.Errorf("create Ingress failed: %w", err)
	}

	return ing, err
}

// getMeta 返回 ingress Meta 信息
func genIngressMeta(ns string, name string) metav1.ObjectMeta {
	return metav1.ObjectMeta{
		Name:      name,
		Namespace: ns,
	}
}

func genIngressSpec(name string, endpoints []string) (*netv1.IngressSpec, error) {
	v1Rules := []netv1.IngressRule{}
	v1TLSs := []netv1.IngressTLS{}

	for _, endpoint := range endpoints {
		v1Rule, v1TLS, err := parseEndpoint(name, endpoint)
		if err != nil {
			return nil, err
		}

		v1Rules = append(v1Rules, *v1Rule)
		if v1TLS != nil {
			v1TLSs = append(v1TLSs, *v1TLS)
		}

	}

	return &netv1.IngressSpec{
		Rules: v1Rules,
		TLS:   v1TLSs,
	}, nil
}

// parseEndpoint 将 endpoint 信息解析成对应的 networkingv1.Ingress 信息
// endpoint: http://www.baidu.com/api/*?backend=service:8080&tls=secret1
func parseEndpoint(ingressName, endpoint string) (*netv1.IngressRule, *netv1.IngressTLS, error) {

	ep, err := url.Parse(endpoint)
	if err != nil {
		return nil, nil, err
	}

	backend := ep.Query().Get("backend")
	if len(backend) == 0 {
		backend = ingressName
	}
	v1Rule := getV1IngreeRule(ep.Host, ep.Path, backend)

	secret := ep.Query().Get("tls")
	v1Tls := getIngressTLS(ep.Host, secret)

	return v1Rule, v1Tls, nil
}

// getV1IngreeRule 返回对应的 networking Ingress Rule 信息
// host: www.baidu.com
// path: /v0/api/*
// backend: my-nginx-3:8088
func getV1IngreeRule(host string, path string, backend string) *netv1.IngressRule {
	ingPaths := getIngressPaths(path, backend)
	return &netv1.IngressRule{
		Host: getHost(host),
		IngressRuleValue: netv1.IngressRuleValue{
			HTTP: &netv1.HTTPIngressRuleValue{
				Paths: ingPaths,
			},
		},
	}
}

// getIngressTLS 返回对应的 networking Ingress TLS 信息
// host: wwww.baidu.com
// sercret: secretName
func getIngressTLS(host string, secret string) *netv1.IngressTLS {
	if secret == "" {
		return nil
	}

	return &netv1.IngressTLS{
		Hosts:      []string{host},
		SecretName: secret,
	}
}

// getIngressPaths return IngresPaths
// path: /v0/api/*
// service: my-nginx-3:80
func getIngressPaths(path string, service string) []netv1.HTTPIngressPath {

	path, typ := parsePath(path)
	svcName, backendPort := extractService(service)

	v1Path := netv1.HTTPIngressPath{
		Path:     path,
		PathType: &typ,
		Backend: netv1.IngressBackend{
			Service: &netv1.IngressServiceBackend{
				Name: svcName,
				Port: backendPort,
			},
		},
	}

	return []netv1.HTTPIngressPath{v1Path}
}

// getHost 返回不带端口的 Host 信息
// hostWithPort: www.baidu.com:80
func getHost(hostWithPort string) string {
	parts := strings.Split(hostWithPort, ":")
	if len(parts) == 0 {
		return ""
	}
	return parts[0]
}

// parsePath 解析并返回 具体Path 及 Path类型
// pathWithType: /v0/api*
func parsePath(pathWithType string) (string, netv1.PathType) {
	// pathWithType= /*
	lens := len(pathWithType)
	if lens == 0 {
		return "/", netv1.PathTypeExact
	}

	if pathWithType[lens-1] == '*' {
		return pathWithType[:lens-1], netv1.PathTypePrefix
	}

	return pathWithType, netv1.PathTypeExact
}

// extractService 提取 service:port
// serviceWithPort: my-nginx-3:8080
func extractService(serviceWithPort string) (svcName string, backendPort netv1.ServiceBackendPort) {

	parts := strings.Split(serviceWithPort, ":")

	if len(parts) == 1 {
		return parts[0], netv1.ServiceBackendPort{
			Number: 80,
		}
	}

	svcName, portName := parts[0], parts[1]
	if len(portName) == 0 {
		portName = "80"
	}
	port, err := strconv.Atoi(portName)
	if err != nil {
		return svcName, netv1.ServiceBackendPort{
			Name: portName,
		}
	}

	return svcName, netv1.ServiceBackendPort{
		Number: int32(port),
	}
}
