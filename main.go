package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/jacobsa/go-serial/serial"
	log "github.com/sirupsen/logrus"
)

var debug *bool = flag.Bool("debug", false, "debug mode")
var serialPort *string = flag.String("port", "/dev/ttyACM3", "serial port name")

const cardServerAddress string = "http://localhost:8080"

func init() {
	flag.Parse()
	if *debug {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}
}
func main() {

	//
	// Open Port
	port, err := serial.Open(serialOptions)

	if err != nil {
		log.Fatalf("Serial.Open: %v", err)
	}
	log.Debugf("Serial port Open!\n")
	defer port.Close()

	// Configure USBtin
	err = configureUSBtin(port)
	if err != nil {
		log.Fatalf("USBtin configuration failed: %v", err)
	}
	log.Debugf("USBtin ready!\n")

	// Step 1: reset state

	var wg sync.WaitGroup
	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		for {
			resp, err := readFrame(port)
			if err != nil {
				log.Errorln(err)
			}
			log.Debugln("Incoming Message", resp)
			if strings.Contains(resp, PositiveResponseRemoteAuthenticationExit) {
				fmt.Println("DCTO Connected")
				time.Sleep(time.Second)
				os.Exit(0)
			}
		}
	}(&wg)

	startRemoteSession(port)
	time.Sleep(time.Second)
	closeAuthenticationProcess(port)
	time.Sleep(time.Second)
	endRemoteSession(port)
	time.Sleep(3 * time.Second)
	fmt.Println("DCTO Not Found")

}
