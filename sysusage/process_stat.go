package sysusage

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"strconv"
)

// IsPidRunning 进程是否运行
// 返回未运行的进程pid
func IsPidRunning(pids ...int) ([]int, bool) {

	if len(pids) == 0 {
		return nil, false
	}

	deadPids := []int{}
	hasDeadPid := false

	for _, pid := range pids {
		_, err := os.Stat("/proc/" + strconv.Itoa(int(pid)))
		if err != nil {

			hasDeadPid = true

			deadPids = append(deadPids, pid)

		}
	}

	if hasDeadPid {
		return deadPids, false
	}

	return nil, true
}

// GetProcessPID 获取进程pid
func GetProcessPID(name string) []int {

	res := []int{}

	// 获取所有进程信息 map[pid]cmdline
	processes := GetAllProcess()

	for k, v := range processes {
		if bytes.Contains(v, []byte(name)) {
			res = append(res, k)
		}
	}

	return res
}

// GetAllProcess 获取所有运行的进程 map[pid]processName
func GetAllProcess() map[int][]byte {

	res := make(map[int][]byte)

	fileInfoList, err := ioutil.ReadDir("/proc/")
	if err != nil {
		log.Fatal(err)
	}

	for _, pidFile := range fileInfoList {

		pidInt, err := strconv.Atoi(pidFile.Name())
		if err != nil {
			continue
		}

		cmdLine, err := ioutil.ReadFile("/proc/" + pidFile.Name() + "/cmdline")

		if err != nil {
			log.Fatal(err)
		}

		if bytes.Compare(cmdLine, []byte("")) == 0 {
			continue
		}

		res[pidInt] = cmdLine

	}

	return res
}
