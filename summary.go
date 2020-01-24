package main

import (
	"fmt"
	"math"
	"sort"
	"strings"

	"github.com/vivalapanda/2020primarybot/clients"
)

func OverallSummary(overallStats clients.RaceStats) (summaryString string) {
	deltas := clients.GetDeltas(overallStats)
	stringDeltas := make([]string, 0)

	for name, flt := range deltas {
		intDelta := int64(math.RoundToEven(100 * flt))
		if intDelta < 0 {
			stringDeltas = append(stringDeltas, fmt.Sprintf("%s: %d", name, intDelta))
		} else if intDelta > 0 {
			stringDeltas = append(stringDeltas, fmt.Sprintf("%s: +%d", name, intDelta))
		}
	}

	return strings.Join(stringDeltas, "\n")
}

type deltaTracker struct {
	delta float64
	state string
	name  string
}

type byDelta []deltaTracker

func BiggestMoverStates(stateRaces map[string]clients.RaceStats) (summaryString string) {
	// Calculate the biggest movers
	biggestMovers := calcBiggestMovers(stateRaces)

	// Format the results
	stringDeltas := make([]string, len(biggestMovers))
	for idx, dt := range biggestMovers {
		intDelta := int64(math.RoundToEven(100 * dt.delta))
		if intDelta < 0 {
			stringDeltas[idx] = fmt.Sprintf("%s in %s: %d", dt.name, dt.state, intDelta)
		} else if intDelta > 0 {
			stringDeltas[idx] = fmt.Sprintf("%s in %s: +%d", dt.name, dt.state, intDelta)
		}
	}

	return strings.Join(stringDeltas, "\n")
}

func calcBiggestMovers(stateRaces map[string]clients.RaceStats) (biggestMovers []deltaTracker) {
	// Store the top 5 biggest movers here
	biggestMovers = make([]deltaTracker, 0, len(stateRaces))
	// Iterate over the states, getting deltas and putting them above
	for state, stats := range stateRaces {
		deltas := clients.GetDeltas(stats)

		for name, delta := range deltas {
			biggestMovers = append(biggestMovers, deltaTracker{delta, state, name})
		}
	}

	// Sort the biggest movers and take the top 5
	sort.Slice(biggestMovers, func(i, j int) bool {
		return biggestMovers[i].delta < biggestMovers[j].delta
	})

	return biggestMovers[:5]
}
