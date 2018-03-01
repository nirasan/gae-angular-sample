package app

import (
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"github.com/labstack/echo"
	"os"
	"net/http"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
)

func OauthStartHandler(e echo.Context) error {
	c := createOauth2Config()
	url := c.AuthCodeURL("state1", oauth2.AccessTypeOnline)
	http.Redirect(e.Response(), e.Request(), url, 302)
	return nil
}

func OauthCallbackHandler(e echo.Context) error {
	ctx := appengine.NewContext(e.Request())

	//TODO check state
	state := e.QueryParam("state")
	log.Debugf(ctx, "state: %v", state)

	code := e.QueryParam("code")

	c := createOauth2Config()
	tok, err := c.Exchange(ctx, code)
	if err != nil {
		panic(err)
	}
	log.Debugf(ctx, "token: %v", tok)

	return e.JSON(http.StatusOK, struct{ Message string }{"ok"})
}

func createOauth2Config() oauth2.Config {
	return oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		Scopes:       []string{"openid"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  google.Endpoint.AuthURL,
			TokenURL: google.Endpoint.TokenURL,
		},
		RedirectURL: "http://localhost:8080/oauth/callback",
	}
}
