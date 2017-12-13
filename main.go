package main

import (
	"encoding/json"
	"fmt"
	"github.com/stianeikeland/go-rpio"
	"log"
	"net/http"
	"regexp"
)

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
		allPins(pins[:], rpio.High)
	case rOff.MatchString(s):
		log.Println("off command issued")
		allPins(pins[:], rpio.Low)
	default:
		httpLogger(w, nil, 400, "Bad Request")
	}
}

func lightsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		b := pinStatus(pins[:])
		log.Println(b)
		fmt.Fprintf(w, "%s", jString(b))
	default:
		http.Error(w, "Method Not Allowed", 405)
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

func jString(v interface{}) string {
	j, _ := json.Marshal(v)
	return string(j)
}

func main() {
	rpioInit()
	http.HandleFunc("/program", programHandler)
	http.HandleFunc("/lights", lightsHandler)
	http.ListenAndServe(":8080", nil)
}
