package main

import (
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/tarm/serial"
	"periph.io/x/conn/v3/driver/driverreg"
	"periph.io/x/conn/v3/gpio"
	"periph.io/x/conn/v3/gpio/gpioreg"
	"periph.io/x/host/v3"
)

/*
This example demonstrates how to scan the valid servo ID range. This is done by sending a ping packet
to all addresses between 0 and 253. A response packet indicates that a servo has been found in the network

Assuming that the network contains one servo with ID=1, the output will be on the form:

2025/08/01 14:09:41 INFO Opening serialport port=/dev/ttyAMA0
2025/08/01 14:09:41 INFO Serial port opened successfully
2025/08/01 14:09:41 INFO DXL DIR pin=GPIO22
PING id:0 - no response
PING id:1 - response packet: FFFF010200FC
PING id:2 - no response
PING id:3 - no response
...

The same protocol can be used for reading and writing to EEPROM and RAM control table parameters
Protocol specifications:
 - 1.0: https://emanual.robotis.com/docs/en/dxl/protocol1/
 - 2.0: https://emanual.robotis.com/docs/en/dxl/protocol2/

This example is using Dynamixel protocol 1.0 for sending ping packets.

Note:
1) The carrier board has to be powered through the "COMPUTE POWAH" connector (5-36V)
2) The dynamixel servos has to be powered through the "DYNAMIXEL POWAH". This power rail only feeds the
   servos via the dynamixel bus. Please refer to the servo manuals regarding min/max voltages. Exceeding
   the rated servo voltages will permanently damage your servos.
*/

const (
	DxlDirPin      = "GPIO22"
	SerialPortName = "/dev/ttyAMA0"
	MaxId          = 253
	// ServoReturnDelayTime is configurable for each servo. Refer to the servo manual for details.
	ServoReturnDelayTime = 250 * time.Microsecond
	LineDriverDelay      = 10 * time.Microsecond
)

func main() {
	host.Init()

	// Initialize Periph library
	if _, err := driverreg.Init(); err != nil {
		slog.Error("Could not initialize peripheral", "err", err)
		os.Exit(1)
	}

	p := gpioreg.ByName("GPIO22")
	if nil == p {
		slog.Error("Could not get GPIO pin")
		os.Exit(1)
	}
	if err := p.Out(gpio.Low); err != nil {
		slog.Error("Unable to configure GPIO pin", "err", err)
		os.Exit(1)
	}

	Scan()
}

func Scan() {
	c := &serial.Config{
		Name:        SerialPortName,
		Baud:        57600,                // Baud rate is a parameter in the servo control table
		ReadTimeout: ServoReturnDelayTime, // Return delay time is a parameter in the servo control table
		Size:        8,
		Parity:      serial.ParityNone,
		StopBits:    1,
	}
	port, err := serial.OpenPort(c)
	if err != nil {
		slog.Error("Error opening serial port", "name", SerialPortName, "err", err)
		os.Exit(1)
	}

	slog.Info("Opening serialport", "port", SerialPortName)
	slog.Info("Serial port opened successfully")
	pin := gpioreg.ByName("GPIO22") // DXL_DIR pin
	if pin == nil {
		slog.Error("Failed to find GPIO pin", "pin", DxlDirPin)
		os.Exit(1)
	}
	slog.Info("DXL DIR", "pin", DxlDirPin)

	// Ping all servos in the valid ID range
	pingPacket := make([]byte, 10)
	responsePacket := make([]byte, 11)
	for id := 0; id < MaxId; id++ {
		pingPacket[0] = 0xFF                                             // Header
		pingPacket[1] = 0xFF                                             // Header
		pingPacket[2] = byte(id)                                         // target ID
		pingPacket[3] = 0x02                                             // length
		pingPacket[4] = 0x01                                             // ping instruction
		pingPacket[5] = ^(pingPacket[2] + pingPacket[3] + pingPacket[4]) //Checksum

		if err := pin.Out(gpio.High); err != nil {
			slog.Error("Failed to pull gpio pin high", "pin", DxlDirPin)
			os.Exit(1)
		}
		time.Sleep(LineDriverDelay)
		if _, err := port.Write(pingPacket); err != nil {
			slog.Error("Unable to write to serial port", "err", err)
			os.Exit(1)
		}
		time.Sleep(LineDriverDelay)
		if err := pin.Out(gpio.Low); err != nil {
			slog.Error("Failed to pull gpio pin low", "pin", DxlDirPin)
			os.Exit(1)
		}

		n, _ := port.Read(responsePacket)
		if n == 6 {
			fmt.Printf("PING id:%d - response packet: %X\n", id, responsePacket[:n])
		} else {
			fmt.Printf("PING id:%d - no response\n", id)
		}
	}
}
