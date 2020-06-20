package system_monitor

import (
	"os"
	"time"
)

// SnitchScheduler represents a scheduler for snitch system metrics logger
type SnitchScheduler struct {
	Interval time.Duration
	stop     chan bool
}

// NewSnitchScheduler initializes a new Snitch Scheduler
func NewSnitchScheduler(interval int32) *SnitchScheduler {
	scheduler := &SnitchScheduler{
		Interval: time.Duration(interval) * time.Second,
		stop:     make(chan bool),
	}

	return scheduler
}

// StartSnitch starts the Snitch monitor
func StartSnitch(scheduler *SnitchScheduler) {
	configPath, _ := os.UserConfigDir()
	if _, err := os.Stat(configPath + "/ghostdb"); os.IsNotExist(err) {
		os.Mkdir(configPath+"/ghostdb", 0777)
	}

	ticker := time.NewTicker(scheduler.Interval)

	for {
		select {
		case <-ticker.C:
			go StartSnitchMonitor()
		case <-scheduler.stop:
			ticker.Stop()
			return
		}
	}
}

// StopSnitch stops the snitch scheduler by passing
// a bool to the scheduler stop channel.
func StopSnitch(scheduler *SnitchScheduler) {
	scheduler.stop <- true
}
