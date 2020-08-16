package common

import (
	"context"
	"log"
	"runtime"
	"strconv"
	"sync"
	"time"
)

// None some interface need implement func sign, but do nothing
const None = "this value do nothing"

// HZ grep 'define HZ' /usr/include/asm*/param.h
const HZ = 100

// ICPUStat cpu状态接口
type ICPUStat interface {
	Pid() string
	Usage(total uint64) float64
	Used() uint64
	Total() uint64
}

// ProcessCPUStat 进程cpu状态
type ProcessCPUStat struct {
	pid    string
	utime  uint64
	stime  uint64
	cutime uint64
	cstime uint64
	used   uint64
	time   int64
	isDead bool
}

// New 进程cpu使用量统计
func (c *ProcessCPUStat) New(opt int32) {
	pidString := strconv.Itoa(int(opt))
	// c = &ProcessCPUStat{opt, 0, 0, 0, 0, 0}
	fileName := "/proc/" + pidString + "/stat"
	byteSlices, err := HeadLineSplitOfFile(fileName)
	if err != nil {
		c.isDead = true
		return
	}

	c.utime = byteToUint(byteSlices[13])
	c.stime = byteToUint(byteSlices[14])
	c.cutime = byteToUint(byteSlices[15])
	c.cstime = byteToUint(byteSlices[16])
	c.time = time.Now().Unix()
	c.used = c.utime + c.stime + c.cutime + c.cstime
}

// CPUStat 系统cpu状态
type CPUStat struct {
	user    uint64
	nice    uint64
	system  uint64
	idle    uint64
	iowait  uint64
	irq     uint64
	softirq uint64
	steal   uint64
	// 以下两个指标分别被包含在 user、nice中
	// guest     uint64
	// guestNice uint64
	used  uint64
	total uint64
}

// New 系统cpu使用量统计
func (c *CPUStat) New(opt string) {
	// c = &CPUStat{0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	fileName := "/proc/stat"
	byteSlices, err := HeadLineSplitOfFile(fileName)
	if err != nil {
		return
	}

	c.user = byteToUint(byteSlices[1])
	c.nice = byteToUint(byteSlices[2])
	c.system = byteToUint(byteSlices[3])
	c.idle = byteToUint(byteSlices[4])
	c.iowait = byteToUint(byteSlices[5])
	c.irq = byteToUint(byteSlices[6])
	c.softirq = byteToUint(byteSlices[7])
	c.steal = byteToUint(byteSlices[8])

	c.total = c.user + c.nice + c.system + c.idle + c.iowait + c.irq + c.softirq + c.steal
	c.used = c.total - c.idle
}

// ProcessesCPUMonitor 进程cpu使用率和系统cpu使用率, 如果需要获取进程cpu使用率
// 系统cpu使用率： reciver -> map["system"]
// 进程cpu使用率： reciver -> map[pid]
/*
	传入pid的方式参考如下

			pidsChan := make(chan []string, 0)
			pidsChan <- pids
			time.Sleep(duration)

*/
func ProcessesCPUMonitor(ctx context.Context, reciver chan map[string]float64, pidschan chan []int32) {

	cpuCoreNum := runtime.NumCPU()

	process := &ProcessCPUStat{}
	system := &CPUStat{}

	oldPorcess := &ProcessCPUStat{}
	oldSystem := &CPUStat{}

	var wg = &sync.WaitGroup{}

	for {

		metric := make(map[string]float64)

		metric["core_num"] = float64(cpuCoreNum)

		select {

		case <-ctx.Done():
			log.Println("stopping")
			return

		case pids := <-pidschan:

			for _, pid := range pids {

				wg.Add(2)

				// 计算系统cpu总使用率
				// currSystemCPUCount := &CPUStat{}
				go func(wg *sync.WaitGroup) {
					system.New(sysPid)

					wg.Done()
				}(wg)

				//开始计算当前pid进程cpu统计量以及与上一次cpu统计量的变化量
				// currProcessCPUCount := &ProcessCPUStat{}
				go func(wg *sync.WaitGroup) {
					process.New(pid)

					wg.Done()
				}(wg)

				wg.Wait()

				metric["sys_cpu_usage"] = float64(system.used-oldSystem.used) * 100 / float64(system.total-oldSystem.total+1)

				// 进程不存在，则只返回系统cpu使用情况
				if process.isDead {
					continue
				}

				metric["pid"] = float64(pid)

				// 计算本次进程cpu和上次进程cpu差异，并与总cpu变化求商得使用率
				//										进程cpu使用量-上次进程cpu使用量									系统cpu使用量-上次系统cpu使用量						系统cpu核数
				metric["process_cpu_usage"] = float64(process.used-oldPorcess.used) * 100 / (float64(process.time-oldPorcess.time+1) * HZ * float64(cpuCoreNum))

				// 将本次进程cpu统计结果设置为旧
				*oldPorcess = *process

				// 将本次系统cpu总量果设置为旧
				*oldSystem = *system

			}
			reciver <- metric
		default:
			time.Sleep(cpuUsageDuration)
		}
	}
}
