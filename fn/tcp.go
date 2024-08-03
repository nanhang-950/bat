package fn 

import (
  "net"
  "time"
  "fmt"
)

// TCP扫描
func TcpScan(ip string, ports []int) bool {
	for _, port := range ports {
		addr := fmt.Sprintf("%s:%d", ip, port)
		conn, err := net.DialTimeout("tcp", addr, time.Second)
		if err == nil {
			conn.Close()
			return true
		}
	}
	return false
}
