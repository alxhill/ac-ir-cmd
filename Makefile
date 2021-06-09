CC=gcc
CFLAGS=-Wall -Wextra -Wpedantic -lm -lpigpio -pthread -lrt

all: ac-ir-cmd irslinger

ac-ir-cmd: main.go
	go build .

irslinger: c/irslinger.c
	$(CC) $(CFLAGS) $< -o $@
