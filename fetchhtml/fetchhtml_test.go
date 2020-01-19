package fetchhtml

import "testing"

// Kinda crappy test because it relies on the million dollar homepage being alive
// /shrug
func TestPollUrlForID(t *testing.T) {
	testString, err := PollUrlForID("http://www.milliondollarhomepage.com/", "note")
	if err != nil {
		t.Errorf("Test failed: %v\n", err)
	}
	if testString[:11] != "The Million" {
		t.Errorf("Test failed: Tag contained unexpected data: %v", testString[:11])
	}
}
