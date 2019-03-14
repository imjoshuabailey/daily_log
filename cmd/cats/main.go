package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func main() {
	cats := "https://cat-fact.herokuapp.com/facts/random?animal=cat&amount"

	var catFact CatInfo
	var err error

	err = Fetch(cats, &catFact)
	if err != nil {
		log.Fatalf("Error fetching information about cats: %s", err.Error())
	}

	log.Printf("This is an interesting fact about cats:\n%v", catFact)
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

// CatInfo is a slice of CatStatement
type CatInfo []CatStatement

// CatStatement holds random facts about cats
type CatStatement struct {
	Text string `json:"text"`
}
