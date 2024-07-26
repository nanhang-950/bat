package main

import (
  "fmt"
  "bat/fn"
  "sync"
  "time"
)
func main(){
  fn.Banner()
  fmt.Println("扫描开始，请耐心等待")
  cidrs:=make([]string,10)
  cidrs=fn.Getlocalip()
  var allIPs []string
  for _,cidr:=range cidrs{
    ips,err:=fn.GenerateIPs(cidr)
    if err!=nil{
      fmt.Println("Error:",err)
      continue
    }
    allIPs=append(allIPs,ips...)
  }
  var wg sync.WaitGroup
  results:=make(chan fn.ScanResult,len(allIPs)*len(fn.DefaultPorts))
  start:=time.Now()
  for _,ip:=range allIPs{
    portsTask:=make(chan int,len(fn.DefaultPorts))
    for _,port:=range fn.DefaultPorts{
      portsTask<-port
    }
    close(portsTask)
    for threads:=0;threads<10;threads++{
      wg.Add(1)
      go fn.Scan(ip,portsTask,results,&wg)
    }
  }
  wg.Wait()
  close(results)
  fn.Savefile(results)
  fmt.Println("扫描报告已生成：result.html")
  fmt.Println("用时：",time.Since(start).String())
}
