package others

import (
	"bufio"
	"fmt"
	ctypes "github.com/bechampion/gohip/types"
	"os"
	"os/exec"
	"os/user"
	"strings"
)

func Others() {
	fmt.Println("hello")
}

func GetPackageManager() []ctypes.ListEntry {
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
		"Debian":   "apt",
		"Ubuntu":   "apt",
		"Centos":   "yum",
		"Fedora":   "dnf",
		"RedHat":   "yum",
		"Arch":     "pacman",
		"Manjaro":  "pacman",
		"Opensuse": "zypper",
		"Suse":     "zypper",
	}

	if pkgManager, found := packageManagers[distro]; found {
		prod := ctypes.ListEntry{
			ProductInfo: ctypes.ProductInfo{
				Prod: ctypes.Prod{
					Vendor:  distro,
					Name:    pkgManager,
					Version: "1",
				},
				IsEnabled: "yes",
			},
		}
		listpkgs = append(listpkgs, prod)
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
			Prod: ctypes.Prod{
				Vendor:  vendor,
				Name:    fw,
				Version: "1",
			},
			IsEnabled: "yes",
		},
	}
	listfw = append(listfw, prod)
	return listfw
}
func GetEncryptedPartitions() []ctypes.DriveEntry {
drives := []ctypes.DriveEntry{}
file, err := os.Open("/proc/mounts")
    if err != nil {
        return []ctypes.DriveEntry{}
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        line := scanner.Text()
        fields := strings.Fields(line)
        if len(fields) < 2 {
            continue
        }
        mountPoint := fields[1]
        if mountPoint == "/" {
				drives = append(drives, ctypes.DriveEntry{
					DriveName: fields[0],
					EncState:  "unncrypted",
				})
        }
    }

    if err := scanner.Err(); err != nil {
        return []ctypes.DriveEntry{}
    }
	return drives

}
func GetUserHomeDir() (string, error) {
	currentUser, err := user.Current()
	if err != nil {
		return "", err
	}
	return currentUser.HomeDir, nil
}
