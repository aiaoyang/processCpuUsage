package sysusage

import (
	"runtime"
	"strconv"
	"time"
)

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
func (c *ProcessCPUStat) New(pid int) {
	// 指针结构 每次调用方法时需要 设置默认值，因为会被中断的进程修改
	c.isDead = false

	if pid <= 0 {
		c.isDead = true
	}

	pidString := strconv.Itoa(pid)

	fileName := "/proc/" + pidString + "/stat"

	byteSlices, err := HeadLineSplitOfFile(fileName)
	if err != nil {
		c.isDead = true
		return
	}
	// 尽量接近文件读入的时间
	// 生成进程cpu快照时的时间
	c.time = time.Now().Unix()

	c.utime = byteToUint(byteSlices[13])
	c.stime = byteToUint(byteSlices[14])
	c.cutime = byteToUint(byteSlices[15])
	c.cstime = byteToUint(byteSlices[16])

	c.used = c.utime + c.stime + c.cutime + c.cstime
}

// SysCPUStat 系统cpu状态
type SysCPUStat struct {
	user    uint64
	nice    uint64
	system  uint64
	idle    uint64
	iowait  uint64
	irq     uint64
	softirq uint64
	steal   uint64
	isDead  bool
	// 以下两个指标分别被包含在 user、nice中
	// guest     uint64
	// guestNice uint64
	used  uint64
	total uint64
}

func (c *SysCPUStat) sub(sub *SysCPUStat) *SysCPUStat {
	return &SysCPUStat{
		used:  c.used - sub.used,
		total: c.total - sub.total,
	}
}

// New 系统cpu使用量统计
func (c *SysCPUStat) New(pid int) {
	// c = &CPUStat{0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	fileName := "/proc/stat"
	byteSlices, err := HeadLineSplitOfFile(fileName)
	if err != nil {
		c.isDead = true
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

func (c *ProcessCPUStat) sub(sub *ProcessCPUStat) *ProcessCPUStat {
	return &ProcessCPUStat{
		used: c.used - sub.used,
		time: c.time - sub.time,
	}
}

///////////////////////////////////////////////
//   					单次cpu使用量统计								//
///////////////////////////////////////////////

// ProcessCPUUsageOnce 统计一次进程cpu使用率
func ProcessCPUUsageOnce(pid int, duration time.Duration) Usage {
	t1 := &ProcessCPUStat{}
	t1.New(pid)
	time.Sleep(duration)
	t2 := &ProcessCPUStat{}
	t2.New(pid)

	if t1.isDead || t2.isDead {
		// 如果初始化失败，返回-1
		return Usage(-1)
	}
	delta := t2.sub(t1)
	return Usage(delta.used) / Usage(delta.time+1) / Usage(runtime.NumCPU())
}

// SysCPUUsageOnce 统计一次系统cpu使用率
func SysCPUUsageOnce(duration time.Duration) Usage {
	t1 := &SysCPUStat{}
	t1.New(sysPid)
	time.Sleep(duration)
	t2 := &SysCPUStat{}
	t2.New(sysPid)
	if t1.isDead || t2.isDead {
		// 如果初始化失败，返回-1
		return Usage(-1)
	}
	delta := t2.sub(t1)
	return Usage(delta.used) / Usage(delta.total+1)
}

// ProcessesCPUMonitor 进程cpu使用率和系统cpu使用率, 如果需要获取进程cpu使用率
// 系统cpu使用率： reciver -> map["system"]
// 进程cpu使用率： reciver -> map[pid]
// 当pidsChan输出负值时，则只监控系统cpu使用情况
/*
	传入pid的方式参考如下

			pidsChan := make(chan []string, 0)
			pidsChan <- pids
			time.Sleep(duration)

*/
// func ProcessesCPUMonitor(ctx context.Context, reciver chan<- metric.Metric, pidschan chan []int) {

// 	cpuCoreNum := runtime.NumCPU()

// 	process := &ProcessCPUStat{}
// 	system := &CPUStat{}

// 	// 保存上一次系统cpu使用状态
// 	oldsys := make(map[int]CPUStat)

// 	// 保存上一次进程cpu使用状态
// 	oldprocess := make(map[int]ProcessCPUStat)

// 	// 同时统计进程与系统cpu使用量
// 	// var wg = &sync.WaitGroup{}

// 	for {
// 		select {

// 		case <-ctx.Done():
// 			log.Println("stopping")
// 			return

// 		case pids := <-pidschan:

// 			for _, pid := range pids {

// 				// 保存 使用率 至指标中
// 				metric := make(metric.Metric)

// 				// 计算系统cpu总使用量
// 				system.New(sysPid)

// 				// 计算进程cpu使用量
// 				process.New(pid)

// 				sysCPUUsage := Usage(system.used-oldsys[pid].used) * 100 / Usage(system.total-oldsys[pid].total+1)

// 				// 统计系统核数
// 				metric.Insert("core_num", Usage(cpuCoreNum))
// 				metric.Insert("sys_cpu_usage", sysCPUUsage)

// 				// 进程不存在，则只返回系统cpu使用情况
// 				if process.isDead {
// 					reciver <- metric
// 					continue
// 				}

// 				// 计算本次进程cpu和上次进程cpu差异，并与总cpu变化求商得使用率
// 				//										进程cpu使用量-上次进程cpu使用量									系统cpu使用量-上次系统cpu使用量						系统cpu核数
// 				// 该计算方式与 top 命令一致，做了一点修改：统计的 cpu 使用率 = top数值/cpu核心数

// 				processCPUUsage := Usage(process.used-oldprocess[pid].used) * 100 / (Usage(process.time-oldprocess[pid].time+1) * HZ * Usage(cpuCoreNum))

// 				metric.Insert("pid", Usage(pid))
// 				metric.Insert("process_cpu_usage", processCPUUsage)

// 				reciver <- metric

// 				// 将本次系统cpu总量果设置为旧
// 				oldsys[pid] = *system

// 				// 将本次进程cpu统计结果设置为旧
// 				oldprocess[pid] = *process

// 			}

// 		default:
// 			time.Sleep(cpuUsageDuration)
// 		}
// 	}
// }
