package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"strconv"
)

func isPidRunning(pids ...int) ([]int, bool) {

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
		return deadPids, true
	}

	return pids, false
}

func getProcessPID(name string) []int {

	res := []int{}

	// 获取所有进程信息 map[pid]cmdline
	processes := getAllProcess()

	for k, v := range processes {
		if bytes.Contains(v, []byte(name)) {
			res = append(res, k)
		}
	}

	return res
}

// 获取所有运行的进程 map[pid]processName
func getAllProcess() map[int][]byte {

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
