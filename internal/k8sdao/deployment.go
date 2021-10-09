package k8sdao

import (
	"context"
	"errors"
	"fmt"
	"strings"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ListDeployments 返回 namespace 下的所有 deployments
// 兼容其他数据来源， 比如使用 informer 保存在本地的 cache， 不返回 DeploymentList 对象 而返回 []Deployment
func ListDeployments(ctx context.Context, namespace string) ([]appsv1.Deployment, error) {
	opts := metav1.ListOptions{}
	v1DepList, err := clientset.AppsV1().Deployments(namespace).List(ctx, opts)
	if err != nil {
		return nil, err
	}

	return v1DepList.Items, nil
}

// GetDeploymentByName 根据 name 获取 deployment
func GetDeploymentByName(ctx context.Context, namespace string, name string) (*appsv1.Deployment, error) {
	opts := metav1.GetOptions{}
	return clientset.AppsV1().Deployments(namespace).Get(ctx, name, opts)
}

// SetDeploymentReplicas 设置 namespace 大小
func SetDeploymentReplicas(ctx context.Context, namespace string, name string, replicas int) error {

	opts := metav1.GetOptions{}
	v1Scale, err := clientset.AppsV1().Deployments(namespace).GetScale(ctx, name, opts)
	if err != nil {
		return err
	}

	if replicas < 0 || replicas > 10 {
		return errors.New("invalid replicas number, must 0 < replica < 10")
	}
	// set new replicas
	v1Scale.Spec.Replicas = int32(replicas)

	upOpts := metav1.UpdateOptions{}
	v1Scale, err = clientset.AppsV1().Deployments(namespace).UpdateScale(ctx, name, v1Scale, upOpts)
	if err != nil {
		return err
	}

	return nil
}

func DeleteDeploymentByName(ctx context.Context, namespace string, name string) error {
	opts := metav1.DeleteOptions{}
	return clientset.AppsV1().Deployments(namespace).Delete(ctx, name, opts)
}

type CreateDeploymentInput struct {
	Name     string
	Replicas *int32
	Image    []string
}

func CreateDeployment(ctx context.Context, namespace string, input CreateDeploymentInput) {
	labels := map[string]string{
		"app": input.Name,
	}
	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      input.Name,
			Namespace: namespace,
			Labels:    labels,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: input.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					Containers: containers(input.Image),
				},
			},
		},
	}
	opts := metav1.CreateOptions{}

	clientset.AppsV1().Deployments(namespace).Create(ctx, dep, opts)
}

func containers(images []string) []corev1.Container {
	containers := make([]corev1.Container, len(images))
	for i, image := range images {
		container := corev1.Container{
			Image: image,
			Name:  imageName(i, image),
		}

		containers[i] = container
	}

	return containers
}

func imageName(i int, image string) string {
	name := strings.ReplaceAll(image, "/", "-")
	name = strings.ReplaceAll(name, ":", "-")
	return fmt.Sprintf("%s-%d", name, i)
}
