#ifndef IRSLINGER_H
#define IRSLINGER_H

#include <string.h>
#include <math.h>
#include <stdio.h>
#include <pigpio.h>

#define MAX_COMMAND_SIZE 512
#define MAX_PULSES 12000

static inline void addPulse(uint32_t onPins, uint32_t offPins, uint32_t duration, gpioPulse_t *irSignal, unsigned int *pulseCount)
{
	unsigned int index = *pulseCount;

	irSignal[index].gpioOn = onPins;
	irSignal[index].gpioOff = offPins;
	irSignal[index].usDelay = duration;

	(*pulseCount)++;
}

// Generates a square wave for duration (microseconds) at frequency (Hz)
// on GPIO pin outPin. dutyCycle is a floating value between 0 and 1.
static inline void carrierFrequency(uint32_t outPin, double frequency, double dutyCycle, double duration, gpioPulse_t *irSignal, unsigned int *pulseCount)
{
	double oneCycleTime = 1000000.0 / frequency; // 1000000 microseconds in a second
	int onDuration = (int)round(oneCycleTime * dutyCycle);
	int offDuration = (int)round(oneCycleTime * (1.0 - dutyCycle));

	int totalCycles = (int)round(duration / oneCycleTime);
	int totalPulses = totalCycles * 2;

	int i;
	for (i = 0; i < totalPulses; i++)
	{
		if (i % 2 == 0)
		{
			// High pulse
			addPulse(1 << outPin, 0, onDuration, irSignal, pulseCount);
		}
		else
		{
			// Low pulse
			addPulse(0, 1 << outPin, offDuration, irSignal, pulseCount);
		}
	}
}

// Generates a low signal gap for duration, in microseconds, on GPIO pin outPin
static inline void gap(double duration, gpioPulse_t *irSignal, unsigned int *pulseCount)
{
	addPulse(0, 0, duration, irSignal, pulseCount);
}

// Transmit generated wave
static inline int transmitWave(uint32_t outPin, gpioPulse_t *irSignal, unsigned int *pulseCount)
{
	// Init pigpio
	if (gpioInitialise() < 0)
	{
		// Initialization failed
		printf("GPIO Initialization failed\n");
		return 1;
	}

	// Setup the GPIO pin as an output pin
	gpioSetMode(outPin, PI_OUTPUT);

	// Start a new wave
	gpioWaveClear();

	gpioWaveAddGeneric(*pulseCount, irSignal);
	int waveID = gpioWaveCreate();

	if (waveID >= 0)
	{
		int result = gpioWaveTxSend(waveID, PI_WAVE_MODE_ONE_SHOT);

		printf("Result: %i\n", result);
	}
	else
	{
		printf("Wave creation failure!\n %i", waveID);
	}

	// Wait for the wave to finish transmitting
	while (gpioWaveTxBusy())
	{
		time_sleep(0.1);
	}
	printf("Wave finished transmitting %i\n", waveID);

	// Delete the wave if it exists
	if (waveID >= 0)
	{
	    printf("Deleted wave %i\n", waveID);
		gpioWaveDelete(waveID);
	}

    printf("Cleaning up\n");
	// Cleanup
	gpioTerminate();
	return 0;
}

static inline int irSling(uint32_t outPin,
	int frequency,
	double dutyCycle,
	int leadingPulseDuration,
	int leadingGapDuration,
	int onePulse,
	int zeroPulse,
	int oneGap,
	int zeroGap,
	int sendTrailingPulse,
	const char *code)
{
	if (outPin > 31)
	{
		// Invalid pin number
		return 1;
	}

	size_t codeLen = strlen(code);

	printf("code size is %zu\n", codeLen);

	if (codeLen > MAX_COMMAND_SIZE)
	{
		// Command is too big
		return 1;
	}

	gpioPulse_t irSignal[MAX_PULSES];
	unsigned int pulseCount = 0;

	// Generate Code
	carrierFrequency(outPin, frequency, dutyCycle, leadingPulseDuration, irSignal, &pulseCount);
	gap(leadingGapDuration, irSignal, &pulseCount);

	for (unsigned int i = 0; i < codeLen; i++)
	{
		if (code[i] == '0')
		{
			carrierFrequency(outPin, frequency, dutyCycle, zeroPulse, irSignal, &pulseCount);
			gap(zeroGap, irSignal, &pulseCount);
		}
		else if (code[i] == '1')
		{
			carrierFrequency(outPin, frequency, dutyCycle, onePulse, irSignal, &pulseCount);
			gap(oneGap, irSignal, &pulseCount);
		}
		else
		{
			printf("Warning: Non-binary digit in command\n");
		}
	}

	if (sendTrailingPulse)
	{
		carrierFrequency(outPin, frequency, dutyCycle, onePulse, irSignal, &pulseCount);
	}

	printf("pulse count is %i\n", pulseCount);
	// End Generate Code

	return transmitWave(outPin, irSignal, &pulseCount);
}

int main(int argc, char *argv[]) {
    if (argc < 2) {
		printf("What are you doing, pass a command string in thicko.\n");
        return -1;
    }

    return irSling(
        17,
        38000,
        0.5,
        9000,
        4500,
        562,
        562,
        1688,
        562,
        1,
        argv[1]);
}

#endif
