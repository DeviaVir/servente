package main

import (
	"bytes"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/justinas/nosurf"

	sdk "github.com/DeviaVir/servente-sdk"
	"github.com/DeviaVir/servente/pkg/models"
	"github.com/DeviaVir/servente/pkg/owners"
)

func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Output(2, trace) // don't report from this function but the originator

	if app.debug {
		http.Error(w, trace, http.StatusInternalServerError)
		return
	}

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

func (app *application) addDefaultData(td *templateData, r *http.Request) *templateData {
	if td == nil {
		td = &templateData{}
	}

	td.CSRFToken = nosurf.Token(r)
	td.CurrentYear = time.Now().Year()
	td.Flash = app.session.PopString(r, "flash")
	td.IsAuthenticated = app.isAuthenticated(r)
	td.SettingsTypes = models.SettingsTypes
	td.AttributesTypes = models.AttributesTypes
	return td
}

func (app *application) isAuthenticated(r *http.Request) bool {
	isAuthenticated, ok := r.Context().Value(contextKeyIsAuthenticated).(bool)
	if !ok {
		return false
	}
	return isAuthenticated
}

func (app *application) addData(w http.ResponseWriter, r *http.Request) (*models.Organization, *models.User, error) {
	// get user object
	userID := app.session.GetInt(r, "authenticatedUserID")
	user, err := app.users.Get(userID)
	if err != nil {
		return nil, nil, err
	}

	// get selected organization identifier
	selectedOrganizationID := app.session.GetString(r, "selectedOrganizationID")
	if selectedOrganizationID == "" {
		return nil, user, models.ErrNoOrg
	}

	// get user's organizations
	orgs, err := app.users.Organizations(user)
	if err != nil {
		return nil, user, err
	}
	// now find the organization ID we are looking for
	var org *models.Organization
	for _, o := range orgs {
		if o.Identifier == selectedOrganizationID {
			org, err = app.organizations.Get(o.Identifier)
			if err != nil {
				return nil, user, err
			}
		}
	}
	if org == nil {
		return nil, user, models.ErrNoOrg
	}

	// return!
	return org, user, nil
}

func (app *application) getOwnerTeams(org *models.Organization) ([]sdk.JSONTeam, error) {
	api, err := owners.Init(org)
	if err != nil {
		return nil, err // no endpoints discovered
	}

	return api.GetTeamsList()
}

func (app *application) userPartOfOwnerTeams(org *models.Organization, user *models.User) ([]sdk.JSONTeam, error) {
	api, err := owners.Init(org)
	if err != nil {
		return nil, err // no endpoints discovered
	}

	return api.UserPartOfTeams(user)
}

func (app *application) render(w http.ResponseWriter, r *http.Request, name string, td *templateData) {
	ts, ok := app.templateCache[name]
	if !ok {
		app.serverError(w, fmt.Errorf("The template %s does not exist", name))
		return
	}

	buf := new(bytes.Buffer)

	err := ts.Execute(buf, app.addDefaultData(td, r))
	if err != nil {
		app.serverError(w, err)
	}

	buf.WriteTo(w)
}
