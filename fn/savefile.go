package fn

import (
	"fmt"
	"html/template"
	"os"
)

func Savefile(results chan ScanResult){
	ipResults := make(map[string][]ScanResult)

	for result := range results {
		ipResults[result.IP] = append(ipResults[result.IP], result)
	}

	file, err := os.Create("result.html")
	if err != nil {
		fmt.Println("Failed to create file:", err)
		return
	}
	defer file.Close()

	tmpl := `
<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="UTF-8">
	<title>Scan Results</title>
	<style>
		table { width: 100%; border-collapse: collapse; }
		th, td { border: 1px solid black; padding: 8px; text-align: left; }
		th { background-color: #4CAF50; color: white; }
	</style>
</head>
<body>
	<h1>扫描报告</h1>
	{{range $ip, $results := .}}
		<h2>IP: {{$ip}}</h2>
		<table>
			<tr>
				<th>端口</th>
				<th>服务</th>
				<th>指纹</th>
			</tr>
			{{range $results}}
			<tr>
				<td>{{.Port}}</td>
				<td>{{.Protocol}}</td>
				<td>{{.State}}</td>
			</tr>
			{{end}}
		</table>
	{{end}}
</body>
</html>`

	t, err := template.New("webpage").Parse(tmpl)
	if err != nil {
		fmt.Println("Failed to parse template:", err)
		return
	}

	if err := t.Execute(file, ipResults); err != nil {
		fmt.Println("Failed to execute template:", err)
	}
}
