package systemd

import (
    "bytes"
    "fmt"
    "os/exec"
    "regexp"
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
    re := regexp.MustCompile(`\s{2,}`)
	cleanout := re.ReplaceAllString(out.String(), " ")
    cd := ClamDetails{}
    cd.Version = strings.Split(strings.Split(cleanout, " ")[1], "/")[0]
    cd.Defver = strings.Split(strings.Split(cleanout, " ")[1], "/")[1]
    m,_ := time.Parse("Jan",strings.Split(cleanout, " ")[2])
    cd.Month = fmt.Sprintf("%d",int(m.Month()))
    cd.Day = strings.Split(cleanout, " ")[3]
    cd.Year = strings.Split(cleanout, " ")[5][:4]
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
