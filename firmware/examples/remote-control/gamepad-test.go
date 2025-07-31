package main

import (
	"encoding/binary"
	"fmt"
	"os"
)

const (
	JS_EVENT_BUTTON = 0x01
	JS_EVENT_AXIS   = 0x02
	JS_EVENT_INIT   = 0x80
)

type jsEvent struct {
	Time   uint32
	Value  int16
	Type   uint8
	Number uint8
}

func main() {
	file, err := os.Open("/dev/input/js0")
	if err != nil {
		fmt.Printf("Failed to open joystick: %v\n", err)
		return
	}
	defer file.Close()

	for {
		var e jsEvent
		err := binary.Read(file, binary.LittleEndian, &e)
		if err != nil {
			fmt.Printf("Failed to read: %v\n", err)
			break
		}

		eventType := e.Type & (^(uint8(JS_EVENT_INIT))) // mask out JS_EVENT_INIT bit

		switch eventType {
		case JS_EVENT_BUTTON:
			state := "released"
			if e.Value != 0 {
				state = "pressed"
			}
			fmt.Printf("Button %d %s\n", e.Number, state)
		case JS_EVENT_AXIS:
			fmt.Printf("Axis %d moved to %d\n", e.Number, e.Value)
		default:
			fmt.Printf("Unknown event type: %02x\n", e.Type)
		}
	}
}
