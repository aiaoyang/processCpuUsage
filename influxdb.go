package main

import (
	"log"
	"os"
	"time"

	influxc "github.com/influxdata/influxdb-client-go"
	"github.com/influxdata/influxdb-client-go/api"
)

func pushToInfluxDB(writeAPI api.WriteAPI,
	tag map[string]string,
	value map[string]interface{},
) {
	hostName, err := os.Hostname()
	if err != nil {
		log.Fatal(err)
	}

	tag["hostname"] = hostName
	tag["env"] = MyLocalConfig.Env

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
