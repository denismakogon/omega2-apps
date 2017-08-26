package api

import (
	"errors"
	"fmt"
	"github.com/ChimeraCoder/anaconda"
	"net/http"
	"os"
	"strconv"
	"time"
)

type TwitterSecret struct {
	Mapper
	ConsumerKey    string `json:"consumer_key"`
	ConsumerSecret string `json:"consumer_secret"`
	APIToken       string `json:"api_key"`
	APITokenSecret string `json:"api_key_secret"`
}

func (twitter *TwitterSecret) FromEnv() (*anaconda.TwitterApi, error) {
	consumerKey := os.Getenv("TwitterConsumerKey")
	if consumerKey == "" {
		return nil, errors.New("Empty TwitterConsumerKey")
	}
	consumerSecret := os.Getenv("TwitterConsumerSecret")
	if consumerSecret == "" {
		return nil, errors.New("TwitterConsumerSecret")
	}
	apiToken := os.Getenv("TwitterAccessToken")
	if apiToken == "" {
		return nil, errors.New("TwitterAccessToken")
	}
	apiTokenSecret := os.Getenv("TwitterAccessTokenSecret")
	if apiTokenSecret == "" {
		return nil, errors.New("TwitterAccessTokenSecret")
	}

	anaconda.SetConsumerKey(consumerKey)
	anaconda.SetConsumerSecret(consumerSecret)
	api := anaconda.NewTwitterApi(apiToken, apiTokenSecret)
	ok, err := api.VerifyCredentials()
	if !ok {
		return nil, errors.New("Unauthorized to Twitter")
	}
	if err != nil {
		return nil, err
	}
	twitter.ConsumerKey = consumerKey
	twitter.ConsumerSecret = consumerSecret
	twitter.APIToken = apiToken
	twitter.APITokenSecret = apiTokenSecret

	api.SetDelay(5 * time.Second)
	return api, nil
}

func (omega *OnionOmega2) GetRecentMentions() (tweets []anaconda.Tweet, err error) {
	tweets, err = omega.TwitterAPI.GetMentionsTimeline(omega.SearchValues)
	if err != nil {
		return nil, err
	}
	if len(tweets) != 0 {
		omega.SearchValues.Set(
			"since_id", strconv.FormatInt(tweets[len(tweets)-1].Id, 10))
	}
	return tweets, nil
}

func (omega *OnionOmega2) ProcessTweets(tweets *[]anaconda.Tweet, httpClient *http.Client, fnAPIURL, fnToken string) error {
	detect, err := http.NewRequest(
		http.MethodPost, fmt.Sprintf("%s/r/whereisit/detect-where", fnAPIURL),
		nil)
	if err != nil {
		return err
	}

	fail, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/r/whereisit/tweet-fail", fnAPIURL), nil)
	if err != nil {
		return err
	}

	for _, tweet := range *tweets {
		omega.PrintTweetInfo(tweet)
		if len(tweet.Entities.Media) != 0 {
			media := tweet.Entities.Media[0]
			user := fmt.Sprintf("@%v", tweet.User.ScreenName)
			if media.Type != "photo" {
				payload := &RequestPayload{
					User:    user,
					TweetID: tweet.IdStr,
				}
				err := doRequest(payload, fail, httpClient, fnToken)
				if err != nil {
					return err
				}
			} else {
				payload := &RequestPayload{
					MediaURL: media.Expanded_url,
					User:     user,
					TweetID:  tweet.IdStr,
					GCloud:   *omega.GCloudAuth,
				}
				err := doRequest(payload, detect, httpClient, fnToken)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (omega *OnionOmega2) PrintTweetInfo(tweet anaconda.Tweet) {
	fmt.Printf(fmt.Sprintf(
		"Found new tweet: %v from @%v.", tweet.Text, tweet.User.ScreenName))
}
