package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	flag "github.com/spf13/pflag"
)

type CliArguments struct {
	urls    []string
	timeout uint
	help    bool
}

const (
	defaultTimeout       = 15
	timeoutArgErrMessage = "Incorrect usage of the timeout argument"

	badCliArgumentsErrorCode = 2
	urlValidationErrorCode   = 3
	timeoutErrorCode         = 228

	helpMessage = `Запрос к нескольким серверам, вернет первый полученный ответ
./hedgedcurl https://motherfuckingwebsite.com/ https://thebestmotherfucking.website/ https://belyaev.work`
	invalidUrlMessage = "Один из URLов не валиден. Формат: https://example.com или http://example.com. Пути и query параметры разрешены."
	timedOutMessage   = "Timed out."

	timeoutFlagMessage = "Время таймаута запроса в секундах"
	helpFlagMessage    = "Вывод текста об использовании утилиты"
)

func main() {
	cliArguments, err := processArguments()
	if err != nil {
		fmt.Printf("error raised when parsing cli arguments: %s\n", err)
		os.Exit(badCliArgumentsErrorCode)
	}

	if cliArguments.help || len(cliArguments.urls) == 0 {
		fmt.Printf("%s\n\n", helpMessage)
		flag.Usage()
	}

	if cliArguments.help {
		os.Exit(0)
	} else if len(cliArguments.urls) == 0 {
		os.Exit(badCliArgumentsErrorCode)
	}

	for _, url := range cliArguments.urls {
		result, err := validateUrl(url)
		if err != nil {
			fmt.Printf("error raised when validating a url: %s\n", err)
			os.Exit(urlValidationErrorCode)
		}

		if !result {
			fmt.Printf("%s. URL: %s", invalidUrlMessage, url)
			os.Exit(urlValidationErrorCode)
		}
	}

	responseChannel := make(chan *http.Response)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(cliArguments.timeout))
	// the docs state that calling CancelFunc after the 1st call, subsequent calls do nothing
	// So we don't have to handle the case, when we call cancel after receiving a response from responseChannel
	defer cancel()

	httpClient := http.Client{}

	for _, url := range cliArguments.urls {
		go makeRequest(ctx, url, &httpClient, responseChannel)
	}

	select {
	case <-ctx.Done(): // means timeout ?
		fmt.Println(timedOutMessage)
		os.Exit(timeoutErrorCode)
	case response := <-responseChannel:
		defer func() {
			err := response.Body.Close()
			if err != nil {
				fmt.Printf("error raised when closing a response body: %s\n", err)
			}
		}()

		err = outputResponse(response)
		if err != nil {
			fmt.Println("error raised when outputting response: ", err)
		}
		cancel()
	}
}

func makeRequest(ctx context.Context, url string, client *http.Client, responseChannel chan *http.Response) {
	request, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		fmt.Printf("error raised when creating request: %s\n", err)
		return
	}

	response, err := client.Do(request)
	if err != nil {
		if ctx.Err() == nil && !errorIsTimeout(err) {
			fmt.Printf("error raised when requesting from %s: %s\n", url, err)
		}
		return
	}

	select {
	case <-ctx.Done():
		err := response.Body.Close()
		if err != nil {
			fmt.Printf("error raised when closing a response body: %s\n", err)
		}
	case responseChannel <- response:
	}
}

func processArguments() (*CliArguments, error) {
	var timeoutInt int
	var help bool

	flag.IntVarP(&timeoutInt, "timeout", "t", defaultTimeout, timeoutFlagMessage)
	flag.BoolVarP(&help, "help", "h", false, helpFlagMessage)

	flag.Parse()

	if timeoutInt <= 0 {
		return nil, fmt.Errorf("error raised while processing given CLI arguments: %s", timeoutArgErrMessage)
	}
	timeout := uint(timeoutInt)

	return &CliArguments{
		flag.Args(),
		timeout,
		help,
	}, nil
}

func outputResponse(response *http.Response) error {
	fmt.Printf("Response received:\n\n")

	fmt.Printf("Status code: %d\n\n", response.StatusCode)

	fmt.Printf("Headers:\n")
	for name, value := range response.Header {
		fmt.Printf("%s: %s\n", name, value)
	}

	if response.StatusCode != http.StatusOK {
		return nil
	}

	fmt.Println("\nBody:")
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("error raised when reading a response body: %w", err)
	}
	bodyString := string(body)
	fmt.Println(bodyString)

	return nil
}
