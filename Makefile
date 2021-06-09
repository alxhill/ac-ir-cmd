CC=gcc
CFLAGS=-Wall -Wextra -Wpedantic -lm -lpigpio -pthread -lrt

all: ac-ir-cmd irslinger

ac-ir-cmd: main.go
	go build .

irslinger: c/irslinger.c
	$(CC) $(CFLAGS) $< -o $@

install: all
	sudo systemctl stop ac-server
	cp irslinger /usr/bin/irslinger
	cp ac-ir-cmd /usr/bin/ac-ir-cmd
	cp ac-server.service /etc/systemd/system/ac-server.service
	sudo sytemctl daemon-reload
