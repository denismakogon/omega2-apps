package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/ChimeraCoder/anaconda"
	"io/ioutil"
	"net/http"
	"net/url"
)

type RequestPayload struct {
	TweetID  string       `json:"tweet_id"`
	MediaURL string       `json:"media_url,omitempty"`
	User     string       `json:"user"`
	TweetFail string `json:"tweet_fail,omitempty"`
	TweetSuccess string `json:"tweet_success,omitempty"`
}

type OnionOmega2 struct {
	TwitterAPI   *anaconda.TwitterApi
	GCloudAuth   *GCloudSecret
	SearchValues url.Values
}

type Mapper struct {
}

func (m *Mapper) ToMap() (map[string]string, error) {
	var inInterface interface{}
	inrec, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(inrec, &inInterface)
	if err != nil {
		return nil, err
	}

	return inInterface.(map[string]string), nil
}

func (m *Mapper) Append(config map[string]string) (map[string]string, error) {
	mMap, err := m.ToMap()
	if err != nil {
		return nil, err
	}
	for key, value := range mMap {
		config[key] = value
	}
	return config, nil
}

func doRequest(payload *RequestPayload, req *http.Request, httpClient *http.Client, fnToken string) error {
	if fnToken != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", fnToken))
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	req.Body = ioutil.NopCloser(bytes.NewReader(body))
	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode == 200 {
		callID := new(CallID)
		err = json.NewDecoder(resp.Body).Decode(&callID)
		if err != nil {
			return err
		}
		fmt.Printf("New detect func submitted. Call ID: %v", callID.ID)
	} else {
		apiErr := new(ErrBody)
		err = json.NewDecoder(resp.Body).Decode(&apiErr)
		if err != nil {
			return err
		}
		fmt.Printf("Error during detect func submittion. Call ID: %v", apiErr.Error.Message)
	}
	return nil
}
