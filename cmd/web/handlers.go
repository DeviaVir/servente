package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/DeviaVir/servente/pkg/forms"
	"github.com/DeviaVir/servente/pkg/models"
)

// ping: healthcheck endpoint
func ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	s, err := app.services.Latest(10)
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.render(w, r, "service/home.page.tmpl", &templateData{
		Services: s,
	})
}

func (app *application) about(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "servente/about.page.tmpl", nil)
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
	form.Required("title", "content", "expires")
	form.MaxLength("title", 100)
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

	path := app.session.PopString(r, "redirectPathAfterLogin")
	if path != "" {
		http.Redirect(w, r, path, http.StatusSeeOther)
		return
	}
	http.Redirect(w, r, "/service/new", http.StatusSeeOther)
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
