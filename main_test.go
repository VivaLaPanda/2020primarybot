package main

import (
	"fmt"
	"testing"

	"github.com/vivalapanda/2020primarybot/clients"
)

func TestOverallSummary(t *testing.T) {
	primaryState, err := clients.GetStateOfRace()
	if err != nil {
		t.Errorf("TestOverallSummary encountered an error: %v\n", err)
	}

	summary := OverallSummary(primaryState.Overall)

	if len(summary) == 0 {
		t.Errorf("TestOverallSummary failed because it produced a 0 len string.\n")
		return
	}

	fmt.Printf("%s\n", summary)
}

func TestBiggestMoverStates(t *testing.T) {
	primaryState, err := clients.GetStateOfRace()
	if err != nil {
		t.Errorf("TestGetOverallSummary encountered an error: %v\n", err)
	}

	biggestMovers := BiggestMoverStates(primaryState.States)

	if len(biggestMovers) < 0 {
		t.Errorf("TestCalcBiggestMovers failed because it produced a 0 len string")
	}

	fmt.Printf("%s\n", biggestMovers)
}
