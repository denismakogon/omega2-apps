package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/denismakogon/omega2-apps/serverless/twitter-daemon/api"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
)

func writeBadResponse(buf *bytes.Buffer, resp *http.Response, errMsg string) {
	resp.StatusCode = 500
	resp.Status = http.StatusText(resp.StatusCode)
	fmt.Fprintln(buf, errMsg)
	fmt.Fprintf(os.Stderr, errMsg)
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
			writeBadResponse(&buf, &res,
				fmt.Sprintf("Unable to read request from STDIN, "+
					"it might be empty. Error: %v", err.Error()))
		} else {
			ok, err := twitterAPI.VerifyCredentials()
			if !ok {
				writeBadResponse(
					&buf, &res, "Unable to authorize to twitter.")
			}
			if err != nil {
				writeBadResponse(&buf, &res,
					fmt.Sprintf("Unable to authorize to twitter. "+
						"Error: %v", err.Error()))
			} else {
				for k, v := range req.Header {
					fmt.Fprintf(os.Stderr, "%s: %s\n", k, v)
				}
				l, _ := strconv.Atoi(req.Header.Get("Content-Length"))

				p := make([]byte, l)
				_, err = r.Read(p)
				if err != nil {
					writeBadResponse(&buf, &res,
						fmt.Sprintf("Unable to read request data. "+
							"Error: %v", err.Error()))
				} else {
					payload := &api.RequestPayload{}
					err = json.Unmarshal(p, payload)
					if err != nil {
						writeBadResponse(&buf, &res,
							fmt.Sprintf("Unable to unmarshal request data. "+
								"Error: %v", err.Error()))
					} else {
						fmt.Fprintf(os.Stderr, fmt.Sprintf("TweetID: %v", payload.TweetIDInt64))
						tweet, err := twitterAPI.GetTweet(payload.TweetIDInt64, nil)
						if err != nil {
							writeBadResponse(&buf, &res,
								fmt.Sprintf("Unable to get tweet. "+
									"Error: %v", err.Error()))
						} else {
							if payload.RecognitionType == "landmark" {
								err = api.ProcessTweetWithLandmark(tweet, httpClient, payload.APIURL, "")
								if err != nil {
									writeBadResponse(&buf, &res,
										fmt.Sprintf("Unable to submit tweet processing. "+
											"Error: %v", err.Error()))
								} else {
									fmt.Fprint(&buf, "OK landmark\n")
								}
							}
							fmt.Fprintf(os.Stderr, "Recognition type: %s\n", payload.RecognitionType)
							if payload.RecognitionType == "emokognition" {
								fmt.Fprintln(os.Stderr, "entering emokognition")
								err = api.ProcessTweetWithEmotion(tweet, httpClient, payload.APIURL, "")
								if err != nil {
									writeBadResponse(&buf, &res,
										fmt.Sprintf("Unable to submit tweet processing. "+
											"Error: %v", err.Error()))
								} else {
									fmt.Fprint(&buf, "OK\n")
								}
							}
						}
					}
				}
			}
		}
		res.Body = ioutil.NopCloser(&buf)
		res.ContentLength = int64(buf.Len())
		res.Write(os.Stdout)
	}
}
