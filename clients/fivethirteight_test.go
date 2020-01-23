package clients

import (
	"testing"
	"time"
)

func TestGetStateOfRace(t *testing.T) {
	res, err := GetStateOfRace()
	if err != nil {
		t.Errorf("TestGetStateOfRace errored: %v\n", err)
		return
	}

	if len(res.Overall) < 1 {
		t.Errorf("TestGetStateOfRace failed because it returned no results: %v\n", res)
		return
	}

	currentTime := time.Now()
	todayString := currentTime.Format("2006-01-02")
	if res.Overall[0].Date != todayString {
		t.Errorf("TestGetStateOfRace failed because most recent result wasn't today: %v\n", res)
	}
}
