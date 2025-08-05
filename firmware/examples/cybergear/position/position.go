package main

import (
	"encoding/binary"
	"log/slog"
	"math"
	"net"
	"os"
	"time"

	"golang.org/x/sys/unix"
)

const (
	hostID  = 0x01
	motorID = 0x7F
)

func createSetControlModeFrame(hostID, motorID, mode uint8) []byte {
	// Build the 29-bit CAN ID
	canID := uint32(SingleParameterWrite)<<24 |
		uint32(hostID)<<8 |
		uint32(motorID) |
		CAN_EFF_FLAG

	// Allocate the frame and pack ID+flags (little-endian)
	frame := make([]byte, unix.CAN_MTU) // CAN_MTU == 16
	binary.LittleEndian.PutUint32(frame[0:4], canID)
	frame[4] = 8 // DLC = 8 bytes
	// Payload: index 0x7005 → bytes 8–9
	frame[8] = 0x70
	frame[9] = 0x05
	// Insert the 1-byte mode value at byte 12
	frame[12] = mode

	return frame
}

func createSetPositionFrame(hostID, motorID uint8, angleDeg float64) []byte {
	// Build the 29-bit CAN ID
	canID := uint32(SingleParameterWrite)<<24 |
		uint32(hostID)<<8 |
		uint32(motorID) |
		CAN_EFF_FLAG

	// 2) Angle: degrees → radians → IEEE-754 float32
	angleRad := float32(angleDeg * math.Pi / 180.0)
	bits := math.Float32bits(angleRad)

	frame := make([]byte, unix.CAN_MTU)
	binary.LittleEndian.PutUint32(frame[0:4], canID) // ID+flags
	frame[4] = 8                                     // DLC = 8

	// 4) Register index 0x7016  (MSB first)
	frame[8] = 0x16
	frame[9] = 0x70

	// 5) Angle payload (little-endian) at bytes 12–15
	binary.LittleEndian.PutUint32(frame[12:16], bits)

	return frame
}

func RunSpeedExample(fd int) {
	slog.Info("--- Running speed example ---")
	slog.Info("(Running 1 rad/s for 3 seconds and -2 rad/s for 3 seconds)")
	EnableMotor(fd, hostID, motorID)
	SetOperationMode(fd, hostID, motorID, SpeedMode)
	SetSpeed(fd, hostID, motorID, 1)
	time.Sleep(3 * time.Second)
	SetSpeed(fd, hostID, motorID, -2)
	time.Sleep(3 * time.Second)
	DisableMotor(fd, hostID, motorID)
	slog.Info("--- Speed example finished")
}

func RunPositionExample(fd int) {
	slog.Info("--- Running positioning example ---")
	slog.Info("(flipping between 0 and +90 degrees 5 times and pausing for each second between each position change)")
	EnableMotor(fd, hostID, motorID)
	SetOperationMode(fd, hostID, motorID, PositionMode)
	for i := 0; i < 5; i++ {
		SetPosition(fd, hostID, motorID, 0 /* radians */)
		time.Sleep(1 * time.Second)
		SetPosition(fd, hostID, motorID, 1.5708 /* radians */)
		time.Sleep(1 * time.Second)
	}
	DisableMotor(fd, hostID, motorID)
	slog.Info("--- Positioning example finished")
}

func main() {
	// 1) Open a PF_CAN/SOCK_RAW socket
	s, err := unix.Socket(unix.AF_CAN, unix.SOCK_RAW, unix.CAN_RAW)
	if err != nil {
		slog.Error("socket", "err", err)
		os.Exit(1)
	}
	defer unix.Close(s)

	// 2) Look up the interface index for "can0"
	iface, err := net.InterfaceByName("can0")
	if err != nil {
		slog.Error("InterfaceByName", "err", err)
		os.Exit(1)
	}

	// 3) Bind the socket to can0
	addr := &unix.SockaddrCAN{Ifindex: iface.Index}
	if err := unix.Bind(s, addr); err != nil {
		slog.Error("bind:", "err", err)
		os.Exit(1)
	}

	RunSpeedExample(s)
	RunPositionExample(s)
}
