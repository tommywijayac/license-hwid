package main

import (
	"flag"
	"log"
	"time"

	"github.com/denisbrodbeck/machineid"

	"github.com/tommywijayac/license-hwid/repo"
)

func main() {
	pGetID := flag.Bool("get_id", false, "get id")
	pGenerate := flag.Bool("generate", false, "generate license")
	pGenerate_HWLabel := flag.String("hwlabel", "", "hardware label")
	pGenerate_Debug := flag.Bool("debug", false, "debug")

	// scripts
	pHelperRSA := flag.Bool("rsa", false, "helper to generate public & private key")

	flag.Parse()

	if *pHelperRSA {
		GenerateRSAKey()
		return // ignore other flags
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
