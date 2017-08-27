package main

import (
	"encoding/json"
	"fmt"
	"github.com/ChimeraCoder/anaconda"
	"net/url"
	"os"
)

type RequestPayload struct {
	TweetID        string `json:"tweet_id"`
	User           string `json:"user"`
	BadImageSource bool   `json:"bad_image_source,omitempty"`
}

func main() {
	r := new(RequestPayload)
	err := json.NewDecoder(os.Stdin).Decode(r)
	if err != nil {
		fmt.Fprintf(os.Stderr, "unable to decode STDIN, got error %v", err.Error())
		return
	}

	consumerKey := os.Getenv("CONSUMER_KEY")
	consumerSecret := os.Getenv("CONSUMER_SECRET")
	apiToken := os.Getenv("API_KEY")
	apiTokenSecret := os.Getenv("API_KEY_SECRET")
	anaconda.SetConsumerKey(consumerKey)
	anaconda.SetConsumerSecret(consumerSecret)

	api := anaconda.NewTwitterApi(apiToken, apiTokenSecret)
	ok, err := api.VerifyCredentials()
	if !ok {
		fmt.Fprint(os.Stderr, "unable to verify credentials\n")
		return
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "unable to verify createndials, got error: %v\n", err.Error())
		return
	}
	v := url.Values{}
	v.Set("in_reply_to_status_id", r.TweetID)
	msg := ""
	if !r.BadImageSource {
		msg = fmt.Sprintf(`Hey, %v! Had to admit, i don't know where is it ¯\_(ツ)_/¯`, r.User)
	} else {
		msg = fmt.Sprintf(`Hey, %v! Seems like bad image ¯\_(ツ)_/¯, try another one...`, r.User)
	}
	_, err = api.PostTweet(msg, v)
	if err != nil {
		fmt.Fprintf(os.Stderr, "unable to post tweet for %v, got error %v", r.User, err.Error())
		return
	}
}
