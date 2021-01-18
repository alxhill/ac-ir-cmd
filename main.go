package main

// #cgo LDFLAGS: -lm -lpigpio -pthread -lrt
// #include <stdio.h>
// #include <stdlib.h>
// #include "irslinger.h"
import "C"
import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"unsafe"

	"./state"
	"./sensor"
)

func main() {
	http.HandleFunc("/set", setState)
	http.HandleFunc("/temp", getTemp)

	fmt.Println("Starting server...")

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func setState(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	decoder := json.NewDecoder(r.Body)
	decodedBody := state.AcState{}
	if err := decoder.Decode(&decodedBody); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if !decodedBody.IsValid() {
		fmt.Printf("Invalid state %#v\n", decodedBody)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	fmt.Printf("Requested State %#v\n", decodedBody)

	sendIrCommand(&decodedBody)
}

func getTemp(w http.ResponseWriter, r *http.Request) {
	temp, err := sensor.GetTempCelcius()

	if err != nil {
		fmt.Printf("Failed to get temp due to err %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte(fmt.Sprintf("%f", temp)))

	if err != nil {
		fmt.Printf("Failed to get or write temp")
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
