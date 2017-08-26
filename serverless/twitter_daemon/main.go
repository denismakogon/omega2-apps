package main

import (
	"crypto/tls"
	"fmt"
	"github.com/denismakogon/omega2-apps/serverless/twitter_daemon/api"
	"net"
	"net/http"
	"net/url"
	"time"
	"sync"
)

func main() {

	twitter := new(api.TwitterSecret)
	twitterAPI, err := twitter.FromEnv()
	if err != nil {
		fmt.Println(err.Error())
		panic(err.Error())
	}

	gc := new(api.GCloudSecret)
	err = gc.FromFile()
	if err != nil {
		fmt.Println(err.Error())
		panic(err.Error())
	}

	fnAPIURL, fnToken, err := api.SetupFunctions(gc, twitter)

	transport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		Dial: (&net.Dialer{
			Timeout:   10 * time.Second,
			KeepAlive: 120 * time.Second,
		}).Dial,
		MaxIdleConnsPerHost: 512,
		TLSHandshakeTimeout: 10 * time.Second,
		TLSClientConfig: &tls.Config{
			ClientSessionCache: tls.NewLRUClientSessionCache(4096),
		},
	}
	httpClient := &http.Client{Transport: transport}

	// get latest 200 tweets
	v := url.Values{}
	v.Set("count", "200")

	omega := api.OnionOmega2{
		TwitterAPI:   twitterAPI,
		SearchValues: v,
		GCloudAuth:   gc,
	}

	wg := new(sync.WaitGroup)

	for {
		// make this async as hell with WaitGroup
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
				go omega.ProcessTweets(wg, tweet, httpClient, fnAPIURL, fnToken)
			}
			wg.Wait()
		}
		time.Sleep(time.Second * 6)
	}
}
