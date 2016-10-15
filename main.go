package main

import (
	// Standard library
	"fmt"
	"os"
	"os/signal"

	// Internal packages
	"github.com/deuill/granola/applet"

	// Third-party packages
	"github.com/rakyll/globalconf"

	// Internal applets
	_ "github.com/deuill/granola/applet/date"
)

func listen(ln chan *applet.Segment) {
	// Listen for and terminate Mash on SIGKILL or SIGINT signals.
	kill := make(chan os.Signal)
	signal.Notify(kill, os.Interrupt, os.Kill)

	// Print JSON header.
	fmt.Print(`{"version": 1, "click_events": true}` + "\n[\n")

	for {
		select {
		case seg := <-ln:
			fmt.Printf("[%s],\n", seg)
		case <-kill:
			return
		}
	}
}

func main() {
	// Initialize configuration, reading from environment variables using a
	// 'GRANOLA_' prefix first, then moving to a static configuration file,
	// usually located in '~/.config/granola/config.ini'.
	conf, err := globalconf.New("granola")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading configuration: %s\n", err)
		os.Exit(1)
	}

	conf.EnvPrefix = "GRANOLA_"
	conf.ParseAll()

	// Initialize requested applets and return message listener.
	ln, err := applet.Init(os.Args[1:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing applets: %s\n", err)
		os.Exit(2)
	}

	listen(ln)
	os.Exit(0)
}
