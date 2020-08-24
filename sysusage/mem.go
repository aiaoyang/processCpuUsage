package sysusage

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"strconv"
)

var totalMem uint64

func init() {
	totalMem = TotalMem()
}

const (
	sysMemFileName = "/proc/meminfo"
	sysPid         = 0
)

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
	isDead   bool
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

// ProcessMemUsageOnce 统计一次内存使用量
func ProcessMemUsageOnce(pid int) Usage {
	processMem := &ProcessMemStat{}
	processMem.New(pid)
	if processMem.isDead {
		return -1
	}
	return processMem.Usage()
}

// SystemMemUsageOnce 统计一次系统内存使用量
func SystemMemUsageOnce() Usage {
	sys := &MemStat{}
	return sys.New(sysPid).Usage()
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

// =============================================================

// New 统计一次进程内存使用情况
func (m *ProcessMemStat) New(pid int) *ProcessMemStat {
	p := strconv.Itoa(pid)
	filename := "/proc/" + p + "/statm"

	line, err := HeadLineSplitOfFile(filename)
	if err != nil {
		m.isDead = true
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
	m.isDead = false
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
func (m *ProcessMemStat) Usage() Usage {
	return Usage(m.Used()) * 100 / Usage(m.Total())
}

//==============================================================

// New 统计一次系统内存使用情况
func (m *MemStat) New(pid int) *MemStat {
	line, err := readMemFile()
	if err != nil {
		return nil
	}
	m.pid = strconv.Itoa(pid)
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
func (m *MemStat) Usage() Usage {
	return Usage(m.Used()) * 100 / Usage(m.Total())
}
