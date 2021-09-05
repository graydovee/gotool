package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"gotool/netutils"
	"sort"
	"time"
)

var (
	tcpScannerHost    string
	tcpScannerPort    int
	tcpScannerNum     int
	tcpScannerTimeout int
)

func getAddress(port int) string {
	return fmt.Sprintf("%s:%d", tcpScannerHost, port)
}

func work(port int, results chan int) {
	address := getAddress(port)
	if netutils.Tcping(address, time.Duration(tcpScannerTimeout) * time.Millisecond) {
		results <- port
		return
	}
	results <- -port
}

func NewTcpScannerCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use: "tcpscan",
		Short: "scan opened tcp port on target ip",
		Long: "scan opened tcp port on target ip",
		Run: func(cmd *cobra.Command, args []string) {
			var portMin = tcpScannerPort
			var portMax = tcpScannerPort + tcpScannerNum
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
		},
	}
	cmd.Flags().StringVarP(&tcpScannerHost, "scanner_addr", "a", "127.0.0.1", "需要扫描的host")
	cmd.Flags().IntVarP(&tcpScannerPort, "port", "p", 1, "需要扫描起始端口")
	cmd.Flags().IntVarP(&tcpScannerNum, "scan_num", "n", 120, "需要扫描端口的数量")
	cmd.Flags().IntVarP(&tcpScannerTimeout, "timeout", "t", 1000, "tcp扫描超时时间")
	return cmd
}
