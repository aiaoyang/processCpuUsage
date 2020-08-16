package main

import (
	"testing"
	"unsafe"
)

func Benchmark_test(b *testing.B) {
	for i := 0; i < b.N; i++ {
		getCPUUsage("1")
	}
}

func Benchmark_ParseInt(b *testing.B) {
	bt := []byte("12345678")
	for i := 0; i < b.N; i++ {
		_ = *(*string)(unsafe.Pointer(&bt))
		// fmt.Println(str)

		// _, err := strconv.ParseFloat(string(bt), 64)
		// if err != nil {
		// 	log.Fatal(err)
		// }
	}
}
