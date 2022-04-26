package restutil

import (
	"errors"
	"strings"
)

const (
	ErrDownServer = "Server is down"
	ErrFailedAuth = "Failed to authenticate"
	ErrInvalidArg = "Invalid argument"
	ErrNotFound   = "Not found"
	ErrEmailTake  = "Email is already taken"
)

func CheckErr(i error) error {
	if any(i.Error(), "Error while dialing") {
		return errors.New(ErrDownServer)
	}
	if any(i.Error(), "signin failed") {
		return errors.New(ErrFailedAuth)
	}
	if any(i.Error(), "email already exists") {
		return errors.New(ErrEmailTake)
	}

	return errors.New(ErrNotFound)
}

func any(i string, what string) bool {
	if strings.ContainsAny(i, what) {
		return true
	}

	return false
}
