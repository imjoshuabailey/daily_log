package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

func main() {
	var spaceStation Location
	var localWeather Environment
	var err error

	stationLocation := "http://api.open-notify.org/iss-now.json"
	err = Fetch(stationLocation, &spaceStation)
	if err != nil {
		log.Fatalf("Error fetching space station: %s", err.Error())
	}

	weatherByLocation := fmt.Sprintf(
		"http://api.openweathermap.org/data/2.5/weather?lat=%s&lon=%s&APPID=4f1ad4c47a1b62d9b6fc17ca508ad497",
		spaceStation.Position.Latitude,
		spaceStation.Position.Longitude,
	)
	err = Fetch(weatherByLocation, &localWeather)
	if err != nil {
		log.Fatalf("Error fetching space station: %s", err.Error())
	}

	log.Printf("Here is the position of the ISS:\n%+v", spaceStation)
	log.Printf("Here is the weather under the space station:\n%+v", localWeather)
}

// Fetch takes a string that is the URL and a Pointer (&) to a payload struct or map
// and fills it with the data from the URL
func Fetch(url string, pointer interface{}) error {
	res, err := http.Get(url)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	buf := bytes.Buffer{}

	if _, err := io.Copy(&buf, res.Body); err != nil {
		return err
	}

	//log.Printf("RAW RESPONSE: %s", buf.String())

	err = json.NewDecoder(&buf).Decode(pointer)
	if err != nil {
		return err
	}

	return nil
}

// Location holds a position and the time the position was acquired
type Location struct {
	//Timestamp int      `json:"timestamp"`
	Position Position `json:"iss_position"`
}

// Position holds longitude and latitude
type Position struct {
	Latitude  string `json:"latitude"`
	Longitude string `json:"longitude"`
}

// Environment holds values for temp, and humiditiy and a description of the weather
type Environment struct {
	WeatherValues WeatherValues `json:"main"`
	Weather       Weather       `json:"weather"`
}

// WeatherValues holds values for temperature and humidity
type WeatherValues struct {
	Temp     float64 `json:"temp"`
	Humidity int     `json:"humidity"`
}

// Weather is a slice of weather description
type Weather []WeatherDescription

// WeatherDescription holds a description of the weather
type WeatherDescription struct {
	Description string `json:"description"`
}

// International Space Station API link
// http://open-notify.org/Open-Notify-API/ISS-Location-Now/

// documentation for parsing json in golang
// docs:
// https://golang.org/pkg/encoding/json/#pkg-examples
// stackoverflow:
// https://stackoverflow.com/questions/17156371/how-to-get-json-response-in-golang

/*

{
    "coord": {
        "lon": 139,
        "lat": 35
    },
    "sys": {
        "country": "JP",
        "sunrise": 1369769524,
        "sunset": 1369821049
    },
    "weather": [
        {
            "id": 804,
            "main": "clouds",
            "description": "overcast clouds",
            "icon": "04n"
        }
    ],
    "main": {
        "temp": 289.5,
        "humidity": 89,
        "pressure": 1013,
        "temp_min": 287.04,
        "temp_max": 292.04
    },
    "wind": {
        "speed": 7.31,
        "deg": 187.002
    },
    "rain": {
        "3h": 0
    },
    "clouds": {
        "all": 92
    },
    "dt": 1369824698,
    "id": 1851632,
    "name": "Shuzenji",
    "cod": 200
}

*/
