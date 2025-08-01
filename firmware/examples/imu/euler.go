package main

import (
	"encoding/binary"
	"fmt"
	"log/slog"
	"os"
	"time"

	"periph.io/x/conn/v3/driver/driverreg"
	"periph.io/x/conn/v3/i2c"
	"periph.io/x/conn/v3/i2c/i2creg"
	"periph.io/x/host/v3"
)

const (
	BNO055_ADDR           = 0x28
	BNO055_OPR_MODE_ADDR  = 0x3D
	BNO055_EULER_H_LSB    = 0x1A
	OPERATION_MODE_CONFIG = 0x00
	OPERATION_MODE_NDOF   = 0x0C
)

func main() {
	// Initialize periph
	if _, err := host.Init(); err != nil {
		slog.Error("Unable to initialize periph library", "err", err)
		os.Exit(1)
	}

	if _, err := driverreg.Init(); err != nil {
		slog.Error("Unable to initialize i2c driver")
		os.Exit(1)
	}

	// Open the I²C bus
	bus, err := i2creg.Open("/dev/i2c-1")
	if err != nil {
		slog.Error("Unable to access I2C bus", "err", err)
		os.Exit(1)
	}
	defer bus.Close()

	dev := &i2c.Dev{Bus: bus, Addr: BNO055_ADDR}

	// Sanity check
	if err := checkSensorID(dev); err != nil {
		slog.Error("Error checking IMU id", "err", err)
		os.Exit(1)
	}

	// Set to CONFIG mode
	writeRegister(dev, BNO055_OPR_MODE_ADDR, OPERATION_MODE_CONFIG)
	time.Sleep(20 * time.Millisecond)

	// Set to NDOF mode
	writeRegister(dev, BNO055_OPR_MODE_ADDR, OPERATION_MODE_NDOF)
	time.Sleep(600 * time.Millisecond) // Let sensor settle

	// Read Euler angles
	for {
		buf := make([]byte, 6)
		if err := dev.Tx([]byte{BNO055_EULER_H_LSB}, buf); err != nil {
			slog.Error("I2C Tx error", "err", err)
			os.Exit(1)
		}

		heading := float32(binary.LittleEndian.Uint16(buf[0:2])) / 16.0
		roll := float32(binary.LittleEndian.Uint16(buf[2:4])) / 16.0
		pitch := float32(binary.LittleEndian.Uint16(buf[4:6])) / 16.0

		slog.Info("BNO055", "heading", fmt.Sprintf("Heading: %.2f°, Roll: %.2f°, Pitch: %.2f°\n", heading, roll, pitch))
		time.Sleep(500 * time.Millisecond)
	}
}

func writeRegister(dev *i2c.Dev, reg, value byte) {
	if err := dev.Tx([]byte{reg, value}, nil); err != nil {
		slog.Error("Failed to write register", "reg", reg, "err", err)
	}
}

func checkSensorID(dev *i2c.Dev) error {
	const BNO055_CHIP_ID_ADDR = 0x00
	const BNO055_CHIP_ID_EXPECTED = 0xA0

	buf := []byte{0}
	if err := dev.Tx([]byte{BNO055_CHIP_ID_ADDR}, buf); err != nil {
		slog.Error("Failed to read chip ID", "err", err)
		return err
	}

	if buf[0] != BNO055_CHIP_ID_EXPECTED {
		err := fmt.Errorf("unexpected chip ID:", "got", buf[0], "expected", BNO055_CHIP_ID_EXPECTED)
		slog.Error("CHIP ID", "err", err)
		return err
	}

	slog.Info("BNO055 sensor detected, chip ID is correct (0xA0)")
	return nil
}
