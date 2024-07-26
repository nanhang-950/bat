package fn

import (
	"fmt"
	"net"
	"os"
)

func main() {
	ips := Getlocalip()
  fmt.Println(ips[0])
  fmt.Println(ips[1])
}

func Getlocalip() []string {
	//获取本地所有网络接口
	interfaces, err := net.Interfaces()
	//错误处理
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	//定义一个切片用于存储ip
	var localIps []string
	//迭代本地接口
	for _, iface := range interfaces {
		//使用Addrs获取每个网络接口的ip
		addrs, err := iface.Addrs()
		//错误处理
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}
		//迭代ip地址
		for _, addr := range addrs {
			//定义ip变量
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			//如果ip为空或ip为回环地址或ip不为ipv4地址则
			if ip == nil || ip.IsLoopback() || ip.To4() == nil {
				continue
			}
			if intranetip(ip) {
				localIps = append(localIps, addr.String())
			}
		}
	}

	return localIps
}

// 判断ip地址为内网ip
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

//生成网段内所有ip
func GenerateIPs(cidr string) ([]string,error){
  var ips []string
  ip,ipNet,err:=net.ParseCIDR(cidr)
  if err!=nil{
    return nil,err
  }
  for ip:=ip.Mask(ipNet.Mask);ipNet.Contains(ip);inc(ip){
    ips=append(ips,ip.String())
  }
  if len(ips)>2{
    return ips[1:len(ips)-1],nil
  }
  return ips,nil
}

//ip地址自增
func inc(ip net.IP){
  for j:=len(ip)-1;j>=0;j--{
    ip[j]++
    if ip[j]>0{
      break
    }
  }
}
