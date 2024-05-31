package others

import (
	"fmt"
	"bufio"
	"os"
	"strings"
	ctypes "types"
	"os/exec"
)

func Others() {
	fmt.Println("hello")
}

func GetPackageManager() ([]ctypes.ListEntry) {
	var listpkgs []ctypes.ListEntry
	file, _ := os.Open("/etc/os-release")
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var distro string
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "ID=") {
			distro = strings.TrimPrefix(line, "ID=")
			distro = strings.Trim(distro, "\"")
			break
		}
	}

	packageManagers := map[string]string{
		"Debian":  "apt",
		"Ubuntu":  "apt",
		"Centos":  "yum",
		"Fedora":  "dnf",
		"RedHat":    "yum",
		"Arch":    "pacman",
		"Manjaro": "pacman",
		"Opensuse": "zypper",
		"Suse":    "zypper",
	}

	if pkgManager, found := packageManagers[distro]; found {
		prod := ctypes.ListEntry{
		ProductInfo: ctypes.ProductInfo{
				Prod: ctypes.Prod {
					Vendor: distro,
						Name: pkgManager,
						Version: "1",
					},
				IsEnabled: "yes",
			},
		}
		listpkgs = append(listpkgs,prod)
		return listpkgs
	}

	return []ctypes.ListEntry{}
}
func GetFirewall() []ctypes.ListEntry {
	var listfw []ctypes.ListEntry
	fw := "none"
    vendor := "none"
	if _, err := exec.LookPath("ufw"); err == nil {
		fw = "ufw"
		vendor = "Canonical Ltd."
	}
	if _, err := exec.LookPath("iptables"); err == nil {
		fw = "iptables"
		vendor = "IPTables"
	}
		prod := ctypes.ListEntry{
		ProductInfo: ctypes.ProductInfo{
				Prod: ctypes.Prod {
					Vendor: vendor,
						Name: fw,
						Version: "1",
					},
				IsEnabled: "yes",
			},
		}
		listfw = append(listfw,prod)
		return listfw
}
