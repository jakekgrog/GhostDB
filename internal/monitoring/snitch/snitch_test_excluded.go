package snitch

import (
	"fmt"
	"os"
	"os/user"
	"testing"
	"time"

	"github.com/ghostdb/ghostdb-cache-node/cache/utils"
)

func TestSnitchMonitor(t *testing.T) {
	sch := NewSnitchScheduler(int32(2))
	go StartSnitch(sch)
	time.Sleep(6 * time.Second)
	StopSnitch(sch)

	usr, _ := user.Current()
	configPath := usr.HomeDir

	utils.AssertEqual(t, fileExists(configPath+SnitchLogFileName), true, "")
	utils.AssertEqual(t, fileNotEmpty(configPath+SnitchLogFileName), true, "")
}

func TestGetSnitchMetrics(t *testing.T) {
	fmt.Println(GetSnitchMetrics())
}

func fileExists(filename string) bool {
	file, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}

	return !file.IsDir()
}

func fileNotEmpty(filename string) bool {
	file, err := os.Stat(filename)
	if err != nil {
		return false
	}

	size := file.Size()
	if size > 0 {
		return true
	}
	return false
}
