package sysusage

import (
	"io/ioutil"
	"strconv"
)

// OpenFD 进程文件打开数量
func OpenFD(pid int) Usage {
	f, err := ioutil.ReadDir("/proc/" + strconv.Itoa(pid) + "/fd")
	if err != nil {
		return -1
	}
	return Usage(len(f))
}
