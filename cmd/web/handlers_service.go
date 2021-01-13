package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/DeviaVir/servente/pkg/forms"
	"github.com/DeviaVir/servente/pkg/models"
)

func (app *application) serviceHome(w http.ResponseWriter, r *http.Request) {
	s, err := app.services.Latest(10)
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.render(w, r, "service/services.page.tmpl", &templateData{
		Services: s,
	})
}

func (app *application) serviceHomeYou(w http.ResponseWriter, r *http.Request) {
	s, err := app.services.Latest(10)
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.render(w, r, "service/services.page.tmpl", &templateData{
		Services: s,
	})
}

func (app *application) serviceShow(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get(":id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	s, err := app.services.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	app.render(w, r, "service/show.page.tmpl", &templateData{
		Service: s,
	})
}

func (app *application) serviceNewForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "service/new.page.tmpl", &templateData{Form: forms.New(nil)})
}

func (app *application) serviceNew(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := forms.New(r.PostForm)
	form.Required("title", "identifier", "status")
	form.MaxLength("identifier", 100)
	form.PermittedValues("status", "1", "2", "3", "4", "5", "6")

	if !form.Valid() {
		app.render(w, r, "service/new.page.tmpl", &templateData{Form: form})
		return
	}

	status, err := strconv.ParseInt(form.Get("status"), 10, 0)
	if err != nil {
		app.serverError(w, err)
		return
	}

	id, err := app.services.Insert(
		form.Get("identifier"),
		form.Get("title"),
		form.Get("description"),
		form.Get("attributes"),
		int(status),
	)
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.session.Put(r, "flash", "Service successfully created!")

	http.Redirect(w, r, fmt.Sprintf("/service/%d", id), http.StatusSeeOther)
}
