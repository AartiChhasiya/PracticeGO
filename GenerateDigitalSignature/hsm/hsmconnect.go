package hsmutil

import (
	"errors"
	"fmt"
	"net"
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

func (h *HSMConnect) NewHSMConnectWithPort(primaryHSMLongIP uint32, secondaryHSMLongIP uint32, port int) {
	h.primaryHSMIP = net.ParseIP(Uint32ToIPv4(primaryHSMLongIP))
	h.secondaryHSMIP = net.ParseIP(Uint32ToIPv4(secondaryHSMLongIP))
	h.tcpPort = port

	fmt.Printf("Primary IP is: %s Secondary IP is: %s Port is: %d \n", h.primaryHSMIP, h.secondaryHSMIP, h.tcpPort)
}

func (h *HSMConnect) Connect() error {
	if h.primaryHSMIP == nil {
		return fmt.Errorf("primary HSM IP must be specified, %w", nil)
	}

	var err error
	h.client, err = net.DialTCP("tcp", nil, &net.TCPAddr{IP: h.primaryHSMIP, Port: h.tcpPort})
	if err != nil {
		if h.secondaryHSMIP == nil {
			return fmt.Errorf("unable to connect to Primary HSM, %w", err)
		}

		h.client, err = net.DialTCP("tcp", nil, &net.TCPAddr{IP: h.secondaryHSMIP, Port: h.tcpPort})
		if err != nil {
			return fmt.Errorf("unable to connect to Primary and secondary HSM, %w", err)
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
