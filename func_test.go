package main

import (
	"fmt"
	"sync"
	"testing"
)

func Benchmark_test(b *testing.B) {
	a := sync.Once{}
	a.Do(func() { fmt.Println("first") })
	a.Do(func() { fmt.Println("second") })
	fmt.Println(&a == nil)
}
