package main

import (
	"fmt"
	"github.com/stianeikeland/go-rpio"
	"os"
)

var pins [4]rpio.Pin

func allPins(p []rpio.Pin, s rpio.State) {
	for _, i := range p {
		rpio.WritePin(i, s)
	}
}

func pinStatus(p []rpio.Pin) (b []bool) {
	for _, i := range p {
		b = append(b, i.Read() == rpio.High)
	}
	return
}

func rpioInit() {
	if err := rpio.Open(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer rpio.Close()

	pins[0] = rpio.Pin(17)
	pins[1] = rpio.Pin(18)
	pins[2] = rpio.Pin(22)
	pins[3] = rpio.Pin(23)
	for _, p := range pins {
		p.Output()
	}
}
