package main

import (
    "net/http"
    "fmt"
)

func init() {
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprint(w, "hello world")
    })
}
