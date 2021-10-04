package k8sdemo

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/tangx/k8sailor/pkg/confk8s"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func Test_GetEvent(t *testing.T) {

	kubeClient := confk8s.Client{
		KubeConfig: "../../cmd/k8sailor/k8sconfig/config.yml",
	}
	kubeClient.Init()
	clientset := kubeClient.Client()

	ctx := context.TODO()
	events, err := clientset.EventsV1().Events("default").List(ctx, v1.ListOptions{})
	if err != nil {
		panic(err)
	}

	for _, event := range events.Items {
		// fmt.Printf("%+v", event)
		b, _ := json.MarshalIndent(event, "", "  ")
		fmt.Printf("%s", b)

		fmt.Println("##########")
	}

	fmt.Println("##########")

	fmt.Println("-==========-")
	events2, err := clientset.CoreV1().Events("default").List(ctx, v1.ListOptions{})
	if err != nil {
		panic(err)
	}

	for _, event := range events2.Items {
		// fmt.Printf("%+v", event)
		b, _ := json.MarshalIndent(event, "", "  ")

		fmt.Printf("%s\n", b)
		fmt.Println("=======")
	}

}
