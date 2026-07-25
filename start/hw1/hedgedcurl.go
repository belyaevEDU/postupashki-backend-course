package main

import (
	"fmt"

	flag "github.com/spf13/pflag"
)

type CliArguments struct {
	urls    []string
	timeout uint16
	help    bool
}

const (
	defaultTimeout       = 16
	timeoutArgErrMessage = "Incorrect usage of the timeout argument"
)

func main() {

}

func processArguments() (*CliArguments, error) {
	var timeoutInt int
	var help bool

	flag.IntVarP(&timeoutInt, "timeout", "t", defaultTimeout, "HTTP request timeout in seconds")
	flag.BoolVarP(&help, "help", "h", false, "Show help text")

	flag.Parse()

	if timeoutInt <= 0 {
		return nil, fmt.Errorf("error raised while processing given CLI arguments: %v", timeoutArgErrMessage)
	}
	timeout := uint16(timeoutInt)

	return &CliArguments{
		flag.Args(),
		timeout,
		help,
	}, nil
}
