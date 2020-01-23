package clients

import "testing"

func TestGetDeltas(t *testing.T) {
	res, err := GetStateOfRace()
	if err != nil {
		t.Errorf("TestGetDeltas failed because of TestGetStateOfRace: %v\n", err)
	}

	GetDeltas(res.Overall)
}
