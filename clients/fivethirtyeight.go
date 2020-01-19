package clients

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"
)

var states = []string{"alabama", "alaska", "arizona", "arkansas", "california", "colorado", "connecticut", "delaware", "florida", "georgia", "guam", "hawaii", "idaho", "illinois", "indiana", "iowa", "kansas", "kentucky", "louisiana", "maine", "district-of-columbia", "democrats-abroad", "maryland", "massachusetts", "michigan", "minnesota", "mississippi", "missouri", "montana", "nebraska", "nevada", "new-hampshire", "new-jersey", "new-mexico", "new-york", "north-carolina", "north-dakota", "northern-marianas", "ohio", "oklahoma", "oregon", "puerto-rico", "pennsylvania", "rhode-island", "south-carolina", "south-dakota", "tennessee", "texas", "utah", "vermont", "virgin-islands", "virginia", "washington", "west-virginia", "wisconsin", "wyoming"}

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

	raceState.overall, err = getOverallRace()
	raceState.states = make(map[string][]CandidateStats)
	for _, elem := range states {
		raceState.states[elem], err = getStateStats(elem)
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

	err = json.Unmarshal([]byte(data), &overall)

	return
}

func getStateStats(state string) (stateStats []CandidateStats, err error) {
	resp, getErr := http.Get("https://projects.fivethirtyeight.com/2020-primary-forecast/" + state + ".json")
	if getErr != nil {
		err = fmt.Errorf("Failed to fetch 538 url; ERR: %v", getErr)
		return
	}

	respReceiver := struct {
		State_chances []CandidateStats
	}{stateStats}

	err = json.NewDecoder(resp.Body).Decode(&respReceiver)
	stateStats = respReceiver.State_chances

	return
}
