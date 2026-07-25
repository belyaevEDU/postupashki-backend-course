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
	timeoutInt := *flag.IntP("timeout", "t", defaultTimeout, "HTTP request timeout in seconds")
	if timeoutInt <= 0 {
		return nil, fmt.Errorf("error raised while processing given CLI arguments: %v", timeoutArgErrMessage)
	}
	timeout := uint16(timeoutInt)
	help := *flag.BoolP("help", "h", false, "Show help text")

	return &CliArguments{
		flag.Args(),
		timeout,
		help,
	}, nil
}
