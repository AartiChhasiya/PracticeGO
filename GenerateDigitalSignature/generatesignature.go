package main

import (
	"crypto/sha256"
	"fmt"

	"github.com/miekg/pkcs11"
)

func main() {
	// Load the PKCS#11 library for the HSM
	p := pkcs11.New("/usr/lib/softhsm/libsofthsm2.so")
	err := p.Initialize()
	if err != nil {
		fmt.Println("Failed to initialize HSM:", err)
		return
	}
	defer p.Destroy()
	defer p.Finalize()

	// Open a session with the HSM
	session, err := p.OpenSession(pkcs11.CKF_SERIAL_SESSION|pkcs11.CKF_RW_SESSION, 0)
	if err != nil {
		fmt.Println("Failed to open session:", err)
		return
	}
	defer session.Close()

	// Login to the HSM
	err = session.Login(pkcs11.CKU_USER, "1234")
	if err != nil {
		fmt.Println("Failed to login:", err)
		return
	}

	// Hash the message
	message := []byte("This is the message to be signed")
	hash := sha256.Sum256(message)

	// Sign the hash using a key from the HSM
	signature, err := session.Sign(hash[:], nil, pkcs11.Mechanism{pkcs11.CKM_SHA256_RSA_PKCS, nil})
	if err != nil {
		fmt.Println("Failed to sign:", err)
		return
	}

	fmt.Println("Signature:", signature)
}
