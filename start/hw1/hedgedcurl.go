package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"

	flag "github.com/spf13/pflag"
)

type CliArguments struct {
	urls    []string
	timeout uint16
	help    bool
}

const (
	defaultTimeout       = 15
	timeoutArgErrMessage = "Incorrect usage of the timeout argument"
	helpMessage          = `Запрос к нескольким серверам, вернет первый полученный ответ
./hedgedcurl https://motherfuckingwebsite.com/ https://thebestmotherfucking.website/ https://belyaev.work`
)

func main() {
	cliArguments, err := processArguments()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if cliArguments.help {
		fmt.Printf("%s\n\n", helpMessage)
		flag.Usage()
		os.Exit(0)
	}
}

func makeQuery(url string, client http.Client, responseChannel chan *http.Response, context context.Context) {
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Printf("error raised when creating request: %v\n", err)
		return
	}

	response, err := client.Do(request.WithContext(context))
	if err != nil && !errorIsTimeout(err) {
		select {
		case <-context.Done():
		default:
			fmt.Printf("error raised when requesting from %v: %v\n", url, err)
		}
	} else {
		responseChannel <- response
	}
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

func errorIsTimeout(err error) bool {
	if err == nil {
		return false
	}

	errNetError, ok := err.(net.Error)
	return ok && errNetError.Timeout()
}
