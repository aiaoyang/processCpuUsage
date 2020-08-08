package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sync"
	"time"

	influxc "github.com/influxdata/influxdb-client-go"
	"github.com/influxdata/influxdb-client-go/api"
)

func pushToInfluxDB(writeAPI api.WriteAPI, pids ...string) {
	hostName, err := os.Hostname()
	if err != nil {
		log.Fatal(err)
	}
	// writeAPI := c.WriteAPI("test", "test")

	wg := &sync.WaitGroup{}
	wg.Add(len(pids))

	pushInflux := func(wg *sync.WaitGroup, hostName, pid string) {
		// create point
		p := influxc.NewPoint(
			"system",
			map[string]string{
				"hostname": hostName,
				"pid":      pid,
				"env":      CLocal.Env,
				"process":  processName(pid),
			},
			map[string]interface{}{
				"cpu_usage": getCPUUsage(pid),
				"mem_usage": getMemUsage(pid),
			},
			time.Now())
		// write asynchronously
		writeAPI.WritePoint(p)
		writeAPI.Flush()
		wg.Done()
	}

	for _, pid := range pids {
		go pushInflux(wg, hostName, pid)
	}
	wg.Wait()
}

func processName(pid string) string {
	content, err := ioutil.ReadFile("/proc/" + pid + "/comm")
	if err != nil {
		log.Fatal(err)
	}
	return fmt.Sprintf("%s", content)
}
