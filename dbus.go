package main

import (
    "fmt"
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

func main() {
    // Connect to the system bus
    conn, err := dbus.SystemBus()
    if err != nil {
        fmt.Printf("Failed to connect to system bus: %v\n", err)
        return
    }

    // Get a reference to the systemd manager object
    obj := conn.Object("org.freedesktop.systemd1", "/org/freedesktop/systemd1")

    // Call the ListUnits method
    var units []Unit
    err = obj.Call("org.freedesktop.systemd1.Manager.ListUnits", 0).Store(&units)
    if err != nil {
        fmt.Printf("Failed to list units: %v\n", err)
        return
    }

    // Filter and print active services
    fmt.Println("Active services:")
    for _, unit := range units {
        if unit.ActiveState == "active" && unit.SubState == "running" {
            fmt.Printf("Name: %s, Description: %s, LoadState: %s, ActiveState: %s, SubState: %s\n",
                unit.Name, unit.Description, unit.LoadState, unit.ActiveState, unit.SubState)
        }
    }
}

