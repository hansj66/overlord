package main

import (
	"log/slog"
	"net"
	"os"
	"time"

	"golang.org/x/sys/unix"
)

const (
	hostID  = 0x01
	motorID = 0x7F
)

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
