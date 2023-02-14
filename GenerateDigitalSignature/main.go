package main

import (
	"fmt"
	hsmutil "generatedigitalsignature/hsm"
)

// https://pkg.go.dev/azoo.dev/utils/dvx/hsm#section-readme
func main() {
	fmt.Println("Welcome to Generate Digital Signature")

	// inputMessage := `{"paymentReturnInstruction": {"TransactionId": "ae07582a-88f1-4f5f-8d1a-a564ae933e68","ReasonForReturn": ""}}`
	inputMessage := `{"MachineName":"CON-IND-09","UserName":"sapan.patibandha","Timestamp":"2018-10-17T11:12:05.4212206+05:30"}`
	// inputMessage := "Hello World!"

	var hsmPrimaryIP uint32 = 3232272392
	var hsmSecondaryIP uint32 = 3232272392
	var hsmPort int = 9000

	hsm := hsmutil.NewHSMSigningWithPort(hsmPrimaryIP, hsmSecondaryIP, hsmPort)
	hsm.GenerateSignature(inputMessage)
}
