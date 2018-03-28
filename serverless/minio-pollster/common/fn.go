package common

import (
	"context"
	"errors"
	"github.com/fnproject/fn_go/client"
	"github.com/fnproject/fn_go/client/apps"
	"github.com/fnproject/fn_go/client/routes"
	"github.com/fnproject/fn_go/models"
	openapi "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
	"net/url"
	"os"
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

func RecreateRoute(ctx context.Context, fnclient *client.Fn, appName, image, routePath, routeType, fformat, cpus string, timeout, idleTimeout int32, memory uint64) error {
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
				Cpus:        cpus,
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

func RedeployFnApp(ctx context.Context, fnclient *client.Fn, app string, config map[string]string) error {
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

func SetupFNClient() (string, string, *client.Fn, error) {
	fnAPIURL := os.Getenv("FN_API_URL")
	if fnAPIURL == "" {
		fnAPIURL = "http://localhost:8080"
	}
	u, err := url.Parse(fnAPIURL)
	if err != nil {
		return "", "", nil, err
	}

	fnToken := os.Getenv("FN_TOKEN")
	fnTransport := openapi.New(u.Host, "/v1", []string{u.Scheme})
	if fnToken != "" {
		fnTransport.DefaultAuthentication = openapi.BearerToken(fnToken)
	}
	// create the API client, with the transport
	fnclient := client.New(fnTransport, strfmt.Default)
	return fnAPIURL, fnToken, fnclient, nil
}
