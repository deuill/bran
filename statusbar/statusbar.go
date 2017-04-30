package statusbar

import (
	// Standard library
	"flag"
	"fmt"
)

// A Statusbar represents a collection of applets that can be polled for updates
// and collected as ordered segments.
type Statusbar struct {
	applets  []*Applet
	listener chan []*Segment
}

// Listen processes messages for all underlying applets and returns a channel
// for receiving updates in the form of segments, which can then be iterated or
// marshalled as needed.
func (s *Statusbar) Listen() chan []*Segment {
	var instances = make(map[string]int)
	var segments = make([]*Segment, len(s.applets))

	for i, applet := range s.applets {
		// Multiple applets of the same type are identified by their instance
		// number, starting from 1.
		instances[applet.name]++

		// Launch a persistent process for applet instance, alternating between
		// fetching messages and waiting for the next iteration. Messages are
		// collected and sent back to the listener as segments.
		go func(applet *Applet, i int) {
			var seg = &Segment{
				Name:     applet.name,
				Instance: fmt.Sprint(instances[applet.name]),
			}

			for {
				msg := applet.runner.Run()
				if msg != nil {
					// Update segment state.
					seg.Message = msg
					segments[i] = seg

					// Send list of segments to listener.
					s.listener <- segments
				}

				applet.runner.Wait()
			}
		}(applet, i)
	}

	return s.listener
}

// New initializes all applets provided and attaches them to a new Statusbar
// instance.
func New(applets ...*Applet) (*Statusbar, error) {
	var err error

	// Assign un-formatted values to each applet's flag-set (if any).
	for _, applet := range applets {
		if applet.flags == nil {
			continue
		}

		applet.flags.VisitAll(func(f *flag.Flag) {
			// Ignore subsequent flags if an error has occured.
			if err != nil {
				return
			}

			// Assign value to corresponding flag.
			v, ok := applet.values[f.Name]
			if ok == false {
				return
			}

			err = f.Value.Set(v)
			if err != nil {
				err = fmt.Errorf("%s.%s: %s", applet.name, f.Name, err)
			}
		})

		if err != nil {
			return nil, err
		}
	}

	// Perform any post-registration initialization routines.
	for _, applet := range applets {
		if err = applet.runner.Init(); err != nil {
			return nil, err
		}
	}

	return &Statusbar{
		applets:  applets,
		listener: make(chan []*Segment, len(applets)),
	}, nil
}
