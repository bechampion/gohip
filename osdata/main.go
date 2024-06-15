package osdata

import (
    "fmt"
    "os"
    "strings"
	ctypes "github.com/bechampion/gohip/types"
    "github.com/shirou/gopsutil/v3/host"
    "github.com/shirou/gopsutil/v3/net"
)

func GetHostname() (string, error) {
    hostname, err := os.Hostname()
    if err != nil {
        return "", err
    }
    return hostname, nil
}

func GetOS() (string, error) {
    info, err := host.Info()
    if err != nil {
        return "", err
    }
    return info.OS,nil
}
func GetOSVersion() (string, error) {
    info, err := host.Info()
    if err != nil {
        return "", err
    }
    return fmt.Sprintf("%s %s", info.Platform, info.PlatformVersion), nil
}

func GetInterfaces() ([]ctypes.NetworkEntry, error) {
    ifaces := []ctypes.NetworkEntry{}
    interfaces, err := net.Interfaces()
    if err != nil {
        return nil, err
    }

    var i ctypes.NetworkEntry
    var ips []ctypes.IPEntry
    var ip ctypes.IPEntry
    var ipss ctypes.IPAddresses
    for _, iface := range interfaces {
	i.Name = iface.Name
        for _, addr := range iface.Addrs {
            if strings.Contains(addr.Addr, ":") {
                // Skip IPv6 addresses for simplicity
                continue
            }
	    ip.Name = addr.Addr
            ips = append(ips, ip)
        }
	ipss.Entries = ips
	i.IPAddress = ipss
        i.MacAddress= iface.HardwareAddr
    ifaces = append(ifaces,i)
    }


    return ifaces, nil
}

