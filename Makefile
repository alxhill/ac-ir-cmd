CC=gcc
CFLAGS=-Wall -Wextra -Wpedantic -lm -lpigpio -pthread -lrt

all: ac-ir-cmd ir-cmd

ac-ir-cmd: main.go
	go build .

ir-cmd: main.c irslinger.h
	$(CC) $(CFLAGS) $< -o $@
