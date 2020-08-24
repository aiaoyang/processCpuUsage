package sysusage

import (
	"fmt"
	"math"
	"runtime"
	"strconv"
	"time"
)

// sec = jiffies / HZ ; here - HZ = number of ticks per second

// HZ hz
const HZ = 100.0

// ProcessCPUStat 进程cpu状态
type ProcessCPUStat struct {
	pid    string
	utime  uint64
	stime  uint64
	cutime uint64
	cstime uint64
	used   uint64
	time   float64
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
	c.time = float64(time.Now().UnixNano()) / float64(1e9)
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
	fmt.Printf("delta time : %.2f\n", delta.time)
	finalUsage := Usage(delta.used) / Usage(delta.time*float64(runtime.NumCPU()))
	isNaN := math.IsNaN(float64(finalUsage))
	if isNaN {
		return -1
	}
	return finalUsage
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
	finalUsage := Usage(delta.used) * 100 / Usage(delta.total)
	isNaN := math.IsNaN(float64(finalUsage))
	if isNaN {
		return -1
	}
	return finalUsage
}
