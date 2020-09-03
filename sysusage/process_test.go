package sysusage

import "testing"

func Test_process_pid(t *testing.T) {
	p := GetProcessPID("vscode")
	t.Log(p)
}
