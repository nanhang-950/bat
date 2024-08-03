package main

import (
	"fmt"
	"github.com/tatsushid/go-fastping"
	"net"
	"os"
	"sync"
	"time"
)

func main() {
	fmt.Println("扫描开始，请耐心等待")

	// 获取本地IP地址
	cidrs := Getlocalip()

	// 对每个子网进行扫描
	var allIPs []string
	for _, cidr := range cidrs {
		ips, err := GenerateIPs(cidr)
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}
		allIPs = append(allIPs, ips...)
	}

	// 定义常见端口列表
	var commonPorts = []int{80, 443, 22, 21, 3389,25,23,137,138,139,3389}

	// 并发扫描
	var wg sync.WaitGroup
	results := make(chan string, len(allIPs))

	start := time.Now()
	for _, ip := range allIPs {
		wg.Add(1)
		go func(ip string) {
			defer wg.Done()
			if IcmpScan(ip) || TcpScan(ip, commonPorts) {
				results <- ip
			}
		}(ip)
	}

	wg.Wait()
	close(results)

	// 输出扫描结果信息
	fmt.Println("存活主机：")
	for ip := range results {
		fmt.Println(ip)
	}
	fmt.Println("用时：", time.Since(start).String())
}

// 获取本地网卡IP地址
func Getlocalip() []string {
	interfaces, err := net.Interfaces()
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	var localIps []string
	for _, iface := range interfaces {
		addrs, err := iface.Addrs()
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}

		for _, addr := range addrs {
			var ip net.IP
			var mask net.IPMask
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
				mask = v.Mask
			case *net.IPAddr:
				ip = v.IP
				mask = ip.DefaultMask()
			}

			if ip == nil || ip.IsLoopback() || ip.To4() == nil {
				continue
			}

			if intranetip(ip) {
				cidr := fmt.Sprintf("%s/%d", ip.String(), maskSize(mask))
				localIps = append(localIps, cidr)
			}
		}
	}
	return localIps
}

// 获取掩码长度
func maskSize(mask net.IPMask) int {
	size, _ := mask.Size()
	return size
}

// 判断IP地址是否为内网地址
func intranetip(ip net.IP) bool {
	blocks := []string{
		"10.0.0.0/8",
		"172.16.0.0/12",
		"192.168.0.0/16",
	}

	for _, cidr := range blocks {
		_, block, _ := net.ParseCIDR(cidr)
		if block.Contains(ip) {
			return true
		}
	}
	return false
}

// 对IP地址进行递增
func IPInc(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}

// 生成CIDR范围内的所有IP地址
func GenerateIPs(cidr string) ([]string, error) {
	ip, ipNet, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, err
	}

	var ips []string
	for ip := ip.Mask(ipNet.Mask); ipNet.Contains(ip); IPInc(ip) {
		if !ip.Equal(ipNet.IP) && !ip.Equal(lastIP(ipNet)) {
			ips = append(ips, ip.String())
		}
	}
	return ips, nil
}

// 获取CIDR的最后一个IP地址
func lastIP(ipNet *net.IPNet) net.IP {
	ip := ipNet.IP
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j] |= ^ipNet.Mask[j]
	}
	return ip
}

// ICMP扫描
func IcmpScan(ip string) bool {
	p := fastping.NewPinger()
	ra, err := net.ResolveIPAddr("ip4:icmp", ip)
	if err != nil {
		fmt.Println("Error resolving IP:", err)
		return false
	}
	p.AddIPAddr(ra)

	found := false
	p.OnRecv = func(addr *net.IPAddr, rtt time.Duration) {
		found = true
	}
	p.OnIdle = func() {}

	p.MaxRTT = time.Second
	err = p.Run()
	if err != nil {
		fmt.Println("Error running ping:", err)
		return false
	}

	return found
}

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
