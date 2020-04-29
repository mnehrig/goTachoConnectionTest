package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

type stateContainer struct {
	currentState        state
	expectingDataFromVU bool
	vuCanAcceptRoutines bool
}

var appState *stateContainer = &stateContainer{
	currentState:        stateUndefined,
	expectingDataFromVU: false,
}

// sendTLPMessage send Transport Layer Protocol Message
func sendTLPMessage(port io.ReadWriteCloser, frameContent string) {
	l := len(frameContent) / 2 // 1 Byte = 0xFF is two Ascii letters long

	if l <= 7 {
		// Send Single Frame Message "0"
		sendFrame(port, "0"+strconv.Itoa(l)+frameContent)
	} else {
		// Send Multi Frame Package
		// Send First Message
		firstFrameContent := "1"                     // 4Bit/1ASCII Trasport Layer Protocol Message Type 1: First Message
		firstFrameContent += fmt.Sprintf("%x-03", l) // 12Bit/3ASCII Trasport Layer Protocol Message Length
		firstFrameContent += frameContent[0:12]      // First 6Bytes,12ASCII Data
		restContent := frameContent[12:]
		log.Debugln("FirstFrame: ", firstFrameContent)
		sendFrame(port, firstFrameContent)
		// Wait for Flow Control Response Parameters

		f, err := receiveFlowControlResponse(port)
		if err != nil {
			log.Debugln(err.Error())
		}
		log.Debugln("Received Flow Control Response")

		// Send Data According to Flow Control Response Parameters
		lenRest := len(restContent)   // in ASCII Letters
		numFullBlocks := lenRest / 14 // Number of Full Type 2 Messages
		lenLastBlock := lenRest % 14  // If lenLastBlock >0 one more Type 2 Message with Rest Data filled up with FF

		// numFullBlocks Type 2 Blocks
		for i := 1; i <= numFullBlocks; i++ {
			time.Sleep(f.Gap())                                           // Wait for the amout of time listed in Flow Control Response
			consecutiveFrameContent := "2"                                // 4Bit/1ASCII Trasport Layer Protocol Message Type 2: Consecutive Data Message
			consecutiveFrameContent += strconv.FormatInt(int64(i%16), 16) // 4Bit Block Repeating Counter
			consecutiveFrameContent += restContent[i*14-14 : i*14]        // 7Bytes/14ASCII of Data.
			log.Debugln("ConsecutiveFrame: ", consecutiveFrameContent)
			sendFrame(port, consecutiveFrameContent) // Send Frame

		}

		i := numFullBlocks + 1

		if lenLastBlock != 0 {
			time.Sleep(f.Gap())                                                    // Wait for the amout of time listed in Flow Control Response
			consecutiveFrameContent := "2"                                         // 4Bit/1ASCII Trasport Layer Protocol Message Type 2: Consecutive Data Message
			consecutiveFrameContent += strconv.FormatInt(int64(i%16), 16)          // 4Bit Block Repeating Counter
			consecutiveFrameContent += restContent[i*14-14 : i*14-14+lenLastBlock] // 7Bytes/14ASCII of Data
			consecutiveFrameContent += strings.Repeat("F", 14-lenLastBlock)        // Fill Up with FFFFFF
			log.Debugln("LastFrame: ", consecutiveFrameContent)
			sendFrame(port, consecutiveFrameContent)
		}
	}
}

type TLPMessageType string

const TLPStandardSinglePacket TLPMessageType = "0"
const TLPFirstMultiPacket TLPMessageType = "1"
const TLPConsecutiveMultiPacket TLPMessageType = "2"
const TLPFlowControlResponse TLPMessageType = "3"

type FlowControlCommand string
type BlockSize string
type SeparationTime string

const FCClear FlowControlCommand = "0"
const FCWait FlowControlCommand = "1"
const FCAbort FlowControlCommand = "2"

type FlowControlFrame struct {
	FC FlowControlCommand // 4Bit / 1ASCII
	BS BlockSize          //
	ST SeparationTime
}

func (f *FlowControlFrame) Gap() time.Duration {

	bs, err := strconv.ParseInt(string(f.BS), 16, 8)
	if err != nil {
		log.Debugln("FlowControlFrame GAP Parsint error", err.Error())
	}
	if bs <= 127 {
		return time.Duration(bs) * time.Millisecond
	}

	bs0 := f.BS[0]
	bs1, err := strconv.ParseInt(string(f.BS[1]), 16, 8)
	if bs0 == 'F' && bs1 > 0 && bs1 <= 9 {
		return time.Duration(bs1*100) * time.Microsecond
	}

	return 0
}

