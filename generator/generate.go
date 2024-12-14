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
)

func createLicense(hwLabel, hwID, pathAddInfo, pathPrivateKey string, isDebug bool) (string, string) {
	var (
		err     error
		hid     []byte
		info    []byte
		content []byte
	)

	hash := crypto.SHA256.New()

	// Add machine id to content
	if _, err := hash.Write([]byte(hwID)); err != nil {
		panic(err)
	}
	hid = []byte(hwID)

	// Add additional info to content (optional)
	if len(pathAddInfo) > 0 {
		info, err = os.ReadFile(pathAddInfo)
		if err != nil {
			panic(err)
		}

		if _, err := hash.Write(info); err != nil {
			panic(err)
		}
	}
	content = hash.Sum(nil)

	// Create signature from license content.
	// It can be verified later with public key pair to ensure we generate the license & not tampered.
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

	sig, err := rsa.SignPSS(rand.Reader, privateKey, crypto.SHA256, content[:], nil)
	if err != nil {
		panic(err)
	}

	// Compose license
	lic := make([]byte, 6)
	binary.LittleEndian.PutUint16(lic[:2], uint16(len(sig)))
	binary.LittleEndian.PutUint16(lic[2:4], uint16(len(hid)))
	binary.LittleEndian.PutUint16(lic[4:6], uint16(len(info)))

	lic = append(lic, sig...)
	lic = append(lic, hid...)
	lic = append(lic, info...)

	if isDebug {
		fmt.Println("string sig", base64.StdEncoding.EncodeToString(sig))
		fmt.Println("string hid", string(hid))
		fmt.Println("string info", string(info))
	}

	lfp := fmt.Sprintf("./license_issued/license.%s", hwLabel)
	err = os.WriteFile(lfp, lic, os.ModePerm)
	if err != nil {
		panic(err)
	}

	return string(info), lfp
}
