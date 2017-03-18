package cpu

import (
	// Standard library
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	// Internal packages
	"github.com/deuill/granola/applet"
)

// CPU represents an applet containing information about CPU usage and temperature.
type CPU struct {
	Interval *int    // The time interval between updates, in seconds.
	Scale    *string // The temperature scale to use. Choices are C, F.
	IconCPU  *string // The CPU icon.
	IconTemp *string // The temperature icon.

	msg  applet.Message
	prev *usage

	stat    *os.File
	thermal *os.File
}

// Type usage represents the current CPU usage metrics.
type usage struct {
	active int
	idle   int
}

// Run returns a message containing the current CPU load and temperature.
func (c *CPU) Run() *applet.Message {
	// Get current CPU usage.
	now := c.usage()

	// Calculate idle and total deltas.
	var idle = now.idle - c.prev.idle
	var total = (now.active + now.idle) - (c.prev.active + c.prev.idle)

	// Calculate usage percentage.
	var percent int

	if total != 0 {
		percent = 100 * (total - idle) / total
	}

	c.msg.Text = *c.IconCPU + " " + strconv.Itoa(percent) + "%"
	c.prev = now

	// Get temperature for CPU.
	if temp := c.temp(); temp > 0 {
		c.msg.Text += " " + *c.IconTemp + " "

		// Convert to different scale if required.
		switch *c.Scale {
		case "F": // Fahrenheit
			temp = int(float64(temp)*1.8) + 32
		}

		c.msg.Text += strconv.Itoa(temp) + "°" + *c.Scale
	}

	return &c.msg
}

// Wait sleeps for a configurable amount of seconds.
func (c *CPU) Wait() {
	time.Sleep(time.Duration(*c.Interval) * time.Second)
}

// Function usage returns the CPU usage statistics, as given in /proc/stat.
func (c *CPU) usage() *usage {
	var stat = make([]int, 8)
	var buf = make([]byte, 128)

	// Read 'cpu' line from /proc/stat
	count, _ := c.stat.Read(buf)
	c.stat.Seek(0, 0)

	fields := strings.Fields(string(buf[:count]))
	if len(fields) < len(stat)+1 || fields[0] != "cpu" {
		return nil
	}

	for i := 0; i < len(stat); i++ {
		stat[i], _ = strconv.Atoi(fields[i+1])

	}

	// Calculate idle and active times for CPU stats.
	return &usage{
		active: stat[0] + stat[1] + stat[2] + stat[5] + stat[6] + stat[7],
		idle:   stat[3] + stat[4],
	}
}

// Function temp returns the current CPU temperature, in Celsius.
func (c *CPU) temp() int {
	var buf = make([]byte, 32)
	count, _ := c.thermal.Read(buf)
	c.thermal.Seek(0, 0)

	num, err := strconv.Atoi(strings.TrimSpace(string(buf[:count])))
	if err != nil {
		return 0
	}

	// The thermal_zone subsystem returns a number with milligrade precision, we
	// return the closest integer part.
	return num / 1000
}

// Init processes post-registration operations.
func (c *CPU) Init() error {
	var err error

	c.stat, err = os.Open("/proc/stat")
	if err != nil {
		return err
	}

	c.thermal, err = os.Open("/sys/class/thermal/thermal_zone0/temp")
	if err != nil {
		return err
	}

	// Calculate initial CPU usage stats for subsequent runs.
	c.prev = c.usage()
	if c.prev == nil {
		return fmt.Errorf("Could not fetch initial CPU usage stats")
	}

	return nil
}

func init() {
	var flags flag.FlagSet
	cpu := &CPU{
		Interval: flags.Int("interval", 5, ""),
		Scale:    flags.String("scale", "C", ""),
		IconCPU:  flags.String("icon-cpu", "", ""),
		IconTemp: flags.String("icon-temp", "", ""),
	}

	applet.Register("cpu", cpu, &flags)
}
