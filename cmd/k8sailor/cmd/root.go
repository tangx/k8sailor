package cmd

import (
	"github.com/go-jarvis/cobrautils"
	"github.com/spf13/cobra"
	"github.com/tangx/k8sailor/cmd/k8sailor/global"
)

var rootCmd = &cobra.Command{
	Use:  "k8sailor",
	Long: "k8s 管理平台",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// 始终显示帮助信息
		_ = cmd.Help()
	},
	Run: func(cmd *cobra.Command, args []string) {
		// 什么也不做
		// k8s.Connent()
	},
}

func init() {
	cobrautils.BindFlags(rootCmd, global.Flags)

	rootCmd.AddCommand(cmdHttpserver)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}
