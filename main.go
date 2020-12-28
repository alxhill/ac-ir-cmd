package ac_ir_cmd

// #cgo LDFLAGS: -lm -lpigpio -pthread -lrt
// #include <stdio.h>
// #include <stdlib.h>
// #include "irslinger.h"
import "C"
import (
	"fmt"
	"os"
	"unsafe"

	"./state"
)

//const ON_CMD = "1000100000100110110100000101"
//const OFF_CMD = "1000100011000000000001010001"

var defaultAcState = &state.AcState{
	Fan:   state.LOW,
	Mode:  state.COOL,
	Power: state.POWER_ON,
	Temp:  72,
}

func main() {
	if len(os.Args) < 5 {
		fmt.Println("Expecting: fan speed (1-3), mode (1-4), power (0/1), Temp (60-86)")
		os.Exit(1)
	}

	acState := state.NewAcState(os.Args[1], os.Args[2], os.Args[3], os.Args[4])
	fmt.Printf("AC State -> Cmd: %s\n", acState.GetCommand())

	sendIrCommand(acState.GetCommand())
}

func sendIrCommand(commandStr string) {
	commandStrC := C.CString(commandStr)
	defer C.free(unsafe.Pointer(commandStrC))

	outPin := C.uint(17)
	frequency := C.int(38000)
	dutyCycle := C.double(0.5)

	leadingPulseDuration := C.int(9000)
	leadingGapDuration := C.int(4500)
	onePulse := C.int(562)
	zeroPulse := C.int(562)
	oneGap := C.int(1688)
	zeroGap := C.int(562)
	sendTrailingPulse := C.int(1)

	result, err := C.irSling(outPin, frequency, dutyCycle, leadingPulseDuration, leadingGapDuration, onePulse, zeroPulse, oneGap, zeroGap, sendTrailingPulse, commandStrC)

	fmt.Printf("Command ran, result: %d, error: %s\n", int(result), err)
}
