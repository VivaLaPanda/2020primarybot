package main

import (
	"fmt"
	"testing"

	"github.com/vivalapanda/2020primarybot/clients"
)

func TestGetOverallSummary(t *testing.T) {
	primaryState, err := clients.GetStateOfRace()
	if err != nil {
		t.Errorf("TestGetOverallSummary encountered an error: %v\n", err)
	}

	summary := GetOverallSummary(primaryState.Overall)

	fmt.Printf("%s", summary)
}
