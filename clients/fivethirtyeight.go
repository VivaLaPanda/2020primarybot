package clients

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"
)

func GetStateOfRace() (raceState PrimaryState, err error) {
	raceState = PrimaryState{}

	raceState.Overall, err = getOverallRace()
	raceState.States = make(map[string]RaceStats)
	for _, elem := range states {
		raceState.States[elem], err = getStateStats(elem)
		if err != nil {
			return
		}
	}

	return
}

func getOverallRace() (overall []CandidateStats, err error) {
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

	var rawStats []OverallStats
	err = json.Unmarshal([]byte(data), &rawStats)

	overall = make([]CandidateStats, len(rawStats))
	for idx, elem := range rawStats {
		overall[idx] = CandidateStats(elem)
	}

	return
}

func getStateStats(state string) (stateStats []CandidateStats, err error) {
	resp, getErr := http.Get("https://projects.fivethirtyeight.com/2020-primary-forecast/" + state + ".json")
	if getErr != nil {
		err = fmt.Errorf("Failed to fetch 538 url; ERR: %v", getErr)
		return
	}

	respReceiver := struct {
		State_chances []StateStats
	}{}

	err = json.NewDecoder(resp.Body).Decode(&respReceiver)
	rawStats := respReceiver.State_chances

	stateStats = make([]CandidateStats, len(rawStats))
	for idx, elem := range rawStats {
		stateStats[idx] = CandidateStats(elem)
	}

	return stateStats, err
}
