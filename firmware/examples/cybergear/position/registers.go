package main

type CommunicationType uint8

// CyberGear communication types
const (
	Broadcast                        CommunicationType = 0x00
	Control                          CommunicationType = 0x01
	Feedback                         CommunicationType = 0x02
	EnableOperation                  CommunicationType = 0x03
	StopOperation                    CommunicationType = 0x04
	SetMechanicalZero                CommunicationType = 0x06
	SetCANId                         CommunicationType = 0x07
	UndocumentedSingleParameterWrite                   = 0x08 // For accessing registers 0x0000 - 0x302F. Undocumented and untested. Run at your own risk.
	UndocumentedSingleParameterRead                    = 0x09 // For accessing registers 0x0000 - 0x302F. Undocumented and untested. Run at your own risk.
	SingleParameterRead              CommunicationType = 0x11 // For 0x70xx registers only
	SingleParameterWrite             CommunicationType = 0x12 // For 0x70xx registers only
	FaultFeedback                    CommunicationType = 0x15
	SetBaudRate                      CommunicationType = 0x16
)

type OperationMode uint8

// CyberGear operation modes
const (
	OperationControlMode OperationMode = 0
	PositionMode         OperationMode = 1
	SpeedMode            OperationMode = 2
	CurrentMode          OperationMode = 3
)

// CyberGear registers
const (
	// uint8, R/W
	// 0: Operation control mode
	// 1: Position mode
	// 2: Speed mode
	// 3: Current mode
	runModeRegister uint16 = 0x7005
	// float32, R/W
	// Range: -23 ~ 23A
	iqRefRegister uint16 = 0x7006
	// float32, R/W
	// Range: -30 ~ 30rad/s
	speedRefRegister uint16 = 0x700A
	// float32, R/W
	// Range: 0~12Nm
	torqueLimitRegister uint16 = 0x700B
	// float32, R/W
	// Default value: 0.125
	currentKpRegister uint16 = 0x7010
	// float32, R/W
	// Default value: 0.0158
	currentKiRegisteruint16 = 0x7011
	// float32, R/W
	// Default value: 0.1
	// Range 0~1.0
	currentFilterGainRegister uint16 = 0x7014
	// float32, R/W
	// Range: radians (in position mode)
	locationRefRegister uint16 = 0x7016
	// float32, R/W
	// Range: 0 ~ 30rad/s (in position mode)
	limitSpeedRegister uint16 = 0x7017
	// float32, R/W
	// Range: 0 ~ 23A (speed / position mode
	limitCurrentRegister uint16 = 0x7018
	// float 32, R
	// Range: radians
	mechPosRegister uint16 = 0x7019
	// float32, R
	// Range: -23 ~ 23A
	iqFilterValueRegister uint16 = 0x701A
	// float32, R
	// Range: -30 ~ 30rad/s
	mechVelocityRegister uint16 = 0x701B
	// float32, R
	// Range: V
	vBusRegister uint16 = 0x701C
	// int16, R/W
	// Range: Number of turns
	rotationRegister uint16 = 0x701D
	// float32, R/W
	// Default value: 30
	positionKpRegister uint16 = 0x701E
	// float32, R/W
	// Default value: 1
	speedKpRegister uint16 = 0x701F
	// float32, R/W
	// Default value: 0.002
	speedKiRegister uint16 = 0x7020
)
