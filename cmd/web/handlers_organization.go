package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/DeviaVir/servente/pkg/forms"
	"github.com/DeviaVir/servente/pkg/models"
)

func (app *application) organizationsHomeForm(w http.ResponseWriter, r *http.Request) {
	orgID := r.URL.Query().Get(":id")
	if orgID == "" {
		app.notFound(w)
		return
	}

	o, err := app.organizations.Get(orgID)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	// get existing data, if applicable
	existingSettings, err := app.organizations.GetSettings(o)
	if err != nil {
		app.serverError(w, err)
		return
	}
	o.Settings = existingSettings
	existingAttributes, err := app.organizations.GetAttributes(o)
	if err != nil {
		app.serverError(w, err)
		return
	}
	o.OrganizationAttributes = existingAttributes

	app.render(w, r, "organization/home.page.tmpl", &templateData{
		Organization: o,
		Form:         forms.New(nil),
	})
}

func (app *application) organizationsHome(w http.ResponseWriter, r *http.Request) {
	userID := app.session.GetInt(r, "authenticatedUserID")
	user, err := app.users.Get(userID)
	if err != nil {
		app.serverError(w, err)
		return
	}

	orgID := r.URL.Query().Get(":id")
	if orgID == "" {
		app.notFound(w)
		return
	}
	o, err := app.organizations.Get(orgID)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	if err := r.ParseForm(); err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := forms.New(r.PostForm)
	form.Required("name", "identifier")
	form.MaxLength("identifier", 100)

	if !form.Valid() {
		app.render(w, r, "organization/home.page.tmpl", &templateData{
			Organization: o,
			Form:         form,
		})
		return
	}

	org, err := app.organizations.Update(
		user,
		o,
		form.Get("name"),
	)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// store organization settings, attributes carry references so they will have to follow after
	// initially saving.
	settingsIdentifiers := r.PostForm["settings[identifier][]"]
	settingsTypes := r.PostForm["settings[type][]"]
	settingsValues := r.PostForm["settings[value][]"]
	for i, identifier := range settingsIdentifiers {
		setting := &models.Setting{
			Key:            identifier,
			Title:          identifier,
			Type:           settingsTypes[i],
			Scope:          "organization",
			OrganizationID: org.ID,
		}

		_, err := app.organizations.UpdateSetting(setting)
		if err != nil {
			app.serverError(w, err)
			return
		}
	}
	// store service attribute settings
	attributesIdentifiers := r.PostForm["attributes[identifier][]"]
	attributesTypes := r.PostForm["attributes[type][]"]
	for i, identifier := range attributesIdentifiers {
		setting := &models.Setting{
			Key:            identifier,
			Title:          identifier,
			Type:           attributesTypes[i],
			Scope:          "service",
			OrganizationID: org.ID,
		}

		_, err := app.organizations.UpdateSetting(setting)
		if err != nil {
			app.serverError(w, err)
			return
		}
	}

	// get existing settings now that they have been saved
	existingSettings, err := app.organizations.GetSettings(org)
	if err != nil {
		app.serverError(w, err)
		return
	}
	org.Settings = existingSettings

	// store organization attribute values connected to the
	// organization settings we created earlier
	for i, identifier := range settingsIdentifiers {
		for _, setting := range existingSettings {
			if setting.Key == identifier && setting.Scope == "organization" {
				_, err = app.organizations.UpdateAttribute(setting, settingsValues[i])
				if err != nil {
					app.serverError(w, err)
					return
				}
			}
		}
	}

	// get existing attributes now that they have been saved
	existingAttributes, err := app.organizations.GetAttributes(o)
	if err != nil {
		app.serverError(w, err)
		return
	}
	org.OrganizationAttributes = existingAttributes

	app.session.Put(r, "flash", "Organization successfully updated!")

	app.render(w, r, "organization/home.page.tmpl", &templateData{
		Organization: org,
		Form:         forms.New(nil),
	})
}

