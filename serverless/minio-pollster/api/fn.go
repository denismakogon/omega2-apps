package api

import (
	"context"
	"errors"
	"github.com/denismakogon/omega2-apps/serverless/twitter-daemon/api"
	"os"
)

func setupEmokognitionV2(ctx context.Context, appName string) (fnAPIURL, fnToken string, err error) {
	fnAPIURL, fnToken, fnClient, err := api.SetupFNClient()
	config := map[string]string{}
	config["FN_API_URL"] = os.Getenv("INTERNAL_FN_API_URL")

	err = api.RedeployFnApp(ctx, fnClient, appName, config)
	if err != nil {
		return "", "", err
	}

	err = api.RecreateRoute(ctx, fnClient, appName,
		"denismakogon/emokognition-v2:0.0.1",
		"/detect",
		"async",
		"json",
		"2000m",
		600, 200, uint64(1024))
	if err != nil {
		return "", "", errors.New(err.Error())
	}

	err = api.RecreateRoute(ctx, fnClient, appName,
		"denismakogon/emotion-recorder:0.0.11",
		"/recorder",
		"async",
		"http",
		"",
		120, 120, uint64(256))
	if err != nil {
		return "", "", errors.New(err.Error())
	}

	err = api.RecreateRoute(ctx, fnClient, appName,
		"denismakogon/emotion-results:0.0.8",
		"/results",
		"sync",
		"json",
		"",
		120, 120, uint64(512))
	if err != nil {
		return "", "", errors.New(err.Error())
	}

	return fnAPIURL, fnToken, nil
}
