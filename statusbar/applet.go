package statusbar

import (
	// Standard library
	"flag"
)

// A Message represents an update for a statusbar segment.
type Message struct {
	Text  string `json:"full_text"`
	Short string `json:"short_text,omitempty"`
}

// A Segment represents a statusbar segment, or a specific position in a statubar.
type Segment struct {
	Name     string `json:"name"`
	Instance string `json:"instance"`
	*Message
}

// A Runner represents the underlying process that produces status messages
// for applets, to be consumed by the containing statusbar.
type Runner interface {
	Init() error
	Run() *Message
	Wait()
}

// Applet represents a self-contained process which sends messages to the
// underlying statusbar for rendering.
type Applet struct {
	name   string
	runner Runner

	flags  *flag.FlagSet
	values map[string]string
}

// Flags attaches the flagset provided to the applet instance.
func (a *Applet) Flags(flags *flag.FlagSet) {
	a.flags = flags
}

// Set provides a map of unformatted values to the applet, mainly for use against
// the applet's internal flag-set.
func (a *Applet) Set(values map[string]string) {
	a.values = values
}

// NewApplet initializes a concrete applet instance from a messenger, under a
// unique name.
func NewApplet(name string, runner Runner) *Applet {
	return &Applet{
		name:   name,
		runner: runner,
		values: make(map[string]string),
	}
}
