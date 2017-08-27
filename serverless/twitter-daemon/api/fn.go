package api

import (
	"context"
	"fmt"
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

func recreateRoute(ctx context.Context, fnclient *client.Functions, appName, image, routePath, routeType string, timeout int32) error {
	cfg := &routes.PostAppsAppRoutesParams{
		App: appName,
		Body: &models.RouteWrapper{
			Route: &models.Route{
				Image:   image,
				Path:    routePath,
				Type:    routeType,
				Timeout: &timeout,
			},
		},
		Context: ctx,
	}
	// Renew or create from scratch route no matter exist it or not
	_, err := fnclient.Routes.DeleteAppsAppRoutesRoute(&routes.DeleteAppsAppRoutesRouteParams{
		App:     appName,
		Route:   routePath,
		Context: ctx,
	})
	if err != nil {
		// we should not fail here in case route does not exist
		fmt.Fprintf(os.Stdout, "Unable to delete route, got error %v", err.Error())
	}
	_, err = fnclient.Routes.PostAppsAppRoutes(cfg)
	if err != nil {
		return err
	}
	return nil
}

func setupAppAndRoutes(fnclient *client.Functions, gcloud *GCloudSecret, twitterSecret *TwitterSecret) error {
	app := "where-is-it"
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	config := map[string]string{}
	config, err := gcloud.Append(config)
	if err != nil {
		return err
	}
	config, err = twitterSecret.Append(config)
	if err != nil {
		return err
	}

	// Renew app or create from scratch
	_, err = fnclient.Apps.DeleteAppsApp(&apps.DeleteAppsAppParams{
		App:     app,
		Context: ctx,
	})
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
		return err
	}

	err = recreateRoute(ctx, fnclient, app, "denismakogon/tweet-fail:0.0.1",
		"/tweet-fail", "async", 60)
	if err != nil {
		return err
	}
	err = recreateRoute(ctx, fnclient, app, "denismakogon/detect-where:0.0.1",
		"/detect-where", "async", 60)
	if err != nil {
		return err
	}
	err = recreateRoute(ctx, fnclient, app, "denismakogon/tweet-success:0.0.1",
		"/tweet-success", "async", 60)
	if err != nil {
		return err
	}
	return nil
}

func SetupFunctions(gc *GCloudSecret, twitterSecret *TwitterSecret) (string, string, error) {
	fnAPIURL := os.Getenv("API_URL")
	if fnAPIURL == "" {
		fnAPIURL = "http://localhost:8080"
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
