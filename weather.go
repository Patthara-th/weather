package weather

import (
	"encoding/json"	
	"net/http"	
	"strconv"	
	"sync"
)


var (
	wg      sync.WaitGroup
)

type Weatherresult struct {
	Coord struct {
		Lon float64 `json:"lon"`
		Lat float64 `json:"lat"`
	} `json:"coord"`

	Weather []struct {
		ID          int    `json:"id"`
		Main        string `json:"main"`
		Description string `json:"description"`
		Icon        string `json:"icon"`
	} `json:"weather"`

	Base string `json:"base"`

	Main struct {
		Temp     float64 `json:"temp"`
		Pressure int `json:"pressure"`
		Humidity int `json:"humidity"`
		TempMin  int `json:"temp_min"`
		TempMax  int `json:"temp_max"`
	} `json:"main"`

	Visibility int `json:"visibility"`

	Wind       struct {
		Speed float64 `json:"speed"`
		Deg   int     `json:"deg"`
	} `json:"wind"`

	Clouds struct {
		All int `json:"all"`
	} `json:"clouds"`

	Dt  int `json:"dt"`

	Sys struct {
		Type    int     `json:"type"`
		ID      int     `json:"id"`
		Message float64 `json:"message"`
		Country string  `json:"country"`
		Sunrise int     `json:"sunrise"`
		Sunset  int     `json:"sunset"`
	} `json:"sys"`
	
	ID   int    `json:"id"`
	Name string `json:"name"`
	Cod  int    `json:"cod"`
}

func Getcity(city string) string {

	URL := "http://localhost:8882/api/v1/weather/"
	res, _ := http.Get(URL + city)
	
	weather := new(Weatherresult)
	json.NewDecoder(res.Body).Decode(weather)	
	
	return weather.Name  + "\n" + strconv.FormatFloat(weather.Main.Temp, 'f', 0, 64) +  "c " + weather.Weather[0].Description + "\n\n"
}

func Getallcity() string {
	
	city := [5]string{"hobart","newyork","kupang","nairobi","bangkok"}

	var result string

	resultCh := make(chan string)
	cityCh := make(chan string)
	
	
	go func() {
		for _,v := range city {
			cityCh <- v
		}
		close(cityCh)
	}()	
	

	wg.Add(5)
	for i := 1; i <= 5; i++ {
		go getConCity(cityCh, resultCh)
	}	

	go func() {
		wg.Wait()
		close(resultCh)
	}()

	for x := range resultCh {
		result += x
	}
	return result
}


func getConCity(queue <-chan string, resultCh chan<- string) {
	for x := range queue {
		resultCh <- Getcity(x)
	}
	wg.Done()
}
