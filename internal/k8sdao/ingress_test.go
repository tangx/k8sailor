package k8sdao

import (
	"fmt"
	"testing"
)

func Test_parseEndpoint(t *testing.T) {
	uRLs := []string{
		"http://www.baidu.com/v0/api",
		"http://www.baidu.com/v1/api*",
	}

	for _, uRL := range uRLs {
		fmt.Println(uRL)
	}
}
