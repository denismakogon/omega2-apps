package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/denismakogon/omega2-apps/serverless/twitter-daemon/api"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
)

func writeBadResponse(buf *bytes.Buffer, resp *http.Response, err error) {
	resp.StatusCode = 500
	resp.Status = http.StatusText(resp.StatusCode)
	fmt.Fprintln(buf, err)
	fmt.Fprintf(os.Stderr, err.Error())
}

func main() {
	httpClient := api.SetupHTTPClient()
	twitter := new(api.TwitterSecret)
	twitterAPI, err := twitter.FromEnv()
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		os.Exit(0)
	}
	for {
		res := http.Response{
			Proto:      "HTTP/1.1",
			ProtoMajor: 1,
			ProtoMinor: 1,
			StatusCode: 200,
			Status:     "OK",
		}
		r := bufio.NewReader(os.Stdin)
		req, err := http.ReadRequest(r)
		var buf bytes.Buffer
		if err != nil {
			writeBadResponse(&buf, &res, err)
		} else {
			ok, err := twitterAPI.VerifyCredentials()
			if !ok {
				writeBadResponse(
					&buf, &res, errors.New("Unable to authorize to twitter."))
			}
			if err != nil {
				writeBadResponse(&buf, &res, err)
			} else {
				l, _ := strconv.Atoi(req.Header.Get("Content-Length"))
				p := make([]byte, l)
				_, err = r.Read(p)
				if err != nil {
					writeBadResponse(&buf, &res, err)
				} else {
					payload := &api.RequestPayload{}
					err = json.Unmarshal(p, payload)
					if err != nil {
						writeBadResponse(&buf, &res, err)
					} else {
						tweet, err := twitterAPI.GetTweet(payload.TweetIDInt64, nil)
						if err == nil {
							err := api.ProcessTweet(tweet, httpClient, payload.APIURL, "")
							if err != nil {
								writeBadResponse(&buf, &res, err)
							} else {
								fmt.Fprint(&buf, "OK\n")
								res.Body = ioutil.NopCloser(&buf)
								res.ContentLength = int64(buf.Len())
								res.Write(os.Stdout)
							}
						}
					}
				}
			}
		}
	}
}
