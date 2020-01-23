package clients

import (
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

func SendTweet(text string) (err error) {
	config := oauth1.NewConfig(consumerKey, consumerSecret)
	token := oauth1.NewToken(accessToken, accessSecret)
	httpClient := config.Client(oauth1.NoContext, token)

	// Twitter client
	client := twitter.NewClient(httpClient)

	// Send a Tweet
	_, _, tweetErr := client.Statuses.Update(text, nil)
	if tweetErr != nil {
		return tweetErr
	}

	return
}
