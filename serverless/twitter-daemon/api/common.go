package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/ChimeraCoder/anaconda"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
)

type RequestPayload struct {
	TweetID      string `json:"tweet_id"`
	MediaURL     string `json:"media_url,omitempty"`
	User         string `json:"user"`
	TweetFail    string `json:"tweet_fail,omitempty"`
	TweetSuccess string `json:"tweet_success,omitempty"`
}

type OnionOmega2 struct {
	TwitterAPI   *anaconda.TwitterApi
	GCloudAuth   *GCloudSecret
	SearchValues url.Values
}

func ToMap(in interface{}) (map[string]interface{}, error) {
	out := make(map[string]interface{})

	v := reflect.ValueOf(in)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	typ := v.Type()
	for i := 0; i < v.NumField(); i++ {
		fi := typ.Field(i)
		if tagValue := fi.Tag.Get("json"); tagValue != "" {
			out[tagValue] = v.Field(i).String()
		}
	}
	return out, nil
}

func Append(obj interface{}, config map[string]string) (map[string]string, error) {
	mMap, err := ToMap(obj)
	if err != nil {
		return nil, err
	}
	for key, value := range mMap {
		config[key] = value.(string)
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
	if resp.StatusCode == http.StatusAccepted || resp.StatusCode == http.StatusOK {
		callID := new(CallID)
		err = json.NewDecoder(resp.Body).Decode(&callID)
		if err != nil {
			return err
		}
		fmt.Printf("New detect func submitted. Call ID: %v", callID.ID)
	} else {
		apiError := new(ErrBody)
		err = json.NewDecoder(resp.Body).Decode(&apiError)
		if err != nil {
			return err
		}
		fmt.Printf("Error during detect func submittion. Call ID: %v", apiError.Error.Message)
	}
	return nil
}
