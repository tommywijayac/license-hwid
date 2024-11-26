package license

import (
	"bytes"
	"crypto"
	"crypto/rsa"
	"crypto/x509"
	"encoding/binary"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/denisbrodbeck/machineid"
)

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
