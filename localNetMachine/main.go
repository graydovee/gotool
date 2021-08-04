//搜索网段内，所有能ping通的机器

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

func getLocalNetMachines(ip string, timeout int) []string {

	//collect ip
	ch := make(chan string)
	ipList := make([]string, 1)

	if timeout < 1 {
		timeout = 1
	}
	go testIp(ip, ch, time.Duration(timeout) * time.Millisecond)
	for s := range ch {
		ipList = append(ipList, s)
	}
	return ipList
}

func testIp(ip string, ch chan string, timeout time.Duration) {
	//test ip
	var wg sync.WaitGroup
	var maxNum = 255
	wg.Add(maxNum)
	ipFormatter := strings.Replace(ip, "*", "%d", 1)
	for i := 1; i <= maxNum; i++ {
		go func(number int) {
			defer wg.Done()
			ip := fmt.Sprintf(ipFormatter, number)
			ping := netutils.Ping(ip,  timeout)
			if ping {
				ch <- ip
			}
		}(i)
	}
	wg.Wait()
	close(ch)
}

func main() {
	ip := flag.String("ip", "192.168.1.*", "test target ip")
	timeout := flag.Int("t", 1000, "ping timeout")
	flag.Parse()
	if strings.Count(*ip, "*") != 1 {
		fmt.Println("ip is illegal")
	}
	machines := getLocalNetMachines(*ip, *timeout)
	sort.Strings(machines)
	for _, ip := range machines {
		fmt.Println(ip)
	}
}
