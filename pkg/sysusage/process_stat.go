package sysusage

import (
	"bytes"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"sync"
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
// TODO:大量进程的情况下需要优化
func GetProcessPID(names ...string) []int {

	res := []int{}

	// 获取所有进程信息 map[pid]cmdline
	processes := GetAllProcess()

	for k, v := range processes {
		for _, name := range names {
			if bytes.Contains(v, []byte(name)) {
				res = append(res, k)
			}
		}
	}

	return res
}

// GetAllProcess 获取所有运行的进程 map[pid]processName
func GetAllProcess() map[int][]byte {

	fileInfoList, err := ioutil.ReadDir("/proc/")
	if err != nil {
		log.Fatal(err)
	}

	mapChan := make(chan map[int][]byte, len(fileInfoList))
	wg := &sync.WaitGroup{}
	wg.Add(len(fileInfoList))

	for _, pidFile := range fileInfoList {
		go func(wg *sync.WaitGroup, ch chan map[int][]byte, pFile fs.FileInfo) {
			defer wg.Done()

			pidInt, err := strconv.Atoi(pFile.Name())
			if err != nil {
				return
			}

			// cmdLine, err := ioutil.ReadFile("/proc/" + pFile.Name() + "/cmdline")

			f, err := os.Open("/proc/" + pFile.Name() + "/cmdline")

			if err != nil {
				return
			}
			defer f.Close()

			buf := make([]byte, 256)
			n, err := f.Read(buf)
			if err != nil {
				return
			}

			if bytes.Equal(buf[:n], []byte("")) {
				return
			}

			mapChan <- map[int][]byte{pidInt: buf[:n]}
		}(wg, mapChan, pidFile)

	}
	wg.Wait()
	close(mapChan)

	res := make(map[int][]byte)
	for _map := range mapChan {
		for k, v := range _map {
			res[k] = v
		}
	}
	return res

}
