package api

import (
	"context"
	"errors"
	"fmt"
	"github.com/fnproject/fn_go/client"
	"github.com/fnproject/fn_go/client/apps"
	"github.com/fnproject/fn_go/client/routes"
	"github.com/fnproject/fn_go/models"
	openapi "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
	"os"
	"time"
)

type CallID struct {
	ID string `json:"call_id"`
}

type ErrMessage struct {
	Message string `json:"message"`
}

type ErrBody struct {
	Error ErrMessage `json:"error"`
}

func recreateRoute(ctx context.Context, fnclient *client.Fn, appName, image, routePath, routeType, fformat string, timeout, idleTimeout int32, memory uint64) error {
	cfg := &routes.PostAppsAppRoutesParams{
		App: appName,
		Body: &models.RouteWrapper{
			Route: &models.Route{
				Image:       image,
				Path:        routePath,
				Type:        routeType,
				Timeout:     &timeout,
				Memory:      memory,
				Format:      fformat,
				IDLETimeout: &idleTimeout,
			},
		},
		Context: ctx,
	}
	_, err := fnclient.Routes.PostAppsAppRoutes(cfg)
	if err != nil {
		return errors.New(err.Error())
	}
	return nil
}

func redeployFnApp(ctx context.Context, fnclient *client.Fn, app string, config map[string]string) error {
	_, err := fnclient.Apps.GetAppsApp(&apps.GetAppsAppParams{
		App:     app,
		Context: ctx,
	})

	if err == nil {
		appRoutes, err := fnclient.Routes.GetAppsAppRoutes(&routes.GetAppsAppRoutesParams{
			App:     app,
			Context: ctx,
		})
		if err != nil {
			return errors.New(err.Error())
		}
		// dropping all routes
		if len(appRoutes.Payload.Routes) != 0 {
			for _, route := range appRoutes.Payload.Routes {
				fnclient.Routes.DeleteAppsAppRoutesRoute(&routes.DeleteAppsAppRoutesRouteParams{
					App:     app,
					Route:   route.Path,
					Context: ctx,
				})
			}
		}
	}
	// deleting app
	fnclient.Apps.DeleteAppsApp(&apps.DeleteAppsAppParams{
		App:     app,
		Context: ctx,
	})
	// creating from scratch
	_, err = fnclient.Apps.PostApps(&apps.PostAppsParams{
		Body: &models.AppWrapper{
			App: &models.App{
				Config: config,
				Name:   app,
			},
		},
		Context: ctx,
	})
	return err
}

func setupEmokognitionAppAndRoutes(fnAPIURL string, fnclient *client.Fn, twitterSecret *TwitterSecret, pgConfig *PostgresConfig) error {
	app := "emokognition"
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()

	config := map[string]string{}
	config, err := Append(twitterSecret, config)
	if err != nil {
		return err
	}
	config, err = Append(pgConfig, config)
	if err != nil {
		return err
	}
	config["FN_API_URL"] = fmt.Sprintf("http://%v", fnAPIURL)

	err = redeployFnApp(ctx, fnclient, app, config)
	if err != nil {
		return err
	}
	err = recreateRoute(ctx, fnclient, app,
		"denismakogon/emotion-recorder:0.0.11",
		"/recorder",
		"async",
		"http",
		60, 120, uint64(256))
	if err != nil {
		return errors.New(err.Error())
	}
	err = recreateRoute(ctx, fnclient, app,
		"denismakogon/emotion-results:0.0.3",
		"/results",
		"sync",
		"default",
		60, 120, uint64(126))
	if err != nil {
		return errors.New(err.Error())
	}
	err = recreateRoute(ctx, fnclient, app,
		"denismakogon/emokognition:0.0.5",
		"/detect",
		"async",
		"default",
		60, 200, uint64(2048))
	if err != nil {
		return errors.New(err.Error())
	}
	err = recreateRoute(ctx, fnclient, app,
		"denismakogon/emokognition-view:0.0.9",
		"/index.html",
		"sync",
		"default",
		60, 120, uint64(256))
	if err != nil {
		return errors.New(err.Error())
	}

	return nil
}

func setupLandmarkAppAndRoutes(fnclient *client.Fn, gcloud *GCloudSecret, twitterSecret *TwitterSecret) error {
	app := "where-is-it"
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()

	config := map[string]string{}
	config, err := Append(gcloud, config)
	if err != nil {
		return err
	}
	config, err = Append(twitterSecret, config)
	if err != nil {
		return err
	}

	err = redeployFnApp(ctx, fnclient, app, config)
	if err != nil {
		return err
	}

	err = recreateRoute(ctx, fnclient, app,
		"denismakogon/tweet-fail:0.0.2",
		"/tweet-fail",
		"async",
		"default",
		60, 120, uint64(256))
	if err != nil {
		return errors.New(err.Error())
	}
	err = recreateRoute(ctx, fnclient, app,
		"denismakogon/detect-task:0.0.5",
		"/detect-where",
		"async",
		"default",
		60, 120, uint64(256))
	if err != nil {
		return errors.New(err.Error())
	}
	err = recreateRoute(ctx, fnclient, app,
		"denismakogon/tweet-success:0.0.2",
		"/tweet-success",
		"async",
		"default",
		60, 120, uint64(256))
	if err != nil {
		return errors.New(err.Error())
	}
	return nil
}

func setupFNClient() (string, string, *client.Fn) {
	fnAPIURL := os.Getenv("API_URL")
	if fnAPIURL == "" {
		fnAPIURL = "localhost:8080"
	}
	fnToken := os.Getenv("FN_TOKEN")
	fnTransport := openapi.New(fnAPIURL, "/v1", []string{"http"})
	if fnToken != "" && os.Getenv("API_URL") != "" {
		fnTransport.DefaultAuthentication = openapi.BearerToken(fnToken)
	}
	// create the API client, with the transport
	fnclient := client.New(fnTransport, strfmt.Default)
	return fnAPIURL, fnToken, fnclient
}

func SetupEmoKognitionFunctions(twitterSecret *TwitterSecret, pgConfig *PostgresConfig) (string, string, error) {
	fnAPIURL, fnToken, fnclient := setupFNClient()
	err := setupEmokognitionAppAndRoutes(fnAPIURL, fnclient, twitterSecret, pgConfig)
	if err != nil {
		return "", "", err
	}
	return fnAPIURL, fnToken, nil
}

func SetupLandmarkRecognitionFunctions(gc *GCloudSecret, twitterSecret *TwitterSecret) (string, string, error) {
	fnAPIURL, fnToken, fnclient := setupFNClient()
	err := setupLandmarkAppAndRoutes(fnclient, gc, twitterSecret)
	if err != nil {
		return "", "", err
	}
	return fnAPIURL, fnToken, nil
}
