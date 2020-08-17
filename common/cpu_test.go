package common

import (
	"context"
	"fmt"
	"log"
	"math"
	"runtime"
	"sync"
	"testing"
	"time"
)

func Test_initMap(t *testing.T) {
	m := make(map[int]CPUStat)
	fmt.Printf("M: %v\n", m[1])
}
func Benchmark_CPUUsage(b *testing.B) {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		time.Sleep(time.Second * 5)
		cancel()
	}()
	reciver := make(chan map[string]float64, 0)
	pidsChan := make(chan []int, 0)

	go func() {
		for {
			pidsChan <- []int{46787}

			// 修改传入pid的时间会反映在测试结果时间中
			time.Sleep(time.Millisecond * 500)
		}
	}()
	go func() {
		for {
			select {
			case rec := <-reciver:
				fmt.Printf("rec: %v\n", rec)
			default:
				continue
			}
		}
	}()

	// =========================================================
	cpuCoreNum := runtime.NumCPU()

	process := &ProcessCPUStat{}
	system := &CPUStat{}

	oldPorcess := &ProcessCPUStat{}
	oldSystem := &CPUStat{}

	var wg = &sync.WaitGroup{}

	for i := 0; i < b.N; i++ {

		metric := make(map[string]float64)

		metric["core_num"] = float64(cpuCoreNum)

		select {

		case <-ctx.Done():
			log.Println("stopping")
			return
		case pids := <-pidsChan:

			for _, pid := range pids {

				wg.Add(2)

				// 计算系统cpu总使用率
				// currSystemCPUCount := &CPUStat{}
				go func(wg *sync.WaitGroup) {
					system.New(sysPid)

					wg.Done()
				}(wg)

				//开始计算当前pid进程cpu统计量以及与上一次cpu统计量的变化量
				// currProcessCPUCount := &ProcessCPUStat{}
				go func(wg *sync.WaitGroup) {
					process.New(pid)

					wg.Done()
				}(wg)

				wg.Wait()

				metric["sys_cpu_usage"] = float64(system.used-oldSystem.used) * 100 / float64(system.total-oldSystem.total+1)

				// 进程不存在，则只返回系统cpu使用情况
				if process.isDead {
					continue
				}

				metric["pid"] = float64(pid)

				// 计算本次进程cpu和上次进程cpu差异，并与总cpu变化求商得使用率
				//										进程cpu使用量-上次进程cpu使用量									系统cpu使用量-上次系统cpu使用量						系统cpu核数
				metric["process_cpu_usage"] = math.Max(0, float64(process.used-oldPorcess.used)*100/(float64(process.time-oldPorcess.time+1)*HZ*float64(cpuCoreNum)))

				// 将本次进程cpu统计结果设置为旧
				*oldPorcess = *process

				// 将本次系统cpu总量果设置为旧
				*oldSystem = *system

			}
			reciver <- metric
		default:
			time.Sleep(time.Millisecond * 100)
		}
	}
	// =========================================================

	// select {}
}

type cupStat struct {
	used  uint64
	total uint64
}

func cpuUsageTtest(pid string) float64 {
	total1 := totalCPU()
	process1 := processCPU(pid)
	time.Sleep(time.Millisecond * 10)
	total2 := totalCPU()
	process2 := processCPU(pid)
	return float64(process2-process1) * 100 / float64(total2.total-total1.total+1)
}
func processCPU(pid string) uint64 {
	f := "/proc/" + pid + "/stat"
	l, err := HeadLineSplitOfFile(f)
	if err != nil {
		log.Fatal(err)
	}
	return byteToUint(l[13]) + byteToUint(l[14]) + byteToUint(l[15]) + byteToUint(l[16])
}
func percent() float64 {
	t1 := totalCPU()
	time.Sleep(1)
	t2 := totalCPU()
	return float64(t2.used-t1.used) * 100 / float64(t2.total-t1.total)
}
func totalCPU() cupStat {
	f := "/proc/stat"
	l, err := HeadLineSplitOfFile(f)
	if err != nil {
		log.Fatal(err)
	}

	total := func() uint64 {
		// 除数 不能为0
		var tmp uint64 = 1
		for i := 1; i < len(l)-2; i++ {
			tmp += byteToUint(l[i])
		}
		return tmp
	}()

	return cupStat{
		total: total,
		used:  total - byteToUint(l[4]),
	}
}
