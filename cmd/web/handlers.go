package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/julienschmidt/httprouter"
	"snippetbox.kavinddd.net/internal/models"
)

// a struct represeneting snippet create form
type sniperCreateForm struct {
	Title       string
	Content     string
	Expires     int
	FieldErrors map[string]string
}

func (app *application) newTemplateData() *templateData {
	return &templateData{
		CurrentYear: time.Now().Year(),
	}
}

func (app *application) home(w http.ResponseWriter, r *http.Request) {

	// this is required using vanilla because if server get requets /someRandomShit
	// and we have no handler for that, it will use this handler
	// but using httprouter, it matches the path / exactly, so only path "/" can use this handler
	// if r.RequestURI != "/" {
	// 	app.notFound(w)
	// 	return
	// }

	snippets, err := app.snippets.Latest()

	if err != nil {
		app.serverError(w, err)
		return
	}

	data := app.newTemplateData()
	data.Snippets = snippets

	app.render(w, http.StatusOK, "home.html", data)

}

func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {

	// vanilla version
	// view?id=10
	// idString := r.URL.Query().Get("id")

	// httprouter version
	// the path is changed, since we can use named parameter, not query parameter
	// view/:id
	params := httprouter.ParamsFromContext(r.Context())
	idString := params.ByName("id")

	id, err := strconv.Atoi(idString)

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

	data := app.newTemplateData()
	data.Snippet = snippet

	app.render(w, http.StatusOK, "view.html", data)

}

func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	// vanilla version, we can get rid of this since we uses method-based router from httprouter
	// if r.Method != http.MethodPost {
	// 	w.Header().Set("Allow", "POST")
	// 	app.clientError(w, http.StatusMethodNotAllowed)
	// 	return
	// }

	// by default, r.ParseForm() will return an error if the size of body is reaching 10 MB
	// you can override the maximum size usiing http.MaxBytesReader(w, r.Body, 4096) bytes (4MB)
	// http.MaxBytesReader(w, r.Body, 4096)

	if err := r.ParseForm(); err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	expires := r.PostForm.Get("expires")
	expiresInt, err := strconv.Atoi(expires)

	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := sniperCreateForm{
		Title:       r.PostForm.Get("title"),
		Content:     r.PostForm.Get("content"),
		Expires:     expiresInt,
		FieldErrors: map[string]string{},
	}

	// form validation - start

	if strings.TrimSpace(form.Title) == "" {
		form.FieldErrors["title"] = "This field cannot be blank"
	}

	if utf8.RuneCountInString(form.Title) > 100 {
		form.FieldErrors["title"] = "This field cannot be more than 100 characters long"
	}

	if strings.TrimSpace(form.Content) == "" {
		form.FieldErrors["content"] = "This field cannot be blank"
	}

	if expiresInt != 1 && expiresInt != 7 && expiresInt != 365 {
		form.FieldErrors["expires"] = "This field must equal to 1, 7 or 365"
	}

	if len(form.FieldErrors) != 0 {
		data := app.newTemplateData()
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "create.html", data)
		return
	}

	// form validation - end

	id, err := app.snippets.Insert(form.Title, form.Content, expiresInt)

	if err != nil {
		app.serverError(w, err)
		return
	}

	app.infoLog.Printf("snippet %d is just created", id)
	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {

	form := sniperCreateForm{
		Expires: 365,
	}

	data := app.newTemplateData()
	data.Form = form

	app.render(w, http.StatusOK, "create.html", data)

}
