package main

import (
	"fmt"
	"testing"
	"time"
)

func Test_gochan(t *testing.T) {
	ch := test_go()

	for i := range ch {
		fmt.Println(i)
	}

}

func test_go() (ch chan int) {
	ch = make(chan int)
	go func() {
		for i := 0; i < 3; i++ {
			ch <- 1
			time.Sleep(time.Second)
		}
		close(ch)
	}()
	return
}
