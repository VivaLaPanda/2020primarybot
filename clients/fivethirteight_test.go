package clients

import (
	"fmt"
	"io/ioutil"
	"testing"
)

func TestGetStateOfRace(t *testing.T) {
	res, err := GetStateOfRace()
	if err != nil {
		t.Errorf("TestGetStateOfRace errored: %v\n", err)
	}

	fmt.Printf("%v", res)

	d1 := []byte(fmt.Sprintf("%v", res))
	ioutil.WriteFile("file.txt", d1, 0644)
}
