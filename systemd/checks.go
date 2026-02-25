package systemd

import (
	"bytes"
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
		return fmt.Errorf("virus definition is too old: %s is more than %d hours old (> 7 days)", clamavDbFile.path, hoursSince)
	} else {
		return nil
	}
}

func DbConfigAgeCheck(details ClamConfDetails) error {
	weekAgo := time.Now().Add(-time.Hour * 24 * 7)

	tooOld := details.DailyCld.Before(weekAgo)

	if tooOld {
		return fmt.Errorf("virus definition is more than 7 days old: %s", details.DailyCld.String())
	}

	return nil
}

type ClamConfDetails struct {
	version  string
	sigs     string
	DailyCld time.Time
}

func parseDate(src string) (time.Time, error) {
	const layout = "Mon Jan 2 15:04:05 2006"
	return time.Parse(layout, src)
}

func parseDailyCvdLine(line string) []string {
	re := regexp.MustCompile(`^daily.c[lv]d: version (.*), sigs: (.*), built on (.*)`)
	return re.FindStringSubmatch(line)
}

func GetClamConfDetails() (ClamConfDetails, error) {
	cmd := exec.Command("clamconf")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return ClamConfDetails{}, fmt.Errorf("failed to run clamconf: %w", err)
	}

	lines := strings.Split(out.String(), "\n")

	for i := range lines {
		line := lines[i]
		finds := parseDailyCvdLine(line)

		if len(finds) > 0 {
			cd := ClamConfDetails{}
			parsedDate, parseErr := parseDate(finds[3])
			if parseErr != nil {
				return ClamConfDetails{}, fmt.Errorf("failed to parse date %q from clamconf output: %w", finds[3], parseErr)
			}
			cd.DailyCld = parsedDate
			cd.version = finds[1]
			cd.sigs = finds[2]
			return cd, nil
		}
	}

	return ClamConfDetails{}, fmt.Errorf("could not determine timestamp for daily.cld in clamconf output")
}
