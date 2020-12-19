package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

func (app *application) home(writer http.ResponseWriter, request *http.Request)  {
	if request.URL.Path != "/" {
		http.NotFound(writer, request)
		return
	}
	files := []string{
		"./ui/html/home.page.bak",
		"./ui/html/base.layout.tmpl",
		"./ui/html/footer.partial.tmpl",
	}
	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.errorLog.Println(err.Error())
		http.Error(writer, "Internal Server Error", 500)
		return
	}
	err = ts.Execute(writer, nil)
	if err != nil {
		app.errorLog.Println(err.Error())
		http.Error(writer, "Internal Server Error", 500)
	}
}

func (app *application) showSnippet(writer http.ResponseWriter, request *http.Request)  {
	id, err := strconv.Atoi(request.URL.Query().Get("id"))
	if err != nil {
		http.NotFound(writer, request)
		return
	}

	fmt.Fprintf(writer, "Display a specific snippet with ID %d...", id)
}

func (app *application) createSnippet(writer http.ResponseWriter, request *http.Request)  {
	if request.Method != http.MethodPost {
		writer.Header().Set("Allow", http.MethodPost)
		http.Error(writer, "Method Not Allowed", 405)
		return
	}
}