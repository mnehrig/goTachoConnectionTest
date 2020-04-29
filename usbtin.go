package main

import (
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

func configureUSBtin(port io.ReadWriteCloser) (err error) {
	// Initializing the USB Device
	// Closing channel
	log.Debugln("Closing Channel")
	cmd := []byte("C\r")
	_, err = port.Write(cmd)
	if err != nil {
		log.Fatalf("Closing Channel Error: %v", err)
	}

	response := make([]byte, 1)
	_, err = port.Read(response)
	if err != nil {
		log.Fatalf("Reading Response Error: %v", err)
	}
	log.Debugf("closed channel: %x\n", response)

	time.Sleep(time.Millisecond * 100)
	// Initializing the USB Device
	log.Debugln("Setting Baudrate")
	cmd = []byte("S5\r")
	_, err = port.Write(cmd)
	if err != nil {
		log.Fatalf("Closing Channel Error: %v", err)
	}

	response = make([]byte, 1)
	_, err = port.Read(response)
	if err != nil {
		log.Fatalf("Reading Response Error: %v", err)
	}
	log.Debugf("set baudrate: %x\n", response)

	//Setting Timestap off
	log.Debugln("Setting Timestamp off")
	cmd = []byte("Z0\r")
	_, err = port.Write(cmd)
	if err != nil {
		log.Fatalf("TimeStam off Error: %v", err)
	}

	response = make([]byte, 1)
	_, err = port.Read(response)
	if err != nil {
		log.Fatalf("Reading Response Error: %v", err)
	}
	log.Debugf("TimeStamp Set off: %x\n", response)

	// Initializing the USB Device
	log.Debugln("Openen Canbus Channel")
	cmd = []byte("O\r")
	_, err = port.Write(cmd)
	if err != nil {
		log.Fatalf("Opening Channel Error: %v", err)
	}

	response = make([]byte, 1)
	_, err = port.Read(response)
	if err != nil {
		log.Fatalf("Reading Response Error: %v", err)
	}
	log.Debugf("canbus channel opened: %x\n", response)
	return err
}

func readFrame(port io.ReadWriter) (response string, err error) {
	buffer := []byte{}
	for {
		r := make([]byte, 1)
		_, err = port.Read(r)
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Debugf("Reading Response Error: %v", err)
		}
		buffer = append(buffer, r[0])
		if r[0] == '\r' {
			break
		}
		if r[0] == 7 {
			break
		}
	}
	response = string(buffer)
	return response, err
}

var internalIDcounter int = 0

// Frame struct
type Frame struct {
	internalID       int
	frameType        string
	canIDTacho       string
	canAddressVU     string
	canAddressFMS    string
	canMessageLength int
	frameContent     string
	padding          string
	frameString      string
}

const frameFormat string = "T%s%s%s%d%s%s\r"

func sendFrame(port io.ReadWriteCloser, frameContent string) {
	length := len(frameContent)
	padding := ""
	if length < canMessageLength*2 {
		padding = strings.Repeat("F", 8*2-len(frameContent))
	}
	message := fmt.Sprintf(frameFormat, canIDTacho, canAddressVU, canAddressFMS, canMessageLength, frameContent, padding)
	log.Debugf("SENDING COMMAND: %s\n%X\n", message, message)
	port.Write([]byte(message))
}

func printFrame(frame Frame) {
	println("[printFrame]", "ID:", frame.internalID, " RAW:", frame.frameString, " TYPE:", frame.frameType, " TACHO_ID:", frame.canIDTacho, " SENDER:", frame.canAddressVU, " RECEIVER:", frame.canAddressFMS, " MESSAGE LENGTH:", frame.canMessageLength, " CONTENT:", frame.frameContent)
}
func unmarshallFrame(message string) (frame Frame, err error) {
	frame.frameString = message

	internalIDcounter = internalIDcounter + 1
	frame.internalID = internalIDcounter

	if len(message) == 0 {
		return frame, errors.New("Empty message")
	}

	frame.frameType = string(message[0:1])

	switch frame.frameType {
	case "T":
		if len(message) != 26 {
			return frame, errors.New("Wrong message Length")
		}

		frame.canIDTacho = message[1:5]
		frame.canAddressVU = message[5:7]
		frame.canAddressFMS = message[7:9]

		var lstring string
		lstring = message[9:10]

		l, err := strconv.Atoi(lstring)
		if err != nil {
			return frame, errors.New("Message Length Byte Error. Value: " + lstring + " Error: " + err.Error())
		}
		frame.canMessageLength = l
		frame.frameContent = message[10 : 10+l*2]
		frame.padding = message[10+l*2:]
	case "Z":
		// Do nothing
	default:
	}

	return frame, err
}
