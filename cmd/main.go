package main

import (
	"log"
	"os"

	"github.com/aiaoyang/processCpuUsage/pkg/metric"
	// _ "net/http/pprof"
)

func main() {
	log.SetFlags(log.Llongfile | log.Ldate)

	hostname, err := os.Hostname()
	if err != nil {
		log.Fatal(err)
	}

	env := "p08_dev"
	job(hostname, env)

}

func job(hostname, env string) {

	// 	ctx := context.TODO()

	// 	influxdbString := fmt.Sprintf("http://%s:%s",
	// 		configs.MyLocalConfig.InfluxDBConfig.Host,
	// 		configs.MyLocalConfig.InfluxDBConfig.Port,
	// 	)

	// 	c := influxc.NewClient(influxdbString, "")
	// 	defer c.Close()

	// 	ok, err := c.Health(ctx)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}

	// 	if ok.Status != "pass" {
	// 		log.Fatal("not connect to influxdb")
	// 	}

	// 	// TODO: influxdb测试数据库，后续添加正式名称
	// 	writeAPI := c.WriteAPI("test", "test")

	// 	go func() {

	// 		time.Sleep(time.Second * 10)

	// 		writeAPI.Flush()

	// 	}()

	// 	sender := sender.NewInfluxDBSender(writeAPI)

	// 	go processMon(sender)
	// 	go systemMon(sender)
	// 	go healthMon()

	// 	select {}
	// }

	// func processMon(sender metric.Sender) {

	// 	hostname, err := os.Hostname()
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}

	// 	env := configs.MyLocalConfig.Env

	// 	duration := time.Millisecond * 300

	// 	processJobMetric := metric.NewCustomMetric(sender)

	// 	processJobMetric.Tag.Insert(metric.HOSTNAME, hostname)
	// 	processJobMetric.Tag.Insert(metric.ENV, env)
	// 	for {

	// 		for _, pid := range configs.NeedMonitorProcessInfo.PIDS {

	// 			processCPU := sysusage.ProcessCPUUsageOnce(pid, duration)

	// 			processMEM := sysusage.ProcessMemUsageOnce(pid)

	// 			load1 := sysusage.Load()[sysusage.Load1]

	// 			fd := sysusage.OpenFD(pid)

	// 			processJobMetric.Tag.Insert(metric.PID, strconv.Itoa(pid))
	// 			processJobMetric.Metric.Insert(metric.PROCESS_CPU, processCPU)
	// 			processJobMetric.Metric.Insert(metric.PROCESS_MEM, processMEM)
	// 			processJobMetric.Metric.Insert(metric.FD, fd)
	// 			processJobMetric.Metric.Insert(metric.LOAD1, load1)

	// 			processJobMetric.Send()
	// 			time.Sleep(time.Millisecond * 100)

	// 		}

	// 		time.Sleep(time.Second * 5)

	// 	}
}

func systemMon(sender metric.Sender) {

	// hostname, err := os.Hostname()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// env := configs.MyLocalConfig.Env

	// duration := time.Millisecond * 300

	// sysJobMetric := metric.NewCustomMetric(sender)

	// sysJobMetric.Tag.Insert(metric.HOSTNAME, hostname)
	// sysJobMetric.Tag.Insert(metric.ENV, env)
	// sysJobMetric.Tag.Insert(metric.PID, strconv.Itoa(0))

	// for {

	// 	sysCPU := sysusage.SysCPUUsageOnce(duration)
	// 	sysMEM := sysusage.SystemMemUsageOnce()

	// 	sysJobMetric.Metric.Insert(metric.SYSCPU, sysCPU)
	// 	sysJobMetric.Metric.Insert(metric.SYSMEM, sysMEM)

	// 	sysJobMetric.Send()

	// 	time.Sleep(time.Second * 5)

	// }
}

func healthMon() {
	// for {

	// 	log.Printf("pid is : %v, name is : %s\n",
	// 		configs.NeedMonitorProcessInfo.PIDS,
	// 		configs.NeedMonitorProcessInfo.Names,
	// 	)

	// 	if _, isRunning := sysusage.IsPidRunning(configs.NeedMonitorProcessInfo.PIDS...); !isRunning {

	// 		tmppid := sysusage.GetProcessPID(configs.NeedMonitorProcessInfo.Names...)

	// 		if len(tmppid) != 0 {

	// 			configs.NeedMonitorProcessInfo.PIDS = tmppid

	// 			log.Printf("changed pid is : %v\n", configs.NeedMonitorProcessInfo.PIDS)

	// 		}
	// 	}

	// 	time.Sleep(time.Second * 5)
	// }
}
