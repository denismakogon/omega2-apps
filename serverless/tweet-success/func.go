package main

import (
	"encoding/json"
	"fmt"
	"github.com/denismakogon/omega2-apps/serverless/twitter-daemon/api"
	"net/url"
	"os"
)

type RequestPayload struct {
	TweetID  string `json:"tweet_id"`
	User     string `json:"user"`
	Landmark string `json:"landmark"`
}

func main() {
	r := new(RequestPayload)
	err := json.NewDecoder(os.Stdin).Decode(r)
	if err != nil {
		fmt.Fprintf(os.Stderr, "unable to decode STDIN, got error %v", err.Error())
		return
	}
	twitter := new(api.TwitterSecret)
	twitterAPI, err := twitter.FromEnv()

	ok, err := twitterAPI.VerifyCredentials()
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
	_, err = twitterAPI.PostTweet(fmt.Sprintf("Hey, %v! That should be %v!", r.User, r.Landmark), v)
	if err != nil {
		fmt.Fprintf(os.Stderr, "unable to post tweet for %v, got error %v", r.User, err.Error())
		return
	}
}
