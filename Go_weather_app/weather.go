package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type WeatherData struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Daily     struct {
		UvIndexMax       []float64 `json:"uv_index_max"`
		DaylightDuration []float64 `json:"daylight_duration"`
	} `json:"daily"`
	Hourly struct {
		Temperature2m            []float64 `json:"temperature_2m"`
		PrecipitationProbability []float64 `json:"precipitation_probability"`
		SurfacePressure          []float64 `json:"surface_pressure"`
		WindSpeed10m             []float64 `json:"wind_speed_10m"`
	} `json:"hourly"`
	Current struct {
		Temperature2m      float64 `json:"temperature_2m"`
		Precipitation      float64 `json:"precipitation"`
		SurfacePressure    float64 `json:"surface_pressure"`
		RelativeHumidity2m float64 `json:"relative_humidity_2m"`
	} `json:"current"`
}

func main() {
	fmt.Println("Weather App")

	// Get the weather data
	weatherData, err := getWeatherData()
	if err != nil {
		fmt.Println("Error getting weather data:", err)
		return
	}

	displayWeatherData(weatherData)
}

func getWeatherData() (WeatherData, error) {
	cacheFile := "weather_cache.json"

	// caching
	if info, err := os.Stat(cacheFile); err == nil {
		if time.Since(info.ModTime()) < 1*time.Hour {
			fmt.Println("Using cached data...")
			data, err := os.ReadFile(cacheFile)
			if err == nil {
				var weatherData WeatherData
				if err := json.Unmarshal(data, &weatherData); err == nil {
					return weatherData, nil
				}
			}
		}
	}

	fmt.Println("Fetching fresh data...")
	// Get the weather data from the API
	response, err := http.Get("https://api.open-meteo.com/v1/forecast?latitude=52.52&longitude=13.41&daily=uv_index_max,daylight_duration&hourly=temperature_2m,precipitation_probability,surface_pressure,wind_speed_10m&current=temperature_2m,precipitation,surface_pressure,relative_humidity_2m&timezone=Europe%2FBerlin")
	if err != nil {
		return WeatherData{}, err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return WeatherData{}, err
	}

	// caching
	_ = os.WriteFile(cacheFile, body, 0644)

	var weatherData WeatherData
	err = json.Unmarshal(body, &weatherData)
	if err != nil {
		return WeatherData{}, err
	}

	return weatherData, nil
}

func displayWeatherData(weatherData WeatherData) {
	fmt.Println("Weather data:")
	fmt.Println("Latitude:", weatherData.Latitude)
	fmt.Println("Longitude:", weatherData.Longitude)
	fmt.Println("Daily:")
	fmt.Println("UvIndexMax:", weatherData.Daily.UvIndexMax)
	fmt.Println("DaylightDuration:", weatherData.Daily.DaylightDuration)
	fmt.Println("Hourly:")
	fmt.Println("Temperature2m:", weatherData.Hourly.Temperature2m)
	fmt.Println("PrecipitationProbability:", weatherData.Hourly.PrecipitationProbability)
	fmt.Println("SurfacePressure:", weatherData.Hourly.SurfacePressure)
	fmt.Println("WindSpeed10m:", weatherData.Hourly.WindSpeed10m)
	fmt.Println("Current:")
	fmt.Println("Temperature2m:", weatherData.Current.Temperature2m)
	fmt.Println("Precipitation:", weatherData.Current.Precipitation)
	fmt.Println("SurfacePressure:", weatherData.Current.SurfacePressure)
	fmt.Println("RelativeHumidity2m:", weatherData.Current.RelativeHumidity2m)
}
