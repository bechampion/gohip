package osdata

import (
	"encoding/json"
	"fmt"
	"github.com/shirou/gopsutil/v3/host"
	ctypes "gohip/types"
	"os"
	"os/exec"
)

type MainInterface struct {
	Dst     string   `json:"dst"`
	Gateway string   `json:"gateway"`
	Dev     string   `json:"dev"`
	PrefSrc string   `json:"prefsrc"`
	Flags   []string `json:"flags"`
	UID     int      `json:"uid"`
	Cache   []string `json:"cache"`
}
type MacAddress struct {
	Address string `json:"address"`
}

func GetMac(dev string) (string, error) {
	cmd := exec.Command("ip", "-json", "link", "show", dev)
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	var address []MacAddress
	err = json.Unmarshal(output, &address)
	if err != nil {
		return "", err
	}
	return address[0].Address, nil
}
func GetInterfaces() ([]ctypes.NetworkEntry, error) {
	iface := ctypes.NetworkEntry{}
	cmd := exec.Command("ip", "-json", "r", "get", "1.1.1.1")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var routes []MainInterface
	err = json.Unmarshal(output, &routes)
	if err != nil {
		return nil, err
	}

	if len(routes) == 0 {
		return nil, fmt.Errorf("no route found for destination %s", "1.1.1.1")
	}

	var ipaddress ctypes.IPEntry
	ipaddress.Name = routes[0].Gateway
	iface.IPAddress.Entries = append(iface.IPAddress.Entries, ipaddress)
	iface.Name = routes[0].Dev
	// Let's get the mac
	mac, err := GetMac(iface.Name)
	if err != nil {
		return nil, err
	}
	iface.Description = "This is your default dev"
	iface.MacAddress = mac
	ifaces := []ctypes.NetworkEntry{}
	ifaces = append(ifaces, iface)
	return ifaces, nil
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
	return info.OS, nil
}
func GetOSVersion() (string, error) {
	info, err := host.Info()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s %s", info.Platform, info.PlatformVersion), nil
}