type UnexpectedFrameError struct {
	err string
	frm Frame
}

func (e *UnexpectedFrameError) Error() string {
	return e.err
}

func (e *UnexpectedFrameError) Frame() Frame {
	return e.frm
}

func receiveFlowControlResponse(port io.ReadWriteCloser) (f FlowControlFrame, err error) {
	for {
		resp, err := readFrame(port)
		if err != nil {
			log.Debugln(err)
			// TODO react on errors
		}
		frame, err := unmarshallFrame(resp)

		switch TLPMessageType(frame.frameContent[0]) {
		case TLPConsecutiveMultiPacket:
			return FlowControlFrame{
				FC: FlowControlCommand(frame.frameContent[1:2]),
				BS: BlockSize(frame.frameContent[2:5]),
				ST: SeparationTime(frame.frameContent[5:8]),
			}, nil
		case TLPStandardSinglePacket:
			// Error
		case TLPFirstMultiPacket:
			// Error

		case TLPFlowControlResponse:

		default:
			// Error
		}
	}
}

func receiveTLPMessage() {

}

func RemoteSessionValue(port io.ReadWriteCloser) {
	function := "22" // Length
	serviceID := ""  // TODO serviceID?
	value := "F900"  // GiagnosticSessionType remoteSession
	frameContent := function + serviceID + value
	sendFrame(port, frameContent)

	/* 	for {
		resp, err := readFrame(port)
		if err != nil {
			log.Debugln(err)
		}
		frame, err := unmarshallFrame(resp)
		if err != nil {
			log.Debugln(err)
		}

		//if frame.Content
		log.Debugln("STEP 1: RESPONSE FROM VU:", frame.frameContent)
	} */

}

func startRemoteSession(port io.ReadWriteCloser) {
	sid := "10"   // DiagnosticSessionControl
	value := "FE" // GiagnosticSessionType remoteSession
	frameContent := sid + value
	sendTLPMessage(port, frameContent)

	/* 	for {
		resp, err := readFrame(port)
		if err != nil {
			log.Debugln("startRemoteSession Read Response Error", err.Error())
		}
		frame, err := unmarshallFrame(resp)
		if err != nil {
			log.Debugln("UnMarshallError", err.Error())
		}
		log.Debugln("STEP 1: RESPONSE FROM VU:", frame.frameString)
		//if frame.Content
	} */

}

func endRemoteSession(port io.ReadWriteCloser) {
	sid := "10"   // DiagnosticSessionControl
	value := "01" // end remoteSession
	frameContent := sid + value
	sendTLPMessage(port, frameContent)

	/* 	for {
	   		resp, err := readFrame(port)
	   		if err != nil {
	   			log.Debugln(err)
	   		}
	   		frame, err := unmarshallFrame(resp)
	   		log.Debugln("STEP 1: RESPONSE FROM VU:", frame.frameContent)

	   		// TODO Check result. Should be positive response sid +0x40
	   		//if frame.Content
	   	}
	*/
}

func closeAuthenticationProcess(port io.ReadWriteCloser) {
	sid := "31"       // DiagnosticSessionControl
	value := "010180" // RemoteAuth
	control := "09"
	frameContent := sid + value + control
	sendTLPMessage(port, frameContent)

	/* 	for {
		resp, err := readFrame(port)
		if err != nil {
			log.Debugln(err)
		}
		frame, err := unmarshallFrame(resp)
		log.Debugln("STEP 1: RESPONSE FROM VU:", frame.frameContent)
		//if frame.Content
	} */

}

