package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func main() {
	dogImages := "https://dog.ceo/api/breeds/image/random"

	var dog Images
	var err error

	err = Fetch(dogImages, &dog)
	if err != nil {
		log.Fatalf("Error fetching dogs: %s", err.Error())
	}

	log.Printf("Check out this dog:\n%v", dog)
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

	return nil
}

// Images holds the URL for an imaage of a dog
type Images struct {
	// Status  string `json:"status"`
	Message string `json:"message"`
}
