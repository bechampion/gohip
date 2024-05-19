package osdata

import (
    "fmt"
    "os"
    "strings"

    "github.com/shirou/gopsutil/v3/host"
    "github.com/shirou/gopsutil/v3/net"
)

type IPEntry struct {
	Name string `xml:"name,attr"`
}
type IPAddresses struct {
	Entries []IPEntry `xml:"entry"`
}
type NetworkEntry struct {
	Name         string     `xml:"name,attr"`
	Description  string     `xml:"description"`
	MacAddress   string     `xml:"mac-address"`
	IPAddress    IPAddresses `xml:"ip-address,omitempty"`
	IPv6Address  IPAddresses `xml:"ipv6-address,omitempty"`
}
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

func GetInterfaces() ([]NetworkEntry, error) {
    ifaces := []NetworkEntry{}
    interfaces, err := net.Interfaces()
    if err != nil {
        return nil, err
    }

    var i NetworkEntry
    var ips []IPEntry
    var ip IPEntry
    var ipss IPAddresses
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

