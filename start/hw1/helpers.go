package main

import (
	"net"
	"net/url"
)

func errorIsTimeout(err error) bool {
	if err == nil {
		return false
	}

	errNetError, ok := err.(net.Error)
	return ok && errNetError.Timeout()
}

func validateUrl(urlString string) (bool, error) {
	_, err := url.ParseRequestURI(urlString)
	// didnt really see an opportunity to use the URI down the line.
	// url.ParseRequestURI assumes that the string provided is from an HTTP request,
	// while url.Parse accepts relative URLs

	if err != nil {
		return false, err // wrapping the error down the line already
	}

	return true, nil
}
