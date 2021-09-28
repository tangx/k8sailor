package pod

import corev1 "k8s.io/api/core/v1"

// 返回 Pod 的镜像列表
func PodImages(podSpec corev1.PodSpec) []string {
	images := containerImages(podSpec.Containers)
	initImages := containerImages(podSpec.InitContainers)

	return append(images, initImages...)
}

// containerImages 返回容器的镜像列表
func containerImages(containers []corev1.Container) []string {
	n := len(containers)
	images := make([]string, n)
	for i := 0; i < n; i++ {
		images[i] = containers[i].Image
	}
	return images
}
