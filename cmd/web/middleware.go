package main

import (
	"fmt"
	"github.com/justinas/nosurf"
	"net/http"
)

func secureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("X-XSS-Protection", "1; mode=block")
		writer.Header().Set("X-Frame-Options", "deny")
		next.ServeHTTP(writer, request)
	})
}

func (app *application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		app.infoLog.Printf("%s - %s %s %s", request.RemoteAddr, request.Proto, request.Method, request.URL.RequestURI())
		next.ServeHTTP(writer, request)
	})
}

func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				writer.Header().Set("Connection", "close")
				app.serverError(writer, fmt.Errorf("%s", err))
			}
		}()
		next.ServeHTTP(writer, request)
	})
}

func (app *application) requireAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func (writer http.ResponseWriter, request *http.Request) {
		if !app.isAuthenticated(request) {
			http.Redirect(writer, request, "/user/login", http.StatusSeeOther)
			return
		}
		writer.Header().Add("Cache-Control", "no-store")
		next.ServeHTTP(writer, request)
	})
}

func (app *application) noSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)

	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path: "/",
		Secure: true,
	})

	return csrfHandler
}