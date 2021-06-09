package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os/exec"

	"alxhill.com/ac-ir-cmd/sensor"
	"alxhill.com/ac-ir-cmd/state"
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
		if r.Method != "GET" {
			w.WriteHeader((http.StatusMethodNotAllowed))
			return
		}

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
	cmd := exec.Command("irslinger", commandStr)
	err := cmd.Run()

	if err != nil {
		fmt.Printf("Command had error %s\n", err)
	} else {
		fmt.Println("Sent IR command successfully")
	}
}
