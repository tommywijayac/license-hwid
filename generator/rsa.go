package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"os"
)

// GenerateRSAKey is a helper function to create pairs of private and public key.
func GenerateRSAKey() {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic(err)
	}

	byPrivate := x509.MarshalPKCS1PrivateKey(privateKey)
	pemPrivate := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: byPrivate,
	})
	_ = os.WriteFile("./secret/rsa", pemPrivate, 0644)

	byPublic, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		panic(err)
	}
	pemPublic := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: byPublic,
	})
	_ = os.WriteFile("./secret/rsa.pub", pemPublic, 0644)
}
