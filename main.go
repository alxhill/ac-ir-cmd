package ac_ir_cmd

// #cgo LDFLAGS: -lm -lpigpio -pthread -lrt
// #include <stdio.h>
// #include <stdlib.h>
// #include "irslinger.h"
import "C"
import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
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
	//if len(os.Args) < 5 {
	//	fmt.Println("Expecting: fan speed (1-3), mode (1-4), power (0/1), Temp (60-86)")
	//	os.Exit(1)
	//}
	//
	//acState := state.NewAcState(os.Args[1], os.Args[2], os.Args[3], os.Args[4])
	//fmt.Printf("AC State -> Cmd: %s\n", acState.GetCommand())
	//
	//sendIrCommand(acState)

	http.HandleFunc("/set", setState)
	http.HandleFunc("/temp", getTemp)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func setState(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	reader, err := r.GetBody()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	requestBody, err := ioutil.ReadAll(reader)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	decodedBody := state.AcState{}

	if err := json.Unmarshal(requestBody, &decodedBody); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	fmt.Printf("Requested Command: ", decodedBody.GetCommand())

	sendIrCommand(&decodedBody)
}

func getTemp(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	_, err := w.Write([]byte("18"))
	if err != nil {
		fmt.Printf("Failed to write temp")
	}
}

func sendIrCommand(acState *state.AcState) {
	commandStr := acState.GetCommand()
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
