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
	//This is gnarly
	re := regexp.MustCompile(`\s{2,}`)
	cleanout := re.ReplaceAllString(out.String(), " ")
	cd := ClamDetails{}
	cd.Version = strings.Split(strings.Split(cleanout, " ")[1], "/")[0]
	cd.Defver = strings.Split(strings.Split(cleanout, " ")[1], "/")[1]
	m, _ := time.Parse("Jan", strings.Split(cleanout, " ")[2])
	cd.Month = fmt.Sprintf("%d", int(m.Month()))
	cd.Day = strings.Split(cleanout, " ")[3]
	cd.Year = strings.Split(cleanout, " ")[5][:4]
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
