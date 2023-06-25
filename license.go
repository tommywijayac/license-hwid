package main

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/binary"
	"encoding/pem"
	"fmt"
	"os"
	"time"
)

// this is a secret method. move it to local file.
func CreateLicense(id, hwlabel, pathPrivateKey string, isDebug bool) ([]byte, string) {
	// machine ID must be hashed to obscure what we actually use as license
	// can't use machineid.ProtectedID since need to use same hash function given to rsa.SignPSS
	hash := crypto.SHA256.New()
	_, err := hash.Write([]byte(id))
	if err != nil {
		panic(err)
	}
	hid := hash.Sum(nil)

	// load private key
	bypempk, err := os.ReadFile(pathPrivateKey)
	if err != nil {
		panic(err)
	}
	block, _ := pem.Decode(bypempk)
	if block == nil {
		panic("fail to parse PEM block private key")
	}
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		panic(err)
	}

	// create signature from the hashed machine ID with it,
	// which can be verified with pub key by app to authenticate that we issue this license
	// in case that user manage to know our verification method (using machine ID) and generate their own license.
	sig, err := rsa.SignPSS(rand.Reader, privateKey, crypto.SHA256, hid[:], nil)
	if err != nil {
		panic(err)
	}

	lic := make([]byte, 2)
	binary.LittleEndian.PutUint16(lic, uint16(len(hid)))
	lic = append(lic, hid...)
	lic = append(lic, sig...)

	if isDebug {
		fmt.Println(len(hid))
		fmt.Println(hid)
		fmt.Println(sig)
	}

	lfp := fmt.Sprintf("./license_issued/license.%s.%d", hwlabel, time.Now().Unix())
	err = os.WriteFile(lfp, lic, os.ModePerm)
	if err != nil {
		panic(err)
	}

	return hid, lfp
}
