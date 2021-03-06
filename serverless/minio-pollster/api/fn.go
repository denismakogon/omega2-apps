package api

import (
	"context"
	"errors"
	"os"

	"github.com/denismakogon/omega2-apps/serverless/minio-pollster/common"
)

func setupEmokognitionV2(ctx context.Context, appName string, config map[string]string) (fnAPIURL, fnToken string, err error) {
	fnAPIURL, fnToken, fnClient, err := common.SetupFNClient()

	config["FN_API_URL"] = os.Getenv("INTERNAL_FN_API_URL")

	err = common.RedeployFnApp(ctx, fnClient, appName, config)
	if err != nil {
		return "", "", err
	}

	err = common.RecreateRoute(ctx, fnClient, appName,
		"denismakogon/emokognition:0.0.8",
		"/detect",
		"async",
		"json",
		"2000m",
		600, 200, uint64(1500))
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
