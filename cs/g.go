package main

import (
	"fmt"

	"net"
)

// 生成指定网段的所有 IP 地址

func GenerateIPs(cidr string) ([]string, error) {

	var ips []string

	ip, ipNet, err := net.ParseCIDR(cidr)

	if err != nil {
		return nil, err
	}

	for ip := ip.Mask(ipNet.Mask); ipNet.Contains(ip); IPInc(ip) {
		if !ip.Equal(ipNet.IP) && !ip.Equal(lastIP(ipNet)) {
			ips = append(ips, ip.String())
		}
	}
	return ips, nil
}

// IP 地址自增
func IPInc(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}

// 计算网段的最后一个 IP 地址
func lastIP(ipNet *net.IPNet) net.IP {
	ip := ipNet.IP
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j] |= ^ipNet.Mask[j]
	}
	return ip
}

func main() {
	// 内网 IP 地址块
	cidrs := []string{
		"10.0.0.0/8",
		"172.16.0.0/12",
		"192.168.0.0/16",
	}

	for _, cidr := range cidrs {

		ips, err := GenerateIPs(cidr)

		if err != nil {

			fmt.Println("Error generating IPs:", err)

			return

		}

		fmt.Printf("IPs in %s:\n", cidr)

		for _, ip := range ips {

			fmt.Println(ip)

		}

	}

}
