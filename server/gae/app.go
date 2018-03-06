package main

import (
	"github.com/nirasan/gae-angular-sample/server/app"
	"net/http"
)

func init() {
	http.Handle("/", app.NewHandler())
}
