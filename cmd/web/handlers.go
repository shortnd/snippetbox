package main

import (
	"errors"
	"fmt"
	"github.com/shortnd/snippetbox/pkg/models"
	"net/http"
	"strconv"
)

func (app *application) home(writer http.ResponseWriter, request *http.Request) {
	if request.URL.Path != "/" {
		app.notFound(writer)
		return
	}
	s, err := app.snippets.Latest()
	if err != nil {
		app.serverError(writer, err)
		return
	}

	for _, snippet := range s {
		fmt.Fprintf(writer, "%v\n", snippet)
	}

	//return json.NewEncoder(writer).Encode(snippets)
	//files := []string{
	//	"./ui/html/home.page.tmpl",
	//	"./ui/html/base.layout.tmpl",
	//	"./ui/html/footer.partial.tmpl",
	//}
	//ts, err := template.ParseFiles(files...)
	//if err != nil {
	//	app.serverError(writer, err)
	//	return
	//}
	//err = ts.Execute(writer, nil)
	//if err != nil {
	//	app.serverError(writer, err)
	//}
}

func (app *application) showSnippet(writer http.ResponseWriter, request *http.Request) {
	id, err := strconv.Atoi(request.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.notFound(writer)
		return
	}

	s, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(writer)
		} else {
			app.serverError(writer, err)
		}
		return
	}
	fmt.Fprintf(writer, "%v", s)
}

func (app *application) createSnippet(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		writer.Header().Set("Allow", http.MethodPost)
		app.clientError(writer, http.StatusMethodNotAllowed)
		return
	}
	title := "0 snail"
	content := "0 snail\nClimb Mount Fuji,\nBut slowly, slowly!,\n\n-Kobayashi Issa"
	expires := "7"

	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(writer, err)
		return
	}

	http.Redirect(writer, request, fmt.Sprintf("/snippet?id=%d", id), http.StatusSeeOther)
}
