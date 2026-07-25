package main

import (
	"fmt"
	"net"
	"regexp"
)

const urlRegex string = `/^https?:\/\/(?:www\.)?[-a-zA-Z0-9@:%._\+~#=]{1,256}\.[a-zA-Z0-9()]{1,6}\b(?:[-a-zA-Z0-9()@:%_\+.~#?&\/=]*)$/`

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
