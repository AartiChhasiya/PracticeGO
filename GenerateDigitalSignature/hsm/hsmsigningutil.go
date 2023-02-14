package hsmutil

import (
	"encoding/base64"
	"encoding/hex"
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

const RSAS_MAX_LENGTH int = 4096 //Max length supported by HSM function

type HSMSigningResult struct {
	IsSuccess    bool
	IsContinue   bool
	SignedData   string
	ErrorMessage string
	CHCommand    string
}

func GenerateSignature() {
	// inputMessage := `{"paymentReturnInstruction": {"TransactionId": "ae07582a-88f1-4f5f-8d1a-a564ae933e68","ReasonForReturn": ""}}`
	// inputMessage := `{"MachineName":"CON-IND-09","UserName":"sapan.patibandha","Timestamp":"2018-10-17T11:12:05.4212206+05:30"}`

	inputMessage := "Hello World!"
	hashOfInputData := hex.EncodeToString([]byte(inputMessage))
	fmt.Println("SignData is", hashOfInputData)

	var signKeyIndex int16 = 23
	connectClient()

	HSMSignature := generateSign(hashOfInputData, signKeyIndex, SHA256)
	fmt.Println("HSMSignature is:", HSMSignature)

	body, err := hex.DecodeString(HSMSignature)
	if err != nil {
		panic(fmt.Sprintf("Unable to generate signature - %s", err.Error()))
	}

	//signature := base64.RawStdEncoding.EncodeToString(body)
	signature := base64.StdEncoding.EncodeToString(body)
	fmt.Println("Digital Signature is: ", signature)
}

func generateSign(hexString string, signingKeyIndex int16, algorithm int) string {
	signedData := ""
	chCommand := ""

	objResult := HSMSigningResult{}
	objResult.IsSuccess = false
	hexChunkList := split(hexString, RSAS_MAX_LENGTH)

	for i := 0; i < len(hexChunkList); i++ {

		var isContinue bool = false
		if len(hexChunkList)-1 > 0 {
			isContinue = true
		}

		objResult = getSignData(
			hexChunkList[i],
			signingKeyIndex,
			isContinue,
			algorithm,
			func() string {
				if chCommand != "" {
					return chCommand
				}
				return ""
			}(),
		)

		chCommand = objResult.CHCommand

		if objResult.IsSuccess {
			if !objResult.IsContinue {
				signedData = objResult.SignedData
				break
			}
		} else {
			panic(fmt.Sprintf("Unable to sign data - %s", objResult.ErrorMessage))
		}
	}

	return signedData
}

func split(str string, chunkSize int) []string {
	var listArray []string
	remaining := 0

	for i := 0; i < len(str)/chunkSize; i++ {
		listArray = append(listArray, str[i*chunkSize:(i+1)*chunkSize])
	}

	if remaining*chunkSize < len(str) {
		listArray = append(listArray, str[remaining*chunkSize:])
	}

	return listArray
}

func getSignData(hexData string, privateKeyIndex int16, isContinue bool, algorithm int, sChcommand string) HSMSigningResult {
	var objResult HSMSigningResult
	objResult.IsSuccess = false

	var commandParam string
	var cmd string

	if strings.TrimSpace(sChcommand) == "" {
		commandParam = "[AORSAS;RC%d;RF%s;RG%d;BN%s;KY1;ZA1;]" // *KY1(BER encoding of the HASH); *ZA1(Padding (Default))
		cmd = fmt.Sprintf(commandParam,
			privateKeyIndex, // %d private key index
			hexData,         // %s Data used to generate the signature
			algorithm,       // %d Hash algorithm
			map[bool]string{true: "1", false: "0"}[isContinue]) // %s Send Data in chunk or not
	} else {
		// This section of command need to build when we have to pass CH parameter.
		commandParam = "[AORSAS;CH%s;RC%d;RF%s;RG%d;BN%s;KY1;ZA1;]" // *KY1(BER encoding of the HASH); *ZA1(Padding (Default))
		cmd = fmt.Sprintf(commandParam,
			sChcommand,      // %s CH command for split data
			privateKeyIndex, // %d private key index
			hexData,         // %s Data used to generate the signature
			algorithm,       // %d Hash algorithm
			map[bool]string{true: "1", false: "0"}[isContinue]) // %s Send Data in chunk or not
	}

	var endChar string = "]"
	response := executeExcrypt(cmd, endChar, !isContinue)
	functionID := ""

	response = response[1 : len(response)-1]
	resultArray := strings.Split(response, ";")

	var message string

	for i := 0; i < len(resultArray)-1; i++ {
		str := resultArray[i][0:2]
		data := resultArray[i][2:]
		switch str {
		case "AO":
			functionID = strings.ToUpper(data)
		case "BB":
			message = data
		case "BN":
			if strings.ToUpper(data) == "CONTINUE" {
				objResult.IsSuccess = true
				objResult.IsContinue = true
			} else {
				objResult.IsContinue = false
			}
		case "RH":
			objResult.IsSuccess = true
			objResult.SignedData = data
		case "CH":
			objResult.CHCommand = data
		}
	}

	if functionID == "ERRO" {
		objResult.IsSuccess = false
		objResult.ErrorMessage = message
	} else if functionID != "RSAS" {
		objResult.IsSuccess = false
		objResult.ErrorMessage = message
	}

	return objResult
}

func connectClient() {
	var hsmPrimaryIP uint32 = 3232272392
	var hsmSecondaryIP uint32 = 3232272392
	var hsmPort int = 9000

	// fmt.Println("IP address is:", hsmPrimaryIP)

	hsmClient.NewHSMConnectWithPort(hsmPrimaryIP, hsmSecondaryIP, hsmPort)

	//Connect to HSM
	err := hsmClient.Connect()

	if err != nil {
		panic(fmt.Sprintf("Unable to connect to Primary HSM - %s", err.Error()))
		// fmt.Errorf("unable to connect to Primary HSM, %w", err)
	}
}

func executeExcrypt(request, endChar string, endConnection bool) string {
	// now := time.Now()
	if !hsmClient.IsConnected() {
		connectClient()
		if !hsmClient.IsConnected() {
			return fmt.Sprintf("unable to Connect Primary %s and Secondary HSM %s", hsmClient.PrimaryHSMIP(), hsmClient.SecondaryHSMIP())
		}
	}

	result, err := hsmClient.PostRequest(request, endChar)
	if err != nil {
		panic(fmt.Sprintf("Unable to post request to HSM - %s", err.Error()))
		// fmt.Errorf("unable to post request to HSM, %w", err)
	}

	if endConnection {
		hsmClient.Disconnect()
	}

	return result
}
