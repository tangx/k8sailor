package demo

import (
	"fmt"
	"net/url"
	"strings"
	"testing"
)

func Test_parseEndpoint(t *testing.T) {
	uRLs := []string{
		"http://www.baidu.com/api?backend=my-nginx-service:8080&tls=secret1",
		"http://www.baidu.com:80/v0/api?ns=default",
		"http://www.baidu.com/v1/api*?ns=default2",
		"http://:8080/download",
		"http:///download",
		"http://www.baidu.com/*",
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
	fmt.Println("u.Query=", u.Query().Get("tls"))
}

func Test_Split(t *testing.T) {
	for _, s := range []string{"", ":", "my-nginx-3", ":8080"} {
		t.Run(s, func(t *testing.T) {
			parts := strings.Split(s, ":")
			fmt.Println("s", s, "=>", len(parts))
			fmt.Println("parts=", parts)
		})
	}
}
