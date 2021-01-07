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
	mux.Get("/services", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.serviceHome))
	mux.Get("/services/you", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.serviceHomeYou))
	mux.Get("/service/new", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.serviceNewForm))
	mux.Post("/service/new", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.serviceNew))
	mux.Get("/service/:id", dynamicMiddleware.ThenFunc(app.serviceShow))

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
