package main

import (
	"encoding/xml"
	"fmt"
	"net/url"
	"os"
	"osdata"
	"time"
	"flag"
	"strings"
	"log"
)

func logCommandAndArgs() {
    // Get the command and arguments
    command := os.Args[0]
    args := strings.Join(os.Args[1:], " ")

    // Create or open the log file
    file, err := os.OpenFile("command.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        log.Fatalf("Failed to open log file: %v", err)
    }
    defer file.Close()

    // Create a logger
    logger := log.New(file, "", log.LstdFlags)

    // Log the command and arguments
    logger.Printf("Command: %s Arguments: %s\n", command, args)
}
type HipReport struct {
	XMLName      xml.Name   `xml:"hip-report"`
	Name         string     `xml:"name,attr"`
	Md5          string     `xml:"md5-sum"`
	User         string     `xml:"user-name"`
	Domain       string     `xml:"domain"`
	HostName     string     `xml:"host-name"`
	HostId       string     `xml:"host-id"`
	Ip           string     `xml:"ip-address"`
	Ip6          string     `xml:"ipv6-address"`
	GenerateTime string     `xml:"generate-time"`
	Version      int        `xml:"hip-report-version"`
	Categories   Categories `xml:"categories"`
}

type Categories struct {
	Entries []CategoryEntry `xml:"entry"`
}

type CategoryEntry struct {
	Name              string             `xml:"name,attr"`
	ClientVersion     string             `xml:"client-version,omitempty"`
	OS                string             `xml:"os,omitempty"`
	OSVendor          string             `xml:"os-vendor,omitempty"`
	Domain            string             `xml:"domain,omitempty"`
	HostName          string             `xml:"host-name,omitempty"`
	HostID            string             `xml:"host-id,omitempty"`
	NetworkInterfaces *NetworkInterfaces `xml:"network-interface,omitempty"`
	List              *List              `xml:"list,omitempty"`
	MissingPatches    *MissingPatches    `xml:"missing-patches,omitempty"`
}

type NetworkInterfaces struct {
	Entries []osdata.NetworkEntry `xml:"entry"`
}

type List struct {
	Entries []ListEntry `xml:"entry"`
}

type ListEntry struct {
	ProductInfo ProductInfo `xml:"ProductInfo"`
}

type ProductInfo struct {
	Prod               Prod    `xml:"Prod"`
	RealTimeProtection string  `xml:"real-time-protection,omitempty"`
	LastFullScanTime   string  `xml:"last-full-scan-time,omitempty"`
	Drives             *Drives `xml:"drives,omitempty"`
	IsEnabled          string  `xml:"is-enabled,omitempty"`
}

type Prod struct {
	Vendor   string `xml:"vendor,attr"`
	Name     string `xml:"name,attr"`
	Version  string `xml:"version,attr"`
	DefVer   string `xml:"defver,attr,omitempty"`
	EngVer   string `xml:"engver,attr,omitempty"`
	DateMon  string `xml:"datemon,attr,omitempty"`
	DateDay  string `xml:"dateday,attr,omitempty"`
	DateYear string `xml:"dateyear,attr,omitempty"`
	ProdType string `xml:"prodType,attr,omitempty"`
	OSType   string `xml:"osType,attr,omitempty"`
}

type Drives struct {
	Entries []DriveEntry `xml:"entry"`
}

type DriveEntry struct {
	DriveName string `xml:"drive-name"`
	EncState  string `xml:"enc-state"`
}

type MissingPatches struct {
	Entries []MissingPatchEntry `xml:"entry,omitempty"`
}

type MissingPatchEntry struct {
	// Add fields if necessary
}

