package system_monitor

import (
	"os"
	"time"
)

// SysMetricsScheduler represents a scheduler for system metrics logger
type SysMetricsScheduler struct {
	Interval time.Duration
	stop     chan bool
}

// NewSysMetricsScheduler initializes a new SysMetrics Scheduler
func NewSysMetricsScheduler(interval int32) *SysMetricsScheduler {
	scheduler := &SysMetricsScheduler{
		Interval: time.Duration(interval) * time.Second,
		stop:     make(chan bool),
	}

	return scheduler
}

// StartSysMetrics starts the SysMetrics monitor
func StartSysMetrics(scheduler *SysMetricsScheduler) {
	configPath, _ := os.UserConfigDir()
	if _, err := os.Stat(configPath + "/ghostdb"); os.IsNotExist(err) {
		os.Mkdir(configPath+"/ghostdb", 0777)
	}

	ticker := time.NewTicker(scheduler.Interval)

	for {
		select {
		case <-ticker.C:
			go StartSysMetricsMonitor()
		case <-scheduler.stop:
			ticker.Stop()
			return
		}
	}
}

// StopSysMetrics stops the sys metrics scheduler by passing
// a bool to the scheduler stop channel.
func StopSysMetrics(scheduler *SysMetricsScheduler) {
	scheduler.stop <- true
}
