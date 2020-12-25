package main

import (
	"errors"
	"fmt"
	"github.com/shortnd/snippetbox/pkg/forms"
	"github.com/shortnd/snippetbox/pkg/models"
	"net/http"
	"strconv"
)

func (app *application) home(writer http.ResponseWriter, request *http.Request) {
	s, err := app.snippets.Latest()
	if err != nil {
		app.serverError(writer, err)
		return
	}

	app.render(writer, request, "home.page.tmpl", &templateData{
		Snippets: s,
	})
}

func (app *application) all(writer http.ResponseWriter, request *http.Request) {
	query := request.URL.Query()
	page, _ := strconv.Atoi(query.Get("page"))
	//fmt.Printf("%d", page)
	if page == 0 {
		page = 1
	}
	s, err := app.snippets.All(page)
	if err != nil {
		app.serverError(writer, err)
		return
	}

	app.render(writer, request, "all.page.tmpl", &templateData{
		Snippets: s,
	})
}

func (app *application) showSnippet(writer http.ResponseWriter, request *http.Request) {
	id, err := strconv.Atoi(request.URL.Query().Get(":id"))
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

	app.render(writer, request, "show.page.tmpl", &templateData{
		Snippet: s,
	})
}

func (app *application) createSnippetForm(writer http.ResponseWriter, request *http.Request) {
	app.render(writer, request, "create.page.tmpl", &templateData{
		Form: forms.New(nil),
	})
}

func (app *application) createSnippet(writer http.ResponseWriter, request *http.Request) {
	err := request.ParseForm()
	if err != nil {
		app.clientError(writer, http.StatusBadRequest)
		return
	}
	form := forms.New(request.PostForm)
	form.Required("title", "content", "expires")
	form.MaxLength("title", 100)
	form.PermittedValues("expires", "365", "7", "1")

	if !form.Valid() {
		app.render(writer, request, "create.page.tmpl", &templateData{
			Form: form,
		})
		return
	}

	id, err := app.snippets.Insert(form.Get("title"), form.Get("content"), form.Get("expires"))
	if err != nil {
		app.serverError(writer, err)
		return
	}

	app.session.Put(request, "flash", "Snippet successfully created!")
	http.Redirect(writer, request, fmt.Sprintf("/snippet/%d", id), http.StatusSeeOther)
}

func (app *application) signupUserForm(writer http.ResponseWriter, request *http.Request) {
	app.render(writer, request, "signup.page.tmpl", &templateData{
		Form: forms.New(nil),
	})
}

func (app *application) signupUser(writer http.ResponseWriter, request *http.Request) {
	err := request.ParseForm()
	if err != nil {
		app.clientError(writer, http.StatusBadRequest)
		return
	}
	form := forms.New(request.PostForm)
	form.Required("name", "email", "password")
	form.MaxLength("name", 255)
	form.MaxLength("email", 255)
	form.MatchPattern("email", forms.EmailRX)
	form.MinLength("password", 10)

	if !form.Valid() {
		app.render(writer, request, "signup.page.tmpl", &templateData{
			Form: form,
		})
		return
	}

	err = app.users.Insert(form.Get("name"), form.Get("email"), form.Get("password"))
	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			form.Errors.Add("email", "Address is already in use")
			app.render(writer, request, "signup.page.tmpl", &templateData{
				Form: form,
			})
		} else {
			app.serverError(writer, err)
		}
		return
	}
	app.session.Put(request, "flash", "Your signup was successful. Please log in.")
	http.Redirect(writer, request, "/user/login", http.StatusSeeOther)
}

func (app *application) loginUserForm(writer http.ResponseWriter, request *http.Request) {
	app.render(writer, request, "login.page.tmpl", &templateData{
		Form: forms.New(nil),
	})
}

func (app *application) loginUser(writer http.ResponseWriter, request *http.Request) {
	err := request.ParseForm()
	if err != nil {
		app.clientError(writer, http.StatusBadRequest)
		return
	}
	form := forms.New(request.PostForm)
	id, err := app.users.Authenticate(form.Get("email"), form.Get("password"))
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			form.Errors.Add("generic", "Email or Password is incorrect")
			app.render(writer, request, "login.page.tmpl", &templateData{
				Form: form,
			})
			return
		} else {
			app.serverError(writer, err)
			return
		}
	}
	app.session.Put(request, "authenticatedUserID", id)
	http.Redirect(writer, request, "/snippet/create", http.StatusSeeOther)
}

func (app *application) logoutUser(writer http.ResponseWriter, request *http.Request) {
	app.session.Remove(request, "authenticatedUserID")
	app.session.Put(request, "flash", "You've been logged out successfully!")
	http.Redirect(writer, request, "/", http.StatusSeeOther)
}
