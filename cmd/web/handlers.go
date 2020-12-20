package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

func (app *application) home(writer http.ResponseWriter, request *http.Request) {
	if request.URL.Path != "/" {
		app.notFound(writer)
		return
	}
	files := []string{
		"./ui/html/home.page.tmpl",
		"./ui/html/base.layout.tmpl",
		"./ui/html/footer.partial.tmpl",
	}
	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(writer, err)
		return
	}
	err = ts.Execute(writer, nil)
	if err != nil {
		app.serverError(writer, err)
		http.Error(writer, "Internal Server Error", 500)
	}
}

func (app *application) showSnippet(writer http.ResponseWriter, request *http.Request) {
	id, err := strconv.Atoi(request.URL.Query().Get("id"))
	if err != nil {
		app.notFound(writer)
		return
	}

	fmt.Fprintf(writer, "Display a specific snippet with ID %d...", id)
}

func (app *application) createSnippet(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		writer.Header().Set("Allow", http.MethodPost)
		app.clientError(writer, http.StatusMethodNotAllowed)
		return
	}
}
