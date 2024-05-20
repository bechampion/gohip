package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"osdata"
	"strings"
	"systemd"
	"time"
	ctypes "types"
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

func main() {
	logCommandAndArgs()
	systemd.FindAVUnit()
	cookie := flag.String("cookie", "", "")
	 _ = flag.String("client-os", "", "")
	//--client-ip seems to be fed from openconect but i don't think it's used
	clientip  := flag.String("client-ip", "", "")
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
	hipReport := ctypes.HipReport{
		Name:         "hip-report",
		GenerateTime: time.Now().Format("01/02/2006 15:04:05"),
		Version:      4,
		User: user,
		HostName: hostname,
		HostId: hostname,
		Md5: *md5,
		Ip: *clientip,
		Domain: domain,
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
									Prod: ctypes.Prod{
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
										Vendor:  "GitLab Inc.",
										Name:    "cryptsetup",
										Version: "2.4.3",
									},
									Drives: &ctypes.Drives{
										Entries: []ctypes.DriveEntry{
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
					List: &ctypes.List{
						Entries: []ctypes.ListEntry{
							{
								ProductInfo: ctypes.ProductInfo{
									Prod: ctypes.Prod{
										Vendor:  "Canonical Ltd.",
										Name:    "UFW",
										Version: "0.36.1",
									},
									IsEnabled: "no",
								},
							},
							{
								ProductInfo: ctypes.ProductInfo{
									Prod: ctypes.Prod{
										Vendor:  "IPTables",
										Name:    "IPTables",
										Version: "1.8.7",
									},
									IsEnabled: "no",
								},
							},
							{
								ProductInfo: ctypes.ProductInfo{
									Prod: ctypes.Prod{
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
					List: &ctypes.List{
						Entries: []ctypes.ListEntry{
							{
								ProductInfo: ctypes.ProductInfo{
									Prod: ctypes.Prod{
										Vendor:  "Canonical Ltd.",
										Name:    "Snap",
										Version: "22.04.1",
									},
									IsEnabled: "n/a",
								},
							},
							{
								ProductInfo: ctypes.ProductInfo{
									Prod: ctypes.Prod{
										Vendor:  "GNU",
										Name:    "Advanced Packaging Tool",
										Version: "2.4.11",
									},
									IsEnabled: "yes",
								},
							},
						},
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
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}

}
