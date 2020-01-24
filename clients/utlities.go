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
type OverallStats struct {
	Candidate string  `json:"candidate"`
	Date      string  `json:"date"`
	Majority  float64 `json:"majority"`
	Plurality float64 `json:"plurality"`
	Running   bool    `json:"running"`
}

func (s OverallStats) Name() string         { return s.Candidate }
func (s OverallStats) GetDate() string      { return s.Date }
func (s OverallStats) WinMajority() float64 { return s.Plurality }

type StateStats struct {
	Candidate         string  `json:"candidate"`
	Date              string  `json:"date"`
	DelegatePlurality float64 `json:"delegate_plurality"`
	VotePlurality     float64 `json:"vote_plurality"`
}

func (s StateStats) Name() string         { return s.Candidate }
func (s StateStats) GetDate() string      { return s.Date }
func (s StateStats) WinMajority() float64 { return s.DelegatePlurality }

type CandidateStats interface {
	Name() string
	GetDate() string
	WinMajority() float64
}

type RaceStats []CandidateStats

func (rs RaceStats) GroupByCandidate() (candidates map[string]RaceStats) {
	candidates = make(map[string]RaceStats)
	notRunning := make([]string, 0)

	for _, stats := range rs {
		// Only add candidate data if they're running
		if stats.WinMajority() > 0 {
			// Check if the candidate is already in the map, if not, add them
			group, ok := candidates[stats.Name()]
			if !ok {
				group = make([]CandidateStats, 0)
			}

			candidates[stats.Name()] = append(group, stats)
		} else {
			notRunning = append(notRunning, stats.Name())
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
		if stats.GetDate() == timestring {
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

		deltas[name] = today.WinMajority() - yesterday.WinMajority()
	}

	return
}
