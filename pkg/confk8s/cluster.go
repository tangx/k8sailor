package confk8s

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type Client struct {
	KubeConfig string `env:""`
	clientset  *kubernetes.Clientset
}

func (c *Client) SetDefaults() {
	if c.KubeConfig == "" {
		c.KubeConfig = "./k8sconfig/config.yml"
	}
}

func (c *Client) Init() {
	if c.clientset == nil {
		c.SetDefaults()

		c.clientset = connnect(&c.KubeConfig)
	}
}

func (c *Client) Client() *kubernetes.Clientset {
	return c.clientset
}

func connnect(kubeconfig *string) *kubernetes.Clientset {

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err)
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	// ping pong test
	_, err = clientset.ServerVersion()
	if err != nil {
		panic(err)
	}

	return clientset
}
