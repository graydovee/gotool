//搜索网段内，所有开启22端口的机器（服务器）

package main

import (
	"flag"
	"fmt"
	"gotool/core/netutils"
	"sort"
	"strings"
	"sync"
	"time"
)

func testIp(ip string, ch chan string, timeout time.Duration) {
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
				if (tcping) {
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
	go testIp(ip, ch, timeout)
	arr := make([]string, 0)
	for s := range ch {
		arr = append(arr, s)
	}
	return arr
}

func main() {
	ip := flag.String("ip", "192.168.10.*", "test target ip")
	timeout := flag.Int("t", 1000, "ping timeout")
	flag.Parse()
	if strings.Count(*ip, "*") != 1 {
		fmt.Println("ip is illegal")
	}
	machines := searchServer(*ip, time.Duration(*timeout) * time.Millisecond)
	sort.Strings(machines)
	for _, ip := range machines {
		fmt.Println(ip)
	}
}