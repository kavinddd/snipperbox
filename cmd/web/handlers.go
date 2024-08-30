package main

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"snippetbox.kavinddd.net/internal/models"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.RequestURI != "/" {
		app.notFound(w)
		return
	}

	snippets, err := app.snippets.Latest()

	if err != nil {
		app.serverError(w, err)
		return
	}

	for _, snippet := range snippets {
		fmt.Fprintf(w, "%+v\n", snippet)
	}

	files := []string{
		"./ui/html/base.html",
		"./ui/html/partials/nav.html",
		"./ui/html/pages/home.html",
	}

	templates, err := template.ParseFiles(files...)

	if err != nil {
		app.serverError(w, err)
		return
	}

	fmt.Println(templates)

	err = templates.ExecuteTemplate(w, "base", nil)
	if err != nil {
		app.serverError(w, err)
		return
	}
}

func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {

	idString := r.URL.Query().Get("id")

	id, err := strconv.Atoi(idString)

	fmt.Println(id)

	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	snippet, err := app.snippets.Get(id)

	if err != nil {

		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
			return
		}

		app.serverError(w, err)
		return
	}

	fmt.Fprint(w, snippet)
}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", "POST")
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	content := "O snail\nClimb Mount Fuji, \nBut slowly slowly!\n\n= Kobashi"
	expires := 7

	id, err := app.snippets.Insert("Test", content, expires)

	if err != nil {
		app.serverError(w, err)
		return
	}

	app.infoLog.Printf("snippet %d is just created", id)
	http.Redirect(w, r, fmt.Sprintf("/snippet/view?id=%d", id), http.StatusTemporaryRedirect)
}
