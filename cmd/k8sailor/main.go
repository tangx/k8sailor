package main

import (
	"github.com/sirupsen/logrus"
	"github.com/tangx/k8sailor/cmd/k8sailor/cmd"
)

func main() {
	cmd.Execute()
}

func init() {
	logrus.SetLevel(logrus.DebugLevel)
}
