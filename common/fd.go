package common

import (
	"io/ioutil"
)

// OpenFD 进程文件打开数量
func OpenFD(pid string) float64 {
	f, err := ioutil.ReadDir("/proc/" + pid + "/fd")
	if err != nil {
		return 0
	}

	return float64(len(f))
}
