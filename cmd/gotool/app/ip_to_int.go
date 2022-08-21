package app

import (
	"fmt"
	"github.com/spf13/cobra"
	"gotool/netutils"
)

func NewIpToIntCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "ip2int [sub]",
		Short: "Print the version number",
		Long:  `All software has versions. This is mine`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) > 0 {
				ip, err := netutils.IpTo32(args[0])
				if err != nil {
					fmt.Println("format error")
				} else {
					fmt.Println(ip)
				}
			}
		},
	}
}
