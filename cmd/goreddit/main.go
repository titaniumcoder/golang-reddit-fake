package main

import (
	"log"
	"net/http"

	"github.com/titaniumcoder/golang-reddit-fake/postgres"
	"github.com/titaniumcoder/golang-reddit-fake/web"
)

func main() {
	store, err := postgres.NewStore("postgres://postgres:secret@localhost/postgres?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	h := web.NewHandler(store)
	http.ListenAndServe(":3000", h)
}
