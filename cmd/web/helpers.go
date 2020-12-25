package main

import (
	"bytes"
	"fmt"
	"github.com/justinas/nosurf"
	"net/http"
	"runtime/debug"
	"time"
)

// The serverError helper writes an error message and stack trace to the errorLog,
// then sends generic 500 Internal Server Error response to user
func (app *application) serverError(writer http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Output(2, trace)

	http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// The clientError helpers sends a specific status code and corresponding description
// to the user.
func (app *application) clientError(writer http.ResponseWriter, status int) {
	http.Error(writer, http.StatusText(status), status)
}

// For consistency, we'll also implement a notFound helper. This is simply a
// convenience wrapper around clientError which sends a 404 Not Found response to
// the user.
func (app *application) notFound(writer http.ResponseWriter) {
	app.clientError(writer, http.StatusNotFound)
}

func (app *application) render(writer http.ResponseWriter, request *http.Request, name string, td *templateData) {
	ts, ok := app.templateCache[name]
	if !ok {
		app.serverError(writer, fmt.Errorf("The template %s does not exist", name))
		return
	}
	buf := new(bytes.Buffer)

	err := ts.Execute(buf, app.addDefaultDate(td, request))
	if err != nil {
		app.serverError(writer, err)
		return
	}

	_, err = buf.WriteTo(writer)
	if err != nil {
		app.serverError(writer, err)
		return
	}
}

func (app *application) addDefaultDate(td *templateData, request *http.Request) *templateData {
	if td == nil {
		td = &templateData{}
	}
	td.CurrentYear = time.Now().Year()
	td.Flash = app.session.PopString(request, "flash")
	td.IsAuthenticated = app.isAuthenticated(request)
	td.CSRFToken = nosurf.Token(request)
	return td
}

func (app *application) isAuthenticated(request *http.Request) bool {
	return app.session.Exists(request, "authenticatedUserID")
}
