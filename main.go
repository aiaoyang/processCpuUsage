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
	ctx := context.TODO()

	MyLocalConfig.watchConfigChange(ctx)

	hostname, err := os.Hostname()
	if err != nil {
		log.Fatal(err)
	}

	env := "p08_dev"
	job(hostname, env)

}

func job(hostname, env string) {

	ctx := context.TODO()

	influxdbString := fmt.Sprintf("http://%s:%s",
		MyLocalConfig.InfluxDBConfig.Host,
		MyLocalConfig.InfluxDBConfig.Port,
	)

	c := influxc.NewClient(influxdbString, "")
	defer c.Close()

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

		time.Sleep(time.Second * 10)

		writeAPI.Flush()

	}()

	sender := sender.NewInfluxDBSender(writeAPI)

	go func() {

		duration := time.Millisecond * 300

		processJobMetric := metric.NewCustomMetric(sender)

		processJobMetric.Tag.Insert(metric.HOSTNAME, hostname)
		processJobMetric.Tag.Insert(metric.ENV, env)
		for {

			for _, pid := range NeedMonitorProcessInfo.PIDS {

				processCPU := sysusage.ProcessCPUUsageOnce(pid, duration)

				processMEM := sysusage.ProcessMemUsageOnce(pid)
				fd := sysusage.OpenFD(pid)

				processJobMetric.Tag.Insert(metric.PID, strconv.Itoa(pid))
				processJobMetric.Metric.Insert(metric.PROCESS_CPU, processCPU)
				processJobMetric.Metric.Insert(metric.PROCESS_MEM, processMEM)
				processJobMetric.Metric.Insert(metric.FD, fd)

				processJobMetric.Send()

			}

			time.Sleep(time.Second * 5)

		}
	}()

	go func() {

		duration := time.Millisecond * 300

		sysJobMetric := metric.NewCustomMetric(sender)

		sysJobMetric.Tag.Insert(metric.HOSTNAME, hostname)
		sysJobMetric.Tag.Insert(metric.ENV, env)
		sysJobMetric.Tag.Insert(metric.PID, strconv.Itoa(0))

		for {

			sysCPU := sysusage.SysCPUUsageOnce(duration)
			sysMEM := sysusage.SystemMemUsageOnce()

			sysJobMetric.Metric.Insert(metric.SYSCPU, sysCPU)
			sysJobMetric.Metric.Insert(metric.SYSMEM, sysMEM)

			sysJobMetric.Send()

			time.Sleep(time.Second * 10)
		}
	}()
	select {}
}
