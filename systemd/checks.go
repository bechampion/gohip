package systemd

import (
	"errors"
	"fmt"
	"os"
	"time"
)

var (
	clamavDbFile ClamavDbFile
)

func init() {
	//clamavDbFile = ClamavDbFile{path: "/etc/ssh"}
	clamavDbFile = ClamavDbFile{path: "/var/lib/clamav/daily.cld"}
}

type ClamavDbFile struct {
	path string
}

func DbTooOldDefault() error {
	return DbTooOld(clamavDbFile)
}

func DbTooOld(clamavDbFile ClamavDbFile) error {
	hoursInWeek := 24 * 7

	fi, err := os.Stat(clamavDbFile.path)

	if err != nil {
		return err
	}

	mtime := fi.ModTime()

	hoursSince := int(time.Since(mtime).Hours())

	if hoursSince > hoursInWeek {
		return errors.New(fmt.Sprintf("virus definition is too old: %s is more than %d hours old (> 7 days)", clamavDbFile.path, hoursSince))
	} else {
		return nil
	}
}
