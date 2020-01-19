package clients

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"
)

type PrimaryState struct {
	overall []CandidateStats
	states  map[string][]CandidateStats
}

// String is the candidate, int is their chances as percentage
type CandidateStats struct {
	Candidate string  `json:"candidate"`
	Date      string  `json:"date"`
	Majority  float64 `json:"majority"`
	Plurality float64 `json:"plurality"`
	Running   bool    `json:"running"`
}

func GetStateOfRace() (raceState PrimaryState, err error) {
	raceState = PrimaryState{}

	raceState.overall, err = getOverallRaceState()

	return
}

func getOverallRaceState() (overall []CandidateStats, err error) {
	resp, getErr := http.Get("https://projects.fivethirtyeight.com/2020-primary-forecast/js/data.js")
	if getErr != nil {
		err = fmt.Errorf("Failed to fetch 538 url; ERR: %v", getErr)
		return
	}

	// Get full data as string
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	fullJs := buf.String()
	resp.Body.Close()

	// Regex to get the data
	re, err := regexp.Compile("parse\\('{\"\":.*?]")
	matches := re.FindAllString(fullJs, -1)
	data := matches[2]
	data = strings.Trim(data, "parse('{\"\":")

	err = json.Unmarshal([]byte(data), &overall)

	return
}
