package systemd

import (
	"bytes"
	"fmt"
	ctypes "github.com/bechampion/gohip/types"
	"log"
	"os"
	"os/exec"
	"regexp"
	"slices"
	"strings"
	"time"
)

var (
	WarningLogger *log.Logger
)

func init() {
	WarningLogger = log.New(os.Stderr, "[WARN] ", log.Ldate|log.Ltime|log.Lshortfile)
}

type ClamDetails struct {
	Version string
	Month   string
	Day     string
	Year    string
	Defver  string
}

func GetClamDetails() (ClamDetails, error) {
	cmd := exec.Command("clamd", "-V")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return ClamDetails{}, err
	}

	// Example output: "ClamAV 0.103.6/26590/Mon Jul  1 10:40:06 2024"
	// Normalize multiple spaces to single space for consistent parsing
	re := regexp.MustCompile(`\s{2,}`)
	cleanout := strings.TrimSpace(re.ReplaceAllString(out.String(), " "))

	// Parse with regex instead of fragile index-based splits
	// Matches: "ClamAV <version>/<defver>/<date>"
	parseRe := regexp.MustCompile(`ClamAV\s+(\S+)/(\d+)/(.+)`)
	matches := parseRe.FindStringSubmatch(cleanout)
	if len(matches) < 4 {
		return ClamDetails{}, fmt.Errorf("unexpected clamd -V output format: %q", cleanout)
	}

	version := matches[1]
	defver := matches[2]
	dateStr := strings.TrimSpace(matches[3])

	// Parse the date portion (e.g. "Mon Jul 1 10:40:06 2024")
	t, err := time.Parse("Mon Jan 2 15:04:05 2006", dateStr)
	if err != nil {
		return ClamDetails{}, fmt.Errorf("failed to parse date from clamd output %q: %w", dateStr, err)
	}

	cd := ClamDetails{
		Version: version,
		Defver:  defver,
		Month:   fmt.Sprintf("%d", int(t.Month())),
		Day:     fmt.Sprintf("%d", t.Day()),
		Year:    fmt.Sprintf("%d", t.Year()),
	}
	return cd, nil
}

func FindClamdProcess() ctypes.Prod {
	cmd := exec.Command("ps", "aux")
	var out bytes.Buffer
	cmd.Stdout = &out

	noClamAV := ctypes.Prod{}

	if err := cmd.Run(); err != nil {
		return noClamAV
	}

	lines := strings.Split(out.String(), "\n")

	isRunning := slices.ContainsFunc(lines, func(process string) bool {
		return strings.Contains(process, "clamd")
	})

	if isRunning {
		if cd, err := GetClamDetails(); err == nil {
			return ctypes.Prod{
				Vendor:   "Cisco Systems, Inc.",
				Name:     "ClamAV",
				Version:  cd.Version,
				DateMon:  cd.Month,
				DateDay:  cd.Day,
				DateYear: cd.Year,
				DefVer:   cd.Defver,
				ProdType: "3",
				OSType:   "1",
			}
		}
	}

	WarningLogger.Println("clamd process not found")
	return noClamAV
}
