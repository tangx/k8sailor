package global

type CmdFlags struct {
	Config string `flag:"config" usage:"k8s 配置授权文件" persistent:"true"`
}

var Flags = &CmdFlags{
	Config: "./k8sconfig/config.yml",
}
