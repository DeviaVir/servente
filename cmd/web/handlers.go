package main

import (
	"net/http"
)

// ping: healthcheck endpoint
func ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "servente/home.page.tmpl", nil)
}

func (app *application) about(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "servente/about.page.tmpl", nil)
}
