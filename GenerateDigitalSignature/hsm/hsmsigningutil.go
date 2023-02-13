package hsmutil

import (
	"errors"
	"fmt"
	"strings"
)

var hsmClient HSMConnect

const (
	NoHash int = 0
	SHA1   int = 1
	MD5    int = 2
	RIPEMD int = 3
	SHA256 int = 4
)

type HSMSigningResult struct {
	IsSuccess    bool
	IsContinue   bool
	SignedData   string
	ErrorMessage string
}

func GenerateSignature() {
	inputMessage := "Hello World!"
	// message := []byte(inputMessage)
	// hash := sha256.Sum256(message)

	var signKeyIndex int16 = 23
	connectClient()

	signature, _ := generateSign(inputMessage, signKeyIndex, SHA256)
	fmt.Println("Digital Signature is: ", signature)
}

func generateSign(hexString string, signingKeyIndex int16, algorithm int) (string, error) {
	var result string
	var hSMSigningResult HSMSigningResult
	hSMSigningResult.IsSuccess = false

	list := split(hexString, 4096)
	for i := 0; i < len(list); i++ {
		hSMSigningResult = GetSignData(list[i], signingKeyIndex, i != len(list)-1, algorithm)
		if hSMSigningResult.IsSuccess {
			if !hSMSigningResult.IsContinue {
				result = hSMSigningResult.SignedData
				break
			}
			continue
		}
		return "", errors.New("Unable to sign data - " + hSMSigningResult.ErrorMessage)
	}

	return result, nil
}

func split(hexString string, chunkSize int) []string {
	var chunks []string
	runes := []rune(hexString)
	for i := 0; i < len(runes); i += chunkSize {
		end := i + chunkSize
		if end > len(runes) {
			end = len(runes)
		}
		chunks = append(chunks, string(runes[i:end]))
	}
	return chunks
}

func GetSignData(hexData string, privateKeyIndex int16, isContinue bool, algorithm int) HSMSigningResult {
	result := HSMSigningResult{}
	result.IsSuccess = false
	format := "[AORSAS;RC%d;RF%s;RG%d;BN%s;KY1;ZA1;]"
	request := fmt.Sprintf(format, privateKeyIndex, hexData, algorithm, map[bool]string{true: "1", false: "0"}[isContinue])
	endChar := "]"

	text := executeExcrypt(request, endChar, !isContinue)

	text = text[1 : len(text)-1]
	array := strings.Split(text, ";")
	errorMessage := ""
	for i := 0; i < len(array)-1; i++ {
		switch array[i][:2] {
		case "AO":
			text = strings.ToUpper(array[i][2:])
		case "BB":
			errorMessage = array[i][2:]
		case "BN":
			if strings.ToUpper(array[i][2:]) == "CONTINUE" {
				result.IsSuccess = true
				result.IsContinue = true
			} else {
				result.IsContinue = false
			}
		case "RH":
			result.IsSuccess = true
			result.SignedData = array[i][2:]
		}
	}

	if text == "ERRO" {
		result.IsSuccess = false
		result.ErrorMessage = errorMessage
	} else if text != "RSAS" {
		result.IsSuccess = false
		result.ErrorMessage = errorMessage
	}

	return result
}

func connectClient() {
	var hsmPrimaryIP int64 = 3232272392
	var hsmSecondaryIP int64 = 3232272392
	var hsmPort int = 9000

	fmt.Println("IP address is:", hsmPrimaryIP)

	hsmClient.NewHSMConnectWithPort(hsmPrimaryIP, hsmSecondaryIP, hsmPort)

	//Connect to HSM
	err := hsmClient.Connect()

	if err != nil {
		fmt.Errorf("Unable to connect to Primary HSM")
	}
}

func executeExcrypt(request, endChar string, endConnection bool) string {
	// now := time.Now()
	if !hsmClient.IsConnected() {
		connectClient()
		if !hsmClient.IsConnected() {
			return fmt.Sprintf("Unable to Connect Primary %s and Secondary HSM %s", hsmClient.PrimaryHSMIP(), hsmClient.SecondaryHSMIP())
		}
	}

	result, err := hsmClient.PostRequest(request, endChar)
	if err != nil {
		fmt.Errorf("Unable to post request to HSM, %w", err)
	}

	if endConnection {
		hsmClient.Disconnect()
	}

	return result
}
