package api

import (
	"errors"
	"fmt"
	"github.com/ChimeraCoder/anaconda"
	"net/http"
	"strconv"
)

type TwitterSecret struct {
	ConsumerKey    string `json:"consumer_key"`
	ConsumerSecret string `json:"consumer_secret"`
	APIToken       string `json:"api_key"`
	APITokenSecret string `json:"api_key_secret"`
}

func (twitter *TwitterSecret) FromEnv() (*anaconda.TwitterApi, error) {
	err := StructFromEnv(twitter)
	if err != nil {
		return nil, err
	}
	anaconda.SetConsumerKey(twitter.ConsumerKey)
	anaconda.SetConsumerSecret(twitter.ConsumerSecret)
	api := anaconda.NewTwitterApi(twitter.APIToken, twitter.APITokenSecret)
	ok, err := api.VerifyCredentials()
	if !ok {
		return nil, errors.New("Unauthorized to Twitter")
	}
	if err != nil {
		return nil, err
	}
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

func (omega *OnionOmega2) ProcessTweets(tweet anaconda.Tweet, httpClient *http.Client, fnAPIURL, fnToken string) error {
	detect, err := http.NewRequest(
		http.MethodPost, fmt.Sprintf("%s/r/where-is-it/detect-where", fnAPIURL),
		nil)
	if err != nil {
		return err
	}

	fail, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/r/where-is-it/tweet-fail", fnAPIURL), nil)
	if err != nil {
		return err
	}

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
				MediaURL:     media.Media_url,
				User:         user,
				TweetID:      tweet.IdStr,
				TweetFail:    fmt.Sprintf("%s/r/where-is-it/tweet-fail", fnAPIURL),
				TweetSuccess: fmt.Sprintf("%s/r/where-is-it/tweet-success", fnAPIURL),
			}
			err := doRequest(payload, detect, httpClient, fnToken)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (omega *OnionOmega2) PrintTweetInfo(tweet anaconda.Tweet) {
	fmt.Println(fmt.Sprintf(
		"[%v] found new tweet: %v from @%v.\n", tweet.CreatedAt, tweet.Text, tweet.User.ScreenName))
}
