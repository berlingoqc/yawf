package config

import (
	"time"
)

// ITask if my main interface for all my kind of Task
// periodic task and channel task for now
type ITask interface {
	GetName() string
	Enable(map[string]interface{})
	Disable()

	Launch()
}

// TaskPool is the struct that hold all my task and started them when
// they need
type TaskPool struct {
	Tasks map[string]ITask
}

// AddPeriodicTask add a new periodic task that will be execute
// task need to be enable later on with the data
func (t *TaskPool) AddPeriodicTask(name string, frequency time.Duration, f func(s chan *Signal, m map[string]interface{})) {
	t.Tasks[name] = &PeriodicTask{
		Name:      name,
		Frequency: frequency,
		Task:      f,
	}
}

// LaunchNeededTask start the periodic task that need to
func (t *TaskPool) LaunchNeededTask() {
	for _, v := range t.Tasks {
		if pt, ok := v.(*PeriodicTask); ok {
			if pt.ReadyToLaunch() {
				v.Launch()
			}
		}
	}
}
