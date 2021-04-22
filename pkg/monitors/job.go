package monitors

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/aiaoyang/processCpuUsage/configs"
	"github.com/aiaoyang/processCpuUsage/pkg/db"
	"github.com/aiaoyang/processCpuUsage/pkg/metric"
	"github.com/aiaoyang/processCpuUsage/pkg/monitors/process"
	"github.com/aiaoyang/processCpuUsage/pkg/monitors/system"
	influxdb2 "github.com/influxdata/influxdb-client-go"
	"github.com/pkg/errors"
)

type CtxValueKey bool

var (
	configPath = new(string)
	cfg        = configs.Config{}
	acmConfig  = configs.AcmInnerConfig{}
)

func init() {
	flag.StringVar(configPath, "c", "config.yaml", "指定默认配置文件")
	flag.Parse()

	if !flag.Parsed() {
		flag.Usage()
		os.Exit(1)
	}

}

func connectInfluxDB() (sender metric.Sender, err error) {
	influxdbString := fmt.Sprintf("http://%s:%s",
		cfg.InfluxDB.Host,
		cfg.InfluxDB.Port,
	)

	c := influxdb2.NewClient(influxdbString, "")
	ok, err := c.Health(context.TODO())
	if err != nil {
		return
	}
	if ok.Status != "pass" {
		err = errors.Errorf("connect to db err: %+v", ok)
		return
	}

	writeApi := c.WriteAPI("test", "test")

	go func() {
		for {
			writeApi.Flush()
			time.Sleep(time.Second * 10)
		}
	}()

	sender = db.NewInfluxDBSender(writeApi)
	return
}

func StartMonitor(hostname, env string) {
	acmChan := make(chan *configs.AcmInnerConfig)
	acmConfig = configs.AcmInnerConfig{}

	ctx := context.WithValue(context.Background(), CtxValueKey(true), &acmConfig)

	cfg = configs.LoadConfig(*configPath)

	err := configs.LoadAcmConfig(ctx, cfg, acmChan)
	if err != nil {
		log.Fatal(err)
	}

	sender, err := connectInfluxDB()
	if err != nil {
		log.Fatal(err)
	}

	go (&acmConfig).Watch(ctx, acmChan)
	go process.Monitor(ctx, sender)
	go system.Monitor(ctx, sender)
	go healthCheck(ctx, acmChan)
}

func healthCheck(ctx context.Context, acmChan chan *configs.AcmInnerConfig) {

}
