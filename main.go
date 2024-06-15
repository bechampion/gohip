package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"osdata"
	"others"
	"strings"
	"systemd"
	"time"
	// ctypes "types"
	ctypes "github.com/bechampion/gohip/types"
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
	systemd.FindClamdProcess()
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
									Prod: systemd.FindClamdProcess(),
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
										// Entries: []ctypes.DriveEntry{
										// 	{
										// 		DriveName: "/",
										// 		EncState:  "unencrypted",
										// 	},
										// 	{
										// 		DriveName: "/home/user",
										// 		EncState:  "unencrypted",
										// 	},
										// 	{
										// 		DriveName: "All",
										// 		EncState:  "unencrypted",
										// 	},
										// },
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
