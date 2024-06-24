package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	osdata "gohip/osdata"
	others "gohip/others"
	systemd "gohip/systemd"
	ctypes "gohip/types"
	"log"
	"net/url"
	"os"
	"strings"
	"time"
)

func logCommandAndArgs() {
	command := os.Args[0]
	args := strings.Join(os.Args[1:], " ")
	file, err := os.OpenFile("command.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	defer file.Close()
	logger := log.New(file, "", log.LstdFlags)
	logger.Printf("Command: %s Arguments: %s\n", command, args)
}

func main() {
	logCommandAndArgs()
	cookie := flag.String("cookie", "", "")
	_ = flag.String("client-os", "", "")
	clientip := flag.String("client-ip", "", "")
	md5 := flag.String("md5", "", "")
	flag.Parse()
	values, err := url.ParseQuery(*cookie)
	if err != nil {
		panic(err)
	}
	user := values.Get("user")
	domain := values.Get("domain")
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
	hipReport := ctypes.HipReport{
		Name:         "hip-report",
		GenerateTime: time.Now().Format("01/02/2006 15:04:05"),
		Version:      4,
		User:         user,
		HostName:     hostname,
		HostId:       hostname,
		Md5:          *md5,
		Ip:           *clientip,
		Domain:       domain,
		Categories: ctypes.Categories{
			Entries: []ctypes.CategoryEntry{
				{
					Name:     "host-info",
					OSVendor: osname,
					HostName: hostname,
					NetworkInterfaces: &ctypes.NetworkInterfaces{
						Entries: interfaces,
					}},
				{
					Name: "anti-malware",
					List: &ctypes.List{
						Entries: []ctypes.ListEntry{
							{
								ProductInfo: ctypes.ProductInfo{
									Prod:               systemd.FindClamdProcess(),
									RealTimeProtection: "no",
									LastFullScanTime:   "n/a",
								},
							},
						},
					},
				},
				{
					Name: "disk-backup",
					List: &ctypes.List{
						Entries: []ctypes.ListEntry{},
					},
				},
				{
					Name: "disk-encryption",
					List: &ctypes.List{
						Entries: []ctypes.ListEntry{
							{
								ProductInfo: ctypes.ProductInfo{
									Prod: ctypes.Prod{
										Vendor:  "SomethingCrypto",
										Name:    "cryptsetup",
										Version: "1.2.3",
									},
									Drives: &ctypes.Drives{
										Entries: others.GetEncryptedPartitions(),
									},
								},
							},
						},
					},
				},
				{
					Name: "firewall",
					List: &ctypes.List{
						Entries: others.GetFirewall(),
					},
				},
				{
					Name: "patch-management",
					List: &ctypes.List{
						Entries: others.GetPackageManager(),
					},
					MissingPatches: &ctypes.MissingPatches{
						Entries: []ctypes.MissingPatchEntry{},
					},
				},
				{
					Name: "data-loss-prevention",
					List: &ctypes.List{
						Entries: []ctypes.ListEntry{},
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

	fmt.Println(string(xmlData))
	os.Exit(0)
}
