package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "gotool",
	Short: "gotool is set of various tools",
	Long:  "gotool is set of various tools",
}

func init() {
	rootCmd.AddCommand(
		NewSearchAliveAddrCommand(), // 搜索网段内，所有能ping通的机器
		NewServerSearchCommand(),    // 搜索网段内，所有开启22端口的机器（服务器）
		NewTcpScannerCommand(),      //扫描开启tcp监听的端口
	)
}

func Execute() error {
	return rootCmd.Execute()
}
