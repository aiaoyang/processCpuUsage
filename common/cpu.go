package common

import (
	"context"
	"log"
	"runtime"
	"strconv"
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
func (c *ProcessCPUStat) New(opt int) {
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
func ProcessesCPUMonitor(ctx context.Context, reciver chan map[string]interface{}, pidschan chan []int) {

	cpuCoreNum := runtime.NumCPU()

	process := &ProcessCPUStat{}
	system := &CPUStat{}

	// 保存上一次系统cpu使用状态
	oldsys := make(map[int]CPUStat)

	// 保存上一次进程cpu使用状态
	oldprocess := make(map[int]ProcessCPUStat)

	// 同时统计进程与系统cpu使用量
	// var wg = &sync.WaitGroup{}

	for {
		select {

		case <-ctx.Done():
			log.Println("stopping")
			return

		case pids := <-pidschan:

			for _, pid := range pids {

				// 保存 使用率 至指标中
				metric := make(map[string]interface{})

				// 统计系统核数
				metric["core_num"] = float64(cpuCoreNum)

				// 计算系统cpu总使用量
				system.New(sysPid)

				// 计算进程cpu使用量
				process.New(pid)

				metric["sys_cpu_usage"] = float64(system.used-oldsys[pid].used) * 100 / float64(system.total-oldsys[pid].total+1)

				// 进程不存在，则只返回系统cpu使用情况
				if process.isDead {
					reciver <- metric
					continue
				}

				metric["pid"] = float64(pid)

				// 计算本次进程cpu和上次进程cpu差异，并与总cpu变化求商得使用率
				//										进程cpu使用量-上次进程cpu使用量									系统cpu使用量-上次系统cpu使用量						系统cpu核数
				// 该计算方式与 top 命令一致，做了一点修改：统计的 cpu 使用率 = top数值/cpu核心数
				metric["process_cpu_usage"] = float64(process.used-oldprocess[pid].used) * 100 / (float64(process.time-oldprocess[pid].time+1) * HZ * float64(cpuCoreNum))

				reciver <- metric

				// 将本次系统cpu总量果设置为旧
				oldsys[pid] = *system

				// 将本次进程cpu统计结果设置为旧
				oldprocess[pid] = *process

			}

		default:
			time.Sleep(cpuUsageDuration)
		}
	}
}
