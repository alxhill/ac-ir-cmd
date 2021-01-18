package sensor

import (
	"gobot.io/x/gobot/drivers/i2c"
	"gobot.io/x/gobot/platforms/raspi"
)

type Celcius int

func GetTempCelcius() (float32, error) {
	r := raspi.NewAdaptor()
	tempSensor := i2c.NewSHT2xDriver(r)

	tempSensor.SetAccuracy(i2c.SHT2xAccuracyHigh)
	temp, err := tempSensor.Temperature()

	if err != nil {
		return 0, err
	}

	return temp, nil
}