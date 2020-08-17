package main

import (
	"context"
	"fmt"
	"log"
	"strconv"

	// _ "net/http/pprof"

	"sync"
	"time"

	"github.com/aiaoyang/processCpuUsage/common"
	influxc "github.com/influxdata/influxdb-client-go"
)

func init() {
	initViper()
	// log.SetFlags(log.Llongfile | log.Ldate)
	// log.SetOutput(os.Stdout)
}

// var pids string

func main() {
	// go http.ListenAndServe("0.0.0.0:10088", nil)

	// log.SetFlags(log.Ltime | log.Llongfile | log.LstdFlags)
	todo := context.TODO()

	go MyLocalConfig.watchConfigChange(todo)

	genericTODO(nil)
}

func genericTODO(alarm func(msg interface{})) {

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

	defer c.Close()
	defer writeAPI.Close()

	reciver := make(chan map[string]interface{}, 0)

	isAlarm := false
	alarmOnce := sync.Once{}
	cancelAlarmOnce := sync.Once{}

	pidsChan := make(chan []int, 0)

	go common.ProcessesCPUMonitor(ctx, reciver, pidsChan)

	for {
		select {
		case <-ctx.Done():
			return

		// 接收到指标参数，则将其推送至influxdb
		case res := <-reciver:
			tag := make(map[string]string)

			tag["pid"] = strconv.Itoa(int(res["pid"].(float64)))

			// 删除不需要res中的pid键
			delete(res, "pid")

			// debug用
			fmt.Printf("recive value : %v\n", res)

			// 推送数据至influxdb
			pushToInfluxDB(writeAPI, tag, res)

		default:
			log.Println("here")
			if AliyunConfigSrv.GetInt("processinfo.status") == 1 {

				// 如果进程仍在运行，则将 pids 发送给 通道 然后让指标收集函数进行处理
				if pids, hasDeadPid := common.IsPidRunning(NeedMonitorProcessInfo.PIDS...); !hasDeadPid {

					pidsChan <- pids

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

					// debug 打印
					fmt.Printf("running %v\n", pids)

				} else { // 如果进程被中断，则循环监听进程列表，直至同名进程启动

					fmt.Printf("pid [%v] is not running\n", pids)

					// TODO: 应该触发一次告警,然后尝试重新获取pid
					alarmOnce.Do(func() {

						isAlarm = true
						/*
							告警触发逻辑写在此处
							alarm()
						*/
						cancelAlarmOnce = sync.Once{}

					})

					// 查询进程pid是否存在
					NeedMonitorProcessInfo.PIDS = common.GetProcessPID(NeedMonitorProcessInfo.Name)

				}
			} else {

				fmt.Println("nothing happend")

			}

			time.Sleep(time.Second * 3)
		}
	}
}
