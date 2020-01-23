package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strings"
	"time"

	"github.com/vivalapanda/2020primarybot/clients"
)

var DEBUG bool

func main() {
	DEBUG = false
	stopCode := false
	ticker := time.NewTicker(time.Hour * 24)
	go func() {
		for _ = range ticker.C {
			var err error
			err = task()

			if err != nil {
				fmt.Printf("%v\n", err)
				stopCode = true
				return
			}
		}
	}()

	fmt.Printf("Eden Radio Status Bot!\n")
	fmt.Printf("Usage:\n  Kill: 'k'\n  Force Poll: 'p'\n")
	for stopCode == false {
		fmt.Printf("Enter your desired operation: ")
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		inText := scanner.Text()
		opChar := inText[0]

		switch opChar {
		case 'k':
			stopCode = true
		case 'p':
			fmt.Println("Forcing poll")
			var err error
			err = task()

			if err != nil {
				fmt.Printf("%v\n", err)
				stopCode = true
				return
			}
		default:
			fmt.Println("I'm sorry, I didn't understand that...")
			fmt.Printf("Usage:\n  Kill: 'k'\n  Force Poll: 'p'\n")
		}

		fmt.Println("What do you want to do?")
	}

	fmt.Printf("Shutting down...\n")
	ticker.Stop()
	fmt.Println("Ticker stopped")
	return
}
func task() (err error) {
	primaryState, err := clients.GetStateOfRace()
	if err != nil {
		return
	}

	GetOverallSummary(primaryState.Overall)

	return nil
}

func GetOverallSummary(overallStats clients.RaceStats) (summaryString string) {
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
