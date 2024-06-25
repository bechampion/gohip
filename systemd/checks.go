package systemd

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"
)

var clamavDbFile ClamavDbFile

func init() {
	clamavDbFile = ClamavDbFile{path: "/var/lib/clamav/daily.cld"}
}

type ClamavDbFile struct {
	path string
}

func DefaultDbAgeCheck() error {
	details, clamavError := GetClamConfDetails()

	if clamavError != nil {
		return clamavError
	}

	return DbConfigAgeCheck(details)
}

func DbFileAgeCheck(clamavDbFile ClamavDbFile) error {
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

func DbConfigAgeCheck(details ClamConfDetails) error {
	weekAgo := time.Now().Add(-time.Hour * 24 * 7)

	tooOld := details.DailyCld.Before(weekAgo)

	if tooOld {
		return errors.New(fmt.Sprintf("virus definition is more than 7 days old: %s", details.DailyCld.String()))
	}

	return nil
}

type ClamConfDetails struct {
	version  string
	sigs     string
	DailyCld time.Time
}

func GetClamConfDetails() (ClamConfDetails, error) {
	cmd := exec.Command("clamconf")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return ClamConfDetails{}, errors.New(fmt.Sprintf("aaa"))
	}

	const layout = "Mon Jan 02 15:04:05 2006"
	lines := strings.Split(out.String(), "\n")
	re := regexp.MustCompile(`^daily.cld: version (.*), sigs: (.*), built on (.*)`)

	for i := range lines {
		line := lines[i]
		finds := re.FindStringSubmatch(line)

		if len(finds) > 0 {
			cd := ClamConfDetails{}
			cd.DailyCld, _ = time.Parse(layout, finds[3])
			cd.version = finds[1]
			cd.sigs = finds[2]
			return cd, nil
		}
	}

	return ClamConfDetails{}, errors.New(fmt.Sprintf("bbb"))
}
