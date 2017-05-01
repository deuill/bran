# Bran - A zero-configuration statusbar for i3

Bran is a statusbar for i3, intended to be simple to set up and use by providing a strong set of defaults, while allowing for customization through the same command-line interface that's used for invoking itself.

[![API Documentation][godoc-svg]][godoc-url] [![MIT License][license-svg]][license-url]

## Building and Installing

Bran is built in Go, and requires the Go toolchain to be [set up][go-setup] before building and installing.

Assuming everything is set up correctly, installation is simply a matter of running `go get`, i.e.:

```
go get -u -v github.com/deuill/bran
```

You should now have a binary named `bran` in your `$GOPATH/bin` directory.

## Usage

Bran contains a small number of core "applets", which provide functionality such as date/time display, volume status via ALSA/Pulseaudio, CPU and memory usage etc.

Both placement and customization of these is defined through the command-line interface for `bran`, for example:

```
bar {
    status_command bran cpu memory date
}
```

This initializes `bran` with three applets, CPU usage, memory usage and date/time, displayed left-to-right as entered. Applets can be customized by providing additional arguments after the applet name:

```
bar {
    status_command bran cpu:"interval=10 scale=F" memory volume:step=10 date
}
```

This sets the update interval and temperature scale for the `cpu` applet to 10 seconds and Fahrenheit respectively, and sets the volume step for the `volume` applet to 10%.

Multiple configuration values need to be wrapped in quotes, due to the way shells handle space-separated strings.

## License

All code in this repository is covered by the terms of the MIT License, the full text of which can be found in the LICENSE file.

[godoc-url]: https://godoc.org/github.com/deuill/bran
[godoc-svg]: https://godoc.org/github.com/deuill/bran?status.svg

[license-url]: https://github.com/deuill/bran/blob/master/LICENSE
[license-svg]: https://img.shields.io/badge/license-MIT-blue.svg

[go-setup]: https://golang.org/doc/install
