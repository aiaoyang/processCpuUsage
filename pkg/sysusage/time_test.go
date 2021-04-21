package sysusage

import (
	"fmt"
	"testing"
	"time"
)

func Test_time_unix(t *testing.T) {
	t1 := time.Now().Unix()
	t2 := time.Now().UnixNano()
	e := 1e9
	fmt.Printf("t1: %d\nt2: %d\ne: %v\n", t1, t2, e/1000)
}