func (app *application) organizationStart(w http.ResponseWriter, r *http.Request) {
	userID := app.session.GetInt(r, "authenticatedUserID")
	user, err := app.users.Get(userID)
	if err != nil {
		app.serverError(w, err)
		return
	}

	organizations, err := app.users.Organizations(user)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// user not yet part of an organization, redirect to creating a new organization
	if len(organizations) < 1 {
		app.session.Put(r, "flash", "You are not currently part of any organization, please create a new one or request to be invited.")
		http.Redirect(w, r, "/organization/new", http.StatusSeeOther)
		return
	}

	selectedOrganizationID := app.session.GetString(r, "selectedOrganizationID")
	// user part of an organization, but an organization is not yet selected, redirect to selecting an organization
	if selectedOrganizationID == "" {
		app.session.Put(r, "flash", "No organization selected, please select an existing one or create a new organization.")
		http.Redirect(w, r, "/organization/selector", http.StatusSeeOther)
		return
	}

	// show the organization
	http.Redirect(w, r, fmt.Sprintf("/organization/%s", selectedOrganizationID), http.StatusSeeOther)
	return
}

func (app *application) organizationNewForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "organization/new.page.tmpl", &templateData{Form: forms.New(nil)})
}

func (app *application) organizationNew(w http.ResponseWriter, r *http.Request) {
	userID := app.session.GetInt(r, "authenticatedUserID")
	user, err := app.users.Get(userID)
	if err != nil {
		app.serverError(w, err)
		return
	}

	if err := r.ParseForm(); err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := forms.New(r.PostForm)
	form.Required("name", "identifier")
	form.MaxLength("identifier", 100)

	if !form.Valid() {
		app.render(w, r, "organization/new.page.tmpl", &templateData{Form: form})
		return
	}

	org, err := app.organizations.Insert(
		user,
		form.Get("identifier"),
		form.Get("name"),
	)
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.session.Put(r, "flash", "Organization successfully created!")
	app.session.Put(r, "selectedOrganizationID", org.Identifier)

	http.Redirect(w, r, fmt.Sprintf("/organization/%s", org.Identifier), http.StatusSeeOther)
}

func (app *application) organizationSelectorForm(w http.ResponseWriter, r *http.Request) {
	userID := app.session.GetInt(r, "authenticatedUserID")
	user, err := app.users.Get(userID)
	if err != nil {
		app.serverError(w, err)
		return
	}

	organizations, err := app.users.Organizations(user)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// user not yet part of an organization, redirect to creating a new organization
	if len(organizations) < 1 {
		app.session.Put(r, "flash", "You are not currently part of any organization, please create a new one or request to be invited.")
		http.Redirect(w, r, "/organization/new", http.StatusSeeOther)
		return
	}

	app.render(w, r, "organization/selector.page.tmpl", &templateData{
		Organizations: organizations,
		Form:          forms.New(nil),
	})
}

func (app *application) organizationSelector(w http.ResponseWriter, r *http.Request) {
	userID := app.session.GetInt(r, "authenticatedUserID")
	user, err := app.users.Get(userID)
	if err != nil {
		app.serverError(w, err)
		return
	}
	organizations, err := app.users.Organizations(user)
	if err != nil {
		app.serverError(w, err)
		return
	}
	if len(organizations) < 1 {
		app.session.Put(r, "flash", "You are not currently part of any organization, please create a new one or request to be invited.")
		http.Redirect(w, r, "/organization/new", http.StatusSeeOther)
		return
	}

	if err := r.ParseForm(); err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	var orgIDs []string
	for _, org := range organizations {
		orgIDs = append(orgIDs, org.Identifier)
	}

	form := forms.New(r.PostForm)
	form.Required("identifier")
	form.MaxLength("identifier", 100)
	form.PermittedValues("identifier", orgIDs...)

	if !form.Valid() {
		app.render(w, r, "organization/selector.page.tmpl", &templateData{
			Organizations: organizations,
			Form:          form,
		})
		return
	}

	app.session.Put(r, "selectedOrganizationID", form.Get("identifier"))

	http.Redirect(w, r, fmt.Sprintf("/organization/%s", form.Get("identifier")), http.StatusSeeOther)
}
