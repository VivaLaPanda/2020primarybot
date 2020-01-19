package clients

import (
	"fmt"
	"testing"
)

func TestGetStateOfRace(t *testing.T) {
	res, err := GetStateOfRace()
	if err != nil {
		t.Errorf("TestGetStateOfRace errored: %v\n", err)
	}

	fmt.Printf("%v", res)
}
