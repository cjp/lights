package main

import (
	"encoding/json"
	"fmt"
	"github.com/stianeikeland/go-rpio"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"regexp"
)

var pins [4]rpio.Pin

func httpLogger(w http.ResponseWriter, e error, status int, message string) {
	log.Println(message, e)
	http.Error(w, message, status)
}

func commandMatcher(w http.ResponseWriter, s string) {
	var rOn = regexp.MustCompile(`^on$`)
	var rOff = regexp.MustCompile(`^off$`)
	switch {
	case rOn.MatchString(s):
		log.Println("on command issued")
	case rOff.MatchString(s):
		log.Println("off command issued")
	default:
		httpLogger(w, nil, 400, "Bad Request")
	}
}

func programHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "PUT":
		decoder := json.NewDecoder(r.Body)
		var t string
		err := decoder.Decode(&t)
		if err != nil {
			httpLogger(w, err, 400, "Bad Request")
		}
		defer r.Body.Close()
		commandMatcher(w, t)
	default:
		http.Error(w, "Method Not Allowed", 405)
	}
}

func onHandler(w http.ResponseWriter, r *http.Request) {
	for _, p := range pins {
		p.High()
	}
	io.WriteString(w, "OK\n")

}

func offHandler(w http.ResponseWriter, r *http.Request) {
	for _, p := range pins {
		p.Low()
	}
	io.WriteString(w, "OK\n")
}

func statusHandler(w http.ResponseWriter, r *http.Request) {
	for i, p := range pins {
		if p.Read() == rpio.Low {
			fmt.Fprintf(w, "%d off\n", i)
		} else {
			fmt.Fprintf(w, "%d on\n", i)
		}
	}
}

func main() {

	rand.Seed(42)

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

	http.HandleFunc("/on/", onHandler)
	http.HandleFunc("/off/", offHandler)
	http.HandleFunc("/status/", statusHandler)
	http.HandleFunc("/program", programHandler)
	http.ListenAndServe(":8080", nil)

}
