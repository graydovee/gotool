package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"gotool/netutils"
	"sort"
	"strings"
	"sync"
	"time"
)

var (
	ssIp      string
	ssTimeout int
)

func testIpOenSsh(ip string, ch chan string, timeout time.Duration) {
	var wg sync.WaitGroup
	var maxNum = 255
	wg.Add(maxNum)
	ipFormatter := strings.Replace(ip, "*", "%d", 1)
	for i := 1; i <= maxNum; i++ {
		go func(number int) {
			defer wg.Done()
			ip := fmt.Sprintf(ipFormatter, number)
			ping := netutils.Ping(ip, timeout)

			if ping {
				addr := ip + ":22"
				tcping := netutils.Tcping(addr, timeout)
				if tcping {
					ch <- ip
				}
			}
		}(i)
	}
	wg.Wait()
	close(ch)
}

func searchServer(ip string, timeout time.Duration) []string {
	ch := make(chan string)
	go testIpOenSsh(ip, ch, timeout)
	arr := make([]string, 0)
	for s := range ch {
		arr = append(arr, s)
	}
	return arr
}

func NewServerSearchCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ss",
		Short: "search server which open port 22 in target ip segment",
		Long:  "search server which open port 22 in target ip segment",
		Run: func(cmd *cobra.Command, args []string) {
			if strings.Count(ssIp, "*") != 1 {
				fmt.Println("ip is illegal")
			}
			machines := searchServer(ssIp, time.Duration(ssTimeout)*time.Millisecond)
			sort.Strings(machines)
			for _, ip := range machines {
				fmt.Println(ip)
			}
		},
	}
	cmd.Flags().StringVar(&ssIp, "ip", "192.168.1.*", "test target ip")
	cmd.Flags().IntVarP(&ssTimeout, "timeout", "t", 1000, "timeout, unit: ms")
	return cmd
}
