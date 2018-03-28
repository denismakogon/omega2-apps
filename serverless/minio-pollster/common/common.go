package common

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type RequestPayload struct {
	MediaURL string `json:"media_url"`
}

func SetupHTTPClient() *http.Client {
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
	return &http.Client{Transport: transport}
}

func StructFromEnv(i interface{}) error {
	v := reflect.ValueOf(i)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	typ := v.Type()
	for i := 0; i < v.NumField(); i++ {
		fi := typ.Field(i)
		if tagValue := fi.Tag.Get("json"); tagValue != "" {
			value := os.Getenv(strings.ToUpper(tagValue))
			if value == "" {
				return fmt.Errorf("missing env var value: %s", strings.ToUpper(tagValue))
			}
			v.FieldByName(fi.Name).SetString(value)
		}
	}
	return nil
}

func StructFromFile(i interface{}, envVar string) error {
	fPath := os.Getenv(envVar)
	if fPath != "" {
		raw, err := ioutil.ReadFile(fPath)
		if err != nil {
			return err
		}
		err = json.Unmarshal(raw, i)
		if err != nil {
			return err
		}
		return nil
	} else {
		return fmt.Errorf("%v env var is not set", envVar)
	}
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
			if v.Field(i).Type() == reflect.TypeOf(true) {
				out[strings.ToUpper(tagValue)] = strconv.FormatBool(v.Field(i).Bool())
			} else {
				out[strings.ToUpper(tagValue)] = v.Field(i).String()
			}
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

func doUncheckedRequest(payload *RequestPayload, req *http.Request, httpClient *http.Client, fnToken string) (*http.Response, error) {
	if fnToken != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", fnToken))
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req.ContentLength = int64(len(body))
	req.Body = ioutil.NopCloser(bytes.NewReader(body))
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func DoRequest(log *logrus.Entry, payload *RequestPayload, req *http.Request,
	httpClient *http.Client, fnToken string) error {
	resp, err := doUncheckedRequest(payload, req, httpClient, fnToken)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusAccepted || resp.StatusCode == http.StatusOK {
		callID := new(CallID)
		err = json.NewDecoder(resp.Body).Decode(&callID)
		if err != nil {
			return err
		}
		log.Infof("New detect func submitted. Call ID: %v", callID.ID)
	} else {
		apiError := new(ErrBody)
		err = json.NewDecoder(resp.Body).Decode(&apiError)
		if err != nil {
			log.Errorf("Error during detect func submittion. Call ID: %s\n. Error: %s",
				apiError.Error.Message, err.Error())
			return err
		}
	}
	return nil
}
