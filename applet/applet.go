package applet

import (
	// Standard library
	"encoding/json"
	"flag"
	"fmt"

	// third-party packages
	"github.com/rakyll/globalconf"
)

// A Message represents an update for a statusbar segment.
type Message struct {
	Text  string `json:"full_text"`
	Short string `json:"short_text,omitempty"`
}

// A Segment represents a statusbar segment, or a specific position in a statubar.
type Segment struct {
	*Message
	Name string `json:"name"`
}

// String returns the JSON representation of a Segment.
func (s *Segment) String() string {
	buf, err := json.Marshal(*s)
	if err != nil {
		return ""
	}

	return string(buf)
}

// Applet represents a self-contained process which sends messages to the
// underlying statusbar for rendering.
type Applet interface {
	Init() error
	Run() *Message
	Wait()
}

// A map of registered and available applets.
var available = make(map[string]Applet)

// Register makes an applet available under a specific name.
func Register(name string, applet Applet, flags *flag.FlagSet) error {
	if _, exists := available[name]; exists {
		return fmt.Errorf("applet '%s' already registered, refusing to overwrite", name)
	}

	available[name] = applet

	// Register configuration flags, if set.
	if flags != nil {
		globalconf.Register(name, flags)
	}

	return nil
}

// Init sets up requested applets in separate goroutines and passes any incoming
// messages to the channel returned.
func Init(applets []string) (chan *Segment, error) {
	var ln = make(chan *Segment, len(applets))

	// Check for invalid applet names.
	for _, name := range applets {
		if _, exists := available[name]; !exists {
			return nil, fmt.Errorf("requested applet '%s' does not exist", name)
		}
	}

	// Run each default applet action in an infinite loop. Applets are supposed
	// to handle waiting in their own actions, otherwise they risk taxing the CPU.
	for _, name := range applets {
		err := available[name].Init()
		if err != nil {
			return nil, err
		}

		go func(name string) {
			// Initialize segment for applet.
			var seg Segment
			seg.Name = name

			for {
				// Get message from applet.
				if msg := available[name].Run(); msg != nil {
					seg.Message = msg
					ln <- &seg
				}

				available[name].Wait()
			}
		}(name)
	}

	return ln, nil
}
