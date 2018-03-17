package main

import (
	"fmt"
	"github.com/ChimeraCoder/anaconda"
	"github.com/denismakogon/omega2-apps/serverless/common"
	"net/http"
	"net/url"
	"os"
	"sync"
	"time"
)

func asyncRunner(omega *common.OnionOmega2, fnAPIURL, fnToken string, proc func(tweet anaconda.Tweet, httpClient *http.Client, fnAPIURL, fnToken string) error) {
	httpClient := common.SetupHTTPClient()

	tweetID := os.Getenv("InitialTweetID")
	if tweetID == "" {
		panic("Initial tweet ID env var is not set, but suppose to be!")
	}
	omega.SetTweetIDToStartFrom(tweetID)
	wg := new(sync.WaitGroup)

	for {
		tweets, err := omega.GetRecentMentions()
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}
		if len(tweets) != 0 {
			wg.Add(len(tweets))
			for _, tweet := range tweets {
				omega.PrintTweetInfo(tweet)
				go func() {
					defer wg.Done()
					err = proc(tweet, httpClient, fnAPIURL, fnToken)
					if err != nil {
						fmt.Fprintln(os.Stderr, err.Error())
						os.Exit(1)
					}
				}()
			}
		}
		time.Sleep(time.Second * 6)
	}
}

func EmotionRecognition() error {
	pgConf := new(common.PostgresConfig)
	err := pgConf.FromEnv()
	if err != nil {
		return err
	}

	twitter := new(common.TwitterSecret)
	twitterAPI, err := twitter.FromEnv()
	if err != nil {
		return err
	}

	fnAPIURL, fnToken, err := common.SetupEmoKognitionFunctions(twitter, pgConf)
	if err != nil {
		return err
	}

	// get latest 10 tweets fro InitialTweet
	v := url.Values{}
	v.Set("count", "200")

	omega := common.OnionOmega2{
		TwitterAPI:   twitterAPI,
		SearchValues: &v,
	}
	asyncRunner(&omega, fnAPIURL, fnToken, common.ProcessTweetWithEmotion)
	return nil
}

func LandmarkRecognition() error {
	twitter := new(common.TwitterSecret)
	twitterAPI, err := twitter.FromEnv()
	if err != nil {
		return err
	}

	gc := new(common.GCloudSecret)
	err = gc.FromEnv()
	if err != nil {
		return err
	}

	fnAPIURL, fnToken, err := common.SetupLandmarkRecognitionFunctions(gc, twitter)
	if err != nil {
		return err
	}

	// get latest 200 tweets
	v := url.Values{}
	v.Set("count", "200")

	omega := common.OnionOmega2{
		TwitterAPI:   twitterAPI,
		SearchValues: &v,
		GCloudAuth:   gc,
	}
	asyncRunner(&omega, fnAPIURL, fnToken, common.ProcessTweetWithLandmark)
	return nil
}

func main() {
	botType := os.Getenv("TwitterBotType")
	if botType == "landmark" {
		err := LandmarkRecognition()
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}
	}
	if botType == "emokognition" {
		err := EmotionRecognition()
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}
	}
	if botType == "" {
		fmt.Fprintln(os.Stderr, "Recognition type is not set.")
		os.Exit(1)
	}
}
