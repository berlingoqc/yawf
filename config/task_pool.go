package config

import (
	"time"
)

type ITask interface {
	Enable(...interface{})
	Disable()

	ReadyToLaunch() bool
	Launch()
}

type TaskPool struct {
	Tasks map[string]ITask
}

func (t *TaskPool) AddPeriodicTask(name string, frequency time.Duration, f func(chan *Signal, ...interface{})) {
	t.Tasks[name] = &PeriodicTask{
		Name:      name,
		Frequency: frequency,
		Task:      f,
	}
}

func (t *TaskPool) LaunchNeededTask() {
	for _, v := range t.Tasks {
		if v.ReadyToLaunch() {
			v.Launch()
		}
	}
}
