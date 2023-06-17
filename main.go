package main

import (
	"flag"
	"log"

	"github.com/denisbrodbeck/machineid"
)

func main() {
	pGetID := flag.Bool("get_id", false, "get id")
	pGenerate := flag.Bool("generate", false, "generate license")
	pGenerate_ID := flag.String("id", "", "id")
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
		if len(*pGenerate_ID) == 0 {
			log.Print("ERR: generate: empty id")
			return
		}

		CreateLicense(*pGenerate_ID, "./rsa", *pGenerate_Debug)
		return
	}

	log.Print("ERR: invalid mode")
}
