package license

import (
	"bytes"
	"crypto"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
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

func ValidateLicense(pathPublicKey, pathLicense string, runtimeInfo []byte) (bool, error) {
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

	// parse license content
	var (
		lensig     uint16
		lencontent uint16
		leninfo    uint16

		sig  []byte
		hid  []byte
		info []byte
	)

	cursor := 0
	fnReadHeader := func(v *uint16, cursor *int) {
		*v = binary.LittleEndian.Uint16(bylicense[*cursor : *cursor+2])
		*cursor += 2
	}
	fnRead := func(length uint16, v *[]byte, cursor *int) {
		*v = bylicense[*cursor : *cursor+int(length)]
		*cursor += int(length)
	}
	fnReadHeader(&lensig, &cursor)
	fnReadHeader(&lencontent, &cursor)
	fnReadHeader(&leninfo, &cursor)
	fnRead(lensig, &sig, &cursor)
	fnRead(lencontent, &hid, &cursor)
	fnRead(leninfo, &info, &cursor)

	// TODO: remove
	fmt.Println("string sig", base64.StdEncoding.EncodeToString(sig))
	fmt.Println("string hid", string(hid))
	fmt.Println("string info", string(info))

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

	// validate license integrity using its signature
	hash := crypto.SHA256.New()
	_, _ = hash.Write(hid) // intentionally making it vague
	_, _ = hash.Write(info)
	content := hash.Sum(nil)
	if err := rsa.VerifyPSS(publicKey.(*rsa.PublicKey), crypto.SHA256, content[:], sig, nil); err != nil {
		return false, err
	}

	// validate against runtime value
	runtimeID, err := machineid.ID()
	if err != nil {
		panic(err)
	}
	if !bytes.Equal(hid, []byte(runtimeID)) {
		return false, errors.New("invalid license")
	}
	if !bytes.Equal(info, runtimeInfo) {
		return false, errors.New("invalid license")
	}

	return true, nil
}
