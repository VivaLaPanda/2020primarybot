package main

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/vivalapanda/edenradiostatus/fetchhtml"
)

var DEBUG bool

func main() {
	DEBUG = false
	stopCode := false
	ticker := time.NewTicker(time.Minute * 30)
	djWas := ""
	go func() {
		for _ = range ticker.C {
			var err error
			djWas, err = task(djWas)

			if err != nil {
				fmt.Printf("%v\n", err)
				stopCode = true
				return
			}
		}
	}()

	fmt.Printf("Eden Radio Status Bot!\n")
	fmt.Printf("Usage:\n  Kill: 'k'\n  Force Poll: 'p'\n  Debug Poll: 'd'\n")
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
			djWas, err = task(djWas)

			if err != nil {
				fmt.Printf("%v\n", err)
				stopCode = true
				return
			}
		case 'd':
			fmt.Println("Forcing poll without tweeting")
			DEBUG = true
			var err error
			djWas, err = task(djWas)
			fmt.Printf("%v, %v\n", djWas, err)
			DEBUG = false
		default:
			fmt.Println("I'm sorry, I didn't understand that...")
			fmt.Printf("Usage:\n  Kill: 'k'\n  Force Poll: 'p'\n  Debug Poll: 'd'\n")
		}

		fmt.Println("What do you want to do?")
	}

	fmt.Printf("Shutting down...\n")
	ticker.Stop()
	fmt.Println("Ticker stopped")
	return
}

func task(djWas string) (string, error) {
	// Get the current DJ
	djIs, err := pollForDj()

	// If we fail try 2 more times
	errCount := 0
	for errCount < 3 && err != nil {
		errCount++
		time.Sleep(time.Second * time.Duration(15*errCount^2))
		djIs, err = pollForDj()
	}

	// Still didn't get a good answer, give up and shut down
	if err != nil {
		err = fmt.Errorf("FATAL ERROR. Couldn't poll Eden Radio: %v", err)
		return "", err
	}

	// Use string to decide on message to tweet
	var tweetText string
	if djIs != djWas {
		if djIs == "Bot-sama" {
			if djWas != "" {
				tweetText = fmt.Sprintf("%v is headed off.", djWas)
			} else {
				tweetText = fmt.Sprintf("Bot waking up. Good Morning!")
			}
		} else {
			tweetText = fmt.Sprintf("LIVE on Eden: %v", djIs)
		}

		djWas = djIs

		// Attempt send
		err = sendTweet(tweetText)

		// If we fail try 2 more times
		errCount = 0
		for errCount < 3 && err != nil {
			errCount++
			time.Sleep(time.Second * time.Duration(15*errCount^2))
			err = sendTweet(tweetText)
		}

		// Still didn't get a good answer, give up and shut down
		if err != nil {
			err = fmt.Errorf("FATAL ERROR. Couldn't send tweet: %v", err)
			return "", err
		}
	}

	// Return the new djWas value
	return djWas, nil
}

func pollForDj() (djName string, err error) {
	djString, pollError := fetchhtml.PollUrlForID("http://edenofthewest.com/", "status-dj")
	if pollError != nil {
		return "", pollError
	}

	return djString, nil
}

func sendTweet(text string) (err error) {
	config := oauth1.NewConfig(consumerKey, consumerSecret)
	token := oauth1.NewToken(accessToken, accessSecret)
	httpClient := config.Client(oauth1.NoContext, token)

	// Twitter client
	client := twitter.NewClient(httpClient)

	// Send a Tweet
	if DEBUG != true {
		_, _, tweetErr := client.Statuses.Update(text, nil)
		if tweetErr != nil {
			return tweetErr
		}
	} else {
		_, _, tweetErr := client.Timelines.HomeTimeline(&twitter.HomeTimelineParams{
			Count: 20,
		})
		if tweetErr != nil {
			return tweetErr
		}
	}

	return
}
