package main

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/binary"
	"encoding/pem"
	"fmt"
	"os"
	"time"

	"github.com/denisbrodbeck/machineid"
)

func createLicense(hwlabel, pathPrivateKey string, isDebug bool) ([]byte, string) {
	id, err := machineid.ID()
	if err != nil {
		panic(err)
	}

	hash := crypto.SHA256.New()
	if _, err := hash.Write([]byte(id)); err != nil {
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
		fmt.Println("lenhid", len(hid))
		fmt.Println("string hid", base64.StdEncoding.EncodeToString(hid))
		fmt.Println("string sig", base64.StdEncoding.EncodeToString(sig))
	}

	lfp := fmt.Sprintf("./license_issued/license.%s.%d", hwlabel, time.Now().Unix())
	err = os.WriteFile(lfp, lic, os.ModePerm)
	if err != nil {
		panic(err)
	}

	return hid, lfp
}
