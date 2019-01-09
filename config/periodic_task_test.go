package config

import (
	"testing"
	"time"
)

func TestPeriodicTask(t *testing.T) {
	i := 0
	p := &PeriodicTask{
		Frequency: 5 * time.Second,
		Task: func(c chan *Signal, args ...interface{}) {
			i += 1
		},
	}

	p.Enable()

	timeChan := time.NewTimer(20 * time.Second)

	for {
		select {
		case <-timeChan.C:
			if i == 0 {
				t.Fatal("Task was never executed")
			}
			return
		default:
			p.LaunchIfReady()
		}
	}

}
