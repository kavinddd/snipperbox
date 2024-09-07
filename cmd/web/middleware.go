package main

import (
	"fmt"
	"net/http"
)

func secureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Security-Policy", "default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com")
		w.Header().Set("Referrer-Policy", "origin-when-cross-origin")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "deny")
		w.Header().Set("X-XSS-Protection", "0")

		// this is a separate line to execute the middleware logic
		// code above this line will be executed first
		next.ServeHTTP(w, r)
		// code under this line will be executed after the last middleware returns
		// notice that this middlware doesn't return
		// return can be useful when you want to end middleware (next.ServeHTTP)
		// for example, authentication and authorization
	})
}

func (app *application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.infoLog.Printf("%s - %s %s %s", r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI())
		next.ServeHTTP(w, r)
	})
}

func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// having this defer, so, it will trigger after other middleware return/exit
		// this middleware can be the first one with the help of defer
		defer func() {
			if recover := recover(); recover != nil {
				w.Header().Set("Connection", "close")

				app.serverError(w, fmt.Errorf("PANIC: %s", recover))
			}
		}()

		next.ServeHTTP(w, r)
	})
}
