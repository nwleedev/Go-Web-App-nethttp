package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/quavious/golang-net-http-server/pkg/forms"
	"github.com/quavious/golang-net-http-server/pkg/models"
)

func (app *App) home(rw http.ResponseWriter, r *http.Request) {
	snippets, err := app.snippets.Latest()
	if err != nil {
		app.serverError(rw, err)
		return
	}

	app.render(rw, r, "home.page.html", &templateData{
		Snippets: snippets,
	})
}

func (app *App) showSnippet(rw http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil || id < 1 {
		app.notFound(rw)
		return
	}
	snippet, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(rw)
		} else {
			app.serverError(rw, err)
		}
		return
	}

	app.render(rw, r, "show.page.html", &templateData{
		Snippet: snippet,
	})
}

func (app *App) createSnippetForm(rw http.ResponseWriter, r *http.Request) {
	app.render(rw, r, "create.page.html", &templateData{
		Form: forms.New(nil),
	})
}

func (app *App) createSnippet(rw http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(rw, http.StatusBadRequest)
		return
	}

	form := forms.New(r.PostForm)
	form.Required("title", "content", "expires")
	form.MaxLength("title", 100)
	form.PermittedValues("expires", "365", "7", "1")

	if !form.Valid() {
		app.render(rw, r, "create.page.html", &templateData{
			Form: form,
		})
		return
	}

	id, err := app.snippets.Insert(form.Get("title"), form.Get("content"), form.Get("expires"))
	if err != nil {
		app.serverError(rw, err)
		return
	}

	app.session.Put(r, "flash", "Snippet successfully created!")

	http.Redirect(rw, r, fmt.Sprintf("/snippet/%d", id), http.StatusSeeOther)
}

func (app *App) signupUserForm(rw http.ResponseWriter, r *http.Request) {
	app.render(rw, r, "signup.page.html", &templateData{
		Form: forms.New(nil),
	})
}

func (app *App) signupUser(rw http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(rw, http.StatusBadRequest)
		return
	}
	form := forms.New(r.PostForm)
	form.Required("name", "email", "password")
	form.MaxLength("name", 255)
	form.MaxLength("email", 255)
	form.MatchesPattern("email", forms.EmailRegExp)
	form.MinLength("password", 10)

	if !form.Valid() {
		app.render(rw, r, "signup.page.html", &templateData{
			Form: form,
		})
	}

	err = app.users.Insert(form.Get("name"), form.Get("email"), form.Get("password"))
	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			form.Errors.Add("email", "Email Address is already in use.")
			app.render(rw, r, "signup.page.html", &templateData{
				Form: form,
			})
		} else {
			app.serverError(rw, err)
		}
	}
	app.session.Put(r, "flash", "Your signup is successful. Please log in.")
	http.Redirect(rw, r, "/user/login", http.StatusSeeOther)
}

func (app *App) loginUserForm(rw http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(rw, "")
}

func (app *App) loginUser(rw http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(rw, "")
}

func (app *App) logoutUser(rw http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(rw, "")
}
