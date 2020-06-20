package system_monitor

import (
	"fmt"
	"os/user"
	"testing"
	"time"

	"github.com/ghostdb/ghostdb-cache-node/utils"
)

func TestSnitchMonitor(t *testing.T) {
	sch := NewSnitchScheduler(int32(2))
	go StartSnitch(sch)
	time.Sleep(6 * time.Second)
	StopSnitch(sch)

	usr, _ := user.Current()
	configPath := usr.HomeDir

	utils.AssertEqual(t, utils.FileExists(configPath+SnitchLogFileName), true, "")
	utils.AssertEqual(t, utils.FileNotEmpty(configPath+SnitchLogFileName), true, "")
}

func TestGetSnitchMetrics(t *testing.T) {
	fmt.Println(GetSnitchMetrics())
}


