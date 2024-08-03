package fn 

import (
  "fmt"
  "time"
  "github.com/tatsushid/go-fastping"
  "net"
)

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
