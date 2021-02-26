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
	org, _, err := app.addData(w, r)
	if err != nil {
		if errors.Is(err, models.ErrNoOrg) {
			app.session.Put(r, "flash", fmt.Sprintf("%s", err))
			http.Redirect(w, r, "/organization/selector", http.StatusSeeOther)
			return
		}
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	s, err := app.services.Latest(org, 0, 10)
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.render(w, r, "service/services.page.tmpl", &templateData{
		Services:     s,
		Organization: org,
	})
}

func (app *application) serviceHomeYou(w http.ResponseWriter, r *http.Request) {
	org, _, err := app.addData(w, r)
	if err != nil {
		if errors.Is(err, models.ErrNoOrg) {
			app.session.Put(r, "flash", fmt.Sprintf("%s", err))
			http.Redirect(w, r, "/organization/selector", http.StatusSeeOther)
			return
		}
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	s, err := app.services.Latest(org, 0, 10)
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.render(w, r, "service/services.page.tmpl", &templateData{
		Services:     s,
		Organization: org,
	})
}

func (app *application) serviceShow(w http.ResponseWriter, r *http.Request) {
	org, _, err := app.addData(w, r)
	if err != nil {
		if errors.Is(err, models.ErrNoOrg) {
			app.session.Put(r, "flash", fmt.Sprintf("%s", err))
			http.Redirect(w, r, "/organization/selector", http.StatusSeeOther)
			return
		}
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

	s, err := app.services.Get(org, id)
	fmt.Println(s)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	app.render(w, r, "service/show.page.tmpl", &templateData{
		Service:      s,
		Organization: org,
	})
}

func (app *application) serviceEditForm(w http.ResponseWriter, r *http.Request) {
	org, _, err := app.addData(w, r)
	if err != nil {
		if errors.Is(err, models.ErrNoOrg) {
			app.session.Put(r, "flash", fmt.Sprintf("%s", err))
			http.Redirect(w, r, "/organization/selector", http.StatusSeeOther)
			return
		}
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

	s, err := app.services.Get(org, id)
	fmt.Println(s)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	app.render(w, r, "service/edit.page.tmpl", &templateData{
		Form:         forms.New(nil),
		Organization: org,
		Service:      s,
	})
}

func (app *application) serviceEdit(w http.ResponseWriter, r *http.Request) {
	org, _, err := app.addData(w, r)
	if err != nil {
		if errors.Is(err, models.ErrNoOrg) {
			app.session.Put(r, "flash", fmt.Sprintf("%s", err))
			http.Redirect(w, r, "/organization/selector", http.StatusSeeOther)
			return
		}
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

	s, err := app.services.Get(org, id)
	fmt.Println(s)
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
	form.Required("title", "status")
	form.PermittedValues("status", "1", "2", "3", "4", "5", "6")

	if !form.Valid() {
		app.render(w, r, "service/edit.page.tmpl", &templateData{Form: form, Organization: org, Service: s})
		return
	}

	status, err := strconv.ParseInt(form.Get("status"), 10, 0)
	if err != nil {
		app.serverError(w, err)
		return
	}

	var serviceAttrs []*models.ServiceAttribute
	for _, setting := range org.Settings {
		if setting.Scope == "service" {
			attrs := r.PostForm[fmt.Sprintf("attributes[%s]", setting.Key)]
			if len(attrs) > 0 {
				existing := false
				for i, existingAttr := range s.ServiceAttributes {
					if existingAttr.Setting.ID == setting.ID {
						s.ServiceAttributes[i].Value = attrs[0]
						existing = true
					}
				}

				if !existing {
					serviceAttr := models.ServiceAttribute{
						Value:          attrs[0],
						Active:         true,
						SettingID:      setting.ID,
						OrganizationID: org.ID,
					}
					serviceAttrs = append(serviceAttrs, &serviceAttr)
				}
			}
		}
	}

	s.Title = form.Get("title")
	s.Description = form.Get("description")
	s.ServiceAttributes = serviceAttrs
	s.Status = int(status)

	id, err = app.services.Update(s)
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.session.Put(r, "flash", "Service successfully updated!")

	http.Redirect(w, r, fmt.Sprintf("/service/%d", id), http.StatusSeeOther)
}

func (app *application) serviceNewForm(w http.ResponseWriter, r *http.Request) {
	org, _, err := app.addData(w, r)
	if err != nil {
		if errors.Is(err, models.ErrNoOrg) {
			app.session.Put(r, "flash", fmt.Sprintf("%s", err))
			http.Redirect(w, r, "/organization/selector", http.StatusSeeOther)
			return
		}
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	app.render(w, r, "service/new.page.tmpl", &templateData{
		Form:         forms.New(nil),
		Organization: org,
	})
}

func (app *application) serviceNew(w http.ResponseWriter, r *http.Request) {
	org, _, err := app.addData(w, r)
	if err != nil {
		if errors.Is(err, models.ErrNoOrg) {
			app.session.Put(r, "flash", fmt.Sprintf("%s", err))
			http.Redirect(w, r, "/organization/selector", http.StatusSeeOther)
			return
		}
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
	form.Required("title", "identifier", "status")
	form.MaxLength("identifier", 100)
	form.PermittedValues("status", "1", "2", "3", "4", "5", "6")

	if !form.Valid() {
		app.render(w, r, "service/new.page.tmpl", &templateData{Form: form, Organization: org})
		return
	}

	status, err := strconv.ParseInt(form.Get("status"), 10, 0)
	if err != nil {
		app.serverError(w, err)
		return
	}

	var serviceAttrs []*models.ServiceAttribute
	for _, setting := range org.Settings {
		if setting.Scope == "service" {
			attrs := r.PostForm[fmt.Sprintf("attributes[%s]", setting.Key)]
			if len(attrs) > 0 {
				serviceAttr := models.ServiceAttribute{
					Value:          attrs[0],
					Active:         true,
					SettingID:      setting.ID,
					OrganizationID: org.ID,
				}
				serviceAttrs = append(serviceAttrs, &serviceAttr)
			}
		}
	}

	id, err := app.services.Insert(
		org,
		form.Get("identifier"),
		form.Get("title"),
		form.Get("description"),
		serviceAttrs,
		int(status),
	)
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.session.Put(r, "flash", "Service successfully created!")

	http.Redirect(w, r, fmt.Sprintf("/service/%d", id), http.StatusSeeOther)
}
