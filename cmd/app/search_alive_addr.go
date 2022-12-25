package app

import (
	"fmt"
	"github.com/grydovee/gotool/netutils"
	"github.com/spf13/cobra"
	"sort"
	"strings"
	"sync"
	"time"
)

var (
	saaIp      string
	saaTimeout int
)

func getLocalNetMachines(ip string, timeout int) []string {

	//collect ip
	ch := make(chan string)
	ipList := make([]string, 1)

	if timeout < 1 {
		timeout = 1
	}
	go testIpCanPing(ip, ch, time.Duration(timeout)*time.Millisecond)
	for s := range ch {
		ipList = append(ipList, s)
	}
	return ipList
}

func testIpCanPing(ip string, ch chan string, timeout time.Duration) {
	//test ip
	var wg sync.WaitGroup
	var maxNum = 255
	wg.Add(maxNum)
	ipFormatter := strings.Replace(ip, "*", "%d", 1)
	for i := 1; i <= maxNum; i++ {
		go func(number int) {
			defer wg.Done()
			ip := fmt.Sprintf(ipFormatter, number)
			ping := netutils.Ping(ip, timeout)
			if ping.Err != nil {
				ch <- ip
			}
		}(i)
	}
	wg.Wait()
	close(ch)
}

func NewSearchAliveAddrCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "saa",
		Short: "search alive address in target ip",
		Run: func(cmd *cobra.Command, args []string) {
			if strings.Count(saaIp, "*") != 1 {
				fmt.Println("ip is illegal")
			}
			machines := getLocalNetMachines(saaIp, saaTimeout)
			sort.Strings(machines)
			for _, ip := range machines {
				fmt.Println(ip)
			}
		},
	}
	cmd.Flags().StringVar(&saaIp, "ip", "192.168.1.*", "test target ip")
	cmd.Flags().IntVarP(&saaTimeout, "timeout", "t", 1000, "timeout, unit: ms")
	return cmd
}