func analyzeTestSession() {
	file, err := os.Open("test.log.txt")

	if err != nil {
		log.Fatalf("failed opening file: %s", err)
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	var txtlines []string

	for scanner.Scan() {
		txtlines = append(txtlines, scanner.Text())
	}

	file.Close()

	for _, eachline := range txtlines {
		frame, err := unmarshallFrame(eachline)
		if err != nil {
			log.Debugln("Error", frame.internalID, frame.frameString, err.Error())
			time.Sleep(time.Millisecond)
			continue
		}
		switch frame.frameType {
		case "T":
			switch frame.canAddressVU {
			case canAddressVU:
				handleVUFrame(frame)
			case canAddressFMS:
				handleFMSFrame(frame)
			default:
				log.Debugln("Error: Unkown Sender")
				printFrame(frame)
			}
		case "Z":
			handleZFrame(frame)
		default:
			log.Debugln("Error: Unkown Message Type")
			printFrame(frame)
		}

		time.Sleep(time.Millisecond)
	}
}

func reset(level resetLevel) {
	switch level {
	case resetLevelTotal:
		log.Debugln("reset level total")
		// TODO Reset Level Total
	default:
		log.Debugln("reset level error")
		// TODO Reset Level Error
	}
}

func handleZFrame(frame Frame) {
	log.Debugln("Handle Z Frame", frame.internalID)
	printFrame(frame)
}

func handleFMSFrame(frame Frame) {
	log.Debugln("Handle FMS Frame", frame.internalID)
	printFrame(frame)
}

type transportLayerState string

const (
	tlsFree         transportLayerState = "tlsFree"
	tlsReceiving    transportLayerState = "tlsReceiving"
	tlsTransmitting transportLayerState = "tlsTransmitting"
)

func handleVUFrame(frame Frame) {
	log.Debugln("Handle VU Frame", frame.internalID)
	printFrame(frame)

	switch frame.frameContent[0] {
	case '0':
		// Standard One Frame Message
		length := frame.frameContent[1]
		log.Debugln("Standard Message Length:", length, "Content:", frame.frameContent[2:2+length])
	case '1':
		// Received First Multi Frame Message
		length := frame.frameContent[1:4]
		message := frame.frameContent[2:]
		log.Debugln("Extended Message Length:", length, "First Content:", message)

		// Respond Transport Layer Ready Message
		// TODO

		// Listen To Data Stream Until All Data Received.
		// TODO
	case '2':
		// Receiving Consecutive MultiFrame Message
		//TODO
	case '3':
		// Received Multi Frame Message Accept Message
		// TODO
		// Rest of First Byte: FC:
		// 0: Clear
		// 1: Wait
		// 2: Overflow/abort

	default:
		// Error Unknown Message Type
		// TODO
	}

}

func handleVUContent(frame Frame) {
	content := frame.frameContent
	if (*appState).expectingDataFromVU {
		log.Debugln("Receiving Data from VU")
		handleVUData(frame)
	} else {
		switch content[0:2] {
		case "02":
			log.Debugln("Request")
			handleVUPositiveResponse(frame)
		case "03":
			log.Debugln("Positive Response")
			handleVURequest(frame)
		default:
			log.Debugln("Error: Maybe Unexpected Data")
		}
	}
}

func handleVUData(frame Frame) {
	log.Debugln("Handle VU Data")
}

func handleVUPositiveResponse(frame Frame) {
	log.Debugln("Handle VU Positive Response")
	service := string(frame.frameContent[2:4])
	switch service {
	case "50":
		log.Debugln("DiagnosticSessionControl")
		handleVUDiagnosticSessionControlPositiveResponse(frame)
	case "7E":
		log.Debugln("DiagnosticSessionControl")
		handleVUTesterPresentPositiveResponse(frame)
	case "71":
		log.Debugln("DiagnosticSessionControl")
		handleVURoutineControlPositiveResponse(frame)
	case "75":
		log.Debugln("DiagnosticSessionControl")
		handleVURequestUploadPositiveResponse(frame)
	case "76":
		log.Debugln("DiagnosticSessionControl")
		handleVUTransferDataPositiveResponse(frame)
	case "77":
		log.Debugln("DiagnosticSessionControl")
		handleVURequestTransferExit(frame)
	default:
		log.Debugln("Error Unknown Service")
	}
}

func handleVURequest(frame Frame) {
	log.Debugln("Handle VU Request")
}

func handleVUDiagnosticSessionControlPositiveResponse(frame Frame) {
	log.Debugln("Handle VU DiagnosticSessionControlPositiveResponse")
	data := frame.frameContent[4:8]
	// There might be an errors in the User Guide Version 02.01.181209
	// page 17/69 data = 10 FF, but recorded VU response is 01 FF
	if data != "01FF" && data != "10FF" {
		log.Debugln("Error Expected different data")
		// TODO Optimize reset
		reset(resetLevelTotal)
	}

}
func handleVUTesterPresentPositiveResponse(frame Frame) {
	log.Debugln("Handle VU Tester Present Positive Response")
}
func handleVURoutineControlPositiveResponse(frame Frame) {
	log.Debugln("Handle VU Routine Control Positive Response")
	statusValue := frame.frameContent[4:6]

	switch statusValue {
	case "02":
		// VUReady
		appState.currentState = stateVUReady
	case "04":
		// VUToCompanyCardReader
		appState.currentState = stateVUToCompanyCardData
		// data
		apdu := frame.frameContent[6:16]
		// action
		handleVUToCompanyCardData(apdu)
		// Loops until RemoteAuthenticationSucceeded
	case "06":
		// RemoteAuthenticationSucceeded
		appState.currentState = stateRemoteAuthenticatedSucceeded
		// Trigger RemoteDownloadDataRequest
		// TODO implement
	case "08":
		//RemoteDownloadAccessGranted
		appState.currentState = stateRemoteDownloadAccessGranted
		// Trigger Data Download
		// TODO implement
	case "0A":
		// RemoteAuthenticationClosed
		// VU Confirms that authentication process is ended/terminated
		// or a previous valid authentication closed
		appState.currentState = stateRemoteAuthenticationClosed
		// Action:
		// TODO implement
	case "0C":
		// APDU Error
		// VU informs the company that 3 consecutive APDU errors have occured
		appState.currentState = stateAPDUError
		// Action:
		// TODO implement
	case "0E":
		// Authentication Error
		// VU informs company that the card authentication has failed
		// Also send with expired company card
		appState.currentState = stateAuthenticationError
		// Action:
		// TODO implement
	case "10":
		// ToManyAuthenticationErrors
		// The VU informs company that 5 consecutive card authentication Errors have occured
		// Also send with expired company card
		appState.currentState = stateTooManyAuthenticationErrors
		// Action:
		// TODO implement
	default:
		// Unkown Control Routine Positive Response Status Error
		log.Debugln("Unkown Control Routine Positive Response Status Error")
		// Action:
		// TODO implement
	}
}
func handleVURequestUploadPositiveResponse(frame Frame) {
	log.Debugln("Handle VU Request Upload Positive Response")
	// Transferred Data should be 10 FF
}
func handleVUTransferDataPositiveResponse(frame Frame) {
	log.Debugln("Handle VU Transfer Data Positive Response")
	bsc := frame.frameContent[4:6]
	wac := frame.frameContent[6:8]
	transferparameter := frame.frameContent[8:10]
	data := frame.frameContent[10:16]

	switch transferparameter {
	// Transfer Request Parameter
	case "01":
		// Overview Data
		log.Debugln("Reveived TRTP Overview Data:", data, "BSC:", bsc, "WAC:", wac)
	case "02":
		// Activities (1 day) Data
		log.Debugln("Reveived TRTP Activities (1 day) Data:", data, "BSC:", bsc, "WAC:", wac)
	case "03":
		// Events Faults Data
		log.Debugln("Reveived TRTP Events Fault Data:", data, "BSC:", bsc, "WAC:", wac)
	case "04":
		// Detailed Speed Data
		log.Debugln("Reveived TRTP Detailed Speed Data:", data, "BSC:", bsc, "WAC:", wac)
	case "05":
		// Technical Data
		log.Debugln("Reveived TRTP Technical Data:", data, "BSC:", bsc, "WAC:", wac)
	case "06":
		// Card Download (slot)
		log.Debugln("Reveived TRTP Card Download (Slot) Data:", data, "BSC:", bsc, "WAC:", wac)

	// Transfer Response Parameter
	case "21":
		// Overview Data
		log.Debugln("Reveived TREP Overview Data:", data, "BSC:", bsc, "WAC:", wac)
	case "22":
		// Activities (1 day) Data
		log.Debugln("Reveived TREP Activities (1 day) Data:", data, "BSC:", bsc, "WAC:", wac)
	case "23":
		// Events Faults Data
		log.Debugln("Reveived TREP Events Faults Data:", data, "BSC:", bsc, "WAC:", wac)
	case "24":
		// Detailed Speed Data
		log.Debugln("Reveived TREP Detailed Speed Data:", data, "BSC:", bsc, "WAC:", wac)
	case "25":
		// Technical Data
		log.Debugln("Reveived TREP Technical Data:", data, "BSC:", bsc, "WAC:", wac)
	default:
		// Unknown Data Type Error:
		log.Debugln("Unknown Data Type Error")
	}
}
func handleVURequestTransferExit(frame Frame) {
	log.Debugln("Handle VU Transfer Exit Positive Response")
	trtp := frame.frameContent[4:6]
	if trtp == "00" {
		log.Debugln("RequestTransferExit Success")
		// TODO Change Application State to ?
		// reset()
	} else {
		// Unknown TRTP Error
		log.Debugln("RequestTransferExit Unknown TRTP Error")
		// TODO What now?
	}
}
func handleVUToCompanyCardData(apdu string) {
	log.Debugln("handleVUToCompanyCardData")
	// TODO implement
	// Send Data to Smart Card
	// Receive Smart Card ADPU and send next CompanyCardToVUData
}
