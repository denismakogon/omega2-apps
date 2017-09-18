package api

import (
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

func (twitter *TwitterSecret) FromFile() (*anaconda.TwitterApi, error) {
	envVar := "TWITTER_APPLICATION_CREDENTIALS"
	err := StructFromFile(twitter, envVar)
	if err != nil {
		return nil, err
	}
	anaconda.SetConsumerKey(twitter.ConsumerKey)
	anaconda.SetConsumerSecret(twitter.ConsumerSecret)
	api := anaconda.NewTwitterApi(twitter.APIToken, twitter.APITokenSecret)
	return api, nil
}

func (twitter *TwitterSecret) FromEnv() (*anaconda.TwitterApi, error) {
	err := StructFromEnv(twitter)
	if err != nil {
		return nil, err
	}
	anaconda.SetConsumerKey(twitter.ConsumerKey)
	anaconda.SetConsumerSecret(twitter.ConsumerSecret)
	api := anaconda.NewTwitterApi(twitter.APIToken, twitter.APITokenSecret)
	return api, nil
}

func (omega *OnionOmega2) GetRecentMentions() (tweets []anaconda.Tweet, err error) {
	tweets, err = omega.TwitterAPI.GetMentionsTimeline(*omega.SearchValues)
	if err != nil {
		return nil, err
	}
	if len(tweets) != 0 {
		// seems like tweets are now ordered from recent to oldest
		since_id := strconv.FormatInt(tweets[0].Id, 10)
		omega.SearchValues.Set("since_id", since_id)
	}
	return tweets, nil
}

func ProcessTweetWithEmotion(tweet anaconda.Tweet, httpClient *http.Client, fnAPIURL, fnToken string) error {
	detect, err := http.NewRequest(
		http.MethodPost, fmt.Sprintf("http://%s/r/emokognition/detect", fnAPIURL),
		nil)
	if err != nil {
		return err
	}
	for _, media := range tweet.Entities.Media {
		payload := &RequestPayload{
			MediaURL: media.Media_url,
		}
		err := DoRequest(payload, detect, httpClient, fnToken)
		if err != nil {
			return err
		}
	}
	return nil
}

func ProcessTweetWithLandmark(tweet anaconda.Tweet, httpClient *http.Client, fnAPIURL, fnToken string) error {
	detect, err := http.NewRequest(
		http.MethodPost, fmt.Sprintf("http://%s/r/landmark/detect-where", fnAPIURL),
		nil)
	if err != nil {
		return err
	}

	fail, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/r/landmark/tweet-fail", fnAPIURL), nil)
	if err != nil {
		return err
	}

	for _, media := range tweet.Entities.Media {
		user := fmt.Sprintf("@%v", tweet.User.ScreenName)
		if media.Type != "photo" {
			payload := &RequestPayload{
				User:    user,
				TweetID: tweet.IdStr,
			}
			err := DoRequest(payload, fail, httpClient, fnToken)
			if err != nil {
				return err
			}
		} else {
			payload := &RequestPayload{
				MediaURL:     media.Media_url,
				User:         user,
				TweetID:      tweet.IdStr,
				TweetFail:    fmt.Sprintf("%s/r/landmark/tweet-fail", fnAPIURL),
				TweetSuccess: fmt.Sprintf("%s/r/landmark/tweet-success", fnAPIURL),
			}
			err := DoRequest(payload, detect, httpClient, fnToken)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (omega *OnionOmega2) PrintTweetInfo(tweet anaconda.Tweet) {
	hasMedia := false
	if len(tweet.Entities.Media) != 0 {
		hasMedia = true
	}
	fmt.Println(fmt.Sprintf(
		"[%v] found new tweet: %v from @%v. Media included? - %v\n",
		tweet.CreatedAt, tweet.Text, tweet.User.ScreenName, hasMedia))
}
