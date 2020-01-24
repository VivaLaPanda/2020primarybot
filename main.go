package main

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
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

	// Get text to tweet
	mainText := OverallSummary(primaryState.Overall)
	statesText := BiggestMoverStates(primaryState.States)

	// Prepare our twitter access
	config := oauth1.NewConfig(consumerKey, consumerSecret)
	token := oauth1.NewToken(accessToken, accessSecret)
	httpClient := config.Client(oauth1.NoContext, token)

	// Twitter client
	client := twitter.NewClient(httpClient)

	// Send a Tweet with he overall data
	tweet, _, tweetErr := client.Statuses.Update(mainText, nil)
	if tweetErr != nil {
		return tweetErr
	}

	// Reply to that tweet with the state data
	tweet, _, tweetErr = client.Statuses.Update(statesText, &twitter.StatusUpdateParams{InReplyToStatusID: tweet.ID})
	if tweetErr != nil {
		return tweetErr
	}

	return nil
}
