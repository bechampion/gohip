package types

import (
	"encoding/xml"
)

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
	Entries []NetworkEntry `xml:"entry"`
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
