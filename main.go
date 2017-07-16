package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"unicode/utf8"
)

func main() {
	response, err := http.Get("https://yandex.ru/pogoda/")
	handleError(err)
	defer response.Body.Close()

	content, err := ioutil.ReadAll(response.Body)
	handleError(err)

	info := crawl(content)
	fmt.Println(formatOutputData(info))
}

// Parses html responses for weather data
func crawl(content []byte) *map[string]string {
	data := make(map[string]string)
	stringContent := string(content)

	temperatureRegex := regexp.MustCompile(`(?i)<div class="current-weather__thermometer current-weather__thermometer_type_now">(.*?)</div>`)
	data["temperature"] = temperatureRegex.FindStringSubmatch(stringContent)[1]

	cityRegex := regexp.MustCompile(`(?i)<h1 class="title title_level_1">(.*?)((&nbsp;)*<.*?>)*</h1>`)
	data["city"] = cityRegex.FindStringSubmatch(stringContent)[1]

	stateRegex := regexp.MustCompile(`(?i)<span class="current-weather__comment">(.*?)</span>`)
	data["state"] = stateRegex.FindStringSubmatch(stringContent)[1]

	windRegex := regexp.MustCompile(`<span class="wind-speed">(.*?)</span>`)
	data["wind"] = windRegex.FindStringSubmatch(stringContent)[1]

	return &data
}

// Prepares data for displaying into console
func formatOutputData(info *map[string]string) (output string) {
	data := *info
	cityLength := utf8.RuneCountInString(data["city"]) + 4

	output = strings.Repeat("#", cityLength)
	output += "\n# "
	output += data["city"] + " #\n"
	output += strings.Repeat("#", cityLength)
	output += "\n\n"
	output += "Temperature: " + data["temperature"] + " 🌡\n"
	output += "State: " + strings.Title(data["state"]) + " "

	states := map[string]string{
		"ясно": "☀️",
		"облачно": "⛅️",
		"дожд": "🌧",
	}

	for state, icon := range states {
		if strings.Contains(data["state"], state) {
			output += icon
			break
		}
	}

	output += "\nWind: " + data["wind"] + " 💨"

	return
}

// Helper for handling errors
func handleError(err error) {
	if err != nil {
		panic(err)
	}
}
