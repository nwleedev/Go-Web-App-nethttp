package main

import (
	"net/http"
	"path/filepath"

	"github.com/go-chi/chi/v5"
	"github.com/justinas/alice"
)

func (app *App) routes() http.Handler {
	standardMiddleware := alice.New(app.recoverPanic, app.logRequests, secureHeaders)
	dynamicMiddleware := alice.New(app.session.Enable, noSurf, app.authenticate)

	mux := chi.NewRouter()

	mux.Method("GET", "/", dynamicMiddleware.ThenFunc(app.home))
	mux.Method("GET", "/snippet/{id}", dynamicMiddleware.ThenFunc(app.showSnippet))
	mux.Method("GET", "/snippet/create", dynamicMiddleware.Append(app.requiredAuthentication).ThenFunc(app.createSnippetForm))
	mux.Method("POST", "/snippet/create", dynamicMiddleware.Append(app.requiredAuthentication).ThenFunc(app.createSnippet))

	mux.Method("GET", "/user/signup", dynamicMiddleware.ThenFunc(app.signupUserForm))
	mux.Method("POST", "/user/signup", dynamicMiddleware.ThenFunc(app.signupUser))
	mux.Method("GET", "/user/login", dynamicMiddleware.ThenFunc(app.loginUserForm))
	mux.Method("POST", "/user/login", dynamicMiddleware.ThenFunc(app.loginUser))
	mux.Method("POST", "/user/logout", dynamicMiddleware.Append(app.requiredAuthentication).ThenFunc(app.logoutUser))

	staticPath, _ := filepath.Abs("./ui/static")
	fileServer := http.FileServer(http.Dir(staticPath))
	mux.Handle("/static/*", http.StripPrefix("/static/", fileServer))

	return standardMiddleware.Then(mux)
}
