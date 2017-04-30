package main

import (
	// Standard library
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"regexp"
	"syscall"

	// Internal packages
	"github.com/deuill/granola/statusbar"

	// Statusbar applets
	"github.com/deuill/granola/cpu"
	"github.com/deuill/granola/date"
	"github.com/deuill/granola/memory"
	"github.com/deuill/granola/volume"
)

var registered = map[string]func() *statusbar.Applet{
	"cpu":    cpu.New,
	"date":   date.New,
	"memory": memory.New,
	"volume": volume.New,
}

var (
	appletDesc   = regexp.MustCompile("([[:alpha:]]+)(?::(.+))?")
	appletConfig = regexp.MustCompile("([[:alnum:]]+)=([[:graph:]]+)")
)

func setup(desc []string) ([]*statusbar.Applet, error) {
	var applets []*statusbar.Applet

	for i := range desc {
		// Parse command-line description for applet.
		cmd := appletDesc.FindStringSubmatch(desc[i])
		name, values := cmd[1], cmd[2]

		if _, ok := registered[name]; ok == false {
			return nil, fmt.Errorf("applet with name '%s' does not exist", name)
		}

		// Initialize applet configuration, if any.
		conf := make(map[string]string)
		for _, v := range appletConfig.FindAllStringSubmatch(values, -1) {
			conf[v[1]] = v[2]
		}

		applet := registered[name]()
		applet.Set(conf)

		applets = append(applets, applet)
	}

	return applets, nil
}

func main() {
	// Initialize applet list from command-line description.
	applets, err := setup(os.Args[1:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing statusbar: %s\n", err)
		os.Exit(1)
	}

	// Initialize statusbar with pre-declared applet definitions.
	bar, err := statusbar.New(applets...)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing statusbar: %s\n", err)
		os.Exit(1)
	}

	// Listen for and terminate Granola on SIGTERM or SIGINT signals.
	halt := make(chan os.Signal)
	signal.Notify(halt, syscall.SIGTERM, syscall.SIGINT)

	// Listen for incoming messages from applets.
	ln := bar.Listen()

	// Print JSON header.
	fmt.Print(`{"version": 1, "click_events": true}` + "\n[\n")

	for {
		select {
		case s := <-ln:
			buf, _ := json.Marshal(s)
			fmt.Println(string(buf) + ",")
		case <-halt:
			os.Exit(0)
		}
	}
}
