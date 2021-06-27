package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/justinas/nosurf"
	"github.com/quavious/golang-net-http-server/pkg/models"
)

func (app *App) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				rw.Header().Set("Connection", "close")
				app.serverError(rw, fmt.Errorf("%s", err))
			}
		}()
		next.ServeHTTP(rw, r)
	})
}

func (app *App) logRequests(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		app.infoLog.Printf("%s - %s %s %s", r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI())
		next.ServeHTTP(rw, r)
	})
}

func secureHeaders(next http.Handler) http.Handler {
	fn := func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Set("X-XSS-Protection", "1; mode=block")
		rw.Header().Set("X-Frame-Options", "deny")
		next.ServeHTTP(rw, r)
	}
	return http.HandlerFunc(fn)
}

func (app *App) requiredAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		if app.isAuthenticated(r) {
			http.Redirect(rw, r, "/user/login", http.StatusSeeOther)
			return
		}
		rw.Header().Add("Cache-Control", "no-store")
		next.ServeHTTP(rw, r)
	})
}

func noSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)
	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   true,
	})
	return csrfHandler
}

func (app *App) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		exists := app.session.Exists(r, "authenticatedUserID")
		if !exists {
			next.ServeHTTP(rw, r)
			return
		}

		user, err := app.users.Get(app.session.GetInt(r, "authenticatedUserID"))
		if errors.Is(err, models.ErrNoRecord) || !user.Active {
			app.session.Remove(r, "authenticatedUserID")
			next.ServeHTTP(rw, r)
			return
		} else if err != nil {
			app.serverError(rw, err)
		}

		ctx := context.WithValue(r.Context(), contextKeyIsAuthenticated, true)
		next.ServeHTTP(rw, r.WithContext(ctx))
	})
}
