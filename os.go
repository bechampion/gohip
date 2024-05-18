package main

import (
    "fmt"
    "log"
    "os"
    "strings"

    "github.com/shirou/gopsutil/v3/host"
    "github.com/shirou/gopsutil/v3/net"
)

func getHostname() (string, error) {
    hostname, err := os.Hostname()
    if err != nil {
        return "", err
    }
    return hostname, nil
}

func getOSVersion() (string, error) {
    info, err := host.Info()
    if err != nil {
        return "", err
    }
    return fmt.Sprintf("%s %s", info.Platform, info.PlatformVersion), nil
}

func getIPAddressAndMAC() ([]string, []string, error) {
    interfaces, err := net.Interfaces()
    if err != nil {
        return nil, nil, err
    }

    var ips []string
    var macs []string
    for _, iface := range interfaces {
        for _, addr := range iface.Addrs {
            if strings.Contains(addr.Addr, ":") {
                // Skip IPv6 addresses for simplicity
                continue
            }
            ips = append(ips, addr.Addr)
        }
        macs = append(macs, iface.HardwareAddr)
    }
    return ips, macs, nil
}

func main() {
    hostname, err := getHostname()
    if err != nil {
        log.Fatalf("Error getting hostname: %v", err)
    }

    osVersion, err := getOSVersion()
    if err != nil {
        log.Fatalf("Error getting OS version: %v", err)
    }

    ips, macs, err := getIPAddressAndMAC()
    if err != nil {
        log.Fatalf("Error getting IP and MAC addresses: %v", err)
    }

    fmt.Printf("Hostname: %s\n", hostname)
    fmt.Printf("OS Version: %s\n", osVersion)
    fmt.Printf("IP Addresses: %v\n", ips)
    fmt.Printf("MAC Addresses: %v\n", macs)
}

