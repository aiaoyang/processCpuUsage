package sysusage

import (
	"testing"
)

func Test_process_pid(t *testing.T) {
	// p := GetProcessPID("vscode")

	_ = GetAllProcess()
	// fmt.Printf("all Process: %+v\n", allProcess)
	// t.Log(p)
}

func Benchmark_getProcess(t *testing.B) {
	for i := 0; i < t.N; i++ {
		GetAllProcess()
	}
}
