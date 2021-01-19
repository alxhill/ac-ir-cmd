package sensor

import (
	"gobot.io/x/gobot/drivers/i2c"
	"gobot.io/x/gobot/platforms/raspi"
)

type Celcius int

func GetTempCelcius() (float32, error) {
	r := raspi.NewAdaptor()
	tempSensor := i2c.NewSHT2xDriver(r, i2c.WithAddress(0x40))
	if err := tempSensor.Start(); err != nil {
		return 0, err
	}

	if err := tempSensor.SetAccuracy(i2c.SHT2xAccuracyHigh); err != nil {
		return 0, err
	}

	temp, err := tempSensor.Temperature()
	if err != nil {
		return 0, err
	}

	return temp, nil
}
