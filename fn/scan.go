package fn

import (
	"fmt"
	"net"
  "sync"
  "time"
)

// 扫描给定的ip地址的端口
func Scan(ip string, portsTask chan int, results chan ScanResult, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		portTask, ok := <-portsTask
		if !ok {
			break
		}

		address := fmt.Sprintf("%s:%d", ip, portTask)

		conn, err := net.DialTimeout("tcp", address, time.Second*6)
		result := ScanResult{
			IP:       ip,
			Port:     portTask,
			Protocol: GetProtocol(portTask),
		}
		if err == nil {
			result.State = "开放"
      fmt.Println(address,result.State)
			conn.Close()
		} else {
			continue
		}
		results <- result
	}
}

