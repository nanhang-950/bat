package fn

import (
	"fmt"
	"html/template"
	"os"
)

//将扫描结果格式化为HTML并写入到result.html文件中
//参数revults是一个chan ScanResult类型的通道
func Savefile(results chan ScanResult){

  //定义一个映射类型的变量，用于存储每个ip地址对应的扫描结果列表
	ipResults := make(map[string][]ScanResult)

  //遍历results通道，将扫描结果按照ip地址分类汇总到ipResults映射
	for result := range results {
		ipResults[result.IP] = append(ipResults[result.IP], result)
	}
  for os:= range results{
    ipResults[]=
  }
  //创建html文件
	file, err := os.Create("result.html")
  //错误处理
	if err != nil {
		fmt.Println("Failed to create file:", err)
		return
	}
  //确保文件在函数结束时关闭
	defer file.Close()

  //定义一个html模板，用于格式化扫描结果
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
		<h2>OS: {{index $results 0. OS}}</h2>
		<table>
			<tr>
				<th>Port</th>
				<th>Protocol</th>
				<th>OS</th>
			</tr>
			{{range $results}}
			<tr>
				<td>{{.Port}}</td>
				<td>{{.Protocol}}</td>
				<td>{{.Bb}}</td>
			</tr>
			{{end}}
		</table>
	{{end}}
</body>
</html>`
  
  //使用template.New创建一个新的模板对象了，并解析定义的HTML模板tmpl。
	t, err := template.New("webpage").Parse(tmpl)

  //错误处理
	if err != nil {
		fmt.Println("Failed to parse template:", err)
		return
	}

  //使用Execute将ipResults数据填充到模板中，并写入到文件中。
  //如果执行失败，则输出错误信息。
	if err := t.Execute(file, ipResults); err != nil {
		fmt.Println("Failed to execute template:", err)
	}
}
