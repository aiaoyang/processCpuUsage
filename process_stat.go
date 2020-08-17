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

// // TODO: 缓存上次统计的结果，不必每次计算使用率都读取两次文件内容
// func getCPUUsage(pid string) float64 {
// 	totalCPUTimeChan := make(chan int64, 2)
// 	totalThreadTimeChan := make(chan int64, 2)

// 	wg := &sync.WaitGroup{}
// 	wg.Add(2)
// 	// 获取第一次cpu信息
// 	go totalCPUTime(wg, totalCPUTimeChan)
// 	go totalThreadTime(pid, wg, totalThreadTimeChan)
// 	wg.Wait()

// 	time.Sleep(1)

// 	wg.Add(2)
// 	// 获取第二次cpu信息
// 	go totalCPUTime(wg, totalCPUTimeChan)
// 	go totalThreadTime(pid, wg, totalThreadTimeChan)
// 	wg.Wait()

// 	// 得到总cpu差值
// 	// 该值一般为 cpu频率*1000
// 	deltaTotalCPUTime := -(<-totalCPUTimeChan - <-totalCPUTimeChan)
// 	if deltaTotalCPUTime == 0 {
// 		deltaTotalCPUTime = 1
// 	}
// 	// 得到总线程cpu差值
// 	deltaTotalThreadCPUTime := -(<-totalThreadTimeChan - <-totalThreadTimeChan)

// 	close(totalCPUTimeChan)
// 	close(totalThreadTimeChan)

// 	for range totalCPUTimeChan {
// 	}
// 	for range totalThreadTimeChan {
// 	}
// 	// 差值求商即为线程cpu使用率
// 	threadCPUUsage := (float64(deltaTotalThreadCPUTime) / float64(deltaTotalCPUTime)) * 100
// 	if threadCPUUsage == 0 {
// 		fmt.Printf("%.2f %%\n", threadCPUUsage)
// 	}
// 	fmt.Printf("%.2f %%\n", threadCPUUsage)
// 	return threadCPUUsage
// }

// // 总cpu时间
// func totalCPUTime(wg *sync.WaitGroup, ch chan int64) {
// 	var totalCPUTime int64 = 0
// 	content, err := ioutil.ReadFile("/proc/stat")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	b1 := bytes.Split(content, []byte("\n"))
// 	if len(b1) < 1 {
// 		log.Fatal("err")
// 	}
// 	res := bytes.Split(b1[0], []byte(" "))

// 	// https://man7.org/linux/man-pages/man5/proc.5.html
// 	// 查看 /proc/stat段落
// 	// splict后第[0]位是进程名，第[1]位为"空",故从第[2]为开始
// 	for i := 2; i < len(res); i++ {
// 		tmp, err := strconv.ParseInt(string(res[i]), 10, 64)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		totalCPUTime += tmp
// 	}
// 	ch <- totalCPUTime
// 	wg.Done()
// }

// // 线程cpu时间
// func totalThreadTime(pid string, wg *sync.WaitGroup, ch chan int64) {
// 	var totalThreadCPUTime int64 = 0
// 	content, err := ioutil.ReadFile("/proc/" + pid + "/stat")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	tmp := bytes.Split(content, []byte("\n"))
// 	if len(tmp) < 1 {
// 		log.Fatal("/proc/pid/stat file content error")
// 	}
// 	res := bytes.Split(tmp[0], []byte(" "))

// 	// https://man7.org/linux/man-pages/man5/proc.5.html
// 	// 查看 /proc/[pid]/stat段落
// 	// 第14位为utime, 15位为stime, 16位为cutime, 17位为cstime
// 	tmpRes := int64(0)
// 	// totalThreadCPUTime += res[14] + res[15] + res[16] + res[17]
// 	{
// 		tmpRes, _ = strconv.ParseInt(string(res[14]), 10, 64)
// 		totalThreadCPUTime += tmpRes
// 	}
// 	//
// 	{
// 		tmpRes, _ = strconv.ParseInt(string(res[15]), 10, 64)
// 		totalThreadCPUTime += tmpRes
// 	}
// 	{
// 		tmpRes, _ = strconv.ParseInt(string(res[16]), 10, 64)
// 		totalThreadCPUTime += tmpRes
// 	}
// 	{
// 		tmpRes, _ = strconv.ParseInt(string(res[17]), 10, 64)
// 		totalThreadCPUTime += tmpRes
// 	}
// 	ch <- totalThreadCPUTime
// 	wg.Done()
// }

// func load() [3]float64 {
// 	var err error
// 	content, err := ioutil.ReadFile("/proc/loadavg")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	tmp := bytes.Split(content, []byte(" "))
// 	l1, err := strconv.ParseFloat(string(tmp[0]), 64)
// 	l5, err := strconv.ParseFloat(string(tmp[1]), 64)
// 	l15, err := strconv.ParseFloat(string(tmp[2]), 64)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	return [3]float64{l1, l5, l15}
// }

// func getMemUsage(pid string) float64 {
// 	totalMem := func() int64 {
// 		f, err := os.Open("/proc/meminfo")
// 		defer f.Close()
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		buf := bufio.NewReader(f)
// 		l1, _, _ := buf.ReadLine()
// 		tmp := bytes.Fields(l1)[1]
// 		tmpRes, _ := strconv.ParseInt(string(tmp), 10, 64)
// 		return tmpRes
// 	}
// 	totalMem()
// 	content, err := ioutil.ReadFile("/proc/" + pid + "/statm")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	tmp := bytes.Split(content, []byte(" "))
// 	tmpRes, _ := strconv.ParseInt(string(tmp[1]), 10, 64)

// 	usage := (float64(tmpRes) * 4 / float64(totalMem())) * 100
// 	return usage
// }
