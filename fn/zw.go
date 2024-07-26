package fn

func GetProtocol(port int) string {
	protocols := map[int]string{
		21:    "FTP",
		22:    "SSH",
		23:    "Telnet",
		25:    "SMTP",
		80:    "HTTP",
		81:    "HTTP",
		135:   "MS RPC",
		139:   "NetBIOS",
		443:   "HTTPS",
		445:   "SMB",
		1433:  "MSSQL",
		1521:  "Oracle DB",
		3306:  "MySQL",
		5100:  "IBM Tivoli",
		5421:  "Oracle DB",
		5432:  "PostgreSQL",
		6379:  "Redis",
		7001:  "WebLogic",
		7002:  "WebLogic",
		8000:  "HTTP",
		8080:  "HTTP",
		8001:  "HTTP",
		8082:  "HTTP",
		8089:  "HTTP",
		9000:  "HTTP",
		9100:  "JetDirect",
		9200:  "Elasticsearch",
		11211: "Memcached",
		27017: "MongoDB",
	}
  if protocol,exists:=protocols[port];exists{
    return protocol;
  }
  return "Unknown"
}
