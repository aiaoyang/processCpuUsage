package monitors

import (
	"log"
	"testing"

	"github.com/aiaoyang/processCpuUsage/configs"
)

func Test_db(t *testing.T) {
	cfg := configs.LoadConfig("config.yaml")
	_, err := connectInfluxDB(cfg)
	if err != nil {
		log.Fatalf("err: %+v", err)
	}
}
