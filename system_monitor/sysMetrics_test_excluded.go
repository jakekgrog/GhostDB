package system_monitor

import (
	"fmt"
	"os/user"
	"testing"
	"time"

	"github.com/ghostdb/ghostdb-cache-node/utils"
)

func TestSysMetricsMonitor(t *testing.T) {
	sch := NewSysMetricsScheduler(int32(2))
	go StartSysMetrics(sch)
	time.Sleep(6 * time.Second)
	StopSysMetrics(sch)

	usr, _ := user.Current()
	configPath := usr.HomeDir

	utils.AssertEqual(t, utils.FileExists(configPath+SysMetricsLogFilename), true, "")
	utils.AssertEqual(t, utils.FileNotEmpty(configPath+SysMetricsLogFilename), true, "")
}

func TestGetSysMetrics(t *testing.T) {
	fmt.Println(GetSysMetrics())
}


