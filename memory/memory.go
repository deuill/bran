package memory

import (
	// Standard library
	"bufio"
	"flag"
	"os"
	"strconv"
	"strings"
	"time"

	// Internal packages
	"github.com/deuill/bran/statusbar"
)

// Memory represents an applet containing information about current memory usage.
type Memory struct {
	Interval *int    // The time interval between updates, in seconds.
	Icon     *string // The memory icon.

	msg  statusbar.Message
	info *os.File
}

// Run returns a message containing the current memory usage.
func (m *Memory) Run() *statusbar.Message {
	var total, free, buffers, cached int
	var seen int

	defer m.info.Seek(0, 0)
	scanner := bufio.NewScanner(m.info)

	for scanner.Scan() {
		// Stop scanning if all required fields have been processed.
		if seen == 4 {
			break
		}

		fields := strings.Fields(scanner.Text())
		if len(fields) < 2 {
			continue
		}

		// Match specific fields to values.
		num, _ := strconv.Atoi(fields[1])

		switch fields[0] {
		case "MemTotal:":
			total = num / 1024
		case "MemFree:":
			free = num / 1024
		case "Buffers:":
			buffers = num / 1024
		case "Cached:":
			cached = num / 1024
		default:
			continue
		}

		seen++
	}

	if scanner.Err() != nil {
		return nil
	}

	var used = total - free - buffers - cached
	var perc = int((float64(used) / float64(total)) * 100)

	m.msg.Text = *m.Icon + " " + strconv.Itoa(perc) + "%"
	return &m.msg
}

// Wait sleeps for a configurable amount of seconds.
func (m *Memory) Wait() {
	time.Sleep(time.Duration(*m.Interval) * time.Second)
}

// Init processes post-registration operations.
func (m *Memory) Init() error {
	var err error

	m.info, err = os.Open("/proc/meminfo")
	return err
}

// New returns a new instance of the memory applet.
func New() *statusbar.Applet {
	var flags flag.FlagSet
	mem := &Memory{
		Interval: flags.Int("interval", 5, ""),
		Icon:     flags.String("icon", "", ""),
	}

	applet := statusbar.NewApplet("memory", mem)
	applet.Flags(&flags)

	return applet
}
