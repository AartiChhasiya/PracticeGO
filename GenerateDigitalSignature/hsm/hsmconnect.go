package hsmutil

import (
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"
)

type HSMConnect struct {
	primaryHSMIP   net.IP
	secondaryHSMIP net.IP
	tcpPort        int
	client         *net.TCPConn
	clientStream   net.Conn
	connected      bool
}

func (h *HSMConnect) PrimaryHSMIP() string {
	return h.primaryHSMIP.String()
}

func (h *HSMConnect) SecondaryHSMIP() string {
	return h.secondaryHSMIP.String()
}

func (h *HSMConnect) IsConnected() bool {
	return h.connected
}

// func (primaryHSMLongIP, secondaryHSMLongIP int64) NewHSMConnect() *HSMConnect {
// 	h := &HSMConnect{
// 		primaryHSMIP:   net.ParseIP(strconv.FormatInt(primaryHSMLongIP, 10)),
// 		secondaryHSMIP: net.ParseIP(strconv.FormatInt(secondaryHSMLongIP, 10)),
// 		tcpPort:        9000,
// 	}

// 	return h
// }

func (h *HSMConnect) NewHSMConnectWithPort(primaryHSMLongIP int64, secondaryHSMLongIP int64, port int) {
	h.primaryHSMIP = net.ParseIP(strconv.FormatInt(primaryHSMLongIP, 10))
	h.secondaryHSMIP = net.ParseIP(strconv.FormatInt(secondaryHSMLongIP, 10))
	h.tcpPort = port
}

func (h *HSMConnect) Connect() error {
	if h.primaryHSMIP == nil {
		return fmt.Errorf("Primary HSM IP must be specified")
	}

	var err error
	h.client, err = net.DialTCP("tcp", nil, &net.TCPAddr{IP: h.primaryHSMIP, Port: h.tcpPort})
	if err != nil {
		if h.secondaryHSMIP == nil {
			return fmt.Errorf("Unable to connect to Primary HSM")
		}

		h.client, err = net.DialTCP("tcp", nil, &net.TCPAddr{IP: h.secondaryHSMIP, Port: h.tcpPort})
		if err != nil {
			return fmt.Errorf("Unable to connect to Primary and secondary HSM")
		}
	}
	h.clientStream = h.client
	h.connected = true

	return nil
}

func (h *HSMConnect) Disconnect() {
	if h.client != nil {
		h.client.Close()
		h.client = nil
	}

	h.client = nil
	h.connected = false
}

func (h *HSMConnect) PostRequest(request, requestEndChar string) (string, error) {
	var resultBuilder strings.Builder

	if h.client != nil {
		requestBytes := []byte(request)
		_, err := h.clientStream.Write(requestBytes)
		if err != nil {
			return "", err
		}

		empty := ""
		num2 := 1
		array := make([]byte, num2)
		for {
			_, err := h.clientStream.Read(array)
			if err != nil {
				return "", err
			}
			empty = string(array)
			resultBuilder.WriteString(empty)
			if empty == requestEndChar {
				break
			}
		}
		return resultBuilder.String(), nil
	}
	return "", errors.New("HSM is not connected")
}

// func (obj *Signing) PostRequest(request, requestEndChar string) (string, error) {
// 	var strBuilder strings.Builder
// 	if obj.client != nil && obj.client.Connected() {
// 		bytes := []byte(request)
// 		_, err := obj.clientStream.Write(bytes)
// 		if err != nil {
// 			return "", err
// 		}
// 		var num int
// 		var empty string
// 		num2 := 1
// 		array := make([]byte, num2)
// 		for {
// 			flag := true
// 			num, err = obj.clientStream.Read(array)
// 			if err != nil {
// 				return "", err
// 			}
// 			empty = string(array)
// 			strBuilder.WriteString(empty)
// 			if empty == requestEndChar {
// 				break
// 			}
// 		}
// 		return strBuilder.String(), nil
// 	}
// 	return "", fmt.Errorf("HSM is not connected")
// }
