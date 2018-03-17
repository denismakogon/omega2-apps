package api

import (
	"context"
	"errors"
	"github.com/denismakogon/omega2-apps/serverless/common"
	"os"
)

func setupEmokognitionV2(ctx context.Context, appName string) (fnAPIURL, fnToken string, err error) {
	fnAPIURL, fnToken, fnClient, err := common.SetupFNClient()
	config := map[string]string{}

	pgConf := new(common.PostgresConfig)
	err = pgConf.FromEnv()
	if err != nil {
		return "", "", err
	}

	config, err = common.Append(pgConf, config)
	if err != nil {
		return "", "", err
	}

	config["FN_API_URL"] = os.Getenv("INTERNAL_FN_API_URL")

	err = common.RedeployFnApp(ctx, fnClient, appName, config)
	if err != nil {
		return "", "", err
	}

	err = common.RecreateRoute(ctx, fnClient, appName,
		"denismakogon/emokognition-v2:0.0.1",
		"/detect-v2",
		"async",
		"json",
		"2000m",
		600, 200, uint64(1024))
	if err != nil {
		return "", "", errors.New(err.Error())
	}

	err = common.RecreateRoute(ctx, fnClient, appName,
		"denismakogon/emotion-recorder:0.0.11",
		"/recorder",
		"async",
		"http",
		"",
		120, 120, uint64(256))
	if err != nil {
		return "", "", errors.New(err.Error())
	}

	err = common.RecreateRoute(ctx, fnClient, appName,
		"denismakogon/emotion-results:0.0.8",
		"/results",
		"sync",
		"json",
		"",
		120, 120, uint64(512))
	if err != nil {
		return "", "", errors.New(err.Error())
	}

	err = common.RecreateRoute(ctx, fnClient, appName,
		"denismakogon/emokognition-view:0.0.13",
		"/index.html",
		"sync",
		"json",
		"",
		120, 200, uint64(512))
	if err != nil {
		return "", "", errors.New(err.Error())
	}

	return fnAPIURL, fnToken, nil
}
