package main

import (
	"encoding/binary"
	"fmt"
	"log/slog"
	"math"

	"golang.org/x/sys/unix"
)

// From <linux/can.h>: bit 31 = EFF flag for 29-bit IDs
const CAN_EFF_FLAG = 0x80000000

type CyberGearFunction struct {
	hostID            uint8
	motorID           uint8
	communicationType CommunicationType
	parameterData     [8]byte
}

// NewCyberGearFunction sets up a new cybergear function (enable / disable, speed / position mode, register access etc)
func NewCyberGearFunction(communicationType, hostID, motorID uint8, parameterData [8]byte) CyberGearFunction {
	return CyberGearFunction{
		hostID:            hostID,
		motorID:           motorID,
		communicationType: CommunicationType(communicationType),
		parameterData:     parameterData,
	}
}

// marshalToRawCANFrame serializes the CyberGearFunctio to a raw extended CAN frame
func (f *CyberGearFunction) marshalToRawCANFrame() []byte {
	rawCANFrame := make([]byte, unix.CAN_MTU)

	header := uint32(f.communicationType)<<24 |
		uint32(f.hostID)<<8 |
		uint32(f.motorID) |
		CAN_EFF_FLAG
	binary.LittleEndian.PutUint32(rawCANFrame[0:4], header)
	rawCANFrame[4] = 8 // DLC
	copy(rawCANFrame[8:16], f.parameterData[:])

	// fmt.Printf("RAW CAN frame: %x (parameter data: %x)\n", rawCANFrame, f.parameterData)
	return rawCANFrame
}

// execute sends a serialized CAN frame representation of the CyberGear function and writes this to a socket.
func (f *CyberGearFunction) execute(socket int) error {

	_, err := unix.Write(socket, f.marshalToRawCANFrame())
	if err != nil {
		slog.Error("write:", "err", err)
		return err
	}

	readAndDecodeMotorResponseFrame(socket)
	return nil
}

// EnableMotor enables the CyberGear operation. This has to be called before any speed / positioning commands are sent
func EnableMotor(socket int, hostID, motorID uint8) error {
	var parameterData [8]byte
	// binary.LittleEndian.PutUint16(parameterData[0:4], VBUSParam)
	function := NewCyberGearFunction(uint8(EnableOperation), hostID, motorID, parameterData)
	if err := function.execute(socket); err != nil {
		slog.Error("EnableMotor failed", "err", err)
		return err
	}
	slog.Info("Motor enabled")
	return nil
}

// Disable motor stops the motor and shuts down.
func DisableMotor(socket int, hostID, motorID uint8) error {
	var parameterData [8]byte
	// binary.LittleEndian.PutUint16(parameterData[0:4], VBUSParam)
	function := NewCyberGearFunction(uint8(StopOperation), hostID, motorID, parameterData)
	if err := function.execute(socket); err != nil {
		slog.Error("EnableMotor failed", "err", err)
		return err
	}
	slog.Info("Motor disabled")
	return nil
}

// SetOperationMode selects operating mode (speed, postition, current, control)
func SetOperationMode(socket int, hostID, motorID uint8, mode OperationMode) error {
	var parameterData [8]byte
	binary.LittleEndian.PutUint16(parameterData[0:4], runModeRegister)
	parameterData[4] = uint8(mode)
	function := NewCyberGearFunction(uint8(SingleParameterWrite), hostID, motorID, parameterData)
	if err := function.execute(socket); err != nil {
		slog.Error("EnableMotor failed", "err", err)
		return err
	}
	slog.Info(fmt.Sprintf("Operation mode is %d", mode))
	return nil
}

// SetSpeed sets the speed of the motor in radians/s
func SetSpeed(socket int, hostID, motorID uint8, speed float32) error {
	var parameterData [8]byte
	binary.LittleEndian.PutUint16(parameterData[0:4], speedRefRegister)
	binary.LittleEndian.PutUint32(parameterData[4:8], math.Float32bits(speed))
	function := NewCyberGearFunction(uint8(SingleParameterWrite), hostID, motorID, parameterData)
	if err := function.execute(socket); err != nil {
		slog.Error("EnableMotor failed", "err", err)
		return err
	}
	slog.Info("SetSpeed", "speed (rad/s)", speed)
	return nil
}

// SetPosition sets the position of the motor in radians
func SetPosition(socket int, hostID, motorID uint8, position float32) error {
	var parameterData [8]byte
	binary.LittleEndian.PutUint16(parameterData[0:4], locationRefRegister)
	binary.LittleEndian.PutUint32(parameterData[4:8], math.Float32bits(position))
	function := NewCyberGearFunction(uint8(SingleParameterWrite), hostID, motorID, parameterData)
	if err := function.execute(socket); err != nil {
		slog.Error("EnableMotor failed", "err", err)
		return err
	}
	slog.Info("SetPosition", "position", position)
	return nil
}
