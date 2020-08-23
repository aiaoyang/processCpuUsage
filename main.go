package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/aiaoyang/processCpuUsage/metric"
	"github.com/aiaoyang/processCpuUsage/sender"
	"github.com/aiaoyang/processCpuUsage/sysusage"

	// _ "net/http/pprof"

	influxc "github.com/influxdata/influxdb-client-go"
)

func init() {
	initViper()
}

func main() {
	todo := context.TODO()

	go MyLocalConfig.watchConfigChange(todo)

	hostname, err := os.Hostname()
	if err != nil {
		log.Fatal(err)
	}

	env := "p08_dev"
	job(hostname, env)

}

func job(hostname, env string) {

	ctx := context.TODO()

	c := influxc.NewClient("http://"+MyLocalConfig.InfluxDBConfig.Host+":"+MyLocalConfig.InfluxDBConfig.Port, "")

	ok, err := c.Health(ctx)
	if err != nil {
		// 如果无法连接到 influxdb 则退出
		log.Fatal(err)
	}

	// 检查 influxdb 连接是否可用
	if ok.Status != "pass" {
		log.Fatal("not connect to influxdb")
	}

	// TODO: influxdb测试数据库，后续添加正式名称
	writeAPI := c.WriteAPI("test", "test")

	go func() {
		log.Println("hello world")

		time.Sleep(time.Second * 10)

		writeAPI.Flush()

	}()
	defer c.Close()
	// defer writeAPI.Close()

	sender := sender.NewInfluxDBSender(writeAPI)

	go func() {
		processJobMetric := metric.NewCustomMetric(sender)

		duration := time.Millisecond * 100

		pids := []int{1}
		processJobMetric.Tag.Insert(metric.HOSTNAME, hostname)
		processJobMetric.Tag.Insert(metric.ENV, env)
		for {

			for _, pid := range pids {

				processCPU := sysusage.ProcessCPUUsageOnce(pid, duration)
				processMEM := sysusage.ProcessMemUsageOnce(pid)
				fd := sysusage.OpenFD(pid)

				processJobMetric.Metric.Insert(metric.PID, metric.MetricType(strconv.Itoa(pid)))
				processJobMetric.Metric.Insert(metric.PROCESS_CPU, processCPU)
				processJobMetric.Metric.Insert(metric.PROCESS_MEM, processMEM)
				processJobMetric.Metric.Insert(metric.FD, fd)

			}

			fmt.Printf("process metric : %v\n", processJobMetric)
			processJobMetric.Send()
			time.Sleep(time.Second * 5)

		}
	}()

	go func() {

		sysJobMetric := metric.NewCustomMetric(sender)

		duration := time.Millisecond * 100

		sysJobMetric.Tag.Insert(metric.HOSTNAME, hostname)
		sysJobMetric.Tag.Insert(metric.ENV, env)

		for {

			fmt.Println("here")
			sysCPU := sysusage.SysCPUUsageOnce(duration)
			sysMEM := sysusage.SystemMemUsageOnce()

			sysJobMetric.Metric.Insert(metric.SYSCPU, sysCPU)
			sysJobMetric.Metric.Insert(metric.SYSMEM, sysMEM)

			fmt.Printf("sys metric : %v\n", sysJobMetric)

			sysJobMetric.Send()

			time.Sleep(time.Second * 10)
		}
	}()
	select {}
}

/*























































 */

// func genericTODO() {

// 	ctx := context.TODO()

// 	c := influxc.NewClient("http://"+MyLocalConfig.InfluxDBConfig.Host+":"+MyLocalConfig.InfluxDBConfig.Port, "")
// 	ok, err := c.Health(ctx)
// 	if err != nil {
// 		// 如果无法连接到 influxdb 则退出
// 		log.Fatal(err)
// 	}

// 	// 检查 influxdb 连接是否可用
// 	if ok.Status != "pass" {
// 		log.Fatal("not connect to influxdb")
// 	}

// 	// TODO: influxdb测试数据库，后续添加正式名称
// 	writeAPI := c.WriteAPI("test", "test")

// 	defer c.Close()
// 	defer writeAPI.Close()

// 	reciver := make(chan metric.Metric, 0)

// 	pidsChan := make(chan []int, 0)

// 	// 告警监控

// 	go sysusage.ProcessCPUUsageOnce(0, time.Second)

// 	for {
// 		select {
// 		case <-ctx.Done():
// 			return

// 		// 接收到指标参数，则将其推送至influxdb
// 		case res := <-reciver:
// 			tag := make(map[string]string)

// 			// 如果有pid这个key，则将其写入tag，否则则是只有系统cpu使用量
// 			pid, ok := res["pid"].(float64)

// 			if ok {

// 				tag["pid"] = strconv.Itoa(int(pid))

// 			}
// 			// else {

// 			// alarmer.Alarm("告警发生: 进程未运行")

// 			// }

// 			// 删除不需要res中的pid键
// 			delete(res, "pid")

// 			// debug用
// 			fmt.Printf("recive value : %v\n", res)

// 			// 推送数据至influxdb
// 			pushToInfluxDB(writeAPI, tag, res)

// 		default:
// 			if AliyunConfigSrv.GetInt("processinfo.status") == 1 {

// 				// 如果进程仍在运行，则将 pids 发送给 通道 然后让指标收集函数进行处理
// 				if pids, hasDeadPid := sysusage.IsPidRunning(NeedMonitorProcessInfo.PIDS...); !hasDeadPid {

// 					pidsChan <- pids

// 					// debug 打印
// 					fmt.Printf("running %v\n", pids)
// 				} else { // 如果进程被中断，则循环监听进程列表，直至同名进程启动

// 					// 负值pid设定为不进行进程监控
// 					pidsChan <- []int{-1}

// 					// 查询进程pid是否存在
// 					NeedMonitorProcessInfo.PIDS = sysusage.GetProcessPID(NeedMonitorProcessInfo.Name)

// 					fmt.Printf("not running %v\n", pids)

// 				}
// 			} else {

// 				fmt.Println("nothing happend")

// 			}

// 			time.Sleep(time.Second * 3)

// 		}
// 	}
// }
