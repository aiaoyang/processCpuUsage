package common

import (
	"strconv"
	"unsafe"
)

const loadfs = "/proc/loadavg"

// Load 返回主机负载情况
// func Load(ctx context.Context, reciver map[string]float64) {
func Load() map[string]float64 {

	metric := make(map[string]float64)

	l, _ := HeadLineSplitOfFile(loadfs)

	metric["load_1"] = byteToFloat64(l[0])
	metric["load_5"] = byteToFloat64(l[1])
	metric["load_15"] = byteToFloat64(l[2])

	return metric
}

func byteToFloat64(b []byte) float64 {
	f64, err := strconv.ParseFloat(*(*string)(unsafe.Pointer(&b)), 64)
	if err != nil {
		panic(err)
	}
	return f64
}
