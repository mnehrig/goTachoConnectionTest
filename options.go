package main

import (
	"github.com/jacobsa/go-serial/serial"
)

var (
	serialOptions serial.OpenOptions = serial.OpenOptions{
		PortName:        *serialPort,
		BaudRate:        115200,
		DataBits:        8,
		StopBits:        1,
		MinimumReadSize: 1,
		//InterCharacterTimeout: uint(time.Second),
	}
)
