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
	Task      func(chan *Signal, ...interface{})

	enable bool
	signal chan *Signal
	ticker *time.Ticker
}

func (p *PeriodicTask) IsEnabled() bool {
	return p.enable
}

func (p *PeriodicTask) Enable(args ...interface{}) {
	p.reset()
	p.enable = true
}

func (p *PeriodicTask) Disable() {
	p.enable = false
	close(p.signal)
	p.ticker.Stop()
	p.ticker = nil
}

func (p *PeriodicTask) ReadyToLaunch() bool {
	select {
	case <-p.ticker.C:
		return true
	default:
		return false
	}
}

func (p *PeriodicTask) Launch() {
	go p.Task(p.signal)
}

func (p *PeriodicTask) LaunchIfReady() {
	if p.ReadyToLaunch() {
		go p.Task(p.signal)
	}

}

func (p *PeriodicTask) reset() {
	p.ticker = time.NewTicker(p.Frequency)
	p.signal = make(chan *Signal)
}
