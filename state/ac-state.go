package state

import (
	"fmt"
	"strconv"
)

type FanSpeed string
type Mode string
type PowerState string
type Fahrenheit uint8

const (
	LOW    FanSpeed = "low"
	MEDIUM          = "medium"
	HIGH            = "high"
)

const (
	FAN         Mode = "fan"
	COOL             = "cool"
	DRY              = "dry"
	MONEY_SAVER      = "money-saver"
)

const (
	POWER_ON  PowerState = "on"
	POWER_OFF            = "off"
)

type AcState struct {
	Fan   FanSpeed   `json:"fan"`
	Mode  Mode       `json:"mode"`
	Power PowerState `json:"power"`
	Temp  Fahrenheit `json:"temp"`
}

func (fan FanSpeed) validate() bool {
	switch fan {
	case LOW:
	case MEDIUM:
	case HIGH:
		return true
	}
	return false
}

func (fan FanSpeed) fanBinary() uint32 {
	switch fan {
	case LOW:
		return 0x0
	case MEDIUM:
		return 0x2
	case HIGH:
		return 0x4
	}
	return 0x0
}

func (mode Mode) validate() bool {
	switch mode {
	case COOL:
	case DRY:
	case FAN:
	case MONEY_SAVER:
		return true
	}
	return false
}

func (mode Mode) modeBinary() uint32 {
	switch mode {
	case COOL:
		return 0x0
	case DRY:
		return 0x1
	case FAN:
		return 0x2
	case MONEY_SAVER:
		return 0x6
	}
	return 0x6
}

func (temp Fahrenheit) validate() bool {
	//return temp <= 75 && temp >= 65
	return true
}

func (temp Fahrenheit) tempBinary() uint32 {
	if temp < 75 {
		var cmdTemp = uint32(temp) - 59
		return cmdTemp << 1
	}
	var cmdTemp = uint32(temp) - 75
	return cmdTemp<<1 | 1
}

func (power PowerState) validate() bool {
	return power == POWER_OFF || power == POWER_ON
}

func (power PowerState) powerBinary() uint32 {
	if power == POWER_OFF {
		panic("not implemented")
	}
	return 0x0
}

func (ac AcState) GetCommand() string {
	if ac.Power == POWER_OFF {
		return "1000100011000000000001010001"
	}

	var startBits uint32 = 0x882

	blockOne := ac.Power.powerBinary()<<3 | ac.Mode.modeBinary()
	blockTwo := ac.Temp.tempBinary() >> 1
	blockThree := (ac.Temp.tempBinary()&1)<<3 | ac.Fan.fanBinary()

	checkSum := (2 + blockOne + blockTwo + blockThree) & 0xF

	cmd := startBits<<16 | blockOne<<12 | blockTwo<<8 | blockThree<<4 | checkSum

	return strconv.FormatUint(uint64(cmd), 2)
}

func (ac AcState) IsValid() bool {
	fmt.Printf("%b %b %b %b", ac.Mode.validate(), ac.Fan.validate(), ac.Power.validate(), ac.Temp.validate())
	return ac.Mode.validate() &&
		ac.Fan.validate() &&
		ac.Power.validate() &&
		ac.Temp.validate()
}

func NewAcState(fan string, mode string, power string, temp string) *AcState {
	fanUint, _ := strconv.ParseUint(fan, 10, 32)
	modeUint, _ := strconv.ParseUint(mode, 10, 32)
	powerUint, _ := strconv.ParseUint(power, 10, 32)
	tempUint, _ := strconv.ParseUint(temp, 10, 32)
	return &AcState{
		FanSpeed(fanUint),
		Mode(modeUint),
		PowerState(powerUint),
		Fahrenheit(tempUint),
	}
}
