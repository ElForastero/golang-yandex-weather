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

	temperatureRegex := regexp.MustCompile(`(?i)<span class="temp__value">(.*?)</span>`)
	data["temperature"] = temperatureRegex.FindStringSubmatch(stringContent)[1]

	cityRegex := regexp.MustCompile(`(?i)<h1 class="title title_level_1">Погода в <span class="string-with-sticky-item">(.*?)<div class=`)
	data["city"] = "Погода в " + cityRegex.FindStringSubmatch(stringContent)[1]

	//stateRegex := regexp.MustCompile(`(?i)<span class="current-weather__comment">(.*?)</span>`)
	//data["state"] = stateRegex.FindStringSubmatch(stringContent)[1]

	windRegex := regexp.MustCompile(`<span class="wind-speed">(.*?)</span> <span class="fact__unit">(.*?), <abbr`)
	data["wind"] = windRegex.FindStringSubmatch(stringContent)[1]
	data["windUnit"] = windRegex.FindStringSubmatch(stringContent)[2]

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
	output += "Temperature: " + data["temperature"] + " 🌡"
	//output += "State: " + strings.Title(data["state"]) + " "

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

	output += "\nWind: " + data["wind"] + " " + data["windUnit"] + " 💨"

	return
}

// Helper for handling errors
func handleError(err error) {
	if err != nil {
		panic(err)
	}
}
