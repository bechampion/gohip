package main

import (
	"errors"
	systemd "github.com/bechampion/gohip/systemd"
)

func RunPreflightChecks() error {
	clamavError := systemd.DefaultDbAgeCheck()

	return errors.Join(clamavError)
}
