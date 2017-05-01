package date

import (
	// Standard library
	"flag"
	"time"

	// Internal packages
	"github.com/deuill/bran/statusbar"
)

// Date represents a message containing the current date.
type Date struct {
	Format *string // The date format to use.
	Icon   *string // Code-point for date icon.
	msg    statusbar.Message
}

// Run returns a message containing the current date, according to the
// pre-configured format.
func (d *Date) Run() *statusbar.Message {
	d.msg.Text = *d.Icon + " " + time.Now().Format(*d.Format)

	return &d.msg
}

// Wait sleeps until the date minute changes.
func (d *Date) Wait() {
	time.Sleep(time.Duration(60-time.Now().Second()) * time.Second)
}

// Init processes post-registration operations.
func (d *Date) Init() error {
	return nil
}

// New returns a new instance of the date applet.
func New() *statusbar.Applet {
	var flags flag.FlagSet
	date := &Date{
		Format: flags.String("format", "Mon 2 Jan, 15:04", ""),
		Icon:   flags.String("icon", "", ""),
	}

	applet := statusbar.NewApplet("date", date)
	applet.Flags(&flags)

	return applet
}
