package sensor

import (
	"gobot.io/x/gobot/drivers/i2c"
	"gobot.io/x/gobot/platforms/raspi"
	"log"
)

type Celcius float32
type RelativeHumidity float32

type Sensors struct {
	adaptor     *raspi.Adaptor
	sht2xDriver *i2c.SHT2xDriver
}

func InitSensors() (*Sensors, error) {
	r := raspi.NewAdaptor()
	driver := i2c.NewSHT2xDriver(r)
	if err := driver.Start(); err != nil {
		return nil, err
	}

	if err := driver.SetAccuracy(i2c.SHT2xAccuracyMedium); err != nil {
		return nil, err
	}

	return &Sensors{
		adaptor:     r,
		sht2xDriver: driver,
	}, nil
}

func (s *Sensors) Temperature() (Celcius, error) {
	temp, err := s.sht2xDriver.Temperature()
	if err != nil {
		return 0, err
	}

	log.Printf("Reading temp %f", temp)
	return Celcius(temp), nil
}

func (s *Sensors) Humidity() (RelativeHumidity, error) {
	humidity, err := s.sht2xDriver.Humidity()
	if err != nil {
		return 0, err
	}

	log.Printf("Reading humidity %f", humidity)
	return RelativeHumidity(humidity), nil
}