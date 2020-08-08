package main

import (
	"context"
	"fmt"
	"log"

	// _ "net/http/pprof"
	"strings"
	"sync"
	"time"

	influxc "github.com/influxdata/influxdb-client-go"
)

var conf = &ALiConfigClientConfig{}

func init() {
	initViper()
	conf.new()
}

func main() {
	log.SetFlags(log.Llongfile | log.Ltime)
	todo := context.TODO()
	go conf.watchConfigChange(todo)
	// go http.ListenAndServe("0.0.0.0:10088", nil)
	genericTODO(nil)
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
			time.Sleep(time.Second)
		}
	}
}
