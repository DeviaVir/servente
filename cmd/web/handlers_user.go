package main

import (
	"errors"
	"net/http"

	"github.com/DeviaVir/servente/pkg/forms"
	"github.com/DeviaVir/servente/pkg/models"
)

func (app *application) userSignupForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "user/signup.page.tmpl", &templateData{Form: forms.New(nil)})
}

func (app *application) userSignup(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := forms.New(r.PostForm)
	form.Required("name", "email", "password")
	form.MaxLength("name", 255)
	form.MaxLength("email", 255)
	form.MatchesPattern("email", forms.EmailRX)
	form.MinLength("password", 10)

	if !form.Valid() {
		app.render(w, r, "user/signup.page.tmpl", &templateData{Form: form})
		return
	}

	if err := app.users.Insert(form.Get("name"), form.Get("email"), form.Get("password")); err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			form.Errors.Add("email", "Address is already in use")
			app.render(w, r, "user/signup.page.tmpl", &templateData{Form: form})
		} else {
			app.serverError(w, err)
		}
		return
	}

	app.session.Put(r, "flash", "Your signup was successful. Please log in.")

	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

func (app *application) userLoginForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "user/login.page.tmpl", &templateData{Form: forms.New(nil)})
}

func (app *application) userLogin(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := forms.New(r.PostForm)
	id, err := app.users.Authenticate(form.Get("email"), form.Get("password"))
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			form.Errors.Add("generic", "Email or Password is incorrect")
			app.render(w, r, "user/login.page.tmpl", &templateData{Form: form})
		} else {
			app.serverError(w, err)
		}
		return
	}

	app.session.Put(r, "authenticatedUserID", id)

	// also store the organization
	if form.Get("organization") != "" {
		user, err := app.users.Get(id)
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

		var confirmedOrg *models.Organization
		for _, org := range organizations {
			if org.Identifier == form.Get("organization") {
				confirmedOrg = org
			}
		}

		if confirmedOrg == nil {
			app.session.Put(r, "flash", "That organization does not exist or has not invited you, please select a different one, create a new one or request to be invited.")
			http.Redirect(w, r, "/organization/selector", http.StatusSeeOther)
			return
		}

		app.session.Put(r, "selectedOrganizationID", confirmedOrg.Identifier)
	}

	path := app.session.PopString(r, "redirectPathAfterLogin")
	if path != "" {
		http.Redirect(w, r, path, http.StatusSeeOther)
		return
	}
	http.Redirect(w, r, "/services", http.StatusSeeOther)
}

func (app *application) userLogout(w http.ResponseWriter, r *http.Request) {
	app.session.Remove(r, "authenticatedUserID")

	app.session.Put(r, "flash", "You've been logged out successfully!")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) userProfile(w http.ResponseWriter, r *http.Request) {
	userID := app.session.GetInt(r, "authenticatedUserID")

	user, err := app.users.Get(userID)
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.render(w, r, "user/profile.page.tmpl", &templateData{User: user})
}

func (app *application) userChangePasswordForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "user/password.page.tmpl", &templateData{Form: forms.New(nil)})
}

func (app *application) userChangePassword(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := forms.New(r.PostForm)
	form.Required("currentPassword", "newPassword", "newPasswordConfirmation")
	form.MinLength("newPassword", 10)
	if form.Get("newPassword") != form.Get("newPasswordConfirmation") {
		form.Errors.Add("newPasswordConfirmation", "Passwords do not match")
	}

	if !form.Valid() {
		app.render(w, r, "user/password.page.tmpl", &templateData{Form: form})
		return
	}

	userID := app.session.GetInt(r, "authenticatedUserID")

	if err := app.users.ChangePassword(userID, form.Get("currentPassword"), form.Get("newPassword")); err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			form.Errors.Add("currentPassword", "Password is incorrect")
			app.render(w, r, "user/password.page.tmpl", &templateData{Form: form})
		} else if err != nil {
			app.serverError(w, err)
		}
		return
	}

	app.session.Put(r, "flash", "You've successfully changed your password!")
	http.Redirect(w, r, "/user/profile", http.StatusSeeOther)
}
