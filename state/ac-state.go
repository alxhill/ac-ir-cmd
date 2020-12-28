package state


import (
	"strconv"
)

type FanSpeed uint8
type Mode uint8
type PowerState uint8
type Farenheit uint8

const (
	LOW FanSpeed = iota
	MEDIUM
	HIGH
)

const (
	FAN Mode = iota
	COOL
	DRY
	MONEY_SAVER
)

const (
	POWER_OFF PowerState = iota
	POWER_ON
)

type AcState struct {
	Fan   FanSpeed
	Mode  Mode
	Power PowerState
	Temp  Farenheit
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

func (temp Farenheit) tempBinary() uint32 {
	if temp < 75 {
		var cmdTemp uint32 = uint32(temp) - 59
		return cmdTemp << 1
	}
	var cmdTemp uint32 = uint32(temp) - 75
	return cmdTemp<<1 | 1
}

func (PowerState) powerBinary() uint32 {
	return 0x0
}

func (ac AcState) GetCommand() string {
	//if ac.Power == POWER_OFF {
	//	return "1000100011000000000001010001"
	//}

	var startBits uint32 = 0x882

	blockOne := ac.Power.powerBinary()<<3 | ac.Mode.modeBinary()
	blockTwo := ac.Temp.tempBinary() >> 1
	blockThree := (ac.Temp.tempBinary()&1)<<3 | ac.Fan.fanBinary()

	checkSum := (2 + blockOne + blockTwo + blockThree) & 0xF

	cmd := startBits<<16 | blockOne<<12 | blockTwo<<8 | blockThree<<4 | checkSum

	return strconv.FormatUint(uint64(cmd), 2)
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
		Farenheit(tempUint),
	}
}
