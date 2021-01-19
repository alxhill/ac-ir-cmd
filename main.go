package main

// #cgo LDFLAGS: -lm -lpigpio -pthread -lrt
// #include <stdio.h>
// #include <stdlib.h>
// #include "irslinger.h"
import "C"
import (
	"alxhill.com/ac-ir-cmd/sensor"
	"alxhill.com/ac-ir-cmd/state"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"unsafe"
)

func main() {
	sensors, err := sensor.InitSensors()
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/set", setState)
	http.HandleFunc("/temp", getTemp(sensors))
	http.HandleFunc("/humidity", getHumidity(sensors))

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

func getTemp(s *sensor.Sensors) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		temp, err := s.Temperature()
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
}

func getHumidity(s *sensor.Sensors) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		humidity, err := s.Humidity()
		if err != nil {
			fmt.Printf("Failed to get humidity due to err %s", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		_, err = w.Write([]byte(fmt.Sprintf("%f", humidity)))

		if err != nil {
			fmt.Printf("Failed to get or write humidity")
		}
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
