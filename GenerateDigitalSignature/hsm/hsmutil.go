package hsmutil

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/binary"
	"fmt"
	"net"
)

func GenerateSignature() {
	message := []byte("Hello World!")
	hash := sha256.Sum256(message)

	hsmPrimaryIP := Uint32ToIPv4(3232272392)
	var hsmPort string = "9000"
	var signKeyIndex string = "23"

	fmt.Println("IP address is:", hsmPrimaryIP)

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

	// Retrieve the private key from the HSM using the index and the established connection
	privateKey, err := retrievePrivateKeyFromHSM(conn, index)
	if err != nil {
		fmt.Println("Error retrieving private key from HSM:", err)
		return
	}

	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, hash[:])
	if err != nil {
		fmt.Println("Error signing message:", err)
		return
	}

	fmt.Printf("Digital Signature: %x\n", signature)
}

func Uint32ToIPv4(ip uint32) string {
	// Create a 4-byte array from the uint32 value
	ipBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(ipBytes, ip)

	// Convert the 4-byte array to a net.IP
	ipAddr := net.IPv4(ipBytes[0], ipBytes[1], ipBytes[2], ipBytes[3])

	// Return the string representation of the IP address
	return ipAddr.String()
}

func retrievePrivateKeyFromHSM(conn net.Conn, index []byte) (*rsa.PrivateKey, error) {
	// Send a request to the HSM to retrieve the private key using the index
	// _, err := conn.Write(index)
	// if err != nil {
	// 	return nil, fmt.Errorf("Error sending request to retrieve private key: %v", err)
	// }

	// Read the response from the HSM containing the private key
	var privateKeyBytes []byte
	_, err := conn.Read(privateKeyBytes)
	if err != nil {
		return nil, fmt.Errorf("error reading response from hsm: %v", err)
	}

	// Parse the private key bytes into an RSA private key structure
	privateKey, err := x509.ParsePKCS1PrivateKey(privateKeyBytes)
	if err != nil {
		return nil, fmt.Errorf("error parsing private key: %v", err)
	}

	return privateKey, nil
}
