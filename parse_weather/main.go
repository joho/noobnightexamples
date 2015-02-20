package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// http://api.openweathermap.org/data/2.5/weather?q=Melbourne,au

func main() {
	weatherResponse, err := getWeatherResponseBody()
	if err != nil {
		fmt.Print(err)
		return
	}

	melbourne := City{}
	err = json.Unmarshal(weatherResponse, &melbourne)
	if err != nil {
		fmt.Print(err)
		return
	}

	fmt.Printf("It is %v in %v right now\n", melbourne.Weather.NormalisedCurrentTemp(), melbourne.Name)

}

type Weather struct {
	CurrentTemp float64 `json:"temp"`
	MaxTemp     float64 `json:"temp_max"`
}

func (w Weather) NormalisedCurrentTemp() float64 {
	return w.CurrentTemp / 10
}

type City struct {
	Name    string  `json:"name"`
	Weather Weather `json:"main"`
}

func getWeatherResponseBody() ([]byte, error) {
	resp, err := http.Get("http://api.openweathermap.org/data/2.5/weather?q=Melbourne,au")
	if err != nil {
		return []byte{}, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, err
	}
	return body, nil
}
