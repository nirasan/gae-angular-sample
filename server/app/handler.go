package app

import (
	"github.com/labstack/echo"
	"net/http"
)


func NewHandler() http.Handler {
	e := echo.New()

	e.GET("/oauth/start", OauthStartHandler)
	e.GET("/oauth/callback", OauthCallbackHandler)

	api := e.Group("/api", AuthorizationMiddleware)
	api.GET("/hello", func(e echo.Context) error {
		return e.JSON(http.StatusOK, struct{ Message string }{"hello authorized"})
	})

	e.GET("/hello", func(e echo.Context) error {
		return e.JSON(http.StatusOK, struct{ Message string }{"hello not authorized"})
	})

	// todo resource
	api.GET("/todo/", TodoListHandler)
	api.POST("/todo/", TodoCreateHandler)
	api.PUT("/todo/", TodoUpdateHandler)
	api.DELETE("/todo/:id", TodoDeleteHandler)

	return e
}

