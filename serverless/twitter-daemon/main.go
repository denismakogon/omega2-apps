package main

import (
	"fmt"
	"github.com/denismakogon/omega2-apps/serverless/twitter-daemon/api"
	"net/http"
	"net/url"
	"os"
	"sync"
	"time"
)

func asyncRunner(omega *api.OnionOmega2, recognitionType, fnAPIURL, fnToken string) {
	httpClient := api.SetupHTTPClient()

	tweetID := os.Getenv("InitialTweetID")
	if tweetID == "" {
		// start to look for tweets from the very beginning
		panic("Initial tweet ID env var is not set, but suppose to be!")
	}
	omega.SetTweetIDToStartFrom(tweetID)
	wg := new(sync.WaitGroup)
	withSchemaAPI := fmt.Sprintf("http://%v", fnAPIURL)
	for {
		ok, err := omega.TwitterAPI.VerifyCredentials()
		if !ok {
			fmt.Println(err.Error())
			panic(err.Error())
		}
		tweets, err := omega.GetRecentMentions()
		if err != nil {
			fmt.Println(err.Error())
			panic(err.Error())
		}
		if len(tweets) != 0 {
			wg.Add(len(tweets))
			for _, tweet := range tweets {
				omega.PrintTweetInfo(tweet)
				go func() {
					defer wg.Done()
					hotTweetDispatch, err := http.NewRequest(
						http.MethodPost, fmt.Sprintf("%s/r/where-is-it/tweet-dispatch", withSchemaAPI),
						nil)
					if err != nil {
						panic(err.Error())
					}
					payload := &api.RequestPayload{
						TweetIDInt64:    tweet.Id,
						APIURL:          withSchemaAPI,
						RecognitionType: recognitionType,
					}
					_, err = api.DoUncheckedRequest(payload, hotTweetDispatch, httpClient, fnToken)
					if err != nil {
						panic(err.Error())
					}
				}()
			}
			wg.Wait()
		}
		time.Sleep(time.Second * 3)
	}
}

func EmotionRecognition() {
	pgConf := new(api.PostgresConfig)
	pgConf.FromFile()
	twitter := new(api.TwitterSecret)
	twitterAPI, err := twitter.FromFile()
	if err != nil {
		panic(err.Error())
	}
	fnAPIURL, fnToken, err := api.SetupEmoKognitionFunctions(twitter, pgConf)
	if err != nil {
		panic(err.Error())
	}

	// get latest 200 tweets fro InitialTweet
	v := url.Values{}
	v.Set("count", "200")

	omega := api.OnionOmega2{
		TwitterAPI:   twitterAPI,
		SearchValues: &v,
	}
	asyncRunner(&omega, "emokognition", fnAPIURL, fnToken)
}

func LandmarkRecognition() {

	twitter := new(api.TwitterSecret)
	twitterAPI, err := twitter.FromFile()
	if err != nil {
		panic(err.Error())
	}

	gc := new(api.GCloudSecret)
	err = gc.FromFile()
	if err != nil {
		panic(err.Error())
	}

	fnAPIURL, fnToken, err := api.SetupLandmarkRecognitionFunctions(gc, twitter)
	if err != nil {
		panic(err.Error())
	}

	// get latest 200 tweets
	v := url.Values{}
	v.Set("count", "200")

	omega := api.OnionOmega2{
		TwitterAPI:   twitterAPI,
		SearchValues: &v,
		GCloudAuth:   gc,
	}
	asyncRunner(&omega, "landmark", fnAPIURL, fnToken)
}

func main() {
	botType := os.Getenv("TwitterBotType")
	if botType == "landmark" {
		LandmarkRecognition()
	}
	if botType == "emokognition" {
		EmotionRecognition()
	}
	if botType == "" {
		panic("Recognition type is not set.")
	}
}
