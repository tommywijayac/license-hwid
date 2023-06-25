package main

import (
	"flag"
	"log"
	"time"
	"whiteboard/license/repo"

	"github.com/denisbrodbeck/machineid"
)

func main() {
	pGetID := flag.Bool("get_id", false, "get id")
	pGenerate := flag.Bool("generate", false, "generate license")
	pGenerate_ID := flag.String("id", "", "id")
	pGenerate_HWLabel := flag.String("hwlabel", "", "hardware label")
	pGenerate_Debug := flag.Bool("debug", false, "debug")

	flag.Parse()

	if *pGetID {
		id, err := machineid.ID()
		if err != nil {
			log.Print("ERR: ", err)
		}

		log.Print("id = ", id)
		return
	}

	if *pGenerate {
		db := repo.New()

		if len(*pGenerate_ID) == 0 {
			log.Print("ERR: generate: empty id")
			return
		}

		hid, lfp := CreateLicense(*pGenerate_ID, *pGenerate_HWLabel, "./rsa", *pGenerate_Debug)

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
