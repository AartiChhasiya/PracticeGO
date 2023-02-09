package hsmutil

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"fmt"
	"net"
)

// func GenerateSignature() {
// 	message := []byte("Hello World!")
// 	hash := sha256.Sum256(message)

// 	var hsmPrimaryIP string = "3232272392"
// 	// var hsmSecondaryIP string = "3232272392"
// 	var hsmPort string = "9000"
// 	// var signKeyIndex int = 23

// 	// conn, err := net.Dial("tcp", "hsm-ip:hsm-port")
// 	conn, err := net.Dial(hsmPrimaryIP, hsmPort)
// 	if err != nil {
// 		fmt.Println("Error connecting to HSM:", err)
// 		return
// 	}
// 	defer conn.Close()

// 	privateKey, err := rsa.SignPKCS1v15(rand.Reader, conn, crypto.SHA256, hash[:])
// 	if err != nil {
// 		fmt.Println("Error signing message:", err)
// 		return
// 	}

// 	fmt.Printf("Digital Signature: %x\n", privateKey)
// }

func GenerateSignature() {
	message := []byte("Hello World!")
	hash := sha256.Sum256(message)

	var hsmPrimaryIP string = "3232272392"
	// var hsmSecondaryIP string = "3232272392"
	var hsmPort string = "9000"
	var signKeyIndex int = 23

	// conn, err := net.Dial("tcp", "hsm-ip:hsm-port")
	conn, err := net.Dial("tcp", hsmPrimaryIP+":"+hsmPort)
	if err != nil {
		fmt.Println("Error connecting to HSM:", err)
		return
	}
	defer conn.Close()

	// Send the index of the private key to the HSM
	index := []byte(signKeyIndex)
	conn.Write(index)

	signature, err := rsa.SignPKCS1v15(rand.Reader, conn, crypto.SHA256, hash[:])
	if err != nil {
		fmt.Println("Error signing message:", err)
		return
	}

	fmt.Printf("Digital Signature: %x\n", signature)
}
