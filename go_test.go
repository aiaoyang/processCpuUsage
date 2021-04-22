package main

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

type testConfig struct {
	mu   sync.Mutex
	data int
}

func Test_gochan(t *testing.T) {
	c := test_go()
	for i := 0; i < 30; i++ {
		c.mu.Lock()
		fmt.Printf("i: %d\n", c.data)
		c.mu.Unlock()
		time.Sleep(time.Nanosecond * 10)
	}
	time.Sleep(time.Second * 2)

	// for i := range ch {
	// 	fmt.Println(i)
	// }

}

func test_go() (c *testConfig) {
	c = &testConfig{}
	go func() {
		for i := 0; i < 30; i++ {
			c.mu.Lock()
			c.data++
			fmt.Printf("goroutine data: %d\n", c.data)
			c.mu.Unlock()
			time.Sleep(time.Nanosecond * 10)
		}
	}()
	return
}
