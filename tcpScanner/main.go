//扫描开启tcp监听的端口

package main

import (
	"flag"
	"fmt"
	"gotool/core/netutils"
	"sort"
	"time"
)

var host string
var port int
var num int

func init() {
	flag.StringVar(&host, "h", "127.0.0.1", "需要扫描的host")
	flag.IntVar(&port, "p", 1, "需要扫描起始端口")
	flag.IntVar(&num, "n", 120, "需要扫描端口的数量")
	flag.Parse()
}

func getAddress(port int) string {
	return fmt.Sprintf("%s:%d", host, port)
}

func work(port int, results chan int) {
	address := getAddress(port)
	if netutils.Tcping(address, time.Second * 1) {
		results <- port
		return
	}
	results <- -port
}


func main() {

	var portMin = port
	var portMax = port + num
	openPorts := make(chan int)

	go func() {
		for i := portMin; i < portMax; i++ {
			go work(i, openPorts)
		}
	}()

	openedPortsArr := make([]int, 0, portMax-portMin)
	closedPortsArr := make([]int, 0, portMax-portMin)

	for i := portMin; i < portMax; i++ {
		port := <-openPorts
		if port <= 0 {
			closedPortsArr = append(closedPortsArr, -port)
		} else {
			openedPortsArr = append(openedPortsArr, port)
		}
	}
	close(openPorts)

	sort.Ints(closedPortsArr)
	sort.Ints(openedPortsArr)
	for _, v := range closedPortsArr {
		fmt.Printf("port %d is closed\n", v)
	}
	for _, v := range openedPortsArr {
		fmt.Printf("port %d is opened\n", v)
	}
}
