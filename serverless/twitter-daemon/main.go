package main

import (
	"fmt"
	"github.com/denismakogon/omega2-apps/serverless/twitter-daemon/api"
	"net/http"
	"net/url"
	"sync"
	"time"
)

func main() {

	twitter := new(api.TwitterSecret)
	twitterAPI, err := twitter.FromEnv()
	if err != nil {
		panic(err.Error())
	}

	gc := new(api.GCloudSecret)
	err = gc.FromFile()
	if err != nil {
		panic(err.Error())
	}

	fnAPIURL, fnToken, err := api.SetupFunctions(gc, twitter)
	if err != nil {
		panic(err.Error())
	}

	httpClient := api.SetupHTTPClient()

	// get latest 200 tweets
	v := url.Values{}
	v.Set("count", "200")

	omega := api.OnionOmega2{
		TwitterAPI:   twitterAPI,
		SearchValues: v,
		GCloudAuth:   gc,
	}

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
						TweetIDInt64: tweet.Id,
						APIURL:       withSchemaAPI,
					}
					_, err = api.DoUncheckedRequest(payload, hotTweetDispatch, httpClient, fnToken)
					if err != nil {
						panic(err.Error())
					}
				}()
			}
			wg.Wait()
		}
		time.Sleep(time.Second * 6)
	}
}
