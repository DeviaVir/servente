package main

import (
	"net/http"

	"github.com/bmizerany/pat"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders)
	dynamicMiddleware := alice.New(app.session.Enable, noSurf, app.authenticate)

	mux := pat.New()
	mux.Get("/", dynamicMiddleware.ThenFunc(app.home))
	mux.Get("/about", dynamicMiddleware.ThenFunc(app.about))

	// Service management
	mux.Get("/services", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.serviceHome))
	mux.Get("/services/you", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.serviceHomeYou))
	mux.Get("/service/new", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.serviceNewForm))
	mux.Post("/service/new", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.serviceNew))
	mux.Get("/service/:id", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.serviceShow))
	mux.Get("/service/edit/:id", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.serviceEditForm))
	mux.Post("/service/edit/:id", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.serviceEdit))

	// Organization settings
	mux.Get("/organization/start", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.organizationStart))
	mux.Get("/organization/new", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.organizationNewForm))
	mux.Post("/organization/new", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.organizationNew))
	mux.Get("/organization/selector", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.organizationSelectorForm))
	mux.Post("/organization/selector", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.organizationSelector))
	mux.Get("/organization/:id", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.organizationsHomeForm))
	mux.Post("/organization/:id", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.organizationsHome))
	mux.Get("/organization/:id/users", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.organizationsUsersForm))
	mux.Post("/organization/:id/users", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.organizationsUsers))

	// User auth
	mux.Get("/user/signup", dynamicMiddleware.ThenFunc(app.userSignupForm))
	mux.Post("/user/signup", dynamicMiddleware.ThenFunc(app.userSignup))
	mux.Get("/user/login", dynamicMiddleware.ThenFunc(app.userLoginForm))
	mux.Post("/user/login", dynamicMiddleware.ThenFunc(app.userLogin))
	mux.Post("/user/logout", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.userLogout))
	mux.Get("/user/profile", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.userProfile))
	mux.Get("/user/change-password", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.userChangePasswordForm))
	mux.Post("/user/change-password", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.userChangePassword))

	mux.Get("/ping", http.HandlerFunc(ping))

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Get("/static/", http.StripPrefix("/static", fileServer))

	return standardMiddleware.Then(mux)
}
