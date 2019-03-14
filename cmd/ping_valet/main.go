package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func main() {
	peopleInSpace := "http://api.open-notify.org/astros.json"
	stationLocation := "http://api.open-notify.org/iss-now.json"

	var namedPeople Astronauts
	var spaceStation Location
	var err error

	err = Fetch(peopleInSpace, &namedPeople)
	if err != nil {
		log.Fatalf("Error fetching astronauts: %s", err.Error())
	}

	err = Fetch(stationLocation, &spaceStation)
	if err != nil {
		log.Fatalf("Error fetching space station: %s", err.Error())
	}

	log.Printf("here are the number and names of everyone in space:\n%+v", namedPeople)
	log.Printf("here is the position of the ISS:\n%+v", spaceStation)
}

// Fetch takes a string that is the URL and a Pointer (&) to a payload struct or map
// and fills it with the data from the URL
func Fetch(url string, pointer interface{}) error {
	res, err := http.Get(url)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	err = json.NewDecoder(res.Body).Decode(pointer)
	if err != nil {
		return err
	}
	// the above is the same as this:
	// if err := json.NewDecoder(res.Body).Decode(pointer); err != nil {
	// 	return err
	// }

	return nil
}

type Astronauts struct {
	Message string `json:"message"`
	Number  int    `json:"number"`
	People  People `json:"people"`
}

type People []Person

type Person struct {
	Name  string `json:"name"`
	Craft string `json:"craft"`
}

type Location struct {
	Timestamp int      `json:"timestamp"`
	Position  Position `json:"iss_position"`
}

type Position struct {
	Latitude  string `json:"latitude"`
	Longitude string `json:"longitude"`
}

// International Space Station API link
// http://open-notify.org/Open-Notify-API/ISS-Location-Now/

// documentation for parsing json in golang
// docs:
// https://golang.org/pkg/encoding/json/#pkg-examples
// stackoverflow:
// https://stackoverflow.com/questions/17156371/how-to-get-json-response-in-golang
