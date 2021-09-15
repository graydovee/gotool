package cmd

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"strconv"
	"strings"
)

func ipTo32(ipStr string) (uint32, error) {
	split := strings.Split(ipStr, ".")
	if len(split) != 4 {
		return 0, errors.New("err")
	}
	var ip uint32 = 0
	for i := 0; i < 4; i++ {
		n, err := strconv.Atoi(split[3-i])
		if err != nil {
			return 0, err
		}
		ip += uint32(n) << (i * 8)
	}
	return ip, nil
}

func NewIpToIntCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "ip2int [sub]",
		Short: "Print the version number",
		Long:  `All software has versions. This is mine`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) > 0 {
				ip, err := ipTo32(args[0])
				if err != nil {
					fmt.Println("format error")
				} else {
					fmt.Println(ip)
				}
			}
		},
	}
}
