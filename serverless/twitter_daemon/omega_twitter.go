package main

import (
	"fmt"
	"github.com/ChimeraCoder/anaconda"
	"time"
	"net/url"
	"strconv"
	"os"
	"github.com/Sirupsen/logrus"
)

type OnionOmega2 struct {
	TwitterAPI *anaconda.TwitterApi
	SearchValues url.Values
}


func (omega *OnionOmega2) GetRecentMentions() (tweets []anaconda.Tweet, err error) {
	tweets, err = omega.TwitterAPI.GetMentionsTimeline(omega.SearchValues)
	if err != nil {
		return nil, err
	}
	if len(tweets) != 0 {
		omega.SearchValues.Set(
		"since_id", strconv.FormatInt(tweets[len(tweets) - 1].Id, 10))
	}
	return tweets, nil
}

func (omega *OnionOmega2) PrintTweetInfo(tweets *[]anaconda.Tweet) {
	if len(*tweets) != 0 {
		for _, tweet := range *tweets {
			fmt.Printf(fmt.Sprintf(
				"Found new tweet: %v from @%v. Has coordinates? - %v\n",
				tweet.Text, tweet.User.ScreenName, tweet.HasCoordinates()))
		}
	}
}


func main() {
	consumerKey := os.Getenv("TwitterConsumerKey")
	consumerSecret := os.Getenv("TwitterConsumerSecret")
	apiToken := os.Getenv("TwitterAccessToken")
	apiTokenSecret := os.Getenv("TwitterAccessTokenSecret")
	anaconda.SetConsumerKey(consumerKey)
	anaconda.SetConsumerSecret(consumerSecret)
	api := anaconda.NewTwitterApi(apiToken, apiTokenSecret)

	v := url.Values{}
	v.Set("count", "200")
	omega := OnionOmega2{TwitterAPI:api, SearchValues:v}
	for {
		tweets, err := omega.GetRecentMentions()
		if err != nil {
			panic(err)
		}
		omega.PrintTweetInfo(&tweets)
		time.Sleep(time.Second * 12)
	}
}
