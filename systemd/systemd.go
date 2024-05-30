package systemd

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
	"time"
	ctypes "types"

	"github.com/godbus/dbus/v5"
)

type Unit struct {
	Name        string
	Description string
	LoadState   string
	ActiveState string
	SubState    string
	Followed    dbus.ObjectPath
	Path        dbus.ObjectPath
	JobId       uint32
	JobType     string
	JobPath     dbus.ObjectPath
}

// Entries: []ctypes.ListEntry{
//      {
//          ProductInfo: ctypes.ProductInfo{
//              Prod: ctypes.Prod{
//                  Vendor:   "Cisco Systems, Inc.",
//                  Name:     "ClamAV",
//                  Version:  "0.103.11",
//                  DefVer:   "27279",
//                  DateMon:  "5",
//                  DateDay:  "15",
//                  DateYear: "2024",
//                  ProdType: "3",
//                  OSType:   "1",
//              },
//              RealTimeProtection: "no",
//              LastFullScanTime:   "n/a",
//          },
//      },
// },

// This would either find ClamAV or Crowdstrike for now at least
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
	cd := ClamDetails{}
	cd.Version = strings.Split(strings.Split(out.String(), " ")[1], "/")[0]
	cd.Defver = strings.Split(strings.Split(out.String(), " ")[1], "/")[1]
	m,_ := time.Parse("January",strings.Split(out.String(), " ")[2])
    // fmt.Printf("dada -> %d\n",int(m.Month()))
    cd.Month = fmt.Sprintf("%d",int(m.Month()))
	cd.Day = strings.Split(out.String(), " ")[3]
	cd.Year = strings.Split(out.String(), " ")[5][:4]
	return cd, nil
}
func FindAVUnit() (ctypes.Prod) {
	conn, err := dbus.SystemBus()
	if err != nil {
		panic("Failed to connect to system bus")
	}

	obj := conn.Object("org.freedesktop.systemd1", "/org/freedesktop/systemd1")
	var units []Unit
	err = obj.Call("org.freedesktop.systemd1.Manager.ListUnits", 0).Store(&units)
	if err != nil {
		fmt.Printf("Failed to list units: %v\n", err)
		return ctypes.Prod{}
	}

	for _, unit := range units {
		if strings.Contains(unit.Name, "clam") {
			if unit.ActiveState == "active" && unit.SubState == "running" {
				cd, err := GetClamDetails()
				if err != nil {
					return ctypes.Prod{}
				}
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
	}
	return ctypes.Prod{}
}
