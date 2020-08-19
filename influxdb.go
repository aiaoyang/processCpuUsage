package main

import (
	"log"
	"os"
	"time"

	"github.com/aiaoyang/processCpuUsage/common"
	influxc "github.com/influxdata/influxdb-client-go"
	"github.com/influxdata/influxdb-client-go/api"
)

func pushToInfluxDB(writeAPI api.WriteAPI,
	alarmer AlarmActor,
	tag map[string]string,
	value map[string]interface{},
) {
	hostName, err := os.Hostname()
	if err != nil {
		log.Fatal(err)
	}

	tag["hostname"] = hostName
	tag["env"] = MyLocalConfig.Env

	fd := common.OpenFD(tag["pid"])
	if fd > 0 {
		value["fd"] = fd
	}
	loads := common.Load()

	value["load_1"] = loads["load_1"]
	value["load_5"] = loads["load_5"]
	value["load_15"] = loads["load_15"]

	value["process_mem_usage"] = common.ProcessMemUsageOnce(tag["pid"])
	value["sys_mem_usage"] = common.SystemMemUsageOnce()

	// create point
	p := influxc.NewPoint(
		"system",
		tag,
		value,
		time.Now(),
	)

	writeAPI.WritePoint(p)
	writeAPI.Flush()
}

// func processName(pid string) string {

// 	content, err := ioutil.ReadFile("/proc/" + pid + "/comm")
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	return fmt.Sprintf("%s", content)
// }
