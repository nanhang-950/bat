package fn

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"time"
)

func GetProtocol(port int) string {
	protocols := map[int]string{
		21:    "FTP",
		22:    "SSH",
		23:    "Telnet",
		25:    "SMTP",
		53:    "DNS",
		67:    "DHCP",
		68:    "DHCP",
		80:    "HTTP",
		81:    "HTTP",
		110:   "POP3",
		119:   "NNTP",
		123:   "NTP",
		135:   "MS RPC",
		139:   "NetBIOS",
		143:   "IMAP",
		161:   "SNMP",
		162:   "SNMP",
		194:   "IRC",
		443:   "HTTPS",
		445:   "SMB",
		465:   "SMTPS",
		514:   "Syslog",
		873:   "rsync",
		993:   "IMAPS",
		995:   "POP3S",
		1080:  "SOCKS",
		1433:  "MSSQL",
		1434:  "MSSQL",
		1521:  "Oracle DB",
		1701:  "L2TP",
		1723:  "PPTP",
		1883:  "MQTT",
		3306:  "MySQL",
		3389:  "RDP",
		5100:  "IBM Tivoli",
		5421:  "Oracle DB",
		5432:  "PostgreSQL",
		5900:  "VNC",
		5984:  "CouchDB",
		6379:  "Redis",
		6667:  "IRC",
		7001:  "WebLogic",
		7002:  "WebLogic",
		7680:  "MSDTC",
		8000:  "HTTP",
		8080:  "HTTP",
		8001:  "HTTP",
		8082:  "HTTP",
		8089:  "HTTP",
		8443:  "HTTPS",
		9000:  "SonarQube",
		9092:  "Kafka",
		9100:  "JetDirect",
		9200:  "Elasticsearch",
		10000: "Webmin",
		11211: "Memcached",
		27017: "MongoDB",
		50000: "SAP",
	}
	if protocol, exists := protocols[port]; exists {
		return protocol
	}
	return "Unknown"
}

// 发送探针并获取响应
func sendProbe(ip string, port int, probe string) (string, error) {
	address := fmt.Sprintf("%s:%d", ip, port)
	conn, err := net.DialTimeout("tcp", address, 2*time.Second)
	if err != nil {
		return "", err
	}
	defer conn.Close()

	// 发送探针数据
	_, err = fmt.Fprintf(conn, probe)
	if err != nil {
		return "", err
	}

	// 读取响应
	conn.SetReadDeadline(time.Now().Add(2 * time.Second))
	reader := bufio.NewReader(conn)
	response, err := reader.ReadString('\n')
	if err != nil && err.Error() != "EOF" {
		return "", err
	}

	return strings.TrimSpace(response), nil
}

// 识别服务
func identifyService(ip string, port int) string {

	probes := map[string]string{
		"HTTP":       "GET / HTTP/1.1\r\nHost: example.com\r\n\r\n",
		"SMTP":       "HELO example.com\r\n",
		"FTP":        "USER anonymous\r\n",
		"POP3":       "USER anonymous\r\n",
		"SSH":        "SSH-2.0-TestClient\r\n",
		"Telnet":     "QUIT\r\n",
		"MySQL":      "SELECT VERSION();\r\n",
		"Redis":      "PING\r\n",
		"MongoDB":    "db.runCommand({ connectionStatus: 1 })\r\n",
		"IMAP":       "A001 LOGIN user pass\r\n",
		"RDP":        "RDP_REQ\r\n",                // 通常需要特定的协议处理
		"SNMP":       "GET /1.3.6.1.2.1.1.1.0\r\n", // 确保 SNMP 的请求格式正确
		"PostgreSQL": "SELECT VERSION();\r\n",
	}

	identifiers := map[string]func(string) bool{
		"HTTP":       func(r string) bool { return strings.HasPrefix(r, "HTTP/") },
		"SMTP":       func(r string) bool { return strings.HasPrefix(r, "220 ") },
		"FTP":        func(r string) bool { return strings.HasPrefix(r, "220 ") },
		"POP3":       func(r string) bool { return strings.HasPrefix(r, "+OK") },
		"SSH":        func(r string) bool { return strings.HasPrefix(r, "SSH-") },
		"Telnet":     func(r string) bool { return strings.HasPrefix(r, "TELNET") },
		"MySQL":      func(r string) bool { return strings.HasPrefix(r, "5.1.") },
		"Redis":      func(r string) bool { return strings.HasPrefix(r, "+PONG") },
		"MongoDB":    func(r string) bool { return strings.Contains(r, "MongoDB") },
		"IMAP":       func(r string) bool { return strings.HasPrefix(r, "* OK") },
		"RDP":        func(r string) bool { return strings.Contains(r, "RDP") },  // 简单匹配，可以更具体
		"SNMP":       func(r string) bool { return strings.Contains(r, "SNMP") }, // 确保 SNMP 请求的正确处理
		"PostgreSQL": func(r string) bool { return strings.Contains(r, "PostgreSQL") },
	}

	for service, probe := range probes {
		response, err := sendProbe(ip, port, probe)
		if err != nil {
			continue
		}

		if response != "" {
			if check, ok := identifiers[service]; ok && check(response) {
				return service
			}
		}
	}
	return ""
}
