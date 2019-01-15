package config

import (
	"time"
)

// Signal est la stucture utilisé pour communiquer entre mes tasks et ma main thread
type Signal interface {
	IsOver() bool
	GetError() error
}

// PeriodicTask represente une tache qui est executer periodiquement
type PeriodicTask struct {
	Name      string
	Frequency time.Duration
	Task      func(chan *Signal, map[string]interface{})

	enable bool
	signal chan *Signal
	ticker *time.Ticker
}

// GetName ...
func (p *PeriodicTask) GetName() string {
	return p.Name
}

// IsEnabled tell if this task is enable
func (p *PeriodicTask) IsEnabled() bool {
	return p.enable
}

// Enable start to schedule the task
func (p *PeriodicTask) Enable(args map[string]interface{}) {
	p.reset()
	p.enable = true
}

// Disable stop the current task
func (p *PeriodicTask) Disable() {
	p.enable = false
	close(p.signal)
	p.ticker.Stop()
	p.ticker = nil
}

// ReadyToLaunch tell if this task need to be execute right now
func (p *PeriodicTask) ReadyToLaunch() bool {
	select {
	case <-p.ticker.C:
		return true
	default:
		return false
	}
}

// Launch start the task one time
func (p *PeriodicTask) Launch() {
	go p.Task(p.signal, nil)
}

// LaunchIfReady start the task if it is ready to start
func (p *PeriodicTask) LaunchIfReady() {
	if p.ReadyToLaunch() {
		go p.Task(p.signal, nil)
	}

}

func (p *PeriodicTask) reset() {
	p.ticker = time.NewTicker(p.Frequency)
	p.signal = make(chan *Signal)
}
