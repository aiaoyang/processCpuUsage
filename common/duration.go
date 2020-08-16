package common

import "time"

// 统计cpu使用量时间间隔
const cpuUsageDuration = time.Millisecond * 100
const memUsageDuration = time.Millisecond * 100

// 统计机器负载时间间隔
const loadDuration = time.Second * 10
