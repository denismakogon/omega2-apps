package api

import (
	"context"
	"errors"
	"github.com/funcy/functions_go/client"
	"github.com/funcy/functions_go/client/apps"
	"github.com/funcy/functions_go/client/routes"
	"github.com/funcy/functions_go/models"
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

func recreateRoute(ctx context.Context, fnclient *client.Functions, appName, image, routePath, routeType, fformat string, timeout, idleTimeout int32) error {
	cfg := &routes.PostAppsAppRoutesParams{
		App: appName,
		Body: &models.RouteWrapper{
			Route: &models.Route{
				Image:       image,
				Path:        routePath,
				Type:        routeType,
				Timeout:     &timeout,
				Memory:      uint64(256),
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

func setupAppAndRoutes(fnclient *client.Functions, gcloud *GCloudSecret, twitterSecret *TwitterSecret) error {
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

	_, err = fnclient.Apps.GetAppsApp(&apps.GetAppsAppParams{
		App:     app,
		Context: ctx,
	})
	// app exists
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
	if err != nil {
		return errors.New(err.Error())
	}

	err = recreateRoute(ctx, fnclient, app,
		"denismakogon/tweet-fail:0.0.2",
		"/tweet-fail",
		"async",
		"default",
		60, 120)
	if err != nil {
		return errors.New(err.Error())
	}
	err = recreateRoute(ctx, fnclient, app,
		"denismakogon/detect-task:0.0.5",
		"/detect-where",
		"async",
		"default",
		60, 120)
	if err != nil {
		return errors.New(err.Error())
	}
	err = recreateRoute(ctx, fnclient, app,
		"denismakogon/tweet-success:0.0.2",
		"/tweet-success",
		"async",
		"default",
		60, 120)
	if err != nil {
		return errors.New(err.Error())
	}
	err = recreateRoute(ctx, fnclient, app,
		"denismakogon/tweet-dispatch:0.0.1",
		"/tweet-dispatch",
		"sync",
		"http",
		60, 120)
	if err != nil {
		return errors.New(err.Error())
	}
	return nil
}

func SetupFunctions(gc *GCloudSecret, twitterSecret *TwitterSecret) (string, string, error) {
	fnAPIURL := os.Getenv("API_URL")
	if fnAPIURL == "" {
		fnAPIURL = "localhost:8080"
	}
	fnToken := os.Getenv("FN_TOKEN")
	fnTransport := openapi.New(fnAPIURL, "/v1", []string{"http"})
	// This means that FN token is not required if FN is local
	if fnToken != "" && os.Getenv("API_URL") != "" {
		fnTransport.DefaultAuthentication = openapi.BearerToken(fnToken)
	}
	// create the API client, with the transport
	fnclient := client.New(fnTransport, strfmt.Default)
	err := setupAppAndRoutes(fnclient, gc, twitterSecret)
	if err != nil {
		return "", "", err
	}
	return fnAPIURL, fnToken, nil
}
