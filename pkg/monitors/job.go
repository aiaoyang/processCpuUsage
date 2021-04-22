package monitors

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/aiaoyang/processCpuUsage/configs"
	"github.com/aiaoyang/processCpuUsage/pkg/db"
	"github.com/aiaoyang/processCpuUsage/pkg/metric"
	"github.com/aiaoyang/processCpuUsage/pkg/monitors/process"
	"github.com/aiaoyang/processCpuUsage/pkg/monitors/system"
	"github.com/golang/glog"
	influxdb2 "github.com/influxdata/influxdb-client-go"
	"github.com/pkg/errors"
)

type CtxValueKey bool

var (
	configPath = new(string)
	// cfg        = configs.Config{}
	acmConfig = configs.AcmInnerConfig{}
)

func init() {
	flag.StringVar(configPath, "c", "config.yaml", "指定默认配置文件")
	flag.Parse()

	if !flag.Parsed() {
		flag.Usage()
		os.Exit(1)
	}

}

func connectInfluxDB(cfg configs.Config) (sender metric.Sender, err error) {
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

	writeApi := c.WriteAPI(cfg.DBName, cfg.TableName)

	go func() {
		for {
			writeApi.Flush()
			time.Sleep(time.Second * 10)
		}
	}()

	sender = db.NewInfluxDBSender(writeApi)
	return
}

func StartMonitor(ctx context.Context, hostname, env string) (err error) {
	defer ctx.Done()

	value := ctx.Value(CtxValueKey(true))
	cfg, ok := value.(configs.Config)
	if !ok {
		err = errors.Errorf("config value not found, got: %+v", value)
		return
	}

	acmChan := make(chan *configs.AcmInnerConfig)

	acmConfig = configs.AcmInnerConfig{}

	err = configs.LoadAcmConfig(ctx, cfg, acmChan)
	if err != nil {
		return
	}

	go (&acmConfig).Watch(ctx, acmChan)
	go pidCheck(ctx, &acmConfig)

	sender, err := connectInfluxDB(cfg)
	if err != nil {
		return
	}

	errChan := make(chan error)

	go process.Monitor(ctx, sender, errChan)
	go system.Monitor(ctx, sender, errChan)

	for e := range errChan {
		glog.Errorf("got err: %+v\n", e)
	}
	return
}

func pidCheck(ctx context.Context, cfg *configs.AcmInnerConfig) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			cfg.Lock()
			glog.Infof("stored pid: %v\n", cfg.StoredPids)
			cfg.Unlock()
			time.Sleep(time.Second)
		}
	}
}
