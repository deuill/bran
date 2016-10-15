package cpu

import (
	// Standard library
	"bufio"
	"flag"
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
	prev *stats
}

// Type stats represents a parsed snapshot from the /proc/stat file.
type stats struct {
	idle   int
	active int
}

// Run returns a message containing the current CPU load and temperature.
func (c *CPU) Run() *applet.Message {
	var usage int

	// Get current CPU usage.
	now := c.stats()

	// Calculate idle and total deltas.
	var idle = now.idle - c.prev.idle
	var total = (now.active + now.idle) - (c.prev.active + c.prev.idle)

	if total != 0 {
		usage = 100 * (total - idle) / total
	}

	c.msg.Text = *c.IconCPU + " " + strconv.Itoa(usage) + "%"
	c.prev = now

	// Get temperature for CPU.
	if temp := c.temp(); temp > -274 {
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

// Function stats returns the CPU stats, as given in /proc/stat.
func (c *CPU) stats() *stats {
	stat, err := os.Open("/proc/stat")
	if err != nil {
		return nil
	}

	defer stat.Close()

	// Read 'cpu' line from /proc/stat
	scanner := bufio.NewScanner(stat)
	var s []int

	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) > 0 && fields[0] == "cpu" {
			for _, f := range fields {
				num, _ := strconv.Atoi(f)
				s = append(s, num)
			}

			break
		}
	}

	if err := scanner.Err(); err != nil || len(s) < 8 {
		return nil
	}

	// Calculate idle and active times for CPU stats.
	return &stats{
		idle:   s[3] + s[4],
		active: s[0] + s[1] + s[2] + s[5] + s[6] + s[7],
	}
}

// Function temp returns the CPU temperature, in Celsius, from the 'thermal_zone'
// subsystem.
func (c *CPU) temp() int {
	temp, err := os.Open("/sys/class/thermal/thermal_zone0/temp")
	if err != nil {
		return -274
	}

	defer temp.Close()

	var buf = make([]byte, 32)
	count, _ := temp.Read(buf)

	num, err := strconv.Atoi(strings.TrimSpace(string(buf[:count])))
	if err != nil {
		return -274
	}

	// The thermal_zone subsystem returns a number with thousand precision, we
	// return the closest integer part.
	return num / 1000
}

func init() {
	flags := flag.NewFlagSet("cpu", flag.ContinueOnError)
	cpu := &CPU{
		Interval: flags.Int("interval", 5, ""),
		Scale:    flags.String("scale", "C", ""),
		IconCPU:  flags.String("icon-cpu", "", ""),
		IconTemp: flags.String("icon-temp", "", ""),
	}

	// Calculate initial CPU usage stats for subsequent runs.
	cpu.prev = cpu.stats()
	if cpu.prev == nil {
		return
	}

	applet.Register("cpu", cpu, flags)
}
