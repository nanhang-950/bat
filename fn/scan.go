package fn

import (
	"fmt"
	"net"
  "sync"
  "time"
)

// 扫描给定的ip地址的端口

//参数：1.要扫描的ip地址，2.一个包含端口号的通道，3.一个包含扫描结果的通道，4.用于同步等待所有goroutine完成
func Scan(ip string, portsTask chan int, results chan ScanResult, wg *sync.WaitGroup) {
  //在函数退出时，调用wg.Done通知WaitGroup当前goroutine已完成
	defer wg.Done()

  //从通道中获取端口
	for {
		portTask, ok := <-portsTask
    //如果通道关闭则跳出循环
		if !ok {
			break
		}

    //构建地址字符串
		address := fmt.Sprintf("%s:%d", ip, portTask)

    //使用DialTimeout尝试连接目标地址，超时时间设置为6秒
		conn, err := net.DialTimeout("tcp", address, time.Second*6)

    //构建结构体
		result := ScanResult{
			IP:       ip,
			Port:     portTask,
			Protocol: GetProtocol(portTask),
      OS:       Getos(ip),
      Bb:       "1",
		}

    //如果连接成功，打印地址和状态并关闭连接
		if err == nil {
      state := "开放"
      fmt.Println(address,state)
			conn.Close()
		} else {
			continue
		}
    //将结果发送到results通道
		results <- result
	}
}

