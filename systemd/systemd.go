package systemd


import (
    "fmt"
    "strings"
    "os/exec"
    "bytes"

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
func GetClamDetails() string {
    cmd := exec.Command("clamd", "-V")
    var out bytes.Buffer
    cmd.Stdout = &out
    err := cmd.Run()
    if err != nil {
        return "", err
    }
    return out.String(), nil
    return "version"
}
func FindAVUnit() {
    conn, err := dbus.SystemBus()
    if err != nil {
        panic("Failed to connect to system bus")
    }

    obj := conn.Object("org.freedesktop.systemd1", "/org/freedesktop/systemd1")
    var units []Unit
    err = obj.Call("org.freedesktop.systemd1.Manager.ListUnits", 0).Store(&units)
    if err != nil {
        fmt.Printf("Failed to list units: %v\n", err)
        return
    }

    for _, unit := range units {
        if strings.Contains(unit.Name, "clam") {
            if unit.ActiveState == "active" && unit.SubState == "running" {
                fmt.Printf("Name: %s, Description: %s, LoadState: %s, ActiveState: %s, SubState: %s\n",
                    unit.Name, unit.Description, unit.LoadState, unit.ActiveState, unit.SubState)
            }
        }
    }
}
