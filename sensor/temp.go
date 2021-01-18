package sensor

import (
	"fmt"
	"log"
	"periph.io/x/conn/v3/i2c"
	"periph.io/x/conn/v3/i2c/i2creg"
)

type Celcius int

func GetTempCelcius() (float32, error) {
//	//r := raspi.NewAdaptor()
//	//tempSensor := i2c.NewSHT2xDriver(r)
//	//
//	//tempSensor.SetAccuracy(i2c.SHT2xAccuracyHigh)
//	//temp, err := tempSensor.Temperature()
//	if _, err := host.Init(); err != nil {
//		return 0, err
//	}
//
	bus, err := i2creg.Open("")
	defer bus.Close()
	if err != nil {
		return 0, err
	}

	// Dev is a valid conn.Conn.
	d := &i2c.Dev{Addr: 40, Bus: bus}

//	// Send a command 0x10 and expect a 5 bytes reply.
	write := []byte{0x10}
	read := make([]byte, 5)
	if err := d.Tx(write, read); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%v\n", read)

	return 0.0, nil
}