func main() {
	logCommandAndArgs()
	cookie := flag.String("cookie", "", "")
	 _ = flag.String("client-os", "", "")
	//--client-ip seems to be fed from openconect but i don't think it's used
	 _ = flag.String("client-ip", "", "")
	md5 := flag.String("md5", "", "")
    flag.Parse()
	values, err := url.ParseQuery(*cookie)
	if err != nil {
		panic(err)
	}
	user := values.Get("user")
	domain := values.Get("domain")
	// Computer doesn't seem to be used , i leave it here for future reference
	 _ = values.Get("computer")

	hostname, err := osdata.GetHostname()
	if err != nil {
		panic(err)
	}
	osname, err := osdata.GetOS()
	if err != nil {
		panic(err)
	}
	interfaces, err := osdata.GetInterfaces()
	if err != nil {
		panic(err)
	}
	hipReport := HipReport{
		Name:         "hip-report",
		GenerateTime: time.Now().Format("01/02/2006 15:04:05"),
		Version:      4,
		User: user,
		HostId: hostname,
		Md5: *md5,
		Domain: domain,
		Categories: Categories{
			Entries: []CategoryEntry{
				{
					Name:     "host-info",
					OSVendor: osname,
					HostName: hostname,
					NetworkInterfaces: &NetworkInterfaces{
						Entries: interfaces,
					}},
				{
					Name: "anti-malware",
					List: &List{
						Entries: []ListEntry{
							{
								ProductInfo: ProductInfo{
									Prod: Prod{
										Vendor:   "Cisco Systems, Inc.",
										Name:     "ClamAV",
										Version:  "0.103.11",
										DefVer:   "27279",
										DateMon:  "5",
										DateDay:  "15",
										DateYear: "2024",
										ProdType: "3",
										OSType:   "1",
									},
									RealTimeProtection: "no",
									LastFullScanTime:   "n/a",
								},
							},
						},
					},
				},
				{
					Name: "disk-backup",
					List: &List{
						Entries: []ListEntry{},
					},
				},
				{
					Name: "disk-encryption",
					List: &List{
						Entries: []ListEntry{
							{
								ProductInfo: ProductInfo{
									Prod: Prod{
										Vendor:  "GitLab Inc.",
										Name:    "cryptsetup",
										Version: "2.4.3",
									},
									Drives: &Drives{
										Entries: []DriveEntry{
											{
												DriveName: "/",
												EncState:  "unencrypted",
											},
											{
												DriveName: "/home/jgarcia",
												EncState:  "unencrypted",
											},
											{
												DriveName: "All",
												EncState:  "unencrypted",
											},
										},
									},
								},
							},
						},
					},
				},
				{
					Name: "firewall",
					List: &List{
						Entries: []ListEntry{
							{
								ProductInfo: ProductInfo{
									Prod: Prod{
										Vendor:  "Canonical Ltd.",
										Name:    "UFW",
										Version: "0.36.1",
									},
									IsEnabled: "no",
								},
							},
							{
								ProductInfo: ProductInfo{
									Prod: Prod{
										Vendor:  "IPTables",
										Name:    "IPTables",
										Version: "1.8.7",
									},
									IsEnabled: "no",
								},
							},
							{
								ProductInfo: ProductInfo{
									Prod: Prod{
										Vendor:  "The Netfilter Project",
										Name:    "nftables",
										Version: "1.0.2",
									},
									IsEnabled: "no",
								},
							},
						},
					},
				},
				{
					Name: "patch-management",
					List: &List{
						Entries: []ListEntry{
							{
								ProductInfo: ProductInfo{
									Prod: Prod{
										Vendor:  "Canonical Ltd.",
										Name:    "Snap",
										Version: "22.04.1",
									},
									IsEnabled: "n/a",
								},
							},
							{
								ProductInfo: ProductInfo{
									Prod: Prod{
										Vendor:  "GNU",
										Name:    "Advanced Packaging Tool",
										Version: "2.4.11",
									},
									IsEnabled: "yes",
								},
							},
						},
					},
					MissingPatches: &MissingPatches{
						Entries: []MissingPatchEntry{},
					},
				},
				{
					Name: "data-loss-prevention",
					List: &List{
						Entries: []ListEntry{},
					},
				},
			},
		},
	}

	xmlData, err := xml.MarshalIndent(hipReport, "", "        ")
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}

	xmlHeader := []byte(xml.Header)
	xmlData = append(xmlHeader, xmlData...)

	err = os.WriteFile("/dev/stdout", xmlData, 0644)
	os.Exit(0)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}

}
