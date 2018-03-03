package app

import (
	"encoding/json"
	"errors"
	"github.com/labstack/echo"
	"github.com/satori/go.uuid"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
	"io/ioutil"
	"net/http"
	"os"
)

var cookieNameState = "STATE"

// Google の OAuth 認証画面へリダイレクトさせるためのハンドラ
func OauthStartHandler(e echo.Context) error {
	// CSRF 対策にランダムな文字列を付与してコールバックの際に検証する
	// ハンドラ間での値の引き回しは Cookie を利用する
	state := uuid.NewV4().String()
	e.SetCookie(&http.Cookie{
		Name:  cookieNameState,
		Value: state,
		Path:  "/",
	})

	c := createOauth2Config()
	url := c.AuthCodeURL(state, oauth2.AccessTypeOnline)

	http.Redirect(e.Response(), e.Request(), url, 302)
	return nil
}

// Google の OAuth 認証が成功した場合にリダイレクトされてくるハンドラ
// 認証情報を使ってユーザーの参照または登録を行う
func OauthCallbackHandler(e echo.Context) error {
	ctx := appengine.NewContext(e.Request())

	// state が同一であるかチェック
	state := e.QueryParam("state")
	if cookie, err := e.Cookie(cookieNameState); err != nil {
		if cookie.Value != state {
			return errors.New("state is not valid")
		}
	}

	// 認証コードを使ってアクセストークンを取得する
	c := createOauth2Config()
	code := e.QueryParam("code")
	tok, err := c.Exchange(ctx, code)
	if err != nil {
		panic(err)
	}
	log.Debugf(ctx, "token: %v", tok)

	// アクセストークンを使って Google のユーザー情報を取得する
	client := c.Client(ctx, tok)
	ret, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		return err
	}
	defer ret.Body.Close()
	data, err := ioutil.ReadAll(ret.Body)
	if err != nil {
		return err
	}
	log.Debugf(ctx, "userinfo: %v", string(data))
	userinfo := struct {
		Sub string `json:"sub"`
	}{}
	if err := json.Unmarshal(data, &userinfo); err != nil {
		return err
	}
	log.Debugf(ctx, "sub: %v", userinfo.Sub)

	// 取得した Google のユーザー情報でアプリケーションのユーザーがいなければ作成する
	key := datastore.NewKey(ctx, "User", userinfo.Sub, 0, nil)
	u := &User{ID: userinfo.Sub}
	err = datastore.RunInTransaction(ctx, func(ctx context.Context) error {
		err := datastore.Get(ctx, key, u)
		if err != nil && err != datastore.ErrNoSuchEntity {
			log.Debugf(ctx, "user exists: %v", u)
			return err
		}
		_, err = datastore.Put(ctx, key, u)
		return err
	}, nil)
	if err != nil {
		log.Errorf(ctx, "Transaction failed: %v", err)
		return err
	}
	log.Debugf(ctx, "user created: %v", u)

	return e.JSON(http.StatusOK, struct{ Message string }{"ok"})
}

func createOauth2Config() oauth2.Config {
	return oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		Scopes:       []string{"openid", "profile"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  google.Endpoint.AuthURL,
			TokenURL: google.Endpoint.TokenURL,
		},
		RedirectURL: "http://localhost:8080/oauth/callback",
	}
}
