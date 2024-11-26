package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/denisbrodbeck/machineid"
	lic "github.com/tommywijayac/license-hwid"
	"github.com/tommywijayac/license-hwid/generator/repo"
)

func main() {
	pGetID := flag.Bool("get_id", false, "get id")
	pGenerate := flag.Bool("generate", false, "generate license")
	pGenerate_HWLabel := flag.String("hwlabel", "", "hardware label")
	pGenerate_Debug := flag.Bool("debug", false, "debug")

	// scripts
	pHelperRSA := flag.Bool("rsa", false, "helper to generate public & private key")

	pHelperVerify := flag.Bool("v", false, "")

	flag.Parse()

	if *pHelperRSA {
		GenerateRSAKey()
		return // ignore other flags
	}

	if *pHelperVerify {
		fmt.Println(lic.ValidateLicense("./secret/rsa.pub", "./license_issued/license.desktop_ubuntu.1732548673"))
		panic("exit")
	}

	if *pGetID {
		id, err := machineid.ID()
		if err != nil {
			log.Print("ERR: ", err)
		}

		log.Printf("id = %s", id)
		return
	}

	if *pGenerate {
		db := repo.New()

		hid, lfp := createLicense(*pGenerate_HWLabel, "./secret/rsa", *pGenerate_Debug)

		db.Add(repo.LicenseLog{
			HashedMachineID: hid,
			LicenseFilepath: lfp,
			HardwareLabel:   *pGenerate_HWLabel,
			CreatedTime:     time.Now().Format(time.RFC3339),
		})

		return
	}

	log.Print("ERR: invalid mode")
}
