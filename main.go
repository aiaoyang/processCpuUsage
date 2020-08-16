package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"strconv"

	// _ "net/http/pprof"
	"strings"
	"sync"
	"time"

	"github.com/aiaoyang/p08_monitor/common"
	influxc "github.com/influxdata/influxdb-client-go"
)

// var conf = &ALiConfigClientConfig{}

func init() {
	initViper()
	flag.StringVar(&pids, "p", "1", "pids of process")
	flag.Parse()
	// conf.new()
}

var pids string

func main() {
	p, err := strconv.Atoi(pids)
	if err != nil {
		log.Fatal(err)
	}
	reciver := make(chan map[string]float64, 0)
	// ctx, cancel := context.WithCancel(context.Background())
	ctx := context.TODO()
	go func() {
		for {
			select {
			case m := <-reciver:
				fmt.Printf("recive: %v\n", m)
			default:
				time.Sleep(time.Second)
			}
		}
	}()
	pidsChan := make(chan []int32, 0)

	go func(pids int) {
		for {

			pidsChan <- []int32{int32(pids)}
			time.Sleep(time.Second * 3)
		}
	}(p)
	// go func() { time.Sleep(time.Second * 5); cancel() }()
	common.ProcessesCPUMonitor(ctx, reciver, pidsChan)

	// go func() {
	// 	time.Sleep(time.Second * 5)
	// 	pids = "1"
	// }()

	// log.SetFlags(log.Llongfile | log.Ltime)
	// todo := context.TODO()
	// go CLocal.watchConfigChange(todo)
	// // go http.ListenAndServe("0.0.0.0:10088", nil)
	// genericTODO(nil)
}

func genericTODO(alarm func(msg interface{})) {
	ctx := context.TODO()
	c := influxc.NewClient("http://"+CLocal.InfluxDBConfig.Host+":"+CLocal.InfluxDBConfig.Port, "")
	ok, err := c.Health(ctx)
	if err != nil {
		log.Fatal(err)
	}
	if ok.Status != "pass" {
		log.Fatal("not connect to influxdb")
	}
	writeAPI := c.WriteAPI("test", "test")
	defer c.Close()
	defer writeAPI.Close()

	isAlarm := false
	alarmOnce := sync.Once{}
	cancelAlarmOnce := sync.Once{}
	for {
		if VCloud.GetInt("processinfo.status") == 1 {
			// fmt.Printf("pids: %s\n", pids)
			if pids, hasDeadPid := isPidRunning(GConfig.PIDS...); !hasDeadPid {

				// 推送数据至influxdb
				pushToInfluxDB(writeAPI, pids...)

				// 如果之前发生告警则触发告警恢复
				if isAlarm {
					cancelAlarmOnce.Do(func() {
						isAlarm = false
						/*
							告警恢复逻辑写在此处
						*/
						alarmOnce = sync.Once{}
					})
				}

				fmt.Printf("running %s\n", strings.Join(pids, ","))
			} else {
				GConfig.PIDS = getProcessPID(GConfig.Name)
				// TODO: 应该触发一次告警,然后尝试重新获取pid
				alarmOnce.Do(func() {
					isAlarm = true
					/*
						告警触发逻辑写在此处
					*/
					cancelAlarmOnce = sync.Once{}
				})
				fmt.Printf("not running %s\n", strings.Join(pids, ","))
				time.Sleep(time.Second)
			}
		} else {
			fmt.Println("nothing happend")
			time.Sleep(time.Second * 5)
		}
	}
}
