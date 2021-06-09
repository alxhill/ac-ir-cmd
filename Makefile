CC=gcc
CFLAGS=-Wall -Wextra -Wpedantic -lm -lpigpio -pthread -lrt

all: ac-ir-cmd irslinger

ac-ir-cmd: main.go
	go build .

irslinger: c/irslinger.c
	$(CC) $(CFLAGS) $< -o $@

install: all
	cp irslinger /usr/bin/irslinger
	cp ac-ir-cmd /usr/bin/ac-ir-cmd
