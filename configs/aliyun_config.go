package configs

import (
	"context"
	"sync"
)

type AcmInnerConfig struct {
	mu           sync.Mutex `yaml:"ommitte"`
	Hostname     string     `yaml:"hostName"`
	StoredPids   []string   `yaml:"storedPids"`
	ProcessRegex string     `yaml:"processRegex"`
	IsMonitorOn  bool       `yaml:"isMonitorOn"`
}

func (c *AcmInnerConfig) Reload(in *AcmInnerConfig) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.Hostname = in.Hostname
	c.IsMonitorOn = in.IsMonitorOn
	c.ProcessRegex = in.ProcessRegex
	c.StoredPids = in.StoredPids
}
func (c *AcmInnerConfig) Lock() {
	c.mu.Lock()
}
func (c *AcmInnerConfig) Unlock() {
	c.mu.Unlock()
}

func (c *AcmInnerConfig) Watch(ctx context.Context, ch <-chan *AcmInnerConfig) {
	for {
		select {
		case <-ctx.Done():
			return
		case in := <-ch:
			if in != nil {
				c.Reload(in)
			} else {
				continue
			}
		}
	}
}
