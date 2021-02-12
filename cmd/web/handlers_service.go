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
	selectedOrganizationID := app.session.GetString(r, "selectedOrganizationID")
	// user part of an organization, but an organization is not yet selected, redirect to selecting an organization
	if selectedOrganizationID == "" {
		app.session.Put(r, "flash", "No organization selected, please select an existing one or create a new organization.")
		http.Redirect(w, r, "/organization/selector", http.StatusSeeOther)
		return
	}
	org, err := app.organizations.Get(selectedOrganizationID)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	s, err := app.organizations.GetServices(org, 0, 10)
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.render(w, r, "service/services.page.tmpl", &templateData{
		Services: s,
	})
}

func (app *application) serviceHomeYou(w http.ResponseWriter, r *http.Request) {
	selectedOrganizationID := app.session.GetString(r, "selectedOrganizationID")
	// user part of an organization, but an organization is not yet selected, redirect to selecting an organization
	if selectedOrganizationID == "" {
		app.session.Put(r, "flash", "No organization selected, please select an existing one or create a new organization.")
		http.Redirect(w, r, "/organization/selector", http.StatusSeeOther)
		return
	}
	org, err := app.organizations.Get(selectedOrganizationID)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	// @TODO: create new endpoint to also select owning team
	s, err := app.organizations.GetServices(org, 0, 10)
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.render(w, r, "service/services.page.tmpl", &templateData{
		Services: s,
	})
}

func (app *application) serviceShow(w http.ResponseWriter, r *http.Request) {
	selectedOrganizationID := app.session.GetString(r, "selectedOrganizationID")
	// user part of an organization, but an organization is not yet selected, redirect to selecting an organization
	if selectedOrganizationID == "" {
		app.session.Put(r, "flash", "No organization selected, please select an existing one or create a new organization.")
		http.Redirect(w, r, "/organization/selector", http.StatusSeeOther)
		return
	}
	org, err := app.organizations.Get(selectedOrganizationID)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	id, err := strconv.Atoi(r.URL.Query().Get(":id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	s, err := app.organizations.GetService(org, id)
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
	userID := app.session.GetInt(r, "authenticatedUserID")
	_, err := app.users.Get(userID)
	if err != nil {
		app.serverError(w, err)
		return
	}

	selectedOrganizationID := app.session.GetString(r, "selectedOrganizationID")
	// user part of an organization, but an organization is not yet selected, redirect to selecting an organization
	if selectedOrganizationID == "" {
		app.session.Put(r, "flash", "No organization selected, please select an existing one or create a new organization.")
		http.Redirect(w, r, "/organization/selector", http.StatusSeeOther)
		return
	}
	org, err := app.organizations.Get(selectedOrganizationID)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}
	fmt.Println(org)

	// @TODO: get service attributes and show in the form

	app.render(w, r, "service/new.page.tmpl", &templateData{Form: forms.New(nil)})
}

func (app *application) serviceNew(w http.ResponseWriter, r *http.Request) {
	userID := app.session.GetInt(r, "authenticatedUserID")
	_, err := app.users.Get(userID)
	if err != nil {
		app.serverError(w, err)
		return
	}

	selectedOrganizationID := app.session.GetString(r, "selectedOrganizationID")
	// user part of an organization, but an organization is not yet selected, redirect to selecting an organization
	if selectedOrganizationID == "" {
		app.session.Put(r, "flash", "No organization selected, please select an existing one or create a new organization.")
		http.Redirect(w, r, "/organization/selector", http.StatusSeeOther)
		return
	}
	org, err := app.organizations.Get(selectedOrganizationID)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}
	fmt.Println(org)

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

	// @TODO: use organizations.UpdateService instead

	id, err := app.services.Insert(
		form.Get("identifier"),
		form.Get("title"),
		form.Get("description"),
		nil, // form.Get("attributes")
		int(status),
	)
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.session.Put(r, "flash", "Service successfully created!")

	http.Redirect(w, r, fmt.Sprintf("/service/%d", id), http.StatusSeeOther)
}
