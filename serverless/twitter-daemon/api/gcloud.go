package api

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
)

type GCloudSecret struct {
	Mapper
	Type                    string `json:"type"`
	ProjectID               string `json:"project_id"`
	PrivateKeyID            string `json:"private_key_id"`
	PrivateKey              string `json:"private_key"`
	ClientEmail             string `json:"client_email"`
	ClientID                string `json:"client_id"`
	AuthURI                 string `json:"auth_uri"`
	TokenURI                string `json:"token_uri"`
	AuthProviderX509CertURL string `json:"auth_provider_x_509_cert_url"`
	ClientX509CertURL       string `json:"client_x509_cert_url"`
}

func (g *GCloudSecret) FromFile() error {
	gcloudCredsPath := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
	if gcloudCredsPath != "" {
		raw, err := ioutil.ReadFile(gcloudCredsPath)
		if err != nil {
			return err
		}
		err = json.Unmarshal(raw, g)
		if err != nil {
			return err
		}
		return nil
	} else {
		return errors.New("GOOGLE_APPLICATION_CREDENTIALS env var is not set")
	}
}
