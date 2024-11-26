package main

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/binary"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
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

func ValidatePublicKey(wantBypempub []byte, pathPublicKey string) (bool, error) {
	parse := func(bypempub []byte) (*rsa.PublicKey, error) {
		block, _ := pem.Decode(bypempub)
		if block == nil {
			return nil, fmt.Errorf("fail to decode public key")
		}
		publicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
		if err != nil {
			return nil, fmt.Errorf("fail to parse public key: %s", err.Error())
		}
		return publicKey.(*rsa.PublicKey), nil
	}

	want, err := parse(wantBypempub)
	if err != nil {
		return false, err
	}

	bypempub, err := os.ReadFile(pathPublicKey)
	if err != nil {
		return false, fmt.Errorf("fail to open public key: %s", err.Error())
	}
	got, err := parse(bypempub)
	if err != nil {
		return false, err
	}

	return want.Equal(got), nil
}

func ValidateLicense(pathPublicKey, pathLicense string) (bool, error) {
	// load license
	license, err := os.Open(pathLicense)
	if err != nil {
		return false, fmt.Errorf("fail to open license: %s", err.Error())
	}
	defer license.Close()

	bylicense, err := io.ReadAll(license)
	if err != nil {
		return false, fmt.Errorf("fail to parse license: %s", err.Error())
	}

	// parse license element
	offset := binary.LittleEndian.Uint16(bylicense[:2])
	hid := bylicense[2 : 2+offset]
	sig := bylicense[2+offset:]

	// load public key
	bypempub, err := os.ReadFile(pathPublicKey)
	if err != nil {
		return false, fmt.Errorf("fail to open public key: %s", err.Error())
	}
	block, _ := pem.Decode(bypempub)
	if block == nil {
		return false, fmt.Errorf("fail to decode public key")
	}
	publicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return false, fmt.Errorf("fail to parse public key: %s", err.Error())
	}

	// validate signature
	if err := rsa.VerifyPSS(publicKey.(*rsa.PublicKey), crypto.SHA256, hid[:], sig, nil); err != nil {
		return false, err
	}

	// validate content
	id, err := machineid.ID()
	if err != nil {
		panic(err)
	}

	hash := crypto.SHA256.New()
	if _, err := hash.Write([]byte(id)); err != nil {
		panic(err)
	}
	wanthid := hash.Sum(nil)

	if !bytes.Equal(wanthid, hid) {
		return false, errors.New("invalid license")
	}
	return true, nil
}
