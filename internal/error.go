package internal

import "errors"

var (
	MissingAuthenticationError = errors.New("missing authentication token")
	InvalidAuthenticationError = errors.New("invalid authentication token")
)