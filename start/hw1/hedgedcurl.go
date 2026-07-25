package main

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"time"

	flag "github.com/spf13/pflag"
)

type CliArguments struct {
	urls    []string
	timeout uint16
	help    bool
}

const (
	defaultTimeout       = 15
	timeoutErrorCode     = 228
	timeoutArgErrMessage = "Incorrect usage of the timeout argument"

	helpMessage = `Запрос к нескольким серверам, вернет первый полученный ответ
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

	responseChannel := make(chan *http.Response)
	context, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(cliArguments.timeout))
	// the docs state that calling CancelFunc after the 1st call, subsequent calls do nothing
	// So we don't have to handle the case, when we call cancel after receiving a response from responseChannel
	defer cancel()

	httpClient := http.Client{}

	for _, url := range cliArguments.urls {
		go makeRequest(url, httpClient, responseChannel, context)
	}

	select {
	case <-context.Done(): // means timeout ?
		fmt.Println("Timed out.")
		os.Exit(timeoutErrorCode)
	case response := <-responseChannel:
		defer response.Body.Close()
		cancel()

		fmt.Println("Response received:")
		outputResponse(response)
	}
}

func makeRequest(url string, client http.Client, responseChannel chan *http.Response, context context.Context) {
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

func outputResponse(response *http.Response) error {
	fmt.Println("Response received:")

	fmt.Printf("Status code: %d\n", response.StatusCode)

	fmt.Println("Headers:")
	for name, value := range response.Header {
		fmt.Printf("%s: %s", name, value)
	}

	if response.StatusCode != http.StatusOK {
		return nil
	}

	fmt.Println("Body:")
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("error raised when reading a response body: %v", err)
	}
	bodyString := string(body)
	fmt.Println(bodyString)

	return nil
}
