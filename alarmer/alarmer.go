package alarmer

import (
	"context"
	"fmt"
	"time"
)

// Reciver 告警信息接收者
type Reciver interface {
	Recive(ctx context.Context)
}

// Sender 消息发送者
type Sender interface {
	Send(msg string)
}

// AlarmActor 告警动作 主接口
type AlarmActor interface {
	Reciver
	Sender
	Alarm(string)
	Recover(string)
	IsAlarming() bool
}

// AlarmType 枚举体 告警类型
type AlarmType int

const (
	CPU  AlarmType = iota // 0
	MEM                   // 1
	FD                    // 2
	NET                   // 3
	DISK                  // 4
)

// Alarmer 告警结构体
type Alarmer struct {
	// mu *sync.Mutex

	alarmChan chan struct{}

	recoverChan chan struct{}

	// 沉默通道周期定时器
	timer *time.Timer

	// 是否正在发生告警
	isAlarming bool

	// 沉默通道周期时间间隔
	AlarmSilenceDuration time.Duration

	// 收到的告警消息
	AlarmMsg string
}

// NewAlarmer 初始化一个alarmer
func NewAlarmer(duration time.Duration) *Alarmer {
	return &Alarmer{
		// mu: &sync.Mutex{},

		alarmChan:   make(chan struct{}, 0),
		recoverChan: make(chan struct{}, 0),

		timer: time.NewTimer(duration),

		AlarmSilenceDuration: duration,
	}
}

// IsAlarming 是否发生告警
func (a *Alarmer) IsAlarming() bool {
	res := a.isAlarming
	return res
}

// Recive 接收告警信息
func (a *Alarmer) Recive(ctx context.Context) {
	fmt.Printf("timer stop result : %v\n", a.timer.Stop())
	for {
		select {

		case <-ctx.Done():
			return

		case <-a.timer.C:
			// 沉默通道周期时间结束，开启告警通知
			a.isAlarming = false

		// 告警通道收到信息，如果信息为告警发生则进行告警，如果信息为告警恢复 则设置告警状态为未发生告警
		case <-a.alarmChan:
			if !a.isAlarming {
				a.alarm()
			}
			time.Sleep(time.Second)
		case <-a.recoverChan:
			a.recover()

		default:
			time.Sleep(time.Second)
		}
	}
}

func (a *Alarmer) resetTimer() {
	a.timer.Stop()
	a.timer.Reset(a.AlarmSilenceDuration)
}

// Alarm Alarm
func (a *Alarmer) Alarm(msg string) {
	a.AlarmMsg = msg
	a.alarmChan <- struct{}{}
}
func (a *Alarmer) alarm() {
	a.isAlarming = true
	a.resetTimer()
	a.Send(a.AlarmMsg)
}

// Recover Recover
func (a *Alarmer) Recover(msg string) {
	if !a.isAlarming {
		return
	}
	a.AlarmMsg = msg
	a.recoverChan <- struct{}{}
}
func (a *Alarmer) recover() {
	a.isAlarming = false
	a.Send(a.AlarmMsg)
}

// Send 发送告警信息
func (a *Alarmer) Send(msg string) {
	fmt.Printf("%v\n", msg)
}
