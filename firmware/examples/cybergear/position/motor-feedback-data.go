package main

import (
	"encoding/binary"
	"fmt"
	"log/slog"

	"golang.org/x/sys/unix"
)

// Fault bit interpretations
const (
	noFault = 0
	fault   = 1
)

// Fault conditions
const (
	notCalibrated           = 1 << 21
	hallEncodingFailure     = 1 << 20
	magneticEncodingFailure = 1 << 19
	overTemperature         = 1 << 18
	overcurrent             = 1 << 17
	undervoltage            = 1 << 16
)

// modes
const (
	resetMode       = 0
	calibrationMode = 1
	runMode         = 2
)

func readAndDecodeMotorResponseFrame(fd int) {
	// A CAN frame is always 16 bytes in classic CAN
	buf := make([]byte, unix.CAN_MTU)
	n, err := unix.Read(fd, buf)
	if err != nil {
		slog.Error("read error:", "err", err)
		return
	}
	if n != unix.CAN_MTU {
		slog.Error("incomplete CAN frame:", "received bytes", n)
		return
	}

	header := binary.LittleEndian.Uint32(buf[0:4])
	motorId := (header >> 8) & 0xFF
	if header&notCalibrated == notCalibrated {
		slog.Warn("Motor is not calibrated")
	}
	if header&hallEncodingFailure == hallEncodingFailure {
		slog.Error("Hall encoder failure")
	}
	if header&magneticEncodingFailure == magneticEncodingFailure {
		slog.Error("Magnetic encoding failure")
	}
	if header&overTemperature == overTemperature {
		slog.Warn("Overtemperature")
	}
	if header&overcurrent == overcurrent {
		slog.Warn("Overcurrent")
	}
	if header&undervoltage == undervoltage {
		slog.Warn("Undervoltage")
	}

	mode := (header >> 22) & 0x03
	switch mode {
	case resetMode:
		slog.Debug(fmt.Sprintf("Motor %02x is in reset mode", motorId))
	case calibrationMode:
		slog.Debug(fmt.Sprintf("Motor %02x is in calibration mode", motorId))
	case runMode:
		slog.Debug(fmt.Sprintf("Motor %02x is in run mode", motorId))
	default:
		slog.Error("Unknown motor mode", "id", motorId, "mode", mode)
	}
}
