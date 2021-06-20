package main

import (
	"bytes"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"
)

func (app *App) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Println(trace)

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *App) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *App) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

func (app *App) addDefaultData(data *templateData, r *http.Request) *templateData {
	if data == nil {
		data = &templateData{}
	}
	data.CurrentYear = time.Now().Year()
	data.Flash = app.session.PopString(r, "flash")
	return data
}

func (app *App) render(rw http.ResponseWriter, r *http.Request, name string, data *templateData) {
	ts, ok := app.templateCache[name]
	if !ok {
		app.serverError(rw, fmt.Errorf("the template %s does not exist", name))
		return
	}
	buf := &bytes.Buffer{}
	err := ts.Execute(buf, app.addDefaultData(data, r))

	if err != nil {
		app.serverError(rw, err)
		return
	}

	buf.WriteTo(rw)
}
