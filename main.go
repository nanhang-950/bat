package main

import (
	"fmt"
	"bat/fn"
	"sync"
	"time"
)

func main() {
	fn.Banner()
	fmt.Println("扫描开始，请耐心等待")

	// 获取网段ip
	cidrs := fn.Getlocalip()

	// 定义一个用于存储ip的切片
	var allIPs []string

	// 遍历每个cidr地址段
	for _, cidr := range cidrs {
		// 使用 GenerateIPs 生成对应的 ip 地址列表
		ips, err := fn.GenerateIPs(cidr)
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}
		// 将生成的 ip 地址追加到 allIPs 切片中
		allIPs = append(allIPs, ips...)
	}

	// 常见端口列表
	var commonPorts = []int{80, 443, 22, 21, 3389, 25, 23, 137, 138, 139, 3389}

	// 并发扫描
	var wg sync.WaitGroup
	results1 := make(chan string, len(allIPs))

	// 控制最大并发数
	const maxConcurrency = 1000
	semaphore := make(chan struct{}, maxConcurrency)

	start := time.Now()
	fmt.Println("存活主机")

	for _, ip := range allIPs {
		wg.Add(1)
		go func(ip string) {
			defer wg.Done()
			semaphore <- struct{}{} // 请求一个槽位
			defer func() { <-semaphore }() // 释放一个槽位
			if fn.IcmpScan(ip) || fn.TcpScan(ip, commonPorts) {
				results1 <- ip
				fmt.Println(ip)
			}
		}(ip)
	}

	wg.Wait()
	close(results1)

	var aliveIPs []string
	for ip := range results1 {
		aliveIPs = append(aliveIPs, ip)
	}

	results2 := make(chan fn.ScanResult, len(aliveIPs)*len(fn.DefaultPorts))

	for _, ip := range aliveIPs {
		portsTask := make(chan int, len(fn.DefaultPorts))
		for _, port := range fn.DefaultPorts {
			portsTask <- port
		}
		close(portsTask)

		for threads := 0; threads < 600; threads++ {
			wg.Add(1)
			go fn.Scan(ip, portsTask, results2, &wg)
		}
	}

	wg.Wait()
	close(results2)

	// 生成报告文件
	fn.Savefile(results2)

	// 输出扫描结果信息
	fmt.Println("扫描报告已生成：result.html")
	fmt.Println("用时：", time.Since(start).String())
}
