package main

import (
	"github.com/labstack/echo"
	"net/http"
	"github.com/nirasan/gae-angular-sample/server/app"
)

func init() {
    // 軽量なウェブアプリケーションフレームワーク echo を使う
    // 素の net/http に比べてパラメータの bind や Json の出力を便利になる
	e := echo.New()

	// ルート定義
	e.GET("/hello", helloHandler)

	e.GET("/oauth/start", app.OauthStartHandler)
	e.GET("/oauth/callback", app.OauthCallbackHandler)

	// 全リクエストを echo で処理する
	http.Handle("/", e)
}

// helloHandler のリクエスト
type helloRequest struct {
    // クエリストリングの name 要素を受け取る
	Name string `query:"name"`
}

// helloHandler のレスポンス
type helloResponse struct {
	Message string
}

// echo のハンドラ型を定義
func helloHandler(c echo.Context) error {
    // 入力を受け取って構造体に入れる
	req := new(helloRequest)
	c.Bind(req)

	// ステータスコード 200 で JSON を返す
	return c.JSON(http.StatusOK, helloResponse{
		Message: "hello " + req.Name,
	})
}
