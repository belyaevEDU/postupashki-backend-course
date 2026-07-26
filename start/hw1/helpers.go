package main

import (
	"net"
)

func errorIsTimeout(err error) bool {
	if err == nil {
		return false
	}

	errNetError, ok := err.(net.Error)
	return ok && errNetError.Timeout()
}
