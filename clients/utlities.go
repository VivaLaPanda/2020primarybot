package clients

import (
	"fmt"
	"time"
)

var states = []string{"alabama", "alaska", "arizona", "arkansas", "california", "colorado", "connecticut", "delaware", "florida", "georgia", "guam", "hawaii", "idaho", "illinois", "indiana", "iowa", "kansas", "kentucky", "louisiana", "maine", "district-of-columbia", "democrats-abroad", "maryland", "massachusetts", "michigan", "minnesota", "mississippi", "missouri", "montana", "nebraska", "nevada", "new-hampshire", "new-jersey", "new-mexico", "new-york", "north-carolina", "north-dakota", "northern-marianas", "ohio", "oklahoma", "oregon", "puerto-rico", "pennsylvania", "rhode-island", "south-carolina", "south-dakota", "tennessee", "texas", "utah", "vermont", "virgin-islands", "virginia", "washington", "west-virginia", "wisconsin", "wyoming"}

type PrimaryState struct {
	Overall RaceStats
	States  map[string]RaceStats
}

// String is the candidate, int is their chances as percentage
type CandidateStats struct {
	Candidate string  `json:"candidate"`
	Date      string  `json:"date"`
	Majority  float64 `json:"majority"`
	Plurality float64 `json:"plurality"`
	Running   bool    `json:"running"`
}

type RaceStats []CandidateStats

func (a RaceStats) Len() int           { return len(a) }
func (a RaceStats) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a RaceStats) Less(i, j int) bool { return a[i].Date < a[j].Date }

func (rs RaceStats) GroupByCandidate() (candidates map[string]RaceStats) {
	candidates = make(map[string]RaceStats)
	notRunning := make([]string, 0)

	for _, stats := range rs {
		// Only add candidate data if they're running
		if stats.Running {
			// Check if the candidate is already in the map, if not, add them
			group, ok := candidates[stats.Candidate]
			if !ok {
				group = make([]CandidateStats, 0)
			}

			candidates[stats.Candidate] = append(group, stats)
		} else {
			notRunning = append(notRunning, stats.Candidate)
		}
	}

	// Clear out all candidates not currently running
	for _, name := range notRunning {
		delete(candidates, name)
	}

	return
}

// Will this fail at midnight?
func (rs RaceStats) getTodayStats() (res CandidateStats, err error) {
	currentTime := time.Now()
	todayString := currentTime.Format("2006-01-02")

	return getAnyStats(rs, todayString)
}

func (rs RaceStats) getYesterdayStats() (res CandidateStats, err error) {
	currentTime := time.Now().AddDate(0, 0, -1)
	yesterdayString := currentTime.Format("2006-01-02")

	return getAnyStats(rs, yesterdayString)
}

func getAnyStats(rs RaceStats, timestring string) (res CandidateStats, err error) {
	for _, stats := range rs {
		if stats.Date == timestring {
			return stats, nil
		}
	}

	return res, fmt.Errorf("Failed to find today in stat list")
}

func GetDeltas(raceState RaceStats) (deltas map[string]float64) {
	deltas = make(map[string]float64)

	// Filter out not running

	// Sort the state by candidate\
	candidates := raceState.GroupByCandidate()

	// Get deltas
	for name, statList := range candidates {
		today, err := statList.getTodayStats()
		yesterday, err := statList.getYesterdayStats()
		if err != nil {
			panic("Error calculating deltas")
		}

		deltas[name] = today.Majority - yesterday.Majority
	}

	return
}
