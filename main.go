package main

import (
	"encoding/xml"
	"fmt"
	"os"
	"time"
)

type HipReport struct {
	XMLName          xml.Name  `xml:"hip-report"`
	GenerateTime     string    `xml:"generate-time"`
	Version          int       `xml:"hip-report-version"`
	Categories       Categories `xml:"categories"`
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
	Entries []NetworkEntry `xml:"entry"`
}

type NetworkEntry struct {
	Name         string     `xml:"name,attr"`
	Description  string     `xml:"description"`
	MacAddress   string     `xml:"mac-address"`
	IPAddress    *IPAddresses `xml:"ip-address,omitempty"`
	IPv6Address  *IPAddresses `xml:"ipv6-address,omitempty"`
}

type IPAddresses struct {
	Entries []IPEntry `xml:"entry"`
}

type IPEntry struct {
	Name string `xml:"name,attr"`
}

type List struct {
	Entries []ListEntry `xml:"entry"`
}

type ListEntry struct {
	ProductInfo ProductInfo `xml:"ProductInfo"`
}

type ProductInfo struct {
	Prod               Prod       `xml:"Prod"`
	RealTimeProtection string     `xml:"real-time-protection,omitempty"`
	LastFullScanTime   string     `xml:"last-full-scan-time,omitempty"`
	Drives             *Drives    `xml:"drives,omitempty"`
	IsEnabled          string     `xml:"is-enabled,omitempty"`
}

type Prod struct {
	Vendor  string `xml:"vendor,attr"`
	Name    string `xml:"name,attr"`
	Version string `xml:"version,attr"`
	DefVer  string `xml:"defver,attr,omitempty"`
	EngVer  string `xml:"engver,attr,omitempty"`
	DateMon string `xml:"datemon,attr,omitempty"`
	DateDay string `xml:"dateday,attr,omitempty"`
	DateYear string `xml:"dateyear,attr,omitempty"`
	ProdType string `xml:"prodType,attr,omitempty"`
	OSType  string `xml:"osType,attr,omitempty"`
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
	hipReport := HipReport{
		GenerateTime: time.Now().Format("01/02/2006 15:04:05"),
		Version:      4,
		Categories: Categories{
			Entries: []CategoryEntry{
				{
					Name:          "host-info",
					OSVendor:      "Linux",
					HostName:      "coco",
					NetworkInterfaces: &NetworkInterfaces{
						Entries: []NetworkEntry{
							{
								Name:        "enp0s31f6",
								Description: "enp0s31f6",
								MacAddress:  "6c:4b:90:5b:7b:b3",
								IPAddress: &IPAddresses{
									Entries: []IPEntry{
										{Name: "192.168.1.144"},
									},
								},
								IPv6Address: &IPAddresses{
									Entries: []IPEntry{
										{Name: "fe80::d874:bc7c:9ca5:b800"},
									},
								},
							},
							{
								Name:        "wlp2s0",
								Description: "wlp2s0",
								MacAddress:  "0c:54:15:0b:4e:2e",
							},
							{
								Name:        "docker0",
								Description: "docker0",
								MacAddress:  "02:42:f7:89:d6:48",
								IPAddress: &IPAddresses{
									Entries: []IPEntry{
										{Name: "172.17.0.1"},
									},
								},
							},
						},
					},
				},
				{
					Name: "anti-malware",
					List: &List{
						Entries: []ListEntry{
							{
								ProductInfo: ProductInfo{
									Prod: Prod{
										Vendor:  "Cisco Systems, Inc.",
										Name:    "ClamAV",
										Version: "0.103.11",
										DefVer:  "27279",
										DateMon: "5",
										DateDay: "15",
										DateYear: "2024",
										ProdType: "3",
										OSType: "1",
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
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}

}

