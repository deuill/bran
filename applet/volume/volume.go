package volume

// #cgo LDFLAGS: -lm -lasound
// #include <alsa/asoundlib.h>
// #include "monitor.h"
import "C"

import (
	// Standard library
	"flag"
	"fmt"

	// Internal packages
	"github.com/deuill/granola/applet"
)

// Volume represents an applet containing the current input/output volume levels.
type Volume struct {
	Icon *string // The volume icon.
	Card *int    // The sound card to monitor, as an index.

	msg applet.Message
}

// Run returns a message containing the current volume levels.
func (v *Volume) Run() *applet.Message {
	v.msg.Text = fmt.Sprintf("%s %d%%", *v.Icon, C.volume())
	return &v.msg
}

// Wait waits for volume level changes.
func (v *Volume) Wait() {
	C.wait()
}

// Init processes post-registration operations.
func (v *Volume) Init() error {
	return nil
}

func init() {
	var flags flag.FlagSet
	volume := &Volume{
		Icon: flags.String("icon", "", ""),
	}

	applet.Register("volume", volume, &flags)
}
