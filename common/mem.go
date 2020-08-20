package common

import (
	"bufio"
	"bytes"
	"context"
	"io"
	"log"
	"os"
	"strconv"
	"time"
)

const (
	sysMemFileName = "/proc/meminfo"
	sysPid         = "0"
)

var totalMem uint64

// IMem 内存接口
type IMem interface {
	Used() uint64
	Total() uint64
	Usage(uint64) float64
}

// ProcessMemStat 进程内存使用量
type ProcessMemStat struct {
	pid      string
	size     uint64
	resident uint64
	shared   uint64
	text     uint64
	lib      uint64
	data     uint64
	dt       uint64
}

// MemStat 系统内存使用状态
type MemStat struct {
	pid       string
	total     uint64
	free      uint64
	available uint64
	buffers   uint64
	cached    uint64
}

// =============================================================

// Gen 统计一次进程内存使用情况
func (m *ProcessMemStat) Gen(pid int) *ProcessMemStat {
	p := strconv.Itoa(pid)
	filename := "/proc/" + p + "/statm"

	line, err := HeadLineSplitOfFile(filename)
	if err != nil {
		log.Printf("err: %s\n", err.Error())
		return nil
	}
	m.pid = p
	m.size = byteToUint(line[0])
	m.resident = byteToUint(line[1])
	m.shared = byteToUint(line[2])
	m.text = byteToUint(line[3])
	m.lib = byteToUint(line[4])
	m.data = byteToUint(line[5])
	m.dt = byteToUint(line[6])
	return m
}

// Used 进程内存使用量
func (m *ProcessMemStat) Used() uint64 {
	// resident为页大小，linux页一般为4KB
	if m == nil {
		return 0
	}
	return m.resident * 4
}

// Total 进程内存使用量
func (m *ProcessMemStat) Total() uint64 {
	return totalMem
}

// Usage 进程内存使用率
func (m *ProcessMemStat) Usage() float64 {
	return float64(m.Used()) * 100 / float64(m.Total())
}

//==============================================================

// Gen 统计一次系统内存使用情况
func (m *MemStat) Gen(pid string) *MemStat {
	line, err := readMemFile()
	if err != nil {
		return nil
	}
	m.pid = pid
	m.total = byteToUint(line[0])
	m.free = byteToUint(line[1])
	m.available = byteToUint(line[2])
	m.buffers = byteToUint(line[3])
	m.cached = byteToUint(line[4])
	return m
}

// Used 系统内存使用量
func (m *MemStat) Used() uint64 {
	return m.total - m.available
}

// Total 系统内存使用量
func (m *MemStat) Total() uint64 {
	return m.total
}

// Usage 系统内存使用率
func (m *MemStat) Usage() float64 {
	return float64(m.Used()) * 100 / float64(m.Total())
}

//==============================================================

func init() {
	totalMem = TotalMem()
}

// ProcessMemUsage 进程内存使用率
func ProcessMemUsage(ctx context.Context, reciver chan map[string]float64, pidChan chan []int) {
	processMem := &ProcessMemStat{}
	sysMem := &MemStat{}
	metric := make(map[string]float64, 0)
	for {
		select {
		case <-ctx.Done():
			return
		default:
			select {
			case pids := <-pidChan:
				for _, pid := range pids {

					metric["system_mem_usage"] = sysMem.Gen(sysPid).Usage()

					processMem = processMem.Gen(int(pid))
					if processMem == nil {

						delete(metric, "pid")
						delete(metric, "process_mem_usage")

						continue
					}

					metric["pid"] = float64(pid)
					metric["process_mem_usage"] = processMem.Usage()

				}
				reciver <- metric
			default:
				time.Sleep(memUsageDuration)
				continue
			}
		}
	}
}

// ProcessMemUsageOnce 统计一次内存使用量
func ProcessMemUsageOnce(pid string) float64 {
	processMem := &ProcessMemStat{}

	pidInt, err := strconv.Atoi(pid)
	if err != nil {
		return 0.0
	}
	return processMem.Gen(pidInt).Usage()
}

// SystemMemUsageOnce 统计一次系统内存使用量
func SystemMemUsageOnce() float64 {
	sys := &MemStat{}
	return sys.Gen(sysPid).Usage()
}

// TotalMem 系统内存总量
func TotalMem() uint64 {
	line, err := HeadLineSplitOfFile(sysMemFileName)
	if err != nil {
		// 不支持非linux系统
		panic(err)
	}
	return byteToUint(line[1])
}

func readMemFile() ([][]byte, error) {
	f, err := os.Open(sysMemFileName)
	defer f.Close()
	if err != nil {
		return nil, err
	}
	sliceLine := make([][]byte, 5)
	buf := bufio.NewReader(f)
	for i := 0; i < 5; i++ {
		line, _, err := buf.ReadLine()
		if err != nil && err != io.EOF {
			return nil, err
		}
		value := bytes.Fields(line)[1]
		sliceLine[i] = value
	}

	return sliceLine, nil
}
