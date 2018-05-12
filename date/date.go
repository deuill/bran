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

	wait time.Duration
	msg  statusbar.Message
}

// Run returns a message containing the current date, according to the
// pre-configured format.
func (d *Date) Run() *statusbar.Message {
	d.msg.Text = *d.Icon + " " + time.Now().Format(*d.Format)
	return &d.msg
}

// Wait sleeps until the date minute changes.
func (d *Date) Wait() {
	switch d.wait {
	case time.Second:
		time.Sleep(1 * time.Second)
	case time.Minute:
		time.Sleep(time.Duration((60 - time.Now().Second())) * time.Second)
	}
}

// Init processes post-registration operations.
func (d *Date) Init() error {
	t, err := time.Parse(*d.Format, *d.Format)
	if err != nil {
		return err
	}

	// Override default wait interval if format precision is less than the default
	// of one minute.
	if t.Second() > 0 {
		d.wait = time.Second
	}

	return nil
}

// New returns a new instance of the date applet.
func New() *statusbar.Applet {
	var flags flag.FlagSet
	date := &Date{
		Format: flags.String("format", "Mon 2 Jan, 15:04", ""),
		Icon:   flags.String("icon", "", ""),
		wait:   time.Minute,
	}

	applet := statusbar.NewApplet("date", date)
	applet.Flags(&flags)

	return applet
}
