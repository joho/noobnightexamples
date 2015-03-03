package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// http://api.openweathermap.org/data/2.5/weather?q=Melbourne,au

func main() {
	http.HandleFunc("/", weatherHandler)
	http.ListenAndServe(":5000", nil)
}

func weatherHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	cityName := r.FormValue("city")

	body, err := getWeatherResponseBody(cityName)
	if err != nil {
		panic(err)
	}

	city := City{}
	err = json.Unmarshal(body, &city)
	if err != nil {
		panic(err)
	}

	fmt.Printf("The weather in %v is %v\n", city.Name, city.Weather.NormalisedCurrentTemp())

	response, err := json.Marshal(&city)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

func getWeatherResponseBody(cityName string) ([]byte, error) {
	url := fmt.Sprintf("http://api.openweathermap.org/data/2.5/weather?q=%v,au", cityName)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("Error getting weather: %v", err)
		return []byte(""), err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading weather: %v", err)
		return []byte(""), err
	}
	defer resp.Body.Close()
	return body, nil
}

type City struct {
	Weather Weather `json:"main"`
	Name    string  `json:"name"`
}

type Weather struct {
	CurrentTemp float64 `json:"temp"`
	MaxTemp     float64 `json:"temp_max"`
}

func (w Weather) NormalisedCurrentTemp() float64 {
	return w.CurrentTemp - 273.15
}
