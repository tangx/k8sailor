package demo

import (
	"fmt"
	"net/url"
	"testing"
)

func Test_parseEndpoint(t *testing.T) {
	uRLs := []string{
		"http://www.baidu.com:80/v0/api?ns=default",
		"http://www.baidu.com/v1/api*?ns=default2",
		"https://querycap.feishu.cn/docs/doccnQO4cxBOZ67wqzC0EVlT3Kg#",
	}

	for _, uRL := range uRLs {
		// fmt.Println(uRL)
		t.Run(uRL, func(t *testing.T) {
			parseEndpoint(uRL)
		})

	}
}

func parseEndpoint(endpoint string) {
	u, err := url.Parse(endpoint)
	if err != nil {
		panic(err)
	}

	fmt.Println("u.Host=", u.Host)
	fmt.Println("u.RawPath=", u.RawPath)
	fmt.Println("u.RawQuery=", u.RawQuery)
	fmt.Println("u.Path=", u.Path)
	fmt.Println("u.Query=", u.Query())
}
