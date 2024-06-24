package main

import (
	"errors"
	systemd "github.com/bechampion/gohip/systemd"
)

func RunPreflightChecks() error {
	clamavError := systemd.DbTooOldDefault()

	err := errors.Join(clamavError, clamavError)

	return err
}
