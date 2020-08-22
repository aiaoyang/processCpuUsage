package sysusage

import (
	"strconv"
	"unsafe"
)

const loadFS = "/proc/loadavg"

// LoadType LoadType
type LoadType int

const (
	Load1  LoadType = 1
	Load5           = 5
	Load15          = 15
)

// Load 返回主机负载情况
// func Load(ctx context.Context, reciver map[string]float64) {
func Load() map[LoadType]Usage {

	metric := make(map[LoadType]Usage)

	l, _ := HeadLineSplitOfFile(loadFS)

	metric[Load1] = Usage(byteToFloat64(l[0]))
	metric[Load5] = Usage(byteToFloat64(l[1]))
	metric[Load15] = Usage(byteToFloat64(l[2]))

	return metric
}

func byteToFloat64(b []byte) float64 {
	f64, err := strconv.ParseFloat(*(*string)(unsafe.Pointer(&b)), 64)
	if err != nil {
		panic(err)
	}
	return f64
}
