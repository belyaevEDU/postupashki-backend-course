package main

import (
	"fmt"
	"net"
	"regexp"
)

const urlRegex string = `^(http:\/\/www\.|https:\/\/www\.|http:\/\/|https:\/\/|\/|\/\/)?[A-z0-9_-]*?[:]?[A-z0-9_-]*?[@]?[A-z0-9]+([\-\.]{1}[a-z0-9]+)*\.[a-z]{2,5}(:[0-9]{1,5})?(\/.*)?$`

func errorIsTimeout(err error) bool {
	if err == nil {
		return false
	}

	errNetError, ok := err.(net.Error)
	return ok && errNetError.Timeout()
}

func validateUrl(url string) (bool, error) {
	regex, err := regexp.Compile(urlRegex)
	if err != nil {
		return false, fmt.Errorf("error raised when compiling regex: %w", err)
	}

	return regex.MatchString(url), nil
}
