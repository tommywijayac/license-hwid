package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/denisbrodbeck/machineid"
	lic "github.com/tommywijayac/license-hwid"
	"github.com/tommywijayac/license-hwid/generator/repo"
)

func main() {
	pGetID := flag.Bool("get_id", false, "get id")
	pGenerate := flag.Bool("generate", false, "generate license")
	pGenerate_HWLabel := flag.String("hwlabel", "", "hardware label")
	pGenerate_HWID := flag.String("hwid", "", "machine id")
	pGenerate_AddInfo := flag.String("info", "", "path to file containing additional info to be validated. must not be sensitive")
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
		info, err := os.ReadFile(*pGenerate_AddInfo)
		if err != nil {
			panic(err)
		}

		fmt.Println(lic.ValidateLicense("./secret/rsa.pub", "./license_issued/license.mac", info))
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

		if len(*pGenerate_HWID) == 0 {
			panic("id must be set")
		}

		info, lfp := createLicense(
			*pGenerate_HWLabel,
			*pGenerate_HWID,
			*pGenerate_AddInfo,
			"./secret/rsa",
			*pGenerate_Debug,
		)

		db.Add(repo.LicenseLog{
			HashedMachineID: *pGenerate_HWID,
			AdditionalInfo:  info,
			LicenseFilepath: lfp,
			HardwareLabel:   *pGenerate_HWLabel,
			CreatedTime:     time.Now().Format(time.RFC3339),
		})

		return
	}

	log.Print("ERR: invalid mode")
}
