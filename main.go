package main

import (
	"log"
	"net/http"

	"github.com/sedind/flow/app"
)

func main() {
	ctx := app.New("config.yaml")
	router := newAppRouter(ctx)

	log.Fatal(http.ListenAndServe(":8080", router))
}
